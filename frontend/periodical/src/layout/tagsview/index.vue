<template>
    <div class="tags-view-bar" @contextmenu.prevent>
        <!-- 左侧可横向滚动的已打开页签列表（严格还原截图3视觉风格） -->
        <div ref="scrollWrapRef" class="tags-scroll-wrap" @wheel.prevent="handleScrollWheel">
            <div
                v-for="tag in tagsViewStore.visitedViews"
                :key="tag.path"
                class="tag-view-item"
                :class="{ active: isActive(tag), 'is-affix': tag.affix }"
                @click="handleClickTag(tag)"
                @contextmenu.prevent="openContextMenu($event, tag)"
            >
                <span class="tag-title">{{ tag.title }}</span>
                <span
                    v-if="!tag.affix"
                    class="tag-close-btn"
                    title="关闭此标签页"
                    @click.stop="handleCloseTag(tag)"
                >
                    <el-icon :size="12"><Close /></el-icon>
                </span>
                <!-- 截图3风格底部蓝色高亮指示横线 -->
                <span v-if="isActive(tag)" class="active-bottom-bar"></span>
            </div>
        </div>

        <!-- 右侧快捷标签管理下拉按钮 -->
        <div class="tags-right-tools">
            <el-dropdown trigger="click" @command="handleQuickCommand">
                <button type="button" class="tags-ops-btn" title="标签页快捷操作">
                    <span>标签操作</span>
                    <el-icon :size="12"><ArrowDown /></el-icon>
                </button>
                <template #dropdown>
                    <el-dropdown-menu>
                        <el-dropdown-item command="refresh" icon="Refresh">重新加载当前页</el-dropdown-item>
                        <el-dropdown-item
                            command="closeCurrent"
                            icon="Close"
                            :disabled="isCurrentAffix"
                        >
                            关闭当前标签页
                        </el-dropdown-item>
                        <el-dropdown-item
                            command="closeOthers"
                            icon="Remove"
                            :disabled="tagsViewStore.visitedViews.length <= 1"
                        >
                            关闭其他标签页
                        </el-dropdown-item>
                        <el-dropdown-item
                            command="closeAll"
                            divided
                            icon="CircleClose"
                            :disabled="tagsViewStore.visitedViews.length <= 1"
                        >
                            关闭全部标签页
                        </el-dropdown-item>
                    </el-dropdown-menu>
                </template>
            </el-dropdown>
        </div>

        <!-- 右键上下文菜单 -->
        <Teleport to="body">
            <ul
                v-if="ctxMenuVisible"
                class="tags-context-menu"
                :style="{ left: `${ctxMenuLeft}px`, top: `${ctxMenuTop}px` }"
            >
                <li @click="handleRefreshSelectedTag">
                    <el-icon><Refresh /></el-icon>
                    <span>重新加载此页</span>
                </li>
                <li v-if="!selectedTag?.affix" @click="handleCloseSelectedTag">
                    <el-icon><Close /></el-icon>
                    <span>关闭此标签页</span>
                </li>
                <li @click="handleCloseOthers">
                    <el-icon><Remove /></el-icon>
                    <span>关闭其他标签页</span>
                </li>
                <li class="danger" @click="handleCloseAll">
                    <el-icon><CircleClose /></el-icon>
                    <span>关闭全部标签页</span>
                </li>
            </ul>
        </Teleport>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import useTagsViewStore, { type TagViewItem } from '@/store/modules/tagsView'
import useLayOutSettingStore from '@/store/modules/setting'

const route = useRoute()
const router = useRouter()
const tagsViewStore = useTagsViewStore()
const settingStore = useLayOutSettingStore()

const scrollWrapRef = ref<HTMLElement | null>(null)
const ctxMenuVisible = ref(false)
const ctxMenuLeft = ref(0)
const ctxMenuTop = ref(0)
const selectedTag = ref<TagViewItem | null>(null)

const isActive = (tag: TagViewItem) => tag.path === route.path

const isCurrentAffix = computed(() => {
    const current = tagsViewStore.visitedViews.find((v) => v.path === route.path)
    return Boolean(current?.affix)
})

// 监听路由变化，自动登记已打开的页面标签与 Keep-Alive 缓存
watch(
    () => route.fullPath,
    () => {
        tagsViewStore.addView(route)
    },
    { immediate: true },
)

const handleScrollWheel = (e: WheelEvent) => {
    if (!scrollWrapRef.value) return
    scrollWrapRef.value.scrollLeft += e.deltaY || e.deltaX
}

const handleClickTag = (tag: TagViewItem) => {
    if (tag.path === route.path) return
    router.push(tag.fullPath || tag.path)
}

const toLastView = (visitedViews: TagViewItem[]) => {
    const latestView = visitedViews.slice(-1)[0]
    if (latestView) {
        router.push(latestView.fullPath || latestView.path)
    } else {
        router.push('/home')
    }
}

const handleCloseTag = (tag: TagViewItem) => {
    if (tag.affix) return
    const wasActive = isActive(tag)
    const remaining = tagsViewStore.delView(tag.path)
    if (wasActive) {
        toLastView(remaining)
    }
}

const triggerRefreshCurrent = (targetTag?: TagViewItem) => {
    const tag = targetTag || tagsViewStore.visitedViews.find((v) => v.path === route.path)
    if (!tag) return
    if (tag.path !== route.path) {
        router.push(tag.fullPath || tag.path).then(() => {
            settingStore.triggerRefreshCurrentView(tag.name)
        })
    } else {
        settingStore.triggerRefreshCurrentView(tag.name)
    }
}

const handleQuickCommand = (cmd: string) => {
    const currentTag = tagsViewStore.visitedViews.find((v) => v.path === route.path)
    if (cmd === 'refresh') {
        triggerRefreshCurrent(currentTag)
    } else if (cmd === 'closeCurrent' && currentTag && !currentTag.affix) {
        handleCloseTag(currentTag)
    } else if (cmd === 'closeOthers') {
        tagsViewStore.delOthersViews(route.path)
    } else if (cmd === 'closeAll') {
        tagsViewStore.delAllViews()
        if (route.path !== '/home') {
            router.push('/home')
        }
    }
}

const openContextMenu = (e: MouseEvent, tag: TagViewItem) => {
    selectedTag.value = tag
    ctxMenuLeft.value = Math.min(e.clientX, window.innerWidth - 165)
    ctxMenuTop.value = e.clientY + 4
    ctxMenuVisible.value = true
}

const closeContextMenu = () => {
    ctxMenuVisible.value = false
}

const handleRefreshSelectedTag = () => {
    if (selectedTag.value) {
        triggerRefreshCurrent(selectedTag.value)
    }
    closeContextMenu()
}

const handleCloseSelectedTag = () => {
    if (selectedTag.value) {
        handleCloseTag(selectedTag.value)
    }
    closeContextMenu()
}

const handleCloseOthers = () => {
    if (selectedTag.value) {
        if (selectedTag.value.path !== route.path) {
            router.push(selectedTag.value.fullPath || selectedTag.value.path)
        }
        tagsViewStore.delOthersViews(selectedTag.value.path)
    }
    closeContextMenu()
}

const handleCloseAll = () => {
    tagsViewStore.delAllViews()
    if (route.path !== '/home') {
        router.push('/home')
    }
    closeContextMenu()
}

onMounted(() => {
    window.addEventListener('click', closeContextMenu)
})

onBeforeUnmount(() => {
    window.removeEventListener('click', closeContextMenu)
})
</script>

<script lang="ts">
export default {
    name: 'TagsView',
}
</script>

<style scoped lang="scss">
.tags-view-bar {
    width: 100%;
    height: 38px;
    background: #ffffff;
    border-bottom: 1px solid #e2e8f0;
    display: flex;
    align-items: stretch;
    justify-content: space-between;
    user-select: none;

    .tags-scroll-wrap {
        flex: 1;
        min-width: 0;
        display: flex;
        align-items: stretch;
        overflow-x: auto;
        overflow-y: hidden;
        scrollbar-width: none;

        &::-webkit-scrollbar {
            display: none;
        }

        .tag-view-item {
            position: relative;
            display: inline-flex;
            align-items: center;
            gap: 8px;
            padding: 0 20px;
            height: 100%;
            font-size: 13px;
            color: #475569;
            background: #ffffff;
            border-right: 1px solid #f1f5f9;
            cursor: pointer;
            white-space: nowrap;
            transition: background-color 0.16s ease, color 0.16s ease;

            &:hover {
                background: #f8fafc;
                color: #1e293b;

                .tag-close-btn {
                    opacity: 1;
                }
            }

            /* 截图3同款激活样式：清透浅蓝背景 + 品牌蓝文字 + 底部2px实心蓝条 */
            &.active {
                background: #eff6ff;
                color: #2563eb;
                font-weight: 600;

                .tag-close-btn {
                    opacity: 1;
                    color: #2563eb;

                    &:hover {
                        background: #dbeafe;
                        color: #1d4ed8;
                    }
                }
            }

            .tag-title {
                line-height: 1;
            }

            .tag-close-btn {
                width: 16px;
                height: 16px;
                border-radius: 4px;
                display: inline-flex;
                align-items: center;
                justify-content: center;
                color: #94a3b8;
                opacity: 0.75;
                transition: all 0.15s ease;

                &:hover {
                    background: #e2e8f0;
                    color: #0f172a;
                }
            }

            .active-bottom-bar {
                position: absolute;
                left: 12px;
                right: 12px;
                bottom: 0;
                height: 2px;
                border-radius: 2px 2px 0 0;
                background: #2563eb;
            }
        }
    }

    .tags-right-tools {
        flex-shrink: 0;
        display: flex;
        align-items: center;
        padding: 0 14px;
        border-left: 1px solid #f1f5f9;
        background: #ffffff;

        .tags-ops-btn {
            display: inline-flex;
            align-items: center;
            gap: 4px;
            border: 1px solid #e2e8f0;
            background: #f8fafc;
            color: #475569;
            font-size: 12px;
            font-weight: 500;
            padding: 4px 10px;
            border-radius: 6px;
            cursor: pointer;
            transition: all 0.15s ease;

            &:hover {
                color: #2563eb;
                border-color: #bfdbfe;
                background: #eff6ff;
            }
        }
    }
}

.tags-context-menu {
    position: fixed;
    z-index: 3000;
    margin: 0;
    padding: 5px;
    list-style: none;
    background: #ffffff;
    border: 1px solid #e2e8f0;
    border-radius: 9px;
    box-shadow: 0 10px 25px -5px rgba(15, 23, 42, 0.12), 0 4px 10px -4px rgba(15, 23, 42, 0.06);
    min-width: 148px;

    li {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 7px 12px;
        font-size: 12.5px;
        font-weight: 500;
        color: #334155;
        border-radius: 6px;
        cursor: pointer;
        transition: all 0.14s ease;

        &:hover {
            background: #eff6ff;
            color: #2563eb;
        }

        &.danger:hover {
            background: #fef2f2;
            color: #dc2626;
        }
    }
}
</style>
