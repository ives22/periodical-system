// 小仓库：layout组件相关配置仓库
import { defineStore } from "pinia";

let useLayOutSettingStore = defineStore('SettingStore', {
    state: () => {
        return {
            fold: false, // 用于控制菜单折叠还是收起
            refsh: false, // 用于控制刷新效果
            refreshTargetName: '', // 指定需要单页精准刷新的视图组件名
        }
    },
    actions: {
        triggerRefreshCurrentView(viewName?: string) {
            this.refreshTargetName = viewName || ''
            this.refsh = !this.refsh
        }
    }
})

export default useLayOutSettingStore