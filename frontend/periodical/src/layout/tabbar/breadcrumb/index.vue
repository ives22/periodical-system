<template>
    <!-- 顶部左侧的静态 -->
    <el-icon style="margin-right: 10px;" @click="changeIcon">
        <component :is="layOutSettingStore.fold ? 'Expand' : 'Fold'"></component>
    </el-icon>
    <!-- 左侧面包屑 -->
    <el-breadcrumb separator-icon="ArrowRight">
        <!-- <el-breadcrumb-item :to="{ path: '/' }">homepage</el-breadcrumb-item> -->
        <!-- 面包屑动态展示路由名字与标题 -->
        <el-breadcrumb-item v-for="(item, index) in $route.matched" :key="index" v-show="item.meta.title"
            :to="item.path">
            <!-- 图标 -->
            <el-icon style="vertical-align: top; margin-right:2px;">
                <component v-if="item.meta.title" :is="item.meta.icon"></component>
            </el-icon>
            <span>{{ item.meta.title }}</span>
        </el-breadcrumb-item>
    </el-breadcrumb>
</template>

<script lang="ts" setup>
import { useRoute } from 'vue-router';
import useLayOutSettingStore from '@/store/modules/setting'
// 获取layout配置相关的仓库
let layOutSettingStore = useLayOutSettingStore()
// 获取路由对象
let $route = useRoute()

// 点击修改图标的方法
const changeIcon = () => {
    layOutSettingStore.fold = !layOutSettingStore.fold
}


</script>

<script lang="ts">
export default {
    name: 'Breadcrumb'
}
</script>

<style scoped></style>