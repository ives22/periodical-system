<template>
    <router-view v-slot="{ Component, route }">
        <transition name="fade" mode="out-in">
            <keep-alive :include="tagsViewStore.cachedViews">
                <component :is="Component" :key="route.path" v-if="flag" />
            </keep-alive>
        </transition>
    </router-view>
</template>

<script setup lang="ts">
import { watch, ref, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import useLayOutSettingStore from '@/store/modules/setting'
import useTagsViewStore from '@/store/modules/tagsView'

const route = useRoute()
const layOutSettingStore = useLayOutSettingStore()
const tagsViewStore = useTagsViewStore()
const flag = ref(true)

watch(() => layOutSettingStore.refsh, () => {
    const targetName = layOutSettingStore.refreshTargetName || (route.name ? String(route.name) : '')
    if (targetName) {
        tagsViewStore.delCachedView(targetName)
    }
    flag.value = false
    nextTick(() => {
        if (targetName && route.meta?.keepAlive !== false) {
            tagsViewStore.addCachedView(targetName)
        }
        layOutSettingStore.refreshTargetName = ''
        flag.value = true
    })
})
</script>

<script lang="ts">
export default {
    name: 'Main'
}
</script>

<style scoped>
.fade-enter-from,
.fade-leave-to {
    opacity: 0;
    transform: translateY(6px);
}

.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.2s cubic-bezier(0.22, 1, 0.36, 1), transform 0.2s cubic-bezier(0.22, 1, 0.36, 1);
}

.fade-enter-to,
.fade-leave-from {
    opacity: 1;
    transform: translateY(0);
}
</style>
