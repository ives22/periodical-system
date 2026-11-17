<template>
    <div class="breadcrumb-bar">
        <button
            class="collapse-btn"
            type="button"
            :title="layOutSettingStore.fold ? '展开导航' : '收起导航'"
            @click="changeIcon"
        >
            <el-icon :size="17">
                <component :is="layOutSettingStore.fold ? 'Expand' : 'Fold'" />
            </el-icon>
        </button>

        <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/home' }">
                <span class="crumb-item root-crumb">
                    <el-icon class="crumb-icon"><HomeFilled /></el-icon>
                    <span>期刊管理系统</span>
                </span>
            </el-breadcrumb-item>

            <el-breadcrumb-item
                v-for="(item, index) in matchedCrumbs"
                :key="item.path || index"
                :to="item.redirect ? item.redirect : item.path"
            >
                <span class="crumb-item" :class="{ 'is-current': index === matchedCrumbs.length - 1 }">
                    <span>{{ item.meta.title }}</span>
                </span>
            </el-breadcrumb-item>
        </el-breadcrumb>
    </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import useLayOutSettingStore from '@/store/modules/setting'

const layOutSettingStore = useLayOutSettingStore()
const $route = useRoute()

const matchedCrumbs = computed(() => {
    return $route.matched.filter((item) => item.meta && item.meta.title)
})

const changeIcon = () => {
    layOutSettingStore.fold = !layOutSettingStore.fold
}
</script>

<script lang="ts">
export default {
    name: 'Breadcrumb',
}
</script>

<style scoped lang="scss">
.breadcrumb-bar {
    display: flex;
    align-items: center;
    gap: 12px;

    .collapse-btn {
        width: 32px;
        height: 32px;
        border-radius: 7px;
        border: 1px solid #e2e8f0;
        background: #ffffff;
        color: #475569;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        transition: all 0.2s ease;

        &:hover {
            color: #0f172a;
            background: #f8fafc;
            border-color: #cbd5e1;
        }
    }

    .crumb-item {
        display: inline-flex;
        align-items: center;
        gap: 5px;
        font-size: 13px;
        font-weight: 500;
        color: #64748b;

        &.root-crumb {
            color: #475569;
        }

        &.is-current {
            color: #1e293b;
            font-weight: 600;
        }

        .crumb-icon {
            font-size: 14px;
            color: #64748b;
        }
    }
}
</style>
