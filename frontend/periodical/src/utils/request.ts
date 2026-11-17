import axios from "axios"
import { ElMessage } from "element-plus"
import useUserStore from '@/store/modules/user'
import router from '@/router'

// 第一步：利用axios对象的create方法，去创建axios实例（其它的配置，基础路径、超时的时间）
const request = axios.create({
    baseURL: import.meta.env.VITE_APP_BASE_API || '/api', // 路径上会携带 /api
    timeout: 15000 // 超时时间 15s（AI 对话单独配置了 60s）
})

// 第二步：request实例添加请求与响应拦截器
request.interceptors.request.use((config) => {
    // 仅当生产环境显式配置了跨域 VITE_SERVE 时才拼接，默认走同源 /api 反向代理
    if (import.meta.env.PROD && import.meta.env.VITE_SERVE) {
        config.baseURL = import.meta.env.VITE_SERVE + (import.meta.env.VITE_APP_BASE_API || '')
    }
    const userStore = useUserStore()
    if (userStore.access_token) {
        config.headers['Authorization'] = `${userStore.access_token}`
    }
    return config
})

// 第三步：响应拦截器
request.interceptors.response.use(
    (response) => {
        const res = response.data
        // 统一拦截未授权或Token失效
        if (res.code === 401 || res.code === 402) {
            const userStore = useUserStore()
            userStore.clearUserInfo()
            ElMessage.error(res.msg || '登录状态已过期，请重新登录')
            router.push('/login')
        } else if (res.code === 403) {
            ElMessage.error(res.msg || '无权限执行此操作')
        }
        return res
    },
    (error) => {
        let message = ''
        const status = error.response?.status
        switch (status) {
            case 401:
                message = "TOKEN过期"
                break
            case 403:
                message = "无权访问"
                break
            case 404:
                message = "请求地址错误"
                break
            case 500:
                message = "服务器出现问题"
                break
            default:
                message = "网络出现问题"
                break
        }

        ElMessage({
            type: 'error',
            message
        })
        return Promise.reject(error)
    }
)

export default request
