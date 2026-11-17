<template>
    <div class="ai-config-container">
        <!-- 页面标题栏 -->
        <div class="page-header">
            <div class="header-meta">
                <h2>AI 大模型渠道与对话审计中心</h2>
                <p>统一配置通义千问 (Qwen)、DeepSeek、OpenAI 等兼容大模型渠道，指定全局默认模型，并审计全员 AI 对话与操作记录</p>
            </div>
            <div class="header-actions">
                <el-button icon="ChatDotRound" @click="router.push('/ai/workspace')">
                    前往 AI 工作台测试
                </el-button>
                <el-button type="primary" icon="Plus" @click="handleAddModel">
                    新增大模型渠道
                </el-button>
            </div>
        </div>

        <el-tabs v-model="activeTab" class="config-tabs">
            <!-- Tab 1: 大模型渠道与全局默认配置 -->
            <el-tab-pane label="大模型渠道配置 (Model Providers)" name="models">
                <div class="provider-presets-banner">
                    <span class="banner-title">支持主流 OpenAI 兼容协议大模型：</span>
                    <span class="preset-tag qwen">通义千问 Qwen (阿里云百炼)</span>
                    <span class="preset-tag deepseek">DeepSeek-V3 / R1</span>
                    <span class="preset-tag openai">OpenAI GPT-4o / GPT-4.1</span>
                    <span class="preset-tag custom">智谱 GLM / Moonshot / 本地 Ollama</span>
                </div>

                <el-card shadow="never" class="table-card">
                    <el-table v-loading="modelLoading" :data="modelList" stripe style="width: 100%;">
                        <el-table-column label="渠道名称" min-width="200">
                            <template #default="{ row }">
                                <div class="model-name-cell">
                                    <span class="m-title">{{ row.name }}</span>
                                </div>
                            </template>
                        </el-table-column>

                        <el-table-column label="供应商" width="115">
                            <template #default="{ row }">
                                <span class="provider-pill" :class="row.provider">
                                    {{ formatProviderLabel(row.provider) }}
                                </span>
                            </template>
                        </el-table-column>

                        <el-table-column prop="model_name" label="模型标识 (Model ID)" width="150">
                            <template #default="{ row }">
                                <code class="model-code">{{ row.model_name }}</code>
                            </template>
                        </el-table-column>

                        <el-table-column prop="base_url" label="API Base URL" min-width="230" show-overflow-tooltip />

                        <el-table-column label="API Key 状态" width="135">
                            <template #default="{ row }">
                                <el-tag v-if="row.has_api_key" type="primary" size="small" effect="light">
                                    {{ row.api_key_masked }}
                                </el-tag>
                                <el-tag v-else type="warning" size="small" effect="light">
                                    未配置密钥
                                </el-tag>
                            </template>
                        </el-table-column>

                        <el-table-column label="全局默认" width="125" align="center">
                            <template #default="{ row }">
                                <span v-if="row.is_default === 1" class="default-active-pill">
                                    <el-icon><StarFilled /></el-icon>
                                    <span>当前默认</span>
                                </span>
                                <button
                                    v-else
                                    type="button"
                                    class="set-default-btn"
                                    @click="handleSetDefault(row)"
                                >
                                    设为默认
                                </button>
                            </template>
                        </el-table-column>

                        <el-table-column label="状态" width="85" align="center">
                            <template #default="{ row }">
                                <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
                                    {{ row.status === 1 ? '启用' : '禁用' }}
                                </el-tag>
                            </template>
                        </el-table-column>

                        <el-table-column label="操作" width="230" align="center" fixed="right">
                            <template #default="{ row }">
                                <div class="ops-group">
                                    <button
                                        type="button"
                                        class="op-chip test"
                                        :disabled="testingId === row.id"
                                        @click="handleTestConnection(row)"
                                    >
                                        <el-icon class="is-loading" v-if="testingId === row.id"><Loading /></el-icon>
                                        <el-icon v-else><Connection /></el-icon>
                                        <span>测试连接</span>
                                    </button>
                                    <button
                                        type="button"
                                        class="op-chip edit"
                                        @click="handleEditModel(row)"
                                    >
                                        <el-icon><Key /></el-icon>
                                        <span>配置密钥</span>
                                    </button>
                                    <el-tooltip
                                        :content="row.is_default === 1 ? '当前默认模型不可直接删除' : '删除该渠道'"
                                        placement="top"
                                        :show-after="300"
                                    >
                                        <button
                                            type="button"
                                            class="op-chip delete"
                                            :disabled="row.is_default === 1"
                                            @click="handleDeleteModel(row)"
                                        >
                                            <el-icon><Delete /></el-icon>
                                            <span>删除</span>
                                        </button>
                                    </el-tooltip>
                                </div>
                            </template>
                        </el-table-column>
                    </el-table>
                </el-card>
            </el-tab-pane>

            <!-- Tab 2: 全员对话与操作审计 (Super Admin Audit) -->
            <el-tab-pane label="全员对话与操作审计 (Audit Logs)" name="audit">
                <el-card shadow="never" class="search-card">
                    <el-form :inline="true" :model="auditSearch" class="audit-search-form">
                        <el-form-item label="操作用户">
                            <el-input
                                v-model="auditSearch.username"
                                placeholder="按用户名搜索..."
                                clearable
                                style="width: 180px;"
                                @keyup.enter="fetchAuditSessions"
                            />
                        </el-form-item>
                        <el-form-item label="对话标题">
                            <el-input
                                v-model="auditSearch.keyword"
                                placeholder="搜索会话标题关键词..."
                                clearable
                                style="width: 220px;"
                                @keyup.enter="fetchAuditSessions"
                            />
                        </el-form-item>
                        <el-form-item label="入口场景">
                            <el-select
                                v-model="auditSearch.scene"
                                placeholder="全部场景"
                                clearable
                                style="width: 150px;"
                                @change="fetchAuditSessions"
                            >
                                <el-option label="AI 智能工作台" value="workspace" />
                                <el-option label="期刊页悬浮助手" value="copilot" />
                            </el-select>
                        </el-form-item>
                        <el-form-item>
                            <el-button type="primary" icon="Search" @click="fetchAuditSessions">查询</el-button>
                            <el-button icon="RefreshLeft" @click="handleResetAudit">重置</el-button>
                        </el-form-item>
                    </el-form>
                </el-card>

                <el-card shadow="never" class="table-card">
                    <el-table v-loading="auditLoading" :data="auditList" stripe style="width: 100%;">
                        <el-table-column prop="id" label="ID" width="70" align="center" />
                        <el-table-column prop="username" label="所属用户账号" width="150">
                            <template #default="{ row }">
                                <span style="font-weight: 600; color: #0f172a;">{{ row.username }}</span>
                            </template>
                        </el-table-column>
                        <el-table-column prop="title" label="会话标题" min-width="240" show-overflow-tooltip />
                        <el-table-column label="触发入口" width="140">
                            <template #default="{ row }">
                                <el-tag :type="row.scene === 'copilot' ? 'warning' : 'primary'" size="small" effect="light">
                                    {{ row.scene === 'copilot' ? '期刊页悬浮窗' : 'AI 智能工作台' }}
                                </el-tag>
                            </template>
                        </el-table-column>
                        <el-table-column prop="model_name" label="响应模型" width="160">
                            <template #default="{ row }">
                                <code class="model-code">{{ row.model_name || '-' }}</code>
                            </template>
                        </el-table-column>
                        <el-table-column prop="message_count" label="消息条数" width="100" align="center" />
                        <el-table-column label="最后活跃时间" width="175" align="center">
                            <template #default="{ row }">
                                {{ formatTime(row.updated_at) }}
                            </template>
                        </el-table-column>
                        <el-table-column label="审计操作" width="150" align="center" fixed="right">
                            <template #default="{ row }">
                                <el-button type="primary" link icon="View" @click="handleInspectSession(row)">
                                    查看完整对话
                                </el-button>
                            </template>
                        </el-table-column>
                    </el-table>

                    <div class="pagination-wrapper">
                        <el-pagination
                            v-model:current-page="auditPage"
                            v-model:page-size="auditPageSize"
                            :page-sizes="[15, 30, 50]"
                            layout="total, sizes, prev, pager, next"
                            :total="auditTotal"
                            @size-change="fetchAuditSessions"
                            @current-change="fetchAuditSessions"
                        />
                    </div>
                </el-card>
            </el-tab-pane>
        </el-tabs>

        <!-- 新增 / 编辑大模型渠道弹窗 -->
        <el-dialog
            v-model="dialogVisible"
            :title="formModel.id ? '编辑大模型渠道与密钥' : '新增大模型渠道'"
            width="640px"
        >
            <el-form ref="formRef" :model="formModel" :rules="formRules" label-width="115px">
                <el-form-item label="快速预设模板">
                    <div class="preset-buttons">
                        <el-button size="small" @click="applyPreset('qwen')">填入通义千问 Qwen</el-button>
                        <el-button size="small" @click="applyPreset('deepseek')">填入 DeepSeek-V3</el-button>
                        <el-button size="small" @click="applyPreset('openai')">填入 OpenAI GPT-4o</el-button>
                    </div>
                </el-form-item>

                <el-row :gutter="16">
                    <el-col :span="14">
                        <el-form-item label="渠道显示名称" prop="name">
                            <el-input v-model="formModel.name" placeholder="例如：通义千问 Qwen-Plus / Sub2Api" />
                        </el-form-item>
                    </el-col>
                    <el-col :span="10">
                        <el-form-item label="供应商类型" prop="provider">
                            <el-select v-model="formModel.provider" style="width: 100%;" @change="handleProviderChange">
                                <el-option label="通义千问 (Qwen)" value="qwen" />
                                <el-option label="深度求索 (DeepSeek)" value="deepseek" />
                                <el-option label="OpenAI" value="openai" />
                                <el-option label="自定义兼容接口 (Custom)" value="custom" />
                            </el-select>
                        </el-form-item>
                    </el-col>
                </el-row>

                <el-form-item label="API Base URL" prop="base_url">
                    <el-input
                        v-model="formModel.base_url"
                        placeholder="如 https://dashscope.aliyuncs.com/compatible-mode/v1"
                    />
                </el-form-item>

                <el-form-item label="API Key 密钥" prop="api_key">
                    <el-input
                        v-model="formModel.api_key"
                        type="password"
                        show-password
                        :placeholder="
                            formModel.id && formModel.has_api_key
                                ? `已配置 (${formModel.api_key_masked})，留空则保持原密钥不变`
                                : '请输入以 sk- 开头的 API 密钥'
                        "
                    />
                </el-form-item>

                <el-form-item label="模型标识" prop="model_name">
                    <div class="model-select-row">
                        <el-select
                            v-model="formModel.model_name"
                            filterable
                            allow-create
                            default-first-option
                            :loading="fetchingModels"
                            placeholder="点击下拉选择或搜索模型（也可直接输入自定义模型名）"
                            class="model-select-input"
                            @focus="handleModelSelectFocus"
                        >
                            <el-option
                                v-for="mName in availableModels"
                                :key="mName"
                                :label="mName"
                                :value="mName"
                            />
                        </el-select>
                        <el-button
                            type="primary"
                            plain
                            icon="Refresh"
                            :loading="fetchingModels"
                            @click="handleFetchRemoteModels(false)"
                        >
                            获取模型列表
                        </el-button>
                    </div>
                    <div v-if="modelsHintText" class="models-fetch-hint" :class="{ remote: isRemoteModels }">
                        <span>{{ modelsHintText }}</span>
                    </div>
                </el-form-item>

                <el-row :gutter="16">
                    <el-col :span="12">
                        <el-form-item label="采样温度">
                            <el-input-number
                                v-model="formModel.temperature"
                                :min="0.05"
                                :max="1.5"
                                :step="0.05"
                                style="width: 100%;"
                            />
                        </el-form-item>
                    </el-col>
                    <el-col :span="12">
                        <el-form-item label="Max Tokens">
                            <el-input-number
                                v-model="formModel.max_tokens"
                                :min="512"
                                :max="32768"
                                :step="512"
                                style="width: 100%;"
                            />
                        </el-form-item>
                    </el-col>
                </el-row>

                <el-form-item label="系统提示词">
                    <el-input
                        v-model="formModel.system_prompt"
                        type="textarea"
                        :rows="3"
                        placeholder="可选：附加的 Agent 角色设定提示词"
                    />
                </el-form-item>

                <el-row :gutter="16">
                    <el-col :span="12">
                        <el-form-item label="设为全局默认">
                            <el-switch v-model="formModel.is_default" :active-value="1" :inactive-value="0" />
                        </el-form-item>
                    </el-col>
                    <el-col :span="12">
                        <el-form-item label="启用状态">
                            <el-switch v-model="formModel.status" :active-value="1" :inactive-value="0" />
                        </el-form-item>
                    </el-col>
                </el-row>
            </el-form>
            <template #footer>
                <el-button @click="dialogVisible = false">取消</el-button>
                <el-button type="primary" :loading="saving" @click="handleSubmitModel">保存配置</el-button>
            </template>
        </el-dialog>

        <!-- 超级管理员查看用户完整对话审计抽屉 -->
        <el-drawer
            v-model="auditDrawerVisible"
            :title="`对话审计详情 · 用户 [${inspectingSession?.username || ''}] · ${inspectingSession?.title || ''}`"
            size="640px"
        >
            <div v-loading="auditMsgLoading" class="audit-drawer-body">
                <div v-if="auditMessages.length > 0" class="audit-msg-list">
                    <div
                        v-for="m in auditMessages"
                        :key="m.id"
                        class="audit-msg-item"
                        :class="m.role"
                    >
                        <div class="a-header">
                            <span class="a-role">{{ m.role === 'user' ? `用户 (${m.username})` : `AI 助手 (${m.model_name})` }}</span>
                            <span class="a-time">{{ formatTime(m.created_at) }}</span>
                        </div>
                        <div v-if="m.tool_steps && m.tool_steps.length > 0" class="a-steps">
                            <span v-for="(st, idx) in m.tool_steps" :key="idx" class="a-step-tag">
                                工具调用: {{ st.title }} ({{ st.summary }})
                            </span>
                        </div>
                        <div class="a-content" v-html="renderGFMMarkdown(m.content)"></div>
                        <AiActionCard
                            v-if="m.action_card"
                            :card="m.action_card"
                            :message-id="m.id"
                            :readonly="true"
                        />
                    </div>
                </div>
                <el-empty v-else description="该会话暂无消息记录" />
            </div>
        </el-drawer>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import AiActionCard from '@/components/AiActionCard/index.vue'
import { renderGFMMarkdown } from '@/utils/markdown'
import {
    reqAiModelList,
    reqSaveAiModel,
    reqSetDefaultAiModel,
    reqTestAiModelConnection,
    reqDeleteAiModel,
    reqFetchRemoteAiModels,
    reqAuditSessions,
    reqSessionMessages,
    type AIModelConfigItem,
    type AIChatSessionItem,
    type AIChatMessageItem,
} from '@/api/ai'

defineOptions({
    name: 'SystemAiConfig',
})

const router = useRouter()
const activeTab = ref('models')

// Tab 1: 模型列表状态
const modelLoading = ref(false)
const modelList = ref<AIModelConfigItem[]>([])
const testingId = ref<number | null>(null)

const dialogVisible = ref(false)
const saving = ref(false)
const formRef = ref()
const fetchingModels = ref(false)
const availableModels = ref<string[]>([])
const hasFetchedOnce = ref(false)
const isRemoteModels = ref(false)
const modelsHintText = ref('')
const formModel = reactive<Partial<AIModelConfigItem>>({
    id: 0,
    name: '',
    provider: 'qwen',
    base_url: 'https://dashscope.aliyuncs.com/compatible-mode/v1',
    api_key: '',
    api_key_masked: '',
    has_api_key: false,
    model_name: 'qwen-plus',
    system_prompt: '',
    temperature: 0.2,
    max_tokens: 4096,
    is_default: 0,
    status: 1,
    remark: '',
})

const formRules = {
    name: [{ required: true, message: '请输入渠道名称', trigger: 'blur' }],
    base_url: [{ required: true, message: '请输入 API Base URL', trigger: 'blur' }],
    model_name: [{ required: true, message: '请输入模型标识', trigger: 'blur' }],
}

// Tab 2: 对话审计状态
const auditLoading = ref(false)
const auditList = ref<AIChatSessionItem[]>([])
const auditTotal = ref(0)
const auditPage = ref(1)
const auditPageSize = ref(15)
const auditSearch = reactive({
    username: '',
    keyword: '',
    scene: '',
})

const auditDrawerVisible = ref(false)
const auditMsgLoading = ref(false)
const inspectingSession = ref<AIChatSessionItem | null>(null)
const auditMessages = ref<AIChatMessageItem[]>([])

const fetchModels = async () => {
    modelLoading.value = true
    try {
        const res = await reqAiModelList()
        if (res.code === 200) {
            modelList.value = res.data || []
        }
    } finally {
        modelLoading.value = false
    }
}

const applyPreset = (provider: 'qwen' | 'deepseek' | 'openai') => {
    if (provider === 'qwen') {
        formModel.name = '通义千问 Qwen-Plus (阿里云百炼)'
        formModel.provider = 'qwen'
        formModel.base_url = 'https://dashscope.aliyuncs.com/compatible-mode/v1'
        formModel.model_name = 'qwen-plus'
        availableModels.value = ['qwen-plus', 'qwen-max', 'qwen-turbo', 'qwen-long', 'qwq-32b']
    } else if (provider === 'deepseek') {
        formModel.name = 'DeepSeek-V3 (深度求索官方)'
        formModel.provider = 'deepseek'
        formModel.base_url = 'https://api.deepseek.com/v1'
        formModel.model_name = 'deepseek-chat'
        availableModels.value = ['deepseek-chat', 'deepseek-reasoner']
    } else if (provider === 'openai') {
        formModel.name = 'OpenAI GPT-4o'
        formModel.provider = 'openai'
        formModel.base_url = 'https://api.openai.com/v1'
        formModel.model_name = 'gpt-4o'
        availableModels.value = ['gpt-4o', 'gpt-4o-mini', 'gpt-4.1', 'gpt-4.1-mini', 'o3-mini']
    }
    hasFetchedOnce.value = false
    modelsHintText.value = ''
}

const handleProviderChange = () => {
    hasFetchedOnce.value = false
}

// 点击模型下拉框或点击【获取模型列表】按钮时向后端请求渠道 /models
const handleFetchRemoteModels = async (silent: boolean = false) => {
    if (fetchingModels.value) return
    fetchingModels.value = true
    try {
        const res = await reqFetchRemoteAiModels({
            id: formModel.id || 0,
            provider: formModel.provider || 'custom',
            base_url: formModel.base_url || '',
            api_key: formModel.api_key || '',
        })
        if (res.code === 200 && res.data) {
            const list = res.data.models || []
            if (formModel.model_name && !list.includes(formModel.model_name)) {
                availableModels.value = [formModel.model_name, ...list]
            } else {
                availableModels.value = list
            }
            if (!formModel.model_name && list.length > 0) {
                formModel.model_name = list[0]
            }
            isRemoteModels.value = !!res.data.from_remote
            modelsHintText.value = res.data.message || ''
            hasFetchedOnce.value = true
            if (!silent) {
                if (res.data.from_remote) {
                    ElMessage.success(res.data.message)
                } else {
                    ElMessage.info(res.data.message)
                }
            }
        }
    } finally {
        fetchingModels.value = false
    }
}

const handleModelSelectFocus = () => {
    if (!hasFetchedOnce.value) {
        handleFetchRemoteModels(true)
    }
}

const handleAddModel = () => {
    Object.assign(formModel, {
        id: 0,
        name: '通义千问 Qwen-Plus',
        provider: 'qwen',
        base_url: 'https://dashscope.aliyuncs.com/compatible-mode/v1',
        api_key: '',
        api_key_masked: '',
        has_api_key: false,
        model_name: 'qwen-plus',
        system_prompt: '你是专业的学术与商业期刊管理 AI 助手，精通期刊选刊推荐、征稿函信息提取、排期与报价维护。',
        temperature: 0.2,
        max_tokens: 4096,
        is_default: 0,
        status: 1,
        remark: '',
    })
    availableModels.value = ['qwen-plus', 'qwen-max', 'qwen-turbo', 'qwen-long']
    hasFetchedOnce.value = false
    isRemoteModels.value = false
    modelsHintText.value = ''
    dialogVisible.value = true
}

const handleEditModel = (row: AIModelConfigItem) => {
    Object.assign(formModel, {
        id: row.id,
        name: row.name,
        provider: row.provider,
        base_url: row.base_url,
        api_key: '',
        api_key_masked: row.api_key_masked,
        has_api_key: row.has_api_key,
        model_name: row.model_name,
        system_prompt: row.system_prompt,
        temperature: row.temperature,
        max_tokens: row.max_tokens,
        is_default: row.is_default,
        status: row.status,
        remark: row.remark,
    })
    availableModels.value = row.model_name ? [row.model_name] : []
    hasFetchedOnce.value = false
    isRemoteModels.value = false
    modelsHintText.value = ''
    dialogVisible.value = true
    // 若该渠道已保存密钥，打开弹窗时自动在后台拉取一次真实可用模型列表
    if (row.has_api_key) {
        handleFetchRemoteModels(true)
    }
}

const handleSubmitModel = async () => {
    await formRef.value?.validate(async (valid: boolean) => {
        if (!valid) return
        saving.value = true
        try {
            const res = await reqSaveAiModel(formModel)
            if (res.code === 200) {
                ElMessage.success('大模型配置已保存')
                dialogVisible.value = false
                fetchModels()
            } else {
                ElMessage.error(res.message || '保存失败')
            }
        } finally {
            saving.value = false
        }
    })
}

const handleSetDefault = async (row: AIModelConfigItem) => {
    const res = await reqSetDefaultAiModel(row.id)
    if (res.code === 200) {
        ElMessage.success(`已将「${row.name}」设为系统全局默认模型`)
        fetchModels()
    } else {
        ElMessage.error(res.message || '设置失败')
    }
}

const handleTestConnection = async (row: AIModelConfigItem) => {
    testingId.value = row.id
    try {
        const res = await reqTestAiModelConnection(row.id)
        if (res.code === 200) {
            if (res.data.success) {
                ElMessage.success(`连接成功！模型响应耗时 ${res.data.latency_ms}ms，回复：“${res.data.reply}”`)
            } else {
                ElMessage.warning(res.data.message)
            }
        }
    } finally {
        testingId.value = null
    }
}

const handleDeleteModel = (row: AIModelConfigItem) => {
    ElMessageBox.confirm(`确定要删除模型渠道「${row.name}」吗？`, '删除确认', {
        type: 'warning',
    })
        .then(async () => {
            const res = await reqDeleteAiModel(row.id)
            if (res.code === 200) {
                ElMessage.success('已删除该模型渠道')
                fetchModels()
            } else {
                ElMessage.error(res.message || '删除失败')
            }
        })
        .catch(() => {})
}

// Tab 2: 审计查询
const fetchAuditSessions = async () => {
    auditLoading.value = true
    try {
        const res = await reqAuditSessions({
            page_number: auditPage.value,
            page_size: auditPageSize.value,
            username: auditSearch.username,
            keyword: auditSearch.keyword,
            scene: auditSearch.scene,
        })
        if (res.code === 200) {
            auditList.value = res.data.items || []
            auditTotal.value = res.data.total || 0
        }
    } finally {
        auditLoading.value = false
    }
}

const handleResetAudit = () => {
    auditSearch.username = ''
    auditSearch.keyword = ''
    auditSearch.scene = ''
    auditPage.value = 1
    fetchAuditSessions()
}

const handleInspectSession = async (row: AIChatSessionItem) => {
    inspectingSession.value = row
    auditDrawerVisible.value = true
    auditMsgLoading.value = true
    try {
        const res = await reqSessionMessages(row.session_id)
        if (res.code === 200) {
            auditMessages.value = res.data || []
        }
    } finally {
        auditMsgLoading.value = false
    }
}

const formatProviderLabel = (p: string) => {
    const map: Record<string, string> = {
        qwen: '通义千问',
        deepseek: 'DeepSeek',
        openai: 'OpenAI',
        custom: '兼容自定义',
    }
    return map[p] || p
}

const formatTime = (ts?: number) => {
    if (!ts) return '-'
    return dayjs(ts * 1000).format('YYYY-MM-DD HH:mm:ss')
}

onMounted(() => {
    fetchModels()
    fetchAuditSessions()
})
</script>

<style scoped lang="scss">
.ai-config-container {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.page-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: #ffffff;
    padding: 20px 24px;
    border-radius: 12px;
    border: 1px solid #e2e8f0;

    .header-meta {
        h2 {
            margin: 0 0 6px;
            font-size: 18px;
            font-weight: 700;
            color: #0f172a;
        }
        p {
            margin: 0;
            font-size: 13px;
            color: #64748b;
        }
    }

    .header-actions {
        display: flex;
        gap: 10px;
    }
}

.provider-presets-banner {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 14px;
    padding: 12px 16px;
    border-radius: 10px;
    background: #f8fafc;
    border: 1px solid #e2e8f0;
    font-size: 13px;

    .banner-title {
        font-weight: 600;
        color: #334155;
    }

    .preset-tag {
        padding: 3px 10px;
        border-radius: 999px;
        font-size: 12px;
        font-weight: 600;

        &.qwen {
            background: #eff6ff;
            color: #1d4ed8;
        }
        &.deepseek {
            background: #f0f9ff;
            color: #0369a1;
        }
        &.openai {
            background: #ecfdf5;
            color: #047857;
        }
        &.custom {
            background: #f1f5f9;
            color: #475569;
        }
    }
}

.model-name-cell {
    display: flex;
    align-items: center;
    gap: 8px;

    .m-title {
        font-weight: 600;
        color: #0f172a;
    }
}

.provider-pill {
    display: inline-block;
    padding: 2px 8px;
    border-radius: 6px;
    font-size: 12px;
    font-weight: 600;
    background: #f1f5f9;
    color: #334155;

    &.qwen {
        background: #eff6ff;
        color: #2563eb;
    }
    &.deepseek {
        background: #f0f9ff;
        color: #0284c7;
    }
    &.openai {
        background: #ecfdf5;
        color: #059669;
    }
}

.model-code {
    padding: 2px 6px;
    border-radius: 4px;
    background: #f1f5f9;
    color: #0f172a;
    font-size: 12.5px;
}

.default-active-pill {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    height: 26px;
    padding: 0 10px;
    border-radius: 999px;
    background: #ecfdf5;
    color: #059669;
    border: 1px solid #a7f3d0;
    font-size: 12px;
    font-weight: 600;
}

.set-default-btn {
    height: 26px;
    padding: 0 10px;
    border-radius: 999px;
    background: #f8fafc;
    color: #64748b;
    border: 1px dashed #cbd5e1;
    font-size: 12px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.15s ease;

    &:hover {
        background: #eff6ff;
        color: #2563eb;
        border-color: #93c5fd;
        border-style: solid;
    }
}

.ops-group {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
    white-space: nowrap;
}

.op-chip {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    height: 27px;
    padding: 0 9px;
    border-radius: 6px;
    border: 1px solid transparent;
    font-size: 12px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.15s ease;
    white-space: nowrap;

    .el-icon {
        font-size: 13px;
    }

    &:disabled {
        opacity: 0.45;
        cursor: not-allowed;
    }

    &.test {
        color: #0284c7;
        background: #f0f9ff;
        border-color: #bae6fd;

        &:hover:not(:disabled) {
            background: #e0f2fe;
            border-color: #7dd3fc;
        }
    }

    &.edit {
        color: #2563eb;
        background: #eff6ff;
        border-color: #bfdbfe;

        &:hover:not(:disabled) {
            background: #dbeafe;
            border-color: #93c5fd;
        }
    }

    &.delete {
        color: #dc2626;
        background: #fef2f2;
        border-color: #fecaca;

        &:hover:not(:disabled) {
            background: #fee2e2;
            border-color: #fca5a5;
        }
    }
}

.search-card {
    margin-bottom: 14px;

    .audit-search-form {
        :deep(.el-form-item) {
            margin-bottom: 0;
        }
    }
}

.pagination-wrapper {
    margin-top: 16px;
    display: flex;
    justify-content: flex-end;
}

.preset-buttons {
    display: flex;
    gap: 8px;
}

.model-select-row {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 8px;

    .model-select-input {
        flex: 1;
    }
}

.models-fetch-hint {
    margin-top: 4px;
    font-size: 12px;
    line-height: 1.4;
    color: #64748b;

    &.remote {
        color: #059669;
        font-weight: 500;
    }
}

.audit-drawer-body {
    padding: 4px;

    .audit-msg-list {
        display: flex;
        flex-direction: column;
        gap: 14px;

        .audit-msg-item {
            padding: 12px 14px;
            border-radius: 10px;
            border: 1px solid #e2e8f0;
            background: #f8fafc;

            &.user {
                background: #eff6ff;
                border-color: #bfdbfe;
            }

            .a-header {
                display: flex;
                justify-content: space-between;
                margin-bottom: 6px;
                font-size: 12px;

                .a-role {
                    font-weight: 700;
                    color: #0f172a;
                }
                .a-time {
                    color: #64748b;
                }
            }

            .a-steps {
                margin-bottom: 6px;
                display: flex;
                flex-wrap: wrap;
                gap: 6px;

                .a-step-tag {
                    font-size: 11.5px;
                    padding: 2px 8px;
                    border-radius: 4px;
                    background: #dcfce7;
                    color: #166534;
                }
            }

            .a-content {
                font-size: 13.5px;
                line-height: 1.6;
                color: #1e293b;
                white-space: pre-wrap;
            }
        }
    }
}
</style>
