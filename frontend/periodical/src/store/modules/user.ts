//  用户相关的仓库
import { defineStore } from 'pinia'

// 引入路由（常量路由）
import { constantRoute } from '@/router/routers'
import { reqLogin, reqLogout } from '@/api/user'



// 创建用户小仓库
let useUserStore = defineStore('User', {
    // 小仓库存储数据的地方
    state: (): any => {
        return {
            access_token: localStorage.getItem('access_token') || '',
            refresh_token: localStorage.getItem('refresh_token') || '',
            menuRoutes: constantRoute,
            loginUser: {
                name: localStorage.getItem('username') || '',
                access_token: localStorage.getItem('access_token') || '',
                refresh_token: localStorage.getItem('refresh_token') || '',
                role: Number(localStorage.getItem('role')) || 0
            }
        }
    },

    // 异步｜逻辑的地方
    actions: {
        async userLogin(data: any) {
            const result = await reqLogin(data)
            if (result.code == 200) {
                // 存储到 state
                this.access_token = result.data.access_token
                this.refresh_token = result.data.refresh_token
                this.loginUser = {
                    name: result.data.username,
                    access_token: result.data.access_token,
                    refresh_token: result.data.refresh_token,
                    role: Number(result.data.role)  // 确保是数字类型
                }
                // 持久化存储
                localStorage.setItem('access_token', result.data.access_token)
                localStorage.setItem('refresh_token', result.data.refresh_token)
                localStorage.setItem('username', result.data.username)
                localStorage.setItem('role', result.data.role.toString())

                // console.log(result)
                // this.access_token = result.data.access_token
                // this.refresh_token = result.data.refresh_token
                // this.loginUser.name = result.data.username
                // this.loginUser.access_token = result.data.access_token
                // this.loginUser.refresh_token = result.data.refresh_token
                // this.loginUser.role = result.data.role

                // console.log(this.loginUser)
                return 'ok'
            } else {
                return Promise.reject(result.data)
            }
        },
        async userLogout() {
            try {
                // 构造退出登录需要的数据
                const data = {
                    access_token: this.access_token,
                    refresh_token: this.refresh_token
                }
                console.log("token", this.access_token)
                console.log("data:", data)
                const result = await reqLogout(data)
                if (result.code == 200) {
                    // 清空状态和本地存储
                    this.clearUserInfo()
                    // // 清空用户信息
                    // this.access_token = ''
                    // this.refresh_token = ''
                    // this.loginUser = {
                    //     name: '',
                    //     access_token: '',
                    //     refresh_token: '',
                    //     role: ''
                    // }
                    return 'ok'
                } else {
                    return Promise.reject(result.message || '退出登录失败')
                }
            } catch (error) {
                console.log("退出登录失败", error)
                return Promise.reject('退出登录失败')
            }
        },

        // 清空用户信息
        clearUserInfo() {
            this.access_token = ''
            this.refresh_token = ''
            this.loginUser = {
                name: '',
                access_token: '',
                refresh_token: '',
                role: ''
            }
            // 清空本地存储
            localStorage.removeItem('access_token')
            localStorage.removeItem('refresh_token')
            localStorage.removeItem('username')
            localStorage.removeItem('role')
        }
    },

    getters: {

    }


})

// 对外暴露获取小仓库的方法
export default useUserStore