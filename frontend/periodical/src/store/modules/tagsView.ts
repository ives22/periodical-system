import { defineStore } from 'pinia'
import type { RouteLocationNormalizedLoaded } from 'vue-router'

export interface TagViewItem {
    path: string
    fullPath: string
    name: string
    title: string
    icon?: string
    affix?: boolean
    keepAlive?: boolean
}

const STORAGE_KEY = 'periodical_visited_views_v1'

const AFFIX_HOME_TAG: TagViewItem = {
    path: '/home',
    fullPath: '/home',
    name: 'Home',
    title: '首页',
    icon: 'HomeFilled',
    affix: true,
    keepAlive: true,
}

function loadSavedViews(): TagViewItem[] {
    try {
        const raw = sessionStorage.getItem(STORAGE_KEY)
        if (!raw) return [{ ...AFFIX_HOME_TAG }]
        const parsed = JSON.parse(raw)
        if (!Array.isArray(parsed) || parsed.length === 0) {
            return [{ ...AFFIX_HOME_TAG }]
        }
        const hasHome = parsed.some((v: TagViewItem) => v.path === '/home')
        if (!hasHome) {
            parsed.unshift({ ...AFFIX_HOME_TAG })
        } else {
            parsed.forEach((v: TagViewItem) => {
                if (v.path === '/home') v.affix = true
            })
        }
        return parsed
    } catch {
        return [{ ...AFFIX_HOME_TAG }]
    }
}

function persistViews(views: TagViewItem[]) {
    try {
        sessionStorage.setItem(STORAGE_KEY, JSON.stringify(views))
    } catch {
        // ignore storage quota errors
    }
}

const useTagsViewStore = defineStore('TagsViewStore', {
    state: () => {
        const initialVisited = loadSavedViews()
        const initialCached = initialVisited
            .filter((v) => v.keepAlive !== false && v.name)
            .map((v) => v.name)
        return {
            visitedViews: initialVisited as TagViewItem[],
            cachedViews: Array.from(new Set(initialCached)) as string[],
        }
    },

    actions: {
        addView(route: RouteLocationNormalizedLoaded) {
            const path = route.path
            const title = route.meta?.title as string | undefined
            if (!path || path === '/login' || !title) return

            const routeName = route.name ? String(route.name) : ''
            const icon = (route.meta?.icon as string) || ''
            const affix = Boolean(route.meta?.affix) || path === '/home'
            const keepAlive = route.meta?.keepAlive !== false

            const existingIdx = this.visitedViews.findIndex((v) => v.path === path)
            if (existingIdx !== -1) {
                this.visitedViews[existingIdx].fullPath = route.fullPath || path
                this.visitedViews[existingIdx].title = title
                if (routeName) this.visitedViews[existingIdx].name = routeName
                if (icon) this.visitedViews[existingIdx].icon = icon
            } else {
                this.visitedViews.push({
                    path,
                    fullPath: route.fullPath || path,
                    name: routeName,
                    title,
                    icon,
                    affix,
                    keepAlive,
                })
            }

            if (keepAlive && routeName && !this.cachedViews.includes(routeName)) {
                this.cachedViews.push(routeName)
            }

            persistViews(this.visitedViews)
        },

        addCachedView(viewName: string) {
            if (viewName && !this.cachedViews.includes(viewName)) {
                this.cachedViews.push(viewName)
            }
        },

        delCachedView(viewName: string) {
            const idx = this.cachedViews.indexOf(viewName)
            if (idx !== -1) {
                this.cachedViews.splice(idx, 1)
            }
        },

        delView(targetPath: string): TagViewItem[] {
            const target = this.visitedViews.find((v) => v.path === targetPath)
            if (!target || target.affix) return this.visitedViews

            this.visitedViews = this.visitedViews.filter((v) => v.path !== targetPath)
            if (target.name) {
                this.delCachedView(target.name)
            }
            persistViews(this.visitedViews)
            return this.visitedViews
        },

        delOthersViews(currentPath: string) {
            this.visitedViews = this.visitedViews.filter(
                (v) => v.affix || v.path === currentPath,
            )
            const allowedNames = new Set(
                this.visitedViews.filter((v) => v.keepAlive !== false && v.name).map((v) => v.name),
            )
            this.cachedViews = this.cachedViews.filter((name) => allowedNames.has(name))
            persistViews(this.visitedViews)
        },

        delAllViews() {
            this.visitedViews = this.visitedViews.filter((v) => v.affix)
            const allowedNames = new Set(
                this.visitedViews.filter((v) => v.keepAlive !== false && v.name).map((v) => v.name),
            )
            this.cachedViews = this.cachedViews.filter((name) => allowedNames.has(name))
            persistViews(this.visitedViews)
        },
    },
})

export default useTagsViewStore
