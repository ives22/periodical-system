<template>
    <template v-for="item in menuList" :key="item.path">
        <!-- 顶层透明 Layout 容器：直接展开其内部子菜单 -->
        <template v-if="item.path === '/' && !item.meta?.title && item.children">
            <Menu :menuList="item.children" />
        </template>

        <!-- 没有子路由的一级菜单项 -->
        <template v-else-if="!item.children">
            <el-menu-item v-if="!item.meta?.hidden" :index="item.path" @click="goRoute">
                <el-icon v-if="item.meta?.icon">
                    <component :is="item.meta.icon"></component>
                </el-icon>
                <template #title>
                    <span>{{ item.meta?.title }}</span>
                </template>
            </el-menu-item>
        </template>

        <!-- 有子路由，且只有一个子路由 -->
        <template v-else-if="item.children && item.children.length === 1 && !item.meta?.title">
            <el-menu-item v-if="!item.children[0].meta?.hidden" :index="item.children[0].path" @click="goRoute">
                <el-icon v-if="item.children[0].meta?.icon">
                    <component :is="item.children[0].meta.icon"></component>
                </el-icon>
                <template #title>
                    <span>{{ item.children[0].meta?.title }}</span>
                </template>
            </el-menu-item>
        </template>

        <!-- 有多个子路由（或带分组标题的二级菜单，如“系统管理”） -->
        <el-sub-menu v-else-if="item.children && item.children.length > 0 && !item.meta?.hidden" :index="item.path">
            <template #title>
                <el-icon v-if="item.meta?.icon">
                    <component :is="item.meta.icon"></component>
                </el-icon>
                <span>{{ item.meta?.title }}</span>
            </template>
            <Menu :menuList="item.children"></Menu>
        </el-sub-menu>
    </template>
</template>

<script lang="ts" setup>
import { useRouter } from 'vue-router'

defineProps(['menuList'])
const $router = useRouter()

const goRoute = (vc: any) => {
    if (vc?.index) {
        $router.push(vc.index)
    }
}
</script>

<script lang="ts">
export default {
    name: 'Menu',
}
</script>

<style scoped></style>
