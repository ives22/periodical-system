// 引入element-plus提供的全部图标组件
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import useUserStore from '@/store/modules/user'
import pinia from '@/store'

export default {
    install(app: any) {
        for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
            app.component(key, component)
        }

        // 注册全局按钮级权限自定义指令 v-permission="'periodical:create'"
        app.directive('permission', {
            mounted(el: HTMLElement, binding: any) {
                const userStore = useUserStore(pinia)
                const permCode = binding.value
                if (permCode && !userStore.hasPermission(permCode)) {
                    el.style.display = 'none'
                }
            },
            updated(el: HTMLElement, binding: any) {
                const userStore = useUserStore(pinia)
                const permCode = binding.value
                if (permCode && !userStore.hasPermission(permCode)) {
                    el.style.display = 'none'
                } else {
                    el.style.display = ''
                }
            }
        })
    }
}
