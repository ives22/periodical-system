<template>
    <div class="layout_container">
        <!-- 左侧导航栏 -->
        <aside class="layout_slider" :class="{ fold: LayOutSettingStore.fold }">
            <!-- 品牌 Logo -->
            <Logo />
            <!-- 导航菜单滚动区 -->
            <el-scrollbar class="scrollbar">
                <el-menu
                    :collapse="LayOutSettingStore.fold"
                    :default-active="$route.path"
                    background-color="#0f172a"
                    text-color="#94a3b8"
                    active-text-color="#ffffff"
                    :collapse-transition="false"
                >
                    <Menu :menuList="userStore.menuRoutes" />
                </el-menu>
            </el-scrollbar>
        </aside>

        <!-- 顶部导航栏（上层面包屑顶栏 + 下层多标签页导航栏） -->
        <header class="layout_tabbar" :class="{ fold: LayOutSettingStore.fold }">
            <div class="tabbar_top_row">
                <Tabbar />
            </div>
            <TagsView />
        </header>

        <!-- 主内容展示区 -->
        <main class="layout_main" :class="{ fold: LayOutSettingStore.fold }">
            <Main />
        </main>
    </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import Logo from './logo/index.vue'
import Menu from './menu/index.vue'
import Tabbar from './tabbar/index.vue'
import TagsView from './tagsview/index.vue'
import Main from './main/index.vue'

import useUserStore from '@/store/modules/user'
import useLayOutSettingStore from '@/store/modules/setting'

let userStore = useUserStore()
let LayOutSettingStore = useLayOutSettingStore()
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
    background-color: $bgCanvas;

    .layout_slider {
        position: fixed;
        top: 0;
        left: 0;
        z-index: 101;
        width: $baseMenuWidth;
        height: 100vh;
        background: $baseMenuBackground;
        border-right: 1px solid rgba(255, 255, 255, 0.06);
        transition: width 0.25s cubic-bezier(0.22, 1, 0.36, 1);

        .scrollbar {
            height: calc(100vh - $baseMenuLogoHeight);
            width: 100%;
            padding: 8px 0;

            :deep(.el-menu) {
                border-right: none;
                padding: 0 12px;
            }

            :deep(.el-menu-item),
            :deep(.el-sub-menu__title) {
                height: 44px;
                line-height: 44px;
                margin: 4px 0;
                padding-left: 14px !important;
                padding-right: 12px !important;
                border-radius: 9px;
                font-size: 14px;
                font-weight: 500;
                transition: background-color 0.2s ease, color 0.2s ease;

                .el-icon:not(.el-sub-menu__icon-arrow) {
                    width: 20px;
                    font-size: 18px;
                    margin-right: 10px;
                    display: inline-flex;
                    align-items: center;
                    justify-content: center;
                    flex-shrink: 0;
                }

                &:hover {
                    background-color: rgba(255, 255, 255, 0.07) !important;
                    color: #f8fafc !important;
                }
            }

            :deep(.el-sub-menu .el-menu-item) {
                height: 40px;
                line-height: 40px;
                padding-left: 26px !important;
                font-size: 13.5px;

                .el-icon:not(.el-sub-menu__icon-arrow) {
                    font-size: 16px;
                    margin-right: 8px;
                }
            }

            :deep(.el-menu-item.is-active) {
                background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%) !important;
                color: #ffffff !important;
                box-shadow: 0 4px 12px rgba(37, 99, 235, 0.32);
            }

            :deep(.el-menu--collapse) {
                padding: 0 10px;
                width: 100%;

                .el-menu-item,
                .el-sub-menu__title {
                    padding: 0 !important;
                    width: 44px;
                    height: 44px;
                    margin: 6px auto;
                    display: flex !important;
                    align-items: center !important;
                    justify-content: center !important;

                    .el-menu-tooltip__trigger {
                        padding: 0 !important;
                        display: flex !important;
                        align-items: center !important;
                        justify-content: center !important;
                        width: 100%;
                        height: 100%;
                    }

                    .el-icon:not(.el-sub-menu__icon-arrow) {
                        margin: 0 !important;
                        width: 20px;
                        height: 20px;
                        font-size: 18px;
                    }

                    .el-sub-menu__icon-arrow,
                    span {
                        display: none !important;
                    }
                }
            }
        }

        &.fold {
            width: $baseMenuMinWidth;
        }
    }

    .layout_tabbar {
        position: fixed;
        z-index: 100;
        width: calc(100% - $baseMenuWidth);
        height: calc(#{$baseTabbarHeight} + #{$baseTagsViewHeight});
        top: 0;
        left: $baseMenuWidth;
        background: rgba(255, 255, 255, 0.94);
        backdrop-filter: blur(12px);
        display: flex;
        flex-direction: column;
        transition: all 0.25s cubic-bezier(0.22, 1, 0.36, 1);

        .tabbar_top_row {
            width: 100%;
            height: $baseTabbarHeight;
            border-bottom: 1px solid $borderColor;
            flex-shrink: 0;
        }

        &.fold {
            width: calc(100vw - $baseMenuMinWidth);
            left: $baseMenuMinWidth;
        }
    }

    .layout_main {
        position: absolute;
        width: calc(100% - $baseMenuWidth);
        height: calc(100vh - #{$baseTabbarHeight} - #{$baseTagsViewHeight});
        left: $baseMenuWidth;
        top: calc(#{$baseTabbarHeight} + #{$baseTagsViewHeight});
        padding: 18px 22px;
        overflow: auto;
        background-color: $bgCanvas;
        transition: all 0.25s cubic-bezier(0.22, 1, 0.36, 1);

        &.fold {
            width: calc(100vw - $baseMenuMinWidth);
            left: $baseMenuMinWidth;
        }
    }
}
</style>
