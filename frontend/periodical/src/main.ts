// 引入模板的全局样式
import '@/styles/index.scss'
import { createApp } from 'vue'
import App from './App.vue'

// 引入element-plus插件和样式
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
// 配置element-plus国际化语言
import zhCn from 'element-plus/es/locale/lang/zh-cn'
// 引入全局组件
import GlobalComponents from '@/components/index'


// 引入路由
import router from '@/router'
// 引入仓库
import pinia from './store'
// 引入permission
import './permission'

const app = createApp(App)
// 注册element-plus插件
app.use(ElementPlus, {
    locale: zhCn
})


// 注册模板路由
app.use(router)
// 注册仓库
app.use(pinia)
// 注册全局组件
app.use(GlobalComponents)

app.mount('#app')
