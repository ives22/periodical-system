import request from '@/utils/request'

export interface AIModelConfigItem {
    id: number
    name: string
    provider: 'qwen' | 'deepseek' | 'openai' | 'custom' | string
    base_url: string
    api_key?: string
    api_key_masked: string
    has_api_key: boolean
    model_name: string
    system_prompt: string
    temperature: number
    max_tokens: number
    is_default: number
    status: number
    remark: string
    created_at?: number
    updated_at?: number
}

export interface AIChatSessionItem {
    id: number
    session_id: string
    user_id: number
    username: string
    title: string
    scene: 'workspace' | 'copilot' | string
    model_name: string
    message_count: number
    is_pinned: number
    created_at: number
    updated_at: number
}

export interface ToolStepItem {
    tool_name: string
    title: string
    args_summary?: string
    summary: string
    duration_ms?: number
    status: 'success' | 'error' | 'denied' | string
}

export interface FieldDiffItem {
    field: string
    label: string
    old_value: string
    new_value: string
}

export interface BatchUpdateTarget {
    periodical_id: number
    periodical_name: string
    level?: string
    form_payload: Record<string, any>
    diffs: FieldDiffItem[]
}

export interface PeriodicalBriefCard {
    periodical_id: number
    name: string
    level: string
    directions: string
    collections: string[]
    cycle: string
    publication_time: string
    price: string
    periodical_page: string
    is_warp: boolean
    invoice_receipt: boolean
    works: boolean
    is_hidden?: boolean
    organization: string
    competent_unit: string
    domestic_number: string
    international_number: string
    audit_time: string
    periodical_description: string
}

export interface ActionCardData {
    card_type: 'query_result' | 'create_confirm' | 'update_confirm' | 'delete_confirm'
    status: 'pending' | 'confirmed' | 'cancelled' | 'readonly'
    title: string
    description: string
    total_count?: number
    items?: PeriodicalBriefCard[]
    filter_payload?: Record<string, any>
    periodical_id?: number
    periodical_name?: string
    form_payload?: Record<string, any>
    category_labels?: Record<string, string[]>
    diffs?: FieldDiffItem[]
    batch_updates?: BatchUpdateTarget[]
    suggestions?: string[]
    executed_at?: number
    result_message?: string
}

export interface AIChatMessageItem {
    id: number
    session_id: string
    user_id: number
    username: string
    role: 'user' | 'assistant'
    content: string
    model_name: string
    tool_steps: ToolStepItem[]
    action_card?: ActionCardData | null
    created_at: number
}

enum API {
    ACTIVE_MODEL_URL = '/v1/ai/active-model',
    SESSIONS_URL = '/v1/ai/sessions',
    CHAT_URL = '/v1/ai/chat',
    CONFIRM_ACTION_URL = '/v1/ai/action/confirm',
    MODELS_URL = '/v1/ai/models',
    AUDIT_SESSIONS_URL = '/v1/ai/audit/sessions',
}

// 获取当前激活的全局默认模型
export const reqActiveDefaultModel = () =>
    request.get<any, { code: number; data: AIModelConfigItem; message: string }>(API.ACTIVE_MODEL_URL)

// 获取当前登录用户的会话列表
export const reqUserSessions = (keyword: string = '') =>
    request.get<any, { code: number; data: AIChatSessionItem[]; message: string }>(API.SESSIONS_URL, {
        params: { keyword },
    })

// 新建会话
export const reqCreateSession = (data: { title?: string; scene?: string }) =>
    request.post<any, { code: number; data: AIChatSessionItem; message: string }>(API.SESSIONS_URL, data)

// 重命名或置顶会话
export const reqUpdateSession = (sessionId: string, data: { title?: string; is_pinned?: number }) =>
    request.put<any, { code: number; data: any; message: string }>(`${API.SESSIONS_URL}/${sessionId}`, data)

// 删除会话
export const reqDeleteSession = (sessionId: string) =>
    request.delete<any, { code: number; data: any; message: string }>(`${API.SESSIONS_URL}/${sessionId}`)

// 获取会话历史消息（个人或超管审计）
export const reqSessionMessages = (sessionId: string) =>
    request.get<any, { code: number; data: AIChatMessageItem[]; message: string }>(
        `${API.SESSIONS_URL}/${sessionId}/messages`,
    )

// 发送自然语言对话指令（允许大模型推理最长 60s）
export const reqSendAiMessage = (data: {
    session_id?: string
    scene?: string
    prompt: string
    page_context?: Record<string, any>
}) =>
    request.post<
        any,
        {
            code: number
            data: {
                session: AIChatSessionItem
                user_message: AIChatMessageItem
                assistant_message: AIChatMessageItem
            }
            message: string
        }
    >(API.CHAT_URL, data, { timeout: 60000 })

// 确认或取消写操作卡片（新增/修改/删除）
export const reqConfirmAiAction = (data: {
    message_id: number
    decision: 'confirm' | 'cancel'
    form_payload?: Record<string, any>
    diffs?: FieldDiffItem[]
    batch_updates?: BatchUpdateTarget[]
}) =>
    request.post<
        any,
        {
            code: number
            data: {
                updated_message: AIChatMessageItem
                followup_notice?: AIChatMessageItem
            }
            message: string
        }
    >(API.CONFIRM_ACTION_URL, data)

// ===== 超级管理员：大模型渠道管理与审计接口 =====

export const reqAiModelList = () =>
    request.get<any, { code: number; data: AIModelConfigItem[]; message: string }>(API.MODELS_URL)

export const reqSaveAiModel = (data: Partial<AIModelConfigItem>) =>
    request.post<any, { code: number; data: AIModelConfigItem; message: string }>(API.MODELS_URL, data)

export const reqSetDefaultAiModel = (id: number) =>
    request.put<any, { code: number; data: any; message: string }>(`${API.MODELS_URL}/${id}/default`)

export const reqTestAiModelConnection = (id: number) =>
    request.post<
        any,
        {
            code: number
            data: {
                success: boolean
                latency_ms: number
                model_name: string
                reply: string
                message: string
            }
            message: string
        }
    >(`${API.MODELS_URL}/${id}/test`, {}, { timeout: 30000 })

export const reqDeleteAiModel = (id: number) =>
    request.delete<any, { code: number; data: any; message: string }>(`${API.MODELS_URL}/${id}`)

export const reqFetchRemoteAiModels = (data: {
    id?: number
    provider?: string
    base_url?: string
    api_key?: string
}) =>
    request.post<
        any,
        {
            code: number
            data: {
                models: string[]
                from_remote: boolean
                message: string
            }
            message: string
        }
    >(`${API.MODELS_URL}/remote-list`, data, { timeout: 20000 })

export const reqAuditSessions = (data: {
    page_number: number
    page_size: number
    username?: string
    keyword?: string
    scene?: string
}) =>
    request.post<
        any,
        {
            code: number
            data: {
                total: number
                items: AIChatSessionItem[]
            }
            message: string
        }
    >(API.AUDIT_SESSIONS_URL, data)
