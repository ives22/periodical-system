import router from '@/router'
import useUserStore from '@/store/modules/user'
import { ElMessage } from 'element-plus'
import pinia from './store'

const userStore = useUserStore(pinia)

router.beforeEach(async (to, from, next) => {
    // 1. 存在token
    if (userStore.access_token) {
        if (to.path === '/login') {
            next({ path: '/' })
            return
        }

        // 若尚未从服务端拉取最新角色与权限信息，则先同步
        if (!userStore.profileLoaded) {
            try {
                await userStore.fetchUserProfile()
            } catch (error) {
                userStore.clearUserInfo()
                next({ path: '/login' })
                return
            }
        } else {
            userStore.updateMenuRoutes()
        }

        // 校验路由级权限码
        const requiredPerm = to.meta?.perm as string | undefined
        if (requiredPerm && !userStore.hasPermission(requiredPerm)) {
            ElMessage.warning('您没有权限访问该页面')
            next({ path: '/home' })
            return
        }

        next()
    } else {
        // 2. 不存在token 未登录
        if (to.path === '/login') {
            next()
        } else {
            next({ path: '/login' })
        }
    }
})
