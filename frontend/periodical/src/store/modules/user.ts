// 用户相关的仓库
import { defineStore } from 'pinia'
import { constantRoute } from '@/router/routers'
import { reqLogin, reqLogout, reqUserProfile } from '@/api/user'
import useTagsViewStore from './tagsView'

// 递归过滤菜单路由树
function filterMenuRoutes(routes: any[], hasPerm: (code?: string) => boolean): any[] {
    const res: any[] = []
    for (const route of routes) {
        const tmp = { ...route }
        const requiredPerm = tmp.meta?.perm
        if (!requiredPerm || hasPerm(requiredPerm)) {
            if (tmp.children && tmp.children.length > 0) {
                tmp.children = filterMenuRoutes(tmp.children, hasPerm)
                if (tmp.children.length > 0) {
                    res.push(tmp)
                }
            } else {
                res.push(tmp)
            }
        }
    }
    return res
}

let useUserStore = defineStore('User', {
    state: (): any => {
        const savedRoles = localStorage.getItem('roles')
        const savedPerms = localStorage.getItem('permissions')
        return {
            access_token: localStorage.getItem('access_token') || '',
            refresh_token: localStorage.getItem('refresh_token') || '',
            profileLoaded: false,
            menuRoutes: constantRoute,
            loginUser: {
                id: Number(localStorage.getItem('user_id')) || 0,
                name: localStorage.getItem('username') || '',
                access_token: localStorage.getItem('access_token') || '',
                refresh_token: localStorage.getItem('refresh_token') || '',
                role: Number(localStorage.getItem('role')) || 0,
                roles: savedRoles ? JSON.parse(savedRoles) : [],
                permissions: savedPerms ? JSON.parse(savedPerms) : [],
                is_super_admin: localStorage.getItem('is_super_admin') === 'true'
            }
        }
    },

    actions: {
        // 判断当前登录用户是否拥有指定权限码
        hasPermission(permCode?: string): boolean {
            if (!permCode) return true
            if (this.loginUser.is_super_admin) return true
            const perms: string[] = this.loginUser.permissions || []
            return perms.includes(permCode)
        },

        // 根据当前用户权限动态生成左侧菜单路由
        updateMenuRoutes() {
            this.menuRoutes = filterMenuRoutes(constantRoute, (code?: string) => this.hasPermission(code))
        },

        // 保存用户登录态与角色权限数据
        applyUserAuthData(data: any) {
            if (data.access_token) {
                this.access_token = data.access_token
                localStorage.setItem('access_token', data.access_token)
            }
            if (data.refresh_token) {
                this.refresh_token = data.refresh_token
                localStorage.setItem('refresh_token', data.refresh_token)
            }

            const roles = data.roles || []
            const permissions = data.permissions || []
            const isSuperAdmin = Boolean(data.is_super_admin) || roles.some((r: any) => r.code === 'super_admin')

            this.loginUser = {
                id: data.user_id || data.id || this.loginUser.id,
                name: data.username || this.loginUser.name,
                access_token: this.access_token,
                refresh_token: this.refresh_token,
                role: Number(data.role) || 0,
                roles,
                permissions,
                is_super_admin: isSuperAdmin
            }

            localStorage.setItem('user_id', String(this.loginUser.id || 0))
            localStorage.setItem('username', this.loginUser.name)
            localStorage.setItem('role', String(this.loginUser.role))
            localStorage.setItem('roles', JSON.stringify(roles))
            localStorage.setItem('permissions', JSON.stringify(permissions))
            localStorage.setItem('is_super_admin', String(isSuperAdmin))

            this.updateMenuRoutes()
        },

        async userLogin(data: any) {
            const result = await reqLogin(data)
            if (result.code == 200) {
                this.applyUserAuthData(result.data)
                this.profileLoaded = true
                return 'ok'
            } else {
                return Promise.reject(result.msg || result.data || '登录失败')
            }
        },

        // 从后端同步当前登录用户的最新角色与权限集
        async fetchUserProfile() {
            const result = await reqUserProfile()
            if (result.code == 200 && result.data) {
                this.applyUserAuthData(result.data)
                this.profileLoaded = true
                return result.data
            } else {
                return Promise.reject(result.msg || '获取用户信息失败')
            }
        },

        async userLogout() {
            try {
                const data = {
                    access_token: this.access_token,
                    refresh_token: this.refresh_token
                }
                if (this.access_token) {
                    await reqLogout(data)
                }
            } catch (error) {
                console.warn('调用退出登录接口异常:', error)
            } finally {
                this.clearUserInfo()
            }
            return 'ok'
        },

        // 清空用户信息
        clearUserInfo() {
            this.access_token = ''
            this.refresh_token = ''
            this.profileLoaded = false
            this.menuRoutes = constantRoute
            this.loginUser = {
                id: 0,
                name: '',
                access_token: '',
                refresh_token: '',
                role: 0,
                roles: [],
                permissions: [],
                is_super_admin: false
            }
            localStorage.removeItem('access_token')
            localStorage.removeItem('refresh_token')
            localStorage.removeItem('user_id')
            localStorage.removeItem('username')
            localStorage.removeItem('role')
            localStorage.removeItem('roles')
            localStorage.removeItem('permissions')
            localStorage.removeItem('is_super_admin')
            try {
                const tagsViewStore = useTagsViewStore()
                tagsViewStore.delAllViews()
            } catch {
                // ignore
            }
        }
    },

    getters: {
        roleNames(state): string {
            const roles = state.loginUser?.roles || []
            if (roles.length > 0) {
                return roles.map((r: any) => r.name).join(' / ')
            }
            return state.loginUser?.role === 1 ? '管理员' : '普通用户'
        }
    }
})

export default useUserStore
