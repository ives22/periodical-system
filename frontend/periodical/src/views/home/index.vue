<template>
    <div class="dashboard-page">
        <!-- 顶部欢迎横幅 -->
        <section class="hero-banner">
            <div class="hero-left">
                <div class="hero-top-tags">
                    <span class="hero-badge">EDITORIAL &amp; PUBLISHING WORKSPACE</span>
                    <span v-if="userStore.roleNames" class="role-pill">{{ userStore.roleNames }}</span>
                </div>
                <h1 class="hero-title">{{ welcomeMessage }}，{{ userStore.loginUser.name }}</h1>
                <p class="hero-subtitle">
                    欢迎回到学术期刊数字化综合管理平台。您可以在此速览馆藏刊物规模，或快速进入期刊管理与 AI 工作台。
                </p>
            </div>
            <div class="hero-right">
                <div class="clock-box">
                    <span class="clock-label">当前系统时间</span>
                    <span class="clock-time">{{ currentTime }}</span>
                </div>
                <div class="hero-actions">
                    <button
                        v-if="userStore.hasPermission('periodical:view')"
                        type="button"
                        class="hero-btn primary"
                        @click="router.push('/periodical')"
                    >
                        <el-icon><Reading /></el-icon>
                        <span>前往期刊目录</span>
                    </button>
                    <button
                        v-if="userStore.hasPermission('ai:workspace')"
                        type="button"
                        class="hero-btn secondary"
                        @click="router.push('/ai/workspace')"
                    >
                        <el-icon><MagicStick /></el-icon>
                        <span>AI 选刊工作台</span>
                    </button>
                </div>
            </div>
        </section>

        <!-- 核心期刊业务 KPI 指标栏 -->
        <el-row :gutter="16" class="metric-row">
            <el-col :xs="24" :sm="12" :md="6">
                <div class="metric-card clickable" @click="router.push('/periodical')">
                    <div class="metric-icon blue">
                        <el-icon :size="22"><Reading /></el-icon>
                    </div>
                    <div class="metric-body">
                        <span class="metric-label">在库期刊总量</span>
                        <div class="metric-value-row">
                            <span class="metric-value">{{ stats.totalCount }}</span>
                            <span class="metric-unit">本</span>
                        </div>
                        <span class="metric-sub">
                            涵盖国家级、省级、学报及国际刊物
                            <template v-if="stats.hiddenCount > 0">（另含已隐藏 {{ stats.hiddenCount }} 本）</template>
                        </span>
                    </div>
                </div>
            </el-col>

            <el-col :xs="24" :sm="12" :md="6">
                <div class="metric-card clickable" @click="router.push('/periodical')">
                    <div class="metric-icon emerald">
                        <el-icon :size="22"><CircleCheckFilled /></el-icon>
                    </div>
                    <div class="metric-body">
                        <span class="metric-label">支持全包刊物</span>
                        <div class="metric-value-row">
                            <span class="metric-value">{{ stats.warpCount }}</span>
                            <span class="metric-unit">本</span>
                        </div>
                        <span class="metric-sub">提供全流程审稿与排期追踪服务</span>
                    </div>
                </div>
            </el-col>

            <el-col :xs="24" :sm="12" :md="6">
                <div class="metric-card clickable" @click="router.push('/periodical')">
                    <div class="metric-icon indigo">
                        <el-icon :size="22"><Tickets /></el-icon>
                    </div>
                    <div class="metric-body">
                        <span class="metric-label">可开杂志社发票</span>
                        <div class="metric-value-row">
                            <span class="metric-value">{{ stats.invoiceCount }}</span>
                            <span class="metric-unit">本</span>
                        </div>
                        <span class="metric-sub">支持正规学术报销与单位财务入账</span>
                    </div>
                </div>
            </el-col>

            <el-col :xs="24" :sm="12" :md="6">
                <div class="metric-card">
                    <div class="metric-icon amber">
                        <el-icon :size="22"><CollectionTag /></el-icon>
                    </div>
                    <div class="metric-body">
                        <span class="metric-label">多维分类体系</span>
                        <div class="metric-value-row">
                            <span class="metric-value">{{ stats.directionCount }}</span>
                            <span class="metric-unit">个收稿方向</span>
                        </div>
                        <span class="metric-sub">
                            已接入 {{ stats.collectionCount }} 大收录数据库 · {{ stats.levelCount }} 类刊物级别
                        </span>
                    </div>
                </div>
            </el-col>
        </el-row>

        <!-- 核心业务工作台快捷入口 -->
        <el-card shadow="never" class="panel-card">
            <template #header>
                <div class="panel-header">
                    <span class="ph-title">核心业务工作台</span>
                </div>
            </template>

            <div class="nav-grid">
                <div
                    v-if="userStore.hasPermission('periodical:view')"
                    class="nav-tile"
                    @click="router.push('/periodical')"
                >
                    <div class="tile-icon primary">
                        <el-icon :size="22"><Reading /></el-icon>
                    </div>
                    <div class="tile-content">
                        <h3>期刊综合管理</h3>
                        <p>按收稿方向、收录网站、期刊级别与出版周期多维筛选及维护期刊库。</p>
                    </div>
                    <el-icon class="tile-arrow"><ArrowRight /></el-icon>
                </div>

                <div
                    v-if="userStore.hasPermission('ai:workspace')"
                    class="nav-tile"
                    @click="router.push('/ai/workspace')"
                >
                    <div class="tile-icon indigo">
                        <el-icon :size="22"><MagicStick /></el-icon>
                    </div>
                    <div class="tile-content">
                        <h3>AI 智能工作台</h3>
                        <p>支持自然语言智能选刊检索、征稿信息自动提取入库与多刊批量修改。</p>
                    </div>
                    <el-icon class="tile-arrow"><ArrowRight /></el-icon>
                </div>

                <div
                    v-if="userStore.hasPermission('user:view')"
                    class="nav-tile"
                    @click="router.push('/system/user')"
                >
                    <div class="tile-icon success">
                        <el-icon :size="22"><User /></el-icon>
                    </div>
                    <div class="tile-content">
                        <h3>团队账号管理</h3>
                        <p>管理系统成员账号状态、重置密码与角色分配。</p>
                    </div>
                    <el-icon class="tile-arrow"><ArrowRight /></el-icon>
                </div>

                <div
                    v-if="userStore.hasPermission('role:view')"
                    class="nav-tile"
                    @click="router.push('/system/role')"
                >
                    <div class="tile-icon warning">
                        <el-icon :size="22"><Lock /></el-icon>
                    </div>
                    <div class="tile-content">
                        <h3>角色与权限配置</h3>
                        <p>按业务岗位配置菜单访问与操作级功能权限。</p>
                    </div>
                    <el-icon class="tile-arrow"><ArrowRight /></el-icon>
                </div>
            </div>
        </el-card>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import useUserStore from '@/store/modules/user'
import { reqPeriodicalList, reqCategorizeList } from '@/api/periodical'

defineOptions({
    name: 'Home',
})

const router = useRouter()
const userStore = useUserStore()
const currentTime = ref('')
const welcomeMessage = ref('')
let timer: any = null

const stats = reactive({
    totalCount: 0,
    hiddenCount: 0,
    warpCount: 0,
    invoiceCount: 0,
    directionCount: 0,
    collectionCount: 0,
    levelCount: 0,
})

const updateTime = () => {
    const now = new Date()
    const hours = now.getHours()
    currentTime.value = now.toLocaleString('zh-CN', { hour12: false })

    if (hours < 12) {
        welcomeMessage.value = '早上好'
    } else if (hours < 18) {
        welcomeMessage.value = '下午好'
    } else {
        welcomeMessage.value = '晚上好'
    }
}

const buildEmptyQueryPayload = (extra: Record<string, any> = {}) => ({
    page_number: 1,
    page_size: 1,
    name: '',
    competent_unit: '',
    organization: '',
    domestic_number: '',
    column_setting: '',
    publish_time_start: '',
    publish_time_end: '',
    invoice: 3,
    is_warp: 3,
    works: 3,
    is_hidden: 2,
    has_base_type_list: [],
    ...extra,
})

const loadDashboardMetrics = async () => {
    try {
        const canManageHidden = userStore.hasPermission('periodical:update')
        const [allRes, warpRes, invRes, catRes, hiddenRes] = await Promise.all([
            reqPeriodicalList(buildEmptyQueryPayload()),
            reqPeriodicalList(buildEmptyQueryPayload({ is_warp: 1 })),
            reqPeriodicalList(buildEmptyQueryPayload({ invoice: 1 })),
            reqCategorizeList({ page_size: 10000, page_number: 1, status: 1 }),
            canManageHidden ? reqPeriodicalList(buildEmptyQueryPayload({ is_hidden: 1 })) : Promise.resolve(null),
        ])

        if (allRes.code === 200 && allRes.data) {
            stats.totalCount = allRes.data.total || 0
        }
        if (hiddenRes && hiddenRes.code === 200 && hiddenRes.data) {
            stats.hiddenCount = hiddenRes.data.total || 0
        }
        if (warpRes.code === 200 && warpRes.data) {
            stats.warpCount = warpRes.data.total || 0
        }
        if (invRes.code === 200 && invRes.data) {
            stats.invoiceCount = invRes.data.total || 0
        }
        if (catRes.code === 200 && catRes.data?.items) {
            const items = catRes.data.items as any[]
            stats.directionCount = items.filter((i) => i.type === 1).length
            stats.levelCount = items.filter((i) => i.type === 2).length
            stats.collectionCount = items.filter((i) => i.type === 3).length
        }
    } catch {
        // ignore
    }
}

onMounted(() => {
    updateTime()
    timer = setInterval(updateTime, 1000)
    loadDashboardMetrics()
})

onBeforeUnmount(() => {
    if (timer) clearInterval(timer)
})
</script>

<style scoped lang="scss">
.dashboard-page {
    display: flex;
    flex-direction: column;
    gap: 18px;
}

.hero-banner {
    background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
    border-radius: 14px;
    padding: 26px 30px;
    color: #ffffff;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 24px;
    box-shadow: 0 10px 25px -5px rgba(15, 23, 42, 0.12);

    .hero-left {
        .hero-top-tags {
            display: flex;
            align-items: center;
            gap: 8px;
            margin-bottom: 10px;
        }

        .hero-badge {
            display: inline-block;
            font-size: 11px;
            font-weight: 600;
            letter-spacing: 0.1em;
            color: #38bdf8;
            background: rgba(56, 189, 248, 0.12);
            padding: 4px 10px;
            border-radius: 999px;
        }

        .role-pill {
            font-size: 11px;
            font-weight: 500;
            color: #cbd5e1;
            background: rgba(255, 255, 255, 0.08);
            border: 1px solid rgba(255, 255, 255, 0.12);
            padding: 3px 10px;
            border-radius: 999px;
        }

        .hero-title {
            font-size: 24px;
            font-weight: 700;
            color: #f8fafc;
            margin-bottom: 8px;
        }

        .hero-subtitle {
            font-size: 13.5px;
            color: #94a3b8;
            line-height: 1.6;
            max-width: 640px;
        }
    }

    .hero-right {
        flex-shrink: 0;
        display: flex;
        flex-direction: column;
        align-items: flex-end;
        gap: 12px;

        .clock-box {
            background: rgba(255, 255, 255, 0.06);
            border: 1px solid rgba(255, 255, 255, 0.1);
            border-radius: 10px;
            padding: 10px 16px;
            text-align: right;
            display: flex;
            align-items: center;
            gap: 10px;

            .clock-label {
                font-size: 11.5px;
                color: #94a3b8;
            }

            .clock-time {
                font-size: 14.5px;
                font-weight: 600;
                color: #f8fafc;
                font-variant-numeric: tabular-nums;
            }
        }

        .hero-actions {
            display: flex;
            align-items: center;
            gap: 10px;

            .hero-btn {
                display: inline-flex;
                align-items: center;
                gap: 6px;
                padding: 8px 15px;
                border-radius: 8px;
                font-size: 13px;
                font-weight: 600;
                cursor: pointer;
                transition: all 0.2s ease;

                &.primary {
                    background: #2563eb;
                    color: #ffffff;
                    border: 1px solid #3b82f6;

                    &:hover {
                        background: #1d4ed8;
                    }
                }

                &.secondary {
                    background: rgba(255, 255, 255, 0.08);
                    color: #f8fafc;
                    border: 1px solid rgba(255, 255, 255, 0.16);

                    &:hover {
                        background: rgba(255, 255, 255, 0.14);
                    }
                }
            }
        }
    }
}

.metric-row {
    .metric-card {
        background: #ffffff;
        border: 1px solid #e2e8f0;
        border-radius: 12px;
        padding: 18px 20px;
        display: flex;
        align-items: flex-start;
        gap: 15px;
        height: 100%;
        box-shadow: 0 1px 2px rgba(15, 23, 42, 0.03);
        transition: all 0.2s ease;

        &.clickable {
            cursor: pointer;

            &:hover {
                border-color: #93c5fd;
                transform: translateY(-1px);
                box-shadow: 0 6px 16px rgba(15, 23, 42, 0.05);
            }
        }

        .metric-icon {
            width: 44px;
            height: 44px;
            border-radius: 10px;
            display: flex;
            align-items: center;
            justify-content: center;
            flex-shrink: 0;

            &.blue {
                background: #eff6ff;
                color: #2563eb;
            }

            &.indigo {
                background: #eef2ff;
                color: #4f46e5;
            }

            &.emerald {
                background: #ecfdf5;
                color: #059669;
            }

            &.amber {
                background: #fffbeb;
                color: #d97706;
            }
        }

        .metric-body {
            display: flex;
            flex-direction: column;
            gap: 4px;

            .metric-label {
                font-size: 12.5px;
                color: #64748b;
                font-weight: 500;
            }

            .metric-value-row {
                display: flex;
                align-items: baseline;
                gap: 4px;

                .metric-value {
                    font-size: 22px;
                    font-weight: 700;
                    color: #0f172a;
                    font-variant-numeric: tabular-nums;
                }

                .metric-unit {
                    font-size: 12.5px;
                    font-weight: 600;
                    color: #475569;
                }
            }

            .metric-sub {
                font-size: 12px;
                color: #94a3b8;
            }
        }
    }
}

.panel-card {
    border-radius: 12px;
    border: 1px solid #e2e8f0;

    .panel-header {
        display: flex;
        align-items: center;
        justify-content: space-between;

        .ph-title {
            font-size: 14.5px;
            font-weight: 700;
            color: #0f172a;
        }
    }
}

.nav-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 14px;

    .nav-tile {
        display: flex;
        align-items: center;
        gap: 16px;
        padding: 16px 18px;
        border-radius: 10px;
        border: 1px solid #e2e8f0;
        background: #ffffff;
        cursor: pointer;
        transition: all 0.2s cubic-bezier(0.22, 1, 0.36, 1);

        &:hover {
            border-color: #93c5fd;
            background: #f8fafc;
            transform: translateY(-1px);
            box-shadow: 0 4px 12px rgba(15, 23, 42, 0.05);

            .tile-arrow {
                transform: translateX(3px);
                color: #2563eb;
            }
        }

        .tile-icon {
            width: 44px;
            height: 44px;
            border-radius: 10px;
            display: flex;
            align-items: center;
            justify-content: center;
            flex-shrink: 0;

            &.primary {
                background: #eff6ff;
                color: #2563eb;
            }

            &.indigo {
                background: #eef2ff;
                color: #4f46e5;
            }

            &.success {
                background: #ecfdf5;
                color: #059669;
            }

            &.warning {
                background: #fffbeb;
                color: #d97706;
            }
        }

        .tile-content {
            flex: 1;

            h3 {
                font-size: 14px;
                font-weight: 600;
                color: #0f172a;
                margin-bottom: 4px;
            }

            p {
                font-size: 12.5px;
                color: #64748b;
                line-height: 1.5;
            }
        }

        .tile-arrow {
            color: #94a3b8;
            transition: all 0.2s ease;
        }
    }
}
</style>
