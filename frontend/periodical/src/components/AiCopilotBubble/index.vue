<template>
    <div class="ai-copilot-wrapper">
        <!-- 展开态：右下角灵动悬浮对话窗 -->
        <transition name="copilot-pop">
            <div v-if="isOpen" class="copilot-window">
                <!-- 顶部栏 -->
                <div class="copilot-header">
                    <div class="header-left">
                        <div class="ai-badge-icon">AI</div>
                        <div class="title-box">
                            <span class="c-title">期刊 AI 伴随助手</span>
                            <span class="c-sub">{{ activeModel?.model_name || 'smart-agent-local' }} · 实时联动列表</span>
                        </div>
                    </div>
                    <div class="header-ops">
                        <el-tooltip content="开启新对话" placement="top">
                            <button type="button" class="hdr-btn" @click="handleResetChat">
                                <el-icon><RefreshRight /></el-icon>
                            </button>
                        </el-tooltip>
                        <el-tooltip content="在 AI 智能工作台中全屏打开" placement="top">
                            <button type="button" class="hdr-btn" @click="handleExpandToWorkspace">
                                <el-icon><FullScreen /></el-icon>
                            </button>
                        </el-tooltip>
                        <el-tooltip content="最小化收起" placement="top">
                            <button type="button" class="hdr-btn" @click="isOpen = false">
                                <el-icon><Close /></el-icon>
                            </button>
                        </el-tooltip>
                    </div>
                </div>

                <!-- 消息内容区 -->
                <div ref="bodyScrollRef" class="copilot-body">
                    <div v-if="messages.length === 0" class="copilot-welcome">
                        <div class="welcome-card">
                            <h4>👋 需要我对当前期刊列表做什么？</h4>
                            <p>支持自然语言<strong>多维选刊并同步左侧筛选、粘贴征稿函自动填表入库、改排期/调价、安全删除</strong>。</p>
                            <div class="scene-btns">
                                <div
                                    class="scene-item"
                                    @click="
                                        sendQuickPrompt(
                                            '帮我筛选知网收录、全包、可开杂志社发票的省级期刊',
                                        )
                                    "
                                >
                                    <span class="s-badge blue">联动筛选</span>
                                    <span class="s-text">“筛选知网收录、全包、可开发票的省级期刊”</span>
                                </div>
                                <div
                                    class="scene-item"
                                    @click="
                                        fillDraftPrompt(
                                            '新增一本期刊《在此输入期刊名称》，级别：省级，出版周期：月刊，网站收录：知网，收稿方向：综合类，出刊时间：2025-10-01，发表费：1800元，支持全包，可开杂志社发票',
                                            '在此输入期刊名称',
                                        )
                                    "
                                >
                                    <span class="s-badge emerald">智能填表</span>
                                    <span class="s-text">填入新刊模板或粘贴征稿介绍自动提取入库（不自动发送）</span>
                                </div>
                                <div
                                    class="scene-item"
                                    @click="
                                        fillDraftPrompt(
                                            '把《在此输入期刊名称》的出刊时间修改为 2025-10-01，发表费调整为 1500 元',
                                            '在此输入期刊名称',
                                        )
                                    "
                                >
                                    <span class="s-badge amber">排期调价</span>
                                    <span class="s-text">填入改刊模板，补充刊名与新价格后发送（不自动发送）</span>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div v-else class="copilot-msg-list">
                        <div
                            v-for="m in messages"
                            :key="m.id"
                            class="c-msg-item"
                            :class="m.role === 'user' ? 'is-user' : 'is-ai'"
                        >
                            <div v-if="m.tool_steps && m.tool_steps.length > 0" class="c-steps">
                                <div
                                    v-for="(st, idx) in m.tool_steps"
                                    :key="idx"
                                    class="c-step-pill"
                                    :class="st.status"
                                >
                                    <span>
                                        ⚡ {{ st.title }}：{{ st.summary }}
                                        <template v-if="st.duration_ms">({{ st.duration_ms }}ms)</template>
                                    </span>
                                </div>
                            </div>

                            <div class="c-bubble" v-html="renderGFMMarkdown(m.content)"></div>

                            <AiActionCard
                                v-if="m.action_card"
                                :card="m.action_card"
                                :message-id="m.id"
                                :show-sync-filter="true"
                                :show-fill-drawer="true"
                                @confirm-action="handleConfirmAction"
                                @sync-filter="(payload) => emit('sync-filter', payload)"
                                @fill-drawer="(payload) => emit('fill-drawer', payload)"
                                @select-journal="(item) => emit('select-journal', item)"
                                @quick-ask="sendQuickPrompt"
                                @start-modify="(item) => { inputPrompt = `把${item.name}的` }"
                                @draft-prompt="(draft) => { inputPrompt = draft }"
                            />

                            <!-- 气泡内追问建议胶囊 -->
                            <div
                                v-if="
                                    m.role === 'assistant' &&
                                    m.action_card?.suggestions &&
                                    m.action_card.suggestions.length > 0
                                "
                                class="c-suggestions"
                            >
                                <span
                                    v-for="sug in m.action_card.suggestions.slice(0, 2)"
                                    :key="sug"
                                    class="c-sug-chip"
                                    @click="sendQuickPrompt(sug)"
                                >
                                    ✨ {{ sug }}
                                </span>
                            </div>
                        </div>

                        <div v-if="sending" class="c-msg-item is-ai">
                            <div class="c-thinking">
                                <span class="dot"></span>
                                <span class="dot"></span>
                                <span class="dot"></span>
                                <span>{{ copilotThinkingText }}</span>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- 底部输入区 -->
                <div class="copilot-footer">
                    <div class="input-row">
                        <el-input
                            ref="copilotInputRef"
                            v-model="inputPrompt"
                            type="textarea"
                            :autosize="{ minRows: 1, maxRows: 3 }"
                            placeholder="输入自然语言增删改查或粘贴征稿信息... (Enter 发送)"
                            @keydown="handleKeydown"
                        />
                        <el-button
                            type="primary"
                            icon="Position"
                            circle
                            :loading="sending"
                            :disabled="!inputPrompt.trim()"
                            @click="handleSend"
                        />
                    </div>
                </div>
            </div>
        </transition>

        <!-- 收起态：右下角灵动悬浮球 -->
        <button type="button" class="copilot-fab" :class="{ active: isOpen }" @click="isOpen = !isOpen">
            <span class="fab-glow"></span>
            <el-icon class="fab-icon"><MagicStick /></el-icon>
            <span class="fab-label">{{ isOpen ? '收起助手' : 'AI 期刊助手' }}</span>
        </button>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import useUserStore from '@/store/modules/user'
import AiActionCard from '@/components/AiActionCard/index.vue'
import { renderGFMMarkdown } from '@/utils/markdown'
import {
    reqActiveDefaultModel,
    reqSendAiMessage,
    reqConfirmAiAction,
    type AIModelConfigItem,
    type AIChatMessageItem,
} from '@/api/ai'

const emit = defineEmits<{
    (e: 'refresh-list'): void
    (e: 'sync-filter', filterPayload: Record<string, any>): void
    (e: 'fill-drawer', formPayload?: Record<string, any>): void
    (e: 'select-journal', item: any): void
}>()

const router = useRouter()
const userStore = useUserStore()
const isOpen = ref(false)
const activeModel = ref<AIModelConfigItem | null>(null)
const sessionId = ref('')
const messages = ref<AIChatMessageItem[]>([])
const inputPrompt = ref('')
const lastSentPrompt = ref('')
const sending = ref(false)
const bodyScrollRef = ref<HTMLElement | null>(null)
const copilotInputRef = ref<any>(null)

const copilotThinkingText = computed(() => {
    const p = lastSentPrompt.value.trim()
    const match = p.match(/《([^》]+)》/)
    const jName = match && match[1] ? `《${match[1]}》` : ''
    if (/新增|录入|添加期刊|入库/.test(p)) {
        return `正在从文本中提取${jName ? ' ' + jName + ' 的' : ''}结构化字段并匹配分类字典...`
    }
    if (/修改|变更|改|更新|调整/.test(p)) {
        return `正在查询${jName || '目标刊物'}原记录并生成字段变更对比单...`
    }
    if (/删除|移除|下架/.test(p)) {
        return `正在核验${jName || '目标刊物'}记录并生成删除确认单...`
    }
    return `正在检索期刊数据库并整理${jName ? ' ' + jName + ' 的' : ''}选刊结果...`
})

const fillDraftPrompt = (text: string, selectToken?: string) => {
    inputPrompt.value = text
    nextTick(() => {
        const el: HTMLTextAreaElement | undefined =
            copilotInputRef.value?.textarea || copilotInputRef.value?.$el?.querySelector('textarea')
        if (!el) {
            copilotInputRef.value?.focus?.()
            return
        }
        el.focus()
        if (selectToken) {
            const idx = inputPrompt.value.indexOf(selectToken)
            if (idx !== -1) {
                el.setSelectionRange(idx, idx + selectToken.length)
                return
            }
        }
        const len = inputPrompt.value.length
        el.setSelectionRange(len, len)
    })
}

const fetchModel = async () => {
    try {
        const res = await reqActiveDefaultModel()
        if (res.code === 200) {
            activeModel.value = res.data
        }
    } catch {
        // ignore
    }
}

const handleResetChat = () => {
    sessionId.value = ''
    messages.value = []
}

const handleExpandToWorkspace = () => {
    if (sessionId.value) {
        router.push({ path: '/ai/workspace', query: { session_id: sessionId.value } })
    } else {
        router.push('/ai/workspace')
    }
}

const sendQuickPrompt = (text: string) => {
    inputPrompt.value = text
    handleSend()
}

const handleKeydown = (e: KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault()
        handleSend()
    }
}

const handleSend = async () => {
    const prompt = inputPrompt.value.trim()
    if (!prompt || sending.value) return

    lastSentPrompt.value = prompt
    inputPrompt.value = ''
    sending.value = true

    const tempUserId = -Date.now()
    const currentUsername = userStore.loginUser?.name || 'admin'
    const optimisticUserMsg: AIChatMessageItem = {
        id: tempUserId,
        session_id: sessionId.value || 'pending',
        user_id: userStore.loginUser?.id || 0,
        username: currentUsername,
        role: 'user',
        content: prompt,
        model_name: '',
        tool_steps: [],
        action_card: null,
        created_at: Math.floor(Date.now() / 1000),
    }
    messages.value.push(optimisticUserMsg)
    scrollToBottom()

    try {
        const res = await reqSendAiMessage({
            session_id: sessionId.value,
            scene: 'copilot',
            prompt,
        })
        if (res.code === 200) {
            sessionId.value = res.data.session.session_id
            const tempIdx = messages.value.findIndex((m) => m.id === tempUserId)
            if (tempIdx !== -1) {
                messages.value.splice(tempIdx, 1, res.data.user_message, res.data.assistant_message)
            } else {
                messages.value.push(res.data.user_message, res.data.assistant_message)
            }
            scrollToBottom()
        } else {
            const tempIdx = messages.value.findIndex((m) => m.id === tempUserId)
            if (tempIdx !== -1) messages.value.splice(tempIdx, 1)
            inputPrompt.value = prompt
            ElMessage.error(res.message || 'AI 响应失败')
        }
    } catch (err: any) {
        const tempIdx = messages.value.findIndex((m) => m.id === tempUserId)
        if (tempIdx !== -1) messages.value.splice(tempIdx, 1)
        inputPrompt.value = prompt
        ElMessage.error(err?.message || '发送请求异常')
    } finally {
        sending.value = false
    }
}

const handleConfirmAction = async ({
    messageId,
    decision,
    formPayload,
    diffs,
    batchUpdates,
}: {
    messageId: number
    decision: 'confirm' | 'cancel'
    formPayload?: Record<string, any>
    diffs?: any[]
    batchUpdates?: any[]
}) => {
    const res = await reqConfirmAiAction({
        message_id: messageId,
        decision,
        form_payload: formPayload,
        diffs,
        batch_updates: batchUpdates,
    })
    if (res.code === 200) {
        const idx = messages.value.findIndex((m) => m.id === messageId)
        if (idx !== -1 && res.data.updated_message) {
            messages.value[idx] = res.data.updated_message
        }
        if (res.data.followup_notice) {
            messages.value.push(res.data.followup_notice)
        }
        if (decision === 'confirm') {
            ElMessage.success('操作已执行，已同步刷新期刊列表')
            emit('refresh-list')
        } else {
            ElMessage.info('已取消操作')
        }
        scrollToBottom()
    } else {
        ElMessage.error(res.message || '操作失败')
    }
}

const scrollToBottom = () => {
    nextTick(() => {
        if (bodyScrollRef.value) {
            bodyScrollRef.value.scrollTop = bodyScrollRef.value.scrollHeight
        }
    })
}

onMounted(() => {
    fetchModel()
})
</script>

<style scoped lang="scss">
.ai-copilot-wrapper {
    position: fixed;
    right: 26px;
    bottom: 24px;
    z-index: 1500;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
}

/* 灵动悬浮窗卡片 */
.copilot-window {
    width: 440px;
    height: 600px;
    max-height: calc(100vh - 110px);
    margin-bottom: 14px;
    background: #ffffff;
    border: 1px solid #cbd5e1;
    border-radius: 16px;
    box-shadow: 0 20px 45px -10px rgba(15, 23, 42, 0.22), 0 8px 18px -6px rgba(15, 23, 42, 0.1);
    display: flex;
    flex-direction: column;
    overflow: hidden;

    .copilot-header {
        padding: 13px 16px;
        background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
        color: #ffffff;
        display: flex;
        align-items: center;
        justify-content: space-between;

        .header-left {
            display: flex;
            align-items: center;
            gap: 10px;

            .ai-badge-icon {
                width: 32px;
                height: 32px;
                border-radius: 9px;
                background: linear-gradient(135deg, #2563eb 0%, #0284c7 100%);
                display: flex;
                align-items: center;
                justify-content: center;
                font-size: 13px;
                font-weight: 800;
                box-shadow: 0 2px 8px rgba(37, 99, 235, 0.4);
            }

            .title-box {
                display: flex;
                flex-direction: column;

                .c-title {
                    font-size: 14px;
                    font-weight: 700;
                    line-height: 1.2;
                }

                .c-sub {
                    font-size: 11px;
                    color: #94a3b8;
                    margin-top: 2px;
                }
            }
        }

        .header-ops {
            display: flex;
            align-items: center;
            gap: 4px;

            .hdr-btn {
                width: 28px;
                height: 28px;
                border-radius: 7px;
                border: none;
                background: rgba(255, 255, 255, 0.08);
                color: #cbd5e1;
                display: flex;
                align-items: center;
                justify-content: center;
                cursor: pointer;
                transition: all 0.15s ease;

                &:hover {
                    background: rgba(255, 255, 255, 0.18);
                    color: #ffffff;
                }
            }
        }
    }

    .copilot-body {
        flex: 1;
        overflow-y: auto;
        padding: 14px;
        background: #f8fafc;
    }

    .copilot-footer {
        padding: 10px 12px;
        background: #ffffff;
        border-top: 1px solid #e2e8f0;

        .input-row {
            display: flex;
            align-items: flex-end;
            gap: 8px;

            :deep(.el-textarea__inner) {
                border-radius: 10px;
                font-size: 13px;
                padding: 8px 12px;
                resize: none;
            }
        }
    }
}

.copilot-welcome {
    .welcome-card {
        background: #ffffff;
        border: 1px solid #e2e8f0;
        border-radius: 12px;
        padding: 16px;

        h4 {
            margin: 0 0 6px;
            font-size: 14.5px;
            color: #0f172a;
        }

        p {
            margin: 0 0 14px;
            font-size: 12.5px;
            color: #64748b;
            line-height: 1.55;
        }

        .scene-btns {
            display: flex;
            flex-direction: column;
            gap: 8px;

            .scene-item {
                padding: 10px 12px;
                border-radius: 9px;
                border: 1px solid #e2e8f0;
                background: #f8fafc;
                cursor: pointer;
                transition: all 0.15s ease;
                display: flex;
                flex-direction: column;
                gap: 4px;

                &:hover {
                    border-color: #93c5fd;
                    background: #eff6ff;
                }

                .s-badge {
                    width: fit-content;
                    font-size: 10.5px;
                    font-weight: 700;
                    padding: 1px 6px;
                    border-radius: 4px;

                    &.blue {
                        background: #dbeafe;
                        color: #1d4ed8;
                    }
                    &.emerald {
                        background: #d1fae5;
                        color: #047857;
                    }
                    &.amber {
                        background: #fef3c7;
                        color: #b45309;
                    }
                }

                .s-text {
                    font-size: 12.5px;
                    color: #334155;
                    font-weight: 500;
                }
            }
        }
    }
}

.copilot-msg-list {
    display: flex;
    flex-direction: column;
    gap: 12px;

    .c-msg-item {
        display: flex;
        flex-direction: column;

        &.is-user {
            align-items: flex-end;

            .c-bubble {
                background: #2563eb;
                color: #ffffff;
                border-color: #2563eb;
            }
        }

        &.is-ai {
            align-items: flex-start;
        }

        .c-steps {
            margin-bottom: 5px;
            display: flex;
            flex-direction: column;
            gap: 4px;

            .c-step-pill {
                font-size: 11.5px;
                padding: 3px 8px;
                border-radius: 6px;
                background: #ecfdf5;
                color: #047857;
                border: 1px solid #a7f3d0;

                &.denied,
                &.error {
                    background: #fef2f2;
                    color: #b91c1c;
                    border-color: #fecaca;
                }
            }
        }

        .c-bubble {
            max-width: 95%;
            padding: 10px 13px;
            border-radius: 11px;
            background: #ffffff;
            border: 1px solid #e2e8f0;
            color: #1e293b;
            font-size: 13px;
            line-height: 1.55;
            word-break: break-word;
        }

        .c-suggestions {
            margin-top: 6px;
            display: flex;
            flex-wrap: wrap;
            gap: 5px;

            .c-sug-chip {
                font-size: 11px;
                padding: 3px 9px;
                border-radius: 999px;
                background: #eff6ff;
                border: 1px solid #bfdbfe;
                color: #1d4ed8;
                cursor: pointer;
                transition: all 0.15s ease;

                &:hover {
                    background: #dbeafe;
                }
            }
        }

        .c-thinking {
            display: inline-flex;
            align-items: center;
            gap: 6px;
            padding: 9px 13px;
            border-radius: 10px;
            background: #ffffff;
            border: 1px solid #e2e8f0;
            font-size: 12.5px;
            color: #64748b;

            .dot {
                width: 5px;
                height: 5px;
                border-radius: 50%;
                background: #2563eb;
            }
        }
    }
}

/* 右下角灵动悬浮触发按钮 */
.copilot-fab {
    position: relative;
    height: 46px;
    padding: 0 18px;
    border-radius: 999px;
    border: 1px solid rgba(255, 255, 255, 0.22);
    background: linear-gradient(135deg, #2563eb 0%, #0284c7 100%);
    color: #ffffff;
    display: inline-flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    box-shadow: 0 10px 25px -4px rgba(37, 99, 235, 0.5), 0 4px 10px -2px rgba(2, 132, 199, 0.35);
    transition: all 0.22s cubic-bezier(0.22, 1, 0.36, 1);

    &:hover {
        transform: translateY(-2px) scale(1.02);
        box-shadow: 0 14px 30px -4px rgba(37, 99, 235, 0.6);
    }

    &.active {
        background: #0f172a;
        box-shadow: 0 8px 20px -4px rgba(15, 23, 42, 0.45);
    }

    .fab-icon {
        font-size: 18px;
    }

    .fab-label {
        font-size: 13.5px;
        font-weight: 700;
        letter-spacing: 0.01em;
    }
}

.copilot-pop-enter-active,
.copilot-pop-leave-active {
    transition: all 0.22s cubic-bezier(0.22, 1, 0.36, 1);
    transform-origin: bottom right;
}

.copilot-pop-enter-from,
.copilot-pop-leave-to {
    opacity: 0;
    transform: scale(0.9) translateY(12px);
}
</style>
