<template>
    <div v-if="card" class="ai-action-card" :class="[`card-${card.card_type}`, `status-${card.status}`]">
        <!-- 卡片头部 -->
        <div class="card-header">
            <div class="header-left">
                <span class="type-dot" :class="card.card_type"></span>
                <span class="card-title">{{ card.title }}</span>
            </div>
            <div class="header-right">
                <span v-if="card.status === 'pending'" class="status-pill pending">待确认执行 · 可直接微调</span>
                <span v-else-if="card.status === 'confirmed'" class="status-pill confirmed">已确认执行</span>
                <span v-else-if="card.status === 'cancelled'" class="status-pill cancelled">已取消</span>

                <!-- 新增确认卡片：卡片内编辑开关 -->
                <el-button
                    v-if="card.card_type === 'create_confirm' && card.status === 'pending' && !readonly"
                    size="small"
                    :type="isInlineEditingCreate ? 'primary' : 'default'"
                    plain
                    icon="Edit"
                    @click="isInlineEditingCreate = !isInlineEditingCreate"
                >
                    {{ isInlineEditingCreate ? '完成微调' : '在卡片内微调' }}
                </el-button>

                <!-- 查询结果专属工具栏：视图切换 + 画布展开 + 同步筛选 -->
                <template v-if="card.card_type === 'query_result' && displayItems.length > 0">
                    <div class="view-switcher">
                        <button
                            type="button"
                            class="switch-btn"
                            :class="{ active: viewMode === 'grid' }"
                            @click="viewMode = 'grid'"
                        >
                            网格
                        </button>
                        <button
                            type="button"
                            class="switch-btn"
                            :class="{ active: viewMode === 'table' }"
                            @click="viewMode = 'table'"
                        >
                            表格
                        </button>
                    </div>

                    <el-button
                        v-if="showOpenCanvas"
                        size="small"
                        type="primary"
                        plain
                        icon="DataBoard"
                        @click="emit('open-canvas', card)"
                    >
                        在右侧画布展开
                    </el-button>

                    <el-button
                        v-if="showSyncFilter && card.filter_payload"
                        size="small"
                        type="primary"
                        plain
                        icon="Filter"
                        @click="emit('sync-filter', card.filter_payload)"
                    >
                        同步至左侧筛选
                    </el-button>
                </template>
            </div>
        </div>

        <div v-if="card.card_type === 'update_confirm' && card.status === 'pending' && !readonly" class="update-summary-bar">
            <span v-if="isBatchUpdate" class="summary-text">
                当前共锁定 <strong>{{ editableBatchTargets.length }}</strong> 本待修改期刊（合计
                <strong>{{ totalBatchDiffCount }}</strong> 项字段变更），支持对每本刊物独立微调、移除某刊或一键同步新值：
            </span>
            <span v-else class="summary-text">
                当前包含 <strong>{{ editableDiffs.length }}</strong> 项待修改字段（其中
                <strong class="changed-num">{{ changedDiffCount }}</strong> 项已修改新值），可随时输入新值、添加或移除字段：
            </span>
            <button
                v-if="!isBatchUpdate && editableDiffs.length > changedDiffCount && changedDiffCount > 0"
                type="button"
                class="clean-unchanged-btn"
                @click="clearUnchangedFields"
            >
                清理未改动项 ({{ editableDiffs.length - changedDiffCount }})
            </button>
        </div>
        <p v-else-if="card.description" class="card-desc">{{ card.description }}</p>

        <!-- 1. 查询结果期刊卡片列表 (query_result) -->
        <div v-if="card.card_type === 'query_result'" v-loading="pageLoading" class="query-result-body">
            <template v-if="displayItems.length > 0">
                <!-- 模式 A：网格卡片视图 -->
                <div v-if="viewMode === 'grid'" class="journal-grid">
                    <div
                        v-for="item in displayItems"
                        :key="item.periodical_id"
                        class="journal-mini-card"
                        @click="emit('select-journal', item)"
                    >
                        <div class="mini-top">
                            <span class="j-name" :title="item.name">{{ item.name }}</span>
                            <span v-if="item.level" class="j-level">{{ item.level }}</span>
                        </div>
                        <div class="mini-badges">
                            <span v-if="item.is_hidden" class="b-tag amber-hidden">已隐藏</span>
                            <span v-if="item.is_warp" class="b-tag emerald">全包</span>
                            <span v-if="item.invoice_receipt" class="b-tag sky">可开发票</span>
                            <span v-for="col in (item.collections || []).slice(0, 3)" :key="col" class="b-tag slate">
                                {{ col }}
                            </span>
                            <span v-if="item.directions" class="b-tag dir" :title="item.directions">
                                {{ item.directions.split('、')[0] }}
                            </span>
                        </div>
                        <div class="mini-bottom">
                            <span class="j-meta">
                                出刊：<strong>{{ item.publication_time || '待定' }}</strong>
                                <template v-if="item.periodical_page"> · {{ formatPageText(item.periodical_page) }}</template>
                            </span>
                            <span class="j-price">
                                {{ item.price && item.price !== '0' ? `¥${item.price}` : '待询' }}
                            </span>
                        </div>
                        <!-- 单刊快捷操作微栏 -->
                        <div class="mini-actions" @click.stop>
                            <span class="act-link" @click="emit('select-journal', item)">详情</span>
                            <span class="act-sep">·</span>
                            <span class="act-link" @click="handleCopyDescription(item)">复制描述</span>
                            <template v-if="!readonly">
                                <span class="act-sep">·</span>
                                <span class="act-link" @click="handleTriggerModify(item)">改刊</span>
                            </template>
                        </div>
                    </div>
                </div>

                <!-- 模式 B：紧凑商务表格视图 -->
                <div v-else class="journal-table-wrap">
                    <table class="journal-compact-table">
                        <thead>
                            <tr>
                                <th style="width: 52px;">ID</th>
                                <th>期刊名称</th>
                                <th style="width: 76px;">级别</th>
                                <th>网站收录</th>
                                <th>收稿方向</th>
                                <th style="width: 110px;">排期 / 页码</th>
                                <th style="width: 86px; text-align: right;">发表费</th>
                                <th style="width: 116px; text-align: right;">操作</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="item in displayItems" :key="item.periodical_id">
                                <td class="col-id">#{{ item.periodical_id }}</td>
                                <td class="col-name" @click="emit('select-journal', item)">
                                    <span class="name-link">{{ item.name }}</span>
                                    <span v-if="item.is_hidden" class="mini-pill amber-hidden">已隐藏</span>
                                    <span v-if="item.is_warp" class="mini-pill emerald">全包</span>
                                </td>
                                <td>
                                    <span v-if="item.level" class="level-tag">{{ item.level }}</span>
                                    <span v-else>-</span>
                                </td>
                                <td>
                                    <div class="col-tags">
                                        <span
                                            v-for="col in (item.collections || []).slice(0, 3)"
                                            :key="col"
                                            class="col-chip"
                                        >
                                            {{ col }}
                                        </span>
                                    </div>
                                </td>
                                <td class="col-dir" :title="item.directions">{{ item.directions || '-' }}</td>
                                <td class="col-pub">
                                    {{ item.publication_time || '待定' }}
                                    <template v-if="item.periodical_page"> · {{ formatPageText(item.periodical_page) }}</template>
                                </td>
                                <td class="col-price">
                                    {{ item.price && item.price !== '0' ? `¥${item.price}` : '待询' }}
                                </td>
                                <td class="col-ops">
                                    <span class="t-act" @click="emit('select-journal', item)">详情</span>
                                    <span class="t-act" @click="handleCopyDescription(item)">复制描述</span>
                                    <span v-if="!readonly" class="t-act" @click="handleTriggerModify(item)">修改</span>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <!-- 卡片内无刷新分页栏 -->
                <div v-if="totalCount > pageSize" class="card-pagination-bar">
                    <span class="page-info">
                        共 <strong>{{ totalCount }}</strong> 本匹配刊物 · 第 {{ currentPage }} / {{ totalPages }} 页
                    </span>
                    <div class="page-btns">
                        <el-button
                            size="small"
                            :disabled="currentPage <= 1 || pageLoading"
                            @click="handlePageChange(currentPage - 1)"
                        >
                            上一页
                        </el-button>
                        <el-button
                            size="small"
                            :disabled="currentPage >= totalPages || pageLoading"
                            @click="handlePageChange(currentPage + 1)"
                        >
                            下一页
                        </el-button>
                    </div>
                </div>
            </template>
            <div v-else class="empty-result">暂无匹配的期刊记录</div>
        </div>

        <!-- 2. 新增期刊入库确认单 (create_confirm，默认开启卡片内直接编辑 + 支持扩展字段增删) -->
        <div v-else-if="card.card_type === 'create_confirm' && editablePayload" class="create-confirm-body">
            <!-- 顶部：已自动匹配的系统分类字典展示条 -->
            <div v-if="hasMatchedCategories" class="matched-cats-panel">
                <span class="mc-title">🏷️ 已匹配分类字典：</span>
                <div class="mc-chips">
                    <span
                        v-for="cName in card.category_labels?.levels || []"
                        :key="`lvl-${cName}`"
                        class="cat-chip level"
                    >
                        <em>级别</em> {{ cName }}
                    </span>
                    <span
                        v-for="cName in card.category_labels?.cycle || []"
                        :key="`cyc-${cName}`"
                        class="cat-chip cycle"
                    >
                        <em>出版周期</em> {{ cName }}
                    </span>
                    <span
                        v-for="cName in card.category_labels?.collections || []"
                        :key="`col-${cName}`"
                        class="cat-chip col"
                    >
                        <em>收录</em> {{ cName }}
                    </span>
                    <span
                        v-for="cName in card.category_labels?.directions || []"
                        :key="`dir-${cName}`"
                        class="cat-chip dir"
                    >
                        <em>方向</em> {{ cName }}
                    </span>
                </div>
            </div>

            <div class="field-grid">
                <div class="field-item">
                    <span class="f-label">期刊名称 <em class="req-star">*</em></span>
                    <input
                        v-if="isInlineEditingCreate && card.status === 'pending' && !readonly"
                        v-model="editablePayload.name"
                        class="inline-field-input"
                        placeholder="如：《数字经济评论》"
                    />
                    <span v-else class="f-val highlight">{{ editablePayload.name || '-' }}</span>
                </div>
                <div class="field-item">
                    <span class="f-label">期刊发表费 (元)</span>
                    <input
                        v-if="isInlineEditingCreate && card.status === 'pending' && !readonly"
                        v-model="editablePayload.price"
                        class="inline-field-input"
                        placeholder="如：2600"
                    />
                    <span v-else class="f-val price">
                        {{ editablePayload.price ? `¥${editablePayload.price}` : '待询' }}
                    </span>
                </div>
                <div class="field-item">
                    <span class="f-label">出刊时间 (YYYY-MM-DD)</span>
                    <input
                        v-if="isInlineEditingCreate && card.status === 'pending' && !readonly"
                        v-model="editablePayload.publication_time"
                        class="inline-field-input"
                        placeholder="如：2025-09-01"
                    />
                    <span v-else class="f-val">{{ editablePayload.publication_time || '-' }}</span>
                </div>
                <div class="field-item">
                    <span class="f-label">期刊页码</span>
                    <input
                        v-if="isInlineEditingCreate && card.status === 'pending' && !readonly"
                        v-model="editablePayload.periodical_page"
                        class="inline-field-input"
                        placeholder="如：168"
                    />
                    <span v-else class="f-val">
                        {{ editablePayload.periodical_page ? `${editablePayload.periodical_page} 页` : '-' }}
                    </span>
                </div>
                <div class="field-item">
                    <span class="f-label">主管单位</span>
                    <input
                        v-if="isInlineEditingCreate && card.status === 'pending' && !readonly"
                        v-model="editablePayload.competent_unit"
                        class="inline-field-input"
                        placeholder="填写主管单位"
                    />
                    <span v-else class="f-val">{{ editablePayload.competent_unit || '-' }}</span>
                </div>
                <div class="field-item">
                    <span class="f-label">主办单位</span>
                    <input
                        v-if="isInlineEditingCreate && card.status === 'pending' && !readonly"
                        v-model="editablePayload.organization"
                        class="inline-field-input"
                        placeholder="填写主办单位"
                    />
                    <span v-else class="f-val">{{ editablePayload.organization || '-' }}</span>
                </div>
                <div class="field-item">
                    <span class="f-label">国内 / 国际刊号</span>
                    <div v-if="isInlineEditingCreate && card.status === 'pending' && !readonly" style="display: flex; gap: 6px;">
                        <input
                            v-model="editablePayload.domestic_number"
                            class="inline-field-input"
                            placeholder="CN刊号"
                        />
                        <input
                            v-model="editablePayload.international_number"
                            class="inline-field-input"
                            placeholder="ISSN刊号"
                        />
                    </div>
                    <span v-else class="f-val">
                        {{ editablePayload.domestic_number || '-' }} /
                        {{ editablePayload.international_number || '-' }}
                    </span>
                </div>
                <div class="field-item">
                    <span class="f-label">服务属性（点击切换）</span>
                    <div v-if="isInlineEditingCreate && card.status === 'pending' && !readonly" class="create-bool-toggles">
                        <button
                            type="button"
                            class="c-bool-chip"
                            :class="{ active: editablePayload.is_warp }"
                            @click="editablePayload.is_warp = !editablePayload.is_warp"
                        >
                            {{ editablePayload.is_warp ? '✓ 支持全包' : '非全包' }}
                        </button>
                        <button
                            type="button"
                            class="c-bool-chip"
                            :class="{ active: editablePayload.invoice_receipt }"
                            @click="editablePayload.invoice_receipt = !editablePayload.invoice_receipt"
                        >
                            {{ editablePayload.invoice_receipt ? '✓ 可开发票' : '不开票' }}
                        </button>
                        <button
                            type="button"
                            class="c-bool-chip"
                            :class="{ active: editablePayload.works }"
                            @click="editablePayload.works = !editablePayload.works"
                        >
                            {{ editablePayload.works ? '✓ 可发作品' : '不可发作品' }}
                        </button>
                        <button
                            type="button"
                            class="c-bool-chip"
                            :class="{ active: editablePayload.is_hidden }"
                            @click="editablePayload.is_hidden = !editablePayload.is_hidden"
                        >
                            {{ editablePayload.is_hidden ? '🙈 已隐藏' : '正常显示' }}
                        </button>
                    </div>
                    <span v-else class="f-val">
                        {{ editablePayload.is_warp ? '全包' : '非全包' }} ·
                        {{ editablePayload.invoice_receipt ? '可开杂志社发票' : '不开票' }} ·
                        {{ editablePayload.works ? '可发作品' : '不可发作品' }} ·
                        {{ editablePayload.is_hidden ? '已隐藏' : '正常显示' }}
                    </span>
                </div>
            </div>

            <!-- 用户补充的扩展入库字段列表（支持行内编辑与单项 × 移除） -->
            <div v-if="createExtraFields.length > 0" class="create-extra-list">
                <div v-for="(ef, eIdx) in createExtraFields" :key="ef.field" class="create-extra-row">
                    <div class="ce-header">
                        <span class="ce-label">{{ ef.label }}</span>
                        <button
                            v-if="card.status === 'pending' && !readonly"
                            type="button"
                            class="ce-remove-btn"
                            title="移除此扩展字段"
                            @click="removeCreateExtraField(eIdx)"
                        >
                            × 移除
                        </button>
                    </div>
                    <textarea
                        v-if="isLongTextField(ef.field) && isInlineEditingCreate && card.status === 'pending' && !readonly"
                        v-model="editablePayload[ef.field]"
                        class="ce-textarea"
                        rows="2"
                        :placeholder="`请输入${ef.label}...`"
                    ></textarea>
                    <input
                        v-else-if="isInlineEditingCreate && card.status === 'pending' && !readonly"
                        v-model="editablePayload[ef.field]"
                        class="inline-field-input"
                        :placeholder="`请输入${ef.label}...`"
                    />
                    <div v-else class="ce-readonly-val">{{ editablePayload[ef.field] || '-' }}</div>
                </div>
            </div>

            <!-- 补充更多入库字段工具栏（渐进式披露胶囊栏） -->
            <div v-if="card.status === 'pending' && !readonly && availableCreateExtraCatalog.length > 0" class="extra-fields-bar">
                <span class="extra-hint">
                    <span class="hint-dot">+</span>
                    <span>补充入库字段</span>
                </span>
                <div class="extra-chips-wrap">
                    <button
                        v-for="cf in availableCreateExtraCatalog.slice(0, 5)"
                        :key="cf.field"
                        type="button"
                        class="add-field-chip"
                        @click="addCreateExtraField(cf)"
                    >
                        <span class="chip-plus">+</span>
                        <span>{{ cf.label }}</span>
                    </button>
                    <el-dropdown
                        v-if="availableCreateExtraCatalog.length > 5"
                        trigger="click"
                        @command="addCreateExtraField"
                    >
                        <button type="button" class="add-field-chip more-dropdown-chip">
                            <span>更多字段 ({{ availableCreateExtraCatalog.length - 5 }})</span>
                            <el-icon style="margin-left: 3px"><ArrowDown /></el-icon>
                        </button>
                        <template #dropdown>
                            <el-dropdown-menu>
                                <el-dropdown-item
                                    v-for="mf in availableCreateExtraCatalog.slice(5)"
                                    :key="mf.field"
                                    :command="mf"
                                >
                                    + {{ mf.label }}
                                </el-dropdown-item>
                            </el-dropdown-menu>
                        </template>
                    </el-dropdown>
                </div>
            </div>
        </div>

        <!-- 3. 修改期刊 Diff 对比确认单 (update_confirm，同时支持单本微调与多本批量修改) -->
        <div v-else-if="card.card_type === 'update_confirm'" class="update-confirm-body">
            <!-- 多本批量修改视图 (editableBatchTargets.length > 1) -->
            <template v-if="isBatchUpdate">
                <div class="batch-targets-list">
                    <div
                        v-for="(bt, tIdx) in editableBatchTargets"
                        :key="bt.periodical_id"
                        class="batch-target-section"
                    >
                        <div class="bt-header">
                            <div class="bt-left">
                                <span class="bt-index">{{ String(tIdx + 1).padStart(2, '0') }}</span>
                                <span class="bt-name">{{ bt.periodical_name }}</span>
                                <span v-if="bt.level" class="bt-level">{{ bt.level }}</span>
                                <span class="bt-count-pill">{{ bt.diffs.length }} 项变更</span>
                            </div>
                            <div v-if="card.status === 'pending' && !readonly" class="bt-right">
                                <button
                                    v-if="editableBatchTargets.length > 1"
                                    type="button"
                                    class="bt-remove-journal-btn"
                                    title="从本次批量修改中剔除此期刊"
                                    @click="removeBatchTarget(tIdx)"
                                >
                                    × 不改此刊
                                </button>
                            </div>
                        </div>

                        <div v-if="bt.diffs.length > 0" class="diff-list">
                            <div
                                v-for="(d, idx) in bt.diffs"
                                :key="`${bt.periodical_id}-${d.field}`"
                                class="diff-item-card"
                                :class="{
                                    'is-longtext': isLongTextField(d.field),
                                    'is-changed': isDiffChanged(d),
                                }"
                            >
                                <template v-if="isLongTextField(d.field)">
                                    <div class="longtext-top">
                                        <div class="lt-left">
                                            <span class="diff-label-badge">{{ d.label }}</span>
                                            <span class="lt-old-summary" :title="d.old_value">
                                                原值：{{ truncateText(d.old_value, 42) }}
                                            </span>
                                        </div>
                                        <div v-if="card.status === 'pending' && !readonly" class="diff-row-actions">
                                            <button
                                                type="button"
                                                class="row-act-btn sync-all"
                                                title="将该新值同步应用到本次批量修改的所有期刊"
                                                @click="syncFieldToAllBatchTargets(d.field, d.label, d.new_value)"
                                            >
                                                ⇄ 同步至全部({{ editableBatchTargets.length }}本)
                                            </button>
                                            <button
                                                v-if="isDiffChanged(d)"
                                                type="button"
                                                class="row-act-btn reset"
                                                @click="resetBatchTargetDiffField(tIdx, d)"
                                            >
                                                ↺ 重置
                                            </button>
                                            <button
                                                v-if="bt.diffs.length > 1"
                                                type="button"
                                                class="row-act-btn remove"
                                                @click="removeBatchTargetDiffField(tIdx, idx)"
                                            >
                                                × 移除
                                            </button>
                                        </div>
                                    </div>
                                    <div class="longtext-editor-wrap">
                                        <textarea
                                            v-if="card.status === 'pending' && !readonly"
                                            v-model="d.new_value"
                                            rows="2"
                                            class="inline-diff-textarea"
                                            :placeholder="`请输入新的${d.label}内容...`"
                                            @input="syncBatchTargetDiffToPayload(tIdx, d.field, d.new_value)"
                                        ></textarea>
                                        <div v-else class="new-val-block">{{ d.new_value }}</div>
                                    </div>
                                </template>

                                <template v-else>
                                    <div class="compact-left">
                                        <span class="diff-label-badge">{{ d.label }}</span>
                                        <span
                                            v-if="card.status === 'pending' && !readonly"
                                            class="change-state-dot"
                                            :class="isDiffChanged(d) ? 'changed' : 'same'"
                                        ></span>
                                    </div>

                                    <div class="compact-middle">
                                        <span class="old-val-pill" :title="`修改前原值：${d.old_value}`">
                                            原值: {{ truncateText(d.old_value, 24) }}
                                        </span>
                                        <el-icon class="diff-arrow"><Right /></el-icon>

                                        <template v-if="card.status === 'pending' && !readonly">
                                            <div v-if="isBoolField(d.field)" class="bool-switch-group">
                                                <button
                                                    type="button"
                                                    class="bool-btn"
                                                    :class="{ active: d.new_value === '是' }"
                                                    @click="setBatchTargetBoolDiff(tIdx, d, true)"
                                                >
                                                    是
                                                </button>
                                                <button
                                                    type="button"
                                                    class="bool-btn"
                                                    :class="{ active: d.new_value !== '是' }"
                                                    @click="setBatchTargetBoolDiff(tIdx, d, false)"
                                                >
                                                    否
                                                </button>
                                            </div>
                                            <input
                                                v-else
                                                v-model="d.new_value"
                                                class="inline-diff-input"
                                                :placeholder="`请输入新的${d.label}`"
                                                @input="syncBatchTargetDiffToPayload(tIdx, d.field, d.new_value)"
                                            />
                                        </template>
                                        <template v-else>
                                            <span class="new-val">{{ d.new_value }}</span>
                                            <span v-if="card.status === 'confirmed'" class="confirmed-effect-tag">
                                                ✓ 已生效
                                            </span>
                                        </template>
                                    </div>

                                    <div v-if="card.status === 'pending' && !readonly" class="diff-row-actions">
                                        <button
                                            type="button"
                                            class="row-act-btn sync-all"
                                            title="将该新值一键同步到全部待改期刊"
                                            @click="syncFieldToAllBatchTargets(d.field, d.label, d.new_value)"
                                        >
                                            ⇄ 同步全部
                                        </button>
                                        <button
                                            v-if="isDiffChanged(d)"
                                            type="button"
                                            class="row-act-btn reset"
                                            title="恢复为修改前原值"
                                            @click="resetBatchTargetDiffField(tIdx, d)"
                                        >
                                            ↺ 重置
                                        </button>
                                        <button
                                            v-if="bt.diffs.length > 1"
                                            type="button"
                                            class="row-act-btn remove"
                                            title="从该期刊修改项中移除此字段"
                                            @click="removeBatchTargetDiffField(tIdx, idx)"
                                        >
                                            × 移除
                                        </button>
                                    </div>
                                </template>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- 批量为全部待改期刊追加修改字段 -->
                <div v-if="card.status === 'pending' && !readonly" class="extra-fields-bar">
                    <span class="extra-hint">
                        <span class="hint-dot">+</span>
                        <span>统一添加修改字段</span>
                    </span>
                    <div class="extra-chips-wrap">
                        <button
                            v-for="ef in primaryFieldCatalog"
                            :key="ef.field"
                            type="button"
                            class="add-field-chip"
                            @click="addBatchExtraDiffField(ef)"
                        >
                            <span class="chip-plus">+</span>
                            <span>{{ ef.label }}</span>
                        </button>

                        <el-dropdown
                            v-if="moreFieldCatalog.length > 0"
                            trigger="click"
                            @command="addBatchExtraDiffField"
                        >
                            <button type="button" class="add-field-chip more-dropdown-chip">
                                <span>更多字段 ({{ moreFieldCatalog.length }})</span>
                                <el-icon style="margin-left: 3px"><ArrowDown /></el-icon>
                            </button>
                            <template #dropdown>
                                <el-dropdown-menu>
                                    <el-dropdown-item
                                        v-for="mf in moreFieldCatalog"
                                        :key="mf.field"
                                        :command="mf"
                                    >
                                        + {{ mf.label }}
                                    </el-dropdown-item>
                                </el-dropdown-menu>
                            </template>
                        </el-dropdown>
                    </div>
                </div>
            </template>

            <!-- 单本修改视图 (editableBatchTargets.length <= 1) -->
            <template v-else>
                <div v-if="editableDiffs.length > 0" class="diff-list">
                    <div
                        v-for="(d, idx) in editableDiffs"
                        :key="d.field"
                        class="diff-item-card"
                        :class="{
                            'is-longtext': isLongTextField(d.field),
                            'is-changed': isDiffChanged(d),
                        }"
                    >
                        <!-- A. 长文本字段（栏目设置、注意事项、期刊描述、发表流程等）：上下宽幅多行编辑布局 -->
                        <template v-if="isLongTextField(d.field)">
                            <div class="longtext-top">
                                <div class="lt-left">
                                    <span class="diff-label-badge">{{ d.label }}</span>
                                    <span
                                        v-if="card.status === 'pending' && !readonly"
                                        class="change-state-tag"
                                        :class="isDiffChanged(d) ? 'changed' : 'same'"
                                    >
                                        {{ isDiffChanged(d) ? '已修改' : '未改动' }}
                                    </span>
                                    <span class="lt-old-summary" :title="d.old_value">
                                        原值：{{ truncateText(d.old_value, 42) }}
                                    </span>
                                    <button
                                        v-if="d.old_value && d.old_value.length > 42"
                                        type="button"
                                        class="toggle-old-btn"
                                        @click="toggleExpandOld(d.field)"
                                    >
                                        {{ expandedOldMap[d.field] ? '收起原值' : '展开完整原值' }}
                                    </button>
                                </div>
                                <div v-if="card.status === 'pending' && !readonly" class="diff-row-actions">
                                    <button
                                        v-if="isDiffChanged(d)"
                                        type="button"
                                        class="row-act-btn reset"
                                        title="恢复为修改前原值"
                                        @click="resetDiffField(d)"
                                    >
                                        ↺ 重置原值
                                    </button>
                                    <button
                                        type="button"
                                        class="row-act-btn remove"
                                        title="从本次修改单中移除该字段"
                                        @click="removeDiffField(idx)"
                                    >
                                        × 移除
                                    </button>
                                </div>
                            </div>

                            <div v-if="expandedOldMap[d.field]" class="lt-old-full-box">
                                <div class="old-box-label">修改前完整原值：</div>
                                <div class="old-box-content">{{ d.old_value }}</div>
                            </div>

                            <div class="longtext-editor-wrap">
                                <textarea
                                    v-if="card.status === 'pending' && !readonly"
                                    v-model="d.new_value"
                                    rows="3"
                                    class="inline-diff-textarea"
                                    :placeholder="`请输入新的${d.label}内容...`"
                                    @input="syncDiffToPayload(d.field, d.new_value)"
                                ></textarea>
                                <div v-else class="new-val-block">{{ d.new_value }}</div>
                            </div>
                        </template>

                        <!-- B. 短属性与布尔字段（出刊时间、发表费、页码、引用率等）：左中右自适应对齐布局 -->
                        <template v-else>
                            <div class="compact-left">
                                <span class="diff-label-badge">{{ d.label }}</span>
                                <span
                                    v-if="card.status === 'pending' && !readonly"
                                    class="change-state-dot"
                                    :class="isDiffChanged(d) ? 'changed' : 'same'"
                                    :title="isDiffChanged(d) ? '已修改新值' : '与原值一致'"
                                ></span>
                            </div>

                            <div class="compact-middle">
                                <span class="old-val-pill" :title="`修改前原值：${d.old_value}`">
                                    原值: {{ truncateText(d.old_value, 24) }}
                                </span>
                                <el-icon class="diff-arrow"><Right /></el-icon>

                                <template v-if="card.status === 'pending' && !readonly">
                                    <!-- 布尔开关字段 -->
                                    <div v-if="isBoolField(d.field)" class="bool-switch-group">
                                        <button
                                            type="button"
                                            class="bool-btn"
                                            :class="{ active: d.new_value === '是' }"
                                            @click="setBoolDiff(d, true)"
                                        >
                                            是
                                        </button>
                                        <button
                                            type="button"
                                            class="bool-btn"
                                            :class="{ active: d.new_value !== '是' }"
                                            @click="setBoolDiff(d, false)"
                                        >
                                            否
                                        </button>
                                    </div>
                                    <!-- 常规单行自适应宽输入框 -->
                                    <input
                                        v-else
                                        v-model="d.new_value"
                                        class="inline-diff-input"
                                        :placeholder="`请输入新的${d.label}`"
                                        @input="syncDiffToPayload(d.field, d.new_value)"
                                    />
                                </template>
                                <template v-else>
                                    <span class="new-val">{{ d.new_value }}</span>
                                    <span v-if="card.status === 'confirmed'" class="confirmed-effect-tag">✓ 已生效</span>
                                </template>
                            </div>

                            <div v-if="card.status === 'pending' && !readonly" class="diff-row-actions">
                                <button
                                    v-if="isDiffChanged(d)"
                                    type="button"
                                    class="row-act-btn reset"
                                    title="恢复为修改前原值"
                                    @click="resetDiffField(d)"
                                >
                                    ↺ 重置
                                </button>
                                <button
                                    type="button"
                                    class="row-act-btn remove"
                                    title="从本次修改单中移除该字段"
                                    @click="removeDiffField(idx)"
                                >
                                    × 移除
                                </button>
                            </div>
                        </template>
                    </div>
                </div>

                <!-- 当用户移除了全部字段时的空状态引导 -->
                <div v-else class="empty-diffs-box">
                    <span>已清空所有待修改项，请从下方点选您需要修改的期刊字段：</span>
                </div>

                <!-- 快捷补充其他字段（常用快捷胶囊 + 全量 22 个字段下拉选择） -->
                <div v-if="card.status === 'pending' && !readonly" class="extra-fields-bar">
                    <span class="extra-hint">
                        <span class="hint-dot">+</span>
                        <span>添加修改字段</span>
                    </span>
                    <div class="extra-chips-wrap">
                        <button
                            v-for="ef in availablePrimaryExtraFields"
                            :key="ef.field"
                            type="button"
                            class="add-field-chip"
                            @click="addExtraDiffField(ef)"
                        >
                            <span class="chip-plus">+</span>
                            <span>{{ ef.label }}</span>
                        </button>

                        <el-dropdown
                            v-if="availableMoreExtraFields.length > 0"
                            trigger="click"
                            @command="addExtraDiffField"
                        >
                            <button type="button" class="add-field-chip more-dropdown-chip">
                                <span>更多字段 ({{ availableMoreExtraFields.length }})</span>
                                <el-icon style="margin-left: 3px"><ArrowDown /></el-icon>
                            </button>
                            <template #dropdown>
                                <el-dropdown-menu>
                                    <el-dropdown-item
                                        v-for="mf in availableMoreExtraFields"
                                        :key="mf.field"
                                        :command="mf"
                                    >
                                        + {{ mf.label }}
                                    </el-dropdown-item>
                                </el-dropdown-menu>
                            </template>
                        </el-dropdown>
                    </div>
                </div>
            </template>
        </div>

        <!-- 4. 删除期刊高危确认单 (delete_confirm) -->
        <div v-else-if="card.card_type === 'delete_confirm'" class="delete-target-box-wrap delete-confirm-body">
            <div v-if="card.items && card.items.length > 0" class="delete-target-box">
                <div class="target-main">
                    <strong>{{ card.items[0].name }}</strong>
                    <span v-if="card.items[0].level" class="j-level">{{ card.items[0].level }}</span>
                </div>
                <div class="target-sub">
                    ID: {{ card.periodical_id }} · 出刊时间: {{ card.items[0].publication_time || '-' }} · 发表费:
                    {{ card.items[0].price ? `¥${card.items[0].price}` : '-' }}
                </div>
            </div>
        </div>

        <!-- 底部操作栏（待确认状态且非只读审计模式） -->
        <div v-if="card.status === 'pending' && !readonly" class="card-actions">
            <el-button size="small" :disabled="submitting" @click="handleDecision('cancel')">取消</el-button>
            <el-button
                v-if="card.card_type === 'create_confirm' && showFillDrawer"
                size="small"
                type="primary"
                plain
                icon="EditPen"
                @click="emit('fill-drawer', editablePayload)"
            >
                填入新增抽屉微调
            </el-button>
            <el-button
                size="small"
                :type="card.card_type === 'delete_confirm' ? 'danger' : 'primary'"
                :loading="submitting"
                :disabled="
                    card.card_type === 'update_confirm' &&
                    (isBatchUpdate ? totalBatchDiffCount === 0 : editableDiffs.length === 0)
                "
                icon="Check"
                @click="handleDecision('confirm')"
            >
                {{
                    card.card_type === 'create_confirm'
                        ? '确认直接入库'
                        : card.card_type === 'update_confirm'
                          ? isBatchUpdate
                              ? `确认批量修改 (${editableBatchTargets.length}本 · 共${totalBatchDiffCount}项)`
                              : `确认执行修改 (${editableDiffs.length}项)`
                          : '确认删除'
                }}
            </el-button>
        </div>

        <!-- 已执行结果与一体化最新期刊快照（支持单本或多本批量展示） -->
        <div v-else-if="card.result_message || (card.status === 'confirmed' && card.items && card.items.length > 0)" class="confirmed-integrated-footer">
            <div v-if="card.result_message" class="card-result-banner" :class="card.status">
                <span>{{ card.result_message }}</span>
            </div>

            <div
                v-if="card.status === 'confirmed' && card.items && card.items.length > 0"
                class="embedded-snapshot-list"
            >
                <div
                    v-for="snapItem in card.items"
                    :key="snapItem.periodical_id"
                    class="embedded-snapshot-card"
                    @click="emit('select-journal', snapItem)"
                >
                    <div class="snap-left">
                        <span class="snap-tag">最新档案</span>
                        <strong class="snap-name">{{ snapItem.name }}</strong>
                        <span v-if="snapItem.is_hidden" class="b-tag amber-hidden" style="padding: 1px 6px; border-radius: 4px; font-size: 11px;">已隐藏</span>
                        <span v-if="snapItem.level" class="snap-level">{{ snapItem.level }}</span>
                        <span v-for="col in (snapItem.collections || []).slice(0, 3)" :key="col" class="snap-col">
                            {{ col }}
                        </span>
                        <span class="snap-meta">
                            出刊: <strong>{{ snapItem.publication_time || '待定' }}</strong>
                            <template v-if="snapItem.periodical_page"> · {{ formatPageText(snapItem.periodical_page) }}</template>
                            · 发表费:
                            <strong class="snap-price">
                                {{ snapItem.price && snapItem.price !== '0' ? `¥${snapItem.price}` : '待询' }}
                            </strong>
                        </span>
                    </div>
                    <div class="snap-actions" @click.stop>
                        <span class="s-act" @click="emit('select-journal', snapItem)">查看完整详情</span>
                        <span class="s-sep">·</span>
                        <span class="s-act" @click="handleCopyDescription(snapItem)">复制描述</span>
                        <template v-if="!readonly">
                            <span class="s-sep">·</span>
                            <span class="s-act" @click="handleTriggerModify(snapItem)">再次修改</span>
                        </template>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { reqPeriodicalList, reqPeriodicalDetail } from '@/api/periodical'
import { copyTextToClipboard } from '@/utils/clipboard'
import type { ActionCardData, BatchUpdateTarget, FieldDiffItem, PeriodicalBriefCard } from '@/api/ai'

const props = withDefaults(
    defineProps<{
        card?: ActionCardData | null
        messageId: number
        readonly?: boolean
        showSyncFilter?: boolean
        showFillDrawer?: boolean
        showOpenCanvas?: boolean
    }>(),
    {
        card: null,
        readonly: false,
        showSyncFilter: false,
        showFillDrawer: false,
        showOpenCanvas: false,
    },
)

const emit = defineEmits<{
    (
        e: 'confirm-action',
        payload: {
            messageId: number
            decision: 'confirm' | 'cancel'
            formPayload?: Record<string, any>
            diffs?: FieldDiffItem[]
            batchUpdates?: BatchUpdateTarget[]
        },
    ): void
    (e: 'sync-filter', filterPayload: Record<string, any>): void
    (e: 'fill-drawer', formPayload?: Record<string, any>): void
    (e: 'select-journal', item: PeriodicalBriefCard): void
    (e: 'open-canvas', card: ActionCardData): void
    (e: 'quick-ask', promptText: string): void
    (e: 'start-modify', item: PeriodicalBriefCard): void
}>()

const submitting = ref(false)
const viewMode = ref<'grid' | 'table'>('grid')
const currentPage = ref(1)
const pageSize = 8
const pageLoading = ref(false)
const pagedItems = ref<PeriodicalBriefCard[] | null>(null)

// 行内可编辑状态（默认直接开启编辑，支持用户在确认卡片上直接微调、添加或移除字段，无需反复发消息）
const isInlineEditingCreate = ref(true)
const editablePayload = ref<Record<string, any>>({})
const originalValuesMap = ref<Record<string, any>>({})
const editableDiffs = ref<FieldDiffItem[]>([])
const editableBatchTargets = ref<BatchUpdateTarget[]>([])
const batchOriginalMap = ref<Record<number, Record<string, any>>>({})
const createExtraFields = ref<{ field: string; label: string }[]>([])
const expandedOldMap = reactive<Record<string, boolean>>({})

const isBatchUpdate = computed(() => editableBatchTargets.value.length > 1)
const totalBatchDiffCount = computed(() =>
    editableBatchTargets.value.reduce((acc, bt) => acc + (bt.diffs ? bt.diffs.length : 0), 0),
)

const CREATE_EXTRA_CATALOG: { field: string; label: string }[] = [
    { field: 'audit_time', label: '审核周期' },
    { field: 'column_setting', label: '栏目设置' },
    { field: 'periodical_description', label: '期刊描述' },
    { field: 'attention_matter', label: '文章注意事项' },
    { field: 'founding_time', label: '创刊年份' },
    { field: 'citation_rate', label: '引用率' },
    { field: 'composite_influence_factor', label: '复合影响因子' },
    { field: 'synthetic_influence_factor', label: '综合影响因子' },
    { field: 'author_info', label: '作者与样刊说明' },
    { field: 'periodical_final_proof', label: '期刊清样' },
    { field: 'article_naming', label: '文章命名规范' },
    { field: 'publication_process', label: '发表流程' },
    { field: 'all_inclusive_process', label: '全包流程' },
    { field: 'internal_processes', label: '内部流程' },
]

const hasMatchedCategories = computed(() => {
    const labels = props.card?.category_labels
    if (!labels) return false
    return Object.values(labels).some((arr) => Array.isArray(arr) && arr.length > 0)
})

const availableCreateExtraCatalog = computed(() =>
    CREATE_EXTRA_CATALOG.filter((cf) => !createExtraFields.value.some((it) => it.field === cf.field)),
)

const addCreateExtraField = (cf: { field: string; label: string }) => {
    if (createExtraFields.value.some((it) => it.field === cf.field)) return
    createExtraFields.value.push(cf)
    if (editablePayload.value[cf.field] === undefined) {
        editablePayload.value[cf.field] = ''
    }
}

const removeCreateExtraField = (idx: number) => {
    const removed = createExtraFields.value[idx]
    if (!removed) return
    createExtraFields.value.splice(idx, 1)
    editablePayload.value[removed.field] = ''
}

const LONG_TEXT_FIELDS = new Set([
    'column_setting',
    'attention_matter',
    'periodical_description',
    'author_info',
    'periodical_final_proof',
    'article_naming',
    'publication_process',
    'all_inclusive_process',
    'internal_processes',
])

const BOOL_FIELDS = new Set(['is_warp', 'invoice_receipt', 'works', 'is_hidden'])

const isLongTextField = (field: string) => LONG_TEXT_FIELDS.has(field)
const isBoolField = (field: string) => BOOL_FIELDS.has(field)

const initEditableState = () => {
    const rawPayload = props.card?.form_payload ? JSON.parse(JSON.stringify(props.card.form_payload)) : {}
    const rawDiffs: FieldDiffItem[] = props.card?.diffs ? JSON.parse(JSON.stringify(props.card.diffs)) : []
    const rawBatch: BatchUpdateTarget[] = props.card?.batch_updates
        ? JSON.parse(JSON.stringify(props.card.batch_updates))
        : []

    // 构造精确的数据库修改前原值快照，确保用户移除某项时能100%还原为原值
    const origMap = JSON.parse(JSON.stringify(rawPayload))
    rawDiffs.forEach((d) => {
        if (BOOL_FIELDS.has(d.field)) {
            origMap[d.field] = d.old_value === '是'
        } else {
            origMap[d.field] = d.old_value === '(空)' ? '' : d.old_value
        }
    })

    const bOrigMap: Record<number, Record<string, any>> = {}
    rawBatch.forEach((bt) => {
        const singleOrig = JSON.parse(JSON.stringify(bt.form_payload || {}))
        ;(bt.diffs || []).forEach((d) => {
            if (BOOL_FIELDS.has(d.field)) {
                singleOrig[d.field] = d.old_value === '是'
            } else {
                singleOrig[d.field] = d.old_value === '(空)' ? '' : d.old_value
            }
        })
        bOrigMap[bt.periodical_id] = singleOrig
    })

    // 若新增卡片中的 periodical_description 仅是极短的“出版周期：xxx”等重复碎片，自动清理以免弹出无意义的描述框
    if (props.card?.card_type === 'create_confirm' && typeof rawPayload.periodical_description === 'string') {
        const descClean = rawPayload.periodical_description.trim()
        if (descClean.length <= 16 && /^(出版周期|级别|网站收录|收稿方向)[：:]/.test(descClean)) {
            rawPayload.periodical_description = ''
        }
    }

    originalValuesMap.value = origMap
    editablePayload.value = rawPayload
    editableDiffs.value = rawDiffs
    editableBatchTargets.value = rawBatch
    batchOriginalMap.value = bOrigMap

    // 自动将新增卡片中已经有值的扩展字段展示出来
    createExtraFields.value = CREATE_EXTRA_CATALOG.filter(
        (cf) => rawPayload[cf.field] !== undefined && rawPayload[cf.field] !== null && String(rawPayload[cf.field]).trim() !== '',
    )
}

initEditableState()

watch(
    () => props.card,
    () => {
        currentPage.value = 1
        pagedItems.value = null
        initEditableState()
    },
    { deep: true },
)

const isDiffChanged = (d: FieldDiffItem) => {
    const normOld = d.old_value === '(空)' ? '' : String(d.old_value || '').trim()
    const normNew = d.new_value === '(空)' ? '' : String(d.new_value || '').trim()
    return normOld !== normNew
}

const changedDiffCount = computed(() => editableDiffs.value.filter((d) => isDiffChanged(d)).length)

const truncateText = (text: string, maxLen: number) => {
    if (!text) return '(空)'
    const clean = String(text).replace(/\s+/g, ' ').trim()
    if (clean.length <= maxLen) return clean
    return clean.slice(0, maxLen) + '...'
}

const toggleExpandOld = (field: string) => {
    expandedOldMap[field] = !expandedOldMap[field]
}

const syncDiffToPayload = (field: string, val: string) => {
    if (BOOL_FIELDS.has(field)) {
        editablePayload.value[field] = val === '是'
    } else {
        editablePayload.value[field] = val
    }
}

const setBoolDiff = (d: FieldDiffItem, boolVal: boolean) => {
    d.new_value = boolVal ? '是' : '否'
    editablePayload.value[d.field] = boolVal
}

// 恢复某一项为修改前原值
const resetDiffField = (d: FieldDiffItem) => {
    if (BOOL_FIELDS.has(d.field)) {
        const origBool = d.old_value === '是'
        d.new_value = origBool ? '是' : '否'
        editablePayload.value[d.field] = origBool
    } else {
        const origStr = d.old_value === '(空)' ? '' : d.old_value
        d.new_value = origStr
        editablePayload.value[d.field] = origStr
    }
}

// 从修改确认单中移除某一项，并将 payload 中该字段还原为数据库初始原值
const removeDiffField = (idx: number) => {
    const removed = editableDiffs.value[idx]
    if (!removed) return
    editableDiffs.value.splice(idx, 1)
    editablePayload.value[removed.field] = originalValuesMap.value[removed.field]
}

// 一键清理所有未改动的预置字段
const clearUnchangedFields = () => {
    const kept: FieldDiffItem[] = []
    editableDiffs.value.forEach((d) => {
        if (isDiffChanged(d)) {
            kept.push(d)
        } else {
            editablePayload.value[d.field] = originalValuesMap.value[d.field]
        }
    })
    editableDiffs.value = kept
}

const primaryFieldCatalog = [
    { field: 'publication_time', label: '出刊时间' },
    { field: 'price', label: '期刊发表费' },
    { field: 'periodical_page', label: '期刊页码' },
    { field: 'is_hidden', label: '是否隐藏' },
    { field: 'level', label: '期刊级别' },
    { field: 'cycle', label: '出版周期' },
    { field: 'audit_time', label: '审核周期' },
    { field: 'column_setting', label: '栏目设置' },
    { field: 'attention_matter', label: '注意事项' },
]

const moreFieldCatalog = [
    { field: 'collections', label: '网站收录' },
    { field: 'directions', label: '收稿方向' },
    { field: 'periodical_description', label: '期刊描述' },
    { field: 'citation_rate', label: '引用率' },
    { field: 'competent_unit', label: '主管单位' },
    { field: 'organization', label: '主办单位' },
    { field: 'domestic_number', label: '国内刊号 (CN)' },
    { field: 'international_number', label: '国际刊号 (ISSN)' },
    { field: 'founding_time', label: '创刊时间' },
    { field: 'composite_influence_factor', label: '复合影响因子' },
    { field: 'synthetic_influence_factor', label: '综合影响因子' },
    { field: 'author_info', label: '作者与样刊说明' },
    { field: 'periodical_final_proof', label: '期刊清样' },
    { field: 'article_naming', label: '文章命名' },
    { field: 'publication_process', label: '发表流程' },
    { field: 'all_inclusive_process', label: '全包流程' },
    { field: 'internal_processes', label: '内部流程' },
    { field: 'is_warp', label: '是否全包' },
    { field: 'invoice_receipt', label: '可开发票' },
    { field: 'works', label: '可发作品' },
]

const availablePrimaryExtraFields = computed(() =>
    primaryFieldCatalog.filter((ef) => !editableDiffs.value.some((d) => d.field === ef.field)),
)

const availableMoreExtraFields = computed(() =>
    moreFieldCatalog.filter((ef) => !editableDiffs.value.some((d) => d.field === ef.field)),
)

const addExtraDiffField = (ef: { field: string; label: string }) => {
    if (editableDiffs.value.some((d) => d.field === ef.field)) return
    const rawOld = originalValuesMap.value[ef.field]
    if (BOOL_FIELDS.has(ef.field)) {
        const boolOld = Boolean(rawOld)
        const nextBool = !boolOld
        editableDiffs.value.push({
            field: ef.field,
            label: ef.label,
            old_value: boolOld ? '是' : '否',
            new_value: nextBool ? '是' : '否',
        })
        editablePayload.value[ef.field] = nextBool
        return
    }
    const oldStr = rawOld !== undefined && rawOld !== null && String(rawOld) !== '' ? String(rawOld) : '(空)'
    const defaultNew = rawOld !== undefined && rawOld !== null && String(rawOld) !== '' ? String(rawOld) : ''
    editableDiffs.value.push({
        field: ef.field,
        label: ef.label,
        old_value: oldStr,
        new_value: defaultNew,
    })
}

// ===== 批量多本期刊修改专属操作方法 =====

const removeBatchTarget = (tIdx: number) => {
    if (editableBatchTargets.value.length <= 1) return
    editableBatchTargets.value.splice(tIdx, 1)
    if (editableBatchTargets.value.length === 1) {
        const remaining = editableBatchTargets.value[0]
        editablePayload.value = remaining.form_payload
        editableDiffs.value = remaining.diffs
        originalValuesMap.value = batchOriginalMap.value[remaining.periodical_id] || {}
    }
}

const syncBatchTargetDiffToPayload = (tIdx: number, field: string, val: string) => {
    const bt = editableBatchTargets.value[tIdx]
    if (!bt || !bt.form_payload) return
    if (BOOL_FIELDS.has(field)) {
        bt.form_payload[field] = val === '是'
    } else {
        bt.form_payload[field] = val
    }
}

const setBatchTargetBoolDiff = (tIdx: number, d: FieldDiffItem, boolVal: boolean) => {
    const bt = editableBatchTargets.value[tIdx]
    if (!bt || !bt.form_payload) return
    d.new_value = boolVal ? '是' : '否'
    bt.form_payload[d.field] = boolVal
}

const resetBatchTargetDiffField = (tIdx: number, d: FieldDiffItem) => {
    const bt = editableBatchTargets.value[tIdx]
    if (!bt || !bt.form_payload) return
    if (BOOL_FIELDS.has(d.field)) {
        const origBool = d.old_value === '是'
        d.new_value = origBool ? '是' : '否'
        bt.form_payload[d.field] = origBool
    } else {
        const origStr = d.old_value === '(空)' ? '' : d.old_value
        d.new_value = origStr
        bt.form_payload[d.field] = origStr
    }
}

const removeBatchTargetDiffField = (tIdx: number, idx: number) => {
    const bt = editableBatchTargets.value[tIdx]
    if (!bt || bt.diffs.length <= 1) return
    const removed = bt.diffs[idx]
    if (!removed) return
    bt.diffs.splice(idx, 1)
    const origForJournal = batchOriginalMap.value[bt.periodical_id]
    if (origForJournal && bt.form_payload) {
        bt.form_payload[removed.field] = origForJournal[removed.field]
    }
}

const syncFieldToAllBatchTargets = (field: string, label: string, newVal: string) => {
    editableBatchTargets.value.forEach((bt) => {
        if (!bt.form_payload) bt.form_payload = {}
        const existingDiff = bt.diffs.find((d) => d.field === field)
        if (existingDiff) {
            existingDiff.new_value = newVal
        } else {
            const origForJournal = batchOriginalMap.value[bt.periodical_id] || bt.form_payload
            const rawOld = origForJournal[field]
            const oldStr = BOOL_FIELDS.has(field)
                ? rawOld
                    ? '是'
                    : '否'
                : rawOld !== undefined && rawOld !== null && String(rawOld) !== ''
                  ? String(rawOld)
                  : '(空)'
            bt.diffs.push({
                field,
                label,
                old_value: oldStr,
                new_value: newVal,
            })
        }
        if (BOOL_FIELDS.has(field)) {
            bt.form_payload[field] = newVal === '是'
        } else {
            bt.form_payload[field] = newVal
        }
    })
    ElMessage.success(`已将「${label} = ${newVal || '(空)'}」同步至全部 ${editableBatchTargets.value.length} 本期刊`)
}

const addBatchExtraDiffField = (ef: { field: string; label: string }) => {
    let addedCount = 0
    editableBatchTargets.value.forEach((bt) => {
        if (bt.diffs.some((d) => d.field === ef.field)) return
        addedCount++
        const origForJournal = batchOriginalMap.value[bt.periodical_id] || bt.form_payload || {}
        const rawOld = origForJournal[ef.field]
        if (BOOL_FIELDS.has(ef.field)) {
            const boolOld = Boolean(rawOld)
            const nextBool = !boolOld
            bt.diffs.push({
                field: ef.field,
                label: ef.label,
                old_value: boolOld ? '是' : '否',
                new_value: nextBool ? '是' : '否',
            })
            if (bt.form_payload) bt.form_payload[ef.field] = nextBool
        } else {
            const oldStr = rawOld !== undefined && rawOld !== null && String(rawOld) !== '' ? String(rawOld) : '(空)'
            const defaultNew = rawOld !== undefined && rawOld !== null && String(rawOld) !== '' ? String(rawOld) : ''
            bt.diffs.push({
                field: ef.field,
                label: ef.label,
                old_value: oldStr,
                new_value: defaultNew,
            })
        }
    })
    if (addedCount > 0) {
        ElMessage.success(`已为 ${addedCount} 本期刊补充「${ef.label}」修改项`)
    }
}

const displayItems = computed<PeriodicalBriefCard[]>(() => {
    if (pagedItems.value) return pagedItems.value
    return props.card?.items || []
})

const totalCount = computed(() => props.card?.total_count || displayItems.value.length)
const totalPages = computed(() => Math.max(1, Math.ceil(totalCount.value / pageSize)))

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

const handlePageChange = async (targetPage: number) => {
    if (targetPage < 1 || targetPage > totalPages.value) return
    if (targetPage === 1 && props.card?.items && props.card.items.length > 0) {
        currentPage.value = 1
        pagedItems.value = null
        return
    }
    const fp = props.card?.filter_payload || {}
    pageLoading.value = true
    try {
        let publishTimeStart = ''
        let publishTimeEnd = ''
        if (fp.publish_month && String(fp.publish_month).length === 7) {
            publishTimeStart = `${fp.publish_month}-01`
            publishTimeEnd = `${fp.publish_month}-31`
        }
        const res = await reqPeriodicalList({
            page_number: targetPage,
            page_size: pageSize,
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
        if (res.code === 200 && res.data?.items) {
            pagedItems.value = res.data.items.map(mapRawPeriodicalToBrief)
            currentPage.value = targetPage
        }
    } finally {
        pageLoading.value = false
    }
}

const formatPageText = (val?: string) => {
    const clean = String(val || '').trim()
    if (!clean) return ''
    return clean.endsWith('页') ? clean : `${clean}页`
}

const handleTriggerModify = (item: PeriodicalBriefCard) => {
    emit('start-modify', item)
}

const handleCopyDescription = async (item: PeriodicalBriefCard) => {
    let desc = String(item.periodical_description || '').trim()
    if (!desc && item.periodical_id) {
        try {
            const res = await reqPeriodicalDetail(item.periodical_id)
            if (res.code === 200 && res.data?.periodical_description) {
                desc = String(res.data.periodical_description).trim()
                item.periodical_description = desc
            }
        } catch {
            // ignore
        }
    }
    if (!desc) {
        ElMessage.warning('暂无期刊描述')
        return
    }
    const ok = await copyTextToClipboard(desc)
    if (ok) {
        ElMessage.success(`已复制 ${item.name} 的期刊描述`)
    } else {
        ElMessage.error('复制失败')
    }
}

const handleDecision = (decision: 'confirm' | 'cancel') => {
    submitting.value = true
    emit('confirm-action', {
        messageId: props.messageId,
        decision,
        formPayload: editablePayload.value,
        diffs: editableDiffs.value,
        batchUpdates: isBatchUpdate.value ? editableBatchTargets.value : undefined,
    })
    setTimeout(() => {
        submitting.value = false
    }, 600)
}
</script>

<style scoped lang="scss">
.ai-action-card {
    margin-top: 12px;
    border: 1px solid #e2e8f0;
    border-radius: 12px;
    background: #ffffff;
    padding: 14px 16px;
    box-shadow: 0 2px 6px rgba(15, 23, 42, 0.04);

    &.card-delete_confirm {
        border-color: #fecaca;
        background: #fffafb;
    }

    &.card-create_confirm,
    &.card-update_confirm {
        border-color: #bfdbfe;
        background: #f8fbff;
    }

    .card-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        flex-wrap: wrap;
        gap: 10px;
        margin-bottom: 8px;

        .header-left {
            display: flex;
            align-items: center;
            gap: 8px;

            .type-dot {
                width: 8px;
                height: 8px;
                border-radius: 50%;
                background: #2563eb;

                &.delete_confirm {
                    background: #dc2626;
                }
                &.create_confirm {
                    background: #10b981;
                }
                &.update_confirm {
                    background: #f59e0b;
                }
            }

            .card-title {
                font-size: 13.5px;
                font-weight: 700;
                color: #0f172a;
            }
        }

        .header-right {
            display: flex;
            align-items: center;
            gap: 8px;
            flex-wrap: wrap;

            .view-switcher {
                display: inline-flex;
                background: #f1f5f9;
                padding: 2px;
                border-radius: 7px;
                border: 1px solid #e2e8f0;

                .switch-btn {
                    border: none;
                    background: transparent;
                    font-size: 11.5px;
                    font-weight: 600;
                    color: #64748b;
                    padding: 3px 9px;
                    border-radius: 5px;
                    cursor: pointer;
                    transition: all 0.15s ease;

                    &.active {
                        background: #ffffff;
                        color: #2563eb;
                        box-shadow: 0 1px 2px rgba(15, 23, 42, 0.08);
                    }
                }
            }
        }

        .status-pill {
            font-size: 11px;
            font-weight: 600;
            padding: 2px 8px;
            border-radius: 999px;

            &.pending {
                background: #fef3c7;
                color: #b45309;
            }
            &.confirmed {
                background: #dcfce7;
                color: #15803d;
            }
            &.cancelled {
                background: #f1f5f9;
                color: #64748b;
            }
        }
    }

    .card-desc {
        margin: 0 0 10px;
        font-size: 12.5px;
        color: #64748b;
        line-height: 1.5;
    }

    .update-summary-bar {
        margin-bottom: 10px;
        display: flex;
        align-items: center;
        justify-content: space-between;
        flex-wrap: wrap;
        gap: 8px;
        font-size: 12.5px;
        color: #475569;

        .changed-num {
            color: #059669;
        }

        .clean-unchanged-btn {
            border: 1px solid #cbd5e1;
            background: #ffffff;
            color: #475569;
            font-size: 11.5px;
            font-weight: 600;
            padding: 2px 9px;
            border-radius: 6px;
            cursor: pointer;
            transition: all 0.15s ease;

            &:hover {
                border-color: #94a3b8;
                background: #f1f5f9;
                color: #0f172a;
            }
        }
    }
}

.journal-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(225px, 1fr));
    gap: 10px;

    .journal-mini-card {
        padding: 11px 12px;
        border: 1px solid #e2e8f0;
        border-radius: 9px;
        background: #f8fafc;
        cursor: pointer;
        transition: all 0.16s ease;
        display: flex;
        flex-direction: column;
        justify-content: space-between;

        &:hover {
            border-color: #93c5fd;
            background: #f0f7ff;
            transform: translateY(-1px);
            box-shadow: 0 4px 12px -4px rgba(37, 99, 235, 0.12);
        }

        .mini-top {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 6px;
            margin-bottom: 6px;

            .j-name {
                font-size: 13px;
                font-weight: 700;
                color: #0f172a;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
            }

            .j-level {
                flex-shrink: 0;
                font-size: 11px;
                padding: 1px 6px;
                border-radius: 999px;
                background: #dbeafe;
                color: #1d4ed8;
                font-weight: 600;
            }
        }

        .mini-badges {
            display: flex;
            flex-wrap: wrap;
            gap: 4px;
            margin-bottom: 8px;

            .b-tag {
                font-size: 10.5px;
                padding: 1px 5px;
                border-radius: 4px;

                &.emerald {
                    background: #d1fae5;
                    color: #047857;
                }
                &.amber-hidden {
                    background: #fef3c7;
                    color: #b45309;
                    border: 1px solid #fcd34d;
                    font-weight: 600;
                }
                &.sky {
                    background: #e0f2fe;
                    color: #0369a1;
                }
                &.slate {
                    background: #e2e8f0;
                    color: #475569;
                }
                &.dir {
                    background: #fef3c7;
                    color: #92400e;
                    max-width: 110px;
                    overflow: hidden;
                    text-overflow: ellipsis;
                    white-space: nowrap;
                }
            }
        }

        .mini-bottom {
            display: flex;
            align-items: center;
            justify-content: space-between;
            font-size: 12px;

            .j-meta {
                color: #64748b;
            }

            .j-price {
                font-weight: 700;
                color: #0f172a;
                font-variant-numeric: tabular-nums;
            }
        }

        .mini-actions {
            margin-top: 7px;
            padding-top: 6px;
            border-top: 1px dashed #e2e8f0;
            display: flex;
            align-items: center;
            gap: 5px;
            font-size: 11px;

            .act-link {
                color: #2563eb;
                font-weight: 600;
                cursor: pointer;

                &:hover {
                    text-decoration: underline;
                }
            }

            .act-sep {
                color: #cbd5e1;
            }
        }
    }
}

.journal-table-wrap {
    width: 100%;
    overflow-x: auto;
    border: 1px solid #e2e8f0;
    border-radius: 8px;

    .journal-compact-table {
        width: 100%;
        border-collapse: collapse;
        font-size: 12.5px;

        th,
        td {
            padding: 8px 10px;
            border-bottom: 1px solid #f1f5f9;
            text-align: left;
        }

        th {
            background: #f8fafc;
            color: #475569;
            font-weight: 600;
            white-space: nowrap;
        }

        tbody tr:last-child td {
            border-bottom: none;
        }

        tbody tr:hover td {
            background: #f8fbff;
        }

        .col-id {
            color: #94a3b8;
            font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
        }

        .col-name {
            cursor: pointer;
            font-weight: 600;
            color: #0f172a;

            .name-link:hover {
                color: #2563eb;
            }

            .mini-pill {
                margin-left: 4px;
                font-size: 10px;
                padding: 1px 5px;
                border-radius: 4px;
                background: #d1fae5;
                color: #047857;

                &.amber-hidden {
                    background: #fef3c7;
                    color: #b45309;
                    border: 1px solid #fcd34d;
                }
            }
        }

        .level-tag {
            font-size: 11px;
            padding: 1px 6px;
            border-radius: 999px;
            background: #dbeafe;
            color: #1d4ed8;
            font-weight: 600;
            white-space: nowrap;
        }

        .col-tags {
            display: flex;
            flex-wrap: wrap;
            gap: 3px;

            .col-chip {
                font-size: 11px;
                padding: 1px 5px;
                border-radius: 4px;
                background: #f1f5f9;
                color: #475569;
            }
        }

        .col-dir {
            max-width: 150px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            color: #64748b;
        }

        .col-pub {
            color: #475569;
            white-space: nowrap;
        }

        .col-price {
            text-align: right;
            font-weight: 700;
            color: #0f172a;
            font-variant-numeric: tabular-nums;
        }

        .col-ops {
            text-align: right;
            white-space: nowrap;

            .t-act {
                color: #2563eb;
                font-weight: 600;
                cursor: pointer;
                margin-left: 8px;

                &:hover {
                    text-decoration: underline;
                }
            }
        }
    }
}

.card-pagination-bar {
    margin-top: 10px;
    padding-top: 9px;
    border-top: 1px solid #f1f5f9;
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 12px;
    color: #64748b;

    .page-btns {
        display: flex;
        align-items: center;
        gap: 6px;
    }
}

.inline-field-input {
    width: 100%;
    height: 34px;
    padding: 0 11px;
    border: 1px solid #cbd5e1;
    border-radius: 7px;
    background: #f8fafc;
    font-size: 13px;
    font-weight: 500;
    color: #0f172a;
    outline: none;
    transition: all 0.16s ease;

    &::placeholder {
        color: #94a3b8;
        font-weight: 400;
    }

    &:hover {
        border-color: #94a3b8;
        background: #ffffff;
    }

    &:focus {
        border-color: #2563eb;
        background: #ffffff;
        box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.12);
    }
}

.create-confirm-body {
    .matched-cats-panel {
        margin-bottom: 10px;
        padding: 8px 12px;
        border-radius: 9px;
        background: #f8fafc;
        border: 1px solid #e2e8f0;
        display: flex;
        align-items: center;
        flex-wrap: wrap;
        gap: 8px;

        .mc-title {
            font-size: 12px;
            font-weight: 600;
            color: #475569;
            flex-shrink: 0;
        }

        .mc-chips {
            display: flex;
            align-items: center;
            flex-wrap: wrap;
            gap: 6px;
        }

        .cat-chip {
            display: inline-flex;
            align-items: center;
            gap: 4px;
            font-size: 11.5px;
            padding: 2px 9px;
            border-radius: 999px;
            background: #eff6ff;
            color: #1d4ed8;
            border: 1px solid #bfdbfe;
            font-weight: 600;

            em {
                font-style: normal;
                font-weight: 500;
                opacity: 0.72;
                font-size: 11px;
            }

            &.level {
                background: #eef2ff;
                color: #4338ca;
                border-color: #c7d2fe;
            }

            &.cycle {
                background: #fffbeb;
                color: #b45309;
                border-color: #fde68a;
            }

            &.col {
                background: #f0f9ff;
                color: #0369a1;
                border-color: #bae6fd;
            }

            &.dir {
                background: #ecfdf5;
                color: #047857;
                border-color: #a7f3d0;
            }
        }
    }

    .field-grid {
        display: grid;
        grid-template-columns: repeat(2, minmax(0, 1fr));
        gap: 12px 16px;
        background: #ffffff;
        border: 1px solid #e2e8f0;
        border-radius: 10px;
        padding: 14px 16px;
        box-shadow: 0 1px 2px rgba(15, 23, 42, 0.02);

        .field-item {
            display: flex;
            flex-direction: column;
            gap: 5px;

            .f-label {
                font-size: 12px;
                font-weight: 600;
                color: #475569;

                .req-star {
                    color: #ef4444;
                    font-style: normal;
                    margin-left: 2px;
                }
            }

            .create-bool-toggles {
                display: flex;
                flex-wrap: wrap;
                gap: 6px;
                min-height: 34px;
                align-items: center;

                .c-bool-chip {
                    height: 32px;
                    display: inline-flex;
                    align-items: center;
                    border: 1px solid #cbd5e1;
                    background: #f8fafc;
                    color: #64748b;
                    font-size: 12px;
                    font-weight: 500;
                    padding: 0 11px;
                    border-radius: 7px;
                    cursor: pointer;
                    transition: all 0.16s ease;

                    &:hover {
                        border-color: #93c5fd;
                        color: #1e293b;
                        background: #ffffff;
                    }

                    &.active {
                        border-color: #34d399;
                        background: #ecfdf5;
                        color: #047857;
                        font-weight: 600;
                        box-shadow: 0 1px 2px rgba(16, 185, 129, 0.08);
                    }
                }
            }

            .f-val {
                font-size: 13px;
                color: #1e293b;
                font-weight: 500;
                word-break: break-all;

                &.highlight {
                    color: #2563eb;
                    font-weight: 700;
                }

                &.price {
                    color: #059669;
                    font-weight: 700;
                }
            }
        }
    }

    .create-extra-list {
        margin-top: 10px;
        display: flex;
        flex-direction: column;
        gap: 8px;

        .create-extra-row {
            background: #ffffff;
            border: 1px solid #e2e8f0;
            border-radius: 9px;
            padding: 10px 14px;
            transition: border-color 0.15s ease;

            &:focus-within {
                border-color: #93c5fd;
            }

            .ce-header {
                display: flex;
                align-items: center;
                justify-content: space-between;
                margin-bottom: 6px;

                .ce-label {
                    font-size: 12px;
                    font-weight: 600;
                    color: #334155;
                }

                .ce-remove-btn {
                    border: 1px solid #e2e8f0;
                    background: #f8fafc;
                    color: #64748b;
                    font-size: 11px;
                    font-weight: 600;
                    padding: 2px 8px;
                    border-radius: 999px;
                    cursor: pointer;
                    transition: all 0.15s ease;

                    &:hover {
                        border-color: #fecaca;
                        background: #fef2f2;
                        color: #dc2626;
                    }
                }
            }

            .ce-textarea {
                width: 100%;
                border: 1px solid #cbd5e1;
                border-radius: 7px;
                background: #f8fafc;
                padding: 7px 10px;
                font-size: 12.5px;
                color: #0f172a;
                font-family: inherit;
                resize: vertical;
                outline: none;
                transition: all 0.16s ease;

                &:focus {
                    background: #ffffff;
                    border-color: #2563eb;
                    box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.12);
                }
            }

            .ce-readonly-val {
                font-size: 12.5px;
                color: #334155;
                white-space: pre-wrap;
            }
        }
    }
}

.update-confirm-body {
    .batch-targets-list {
        display: flex;
        flex-direction: column;
        gap: 12px;

        .batch-target-section {
            background: #ffffff;
            border: 1px solid #cbd5e1;
            border-radius: 10px;
            padding: 11px 13px;
            box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03);
            transition: border-color 0.16s ease;

            &:hover {
                border-color: #93c5fd;
            }

            .bt-header {
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: 8px;
                padding-bottom: 9px;
                margin-bottom: 9px;
                border-bottom: 1px dashed #e2e8f0;

                .bt-left {
                    display: flex;
                    align-items: center;
                    flex-wrap: wrap;
                    gap: 8px;

                    .bt-index {
                        font-size: 11px;
                        font-weight: 800;
                        font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
                        color: #1d4ed8;
                        background: #dbeafe;
                        padding: 2px 7px;
                        border-radius: 6px;
                    }

                    .bt-name {
                        font-size: 13.5px;
                        font-weight: 700;
                        color: #0f172a;
                    }

                    .bt-level {
                        font-size: 11px;
                        font-weight: 600;
                        color: #1e40af;
                        background: #eff6ff;
                        border: 1px solid #bfdbfe;
                        padding: 1px 7px;
                        border-radius: 999px;
                    }

                    .bt-count-pill {
                        font-size: 11px;
                        font-weight: 600;
                        color: #047857;
                        background: #ecfdf5;
                        border: 1px solid #a7f3d0;
                        padding: 1px 7px;
                        border-radius: 999px;
                    }
                }

                .bt-right {
                    .bt-remove-journal-btn {
                        border: 1px solid #e2e8f0;
                        background: #f8fafc;
                        color: #64748b;
                        font-size: 11.5px;
                        font-weight: 600;
                        padding: 3px 9px;
                        border-radius: 999px;
                        cursor: pointer;
                        transition: all 0.15s ease;

                        &:hover {
                            color: #dc2626;
                            border-color: #fecaca;
                            background: #fef2f2;
                        }
                    }
                }
            }
        }
    }

    .diff-list {
        display: flex;
        flex-direction: column;
        gap: 8px;

        .diff-item-card {
            background: #ffffff;
            border: 1px solid #e2e8f0;
            border-radius: 9px;
            padding: 9px 12px;
            transition: all 0.15s ease;

            &.is-changed {
                border-color: #a7f3d0;
                background: #fcfffd;
            }

            &:not(.is-longtext) {
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: 12px;
            }

            &.is-longtext {
                display: flex;
                flex-direction: column;
                gap: 8px;
            }

            .diff-label-badge {
                font-size: 12.5px;
                font-weight: 700;
                color: #1e293b;
                white-space: nowrap;
                flex-shrink: 0;
            }

            .change-state-dot {
                width: 7px;
                height: 7px;
                border-radius: 50%;
                flex-shrink: 0;

                &.changed {
                    background: #10b981;
                    box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.2);
                }
                &.same {
                    background: #cbd5e1;
                }
            }

            .change-state-tag {
                font-size: 10.5px;
                font-weight: 600;
                padding: 1px 6px;
                border-radius: 999px;
                white-space: nowrap;

                &.changed {
                    background: #d1fae5;
                    color: #047857;
                }
                &.same {
                    background: #f1f5f9;
                    color: #64748b;
                }
            }

            .compact-left {
                width: 118px;
                flex-shrink: 0;
                display: flex;
                align-items: center;
                gap: 6px;
            }

            .compact-middle {
                flex: 1;
                min-width: 0;
                display: flex;
                align-items: center;
                gap: 10px;

                .old-val-pill {
                    max-width: 190px;
                    overflow: hidden;
                    text-overflow: ellipsis;
                    white-space: nowrap;
                    font-size: 12px;
                    color: #64748b;
                    background: #f8fafc;
                    border: 1px solid #e2e8f0;
                    padding: 2px 8px;
                    border-radius: 5px;
                    flex-shrink: 0;
                }

                .diff-arrow {
                    color: #94a3b8;
                    flex-shrink: 0;
                }

                .new-val {
                    color: #047857;
                    font-weight: 700;
                    background: #ecfdf5;
                    padding: 2px 8px;
                    border-radius: 5px;
                    font-size: 12.5px;
                }

                .inline-diff-input {
                    flex: 1;
                    min-width: 160px;
                    max-width: 340px;
                    height: 30px;
                    padding: 0 10px;
                    border: 1px solid #6ee7b7;
                    border-radius: 6px;
                    background: #f0fdf4;
                    color: #065f46;
                    font-weight: 600;
                    font-size: 13px;
                    outline: none;
                    transition: all 0.15s ease;

                    &:focus {
                        background: #ffffff;
                        border-color: #10b981;
                        box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.16);
                    }
                }

                .bool-switch-group {
                    display: inline-flex;
                    border: 1px solid #cbd5e1;
                    border-radius: 6px;
                    overflow: hidden;

                    .bool-btn {
                        border: none;
                        background: #ffffff;
                        padding: 3px 12px;
                        font-size: 12px;
                        font-weight: 600;
                        color: #64748b;
                        cursor: pointer;

                        &.active {
                            background: #10b981;
                            color: #ffffff;
                        }
                    }
                }
            }

            .longtext-top {
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: 10px;
                flex-wrap: wrap;

                .lt-left {
                    display: flex;
                    align-items: center;
                    gap: 8px;
                    min-width: 0;
                    flex-wrap: wrap;

                    .lt-old-summary {
                        font-size: 12px;
                        color: #64748b;
                        max-width: 380px;
                        overflow: hidden;
                        text-overflow: ellipsis;
                        white-space: nowrap;
                    }

                    .toggle-old-btn {
                        border: none;
                        background: #eff6ff;
                        color: #2563eb;
                        font-size: 11px;
                        font-weight: 600;
                        padding: 2px 7px;
                        border-radius: 4px;
                        cursor: pointer;

                        &:hover {
                            background: #dbeafe;
                        }
                    }
                }
            }

            .lt-old-full-box {
                padding: 8px 10px;
                border-radius: 6px;
                background: #f8fafc;
                border: 1px dashed #cbd5e1;
                font-size: 12px;
                color: #475569;
                line-height: 1.55;

                .old-box-label {
                    font-size: 11px;
                    font-weight: 700;
                    color: #64748b;
                    margin-bottom: 3px;
                }

                .old-box-content {
                    white-space: pre-wrap;
                    word-break: break-word;
                }
            }

            .longtext-editor-wrap {
                width: 100%;

                .inline-diff-textarea {
                    width: 100%;
                    min-height: 70px;
                    padding: 8px 10px;
                    border: 1px solid #6ee7b7;
                    border-radius: 7px;
                    background: #f0fdf4;
                    color: #064e3b;
                    font-size: 13px;
                    line-height: 1.55;
                    font-family: inherit;
                    resize: vertical;
                    outline: none;
                    transition: all 0.15s ease;

                    &:focus {
                        background: #ffffff;
                        border-color: #10b981;
                        box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.16);
                    }
                }

                .new-val-block {
                    padding: 8px 10px;
                    border-radius: 6px;
                    background: #ecfdf5;
                    color: #065f46;
                    font-size: 12.5px;
                    line-height: 1.55;
                    white-space: pre-wrap;
                    word-break: break-word;
                }
            }

            .diff-row-actions {
                display: flex;
                align-items: center;
                gap: 6px;
                flex-shrink: 0;

                .row-act-btn {
                    border: 1px solid #e2e8f0;
                    background: #f8fafc;
                    font-size: 11.5px;
                    font-weight: 600;
                    padding: 3px 8px;
                    border-radius: 5px;
                    cursor: pointer;
                    transition: all 0.14s ease;
                    white-space: nowrap;

                    &.sync-all {
                        color: #047857;
                        border-color: #a7f3d0;
                        background: #ecfdf5;

                        &:hover {
                            background: #d1fae5;
                            border-color: #34d399;
                        }
                    }

                    &.reset {
                        color: #2563eb;
                        border-color: #bfdbfe;
                        background: #eff6ff;

                        &:hover {
                            background: #dbeafe;
                        }
                    }

                    &.remove {
                        color: #64748b;

                        &:hover {
                            color: #dc2626;
                            border-color: #fecaca;
                            background: #fef2f2;
                        }
                    }
                }
            }
        }
    }

    .empty-diffs-box {
        padding: 16px;
        text-align: center;
        background: #ffffff;
        border: 1px dashed #cbd5e1;
        border-radius: 8px;
        font-size: 12.5px;
        color: #64748b;
    }
}

/* 通用渐进式字段补充工具栏（同时应用于新增入库确认单 create_confirm 与修改确认单 update_confirm） */
.extra-fields-bar {
    margin-top: 12px;
    padding: 9px 12px;
    border-radius: 9px;
    background: #f8fafc;
    border: 1px dashed #cbd5e1;
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 10px;

    .extra-hint {
        display: inline-flex;
        align-items: center;
        gap: 5px;
        font-size: 12px;
        font-weight: 600;
        color: #475569;
        flex-shrink: 0;

        .hint-dot {
            width: 16px;
            height: 16px;
            border-radius: 50%;
            background: #dbeafe;
            color: #2563eb;
            font-size: 12px;
            font-weight: 700;
            display: inline-flex;
            align-items: center;
            justify-content: center;
            line-height: 1;
        }
    }

    .extra-chips-wrap {
        display: flex;
        align-items: center;
        flex-wrap: wrap;
        gap: 7px;
    }

    .add-field-chip {
        height: 27px;
        display: inline-flex;
        align-items: center;
        gap: 3px;
        border: 1px solid #cbd5e1;
        background: #ffffff;
        color: #334155;
        font-size: 12px;
        font-weight: 500;
        padding: 0 11px;
        border-radius: 999px;
        cursor: pointer;
        box-shadow: 0 1px 2px rgba(15, 23, 42, 0.03);
        transition: all 0.16s ease;

        .chip-plus {
            color: #2563eb;
            font-weight: 700;
            font-size: 13px;
            line-height: 1;
        }

        &:hover {
            background: #eff6ff;
            border-color: #60a5fa;
            color: #1d4ed8;
            transform: translateY(-1px);
            box-shadow: 0 3px 6px -2px rgba(37, 99, 235, 0.14);
        }

        &.more-dropdown-chip {
            background: #eff6ff;
            border-color: #93c5fd;
            color: #1d4ed8;
            font-weight: 600;

            &:hover {
                background: #dbeafe;
                border-color: #3b82f6;
            }
        }
    }
}

.delete-confirm-body {
    .delete-target-box {
        padding: 10px 12px;
        border-radius: 8px;
        background: #fef2f2;
        border: 1px solid #fecaca;

        .target-main {
            display: flex;
            align-items: center;
            gap: 8px;
            color: #991b1b;
            font-size: 13.5px;
            margin-bottom: 4px;

            .j-level {
                font-size: 11px;
                padding: 1px 6px;
                border-radius: 999px;
                background: #fee2e2;
                color: #b91c1c;
            }
        }

        .target-sub {
            font-size: 12px;
            color: #7f1d1d;
        }
    }
}

.card-actions {
    margin-top: 12px;
    padding-top: 10px;
    border-top: 1px solid #e2e8f0;
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 8px;
}

.confirmed-effect-tag {
    font-size: 11px;
    font-weight: 700;
    padding: 1px 7px;
    border-radius: 999px;
    background: #dcfce7;
    color: #15803d;
}

.confirmed-integrated-footer {
    margin-top: 10px;
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.embedded-snapshot-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.card-result-banner {
    padding: 6px 10px;
    border-radius: 6px;
    background: #f1f5f9;
    font-size: 12px;
    color: #334155;
    font-weight: 500;

    &.confirmed {
        background: #ecfdf5;
        color: #065f46;
        border: 1px solid #a7f3d0;
    }
}

.embedded-snapshot-card {
    padding: 10px 14px;
    border-radius: 9px;
    background: #ffffff;
    border: 1px solid #cbd5e1;
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: 10px;
    cursor: pointer;
    transition: all 0.15s ease;

    &:hover {
        border-color: #60a5fa;
        background: #f8fbff;
    }

    .snap-left {
        display: flex;
        align-items: center;
        flex-wrap: wrap;
        gap: 8px;
        font-size: 12.5px;

        .snap-tag {
            font-size: 10.5px;
            font-weight: 700;
            padding: 2px 6px;
            border-radius: 4px;
            background: #2563eb;
            color: #ffffff;
        }

        .snap-name {
            font-size: 13.5px;
            color: #0f172a;
        }

        .snap-level {
            font-size: 11px;
            padding: 1px 6px;
            border-radius: 999px;
            background: #dbeafe;
            color: #1d4ed8;
            font-weight: 600;
        }

        .snap-col {
            font-size: 11px;
            padding: 1px 5px;
            border-radius: 4px;
            background: #f1f5f9;
            color: #475569;
        }

        .snap-meta {
            color: #475569;

            .snap-price {
                color: #059669;
            }
        }
    }

    .snap-actions {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 12px;

        .s-act {
            color: #2563eb;
            font-weight: 600;
            cursor: pointer;

            &:hover {
                text-decoration: underline;
            }
        }

        .s-sep {
            color: #cbd5e1;
        }
    }
}

.empty-result {
    padding: 12px;
    text-align: center;
    color: #94a3b8;
    font-size: 12.5px;
}
</style>
