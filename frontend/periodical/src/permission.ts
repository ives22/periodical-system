
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
        } else {
            next()
        }
    } else {
        // 2. 不存在token 未登录
        if (to.path === '/login') {
            next()
        } else {
            // next({ path: '/login', query: { redirect: to.path } })
            next({ path: '/login' })
        }
    }
})

