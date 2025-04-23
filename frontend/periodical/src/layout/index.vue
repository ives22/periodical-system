<template>
    <div class="layout_container">
        <!-- 左侧菜单 -->
        <div class="layout_slider" :class="{ fold: LayOutSettingStore.fold ? true : false }">
            <!-- 项目logo -->
            <Logo></Logo>
            <!-- 滚动组件 -->
            <el-scrollbar class="scrollbar">
                <!-- 根据路由动态生成菜单 -->
                <el-menu :collapse="LayOutSettingStore.fold ? true : false" :default-active="$route.path"
                    background-color="#001529" text-color="white">
                    <Menu :menuList="userStore.menuRoutes"></Menu>
                </el-menu>
            </el-scrollbar>
        </div>

        <!-- 顶部导航 -->
        <div class="layout_tabbar" :class="{ fold: LayOutSettingStore.fold ? true : false }">
            <Tabbar></Tabbar>
        </div>

        <!-- 内容展示区 -->
        <div class="layout_main" :class="{ fold: LayOutSettingStore.fold ? true : false }">
            <Main></Main>
        </div>
    </div>
</template>

<script setup lang="ts">
// 获取路由对象
import { useRoute } from 'vue-router';
// 引入左侧菜单logo子组件
import Logo from './logo/index.vue'
// 引入菜单组件
import Menu from './menu/index.vue'
// 引入顶部tabbar组件
import Tabbar from './tabbar/index.vue'
// 右侧内容的展示区
import Main from './main/index.vue'

import useUserStore from '@/store/modules/user';
import useLayOutSettingStore from '@/store/modules/setting'
let userStore = useUserStore();
let LayOutSettingStore = useLayOutSettingStore()

// 获取路由对象
let $route = useRoute()
</script>

<script lang="ts">
export default {
    name: 'Layout'
}

</script>

<style scoped lang="scss">
.layout_container {
    width: 100%;
    height: 100vh;

    .layout_slider {
        width: $baseMenuWidth;
        height: 100vh;
        background: $baseMenuBackground;
        color: white;
        transition: all 0.3s;

        .scrollbar {
            height: calc(100vh - $baseMenuLogoHeight);
            width: 100%;

            .el-menu {
                border-right: none;
            }
        }

        &.fold {
            width: $baseMenuMinWidth;
        }
    }

    .layout_tabbar {
        position: fixed;
        width: calc(100% - $baseMenuWidth);
        height: $baseTabbarHeight;
        top: 0px;
        left: $baseMenuWidth;
        transition: all 0.3s;

        &.fold {
            width: calc(100vw - $baseMenuMinWidth);
            left: $baseMenuMinWidth;
        }
    }

    .layout_main {
        position: absolute;
        width: calc(100% - $baseMenuWidth);
        height: calc(100vh - $baseTabbarHeight);
        left: $baseMenuWidth;
        top: $baseTabbarHeight;
        padding: 15px;
        overflow: auto;
        transition: all 0.3s;
        background-color: #f1f1f1;

        &.fold {
            width: calc(100vw - $baseMenuMinWidth);
            left: $baseMenuMinWidth;
        }
    }
}
</style>
