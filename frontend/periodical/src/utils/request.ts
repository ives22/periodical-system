import axios from "axios"
import { ElMessage } from "element-plus"
import useUserStore from '@/store/modules/user';


// 第一步：利用axios对象的create方法，去创建axios实例（其它的配置，基础路径、超时的时间）
const request = axios.create({
    baseURL: import.meta.env.VITE_APP_BASE_API, // 路径上会携带 /api
    timeout: 5000 // 超时时间
})

// 第二步：request实例添加请求与响应拦截器
request.interceptors.request.use((config) => {
    // 在生产环境下添加完整的基础URL
    if (process.env.NODE_ENV === 'production') {
        config.baseURL = import.meta.env.VITE_SERVE + import.meta.env.VITE_APP_BASE_API
    }
    const userStore = useUserStore();
    if (userStore.access_token) {
        config.headers['Authorization'] = `${userStore.access_token}`;
    }
    return config
})

// 第三部：响应拦截器
request.interceptors.response.use((response) => {
    // 成功回调
    if (response.data.code == 402) {
        console.log("登录失败")
    }
    return response.data
},
    (error) => {
        // 失败回调：处理http网络错误的
        // 定义一个变量，存储网络错误信息
        let message = ''
        // http状态码
        const status = error.response.status
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

        // 提示错误信息
        ElMessage({
            type: 'error',
            message
        })
        return Promise.reject(error)
    }
)


// 对外暴露
export default request
