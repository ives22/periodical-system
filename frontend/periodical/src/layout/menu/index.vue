<template>
    <template v-for="(item, index) in menuList" :key="item.path">
        <!-- 没有子路由 -->
        <template v-if="!item.children">
            <el-menu-item v-if="!item.meta.hidden" :index="item.path" @click="goRoute">
                <el-icon>
                    <component :is="item.meta.icon"></component>
                </el-icon>
                <template #title>
                    <span>{{ item.meta.title }}</span>
                </template>
            </el-menu-item>
        </template>

        <!-- 有子路由，且只有一个路由 -->
        <template v-if="item.children && item.children.length == 1">
            <el-menu-item v-if="!item.children[0].meta.hidden" :index="item.children[0].path" @click="goRoute">
                <el-icon>
                    <component :is="item.children[0].meta.icon"></component>
                </el-icon>
                <template #title>
                    <span>{{ item.children[0].meta.title }}</span>
                </template>
            </el-menu-item>
        </template>


        <!-- 有子路由，不止一个子路由 -->
        <el-sub-menu v-if="item.children && item.children.length > 1" :index="item.path">
            <template #title>
                <el-icon>
                    <component :is="item.meta.icon"></component>
                </el-icon>
                <span>{{ item.meta.title }}</span>
            </template>
            <Menu :menuList="item.children"></Menu>
        </el-sub-menu>

    </template>


</template>

<script lang="ts" setup>
// 获取父组件传递过来的全部路由数组
defineProps(['menuList']);
// 获取路由器对象
import { useRouter } from 'vue-router';
let $router = useRouter()

// 点击菜单的回调
const goRoute = (vc: any) => {
    console.log(vc.index)
    // 路由填装
    $router.push(vc.index)

}

</script>

<script lang="ts">
export default {
    name: "Menu"
}
</script>

<style scoped></style>