<template>
    <div class="ai-workspace-page">
        <!-- 左侧历史会话边栏 -->
        <aside class="session-sidebar">
            <div class="sidebar-top">
                <el-button type="primary" icon="Plus" class="new-chat-btn" @click="handleStartNewChat">
                    开启新对话
                </el-button>
                <el-input
                    v-model="searchKeyword"
                    placeholder="搜索历史对话..."
                    prefix-icon="Search"
                    clearable
                    size="small"
                    class="session-search"
                    @input="fetchSessions"
                />
            </div>

            <div v-loading="sessionLoading" class="session-list">
                <template v-if="sessionList.length > 0">
                    <div
                        v-for="s in sessionList"
                        :key="s.session_id"
                        class="session-item"
                        :class="{ active: currentSessionId === s.session_id }"
                        @click="handleSelectSession(s)"
                    >
                        <div class="session-icon">
                            <el-icon v-if="s.is_pinned === 1" style="color: #f59e0b;"><StarFilled /></el-icon>
                            <el-icon v-else><ChatDotRound /></el-icon>
                        </div>
                        <div class="session-info">
                            <div class="s-title">{{ s.title }}</div>
                            <div class="s-meta">
                                <span>{{ s.message_count }} 条消息</span>
                                <span>·</span>
                                <span>{{ formatRelativeTime(s.updated_at) }}</span>
                            </div>
                        </div>
                        <div class="session-ops" @click.stop>
                            <el-dropdown trigger="click" @command="(cmd: string) => handleSessionCommand(cmd, s)">
                                <span class="more-trigger">
                                    <el-icon><MoreFilled /></el-icon>
                                </span>
                                <template #dropdown>
                                    <el-dropdown-menu>
                                        <el-dropdown-item command="pin">
                                            {{ s.is_pinned === 1 ? '取消置顶' : '置顶会话' }}
                                        </el-dropdown-item>
                                        <el-dropdown-item command="rename">重命名标题</el-dropdown-item>
                                        <el-dropdown-item command="delete" divided style="color: #dc2626;">
                                            删除会话
                                        </el-dropdown-item>
                                    </el-dropdown-menu>
                                </template>
                            </el-dropdown>
                        </div>
                    </div>
                </template>
                <div v-else class="empty-sessions">
                    <p>暂无历史对话</p>
                    <span>输入自然语言即可管理与检索期刊</span>
                </div>
            </div>

            <!-- 底部当前全局默认模型状态栏 -->
            <div class="sidebar-footer">
                <div class="model-status-box">
                    <div class="status-left">
                        <span class="pulse-dot" :class="{ active: activeModel?.has_api_key }"></span>
                        <div class="model-meta">
                            <span class="m-label">当前全局默认引擎</span>
                            <span class="m-name">{{ activeModel?.name || '内置智能解析引擎' }}</span>
                        </div>
                    </div>
                    <el-tooltip
                        v-if="userStore.hasPermission('ai:config')"
                        content="前往配置大模型与对话审计"
                        placement="top"
                    >
                        <el-button
                            link
                            type="primary"
                            icon="Setting"
                            @click="router.push('/system/ai-config')"
                        />
                    </el-tooltip>
                </div>
            </div>
        </aside>

        <!-- 中间沉浸式 AI 对话主舞台 -->
        <section class="chat-stage">
            <!-- 舞台顶部状态栏 -->
            <header class="stage-header">
                <div class="header-left">
                    <h2 class="stage-title">{{ currentSessionTitle }}</h2>
                    <span class="model-badge">
                        <el-icon><Cpu /></el-icon>
                        {{ activeModel?.model_name || 'smart-agent-local' }}
                    </span>
                </div>
                <div class="header-right">
                    <el-button
                        v-if="canvasCard || previewJournal"
                        size="small"
                        :type="canvasOpen ? 'primary' : 'default'"
                        plain
                        icon="DataBoard"
                        @click="canvasOpen = !canvasOpen"
                    >
                        {{ canvasOpen ? '收起右侧画布' : '展开数据画布' }}
                    </el-button>
                    <el-button size="small" icon="Compass" @click="router.push('/periodical')">
                        前往期刊目录列表
                    </el-button>
                </div>
            </header>

            <!-- 消息滚动区 -->
            <div ref="msgScrollRef" v-loading="msgLoading" class="message-viewport">
                <!-- 空状态欢迎与快捷指令卡片 -->
                <div v-if="messages.length === 0" class="welcome-hero">
                    <div class="hero-badge">AI AGENT FOR PERIODICAL SYSTEM</div>
                    <h1>你好，{{ userStore.loginUser.name || '管理员' }} 👋</h1>
                    <p class="hero-sub">
                        我可以帮你用自然语言<strong>多维选刊检索、粘贴征稿函自动提取入库、修改排期与发表费、安全删除期刊</strong>。所有写操作均经过人机协同二次确认，并严格遵循当前账号的 RBAC 权限。
                    </p>

                    <div class="prompt-cards-grid">
                        <div
                            v-for="tpl in quickPrompts"
                            :key="tpl.title"
                            class="prompt-card"
                            @click="handleClickWelcomeCard(tpl)"
                        >
                            <div class="p-header">
                                <div class="p-header-left">
                                    <span class="p-tag" :class="tpl.tagClass">{{ tpl.tag }}</span>
                                    <span class="p-title">{{ tpl.title }}</span>
                                </div>
                                <span class="p-mode-hint" :class="tpl.actionMode">
                                    {{ tpl.actionMode === 'send' ? '⚡ 点击直接检索' : '✍️ 填入输入框编辑' }}
                                </span>
                            </div>
                            <p class="p-text">{{ tpl.desc }}</p>
                        </div>
                    </div>
                </div>

                <!-- 历史与实时消息流 -->
                <div v-else class="message-list">
                    <div
                        v-for="(msg, mIdx) in messages"
                        :key="msg.id"
                        class="message-row"
                        :class="msg.role === 'user' ? 'is-user' : 'is-assistant'"
                    >
                        <div class="avatar-col">
                            <div v-if="msg.role === 'user'" class="user-avatar">
                                {{ (msg.username || userStore.loginUser.name || 'U').slice(0, 1).toUpperCase() }}
                            </div>
                            <div v-else class="ai-avatar">AI</div>
                        </div>

                        <div class="bubble-col">
                            <div class="msg-meta">
                                <span class="sender-name">
                                    {{ msg.role === 'user' ? msg.username || '我' : '期刊智能助手' }}
                                </span>
                                <span v-if="msg.role === 'assistant' && msg.model_name" class="msg-model-tag">
                                    {{ msg.model_name }}
                                </span>
                                <span class="send-time">{{ formatFullTime(msg.created_at) }}</span>
                            </div>

                            <!-- 一体化响应容器：将工具轨迹、正文说明与交互确认/结果卡片融合在同一个整体内 -->
                            <div class="ai-turn-container" :class="{ 'is-user-turn': msg.role === 'user' }">
                                <!-- 可折叠 Agent 工具调用轨迹手风琴 (内嵌于容器顶部眉栏) -->
                                <div
                                    v-if="msg.tool_steps && msg.tool_steps.length > 0"
                                    class="agent-trace-accordion"
                                    :class="{ expanded: isTraceExpanded(msg.id) }"
                                >
                                    <div class="trace-header" @click="toggleTrace(msg.id)">
                                        <div class="trace-left">
                                            <el-icon
                                                v-if="msg.tool_steps[0].status === 'success'"
                                                class="trace-status-icon success"
                                            >
                                                <CircleCheckFilled />
                                            </el-icon>
                                            <el-icon v-else class="trace-status-icon error">
                                                <WarningFilled />
                                            </el-icon>
                                            <span class="trace-title">
                                                已调用 {{ msg.tool_steps.length }} 个工具 · {{ msg.tool_steps[0].title }}
                                            </span>
                                            <span class="trace-summary-badge">{{ msg.tool_steps[0].summary }}</span>
                                        </div>
                                        <div class="trace-right">
                                            <span v-if="msg.tool_steps[0].duration_ms" class="trace-ms">
                                                {{ msg.tool_steps[0].duration_ms }}ms
                                            </span>
                                            <el-icon class="trace-arrow">
                                                <ArrowDown />
                                            </el-icon>
                                        </div>
                                    </div>
                                    <div v-if="isTraceExpanded(msg.id)" class="trace-detail-list">
                                        <div v-for="(step, sIdx) in msg.tool_steps" :key="sIdx" class="trace-step-row">
                                            <div class="step-meta-line">
                                                <span class="fn-badge">fn: {{ step.tool_name }}</span>
                                                <span class="step-state" :class="step.status">
                                                    {{ step.status === 'success' ? '执行成功' : '已拦截/异常' }}
                                                </span>
                                                <span v-if="step.duration_ms" class="step-dur">
                                                    耗时 {{ step.duration_ms }}ms
                                                </span>
                                            </div>
                                            <div v-if="step.args_summary" class="step-kv">
                                                <strong>调用条件：</strong>{{ step.args_summary }}
                                            </div>
                                            <div class="step-kv">
                                                <strong>返回摘要：</strong>{{ step.summary }}
                                            </div>
                                        </div>
                                    </div>
                                </div>

                                <!-- 消息文本区 (支持完整 GFM 表格、标题、列表、代码高亮及打字机效果) -->
                                <div
                                    class="msg-bubble"
                                    v-html="renderGFMMarkdown(getDisplayContent(msg))"
                                ></div>

                                <!-- 内嵌式结构化操作/结果卡片 -->
                                <div v-if="msg.action_card" class="embedded-card-slot">
                                    <AiActionCard
                                        :card="msg.action_card"
                                        :message-id="msg.id"
                                        :show-open-canvas="true"
                                        @confirm-action="handleConfirmAction"
                                        @select-journal="handlePreviewJournal"
                                        @open-canvas="handleOpenCanvas"
                                        @quick-ask="useQuickPrompt"
                                        @start-modify="handleStartModify"
                                    />
                                </div>
                            </div>

                            <!-- Assistant 消息底部微工具栏 (复制 / 重新生成 / 打开画布) -->
                            <div v-if="msg.role === 'assistant'" class="msg-action-bar">
                                <button type="button" class="bar-btn" @click="handleCopyMessage(msg.content)">
                                    <el-icon><CopyDocument /></el-icon>
                                    <span>复制内容</span>
                                </button>
                                <button
                                    type="button"
                                    class="bar-btn"
                                    :disabled="sending"
                                    @click="handleRegenerate(mIdx)"
                                >
                                    <el-icon><RefreshRight /></el-icon>
                                    <span>重新生成</span>
                                </button>
                                <button
                                    v-if="msg.action_card && msg.action_card.card_type === 'query_result'"
                                    type="button"
                                    class="bar-btn primary-text"
                                    @click="handleOpenCanvas(msg.action_card)"
                                >
                                    <el-icon><DataBoard /></el-icon>
                                    <span>在右侧画布全屏浏览 ({{ msg.action_card.total_count || 0 }}本)</span>
                                </button>
                            </div>

                            <!-- 上下文智能追问建议胶囊 (仅在最后一条 AI 消息底部展示) -->
                            <div
                                v-if="
                                    msg.role === 'assistant' &&
                                    mIdx === messages.length - 1 &&
                                    !sending &&
                                    getSuggestions(msg).length > 0
                                "
                                class="followup-suggestions"
                            >
                                <span class="suggest-label">✨ 继续追问：</span>
                                <button
                                    v-for="sug in getSuggestions(msg)"
                                    :key="sug"
                                    type="button"
                                    class="suggest-pill"
                                    @click="useQuickPrompt(sug)"
                                >
                                    <span>{{ sug }}</span>
                                    <el-icon><Right /></el-icon>
                                </button>
                            </div>
                        </div>
                    </div>

                    <!-- 意图感知的阶段式思考进度卡片 -->
                    <div v-if="sending" class="message-row is-assistant">
                        <div class="avatar-col">
                            <div class="ai-avatar">AI</div>
                        </div>
                        <div class="bubble-col">
                            <div class="msg-meta">
                                <span class="sender">期刊智能助手</span>
                                <span class="m-pill">{{ activeModel?.model_name || 'smart-agent-local' }}</span>
                                <span class="thinking-live-pill">正在处理</span>
                            </div>
                            <div class="agent-thinking-card">
                                <div class="thinking-main-row">
                                    <div class="thinking-left">
                                        <span class="thinking-orb">
                                            <span class="orb-ring"></span>
                                            <span class="orb-core"></span>
                                        </span>
                                        <span class="thinking-stage-text">{{ currentThinkingDetail }}</span>
                                    </div>
                                    <span class="thinking-timer">{{ (thinkingElapsedMs / 1000).toFixed(1) }}s</span>
                                </div>
                                <div class="thinking-steps-track">
                                    <template v-for="(st, sIdx) in thinkingSteps" :key="st.label">
                                        <div
                                            class="t-step-node"
                                            :class="{
                                                done: sIdx < activeThinkingStepIndex,
                                                active: sIdx === activeThinkingStepIndex,
                                                pending: sIdx > activeThinkingStepIndex,
                                            }"
                                        >
                                            <span class="node-indicator">
                                                <template v-if="sIdx < activeThinkingStepIndex">✓</template>
                                                <template v-else-if="sIdx === activeThinkingStepIndex">
                                                    <span class="pulse-mini"></span>
                                                </template>
                                                <template v-else>{{ sIdx + 1 }}</template>
                                            </span>
                                            <span class="node-label">{{ st.label }}</span>
                                        </div>
                                        <span v-if="sIdx < thinkingSteps.length - 1" class="t-step-divider"></span>
                                    </template>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- 底部输入区 -->
            <footer class="composer-footer">
                <!-- 1. 当用户点击某本具体期刊的【改刊/修改此刊】时：显示该期刊的专属修改引导栏 -->
                <div v-if="modifyingTarget" class="target-modify-guide">
                    <div class="guide-top">
                        <div class="guide-title">
                            <span class="target-badge">🎯 已锁定修改目标</span>
                            <strong>{{ modifyingTarget.name }}</strong>
                            <span class="current-meta">
                                （当前排期：{{ modifyingTarget.publication_time || '待定' }} · 发表费：{{
                                    modifyingTarget.price ? `¥${modifyingTarget.price}` : '待询'
                                }}
                                <template v-if="modifyingTarget.periodical_page">
                                    · {{ modifyingTarget.periodical_page }}页
                                </template>
                                ）
                            </span>
                        </div>
                        <button type="button" class="cancel-target-btn" @click="modifyingTarget = null">
                            取消锁定 ×
                        </button>
                    </div>
                    <div class="guide-fields">
                        <span class="guide-hint">点击要修改的字段（填入输入框后填写新值，或直接输入自然语言）：</span>
                        <button
                            v-for="mf in modifyFieldTemplates"
                            :key="mf.label"
                            type="button"
                            class="mf-chip"
                            @click="appendModifyFieldDraft(mf.template)"
                        >
                            {{ mf.label }}
                        </button>
                        <button
                            type="button"
                            class="mf-chip direct-card"
                            @click="openDirectModifyCard(modifyingTarget)"
                        >
                            ⚡ 直接弹出可编辑修改卡片
                        </button>
                    </div>
                </div>

                <!-- 2. 当用户点击【➕ 录入新刊】时：展开新刊录入辅助面板（绝不自动发送假数据！） -->
                <div v-else-if="activeGuideMode === 'create'" class="target-modify-guide create-guide-theme">
                    <div class="guide-top">
                        <div class="guide-title">
                            <span class="target-badge emerald">📗 新刊录入辅助引导</span>
                            <strong>支持粘贴征稿函原文或按槽位填写刊物信息</strong>
                            <span class="current-meta">（请在下方输入框补充真实刊物信息后再点击【发送指令】）</span>
                        </div>
                        <div class="guide-top-ops">
                            <button type="button" class="tpl-action-btn" @click="setCreatePasteMode">
                                📄 切换粘贴征稿函模式
                            </button>
                            <button type="button" class="tpl-action-btn" @click="setCreateStructuredTemplate">
                                📋 填入完整属性模板
                            </button>
                            <button type="button" class="cancel-target-btn" @click="activeGuideMode = 'none'">
                                收起引导 ×
                            </button>
                        </div>
                    </div>
                    <div class="guide-fields">
                        <span class="guide-hint">点击追加新刊属性槽位：</span>
                        <button
                            v-for="cf in createFieldTemplates"
                            :key="cf.label"
                            type="button"
                            class="mf-chip emerald-chip"
                            @click="appendCreateFieldDraft(cf.template, cf.selectToken)"
                        >
                            {{ cf.label }}
                        </button>
                    </div>
                </div>

                <!-- 3. 当用户点击【✏️ 修改期刊排期与价格】时：展开通用改刊引导面板（绝不自动发送假数据！） -->
                <div v-else-if="activeGuideMode === 'modify'" class="target-modify-guide modify-general-theme">
                    <div class="guide-top">
                        <div class="guide-title">
                            <span class="target-badge amber">✏️ 期刊修改辅助引导</span>
                            <strong>请指定要修改的期刊名称与目标属性</strong>
                            <span class="current-meta">（可点选下方已检索出的刊物直接锁定，或在输入框内修改刊名与数值）</span>
                        </div>
                        <button type="button" class="cancel-target-btn" @click="activeGuideMode = 'none'">
                            收起引导 ×
                        </button>
                    </div>
                    <div v-if="recentCandidateJournals.length > 0" class="guide-fields">
                        <span class="guide-hint">点选当前对话中的刊物立即锁定改刊：</span>
                        <button
                            v-for="item in recentCandidateJournals"
                            :key="item.periodical_id"
                            type="button"
                            class="mf-chip"
                            @click="handleSelectCandidateToModify(item)"
                        >
                            🎯 {{ item.name }}
                        </button>
                    </div>
                    <div class="guide-fields">
                        <span class="guide-hint">点击追加要修改的字段槽位：</span>
                        <button
                            v-for="mf in modifyFieldTemplates"
                            :key="mf.label"
                            type="button"
                            class="mf-chip amber-chip"
                            @click="appendGeneralModifyFieldDraft(mf.template)"
                        >
                            {{ mf.label }}
                        </button>
                    </div>
                </div>

                <!-- 默认快捷操作胶囊栏 -->
                <div v-else class="quick-chips">
                    <span
                        v-for="chip in bottomChips"
                        :key="chip.label"
                        class="chip-item"
                        :class="{ 'is-draft-chip': chip.actionMode !== 'send' }"
                        @click="handleBottomChipClick(chip)"
                    >
                        {{ chip.label }}
                    </span>
                </div>

                <div
                    class="composer-box"
                    :class="{
                        'is-targeting': !!modifyingTarget || activeGuideMode === 'modify',
                        'is-creating': activeGuideMode === 'create',
                    }"
                >
                    <el-input
                        ref="composerInputRef"
                        v-model="inputPrompt"
                        type="textarea"
                        :autosize="{ minRows: 2, maxRows: 5 }"
                        :placeholder="
                            modifyingTarget
                                ? `请输入要把 ${modifyingTarget.name} 修改为什么内容，例如：“把 ${modifyingTarget.name} 的发表费改为 1200 元，出刊时间改为 2025-10-15”...`
                                : activeGuideMode === 'create'
                                  ? '请在此直接粘贴微信/文档中的期刊征稿介绍原文，或补全上方填入的新刊名称、级别、收录网站、排期与发表费后按 Enter 发送...'
                                  : '输入自然语言指令，例如：“查找知网收录的省级教育类期刊” 或直接粘贴一段期刊征稿介绍自动提取入库... (Enter 发送，Shift+Enter 换行)'
                        "
                        @keydown="handleKeydown"
                    />
                    <div class="composer-actions">
                        <span class="hint-text">所有新增 / 修改 / 删除操作均需您点击卡片二次确认后才会生效</span>
                        <el-button
                            type="primary"
                            icon="Position"
                            :loading="sending"
                            :disabled="!inputPrompt.trim()"
                            @click="handleSendMessage"
                        >
                            发送指令
                        </el-button>
                    </div>
                </div>
            </footer>
        </section>

        <!-- 右侧可展开 Artifact 数据与 26 项完整档案分屏画布 -->
        <transition name="canvas-slide">
            <aside v-if="canvasOpen" class="artifact-canvas-pane">
                <!-- 画布顶部切换与操作栏 -->
                <header class="canvas-header">
                    <div class="canvas-tabs">
                        <button
                            type="button"
                            class="c-tab"
                            :class="{ active: canvasTab === 'dataset' }"
                            @click="canvasTab = 'dataset'"
                        >
                            <el-icon><Grid /></el-icon>
                            <span>检索数据集 ({{ canvasTotal }})</span>
                        </button>
                        <button
                            v-if="previewJournal"
                            type="button"
                            class="c-tab"
                            :class="{ active: canvasTab === 'detail' }"
                            @click="canvasTab = 'detail'"
                        >
                            <el-icon><Document /></el-icon>
                            <span>单刊完整档案 · {{ previewJournal.name.slice(0, 10) }}</span>
                        </button>
                    </div>
                    <div class="canvas-ops">
                        <el-tooltip content="关闭右侧画布" placement="top">
                            <button type="button" class="close-canvas-btn" @click="canvasOpen = false">
                                <el-icon><Close /></el-icon>
                            </button>
                        </el-tooltip>
                    </div>
                </header>

                <!-- 画布内容 1：全量分页检索结果表格与多刊报价单生成器 -->
                <div v-if="canvasTab === 'dataset'" v-loading="canvasLoading" class="canvas-dataset-body">
                    <div class="dataset-toolbar">
                        <el-input
                            v-model="canvasFilterKeyword"
                            size="small"
                            clearable
                            prefix-icon="Search"
                            placeholder="在当前结果页按刊名/方向/收录快速过滤..."
                            class="canvas-search"
                        />
                        <el-button
                            size="small"
                            type="primary"
                            plain
                            :disabled="selectedCanvasItems.length === 0"
                            @click="handleCopyBatchQuote"
                        >
                            复制已选报价单 ({{ selectedCanvasItems.length }})
                        </el-button>
                    </div>

                    <div v-if="canvasCard?.description" class="dataset-condition-bar">
                        <span>{{ canvasCard.description }}</span>
                        <span class="total-badge">共 {{ canvasTotal }} 本</span>
                    </div>

                    <div class="dataset-table-scroll">
                        <table class="canvas-data-table">
                            <thead>
                                <tr>
                                    <th style="width: 34px; text-align: center;">
                                        <input
                                            type="checkbox"
                                            :checked="isAllCanvasSelected"
                                            @change="toggleSelectAllCanvas"
                                        />
                                    </th>
                                    <th>期刊名称 / 属性</th>
                                    <th style="width: 72px;">级别</th>
                                    <th>网站收录</th>
                                    <th style="width: 96px;">排期</th>
                                    <th style="width: 78px; text-align: right;">发表费</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr
                                    v-for="item in filteredCanvasItems"
                                    :key="item.periodical_id"
                                    :class="{ active: previewJournal?.periodical_id === item.periodical_id }"
                                    @click="handleSelectCanvasRow(item)"
                                >
                                    <td style="text-align: center;" @click.stop>
                                        <input
                                            type="checkbox"
                                            :checked="isCanvasItemSelected(item.periodical_id)"
                                            @change="toggleSelectCanvasItem(item)"
                                        />
                                    </td>
                                    <td>
                                        <div class="c-name-line">
                                            <span class="c-book-name">{{ item.name }}</span>
                                            <span v-if="item.is_warp" class="c-mini-tag emerald">全包</span>
                                            <span v-if="item.invoice_receipt" class="c-mini-tag sky">发票</span>
                                        </div>
                                        <div class="c-sub-dir">{{ item.directions || '综合方向' }}</div>
                                    </td>
                                    <td>
                                        <span v-if="item.level" class="c-level-pill">{{ item.level }}</span>
                                        <span v-else>-</span>
                                    </td>
                                    <td>
                                        <div class="c-col-wrap">
                                            <span
                                                v-for="col in (item.collections || []).slice(0, 3)"
                                                :key="col"
                                                class="c-col-tag"
                                            >
                                                {{ col }}
                                            </span>
                                        </div>
                                    </td>
                                    <td class="c-pub">
                                        {{ item.publication_time || '待定' }}
                                        <div v-if="item.periodical_page" class="c-page">{{ item.periodical_page }}页</div>
                                    </td>
                                    <td class="c-price">
                                        {{ item.price && item.price !== '0' ? `¥${item.price}` : '待询' }}
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>

                    <!-- 画布底部分页与智能联动 -->
                    <footer class="canvas-footer">
                        <span class="p-stat">第 {{ canvasPage }} / {{ canvasTotalPages }} 页</span>
                        <div class="p-controls">
                            <el-button
                                size="small"
                                :disabled="canvasPage <= 1 || canvasLoading"
                                @click="loadCanvasPage(canvasPage - 1)"
                            >
                                上一页
                            </el-button>
                            <el-button
                                size="small"
                                :disabled="canvasPage >= canvasTotalPages || canvasLoading"
                                @click="loadCanvasPage(canvasPage + 1)"
                            >
                                下一页
                            </el-button>
                        </div>
                    </footer>
                </div>

                <!-- 画布内容 2：单刊 26 项完整档案视图（对齐期刊管理页详情弹窗全部信息） -->
                <div
                    v-else-if="canvasTab === 'detail' && previewJournal"
                    v-loading="detailLoading"
                    class="canvas-detail-body"
                >
                    <!-- 区块 1：核心排期、报价与服务属性 -->
                    <div class="detail-hero-card">
                        <div class="d-top">
                            <div class="d-title-wrap">
                                <h3>{{ fullDetail.name || previewJournal.name }}</h3>
                                <span class="d-id-badge">ID #{{ fullDetail.periodical_id || previewJournal.periodical_id }}</span>
                            </div>
                            <span v-if="detailLevel" class="d-level">{{ detailLevel }}</span>
                        </div>
                        <div class="d-badges">
                            <span v-if="fullDetail.is_hidden" class="d-tag amber" style="font-weight: 700;">🙈 已隐藏（AI默认屏蔽）</span>
                            <span v-if="fullDetail.is_warp" class="d-tag emerald">支持全包</span>
                            <span v-else class="d-tag gray">非全包</span>
                            <span v-if="fullDetail.invoice_receipt" class="d-tag sky">可开杂志社发票</span>
                            <span v-else class="d-tag gray">不可开发票</span>
                            <span v-if="fullDetail.works" class="d-tag amber">可发作品</span>
                            <span v-for="c in detailCollections" :key="c" class="d-tag slate">{{ c }}</span>
                        </div>
                        <div class="d-price-row">
                            <div>
                                <span class="label">参考发表费：</span>
                                <strong class="price-val">
                                    {{ fullDetail.price && fullDetail.price !== '0' ? `¥${fullDetail.price}` : '待询' }}
                                </strong>
                            </div>
                            <div>
                                <span class="label">近期排期：</span>
                                <strong>{{ formatDateShort(fullDetail.publication_time) || '待定' }}</strong>
                            </div>
                            <div>
                                <span class="label">页码：</span>
                                <strong>{{ fullDetail.periodical_page ? `${fullDetail.periodical_page}页` : '-' }}</strong>
                            </div>
                        </div>
                    </div>

                    <!-- 区块 2：出版与学术指标（12 项完整属性） -->
                    <div class="detail-section-card">
                        <div class="sec-header">
                            <span class="sec-title">📊 出版与学术指标档案</span>
                            <span v-if="fullDetail.updated_at" class="sec-sub">
                                更新于 {{ formatFullTime(fullDetail.updated_at) }}
                            </span>
                        </div>
                        <div class="detail-kv-grid">
                            <div class="kv-box full">
                                <span class="k">收稿方向</span>
                                <span class="v">{{ detailDirections || '-' }}</span>
                            </div>
                            <div class="kv-box">
                                <span class="k">主管单位</span>
                                <span class="v">{{ fullDetail.competent_unit || '-' }}</span>
                            </div>
                            <div class="kv-box">
                                <span class="k">主办单位</span>
                                <span class="v">{{ fullDetail.organization || '-' }}</span>
                            </div>
                            <div class="kv-box">
                                <span class="k">国内刊号 (CN)</span>
                                <span class="v">{{ fullDetail.domestic_number || '-' }}</span>
                            </div>
                            <div class="kv-box">
                                <span class="k">国际刊号 (ISSN)</span>
                                <span class="v">{{ fullDetail.international_number || '-' }}</span>
                            </div>
                            <div class="kv-box">
                                <span class="k">出版周期</span>
                                <span class="v">{{ detailCycle || '-' }}</span>
                            </div>
                            <div class="kv-box">
                                <span class="k">创刊时间</span>
                                <span class="v">{{ fullDetail.founding_time ? `${fullDetail.founding_time}年` : '-' }}</span>
                            </div>
                            <div class="kv-box">
                                <span class="k">审核周期</span>
                                <span class="v">{{ fullDetail.audit_time || '-' }}</span>
                            </div>
                            <div class="kv-box">
                                <span class="k">引用率</span>
                                <span class="v">{{ fullDetail.citation_rate || '-' }}</span>
                            </div>
                            <div class="kv-box">
                                <span class="k">复合影响因子</span>
                                <span class="v">{{ fullDetail.composite_influence_factor || '-' }}</span>
                            </div>
                            <div class="kv-box">
                                <span class="k">综合影响因子</span>
                                <span class="v">{{ fullDetail.synthetic_influence_factor || '-' }}</span>
                            </div>
                        </div>
                    </div>

                    <!-- 区块 3：征稿规范、栏目与样刊说明（6 大长文本字段，均带一键复制） -->
                    <div class="detail-section-card">
                        <div class="sec-header">
                            <span class="sec-title">📑 征稿要求与栏目规范</span>
                        </div>
                        <div class="rich-fields-list">
                            <div
                                v-for="rf in richTextFieldsGroup1"
                                :key="rf.key"
                                class="rich-field-item"
                            >
                                <div class="rf-top">
                                    <span class="rf-label">{{ rf.label }}</span>
                                    <button
                                        v-if="fullDetail[rf.key]"
                                        type="button"
                                        class="rf-copy-btn"
                                        @click="handleCopyFieldText(rf.label, fullDetail[rf.key])"
                                    >
                                        复制
                                    </button>
                                </div>
                                <div class="rf-val">{{ fullDetail[rf.key] || '未填写' }}</div>
                            </div>
                        </div>
                    </div>

                    <!-- 区块 4：发表与内部流转流程（3 大流程字段） -->
                    <div class="detail-section-card">
                        <div class="sec-header">
                            <span class="sec-title">🔄 发表与内部流转流程</span>
                        </div>
                        <div class="rich-fields-list">
                            <div
                                v-for="rf in richTextFieldsGroup2"
                                :key="rf.key"
                                class="rich-field-item"
                            >
                                <div class="rf-top">
                                    <span class="rf-label">{{ rf.label }}</span>
                                    <button
                                        v-if="fullDetail[rf.key]"
                                        type="button"
                                        class="rf-copy-btn"
                                        @click="handleCopyFieldText(rf.label, fullDetail[rf.key])"
                                    >
                                        复制
                                    </button>
                                </div>
                                <div class="rf-val">{{ fullDetail[rf.key] || '未填写' }}</div>
                            </div>
                        </div>
                    </div>

                    <!-- 底部操作栏 -->
                    <div class="detail-actions">
                        <el-button size="small" icon="Back" @click="canvasTab = 'dataset'">
                            返回结果列表
                        </el-button>
                        <el-button
                            size="small"
                            type="primary"
                            plain
                            icon="EditPen"
                            @click="handleStartModify(previewJournal)"
                        >
                            定制修改此刊
                        </el-button>
                        <el-button
                            size="small"
                            type="primary"
                            icon="Edit"
                            @click="openDirectModifyCard(previewJournal)"
                        >
                            生成可编辑修改单
                        </el-button>
                    </div>
                </div>
            </aside>
        </transition>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, nextTick, onMounted, onActivated, onBeforeUnmount } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import useUserStore from '@/store/modules/user'
import AiActionCard from '@/components/AiActionCard/index.vue'
import { renderGFMMarkdown } from '@/utils/markdown'
import { copyTextToClipboard } from '@/utils/clipboard'
import { reqPeriodicalList, reqPeriodicalDetail } from '@/api/periodical'
import {
    reqActiveDefaultModel,
    reqUserSessions,
    reqSessionMessages,
    reqSendAiMessage,
    reqConfirmAiAction,
    reqUpdateSession,
    reqDeleteSession,
    type AIModelConfigItem,
    type AIChatSessionItem,
    type AIChatMessageItem,
    type ActionCardData,
    type PeriodicalBriefCard,
} from '@/api/ai'

defineOptions({
    name: 'AiWorkspace',
})

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const activeModel = ref<AIModelConfigItem | null>(null)
const sessionList = ref<AIChatSessionItem[]>([])
const currentSessionId = ref<string>('')
const searchKeyword = ref('')
const sessionLoading = ref(false)

const messages = ref<AIChatMessageItem[]>([])
const msgLoading = ref(false)
const sending = ref(false)
const inputPrompt = ref('')
const msgScrollRef = ref<HTMLElement | null>(null)
const composerInputRef = ref<any>(null)

// 当前锁定的待修改期刊对象（点击【改刊】或【定制修改此刊】时激活，绝不自动发送假指令）
const modifyingTarget = ref<PeriodicalBriefCard | null>(null)
// 当前底部辅助引导面板模式：'none' | 'create' (新刊录入引导) | 'modify' (通用修改引导)
const activeGuideMode = ref<'none' | 'create' | 'modify'>('none')

const modifyFieldTemplates = [
    { label: '📅 改出刊时间', template: '出刊时间修改为 ' },
    { label: '💰 改发表费', template: '发表费修改为 ' },
    { label: '📄 改期刊页码', template: '页码修改为 ' },
    { label: '🙈 设为隐藏', template: '设为隐藏' },
    { label: '👁️ 取消隐藏', template: '取消隐藏' },
    { label: '🏷️ 改期刊级别', template: '期刊级别修改为 ' },
    { label: '📆 改出版周期', template: '出版周期修改为 ' },
    { label: '⏱️ 改审核周期', template: '审核周期修改为 ' },
    { label: '📑 改栏目设置', template: '栏目设置修改为：' },
    { label: '⚠️ 改注意事项', template: '注意事项修改为：' },
]

const createFieldTemplates = [
    { label: '🏷️ + 期刊级别与出版周期', template: '级别：省级，出版周期：月刊', selectToken: '省级' },
    { label: '🌐 + 网站收录', template: '网站收录：知网、万方', selectToken: '知网、万方' },
    { label: '🧭 + 收稿方向', template: '收稿方向：综合类', selectToken: '综合类' },
    { label: '📅 + 排期与发表费', template: '出刊时间：2025-10-01，发表费：1800元，120页', selectToken: '1800' },
    { label: '🏛️ + 主管与主办单位', template: '主管单位：请填写主管单位，主办单位：请填写主办单位', selectToken: '请填写主管单位' },
    { label: '🔢 + 国内/国际刊号', template: '国内刊号：11-xxxx/G，国际刊号：1000-xxxx', selectToken: '11-xxxx/G' },
    { label: '✅ + 全包与发票', template: '支持全包，可开杂志社发票', selectToken: '' },
]

// 工具执行轨迹折叠状态与渐进式打字机状态
const expandedTraces = reactive<Record<number, boolean>>({})
const typingMap = reactive<Record<number, string>>({})
let typingTimer: ReturnType<typeof setInterval> | null = null

// 思考阶段计时器与意图感知步骤状态
const thinkingElapsedMs = ref(0)
const lastSentPrompt = ref('')
let thinkingTimer: ReturnType<typeof setInterval> | null = null

interface ThinkingStepItem {
    label: string
    detail: string
}

const thinkingSteps = computed<ThinkingStepItem[]>(() => {
    const p = lastSentPrompt.value.trim()
    const allBookMatches = Array.from(p.matchAll(/《([^》]+)》/g))
        .map((m) => m[1]?.trim())
        .filter(Boolean)
    const uniqueBooks = Array.from(new Set(allBookMatches))
    const bookLabel =
        uniqueBooks.length > 1
            ? `《${uniqueBooks[0]}》等 ${uniqueBooks.length} 本期刊`
            : uniqueBooks.length === 1
              ? `《${uniqueBooks[0]}》`
              : ''

    if (/(新增|录入|添加|入库|提取.*字段)/.test(p)) {
        return [
            {
                label: '解析征稿要素',
                detail: `正在从文本中提取${bookLabel || '新期刊'}的刊名、刊号、排期与费用信息...`,
            },
            {
                label: '映射分类字典',
                detail: '正在自动匹配系统期刊级别、网站收录与收稿方向字典...',
            },
            {
                label: '生成入库确认单',
                detail: `正在构建可编辑的${bookLabel || '新刊'}入库确认卡片供您核对...`,
            },
        ]
    }

    if (/(修改确认单|把.*改|修改|调整|变更|排期.*改|发表费.*改)/.test(p)) {
        return [
            {
                label: '定位目标刊物',
                detail: `正在期刊库中精准定位${bookLabel || '目标刊物'}的当前库存档案...`,
            },
            {
                label: '计算字段对比 (Diff)',
                detail: '正在比对修改前后字段差异并校验数据格式...',
            },
            {
                label: '生成修改确认单',
                detail: `正在构建${bookLabel || '期刊'}的可交互修改预览卡片...`,
            },
        ]
    }

    if (/(删除|移除)/.test(p)) {
        return [
            {
                label: '核验库存档案',
                detail: `正在查询并核实${bookLabel || '目标期刊'}的库存记录...`,
            },
            {
                label: '启用安全护栏',
                detail: '正在构建高危删除二次确认卡片，等待您核对确认...',
            },
        ]
    }

    if (/(报价方案|推荐文案|推荐语|完整详情)/.test(p)) {
        return [
            {
                label: '调取刊物档案',
                detail: `正在读取${bookLabel || '相关期刊'}的排期、版面与收录参数...`,
            },
            {
                label: '生成专业内容',
                detail: '正在由大模型整理结构化内容与投审稿建议...',
            },
            {
                label: '润色排版输出',
                detail: '正在优化排版格式，即将为您呈现结果...',
            },
        ]
    }

    return [
        {
            label: '解析选刊条件',
            detail: '正在识别您的选刊要求与多维分类约束...',
        },
        {
            label: '检索期刊数据库',
            detail: '正在查询匹配的刊物排期、收录网站与参考发表费...',
        },
        {
            label: '提炼选刊洞察',
            detail: '正在由大模型汇总代表刊物亮点并生成交互数据卡片...',
        },
    ]
})

const activeThinkingStepIndex = computed(() => {
    const sec = thinkingElapsedMs.value / 1000
    const maxIdx = thinkingSteps.value.length - 1
    if (sec < 1.0) return 0
    if (sec < 2.6) return Math.min(1, maxIdx)
    return maxIdx
})

const currentThinkingDetail = computed(() => {
    const step = thinkingSteps.value[activeThinkingStepIndex.value]
    return step ? step.detail : '正在处理您的指令...'
})

// 右侧 Artifact 分屏画布状态与单刊 26 项完整详情状态
const canvasOpen = ref(false)
const canvasTab = ref<'dataset' | 'detail'>('dataset')
const canvasCard = ref<ActionCardData | null>(null)
const canvasItems = ref<PeriodicalBriefCard[]>([])
const canvasTotal = ref(0)
const canvasPage = ref(1)
const canvasPageSize = 15
const canvasLoading = ref(false)
const canvasFilterKeyword = ref('')
const selectedCanvasItems = ref<PeriodicalBriefCard[]>([])
const previewJournal = ref<PeriodicalBriefCard | null>(null)
const detailLoading = ref(false)
const fullDetail = ref<Record<string, any>>({})

const richTextFieldsGroup1 = [
    { key: 'periodical_description', label: '期刊描述' },
    { key: 'column_setting', label: '栏目设置' },
    { key: 'attention_matter', label: '注意事项' },
    { key: 'author_info', label: '作者数量与样刊相关' },
    { key: 'periodical_final_proof', label: '期刊清样' },
    { key: 'article_naming', label: '文章命名' },
]

const richTextFieldsGroup2 = [
    { key: 'publication_process', label: '发表流程' },
    { key: 'all_inclusive_process', label: '全包流程' },
    { key: 'internal_processes', label: '内部流程' },
]

const detailDirections = computed(() => {
    if (Array.isArray(fullDetail.value.categorizes)) {
        const arr = fullDetail.value.categorizes.filter((c: any) => c.type === 1).map((c: any) => c.name)
        if (arr.length > 0) return arr.join('、')
    }
    return previewJournal.value?.directions || ''
})

const detailLevel = computed(() => {
    if (Array.isArray(fullDetail.value.categorizes)) {
        const found = fullDetail.value.categorizes.find((c: any) => c.type === 2)
        if (found) return found.name
    }
    return previewJournal.value?.level || ''
})

const detailCollections = computed<string[]>(() => {
    if (Array.isArray(fullDetail.value.categorizes)) {
        const arr = fullDetail.value.categorizes.filter((c: any) => c.type === 3).map((c: any) => c.name)
        if (arr.length > 0) return arr
    }
    return previewJournal.value?.collections || []
})

const detailCycle = computed(() => {
    if (Array.isArray(fullDetail.value.categorizes)) {
        const found = fullDetail.value.categorizes.find((c: any) => c.type === 4)
        if (found) return found.name
    }
    return previewJournal.value?.cycle || ''
})

const canvasTotalPages = computed(() => Math.max(1, Math.ceil(canvasTotal.value / canvasPageSize)))

const filteredCanvasItems = computed(() => {
    const kw = canvasFilterKeyword.value.trim().toLowerCase()
    if (!kw) return canvasItems.value
    return canvasItems.value.filter(
        (it) =>
            it.name.toLowerCase().includes(kw) ||
            (it.directions || '').toLowerCase().includes(kw) ||
            (it.collections || []).some((c) => c.toLowerCase().includes(kw)) ||
            (it.level || '').toLowerCase().includes(kw),
    )
})

const isAllCanvasSelected = computed(() => {
    if (filteredCanvasItems.value.length === 0) return false
    return filteredCanvasItems.value.every((it) =>
        selectedCanvasItems.value.some((s) => s.periodical_id === it.periodical_id),
    )
})

const currentSessionTitle = computed(() => {
    if (!currentSessionId.value) return '新建智能期刊对话'
    const found = sessionList.value.find((s) => s.session_id === currentSessionId.value)
    return found ? found.title : '智能期刊对话'
})

// 提取当前会话最近检索到的候选期刊（最多6本，供通用修改引导栏一键点选锁定）
const recentCandidateJournals = computed<PeriodicalBriefCard[]>(() => {
    if (canvasItems.value.length > 0) {
        return canvasItems.value.slice(0, 6)
    }
    for (let i = messages.value.length - 1; i >= 0; i--) {
        const items = messages.value[i].action_card?.items
        if (items && items.length > 0) {
            return items.slice(0, 6)
        }
    }
    return []
})

const quickPrompts = [
    {
        tag: '智能检索',
        tagClass: 'blue',
        title: '多条件精准选刊推荐',
        desc: '帮我检索知网收录、全包、可开杂志社发票的省级期刊，按价格从低到高排序。',
        prompt: '帮我检索知网收录、全包、可开杂志社发票的省级期刊，按价格从低到高排序。',
        actionMode: 'send' as const,
        selectToken: '',
    },
    {
        tag: '文本入库',
        tagClass: 'emerald',
        title: '征稿函一键提取新增',
        desc: '直接粘贴微信/文档中的期刊征稿介绍原文，或按模板填入刊名、收录与费用自动提取入库。',
        prompt: '新增一本期刊《在此输入期刊名称》，级别：省级，出版周期：月刊，网站收录：知网，收稿方向：综合类，出刊时间：2025-10-01，发表费：1800元，页码：120页，支持全包，可开杂志社发票',
        actionMode: 'create_guide' as const,
        selectToken: '在此输入期刊名称',
    },
    {
        tag: '排期调价',
        tagClass: 'amber',
        title: '自然语言修改期刊信息',
        desc: '指定目标期刊名称，快速修改出刊排期、发表费、版面页码、栏目设置或注意事项。',
        prompt: '把《在此输入期刊名称》的出刊时间修改为 2025-10-01，发表费调整为 1600 元',
        actionMode: 'modify_guide' as const,
        selectToken: '在此输入期刊名称',
    },
    {
        tag: '安全管控',
        tagClass: 'rose',
        title: '带二次确认的期刊删除',
        desc: '输入要删除的期刊名称，系统将先定位档案并生成高危删除确认单供您二次核对。',
        prompt: '帮我删除《在此输入要删除的期刊名称》这本期刊',
        actionMode: 'draft' as const,
        selectToken: '在此输入要删除的期刊名称',
    },
]

const bottomChips = [
    { label: '🔍 查国家级期刊', prompt: '帮我查询所有国家级期刊', actionMode: 'send' as const },
    { label: '📚 查知网全包刊物', prompt: '帮我找一下知网收录且支持全包的期刊', actionMode: 'send' as const },
    { label: '💰 查低价刊物', prompt: '帮我查找收稿方向为低价刊物的期刊', actionMode: 'send' as const },
    { label: '➕ 录入新刊', prompt: '', actionMode: 'create_guide' as const },
    { label: '✏️ 修改期刊排期与价格', prompt: '', actionMode: 'modify_guide' as const },
]

const isTraceExpanded = (msgId: number) => Boolean(expandedTraces[msgId])
const toggleTrace = (msgId: number) => {
    expandedTraces[msgId] = !expandedTraces[msgId]
}

const getDisplayContent = (msg: AIChatMessageItem) => {
    if (typingMap[msg.id] !== undefined) {
        return typingMap[msg.id]
    }
    return msg.content
}

const startTypewriter = (msg: AIChatMessageItem) => {
    const full = msg.content || ''
    if (full.length <= 24) return
    if (typingTimer) clearInterval(typingTimer)
    let cursor = 18
    typingMap[msg.id] = full.slice(0, cursor)
    const stepSize = Math.max(12, Math.ceil(full.length / 28))
    typingTimer = setInterval(() => {
        cursor += stepSize
        if (cursor >= full.length) {
            delete typingMap[msg.id]
            if (typingTimer) {
                clearInterval(typingTimer)
                typingTimer = null
            }
            scrollToBottom()
        } else {
            typingMap[msg.id] = full.slice(0, cursor)
            scrollToBottom()
        }
    }, 24)
}

const getSuggestions = (msg: AIChatMessageItem): string[] => {
    if (msg.action_card?.suggestions && msg.action_card.suggestions.length > 0) {
        return msg.action_card.suggestions
    }
    if (msg.action_card?.card_type === 'query_result') {
        return [
            '筛选其中知网收录且支持全包的期刊',
            '按发表费从低到高排序推荐高性价比刊物',
            '基于前3本期刊生成一份对外客户选刊报价方案',
        ]
    }
    return []
}

const fetchActiveModel = async () => {
    try {
        const res = await reqActiveDefaultModel()
        if (res.code === 200) {
            activeModel.value = res.data
        }
    } catch {
        // ignore
    }
}

const fetchSessions = async () => {
    sessionLoading.value = true
    try {
        const res = await reqUserSessions(searchKeyword.value)
        if (res.code === 200) {
            sessionList.value = res.data || []
        }
    } finally {
        sessionLoading.value = false
    }
}

const loadMessages = async (sessionId: string) => {
    if (!sessionId) {
        messages.value = []
        return
    }
    msgLoading.value = true
    try {
        const res = await reqSessionMessages(sessionId)
        if (res.code === 200) {
            messages.value = res.data || []
            const lastQueryMsg = [...messages.value]
                .reverse()
                .find((m) => m.action_card?.card_type === 'query_result')
            if (lastQueryMsg?.action_card) {
                canvasCard.value = lastQueryMsg.action_card
                canvasItems.value = lastQueryMsg.action_card.items || []
                canvasTotal.value = lastQueryMsg.action_card.total_count || canvasItems.value.length
            }
            scrollToBottom()
        }
    } finally {
        msgLoading.value = false
    }
}

const handleSelectSession = (s: AIChatSessionItem) => {
    currentSessionId.value = s.session_id
    modifyingTarget.value = null
    activeGuideMode.value = 'none'
    loadMessages(s.session_id)
}

const handleStartNewChat = () => {
    currentSessionId.value = ''
    messages.value = []
    modifyingTarget.value = null
    activeGuideMode.value = 'none'
    canvasOpen.value = false
}

const useQuickPrompt = (text: string) => {
    inputPrompt.value = text
    handleSendMessage()
}

// 聚焦底部输入框并自动高亮选中待填占位词（方便用户直接打字替换，绝不自动发送假数据）
const focusAndSelectToken = (token?: string) => {
    nextTick(() => {
        const el: HTMLTextAreaElement | undefined =
            composerInputRef.value?.textarea || composerInputRef.value?.$el?.querySelector('textarea')
        if (!el) {
            composerInputRef.value?.focus?.()
            return
        }
        el.focus()
        if (token) {
            const idx = inputPrompt.value.indexOf(token)
            if (idx !== -1) {
                el.setSelectionRange(idx, idx + token.length)
                return
            }
        }
        const len = inputPrompt.value.length
        el.setSelectionRange(len, len)
    })
}

// 开启【新刊录入辅助引导】模式（绝不自动发送！）
const openCreateGuide = (withFullTemplate = false) => {
    modifyingTarget.value = null
    activeGuideMode.value = 'create'
    if (withFullTemplate) {
        setCreateStructuredTemplate()
    } else if (!inputPrompt.value.trim() || !inputPrompt.value.includes('新增')) {
        inputPrompt.value = '新增一本期刊《在此输入期刊名称》，'
        focusAndSelectToken('在此输入期刊名称')
    } else {
        focusAndSelectToken()
    }
}

const setCreatePasteMode = () => {
    inputPrompt.value = '请帮我从以下征稿介绍中提取字段并录入新刊：\n【请在此处直接粘贴期刊征稿介绍原文】'
    focusAndSelectToken('【请在此处直接粘贴期刊征稿介绍原文】')
}

const setCreateStructuredTemplate = () => {
    inputPrompt.value =
        '新增一本期刊《在此输入期刊名称》，级别：省级，出版周期：月刊，网站收录：知网，收稿方向：综合类，出刊时间：2025-10-01，发表费：1800元，页码：120页，支持全包，可开杂志社发票'
    focusAndSelectToken('在此输入期刊名称')
}

const appendCreateFieldDraft = (fieldTpl: string, selectToken?: string) => {
    const cur = inputPrompt.value.trim()
    if (!cur) {
        inputPrompt.value = `新增一本期刊《在此输入期刊名称》，${fieldTpl}`
        focusAndSelectToken(selectToken || '在此输入期刊名称')
        return
    }
    const sep = cur.endsWith('，') || cur.endsWith(',') || cur.endsWith('\n') ? '' : '，'
    inputPrompt.value = `${cur}${sep}${fieldTpl}`
    focusAndSelectToken(selectToken)
}

// 开启【期刊修改辅助引导】模式（绝不自动发送！）
const openGeneralModifyGuide = () => {
    modifyingTarget.value = null
    activeGuideMode.value = 'modify'
    if (!inputPrompt.value.trim() || !inputPrompt.value.includes('把')) {
        inputPrompt.value = '把《在此输入期刊名称》的出刊时间修改为 2025-10-01，发表费调整为 1600 元'
        focusAndSelectToken('在此输入期刊名称')
    } else {
        focusAndSelectToken()
    }
}

const handleSelectCandidateToModify = (item: PeriodicalBriefCard) => {
    activeGuideMode.value = 'none'
    handleStartModify(item)
}

const appendGeneralModifyFieldDraft = (fieldTpl: string) => {
    const cur = inputPrompt.value.trim()
    if (!cur) {
        inputPrompt.value = `把《在此输入期刊名称》的${fieldTpl}`
        focusAndSelectToken('在此输入期刊名称')
        return
    }
    const sep = cur.endsWith('的') || cur.endsWith('，') || cur.endsWith(',') ? '' : '，'
    inputPrompt.value = `${cur}${sep}${fieldTpl}`
    focusAndSelectToken()
}

const handleBottomChipClick = (chip: { label: string; prompt: string; actionMode: string }) => {
    if (chip.actionMode === 'send') {
        useQuickPrompt(chip.prompt)
    } else if (chip.actionMode === 'create_guide') {
        openCreateGuide(false)
    } else if (chip.actionMode === 'modify_guide') {
        openGeneralModifyGuide()
    } else {
        inputPrompt.value = chip.prompt
        focusAndSelectToken()
    }
}

const handleClickWelcomeCard = (tpl: {
    prompt: string
    actionMode: 'send' | 'create_guide' | 'modify_guide' | 'draft'
    selectToken?: string
}) => {
    if (tpl.actionMode === 'send') {
        useQuickPrompt(tpl.prompt)
    } else if (tpl.actionMode === 'create_guide') {
        openCreateGuide(true)
    } else if (tpl.actionMode === 'modify_guide') {
        openGeneralModifyGuide()
    } else {
        activeGuideMode.value = 'none'
        modifyingTarget.value = null
        inputPrompt.value = tpl.prompt
        focusAndSelectToken(tpl.selectToken)
    }
}

// 点击【改刊 / 定制修改此刊】时：锁定目标期刊并展示字段引导，绝不自动编造假指令发送！
const handleStartModify = (item: PeriodicalBriefCard) => {
    activeGuideMode.value = 'none'
    modifyingTarget.value = item
    if (!inputPrompt.value.includes(item.name)) {
        inputPrompt.value = `把${item.name}的`
    }
    focusAndSelectToken()
}

const appendModifyFieldDraft = (fieldTpl: string) => {
    if (!modifyingTarget.value) return
    const prefix = `把${modifyingTarget.value.name}的`
    if (!inputPrompt.value.trim()) {
        inputPrompt.value = `${prefix}${fieldTpl}`
    } else if (inputPrompt.value.trim() === prefix) {
        inputPrompt.value = `${prefix}${fieldTpl}`
    } else {
        inputPrompt.value = `${inputPrompt.value.trim()}，${fieldTpl}`
    }
    nextTick(() => {
        composerInputRef.value?.focus?.()
    })
}

// 一键为目标期刊打开可编辑修改卡片（无需先猜日期，直接在生成的卡片里行内填表修改）
const openDirectModifyCard = (item: PeriodicalBriefCard) => {
    modifyingTarget.value = null
    useQuickPrompt(`请为${item.name}生成修改确认单，我直接在卡片上填写要修改的字段`)
}

const handleKeydown = (e: KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault()
        handleSendMessage()
    }
}

const startThinkingClock = () => {
    thinkingElapsedMs.value = 0
    if (thinkingTimer) clearInterval(thinkingTimer)
    thinkingTimer = setInterval(() => {
        thinkingElapsedMs.value += 100
    }, 100)
}

const stopThinkingClock = () => {
    if (thinkingTimer) {
        clearInterval(thinkingTimer)
        thinkingTimer = null
    }
}

const handleSendMessage = async () => {
    const prompt = inputPrompt.value.trim()
    if (!prompt || sending.value) return

    lastSentPrompt.value = prompt
    inputPrompt.value = ''
    modifyingTarget.value = null
    activeGuideMode.value = 'none'
    sending.value = true

    // 乐观更新 (Optimistic UI)：发送瞬间立即将用户的提问气泡渲染上屏，不等接口返回
    const tempUserId = -Date.now()
    const optimisticUserMsg: AIChatMessageItem = {
        id: tempUserId,
        session_id: currentSessionId.value || 'pending',
        user_id: 0,
        username: userStore.loginUser?.name || 'admin',
        role: 'user',
        content: prompt,
        model_name: '',
        tool_steps: [],
        action_card: null,
        created_at: Math.floor(Date.now() / 1000),
    }
    messages.value.push(optimisticUserMsg)
    scrollToBottom()
    startThinkingClock()

    try {
        const res = await reqSendAiMessage({
            session_id: currentSessionId.value,
            scene: 'workspace',
            prompt,
        })
        if (res.code === 200) {
            const { session, user_message, assistant_message } = res.data
            if (!currentSessionId.value) {
                currentSessionId.value = session.session_id
            }
            // 将乐观渲染的临时用户消息替换为后端持久化的真实 user_message，再追加 AI 回复
            const tempIdx = messages.value.findIndex((m) => m.id === tempUserId)
            if (tempIdx !== -1) {
                messages.value.splice(tempIdx, 1, user_message, assistant_message)
            } else {
                messages.value.push(user_message, assistant_message)
            }
            startTypewriter(assistant_message)
            if (assistant_message.action_card?.card_type === 'query_result') {
                canvasCard.value = assistant_message.action_card
                canvasItems.value = assistant_message.action_card.items || []
                canvasTotal.value = assistant_message.action_card.total_count || canvasItems.value.length
                canvasPage.value = 1
            }
            await fetchSessions()
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
        ElMessage.error(err?.message || '发送请求异常，请重试')
    } finally {
        stopThinkingClock()
        sending.value = false
    }
}

const handleCopyMessage = async (content: string) => {
    const ok = await copyTextToClipboard(content)
    if (ok) {
        ElMessage.success('已复制回复内容')
    } else {
        ElMessage.info('请手动选择复制')
    }
}

const handleCopyFieldText = async (label: string, text: string) => {
    const ok = await copyTextToClipboard(text)
    if (ok) {
        ElMessage.success(`已复制【${label}】内容`)
    } else {
        ElMessage.info('请手动选择复制')
    }
}

const handleRegenerate = (msgIdx: number) => {
    for (let i = msgIdx - 1; i >= 0; i--) {
        if (messages.value[i].role === 'user') {
            inputPrompt.value = messages.value[i].content
            handleSendMessage()
            return
        }
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
    try {
        const res = await reqConfirmAiAction({
            message_id: messageId,
            decision,
            form_payload: formPayload,
            diffs,
            batch_updates: batchUpdates,
        })
        if (res.code === 200) {
            const idx = messages.value.findIndex((m) => m.id === messageId)
            const confirmedCard = res.data.updated_message?.action_card
            if (idx !== -1 && res.data.updated_message) {
                messages.value[idx] = res.data.updated_message
            }
            if (res.data.followup_notice) {
                messages.value.push(res.data.followup_notice)
            }
            if (decision === 'confirm') {
                ElMessage.success('数据库已更新生效，已为您生成变更明细回执')
                // 若右侧画布正打开着该期刊的详情或数据集，立即同步刷新最新数据！
                const targetPid =
                    confirmedCard?.periodical_id ||
                    res.data.followup_notice?.action_card?.periodical_id ||
                    previewJournal.value?.periodical_id
                if (targetPid && previewJournal.value?.periodical_id === targetPid) {
                    await loadFullJournalDetail(targetPid)
                }
                if (canvasOpen.value && canvasCard.value) {
                    loadCanvasPage(canvasPage.value)
                }
            } else {
                ElMessage.info('已取消该操作')
            }
            scrollToBottom()
        } else {
            ElMessage.error(res.message || '操作执行失败')
        }
    } catch (err: any) {
        ElMessage.error(err?.message || '操作失败')
    }
}

// ==================== 右侧 Artifact 分屏画布与 26 项完整详情逻辑 ====================

function mapRawPeriodicalToBrief(p: any): PeriodicalBriefCard {
    let level = ''
    let cycle = ''
    const dirs: string[] = []
    const cols: string[] = []
    if (Array.isArray(p.categorizes)) {
        p.categorizes.forEach((c: any) => {
            if (c.type === 1) dirs.push(c.name)
            else if (c.type === 2) level = c.name
            else if (c.type === 3) cols.push(c.name)
            else if (c.type === 4) cycle = c.name
        })
    }
    const pubTime = p.publication_time ? String(p.publication_time).slice(0, 10) : ''
    return {
        periodical_id: p.periodical_id,
        name: p.name,
        level,
        directions: dirs.join('、'),
        collections: cols,
        cycle,
        publication_time: pubTime,
        price: p.price || '',
        periodical_page: p.periodical_page || '',
        is_warp: Boolean(p.is_warp),
        invoice_receipt: Boolean(p.invoice_receipt),
        works: Boolean(p.works),
        is_hidden: Boolean(p.is_hidden),
        organization: p.organization || '',
        competent_unit: p.competent_unit || '',
        domestic_number: p.domestic_number || '',
        international_number: p.international_number || '',
        audit_time: p.audit_time || '',
        periodical_description: p.periodical_description || '',
    }
}

const loadFullJournalDetail = async (periodicalId: number, fallbackItem?: PeriodicalBriefCard) => {
    if (!periodicalId) return
    if (fallbackItem) {
        previewJournal.value = fallbackItem
        fullDetail.value = { ...fallbackItem }
    }
    detailLoading.value = true
    try {
        const res = await reqPeriodicalDetail(periodicalId)
        if (res.code === 200 && res.data) {
            fullDetail.value = res.data
            previewJournal.value = mapRawPeriodicalToBrief(res.data)
        }
    } finally {
        detailLoading.value = false
    }
}

const loadCanvasPage = async (page: number) => {
    if (!canvasCard.value) return
    const fp = canvasCard.value.filter_payload || {}
    canvasLoading.value = true
    try {
        let publishTimeStart = ''
        let publishTimeEnd = ''
        if (fp.publish_month && String(fp.publish_month).length === 7) {
            publishTimeStart = `${fp.publish_month}-01`
            publishTimeEnd = `${fp.publish_month}-31`
        }
        const res = await reqPeriodicalList({
            page_number: page,
            page_size: canvasPageSize,
            name: fp.name || '',
            competent_unit: fp.competent_unit || '',
            organization: fp.organization || '',
            domestic_number: fp.domestic_number || '',
            column_setting: fp.column_setting || '',
            publish_time_start: publishTimeStart,
            publish_time_end: publishTimeEnd,
            invoice: fp.invoice || 0,
            is_warp: fp.is_warp || 0,
            works: fp.works || 0,
            is_hidden: fp.is_hidden ?? 2,
            has_base_type_list: fp.has_base_type_list || [],
        })
        if (res.code === 200 && res.data) {
            canvasItems.value = (res.data.items || []).map(mapRawPeriodicalToBrief)
            canvasTotal.value = res.data.total || canvasItems.value.length
            canvasPage.value = page
        }
    } finally {
        canvasLoading.value = false
    }
}

const handleOpenCanvas = (card: ActionCardData) => {
    canvasCard.value = card
    canvasTab.value = 'dataset'
    canvasOpen.value = true
    canvasFilterKeyword.value = ''
    selectedCanvasItems.value = []
    loadCanvasPage(1)
}

const handlePreviewJournal = (item: PeriodicalBriefCard) => {
    canvasTab.value = 'detail'
    canvasOpen.value = true
    loadFullJournalDetail(item.periodical_id, item)
}

const handleSelectCanvasRow = (item: PeriodicalBriefCard) => {
    canvasTab.value = 'detail'
    loadFullJournalDetail(item.periodical_id, item)
}

const isCanvasItemSelected = (id: number) =>
    selectedCanvasItems.value.some((it) => it.periodical_id === id)

const toggleSelectCanvasItem = (item: PeriodicalBriefCard) => {
    const idx = selectedCanvasItems.value.findIndex((it) => it.periodical_id === item.periodical_id)
    if (idx === -1) {
        selectedCanvasItems.value.push(item)
    } else {
        selectedCanvasItems.value.splice(idx, 1)
    }
}

const toggleSelectAllCanvas = () => {
    if (isAllCanvasSelected.value) {
        selectedCanvasItems.value = []
    } else {
        selectedCanvasItems.value = [...filteredCanvasItems.value]
    }
}

const handleCopyBatchQuote = async () => {
    if (selectedCanvasItems.value.length === 0) return
    const header = `📋 【选刊推荐与报价方案】（共优选 ${selectedCanvasItems.value.length} 本期刊）\n----------------------------------------`
    const rows = selectedCanvasItems.value.map((item, idx) => {
        const tags = [
            item.level || '普刊',
            ...(item.collections || []),
            item.is_warp ? '全包' : '',
            item.invoice_receipt ? '可开发票' : '',
        ]
            .filter(Boolean)
            .join(' · ')
        return `${idx + 1}. ${item.name} [${tags}]\n   • 收稿方向：${item.directions || '-'}\n   • 近期排期：${item.publication_time || '待定'}${item.periodical_page ? ` (${item.periodical_page}页)` : ''} | 发表费用：${item.price && item.price !== '0' ? `¥${item.price}` : '详询'}`
    })
    const ok = await copyTextToClipboard([header, ...rows].join('\n\n'))
    if (ok) {
        ElMessage.success(`已成功复制 ${selectedCanvasItems.value.length} 本期刊的客户报价方案`)
    } else {
        ElMessage.info('请手动复制')
    }
}

const handleSessionCommand = async (cmd: string, s: AIChatSessionItem) => {
    if (cmd === 'pin') {
        const nextPin = s.is_pinned === 1 ? 0 : 1
        const res = await reqUpdateSession(s.session_id, { is_pinned: nextPin })
        if (res.code === 200) {
            ElMessage.success(nextPin === 1 ? '已置顶会话' : '已取消置顶')
            fetchSessions()
        }
    } else if (cmd === 'rename') {
        ElMessageBox.prompt('请输入新的会话标题', '重命名会话', {
            inputValue: s.title,
            confirmButtonText: '保存',
            cancelButtonText: '取消',
        })
            .then(async ({ value }) => {
                if (!value || !value.trim()) return
                const res = await reqUpdateSession(s.session_id, { title: value.trim() })
                if (res.code === 200) {
                    ElMessage.success('重命名成功')
                    fetchSessions()
                }
            })
            .catch(() => {})
    } else if (cmd === 'delete') {
        ElMessageBox.confirm(`确定要删除会话「${s.title}」及其全部消息记录吗？`, '删除确认', {
            type: 'warning',
            confirmButtonText: '删除',
            cancelButtonText: '取消',
        })
            .then(async () => {
                const res = await reqDeleteSession(s.session_id)
                if (res.code === 200) {
                    ElMessage.success('会话已删除')
                    if (currentSessionId.value === s.session_id) {
                        handleStartNewChat()
                    }
                    fetchSessions()
                }
            })
            .catch(() => {})
    }
}

const scrollToBottom = () => {
    nextTick(() => {
        if (msgScrollRef.value) {
            msgScrollRef.value.scrollTop = msgScrollRef.value.scrollHeight
        }
    })
}

const formatDateShort = (val?: string) => {
    if (!val) return ''
    return String(val).slice(0, 10)
}

const formatRelativeTime = (ts: number) => {
    if (!ts) return ''
    return dayjs(ts * 1000).format('MM-DD HH:mm')
}

const formatFullTime = (ts: number) => {
    if (!ts) return ''
    return dayjs(ts * 1000).format('YYYY-MM-DD HH:mm:ss')
}

onMounted(async () => {
    await fetchActiveModel()
    await fetchSessions()
    const querySessionId = route.query.session_id as string
    if (querySessionId) {
        currentSessionId.value = querySessionId
        loadMessages(querySessionId)
    } else if (sessionList.value.length > 0) {
        currentSessionId.value = sessionList.value[0].session_id
        loadMessages(sessionList.value[0].session_id)
    }
})

// Keep-Alive 切回激活钩子：无损保留草稿、对话与右侧画布，仅在带审计 session_id 跳转时切换会话
onActivated(() => {
    fetchActiveModel()
    const querySessionId = route.query.session_id as string
    if (querySessionId && querySessionId !== currentSessionId.value) {
        currentSessionId.value = querySessionId
        loadMessages(querySessionId)
    }
})

onBeforeUnmount(() => {
    stopThinkingClock()
    if (typingTimer) clearInterval(typingTimer)
})
</script>

<style scoped lang="scss">
.ai-workspace-page {
    display: flex;
    height: calc(100vh - 128px);
    background: #ffffff;
    border: 1px solid #e2e8f0;
    border-radius: 14px;
    overflow: hidden;
    box-shadow: 0 2px 8px rgba(15, 23, 42, 0.04);
}

/* 左侧历史会话栏 */
.session-sidebar {
    width: 262px;
    flex-shrink: 0;
    background: #f8fafc;
    border-right: 1px solid #e2e8f0;
    display: flex;
    flex-direction: column;

    .sidebar-top {
        padding: 16px;
        border-bottom: 1px solid #e2e8f0;
        display: flex;
        flex-direction: column;
        gap: 10px;

        .new-chat-btn {
            width: 100%;
            height: 38px;
            font-weight: 600;
            border-radius: 8px;
        }
    }

    .session-list {
        flex: 1;
        overflow-y: auto;
        padding: 10px;
        display: flex;
        flex-direction: column;
        gap: 4px;

        .session-item {
            display: flex;
            align-items: center;
            gap: 10px;
            padding: 10px 12px;
            border-radius: 9px;
            cursor: pointer;
            transition: all 0.15s ease;
            border: 1px solid transparent;

            &:hover {
                background: #f1f5f9;

                .session-ops {
                    opacity: 1;
                }
            }

            &.active {
                background: #eff6ff;
                border-color: #bfdbfe;

                .s-title {
                    color: #1d4ed8;
                    font-weight: 600;
                }
            }

            .session-icon {
                color: #64748b;
                font-size: 16px;
                display: flex;
                align-items: center;
            }

            .session-info {
                flex: 1;
                min-width: 0;

                .s-title {
                    font-size: 13px;
                    color: #1e293b;
                    font-weight: 500;
                    overflow: hidden;
                    text-overflow: ellipsis;
                    white-space: nowrap;
                    margin-bottom: 3px;
                }

                .s-meta {
                    display: flex;
                    align-items: center;
                    gap: 4px;
                    font-size: 11px;
                    color: #94a3b8;
                }
            }

            .session-ops {
                opacity: 0;
                transition: opacity 0.15s ease;

                .more-trigger {
                    padding: 4px;
                    color: #64748b;
                    cursor: pointer;
                }
            }
        }

        .empty-sessions {
            padding: 40px 16px;
            text-align: center;
            color: #94a3b8;

            p {
                margin: 0 0 4px;
                font-size: 13px;
                font-weight: 600;
                color: #64748b;
            }

            span {
                font-size: 12px;
            }
        }
    }

    .sidebar-footer {
        padding: 12px 14px;
        border-top: 1px solid #e2e8f0;
        background: #ffffff;

        .model-status-box {
            display: flex;
            align-items: center;
            justify-content: space-between;

            .status-left {
                display: flex;
                align-items: center;
                gap: 8px;
                min-width: 0;

                .pulse-dot {
                    width: 8px;
                    height: 8px;
                    border-radius: 50%;
                    background: #10b981;
                    flex-shrink: 0;

                    &.active {
                        box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.2);
                    }
                }

                .model-meta {
                    display: flex;
                    flex-direction: column;
                    min-width: 0;

                    .m-label {
                        font-size: 10.5px;
                        color: #94a3b8;
                    }

                    .m-name {
                        font-size: 12px;
                        font-weight: 600;
                        color: #334155;
                        overflow: hidden;
                        text-overflow: ellipsis;
                        white-space: nowrap;
                    }
                }
            }
        }
    }
}

/* 中间对话主舞台 */
.chat-stage {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    background: #ffffff;

    .stage-header {
        height: 56px;
        padding: 0 24px;
        border-bottom: 1px solid #f1f5f9;
        display: flex;
        align-items: center;
        justify-content: space-between;
        flex-shrink: 0;

        .header-left {
            display: flex;
            align-items: center;
            gap: 12px;
            min-width: 0;

            .stage-title {
                margin: 0;
                font-size: 15px;
                font-weight: 700;
                color: #0f172a;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
            }

            .model-badge {
                display: inline-flex;
                align-items: center;
                gap: 4px;
                padding: 2px 10px;
                border-radius: 999px;
                background: #eff6ff;
                color: #2563eb;
                font-size: 12px;
                font-weight: 600;
                flex-shrink: 0;
            }
        }

        .header-right {
            display: flex;
            align-items: center;
            gap: 8px;
        }
    }

    .message-viewport {
        flex: 1;
        overflow-y: auto;
        padding: 24px 28px;
        background: #fcfdff;
    }
}

.welcome-hero {
    width: 100%;
    max-width: clamp(860px, 84%, 1180px);
    margin: 28px auto 0;
    text-align: center;

    .hero-badge {
        display: inline-block;
        padding: 4px 12px;
        border-radius: 999px;
        background: #eff6ff;
        color: #2563eb;
        font-size: 11px;
        font-weight: 700;
        letter-spacing: 0.08em;
        margin-bottom: 12px;
    }

    h1 {
        margin: 0 0 10px;
        font-size: 24px;
        font-weight: 700;
        color: #0f172a;
    }

    .hero-sub {
        margin: 0 auto 28px;
        max-width: 620px;
        font-size: 14px;
        line-height: 1.65;
        color: #64748b;
    }

    .prompt-cards-grid {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 14px;
        text-align: left;

        .prompt-card {
            padding: 16px 18px;
            border-radius: 12px;
            border: 1px solid #e2e8f0;
            background: #ffffff;
            cursor: pointer;
            transition: all 0.18s ease;

            &:hover {
                border-color: #93c5fd;
                box-shadow: 0 6px 16px -4px rgba(37, 99, 235, 0.12);
                transform: translateY(-2px);
            }

            .p-header {
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: 8px;
                margin-bottom: 8px;

                .p-header-left {
                    display: flex;
                    align-items: center;
                    gap: 8px;
                    min-width: 0;
                }

                .p-mode-hint {
                    font-size: 11px;
                    font-weight: 600;
                    padding: 2px 8px;
                    border-radius: 999px;
                    background: #f1f5f9;
                    color: #475569;
                    flex-shrink: 0;

                    &.send {
                        background: #eff6ff;
                        color: #2563eb;
                    }
                }

                .p-tag {
                    font-size: 11px;
                    font-weight: 700;
                    padding: 2px 8px;
                    border-radius: 5px;

                    &.blue {
                        background: #eff6ff;
                        color: #2563eb;
                    }
                    &.emerald {
                        background: #ecfdf5;
                        color: #059669;
                    }
                    &.amber {
                        background: #fffbeb;
                        color: #d97706;
                    }
                    &.rose {
                        background: #fef2f2;
                        color: #dc2626;
                    }
                }

                .p-title {
                    font-size: 14px;
                    font-weight: 700;
                    color: #0f172a;
                }
            }

            .p-text {
                margin: 0;
                font-size: 12.5px;
                color: #475569;
                line-height: 1.55;
            }
        }
    }
}

.message-list {
    width: 100%;
    max-width: clamp(980px, 92%, 1440px);
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    gap: 22px;

    .message-row {
        display: flex;
        gap: 12px;
        align-items: flex-start;

        &.is-user {
            flex-direction: row-reverse;

            .bubble-col {
                align-items: flex-end;
                max-width: 78%;
            }

            .msg-bubble {
                background: #2563eb;
                color: #ffffff;
                border-color: #2563eb;

                :deep(.md-book-tag) {
                    color: #dbeafe;
                }
            }
        }

        .avatar-col {
            flex-shrink: 0;

            .user-avatar,
            .ai-avatar {
                width: 34px;
                height: 34px;
                border-radius: 9px;
                display: flex;
                align-items: center;
                justify-content: center;
                font-size: 13px;
                font-weight: 700;
                color: #ffffff;
            }

            .user-avatar {
                background: #334155;
            }

            .ai-avatar {
                background: linear-gradient(135deg, #2563eb 0%, #0284c7 100%);
                box-shadow: 0 2px 8px rgba(37, 99, 235, 0.25);
            }
        }

        .bubble-col {
            flex: 1;
            min-width: 0;
            max-width: 95%;
            display: flex;
            flex-direction: column;

            .msg-meta {
                display: flex;
                align-items: center;
                gap: 8px;
                margin-bottom: 6px;
                font-size: 12px;

                .sender-name {
                    font-weight: 600;
                    color: #334155;
                }

                .msg-model-tag {
                    font-size: 10.5px;
                    padding: 1px 6px;
                    border-radius: 4px;
                    background: #f1f5f9;
                    color: #64748b;
                    font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
                }

                .send-time {
                    color: #94a3b8;
                }
            }

            /* 一体化响应容器：工具调用眉栏 + 正文说明 + 结构化操作卡片融为一个整体 */
            .ai-turn-container {
                width: 100%;
                border: 1px solid #e2e8f0;
                border-radius: 14px;
                background: #ffffff;
                overflow: hidden;
                box-shadow: 0 2px 10px -2px rgba(15, 23, 42, 0.05);

                &.is-user-turn {
                    width: auto;
                    border: none;
                    background: transparent;
                    box-shadow: none;
                    overflow: visible;

                    .msg-bubble {
                        border-radius: 12px;
                        border: 1px solid #2563eb;
                        padding: 10px 16px;
                    }
                }
            }

            /* 可折叠工具执行轨迹手风琴（内嵌于容器顶部眉栏） */
            .agent-trace-accordion {
                margin-bottom: 0;
                border-bottom: 1px solid #dcfce7;
                background: #f0fdf4;
                overflow: hidden;
                transition: all 0.15s ease;

                .trace-header {
                    padding: 7px 16px;
                    display: flex;
                    align-items: center;
                    justify-content: space-between;
                    gap: 10px;
                    cursor: pointer;
                    user-select: none;

                    &:hover {
                        background: #dcfce7;
                    }

                    .trace-left {
                        display: flex;
                        align-items: center;
                        gap: 6px;
                        min-width: 0;
                        flex-wrap: wrap;

                        .trace-status-icon {
                            font-size: 14px;
                            &.success {
                                color: #16a34a;
                            }
                            &.error {
                                color: #dc2626;
                            }
                        }

                        .trace-title {
                            font-size: 12px;
                            font-weight: 700;
                            color: #166534;
                        }

                        .trace-summary-badge {
                            font-size: 11.5px;
                            color: #15803d;
                        }
                    }

                    .trace-right {
                        display: flex;
                        align-items: center;
                        gap: 6px;
                        flex-shrink: 0;

                        .trace-ms {
                            font-size: 11px;
                            font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
                            color: #15803d;
                            background: rgba(255, 255, 255, 0.7);
                            padding: 1px 5px;
                            border-radius: 4px;
                        }

                        .trace-arrow {
                            font-size: 12px;
                            color: #15803d;
                            transition: transform 0.2s ease;
                        }
                    }
                }

                &.expanded .trace-arrow {
                    transform: rotate(180deg);
                }

                .trace-detail-list {
                    padding: 8px 16px;
                    background: #f8fafc;
                    border-top: 1px solid #dcfce7;
                    display: flex;
                    flex-direction: column;
                    gap: 6px;

                    .trace-step-row {
                        font-size: 12px;
                        color: #334155;
                        display: flex;
                        flex-direction: column;
                        gap: 3px;

                        .step-meta-line {
                            display: flex;
                            align-items: center;
                            gap: 8px;

                            .fn-badge {
                                font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
                                font-size: 11px;
                                padding: 1px 6px;
                                border-radius: 4px;
                                background: #eff6ff;
                                color: #1d4ed8;
                                font-weight: 600;
                            }

                            .step-state {
                                font-size: 11px;
                                font-weight: 600;
                                color: #16a34a;
                            }

                            .step-dur {
                                font-size: 11px;
                                color: #94a3b8;
                            }
                        }

                        .step-kv {
                            color: #475569;
                            line-height: 1.45;
                        }
                    }
                }
            }

            .msg-bubble {
                padding: 14px 18px;
                background: #ffffff;
                color: #1e293b;
                font-size: 14px;
                line-height: 1.65;
                word-break: break-word;
            }

            .embedded-card-slot {
                padding: 0 16px 16px;
                background: #ffffff;

                :deep(.ai-action-card) {
                    margin-top: 0;
                    box-shadow: none;
                }
            }

            /* 消息底部操作栏 */
            .msg-action-bar {
                margin-top: 6px;
                display: flex;
                align-items: center;
                gap: 12px;
                flex-wrap: wrap;

                .bar-btn {
                    display: inline-flex;
                    align-items: center;
                    gap: 4px;
                    border: none;
                    background: transparent;
                    padding: 2px 4px;
                    font-size: 12px;
                    color: #64748b;
                    cursor: pointer;
                    border-radius: 4px;
                    transition: all 0.15s ease;

                    &:hover:not(:disabled) {
                        color: #2563eb;
                        background: #eff6ff;
                    }

                    &.primary-text {
                        color: #2563eb;
                        font-weight: 600;
                    }
                }
            }

            /* 智能追问建议胶囊 */
            .followup-suggestions {
                margin-top: 10px;
                display: flex;
                align-items: center;
                flex-wrap: wrap;
                gap: 8px;

                .suggest-label {
                    font-size: 12px;
                    font-weight: 600;
                    color: #64748b;
                }

                .suggest-pill {
                    display: inline-flex;
                    align-items: center;
                    gap: 5px;
                    padding: 5px 11px;
                    border-radius: 999px;
                    border: 1px solid #dbeafe;
                    background: #f8fbff;
                    color: #1d4ed8;
                    font-size: 12px;
                    font-weight: 500;
                    cursor: pointer;
                    transition: all 0.15s ease;

                    &:hover {
                        background: #eff6ff;
                        border-color: #93c5fd;
                        transform: translateX(2px);
                    }
                }
            }
        }
    }
}

.thinking-live-pill {
    font-size: 11px;
    font-weight: 600;
    padding: 1px 8px;
    border-radius: 999px;
    background: #eff6ff;
    color: #2563eb;
    border: 1px solid #bfdbfe;
}

.agent-thinking-card {
    width: 100%;
    max-width: 680px;
    padding: 14px 18px;
    border-radius: 12px;
    background: linear-gradient(180deg, #ffffff 0%, #f8fbff 100%);
    border: 1px solid #dbeafe;
    box-shadow: 0 4px 14px -4px rgba(37, 99, 235, 0.08);
    display: flex;
    flex-direction: column;
    gap: 12px;

    .thinking-main-row {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 12px;

        .thinking-left {
            display: flex;
            align-items: center;
            gap: 10px;
            min-width: 0;
        }

        .thinking-orb {
            position: relative;
            width: 18px;
            height: 18px;
            display: inline-flex;
            align-items: center;
            justify-content: center;
            flex-shrink: 0;

            .orb-ring {
                position: absolute;
                inset: 0;
                border-radius: 50%;
                border: 2px solid #93c5fd;
                border-top-color: #2563eb;
                animation: spin-orb 0.9s linear infinite;
            }

            .orb-core {
                width: 6px;
                height: 6px;
                border-radius: 50%;
                background: #2563eb;
            }
        }

        .thinking-stage-text {
            font-size: 13.5px;
            font-weight: 600;
            color: #1e293b;
            line-height: 1.4;
        }

        .thinking-timer {
            font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
            font-size: 11.5px;
            color: #2563eb;
            background: #eff6ff;
            border: 1px solid #bfdbfe;
            padding: 2px 8px;
            border-radius: 999px;
            font-weight: 600;
            flex-shrink: 0;
        }
    }

    .thinking-steps-track {
        display: flex;
        align-items: center;
        flex-wrap: wrap;
        gap: 6px;
        padding-top: 10px;
        border-top: 1px dashed #e2e8f0;

        .t-step-node {
            display: inline-flex;
            align-items: center;
            gap: 5px;
            padding: 3px 10px;
            border-radius: 999px;
            font-size: 11.5px;
            font-weight: 500;
            transition: all 0.2s ease;

            .node-indicator {
                width: 15px;
                height: 15px;
                border-radius: 50%;
                display: inline-flex;
                align-items: center;
                justify-content: center;
                font-size: 10px;
                font-weight: 700;
                line-height: 1;
            }

            &.done {
                background: #f0fdf4;
                color: #15803d;
                border: 1px solid #bbf7d0;

                .node-indicator {
                    background: #16a34a;
                    color: #ffffff;
                }
            }

            &.active {
                background: #eff6ff;
                color: #1d4ed8;
                font-weight: 600;
                border: 1px solid #93c5fd;

                .node-indicator {
                    background: #dbeafe;
                }

                .pulse-mini {
                    width: 6px;
                    height: 6px;
                    border-radius: 50%;
                    background: #2563eb;
                    animation: pulse-dot 1.2s ease-in-out infinite;
                }
            }

            &.pending {
                background: #f8fafc;
                color: #94a3b8;
                border: 1px solid #e2e8f0;

                .node-indicator {
                    background: #e2e8f0;
                    color: #64748b;
                }
            }
        }

        .t-step-divider {
            width: 14px;
            height: 1px;
            background: #cbd5e1;
        }
    }
}

@keyframes spin-orb {
    from {
        transform: rotate(0deg);
    }
    to {
        transform: rotate(360deg);
    }
}

@keyframes pulse-dot {
    0%,
    100% {
        transform: scale(0.85);
        opacity: 0.6;
    }
    50% {
        transform: scale(1.2);
        opacity: 1;
    }
}

.composer-footer {
    padding: 12px 28px 18px;
    background: #ffffff;
    border-top: 1px solid #f1f5f9;

    /* 锁定修改目标期刊时的专属引导栏 */
    .target-modify-guide {
        width: 100%;
        max-width: clamp(980px, 92%, 1440px);
        margin: 0 auto 10px;
        padding: 10px 14px;
        border-radius: 10px;
        background: #f0f7ff;
        border: 1px solid #bfdbfe;
        display: flex;
        flex-direction: column;
        gap: 8px;

        &.create-guide-theme {
            background: #f0fdf4;
            border-color: #86efac;
        }

        &.modify-general-theme {
            background: #fffbeb;
            border-color: #fcd34d;
        }

        .guide-top {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 8px;

            .guide-title {
                display: flex;
                align-items: center;
                gap: 6px;
                font-size: 12.5px;
                color: #1e293b;
                flex-wrap: wrap;

                .target-badge {
                    padding: 2px 8px;
                    border-radius: 999px;
                    background: #2563eb;
                    color: #ffffff;
                    font-size: 11px;
                    font-weight: 700;

                    &.emerald {
                        background: #059669;
                    }

                    &.amber {
                        background: #d97706;
                    }
                }

                .current-meta {
                    color: #64748b;
                    font-size: 12px;
                }
            }

            .guide-top-ops {
                display: flex;
                align-items: center;
                gap: 8px;
                flex-shrink: 0;

                .tpl-action-btn {
                    border: 1px solid #6ee7b7;
                    background: #ffffff;
                    color: #047857;
                    font-size: 11.5px;
                    font-weight: 600;
                    padding: 3px 9px;
                    border-radius: 6px;
                    cursor: pointer;
                    transition: all 0.15s ease;

                    &:hover {
                        background: #d1fae5;
                    }
                }
            }

            .cancel-target-btn {
                border: none;
                background: transparent;
                color: #64748b;
                font-size: 12px;
                cursor: pointer;

                &:hover {
                    color: #dc2626;
                }
            }
        }

        .guide-fields {
            display: flex;
            align-items: center;
            flex-wrap: wrap;
            gap: 6px;

            .guide-hint {
                font-size: 11.5px;
                color: #475569;
            }

            .mf-chip {
                padding: 3px 10px;
                border-radius: 999px;
                border: 1px solid #93c5fd;
                background: #ffffff;
                color: #1d4ed8;
                font-size: 11.5px;
                font-weight: 600;
                cursor: pointer;
                transition: all 0.15s ease;

                &:hover {
                    background: #dbeafe;
                }

                &.direct-card,
                &.emerald-chip {
                    background: #ffffff;
                    border-color: #6ee7b7;
                    color: #047857;

                    &:hover {
                        background: #d1fae5;
                    }
                }

                &.amber-chip {
                    background: #ffffff;
                    border-color: #fcd34d;
                    color: #b45309;

                    &:hover {
                        background: #fef3c7;
                    }
                }
            }
        }
    }

    .quick-chips {
        width: 100%;
        max-width: clamp(980px, 92%, 1440px);
        margin: 0 auto 10px;
        display: flex;
        flex-wrap: wrap;
        gap: 8px;

        .chip-item {
            padding: 4px 10px;
            border-radius: 999px;
            background: #f8fafc;
            border: 1px solid #e2e8f0;
            color: #475569;
            font-size: 12px;
            cursor: pointer;
            transition: all 0.15s ease;

            &:hover {
                background: #eff6ff;
                border-color: #bfdbfe;
                color: #2563eb;
            }
        }
    }

    .composer-box {
        width: 100%;
        max-width: clamp(980px, 92%, 1440px);
        margin: 0 auto;
        border: 1px solid #cbd5e1;
        border-radius: 12px;
        padding: 10px 12px;
        background: #ffffff;
        transition: border-color 0.2s ease, box-shadow 0.2s ease;

        &.is-targeting,
        &:focus-within {
            border-color: #2563eb;
            box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
        }

        &.is-creating {
            border-color: #10b981;
            box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.12);
        }

        :deep(.el-textarea__inner) {
            border: none !important;
            box-shadow: none !important;
            padding: 2px 4px;
            font-size: 14px;
            resize: none;
        }

        .composer-actions {
            margin-top: 8px;
            display: flex;
            align-items: center;
            justify-content: space-between;

            .hint-text {
                font-size: 12px;
                color: #94a3b8;
            }
        }
    }
}

/* ==================== 右侧 Artifact 数据与 26 项完整档案分屏画布 ==================== */
.artifact-canvas-pane {
    width: 530px;
    flex-shrink: 0;
    border-left: 1px solid #e2e8f0;
    background: #f8fafc;
    display: flex;
    flex-direction: column;
    overflow: hidden;

    .canvas-header {
        height: 56px;
        padding: 0 14px;
        background: #ffffff;
        border-bottom: 1px solid #e2e8f0;
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 8px;

        .canvas-tabs {
            display: flex;
            align-items: center;
            gap: 6px;

            .c-tab {
                display: inline-flex;
                align-items: center;
                gap: 6px;
                padding: 6px 12px;
                border-radius: 8px;
                border: 1px solid transparent;
                background: transparent;
                font-size: 13px;
                font-weight: 600;
                color: #64748b;
                cursor: pointer;
                transition: all 0.15s ease;

                &.active {
                    background: #eff6ff;
                    color: #2563eb;
                    border-color: #bfdbfe;
                }
            }
        }

        .close-canvas-btn {
            width: 30px;
            height: 30px;
            border-radius: 7px;
            border: 1px solid #e2e8f0;
            background: #ffffff;
            color: #64748b;
            display: flex;
            align-items: center;
            justify-content: center;
            cursor: pointer;

            &:hover {
                color: #0f172a;
                background: #f1f5f9;
            }
        }
    }

    .canvas-dataset-body {
        flex: 1;
        display: flex;
        flex-direction: column;
        overflow: hidden;

        .dataset-toolbar {
            padding: 12px 14px 8px;
            display: flex;
            align-items: center;
            gap: 8px;

            .canvas-search {
                flex: 1;
            }
        }

        .dataset-condition-bar {
            margin: 0 14px 8px;
            padding: 6px 10px;
            border-radius: 7px;
            background: #eff6ff;
            border: 1px solid #dbeafe;
            font-size: 12px;
            color: #1e40af;
            display: flex;
            align-items: center;
            justify-content: space-between;

            .total-badge {
                font-weight: 700;
            }
        }

        .dataset-table-scroll {
            flex: 1;
            overflow-y: auto;
            margin: 0 14px;
            background: #ffffff;
            border: 1px solid #e2e8f0;
            border-radius: 9px;

            .canvas-data-table {
                width: 100%;
                border-collapse: collapse;
                font-size: 12.5px;

                th,
                td {
                    padding: 9px 8px;
                    border-bottom: 1px solid #f1f5f9;
                    text-align: left;
                }

                th {
                    position: sticky;
                    top: 0;
                    z-index: 2;
                    background: #f8fafc;
                    color: #475569;
                    font-weight: 600;
                }

                tbody tr {
                    cursor: pointer;
                    transition: background 0.12s ease;

                    &:hover td {
                        background: #f8fbff;
                    }

                    &.active td {
                        background: #eff6ff;
                    }
                }

                .c-name-line {
                    display: flex;
                    align-items: center;
                    gap: 4px;
                    flex-wrap: wrap;

                    .c-book-name {
                        font-weight: 700;
                        color: #0f172a;
                    }

                    .c-mini-tag {
                        font-size: 10px;
                        padding: 1px 4px;
                        border-radius: 4px;

                        &.emerald {
                            background: #d1fae5;
                            color: #047857;
                        }
                        &.sky {
                            background: #e0f2fe;
                            color: #0369a1;
                        }
                    }
                }

                .c-sub-dir {
                    font-size: 11px;
                    color: #94a3b8;
                    margin-top: 2px;
                    max-width: 155px;
                    overflow: hidden;
                    text-overflow: ellipsis;
                    white-space: nowrap;
                }

                .c-level-pill {
                    font-size: 11px;
                    padding: 1px 6px;
                    border-radius: 999px;
                    background: #dbeafe;
                    color: #1d4ed8;
                    font-weight: 600;
                    white-space: nowrap;
                }

                .c-col-wrap {
                    display: flex;
                    flex-wrap: wrap;
                    gap: 3px;

                    .c-col-tag {
                        font-size: 10.5px;
                        padding: 1px 5px;
                        border-radius: 4px;
                        background: #f1f5f9;
                        color: #475569;
                    }
                }

                .c-pub {
                    font-size: 11.5px;
                    color: #475569;

                    .c-page {
                        font-size: 10.5px;
                        color: #94a3b8;
                    }
                }

                .c-price {
                    text-align: right;
                    font-weight: 700;
                    color: #0f172a;
                    font-variant-numeric: tabular-nums;
                }
            }
        }

        .canvas-footer {
            padding: 10px 14px;
            background: #ffffff;
            border-top: 1px solid #e2e8f0;
            margin-top: 10px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            font-size: 12px;
            color: #64748b;

            .p-controls {
                display: flex;
                gap: 6px;
            }
        }
    }

    .canvas-detail-body {
        flex: 1;
        overflow-y: auto;
        padding: 14px;
        display: flex;
        flex-direction: column;
        gap: 12px;

        .detail-hero-card {
            background: #ffffff;
            border: 1px solid #e2e8f0;
            border-radius: 10px;
            padding: 14px 16px;

            .d-top {
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: 8px;
                margin-bottom: 8px;

                .d-title-wrap {
                    display: flex;
                    align-items: center;
                    gap: 8px;

                    h3 {
                        margin: 0;
                        font-size: 16px;
                        color: #0f172a;
                    }

                    .d-id-badge {
                        font-size: 11px;
                        padding: 1px 6px;
                        border-radius: 4px;
                        background: #f1f5f9;
                        color: #64748b;
                        font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
                    }
                }

                .d-level {
                    font-size: 11.5px;
                    padding: 2px 8px;
                    border-radius: 999px;
                    background: #dbeafe;
                    color: #1d4ed8;
                    font-weight: 700;
                }
            }

            .d-badges {
                display: flex;
                flex-wrap: wrap;
                gap: 6px;
                margin-bottom: 12px;

                .d-tag {
                    font-size: 11px;
                    padding: 2px 8px;
                    border-radius: 5px;
                    font-weight: 600;

                    &.emerald {
                        background: #d1fae5;
                        color: #047857;
                    }
                    &.sky {
                        background: #e0f2fe;
                        color: #0369a1;
                    }
                    &.amber {
                        background: #fef3c7;
                        color: #b45309;
                    }
                    &.slate {
                        background: #eff6ff;
                        color: #1d4ed8;
                    }
                    &.gray {
                        background: #f1f5f9;
                        color: #64748b;
                    }
                }
            }

            .d-price-row {
                padding-top: 10px;
                border-top: 1px dashed #e2e8f0;
                display: flex;
                align-items: center;
                justify-content: space-between;
                flex-wrap: wrap;
                gap: 8px;
                font-size: 12.5px;

                .label {
                    color: #64748b;
                }

                .price-val {
                    font-size: 15px;
                    color: #059669;
                }
            }
        }

        .detail-section-card {
            background: #ffffff;
            border: 1px solid #e2e8f0;
            border-radius: 10px;
            padding: 12px 14px;

            .sec-header {
                display: flex;
                align-items: center;
                justify-content: space-between;
                margin-bottom: 10px;
                padding-bottom: 6px;
                border-bottom: 1px solid #f1f5f9;

                .sec-title {
                    font-size: 13px;
                    font-weight: 700;
                    color: #0f172a;
                }

                .sec-sub {
                    font-size: 11px;
                    color: #94a3b8;
                }
            }

            .detail-kv-grid {
                display: grid;
                grid-template-columns: repeat(2, 1fr);
                gap: 8px;

                .kv-box {
                    background: #f8fafc;
                    border: 1px solid #f1f5f9;
                    border-radius: 7px;
                    padding: 7px 10px;
                    display: flex;
                    flex-direction: column;
                    gap: 2px;

                    &.full {
                        grid-column: span 2;
                    }

                    .k {
                        font-size: 11px;
                        color: #64748b;
                    }

                    .v {
                        font-size: 12.5px;
                        font-weight: 600;
                        color: #1e293b;
                        word-break: break-word;
                    }
                }
            }

            .rich-fields-list {
                display: flex;
                flex-direction: column;
                gap: 8px;

                .rich-field-item {
                    background: #f8fafc;
                    border: 1px solid #f1f5f9;
                    border-radius: 7px;
                    padding: 8px 10px;

                    .rf-top {
                        display: flex;
                        align-items: center;
                        justify-content: space-between;
                        margin-bottom: 4px;

                        .rf-label {
                            font-size: 11.5px;
                            font-weight: 700;
                            color: #334155;
                        }

                        .rf-copy-btn {
                            border: none;
                            background: #eff6ff;
                            color: #2563eb;
                            font-size: 11px;
                            font-weight: 600;
                            padding: 1px 7px;
                            border-radius: 4px;
                            cursor: pointer;

                            &:hover {
                                background: #dbeafe;
                            }
                        }
                    }

                    .rf-val {
                        font-size: 12.5px;
                        color: #475569;
                        line-height: 1.55;
                        white-space: pre-wrap;
                        word-break: break-word;
                    }
                }
            }
        }

        .detail-actions {
            position: sticky;
            bottom: 0;
            background: #f8fafc;
            padding: 10px 0 2px;
            display: flex;
            align-items: center;
            justify-content: flex-end;
            gap: 8px;
        }
    }
}

.canvas-slide-enter-active,
.canvas-slide-leave-active {
    transition: all 0.22s cubic-bezier(0.22, 1, 0.36, 1);
}

.canvas-slide-enter-from,
.canvas-slide-leave-to {
    opacity: 0;
    transform: translateX(30px);
}
</style>
