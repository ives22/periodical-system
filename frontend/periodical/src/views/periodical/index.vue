<!--
 * new page
 * @author: jalen
 * @since: 2024-12-12
 * index.vue
-->
<template>
    <div class="periodical-page">
        <div class="container">
            <div class="column">
                <div class="column-left">
                    <div class="filter-panel-header">
                        <span class="panel-title">多维条件筛选</span>
                        <el-button link type="primary" size="small" @click="handleReset">清空筛选</el-button>
                    </div>

                    <el-form class="date-filter-box">
                        <el-row :gutter="10">
                            <el-col :span="13">
                                <span class="form-label">出刊月份</span>
                                <el-form-item>
                                    <el-date-picker v-model="searchForm.publish_month" type="month" placeholder="选择月份"
                                        format="YYYY-MM" value-format="YYYY-MM" :clearable="true"
                                        @change="handleMonthChange" class="month-picker" />
                                </el-form-item>
                            </el-col>
                            <el-col :span="11">
                                <span class="form-label">出刊旬期</span>
                                <el-form-item>
                                    <el-select v-model="searchForm.publish_period" placeholder="全部"
                                        :disabled="!searchForm.publish_month" @change="handlePeriodChange" clearable
                                        class="period-select">
                                        <el-option label="月初" value="1" />
                                        <el-option label="上旬" value="10" />
                                        <el-option label="中旬" value="15" />
                                        <el-option label="下旬" value="20" />
                                        <el-option label="月底" value="end" />
                                    </el-select>
                                </el-form-item>
                            </el-col>
                        </el-row>
                    </el-form>

                    <!-- 收稿方向 -->
                    <div class="direction-wrapper">
                        <div class="direction-header">
                            <div class="direction-title">
                                收稿方向
                                <span v-if="checkedCities.length" class="count-badge">{{ checkedCities.length }}</span>
                            </div>
                            <el-checkbox v-model="checkAll" :indeterminate="isIndeterminate"
                                @change="handleCheckAllChange" class="check-all">
                                全部
                            </el-checkbox>
                        </div>
                        <el-checkbox-group v-model="checkedCities" @change="handleCheckedCitiesChange"
                            class="direction-group">
                            <el-checkbox v-for="item in directionList" :key="item.categorize_id"
                                :value="{ base_id: item.categorize_id, type: item.type }" class="direction-item">
                                <el-tooltip :content="item.name" placement="top" :show-after="500">
                                    <span class="checkbox-label">{{ item.name }}</span>
                                </el-tooltip>
                            </el-checkbox>
                        </el-checkbox-group>
                    </div>

                    <!-- 网站收录 -->
                    <div class="direction-wrapper">
                        <div class="direction-header">
                            <div class="direction-title">
                                网站收录
                                <span v-if="checkedCollection.length" class="count-badge">{{ checkedCollection.length }}</span>
                            </div>
                            <el-checkbox v-model="collectionAll" :indeterminate="collectionIsIndeterminate"
                                @change="handleCollectionCheckAllChange" class="check-all">
                                全部
                            </el-checkbox>
                        </div>
                        <el-checkbox-group v-model="checkedCollection" @change="handleCheckedCollectionChange"
                            class="direction-group">
                            <el-checkbox v-for="item in collectionList" :key="item.categorize_id"
                                :value="{ base_id: item.categorize_id, type: item.type }" class="direction-item">
                                <el-tooltip :content="item.name" placement="top" :show-after="500">
                                    <span class="checkbox-label">{{ item.name }}</span>
                                </el-tooltip>
                            </el-checkbox>
                        </el-checkbox-group>
                    </div>

                    <!-- 期刊级别 -->
                    <div class="direction-wrapper">
                        <div class="direction-header">
                            <div class="direction-title">
                                期刊级别
                                <span v-if="checkedLevels.length" class="count-badge">{{ checkedLevels.length }}</span>
                            </div>
                            <el-checkbox v-model="levelCheckAll" :indeterminate="levelIsIndeterminate"
                                @change="handleLevelCheckAllChange" class="check-all">
                                全部
                            </el-checkbox>
                        </div>
                        <el-checkbox-group v-model="checkedLevels" @change="handleCheckedLevelsChange"
                            class="direction-group">
                            <el-checkbox v-for="item in levelList" :key="item.categorize_id"
                                :value="{ base_id: item.categorize_id, type: item.type }" class="direction-item">
                                <el-tooltip :content="item.name" placement="top" :show-after="500">
                                    <span class="checkbox-label">{{ item.name }}</span>
                                </el-tooltip>
                            </el-checkbox>
                        </el-checkbox-group>
                    </div>

                    <!-- 出版周期 -->
                    <div class="direction-wrapper">
                        <div class="direction-header">
                            <div class="direction-title">
                                出版周期
                                <span v-if="checkedCycles.length" class="count-badge">{{ checkedCycles.length }}</span>
                            </div>
                            <el-checkbox v-model="cycleCheckAll" :indeterminate="cycleIsIndeterminate"
                                @change="handleCycleCheckAllChange" class="check-all">
                                全部
                            </el-checkbox>
                        </div>
                        <el-checkbox-group v-model="checkedCycles" @change="handleCheckedCyclesChange"
                            class="direction-group">
                            <el-checkbox v-for="item in cycleList" :key="item.categorize_id"
                                :value="{ base_id: item.categorize_id, type: item.type }" class="direction-item">
                                <el-tooltip :content="item.name" placement="top" :show-after="500">
                                    <span class="checkbox-label">{{ item.name }}</span>
                                </el-tooltip>
                            </el-checkbox>
                        </el-checkbox-group>
                    </div>

                    <!-- 学术期刊 -->
                    <div class="direction-wrapper">
                        <div class="direction-header">
                            <div class="direction-title">学术期刊批次</div>
                        </div>
                        <el-checkbox-group v-model="checkedAcademicJournals" @change="handleAcademicJournalsChange"
                            class="direction-group">
                            <el-checkbox :value="1" class="direction-item">第一批学术</el-checkbox>
                            <el-checkbox :value="2" class="direction-item">第二批学术</el-checkbox>
                        </el-checkbox-group>
                    </div>

                    <!-- 发票选项 -->
                    <div class="direction-wrapper">
                        <div class="direction-header">
                            <div class="direction-title">发票服务</div>
                        </div>
                        <el-checkbox-group v-model="canIssueInvoice" @change="handleInvoiceChange">
                            <el-checkbox :value="1" class="direction-item" style="width: 100%;">可开杂志社发票</el-checkbox>
                        </el-checkbox-group>
                    </div>

                    <!-- 可否发布作品 & 全包筛选 紧凑双列 -->
                    <div class="direction-wrapper compact-dual-row">
                        <div class="dual-col">
                            <div class="direction-header">
                                <div class="direction-title">可发作品</div>
                            </div>
                            <el-radio-group v-model="canPublishWorks" @change="handlePublishWorksChange">
                                <el-radio :value="true">可</el-radio>
                                <el-radio :value="false">否</el-radio>
                            </el-radio-group>
                        </div>
                        <div class="dual-col">
                            <div class="direction-header">
                                <div class="direction-title">是否全包</div>
                            </div>
                            <el-radio-group v-model="isWrap" @change="handleWrapChange">
                                <el-radio :value="true">是</el-radio>
                                <el-radio :value="false">否</el-radio>
                            </el-radio-group>
                        </div>
                    </div>

                    <!-- 复合影响因子 -->
                    <div class="direction-wrapper">
                        <div class="direction-header">
                            <div class="direction-title">复合影响因子</div>
                        </div>
                        <el-radio-group v-model="compositeImpactFactor" class="direction-group"
                            @change="handleCompositeImpactFactorChange">
                            <el-radio :value="0">0以上</el-radio>
                            <el-radio :value="0.3">0.3以上</el-radio>
                            <el-radio :value="0.5">0.5以上</el-radio>
                        </el-radio-group>
                    </div>

                    <!-- 隐藏状态筛选（仅具备编辑权限的管理员可见） -->
                    <div v-if="userStore.hasPermission('periodical:update')" class="direction-wrapper">
                        <div class="direction-header">
                            <div class="direction-title">隐藏状态</div>
                        </div>
                        <el-radio-group v-model="hiddenFilter" class="direction-group" @change="handleHiddenFilterChange">
                            <el-radio :value="0">全部</el-radio>
                            <el-radio :value="2">未隐藏</el-radio>
                            <el-radio :value="1">已隐藏</el-radio>
                        </el-radio-group>
                    </div>

                </div>

                <div class="column-right">
                    <!-- 高级搜索面板 -->
                    <el-form :model="searchForm" ref="searchFormRef" class="search-form" label-width="72px"
                        @submit.prevent="handleSearch">
                        <el-row :gutter="16" align="middle">
                            <el-col :xs="24" :sm="12" :md="6">
                                <el-form-item label="期刊名称" prop="name">
                                    <el-input v-model="searchForm.name" placeholder="搜索期刊名称..." clearable
                                        prefix-icon="Search" @keyup.enter="handleSearch" />
                                </el-form-item>
                            </el-col>
                            <el-col :xs="24" :sm="12" :md="6">
                                <el-form-item label="主办单位" prop="sponsor">
                                    <el-input v-model="searchForm.organization" placeholder="输入主办单位..." clearable
                                        @keyup.enter="handleSearch" />
                                </el-form-item>
                            </el-col>
                            <el-col :xs="24" :sm="12" :md="6">
                                <el-form-item label="主管单位" prop="supervisor">
                                    <el-input v-model="searchForm.competent_unit" placeholder="输入主管单位..." clearable
                                        @keyup.enter="handleSearch" />
                                </el-form-item>
                            </el-col>

                            <el-col :xs="24" :sm="12" :md="6">
                                <div class="search-actions">
                                    <el-button type="primary" icon="Search" @click="handleSearch">搜索</el-button>
                                    <el-button icon="RefreshLeft" @click="handleReset">重置</el-button>
                                    <el-button class="expand-toggle" link type="primary" @click="isExpand = !isExpand">
                                        {{ isExpand ? '收起' : '更多' }}
                                        <el-icon class="el-icon--right">
                                            <component :is="isExpand ? 'ArrowUp' : 'ArrowDown'" />
                                        </el-icon>
                                    </el-button>
                                </div>
                            </el-col>
                        </el-row>

                        <!-- 展开/收起的更多筛选项 -->
                        <div v-if="isExpand" class="expanded-search-row">
                            <el-row :gutter="16">
                                <el-col :xs="24" :sm="12" :md="6">
                                    <el-form-item label="国内刊号" prop="domestic_number">
                                        <el-input v-model="searchForm.domestic_number" placeholder="如：23-1558/C" clearable
                                            @keyup.enter="handleSearch" />
                                    </el-form-item>
                                </el-col>
                                <el-col :xs="24" :sm="12" :md="6">
                                    <el-form-item label="国际刊号" prop="international_number">
                                        <el-input v-model="searchForm.international_number" placeholder="如：1674-7879"
                                            clearable @keyup.enter="handleSearch" />
                                    </el-form-item>
                                </el-col>
                                <el-col :xs="24" :sm="12" :md="6">
                                    <el-form-item label="栏目设置" prop="column_setting">
                                        <el-input v-model="searchForm.column_setting" placeholder="搜索栏目关键字..." clearable
                                            @keyup.enter="handleSearch" />
                                    </el-form-item>
                                </el-col>
                            </el-row>
                        </div>
                    </el-form>

                    <!-- 表格工具栏 -->
                    <div class="table-toolbar">
                        <div class="toolbar-left">
                            <h3 class="table-title">期刊目录列表</h3>
                            <span class="total-pill">共 <strong>{{ total }}</strong> 本刊物</span>
                        </div>
                        <div class="toolbar-right" v-if="userStore.hasPermission('periodical:create')">
                            <el-button type="primary" icon="Plus" @click="handleAdd">
                                新增期刊
                            </el-button>
                        </div>
                    </div>

                    <!-- 期刊列表表格 -->
                    <div class="table-container">
                        <el-table :data="periodicalArr" stripe @sort-change="sortChange" style="width: 100%"
                            :row-class-name="getTableRowClassName"
                            v-loading="tableLoading">
                            <el-table-column prop="name" label="期刊名称" min-width="230" show-overflow-tooltip>
                                <template #default="{ row }">
                                    <div class="journal-title-cell">
                                        <el-tooltip :content="row.periodical_description || '点击查看期刊详情'" placement="bottom"
                                            effect="dark" popper-class="custom-tooltip">
                                            <span class="journal-link" :class="{ 'is-hidden-title': row.is_hidden }" @click="handlePeriodicalClick(row)">
                                                {{ row.name }}
                                            </span>
                                        </el-tooltip>
                                        <div class="journal-meta-badges">
                                            <span v-if="row.is_hidden" class="mini-badge amber-hidden">已隐藏</span>
                                            <span v-if="row.is_warp" class="mini-badge emerald">全包</span>
                                            <span v-if="row.invoice_receipt" class="mini-badge sky">可开发票</span>
                                            <span
                                                v-for="colName in getCollections(row.categorizes)"
                                                :key="colName"
                                                class="mini-badge slate"
                                            >
                                                {{ colName }}
                                            </span>
                                        </div>
                                    </div>
                                </template>
                            </el-table-column>

                            <el-table-column label="期刊级别" width="125" show-overflow-tooltip>
                                <template #default="{ row }">
                                    <span
                                        v-if="getPeriodicalLevel(row.categorizes)"
                                        class="level-pill"
                                        :class="getLevelBadgeClass(getPeriodicalLevel(row.categorizes))"
                                    >
                                        {{ getPeriodicalLevel(row.categorizes) }}
                                    </span>
                                    <span v-else class="empty-dash">-</span>
                                </template>
                            </el-table-column>

                            <el-table-column label="期刊页码" prop="periodical_page" width="115" sortable="custom"
                                show-overflow-tooltip>
                                <template #default="{ row }">
                                    <span v-if="row.periodical_page" class="page-chip">{{ row.periodical_page }} 页</span>
                                    <span v-else class="empty-dash">-</span>
                                </template>
                            </el-table-column>

                            <el-table-column prop="publish_time" sortable="custom" label="出刊时间" width="145" show-overflow-tooltip>
                                <template #default="{ row }">
                                    <span v-if="formatDate(row.publication_time)" class="date-chip">
                                        {{ formatDate(row.publication_time) }}
                                    </span>
                                    <span v-else class="empty-dash">-</span>
                                </template>
                            </el-table-column>

                            <el-table-column prop="price" label="期刊发表费" width="155" sortable="custom" show-overflow-tooltip>
                                <template #default="{ row }">
                                    <span v-if="row.price && String(row.price) !== '0'" class="price-text">
                                        <em class="currency">¥</em>{{ row.price }}
                                    </span>
                                    <span v-else class="price-muted">待询</span>
                                </template>
                            </el-table-column>

                            <el-table-column prop="updated_at" label="更新时间" width="165" sortable="custom" show-overflow-tooltip>
                                <template #default="{ row }">
                                    <span class="time-muted">{{ formatTimestamp(row.updated_at) || '-' }}</span>
                                </template>
                            </el-table-column>

                            <el-table-column v-if="userStore.hasPermission('periodical:update')" label="是否隐藏" width="105" align="center">
                                <template #default="{ row }">
                                    <el-switch
                                        :model-value="Boolean(row.is_hidden)"
                                        :loading="Boolean(hiddenLoadingMap[row.periodical_id])"
                                        inline-prompt
                                        active-text="隐藏"
                                        inactive-text="正常"
                                        style="--el-switch-on-color: #f59e0b; --el-switch-off-color: #cbd5e1;"
                                        @change="(val: any) => handleToggleHidden(row, Boolean(val))"
                                    />
                                </template>
                            </el-table-column>

                            <el-table-column label="操作" width="160" align="center" fixed="right">
                                <template #default="{ row }">
                                    <div class="operation-wrapper">
                                        <el-tooltip content="查看详情" placement="top" :show-after="300">
                                            <button type="button" class="action-icon-btn view" @click="handleView(row)">
                                                <el-icon><View /></el-icon>
                                            </button>
                                        </el-tooltip>
                                        <el-tooltip v-if="userStore.hasPermission('periodical:copy')" content="复制期刊信息" placement="top" :show-after="300">
                                            <button
                                                type="button"
                                                class="action-icon-btn copy"
                                                @click="handleCopy(row)"
                                            >
                                                <el-icon><DocumentCopy /></el-icon>
                                            </button>
                                        </el-tooltip>
                                        <el-tooltip v-if="userStore.hasPermission('periodical:update')" content="编辑期刊" placement="top" :show-after="300">
                                            <button
                                                type="button"
                                                class="action-icon-btn edit"
                                                @click="handleEdit(row)"
                                            >
                                                <el-icon><Edit /></el-icon>
                                            </button>
                                        </el-tooltip>
                                        <el-tooltip v-if="userStore.hasPermission('periodical:delete')" content="删除期刊" placement="top" :show-after="300">
                                            <button
                                                type="button"
                                                class="action-icon-btn delete"
                                                @click="handleDelete(row)"
                                            >
                                                <el-icon><Delete /></el-icon>
                                            </button>
                                        </el-tooltip>
                                    </div>
                                </template>
                            </el-table-column>
                        </el-table>
                    </div>

                    <!-- 分页 -->
                    <div class="pagination-container">
                        <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize"
                            :page-sizes="[20, 30, 50]" layout="total, sizes, prev, pager, next, jumper" :total="total"
                            @size-change="handleSizeChange" @current-change="handleCurrentChange" />
                    </div>

                    <!-- 新增期刊抽屉 -->
                    <el-drawer v-model="drawerVisible" :title="isEditing ? '编辑期刊' : '新增期刊'" direction="rtl" size="70%">
                        <el-form ref="addFormRef" :model="addForm" :rules="rules" label-width="100px"
                            class="drawer-form">
                            <el-row :gutter="20">
                                <!-- 第一行 -->
                                <el-col :span="12">
                                    <el-form-item label="期刊名称" prop="name">
                                        <el-input v-model="addForm.name" placeholder="请输入期刊名称" />
                                    </el-form-item>
                                </el-col>
                                <el-col :span="12">
                                    <el-form-item label="主办单位" prop="organization">
                                        <el-input v-model="addForm.organization" placeholder="请输入主办单位" />
                                    </el-form-item>
                                </el-col>

                                <!-- 第二行 -->
                                <el-col :span="12">
                                    <el-form-item label="主管单位" prop="competent_unit">
                                        <el-input v-model="addForm.competent_unit" placeholder="请输入主管单位" />
                                    </el-form-item>
                                </el-col>
                                <el-col :span="12">
                                    <el-form-item label="国内刊号" prop="domestic_number">
                                        <el-input v-model="addForm.domestic_number" placeholder="请输入国内刊号" />
                                    </el-form-item>
                                </el-col>

                                <!-- 第三行 -->
                                <el-col :span="12">
                                    <el-form-item label="国际刊号" prop="international_number">
                                        <el-input v-model="addForm.international_number" placeholder="请输入国际刊号" />
                                    </el-form-item>
                                </el-col>
                                <el-col :span="12">
                                    <el-form-item label="审核周期" prop="audit_time">
                                        <el-input v-model="addForm.audit_time" placeholder="请输入审核周期" />
                                    </el-form-item>
                                </el-col>

                                <!-- 第四行 -->
                                <el-col :span="12">
                                    <el-form-item label="期刊发表费" prop="price">
                                        <el-input v-model="addForm.price" placeholder="请输入期刊发表费" />
                                    </el-form-item>
                                </el-col>
                                <el-col :span="12">
                                    <el-form-item label="期刊页码" prop="periodical_page">
                                        <!-- <el-input v-model="addForm.periodical_page" type="number" -->
                                        <el-input v-model="addForm.periodical_page" placeholder="请输入期刊页码" />
                                    </el-form-item>
                                </el-col>

                                <!-- 第五行 -->
                                <el-col :span="12">
                                    <el-form-item label="引用率" prop="citation_rate">
                                        <el-input v-model="addForm.citation_rate" placeholder="请输入引用率" />
                                    </el-form-item>
                                </el-col>
                                <el-col :span="12">
                                    <el-form-item label="复合影响因子" prop="composite_influence_factor">
                                        <el-input v-model="addForm.composite_influence_factor"
                                            placeholder="请输入复合影响因子" />
                                    </el-form-item>
                                </el-col>

                                <!-- 第六行 -->
                                <el-col :span="12">
                                    <el-form-item label="综合影响因子" prop="synthetic_influence_factor">
                                        <el-input v-model="addForm.synthetic_influence_factor"
                                            placeholder="请输入综合影响因子" />
                                    </el-form-item>
                                </el-col>
                                <el-col :span="12">
                                    <el-form-item label="出刊时间" prop="publish_date">
                                        <el-date-picker v-model="addForm.publish_date" type="date" placeholder="请选择出刊时间"
                                            format="YYYY-MM-DD" value-format="YYYY-MM-DD" />
                                    </el-form-item>
                                </el-col>

                                <!-- 第七行 -->
                                <el-col :span="12">
                                    <el-form-item label="创刊时间" prop="founding_time">
                                        <el-date-picker v-model="addForm.founding_time" type="year"
                                            placeholder="请选择创刊年份" format="YYYY" value-format="YYYY" />
                                    </el-form-item>
                                </el-col>
                                <el-col :span="12">
                                    <el-form-item label="收稿方向" prop="directions">
                                        <el-select v-model="addForm.directions" multiple placeholder="请选择收稿方向"
                                            class="full-width">
                                            <el-option v-for="item in directionList" :key="item.categorize_id"
                                                :label="item.name" :value="item.categorize_id" />
                                        </el-select>
                                    </el-form-item>
                                </el-col>

                                <!-- 第八行 -->
                                <el-col :span="12">
                                    <el-form-item label="网站收录" prop="collections">
                                        <el-select v-model="addForm.collections" multiple placeholder="请选择网站收录"
                                            class="full-width">
                                            <el-option v-for="item in collectionList" :key="item.categorize_id"
                                                :label="item.name" :value="item.categorize_id" />
                                        </el-select>
                                    </el-form-item>
                                </el-col>
                                <el-col :span="12">
                                    <el-form-item label="期刊级别" prop="levels">
                                        <el-select v-model="addForm.levels" multiple placeholder="请选择期刊级别"
                                            class="full-width">
                                            <el-option v-for="item in levelList" :key="item.categorize_id"
                                                :label="item.name" :value="item.categorize_id" />
                                        </el-select>
                                    </el-form-item>
                                </el-col>

                                <!-- 第九行 -->
                                <el-col :span="12">
                                    <el-form-item label="出版周期" prop="cycle">
                                        <el-select v-model="addForm.cycle" placeholder="请选择出版周期" class="full-width">
                                            <el-option v-for="item in cycleList" :key="item.categorize_id"
                                                :label="item.name" :value="item.categorize_id" />
                                        </el-select>
                                    </el-form-item>
                                </el-col>
                                <el-col :span="12">
                                    <el-form-item label="期刊批次" prop="periodical_batch">
                                        <el-select v-model="addForm.periodical_batch" placeholder="请选择期刊批次"
                                            class="full-width">
                                            <el-option v-for="item in periodicalBatchOptions" :key="item.value"
                                                :label="item.label" :value="item.value" />
                                        </el-select>
                                    </el-form-item>
                                </el-col>

                                <!-- 第十行 -->
                                <el-col :span="6">
                                    <el-form-item label="社内发票" prop="invoice_receipt">
                                        <el-radio-group v-model="addForm.invoice_receipt">
                                            <el-radio :value=1>可开取</el-radio>
                                            <el-radio :value=0>不可开取</el-radio>
                                        </el-radio-group>
                                    </el-form-item>
                                </el-col>
                                <el-col :span="6">
                                    <el-form-item label="是否全包" prop="is_warp">
                                        <el-radio-group v-model="addForm.is_warp">
                                            <el-radio :value=1>是</el-radio>
                                            <el-radio :value=0>否</el-radio>
                                        </el-radio-group>
                                    </el-form-item>
                                </el-col>
                                <el-col :span="6">
                                    <el-form-item label="可发作品" prop="works">
                                        <el-radio-group v-model="addForm.works">
                                            <el-radio :value=1>是</el-radio>
                                            <el-radio :value=0>否</el-radio>
                                        </el-radio-group>
                                    </el-form-item>
                                </el-col>
                                <el-col :span="6">
                                    <el-form-item label="是否隐藏" prop="is_hidden">
                                        <el-radio-group v-model="addForm.is_hidden">
                                            <el-radio :value=0>正常显示</el-radio>
                                            <el-radio :value=1>隐藏</el-radio>
                                        </el-radio-group>
                                    </el-form-item>
                                </el-col>
                            </el-row>

                            <!-- 单独占一行的表单项 -->
                            <el-form-item label="文章注意事项" prop="attention_matter">
                                <el-input v-model="addForm.attention_matter" type="textarea" :rows="2"
                                    placeholder="请输入文章注意事项" />
                            </el-form-item>
                            <el-form-item label="栏目设置" prop="column_setting">
                                <el-input v-model="addForm.column_setting" type="textarea" :rows="2"
                                    placeholder="请输入栏目设置" />
                            </el-form-item>
                            <el-form-item label="期刊描述" prop="periodical_description">
                                <el-input v-model="addForm.periodical_description" type="textarea" :rows="2"
                                    placeholder="请输入期刊描述" />
                            </el-form-item>

                            <el-form-item label="作者数量与样刊相关" prop="author_info">
                                <el-input v-model="addForm.author_info" type="textarea" :rows="2"
                                    placeholder="请输入作者数量与样刊相关" />
                            </el-form-item>
                            <el-form-item label="期刊清样" prop="periodical_final_proof">
                                <el-input v-model="addForm.periodical_final_proof" type="textarea" :rows="2"
                                    placeholder="请输入期刊清样" />
                            </el-form-item>
                            <el-form-item label="文章命名" prop="article_naming">
                                <el-input v-model="addForm.article_naming" type="textarea" :rows="2"
                                    placeholder="请输入文章命名" />
                            </el-form-item>

                            <el-form-item label="发表流程" prop="publication_process">
                                <el-input v-model="addForm.publication_process" type="textarea" :rows="2"
                                    placeholder="请输入发表流程" />
                            </el-form-item>
                            <el-form-item label="全包流程" prop="all_inclusive_process">
                                <el-input v-model="addForm.all_inclusive_process" type="textarea" :rows="2"
                                    placeholder="请输入全包流程" />
                            </el-form-item>
                            <el-form-item label="内部流程" prop="internal_processes">
                                <el-input v-model="addForm.internal_processes" type="textarea" :rows="2"
                                    placeholder="请输入内部流程" />
                            </el-form-item>
                        </el-form>

                        <template #footer>
                            <div class="drawer-footer">
                                <el-button @click="drawerVisible = false">取消</el-button>
                                <el-button type="primary" @click="handleSubmit">确定</el-button>
                            </div>
                        </template>
                    </el-drawer>
                </div>
            </div>
        </div>

        <!-- 期刊详情对话框 -->
        <el-dialog v-model="dialogVisible" :title="detailData.name || '期刊详情'" width="85%" :before-close="handleClose"
            class="detail-dialog" destroy-on-close>
            <template #default>
                <div class="dialog-content">
                    <div class="detail-container">
                        <div class="detail-column left-column">
                            <div class="detail-item">
                                <span class="label">更新时间：</span>
                                <span class="value">{{ formatTimestamp(detailData.updated_at) }}</span>
                            </div>
                            <div class="detail-item">
                                <span class="label">收稿方向：</span>
                                <span class="value">{{ getDirections(detailData.categorizes) }}</span>
                            </div>

                            <div class="detail-item">
                                <span class="label">主管单位：</span>
                                <span class="value">{{ detailData.competent_unit }}</span>
                            </div>
                            <div class="detail-item">
                                <span class="label">主办单位：</span>
                                <span class="value">{{ detailData.organization }}</span>
                            </div>
                            <div class="detail-item">
                                <span class="label">国内刊号：</span>
                                <span class="value">{{ detailData.domestic_number }}</span>
                            </div>
                            <div class="detail-item">
                                <span class="label">国际刊号：</span>
                                <span class="value">{{ detailData.international_number }}</span>
                            </div>
                            <div class="detail-item">
                                <span class="label">创刊时间：</span>
                                <span class="value">{{ detailData.founding_time }}年</span>
                            </div>
                            <div class="detail-item">
                                <span class="label">出版周期：</span>
                                <span class="value">{{ getPeriodicalCycle(detailData.categorizes) }}</span>
                            </div>
                            <div class="detail-item">
                                <span class="label">期刊级别：</span>
                                <span class="value">{{ getPeriodicalLevel(detailData.categorizes) }}</span>
                            </div>
                            <div class="detail-item">
                                <span class="label">是否全包：</span>
                                <span class="value">{{ detailData.is_warp ? '是' : '否' }}</span>
                            </div>
                            <div class="detail-item">
                                <span class="label">引用率：</span>
                                <span class="value">{{ detailData.citation_rate }}</span>
                            </div>
                            <div class="detail-item">
                                <span class="label">复合影响因子：</span>
                                <span class="value">{{ detailData.composite_influence_factor }}</span>
                            </div>
                            <div class="detail-item">
                                <span class="label">综合影响因子：</span>
                                <span class="value">{{ detailData.synthetic_influence_factor }}</span>
                            </div>
                            <div class="detail-item">
                                <span class="label">审核周期：</span>
                                <span class="value">{{ detailData.audit_time }}</span>
                            </div>
                            <div class="detail-item">
                                <span class="label">是否可开杂志社发票：</span>
                                <span class="value">{{ detailData.invoice_receipt ? '是' : '否' }}</span>
                            </div>
                            <div v-if="userStore.hasPermission('periodical:update')" class="detail-item">
                                <span class="label">隐藏状态：</span>
                                <span class="value">
                                    <span v-if="detailData.is_hidden" style="color: #d97706; font-weight: 600;">已隐藏（AI默认检索已屏蔽）</span>
                                    <span v-else style="color: #059669;">正常显示</span>
                                </span>
                            </div>
                        </div>
                        <div class="detail-column right-column">
                            <div class="detail-item">
                                <span class="label">期刊描述：</span>
                                <span class="value">
                                    {{ detailData.periodical_description }}
                                    <el-button v-if="detailData.periodical_description" link type="primary"
                                        class="copy-btn" @click="handleCopyText(detailData.periodical_description)">
                                        <template #icon>
                                            <el-icon>
                                                <CopyDocument />
                                            </el-icon>
                                        </template>
                                        复制
                                    </el-button>
                                </span>
                            </div>
                            <div class="detail-item">
                                <span class="label">栏目设置：</span>
                                <span class="value">
                                    {{ detailData.column_setting }}
                                    <el-button v-if="detailData.column_setting" link type="primary" class="copy-btn"
                                        @click="handleCopyText(detailData.column_setting)">
                                        <template #icon>
                                            <el-icon>
                                                <CopyDocument />
                                            </el-icon>
                                        </template>
                                        复制
                                    </el-button>
                                </span>
                            </div>
                            <div class="detail-item">
                                <span class="label">注意事项：</span>
                                <span class="value">
                                    <span v-html="formatAttentionMatter(detailData.attention_matter)"></span>
                                    <el-button v-if="detailData.attention_matter" link type="primary" class="copy-btn"
                                        @click="handleCopyText(detailData.attention_matter)">
                                        <template #icon>
                                            <el-icon>
                                                <CopyDocument />
                                            </el-icon>
                                        </template>
                                        复制
                                    </el-button>
                                </span>
                            </div>
                            <div class="detail-item">
                                <span class="label">作者数量与样刊相关：</span>
                                <span v-html="formatAttentionMatter(detailData.author_info)"></span>
                                <el-button v-if="detailData.author_info" link type="primary" class="copy-btn"
                                    @click="handleCopyText(detailData.author_info)">
                                    <template #icon>
                                        <el-icon>
                                            <CopyDocument />
                                        </el-icon>
                                    </template>
                                    复制
                                </el-button>
                            </div>
                            <div class="detail-item">
                                <span class="label">期刊清样：</span>
                                <span v-html="formatAttentionMatter(detailData.periodical_final_proof)"></span>
                                <el-button v-if="detailData.periodical_final_proof" link type="primary" class="copy-btn"
                                    @click="handleCopyText(detailData.periodical_final_proof)">
                                    <template #icon>
                                        <el-icon>
                                            <CopyDocument />
                                        </el-icon>
                                    </template>
                                    复制
                                </el-button>
                            </div>
                            <div class="detail-item">
                                <span class="label">文章命名：</span>
                                <span v-html="formatAttentionMatter(detailData.article_naming)"></span>
                                <el-button v-if="detailData.article_naming" link type="primary" class="copy-btn"
                                    @click="handleCopyText(detailData.article_naming)">
                                    <template #icon>
                                        <el-icon>
                                            <CopyDocument />
                                        </el-icon>
                                    </template>
                                    复制
                                </el-button>
                            </div>
                            <div class="detail-item">
                                <span class="label">发表流程：</span>
                                <span v-html="formatAttentionMatter(detailData.publication_process)"></span>
                                <el-button v-if="detailData.publication_process" link type="primary" class="copy-btn"
                                    @click="handleCopyText(detailData.publication_process)">
                                    <template #icon>
                                        <el-icon>
                                            <CopyDocument />
                                        </el-icon>
                                    </template>
                                    复制
                                </el-button>
                            </div>
                            <div class="detail-item">
                                <span class="label">全包流程：</span>
                                <span v-html="formatAttentionMatter(detailData.all_inclusive_process)"></span>
                                <el-button v-if="detailData.all_inclusive_process" link type="primary" class="copy-btn"
                                    @click="handleCopyText(detailData.all_inclusive_process)">
                                    <template #icon>
                                        <el-icon>
                                            <CopyDocument />
                                        </el-icon>
                                    </template>
                                    复制
                                </el-button>
                            </div>
                            <div class="detail-item">
                                <span class="label">内部流程：</span>
                                <span class="value">
                                    {{ detailData.internal_processes }}
                                    <el-button v-if="detailData.internal_processes" link type="primary" class="copy-btn"
                                        @click="handleCopyText(detailData.internal_processes)">
                                        <template #icon>
                                            <el-icon>
                                                <CopyDocument />
                                            </el-icon>
                                        </template>
                                        复制
                                    </el-button>
                                </span>
                            </div>
                        </div>
                    </div>
                </div>
            </template>
        </el-dialog>

        <!-- 右下角灵动 AI 期刊伴随助手气泡 -->
        <AiCopilotBubble
            @refresh-list="getPeriodicalList"
            @sync-filter="handleAiSyncFilter"
            @fill-drawer="handleAiFillDrawer"
            @select-journal="handleAiSelectJournal"
        />
    </div>
</template>

<script setup lang="ts">
import { onMounted, onActivated, reactive, ref } from 'vue';
import { useRoute } from 'vue-router';
import { reqCategorizeList, reqPeriodicalList, reqPeriodicalDetail, reqCreatePeriodical, reqDeletePeriodical, reqUpdatePeriodical, reqUpdatePeriodicalHidden } from '@/api/periodical';
import { ElMessage, ElMessageBox } from 'element-plus';
import { View, DocumentCopy, CopyDocument, Plus } from '@element-plus/icons-vue'
import dayjs from 'dayjs' // 需要先安装 dayjs: npm install dayjs
import useUserStore from '@/store/modules/user';
import router from '@/router';
import AiCopilotBubble from '@/components/AiCopilotBubble/index.vue';
import { copyTextToClipboard } from '@/utils/clipboard';

defineOptions({
    name: 'PeriodicalManage',
})

const route = useRoute()
const userStore = useUserStore()

let categorizeListParams = reactive<any>({
    page_size: 10000,
    page_number: 1,
    status: 1
})  // 获取通用分类列表的请求参数

let periodicalListParams = reactive<any>({
    page_size: '',
    page_number: '',
    has_base_type_list: [],
    publish_month: ''
})  // 获取期刊列表的请求参数

let periodicalArr = ref<any>([]) // 存储已有期刊的数据
let directionList = ref<any>([])  // 存储接口返回的收稿方向
let collectionList = ref<any>([]) // 存储接口返回的网站收录
let levelList = ref<any>([]) // 存储期刊级别数据
let cycleList = ref<any>([]) // 存储出版周期数据

interface CheckedItem {
    base_id: number;
    type: string;
}

// 收录方向相关的响应式变量 
const checkAll = ref(false)
const isIndeterminate = ref(false)
const checkedCities = ref<CheckedItem[]>([])

// 网站收录相关的响应式变量
const collectionAll = ref(false)
const collectionIsIndeterminate = ref(false)
const checkedCollection = ref<CheckedItem[]>([])

// 期刊级别相关的响应式变量
const levelCheckAll = ref(false)
const levelIsIndeterminate = ref(false)
const checkedLevels = ref<CheckedItem[]>([])

// 出版周期相关的响应式变量
const cycleCheckAll = ref(false)
const cycleIsIndeterminate = ref(false)
const checkedCycles = ref<CheckedItem[]>([])

// 作品可、否搜索
const canPublishWorks = ref<boolean | null>(null)
// 全包、可搜索
const isWrap = ref<boolean | null>(null)
// 存储选中的学术期刊
const checkedAcademicJournals = ref<number[]>([])
// 是否可开杂志社发票
const canIssueInvoice = ref<number[]>([])
// 复合影响因子
const compositeImpactFactor = ref<number | null>(null)  // 初始值可以为 null
// 隐藏状态筛选：0 全部，2 未隐藏，1 已隐藏
const hiddenFilter = ref<number>(0)
// 行内切换隐藏状态的 Loading 映射
const hiddenLoadingMap = reactive<Record<number, boolean>>({})


// 定义默认值   文章命名、发表流程、全包流程、内部流程 新增时的默认值
const defaultValues = {
    articleNaming: '（1）初次审核：送审-日期-业务员姓名-客户姓名-版面-《刊物名字》-刊期-文章题目\n（2）修改后再审：修改送审 - 日期 - 业务员姓名 - 客户姓名 - 版面 -《刊物名字》-刊期 - 文章题目\n（3）反馈内勤已汇款：已汇款 - 日期 - 业务员姓名 - 客户姓名 - 版面 -《刊物名字》-刊期 - 文章题目\n（4）需替换文章：（替换）日期 - 已汇款 - 日期 - 业务员姓名 - 客户姓名 - 版面 -《刊物名字》-刊期 - 文章题目',
    publicationProcess: '作者投稿—安排审稿（审核费用300）—通过后发用稿通知单—支付版面费（抵扣300）—核款定版—安排排版—出刊—邮寄（到付，提前支付邮费包邮更划算）',
    allInclusiveProcess: '30%定金—提供具体方向—出题（3-5个左右）、选题—文章写好之后审核—通过出录用通知—付70%尾款（给全文）—安排排版—出刊—邮寄（到付，提前支付邮费包邮更划算）',
    internalProcess: '根据作者需求推刊-沟通选定意向刊物-收取客户稿件-整理稿件基础格式和版面，以及和作者核对信息无误-收定金300审核-发内勤（刘QQ）对接审核-通过传录用通知单（未过提供修改建议或拒稿）-收余下费用-告知内勤（刘QQ）已汇款文件名+领款（欧QQ）-更改作者备注-统计文章信息地址等在表格-出刊核对地址发单号-上网做好维护。'
}


// 更新选中状态
const handleCheckAllChange = (val: boolean) => {
    checkedCities.value = val
        ? directionList.value.map((item: { categorize_id: any; type: any; }) => ({
            base_id: item.categorize_id,
            type: item.type
        }))
        : []
    isIndeterminate.value = false

    // 更新期刊列表查询参数
    periodicalListParams.has_base_type_list = [
        ...checkedCities.value,
        ...checkedCollection.value,
        ...checkedLevels.value,
        ...checkedCycles.value
    ]
    getPeriodicalList() // 调用接口获取最新数据
}

const handleCheckedCitiesChange = (value: CheckedItem[]) => {
    const checkedCount = value.length
    checkAll.value = checkedCount === directionList.value.length
    isIndeterminate.value = checkedCount > 0 && checkedCount < directionList.value.length

    periodicalListParams.has_base_type_list = [
        ...value,
        ...checkedCollection.value,
        ...checkedLevels.value,
        ...checkedCycles.value
    ]
    getPeriodicalList()
}

// 网站收录全选处理
const handleCollectionCheckAllChange = (val: boolean) => {
    checkedCollection.value = val
        ? collectionList.value.map((item: { categorize_id: any; type: any; }) => ({
            base_id: item.categorize_id,
            type: item.type
        }))
        : []
    collectionIsIndeterminate.value = false

    // 更新期刊列表查询参数
    periodicalListParams.has_base_type_list = [
        ...checkedCities.value,
        ...checkedCollection.value,
        ...checkedLevels.value,
        ...checkedCycles.value
    ]
    getPeriodicalList() // 调用接口获取最新数据
}

// 网站收录选择变化处理
const handleCheckedCollectionChange = (value: CheckedItem[]) => {
    const checkedCount = value.length
    collectionAll.value = checkedCount === collectionList.value.length
    collectionIsIndeterminate.value = checkedCount > 0 && checkedCount < collectionList.value.length

    // 更新期刊列表查询参数
    // 合并所有选择结果
    periodicalListParams.has_base_type_list = [
        ...checkedCities.value,
        ...value,
        ...checkedLevels.value,
        ...checkedCycles.value
    ]
    getPeriodicalList()
}

// 期刊级别全选处理
const handleLevelCheckAllChange = (val: boolean) => {
    checkedLevels.value = val
        ? levelList.value.map((item: { categorize_id: any; type: any; }) => ({
            base_id: item.categorize_id,
            type: item.type
        }))
        : []
    levelIsIndeterminate.value = false

    // 更新期刊列表查询参数
    periodicalListParams.has_base_type_list = [
        ...checkedCities.value,
        ...checkedCollection.value,
        ...checkedLevels.value,
        ...checkedCycles.value
    ]
    getPeriodicalList() // 调用接口获取最新数据
}

// 期刊级别选择变化处理
const handleCheckedLevelsChange = (value: CheckedItem[]) => {
    const checkedCount = value.length
    levelCheckAll.value = checkedCount === levelList.value.length
    levelIsIndeterminate.value = checkedCount > 0 && checkedCount < levelList.value.length

    // 更新期刊列表查询参数，合并所有选择结果
    periodicalListParams.has_base_type_list = [
        ...checkedCities.value,
        ...checkedCollection.value,
        ...checkedLevels.value,
        ...checkedCycles.value
    ]
    getPeriodicalList()
}

// 出版周期全选处理
const handleCycleCheckAllChange = (val: boolean) => {
    checkedCycles.value = val
        ? cycleList.value.map((item: { categorize_id: any; type: any; }) => ({
            base_id: item.categorize_id,
            type: item.type
        }))
        : []
    cycleIsIndeterminate.value = false

    // 更新期刊列表查询参数
    periodicalListParams.has_base_type_list = [
        ...checkedCities.value,
        ...checkedCollection.value,
        ...checkedLevels.value,
        ...checkedCycles.value
    ]
    getPeriodicalList() // 调用接口获取最新数据
}

// 出版周期选择变化处理
const handleCheckedCyclesChange = (value: CheckedItem[]) => {
    const checkedCount = value.length
    cycleCheckAll.value = checkedCount === cycleList.value.length
    cycleIsIndeterminate.value = checkedCount > 0 && checkedCount < cycleList.value.length

    // 更新期刊列表查询参数，合并所有选择结果
    periodicalListParams.has_base_type_list = [
        ...checkedCities.value,
        ...checkedCollection.value,
        ...checkedLevels.value,
        ...value
    ]
    getPeriodicalList()
}

// 根据作品进行搜索  可 否
const handlePublishWorksChange = (value: boolean) => {
    // 根据选择的值进行搜索或其他操作
    console.log('可否发布作品选择:', value, canPublishWorks.value);
    // 这里可以调用搜索函数
    handleSearch(); // 假设您有一个搜索函数
}

// 根据全包进行搜索，可 否
const handleWrapChange = (value: boolean) => {
    periodicalListParams.is_warp = value
    console.log('全包---作品选择:', value, isWrap.value);
    handleSearch()
    // getPeriodicalList()
}

// 根据期刊批次进行搜索
const handleAcademicJournalsChange = (value: number[]) => {
    console.log('选中的学术期刊:', value, checkedAcademicJournals.value);
    // 这里可以调用搜索函数
    handleSearch(); // 假设您有一个搜索函数
}

// 是否可开杂志社发票选项
const handleInvoiceChange = (value: boolean) => {
    console.log('发票选择:', value, canIssueInvoice.value);
    // 这里可以调用搜索函数或其他操作
    handleSearch(); // 假设您有一个搜索函数
}

// 复合影响因子选项
const handleCompositeImpactFactorChange = (value: number) => {
    console.log('选择的复合影响因子:', value);
    // 这里可以调用搜索函数或其他操作
    handleSearch(); // 假设您有一个搜索函数
}

// 隐藏状态筛选变化
const handleHiddenFilterChange = () => {
    handleSearch()
}

// 表格行样式（已隐藏的期刊柔和视觉区分）
const getTableRowClassName = ({ row }: { row: any }) => {
    return row?.is_hidden ? 'row-is-hidden' : ''
}

// 行内一键切换期刊是否隐藏
const handleToggleHidden = async (row: any, targetHidden: boolean) => {
    if (!row || !row.periodical_id) return
    hiddenLoadingMap[row.periodical_id] = true
    try {
        const res = await reqUpdatePeriodicalHidden(row.periodical_id, targetHidden)
        if (res.code === 200) {
            row.is_hidden = targetHidden
            row.updated_at = Math.floor(Date.now() / 1000)
            ElMessage.success(
                targetHidden
                    ? `已将 ${row.name} 设为隐藏（AI 日常检索将自动过滤）`
                    : `已取消隐藏，${row.name} 已恢复正常显示`,
            )
            // 若当前正在按“未隐藏”或“已隐藏”过滤视图，立即刷新列表保持一致
            if (hiddenFilter.value === 1 || hiddenFilter.value === 2) {
                getPeriodicalList()
            }
        } else {
            ElMessage.error(res.message || '更新隐藏状态失败')
        }
    } catch {
        ElMessage.error('更新隐藏状态失败')
    } finally {
        hiddenLoadingMap[row.periodical_id] = false
    }
}

// 获取分类列表
const getCategorizeList = async () => {
    let result = await reqCategorizeList(categorizeListParams)
    if (result.code == 200) {
        directionList.value = []
        collectionList.value = []
        levelList.value = []
        cycleList.value = []

        for (let item of result.data.items) {
            if (item.type === 1) {  // type 1  收稿方向
                directionList.value.push(item)
            } else if (item.type === 2) {  // type 2 期刊级别
                levelList.value.push(item)
            } else if (item.type === 3) {  // type 3 网站收录
                // levelList.value.push(item)
                collectionList.value.push(item)
            } else if (item.type === 4) {  // type 4 期刊级别
                cycleList.value.push(item)
            } else if (item.type === 5) { // 假设type=5为出版周期，请根据实际情况调整
                console.log(111)
            }
        }
    } else if (result.code == 402) {
        ElMessage.error('token已失效，请重新登录')
        await userStore.userLogout()
        router.push("/login")
    } else {
        ElMessage.error('获取分类列表失败')
    }
}

// 搜索表单相关
const searchForm = reactive({
    name: '',
    organization: '',
    competent_unit: '',
    domestic_number: '',
    international_number: '',
    column_setting: '', // 栏目设置搜索
    // works: '',   // 根据作品 可 否进行搜索
    publish_month: '',
    publish_period: '',
    publish_time_start: '',  // 开始时间
    publish_time_end: ''     // 结束时间
})
const searchFormRef = ref()
const isExpand = ref(false)

// 表格相关
const tableLoading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 定义期刊批次选项
const periodicalBatchOptions = [
    { label: '无', value: 0 },
    { label: '第一批次', value: 1 },
    { label: '第二批次', value: 2 }
]

// 搜索表单处理函数
const handleSearch = () => {
    currentPage.value = 1
    periodicalListParams.publish_month = searchForm.publish_month
    periodicalListParams.publish_period = searchForm.publish_period
    periodicalListParams.publish_time_start = searchForm.publish_time_start
    periodicalListParams.publish_time_end = searchForm.publish_time_end
    // periodicalListParams.works = searchForm.works
    // 根据作品 可否进行搜索
    if (canPublishWorks.value == true) {
        periodicalListParams.works = 1
    } else if (canPublishWorks.value == false) {
        periodicalListParams.works = 2
    }
    // 根据全包 是否进行搜索
    if (isWrap.value == true) {
        periodicalListParams.is_warp = 1
    } else if (isWrap.value == false) {
        periodicalListParams.is_warp = 2
    }
    // 根据期刊批次进行搜索
    if (checkedAcademicJournals.value.length > 0) {
        periodicalListParams.periodical_batch = checkedAcademicJournals.value
    } else {
        periodicalListParams.periodical_batch = []
    }
    // 是否可开杂志社发票
    if (canIssueInvoice.value.length > 0) {
        periodicalListParams.invoice = 1
    } else {
        periodicalListParams.invoice = 3
    }

    // 复合影响因子筛选条件
    if (compositeImpactFactor.value != null) {
        console.log("复合影响因子值：", compositeImpactFactor.value)
        if (compositeImpactFactor.value == 0) {
            periodicalListParams.composite_influence_factor = -99   // 表示选择的0以上的，为了解决后端默认float为0的问题
        } else {
            periodicalListParams.composite_influence_factor = compositeImpactFactor.value
        }
    } else {
        periodicalListParams.composite_influence_factor = 0
    }

    // 隐藏状态筛选（0-全部，2-未隐藏，1-已隐藏）
    periodicalListParams.is_hidden = hiddenFilter.value || 0

    getPeriodicalList()
}

const handleReset = () => {
    // 重置表单字段
    searchFormRef.value?.resetFields()
    // 手动清空所有字段
    Object.keys(searchForm).forEach(key => {
        searchForm[key] = ''
    })
    searchFormRef.value?.resetFields()
    currentPage.value = 1
    searchForm.publish_month = ''
    searchForm.publish_period = ''
    searchForm.publish_time_start = ''
    searchForm.publish_time_end = ''

    // 重置收稿方向
    checkAll.value = false
    isIndeterminate.value = false
    checkedCities.value = []

    // 重置网站收录
    collectionAll.value = false
    collectionIsIndeterminate.value = false
    checkedCollection.value = []

    // 重置期刊级别
    levelCheckAll.value = false
    levelIsIndeterminate.value = false
    checkedLevels.value = []

    // 重置出版周期
    cycleCheckAll.value = false
    cycleIsIndeterminate.value = false
    checkedCycles.value = []
    // 重置学术期刊
    checkedAcademicJournals.value = []

    // 重置发票选项
    canIssueInvoice.value = []

    // 重置作品选项
    canPublishWorks.value = null

    // 重置全包选项
    isWrap.value = null

    // 重置复合影响因子
    compositeImpactFactor.value = null

    // 重置隐藏状态筛选
    hiddenFilter.value = 0

    // 重置查询参数
    periodicalListParams.has_base_type_list = []
    periodicalListParams.works = null
    periodicalListParams.is_warp = null
    periodicalListParams.invoice = null
    periodicalListParams.is_hidden = 0
    periodicalListParams.periodical_batch = []
    periodicalListParams.composite_influence_factor = 0

    // 重新获取数据
    getPeriodicalList()
    handleSearch()
}

// 分页处理函数
const handleSizeChange = (val: number) => {
    pageSize.value = val
    getPeriodicalList()
}

const handleCurrentChange = (val: number) => {
    currentPage.value = val
    getPeriodicalList()
}

// 表格展示，处理出刊时间函数，进行格式化
// const formatDate = (dateString: string) => {
//     if (!dateString) return '';
//     const date = new Date(dateString);
//     return date.toISOString().split('T')[0]; // 格式化为 YYYY-MM-DD
// }

const formatDate = (dateString: string) => {
    if (!dateString) return '';
    const date = new Date(dateString);
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = date.getDate()

    // 判断旬
    let period = ''
    if (day <= 5) {
        period = '月初'
    } else if (day <= 10) {
        period = '上旬'
    } else if (day <= 20) {
        period = '中旬'
    } else if (day <= 28) {
        period = '下旬'
    } else {
        period = '月底'
    }

    return `${year}-${month} ${period}`
    // return date.toISOString().split('T')[0]; // 格式化为 YYYY-MM-DD
}

// 表格操作函数
const handlePeriodicalClick = (row: any) => {
    // 处理期刊名称点击
    console.log('期刊点击', row)
    handleView(row)
}

const dialogVisible = ref(false)
const detailData = ref<any>({})

// 查看详情处理函数
const handleView = async (row: any) => {
    try {
        const result = await getPeriodicalDetail(row.periodical_id)
        if (result.code === 200) {
            console.log("111")
            detailData.value = result.data
            dialogVisible.value = true
        } else {
            ElMessage.error('获取期刊详情失败')
        }
    } catch (error) {
        ElMessage.error('获取期刊详情失败')
    }
}

// 对话框关闭前的处理函数
const handleClose = () => {
    dialogVisible.value = false
    detailData.value = {}
}

// 获取出版周期
const getPeriodicalCycle = (categorizes: any[]) => {
    if (!categorizes || !Array.isArray(categorizes)) return ''
    const cycleItem = categorizes.find(item => item.type === 4)
    return cycleItem ? cycleItem.name : ''
}

// 统一的复制方法（兼容 HTTPS 与 HTTP 非安全上下文及 el-dialog 焦点锁定）
const copyToClipboard = async (text: string) => {
    if (!text) {
        ElMessage.warning('暂无内容可复制')
        return
    }

    const ok = await copyTextToClipboard(text)
    if (ok) {
        ElMessage.success('复制成功')
    } else {
        ElMessage.error('复制失败')
    }
}

const handleCopy = (row: any) => {
    if (row.periodical_description) {
        copyToClipboard(row.periodical_description)
    } else {
        ElMessage.warning('暂无期刊描述')
    }
}

const handleDelete = async (row: any) => {
    console.log(row);
    try {
        await ElMessageBox.confirm('确定要删除吗？', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning',
        });

        console.log("点击了确定删除");

        const result: any = await deletePeriodical(row.periodical_id);
        if (result.code === 200) {
            ElMessage({
                type: 'success',
                message: '删除成功',
            });
            // 这里可以添加刷新列表的逻辑，例如：
            getPeriodicalList(); // 假设您有一个获取期刊列表的函数
        } else {
            ElMessage.error('删除期刊失败');
        }
    } catch (error) {
        if (error === 'cancel') {
            console.log("点击了取消删除");
            ElMessage({
                type: 'info',
                message: '删除已取消',
            });
        } else {
            ElMessage.error('删除期刊失败');
        }
    }
}

// getPeriodicalDetail 获取期刊详情，id标签期刊id
const getPeriodicalDetail = async (id: number) => {
    let result: any = await reqPeriodicalDetail(id)
    if (result.code == 200) {
        console.log(result)
    }
    return result
}

// deletePeriodical 删除期刊
const deletePeriodical = async (id: number) => {
    let result: any = await reqDeletePeriodical(id)
    if (result.code == 200) {
        console.log(result)
    }
    return result
}


// 修改获取期刊列表函数
const getPeriodicalList = async () => {
    tableLoading.value = true
    try {
        periodicalListParams.page_size = pageSize.value
        periodicalListParams.page_number = currentPage.value
        // 添加搜索条件
        Object.assign(periodicalListParams, searchForm)

        let result = await reqPeriodicalList(periodicalListParams)
        if (result.code == 200) {
            periodicalArr.value = result.data.items
            total.value = result.data.total || 0
        }
        if (result.code == 402) {
            ElMessage.error('token已失效，请重新登录')
            await userStore.userLogout()
            router.push("/login")
        }
    } finally {
        tableLoading.value = false
    }
}

// 获取期刊级别
const getPeriodicalLevel = (categorizes: any[]) => {
    if (!categorizes || !Array.isArray(categorizes)) return ''

    const levelItem = categorizes.find(item => item.type === 2)
    return levelItem ? levelItem.name : ''
}

// 获取期刊收录网站列表（用于表格微徽标展示，最多取前3个）
const getCollections = (categorizes: any[]): string[] => {
    if (!categorizes || !Array.isArray(categorizes)) return []
    return categorizes
        .filter(item => item.type === 3 && item.name)
        .map(item => item.name)
        .slice(0, 3)
}

// 根据期刊级别返回语义化胶囊类名
const getLevelBadgeClass = (levelName: string): string => {
    if (!levelName) return 'level-default'
    if (levelName.includes('国家') || levelName.includes('核心')) return 'level-national'
    if (levelName.includes('国际') || levelName.includes('SCI') || levelName.includes('EI')) return 'level-intl'
    if (levelName.includes('学报') || levelName.includes('学术')) return 'level-academic'
    if (levelName.includes('省')) return 'level-provincial'
    return 'level-default'
}

// 格式化时间戳
const formatTimestamp = (timestamp: number): string => {
    if (!timestamp) return ''
    return dayjs(timestamp * 1000).format('YYYY-MM-DD HH:mm:ss')
}



// 期刊表格 排序
const sortChange = (item: any) => {
    // console.log("根据出刊时间排序。。。 ", item)
    console.log("排序 ", item, item.order, item.prop)
    if (item.order != null) {
        if (item.prop === "publish_time") {
            periodicalListParams.order_by = "publication_time"

        } else if (item.prop === "price") {
            periodicalListParams.order_by = "price"
        } else if (item.prop === "periodical_page") {
            periodicalListParams.order_by = "periodical_page"
        } else if (item.prop === "updated_at") {
            periodicalListParams.order_by = "updated_at"
        }

        // 调整排序 方式（顺序还是逆序）
        if (item.order == "ascending") {
            periodicalListParams.order_desc = "ascending"
        } else if (item.order == "descending") {
            periodicalListParams.order_desc = "descending"
        }
    } else {
        periodicalListParams.order_by = ""
        periodicalListParams.order_desc = ""
    }

    getPeriodicalList()
}

// 获取收稿方向
const getDirections = (categorizes: any[]): string => {
    if (!categorizes || !Array.isArray(categorizes)) return ''

    const directions = categorizes
        .filter(item => item.type === 1)
        .map(item => item.name)

    return directions.join('、')
}

// 处理注意事项的换行
const formatAttentionMatter = (text: string): string => {
    if (!text) return ''
    // 将 \n 替换为 HTML 的换行标签，并保留原有的换行样式
    return text.replace(/\n/g, '<br>')
}

// 检查是否携带来自 AI 工作台的一次性联动筛选参数
const checkRouteAiFilter = () => {
    if (route.query.from_ai === '1') {
        if (typeof route.query.name === 'string') {
            searchForm.name = route.query.name
        }
        if (typeof route.query.competent_unit === 'string') {
            searchForm.competent_unit = route.query.competent_unit
        }
        if (typeof route.query.organization === 'string') {
            searchForm.organization = route.query.organization
        }
        handleSearch()
        router.replace({ path: '/periodical', query: {} })
        return true
    }
    return false
}

// 组件挂载完毕钩子
onMounted(() => {
    getCategorizeList()  // 页面挂载时候获取分类列表
    if (!checkRouteAiFilter()) {
        getPeriodicalList()
    }
})

// Keep-Alive 激活钩子：普通切回标签时保留用户筛选与弹窗状态；若从 AI 工作台带参数跳转则联动筛选
onActivated(() => {
    checkRouteAiFilter()
})

// 复制文本到剪贴板
const handleCopyText = (text: string) => {
    copyToClipboard(text)
}

// 处理月份变化
const handleMonthChange = (val: string | null) => {
    // 清空出刊时间的选择
    searchForm.publish_period = ''

    if (!val) {
        // 清空时间范围
        searchForm.publish_time_start = ''
        searchForm.publish_time_end = ''
        handleSearch()
        return
    }

    // 设置选择月份的起止时间
    const monthDate = dayjs(val)
    searchForm.publish_time_start = monthDate.startOf('month').format('YYYY-MM-DD')
    searchForm.publish_time_end = monthDate.endOf('month').format('YYYY-MM-DD')

    handleSearch()
}

// 处理出刊时间变化
const handlePeriodChange = (val: string | null) => {
    if (!val || !searchForm.publish_month) {
        searchForm.publish_time_start = ''
        searchForm.publish_time_end = ''
        handleSearch()
        return
    }

    const monthDate = dayjs(searchForm.publish_month)

    // 根据不同时间段设置起止时间
    switch (val) {
        case '1': // 月初
            searchForm.publish_time_start = monthDate.date(1).format('YYYY-MM-DD')
            searchForm.publish_time_end = monthDate.date(5).format('YYYY-MM-DD')
            break
        case '10': // 上旬
            searchForm.publish_time_start = monthDate.date(1).format('YYYY-MM-DD')
            searchForm.publish_time_end = monthDate.date(10).format('YYYY-MM-DD')
            break
        case '15': // 中旬
            searchForm.publish_time_start = monthDate.date(11).format('YYYY-MM-DD')
            searchForm.publish_time_end = monthDate.date(20).format('YYYY-MM-DD')
            break
        case '20': // 下旬
            searchForm.publish_time_start = monthDate.date(21).format('YYYY-MM-DD')
            // searchForm.publish_time_end = monthDate.endOf('month').format('YYYY-MM-DD')
            searchForm.publish_time_end = monthDate.endOf('month').format('YYYY-MM-DD')
            break
        case 'end': // 月底
            const endDay = monthDate.endOf('month').date()
            searchForm.publish_time_start = monthDate.date(26).format('YYYY-MM-DD')
            searchForm.publish_time_end = monthDate.date(endDay).format('YYYY-MM-DD')
            break
    }

    getPeriodicalList()
}

// 抽屉显示状态
const drawerVisible = ref(false)

// 新增表单数据
const addForm = reactive({
    name: '',
    directions: [],  // 收稿方向
    collections: [],  // 网站收录
    levels: [],  // 期刊级别
    cycle: '',  // 出版周期
    organization: '',  // 主办单位
    competent_unit: '',  // 主管单位
    domestic_number: '', // 国内刊号
    international_number: '', // 国际刊号
    audit_time: '', // 审核周期
    citation_rate: '', // 引用率
    publish_date: '',  // 出刊时间
    founding_time: '', // 创刊时间
    column_setting: '', // 栏目设置
    invoice_receipt: Number(0), // 是否可以开杂志社发票
    periodical_description: '',  // 期刊描述
    attention_matter: '',// 文章注意事项
    article_naming: '', // 文章命名
    internal_processes: '', // 内部流程
    publication_process: '', // 发表流程
    all_inclusive_process: '', // 全包流程
    price: '', // 期刊发表费用
    periodical_page: '', // 请输入期刊页码
    is_warp: Number(0), // 是否全包
    composite_influence_factor: '', //复合影响因子
    synthetic_influence_factor: '', // 综合影响因子
    works: Number(0), // 可发作品
    is_hidden: Number(0), // 是否隐藏
    periodical_batch: Number(0), // 期刊批次
    periodical_id: '', //期刊ID
    periodical_final_proof: '', // 期刊清样
    author_info: '',  // 作者数量于样刊相关
})

// addFormSubmit 新增期刊时候，点击确定时候传递的参数
const addFormSubmit = reactive<any>({
    name: '',  // 期刊名字
    publication_time: '', // 出刊时间
    competent_unit: '', // 主管单位
    organization: '', // 主办单位
    domestic_number: '', // 国内刊号
    international_number: '', //国际刊号
    founding_time: '', // 创刊时间
    citation_rate: '', // 引用率
    periodical_description: '', //期刊描述
    column_setting: '',  // 栏目设置
    audit_time: '', // 审核周期
    invoice_receipt: false,  // 是否可以开杂志社发票
    attention_matter: '', // 文章注意事项
    article_naming: '', // 文章命名
    internal_processes: '', // 内部流程
    publication_process: '', // 发表流程
    all_inclusive_process: '', // 全包流程
    price: '', // 期刊发表费用
    periodical_page: '', // 请输入期刊页码
    is_warp: false, // 是否全包
    composite_influence_factor: '', //复合影响因子
    synthetic_influence_factor: '', // 综合影响因子
    works: false, // 可发作品
    is_hidden: false, // 是否隐藏
    periodical_batch: '', // 期刊批次
    periodical_final_proof: '', // 期刊清样
    author_info: '',  // 作者数量于样刊相关
    categorize_ids: [],   // 类型信息（收稿方向、网站收录、期刊级别、出版周期(半月刊、月刊)）
})

// 验证数字（包括小数）的通用方法
const validateNumber = (rule: any, value: any, callback: any) => {
    if (value === '' || value === null || value === undefined) {
        callback();
        return;
    }
    // 使用正则表达式匹配数字（包括小数）
    const reg = /^[0-9]+(\.[0-9]{1,4})?$/;
    if (!reg.test(value)) {
        callback(new Error('请输入有效的数字，最多四位小数'));
    } else {
        callback();
    }
}

// 表单校验规则
const rules = {
    name: [
        { required: true, message: '请输入期刊名称', trigger: 'blur' }
    ],
    directions: [
        { required: true, message: '请选择收稿方向', trigger: 'change' }
        // { message: '请选择收稿方向', trigger: 'change' }
    ],
    // publish_date: [
    //     { message: '请选择出刊时间', trigger: 'change' }
    // ],
    // founding_time: [
    //     { message: '请选择创刊年份', trigger: 'change' }
    // ],
    // cycle: [
    //     { message: '请选择出版周期', trigger: 'change' }
    // ],
    // price: [
    //     { validator: validateNumber, trigger: 'blur' }
    // ],
    // composite_influence_factor: [
    //     { validator: validateNumber, trigger: 'blur' }
    // ],
    // synthetic_influence_factor: [
    //     { validator: validateNumber, trigger: 'blur' }
    // ]
}

// 新增按钮点击处理
const handleAdd = () => {
    // 重置表单数据
    Object.assign(addForm, {
        name: '',
        directions: [],  // 收稿方向
        collections: [],  // 网站收录
        levels: [],  // 期刊级别
        cycle: '',  // 出版周期
        organization: '',  // 主办单位
        competent_unit: '',  // 主管单位
        domestic_number: '', // 国内刊号
        international_number: '', // 国际刊号
        audit_time: '', // 审核周期
        citation_rate: '', // 引用率
        publish_date: '',  // 出刊时间
        founding_time: '', // 创刊时间
        column_setting: '', // 栏目设置
        invoice_receipt: '', // 是否可以开杂志社发票
        periodical_description: '',  // 期刊描述
        attention_matter: '',// 文章注意事项
        price: '', // 期刊发表费用
        periodical_page: '', // 请输入期刊页码
        is_warp: '', // 是否全包
        composite_influence_factor: '', //复合影响因子
        synthetic_influence_factor: '', // 综合影响因子
        works: 0, // 可发作品
        is_hidden: 0, // 默认正常显示
        periodical_batch: '', // 期刊批次
        periodical_final_proof: '', // 期刊清样
        author_info: '',  // 作者数量于样刊相关
        article_naming: defaultValues.articleNaming, // 文章命名
        internal_processes: defaultValues.internalProcess,  // 内部流程
        publication_process: defaultValues.publicationProcess, // 发表流程
        all_inclusive_process: defaultValues.allInclusiveProcess, // 全包流程
    });

    isEditing.value = false; // 设置为新增模式
    drawerVisible.value = true;
}

const addFormRef = ref()

// 用于控制form表单是 编辑还是修改
const isEditing = ref(false);

// 表单提交处理
const handleSubmit = async () => {
    addFormRef.value?.validate(async (valid: boolean) => {
        if (!valid) {
            ElMessage.error('请填写完整表单信息')
            return
        }

        try {
            // 将表单数据填充到提交对象中
            // 处理期刊名称：如果没有《》则添加
            const name = addForm.name.trim()
            addFormSubmit.name = name
            // addFormSubmit.name = name.startsWith('《') && name.endsWith('》')
            //     ? name
            //     : `《${name}》`

            // 处理创刊年份：添加"年"
            addFormSubmit.founding_time = addForm.founding_time?.replace('年', '') || addForm.founding_time


            // 处理发票选项：转换为 boolean
            // addFormSubmit.invoice_receipt = addForm.invoice_receipt === 1
            // addFormSubmit.invoice_receipt = Boolean(addForm.invoice_receipt)
            addFormSubmit.invoice_receipt = addForm.invoice_receipt == 1 ? true : false
            addFormSubmit.is_warp = addForm.is_warp == 1 ? true : false
            addFormSubmit.works = addForm.works == 1 ? true : false
            addFormSubmit.is_hidden = addForm.is_hidden == 1 ? true : false

            // 处理浮点数转换
            addFormSubmit.composite_influence_factor = addForm.composite_influence_factor;
            addFormSubmit.synthetic_influence_factor = addForm.synthetic_influence_factor;
            addFormSubmit.price = addForm.price;

            addFormSubmit.periodical_batch = parseInt(addForm.periodical_batch) || 0;

            // 复制其他字段
            addFormSubmit.publication_time = addForm.publish_date
            addFormSubmit.competent_unit = addForm.competent_unit
            addFormSubmit.organization = addForm.organization
            addFormSubmit.domestic_number = addForm.domestic_number
            addFormSubmit.international_number = addForm.international_number
            addFormSubmit.citation_rate = addForm.citation_rate
            addFormSubmit.periodical_description = addForm.periodical_description
            addFormSubmit.column_setting = addForm.column_setting
            addFormSubmit.audit_time = addForm.audit_time
            addFormSubmit.attention_matter = addForm.attention_matter
            addFormSubmit.article_naming = addForm.article_naming
            addFormSubmit.internal_processes = addForm.internal_processes
            addFormSubmit.publication_process = addForm.publication_process
            addFormSubmit.all_inclusive_process = addForm.all_inclusive_process
            addFormSubmit.periodical_page = addForm.periodical_page
            addFormSubmit.periodical_final_proof = addForm.periodical_final_proof
            addFormSubmit.author_info = addForm.author_info

            // 合并所有分类ID到categorize_ids
            addFormSubmit.categorize_ids = [
                ...addForm.directions,      // 收稿方向
                ...addForm.collections,     // 网站收录
                ...addForm.levels,          // 期刊级别
                addForm.cycle              // 出版周期（单选，不是数组）
            ].filter(Boolean)  // 过滤掉可能的空值

            // 判断是新增还是编辑
            if (isEditing.value) {
                // 编辑期刊
                addFormSubmit.periodical_id = addForm.periodical_id
                addFormSubmit.update_mode = 0   // 更新模式
                const result = await reqUpdatePeriodical(addFormSubmit.periodical_id, addFormSubmit)
                if (result.code === 200) {
                    ElMessage.success('编辑期刊成功')
                    drawerVisible.value = false // 关闭抽屉
                    getPeriodicalList() // 刷新期刊列表
                    // 重置表单
                    addFormRef.value?.resetFields()
                    Object.assign(addForm, {
                        name: '',
                        directions: [],
                        collections: [],
                        levels: [],
                        cycle: '',
                        organization: '',
                        competent_unit: '',
                        domestic_number: '',
                        international_number: '',
                        founding_time: '',
                        citation_rate: '',
                        periodical_description: '',
                        column_setting: '',
                        influence_factor: '',
                        audit_time: '',
                        invoice_receipt: 1,
                        attention_matter: '',
                        article_naming: '',
                        internal_processes: '',
                        publication_process: '',
                        all_inclusive_process: '',
                        publish_date: ''
                    })
                } else {
                    ElMessage.error(result.message || '编辑期刊失败')
                }
            } else {
                // 新增期刊
                // 调用创建期刊API
                const result = await reqCreatePeriodical(addFormSubmit)

                if (result.code === 200) {
                    ElMessage.success('创建期刊成功')
                    drawerVisible.value = false  // 关闭抽屉
                    getPeriodicalList()  // 刷新期刊列表
                    // 重置表单
                    addFormRef.value?.resetFields()
                    Object.assign(addForm, {
                        name: '',
                        directions: [],
                        collections: [],
                        levels: [],
                        cycle: '',
                        organization: '',
                        competent_unit: '',
                        domestic_number: '',
                        international_number: '',
                        founding_time: '',
                        citation_rate: '',
                        periodical_description: '',
                        column_setting: '',
                        influence_factor: '',
                        audit_time: '',
                        invoice_receipt: 1,
                        attention_matter: '',
                        article_naming: '',
                        internal_processes: '',
                        publication_process: '',
                        all_inclusive_process: '',
                        publish_date: ''
                    })
                } else {
                    ElMessage.error(result.message || '创建期刊失败')
                }
            }
        } catch (error) {
            console.error('创建期刊出错：', error)
            ElMessage.error('创建期刊失败')
        }
    })
}

const handleEdit = (row: any) => {

    // 处理分类数据
    if (row.categorizes && Array.isArray(row.categorizes)) {
        console.log("row.categorizes", row.categorizes)
        addForm.directions = row.categorizes
            .filter((item: { type: number; }) => item.type === 1)
            .map((item: { categorize_id: any; }) => item.categorize_id);
        addForm.collections = row.categorizes
            .filter((item: { type: number; }) => item.type === 3)
            .map((item: { categorize_id: any; }) => item.categorize_id);
        addForm.levels = row.categorizes
            .filter((item: { type: number; }) => item.type === 2)
            .map((item: { categorize_id: any; }) => item.categorize_id);

        // 出版周期是单选，所以取第一个匹配的值
        const cycleItem = row.categorizes.find((item: { type: number; }) => item.type === 4);
        addForm.cycle = cycleItem ? cycleItem.categorize_id : '';
    }

    // 处理出刊时间
    if (row.publication_time) {
        // 将 ISO 格式的日期字符串转换为 YYYY-MM-DD 格式
        addForm.publish_date = dayjs(row.publication_time).format('YYYY-MM-DD');
    }

    // 将当前期刊信息填充到 addForm
    addForm.name = row.name;
    addForm.organization = row.organization;
    addForm.competent_unit = row.competent_unit;
    addForm.domestic_number = row.domestic_number;
    addForm.international_number = row.international_number;
    addForm.audit_time = row.audit_time;
    addForm.price = row.price;
    addForm.periodical_page = row.periodical_page;
    addForm.citation_rate = row.citation_rate;
    addForm.article_naming = row.article_naming;
    addForm.composite_influence_factor = row.composite_influence_factor;
    addForm.synthetic_influence_factor = row.synthetic_influence_factor;
    addForm.founding_time = row.founding_time;
    addForm.periodical_batch = Number(row.periodical_batch);
    // 确保这些字段的值是数字类型
    // 使用 Number 函数确保是数值类型
    addForm.invoice_receipt = Number(row.invoice_receipt === true || row.invoice_receipt === 1 ? 1 : 0);
    addForm.is_warp = Number(row.is_warp === true || row.is_warp === 1 ? 1 : 0);
    addForm.works = Number(row.works === true || row.works === 1 ? 1 : 0);
    addForm.is_hidden = Number(row.is_hidden === true || row.is_hidden === 1 ? 1 : 0);

    addForm.attention_matter = row.attention_matter;
    addForm.column_setting = row.column_setting;
    addForm.periodical_description = row.periodical_description;
    addForm.internal_processes = row.internal_processes;
    addForm.publication_process = row.publication_process;
    addForm.all_inclusive_process = row.all_inclusive_process;
    addForm.periodical_id = row.periodical_id;
    addForm.periodical_final_proof = row.periodical_final_proof;
    addForm.author_info = row.author_info;

    // 打开抽屉
    drawerVisible.value = true;
    isEditing.value = true; // 设置为编辑模式
}

// ===== AI 伴随助手双向联动处理函数 =====

// 1. 将 AI 解析出的自然语言筛选条件一键同步到左侧多维筛选栏与顶部搜索框
const handleAiSyncFilter = (filterPayload: Record<string, any>) => {
    if (!filterPayload) return
    searchForm.name = filterPayload.name || ''
    searchForm.competent_unit = filterPayload.competent_unit || ''
    searchForm.organization = filterPayload.organization || ''
    searchForm.domestic_number = filterPayload.domestic_number || ''
    searchForm.column_setting = filterPayload.column_setting || ''

    if (filterPayload.publish_month) {
        searchForm.publish_month = filterPayload.publish_month
        const mDate = dayjs(filterPayload.publish_month)
        searchForm.publish_time_start = mDate.startOf('month').format('YYYY-MM-DD')
        searchForm.publish_time_end = mDate.endOf('month').format('YYYY-MM-DD')
    }

    canIssueInvoice.value = filterPayload.invoice === 1 ? [1] : []
    periodicalListParams.invoice = filterPayload.invoice || 3

    if (filterPayload.is_warp === 1) {
        isWrap.value = true
        periodicalListParams.is_warp = 1
    } else {
        isWrap.value = null
        periodicalListParams.is_warp = 3
    }

    if (filterPayload.works === 1) {
        canPublishWorks.value = true
        periodicalListParams.works = 1
    } else {
        canPublishWorks.value = null
        periodicalListParams.works = 3
    }

    if (filterPayload.is_hidden === 1 || filterPayload.is_hidden === 2) {
        hiddenFilter.value = filterPayload.is_hidden
        periodicalListParams.is_hidden = filterPayload.is_hidden
    } else {
        hiddenFilter.value = 0
        periodicalListParams.is_hidden = 0
    }

    const baseList = Array.isArray(filterPayload.has_base_type_list) ? filterPayload.has_base_type_list : []
    checkedCities.value = baseList.filter((item: any) => item.type === 1)
    checkedLevels.value = baseList.filter((item: any) => item.type === 2)
    checkedCollection.value = baseList.filter((item: any) => item.type === 3)
    checkedCycles.value = baseList.filter((item: any) => item.type === 4)
    periodicalListParams.has_base_type_list = baseList

    currentPage.value = 1
    getPeriodicalList()
    ElMessage.success('已将 AI 筛选条件同步至左侧多维面板并刷新列表')
}

// 2. 将 AI 从非结构化文本提取出的 20+ 个期刊字段自动预填进“新增期刊”抽屉
const handleAiFillDrawer = (formPayload?: Record<string, any>) => {
    if (!formPayload) return
    handleAdd()
    addForm.name = formPayload.name || ''
    addForm.directions = Array.isArray(formPayload.directions) ? formPayload.directions : []
    addForm.collections = Array.isArray(formPayload.collections) ? formPayload.collections : []
    addForm.levels = Array.isArray(formPayload.levels) ? formPayload.levels : []
    addForm.cycle = formPayload.cycle || ''
    addForm.organization = formPayload.organization || ''
    addForm.competent_unit = formPayload.competent_unit || ''
    addForm.domestic_number = formPayload.domestic_number || ''
    addForm.international_number = formPayload.international_number || ''
    addForm.audit_time = formPayload.audit_time || ''
    addForm.citation_rate = formPayload.citation_rate || ''
    addForm.publish_date = formPayload.publication_time || ''
    addForm.founding_time = formPayload.founding_time || ''
    addForm.column_setting = formPayload.column_setting || ''
    addForm.invoice_receipt = formPayload.invoice_receipt ? 1 : 0
    addForm.is_warp = formPayload.is_warp ? 1 : 0
    addForm.works = formPayload.works ? 1 : 0
    addForm.is_hidden = formPayload.is_hidden ? 1 : 0
    addForm.price = formPayload.price || ''
    addForm.periodical_page = formPayload.periodical_page || ''
    addForm.composite_influence_factor = formPayload.composite_influence_factor || ''
    addForm.synthetic_influence_factor = formPayload.synthetic_influence_factor || ''
    addForm.periodical_description = formPayload.periodical_description || ''
    addForm.attention_matter = formPayload.attention_matter || ''
    addForm.author_info = formPayload.author_info || ''
    ElMessage.success('已将 AI 提取的期刊信息填入新增抽屉，请核对后点击确定提交')
}

// 3. 点击 AI 卡片中的期刊项直接打开详情弹窗
const handleAiSelectJournal = (item: any) => {
    if (!item || !item.periodical_id) return
    handleView({ periodical_id: item.periodical_id })
}

</script>

<style scoped lang="scss">
.periodical-page {
    width: 100%;
}

.column {
    width: 100%;
    height: 100%;
    display: flex;
    gap: 16px;
    align-items: flex-start;

    .column-left {
        width: 288px;
        flex-shrink: 0;
        padding: 16px;
        background: #ffffff;
        border: 1px solid #e2e8f0;
        border-radius: 12px;
        box-shadow: 0 1px 3px rgba(15, 23, 42, 0.04);

        .filter-panel-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding-bottom: 12px;
            margin-bottom: 14px;
            border-bottom: 1px solid #f1f5f9;

            .panel-title {
                font-size: 14px;
                font-weight: 700;
                color: #0f172a;
                letter-spacing: -0.01em;
            }
        }

        .date-filter-box {
            background: #f8fafc;
            border: 1px solid #f1f5f9;
            border-radius: 8px;
            padding: 10px 10px 2px;
        }

        .form-label {
            display: block;
            margin-bottom: 5px;
            color: #475569;
            font-size: 12px;
            font-weight: 600;
        }

        .el-form {
            .el-row {
                margin-bottom: 8px;
            }

            .el-form-item {
                margin-bottom: 0;

                :deep(.el-form-item__content) {
                    .month-picker,
                    .period-select {
                        width: 100%;
                    }
                }
            }
        }

        .direction-wrapper {
            margin-top: 14px;
            padding-top: 12px;
            border-top: 1px solid #f1f5f9;

            &.compact-dual-row {
                display: grid;
                grid-template-columns: 1fr 1fr;
                gap: 10px;

                .dual-col {
                    background: #f8fafc;
                    border: 1px solid #f1f5f9;
                    border-radius: 8px;
                    padding: 8px 10px;

                    .direction-header {
                        margin-bottom: 4px;
                    }

                    :deep(.el-radio-group) {
                        display: flex;
                        gap: 10px;
                    }

                    :deep(.el-radio) {
                        margin-right: 0;
                        height: 24px;
                    }
                }
            }

            .direction-header {
                display: flex;
                align-items: center;
                margin-bottom: 8px;
                justify-content: space-between;

                .direction-title {
                    font-size: 13px;
                    font-weight: 600;
                    color: #0f172a;
                    display: flex;
                    align-items: center;
                    gap: 6px;

                    &::before {
                        content: '';
                        width: 3px;
                        height: 12px;
                        border-radius: 2px;
                        background: #2563eb;
                    }

                    .count-badge {
                        display: inline-flex;
                        align-items: center;
                        justify-content: center;
                        min-width: 18px;
                        height: 18px;
                        padding: 0 5px;
                        border-radius: 999px;
                        background: #eff6ff;
                        color: #2563eb;
                        font-size: 11px;
                        font-weight: 700;
                    }
                }

                .check-all {
                    margin-bottom: 0;
                    height: 22px;

                    :deep(.el-checkbox__label) {
                        font-size: 12px;
                        color: #64748b;
                    }
                }
            }

            .direction-group {
                display: flex;
                flex-wrap: wrap;
                gap: 4px 6px;

                .direction-item {
                    width: calc(50% - 3px);
                    margin-right: 0;
                    padding: 0 6px;
                    border-radius: 6px;
                    transition: all 0.15s ease;

                    &:hover {
                        background-color: #f8fafc;
                    }

                    &.is-checked {
                        background-color: #eff6ff;
                    }

                    :deep(.el-checkbox__label) {
                        width: 100%;
                        font-size: 12.5px;
                        color: #475569;
                    }

                    .checkbox-label {
                        display: inline-block;
                        width: 100%;
                        overflow: hidden;
                        text-overflow: ellipsis;
                        white-space: nowrap;
                    }
                }

                :deep(.el-checkbox) {
                    height: 28px;
                    margin-right: 0;
                }
            }
        }
    }

    .column-right {
        flex-grow: 1;
        min-width: 0;
        padding: 20px 22px;
        background: #ffffff;
        border: 1px solid #e2e8f0;
        border-radius: 12px;
        box-shadow: 0 1px 3px rgba(15, 23, 42, 0.04);

        .search-form {
            padding: 16px 18px;
            background: linear-gradient(180deg, #f8fafc 0%, #f1f5f9 100%);
            border: 1px solid #e2e8f0;
            border-radius: 10px;
            margin-bottom: 18px;

            :deep(.el-form-item) {
                margin-bottom: 0;
                width: 100%;

                .el-form-item__label {
                    font-size: 13px;
                    font-weight: 600;
                    color: #334155;
                }

                .el-input {
                    width: 100%;
                }
            }

            .search-actions {
                display: flex;
                align-items: center;
                justify-content: flex-end;
                gap: 8px;

                .expand-toggle {
                    margin-left: 2px;
                    font-weight: 500;
                }
            }

            .expanded-search-row {
                margin-top: 14px;
                padding-top: 14px;
                border-top: 1px dashed #cbd5e1;
            }
        }

        .table-toolbar {
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 14px;

            .toolbar-left {
                display: flex;
                align-items: center;
                gap: 10px;

                .table-title {
                    margin: 0;
                    font-size: 15px;
                    font-weight: 700;
                    color: #0f172a;
                    letter-spacing: -0.01em;
                }

                .total-pill {
                    display: inline-flex;
                    align-items: center;
                    gap: 4px;
                    padding: 2px 10px;
                    border-radius: 999px;
                    background: #f1f5f9;
                    border: 1px solid #e2e8f0;
                    font-size: 12px;
                    color: #64748b;

                    strong {
                        color: #2563eb;
                        font-weight: 700;
                        font-variant-numeric: tabular-nums;
                    }
                }
            }
        }

        .pagination-container {
            margin-top: 18px;
            padding-top: 14px;
            border-top: 1px solid #f1f5f9;
            display: flex;
            justify-content: flex-end;
        }

        .drawer-form {
            padding: 20px;

            :deep(.el-form-item) {
                margin-bottom: 18px;

                .el-form-item__label {
                    pointer-events: none;
                    font-weight: 500;
                    color: #334155;
                }

                .el-date-editor,
                .el-select,
                .el-input {
                    width: 100%;
                }
            }

            .el-row {
                margin-bottom: 0;
            }
        }

        .drawer-footer {
            padding: 12px 20px;
            text-align: right;
            border-top: 1px solid #e2e8f0;
        }
    }
}

/* 表格内部高级排版样式 */
.table-container {
    width: 100%;
    background: #fff;
    border: 1px solid #e2e8f0;
    border-radius: 10px;
    overflow: hidden;

    :deep(.el-table) {
        --el-table-border-color: #f1f5f9;
        --el-table-header-bg-color: #f8fafc;
        --el-table-row-hover-bg-color: #f1f5f9;

        th.el-table__cell {
            background-color: #f8fafc !important;
            color: #475569;
            font-weight: 600;
            font-size: 12.5px;
            height: 44px;
            border-bottom: 1px solid #e2e8f0;
        }

        td.el-table__cell {
            padding: 11px 0;
            border-bottom: 1px solid #f1f5f9;
        }

        tr.row-is-hidden td.el-table__cell {
            background-color: #fffbeb !important;
        }

        &::before {
            display: none;
        }
    }
}

.journal-title-cell {
    display: flex;
    flex-direction: column;
    gap: 5px;
    padding: 2px 0;

    .journal-link {
        color: #1e293b;
        font-weight: 600;
        font-size: 14px;
        cursor: pointer;
        transition: color 0.15s ease;
        width: fit-content;

        &:hover {
            color: #2563eb;
            text-decoration: underline;
            text-underline-offset: 3px;
        }

        &.is-hidden-title {
            color: #64748b;
        }
    }

    .journal-meta-badges {
        display: flex;
        flex-wrap: wrap;
        gap: 4px;

        .mini-badge {
            display: inline-flex;
            align-items: center;
            height: 18px;
            padding: 0 6px;
            border-radius: 4px;
            font-size: 11px;
            font-weight: 500;
            line-height: 1;

            &.amber-hidden {
                background: #fef3c7;
                color: #b45309;
                border: 1px solid #fcd34d;
                font-weight: 600;
            }

            &.emerald {
                background: #ecfdf5;
                color: #059669;
                border: 1px solid #a7f3d0;
            }

            &.sky {
                background: #f0f9ff;
                color: #0284c7;
                border: 1px solid #bae6fd;
            }

            &.slate {
                background: #f8fafc;
                color: #64748b;
                border: 1px solid #e2e8f0;
            }
        }
    }
}

.level-pill {
    display: inline-flex;
    align-items: center;
    padding: 3px 10px;
    border-radius: 999px;
    font-size: 12px;
    font-weight: 600;
    line-height: 1.4;

    &.level-national {
        background: #eff6ff;
        color: #1d4ed8;
        border: 1px solid #bfdbfe;
    }

    &.level-intl {
        background: #f5f3ff;
        color: #6d28d9;
        border: 1px solid #ddd6fe;
    }

    &.level-academic {
        background: #ecfeff;
        color: #0e7490;
        border: 1px solid #a5f3fc;
    }

    &.level-provincial {
        background: #f8fafc;
        color: #334155;
        border: 1px solid #cbd5e1;
    }

    &.level-default {
        background: #f8fafc;
        color: #475569;
        border: 1px solid #e2e8f0;
    }
}

.page-chip {
    display: inline-block;
    padding: 2px 8px;
    border-radius: 6px;
    background: #f8fafc;
    color: #334155;
    font-size: 12.5px;
    font-weight: 500;
    font-variant-numeric: tabular-nums;
}

.date-chip {
    display: inline-block;
    color: #334155;
    font-size: 13px;
    font-weight: 500;
    font-variant-numeric: tabular-nums;
}

.price-text {
    color: #0f172a;
    font-weight: 700;
    font-size: 13.5px;
    font-variant-numeric: tabular-nums;

    .currency {
        font-style: normal;
        font-size: 12px;
        color: #64748b;
        margin-right: 2px;
        font-weight: 500;
    }
}

.price-muted {
    color: #94a3b8;
    font-size: 12.5px;
}

.time-muted {
    color: #64748b;
    font-size: 12.5px;
    font-variant-numeric: tabular-nums;
}

.empty-dash {
    color: #cbd5e1;
    font-weight: 500;
}

.operation-wrapper {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
}

.action-icon-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    padding: 0;
    border-radius: 7px;
    border: 1px solid #e2e8f0;
    background: #f8fafc;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.16s ease;
    flex-shrink: 0;

    .el-icon {
        font-size: 14px;
    }

    &:hover {
        transform: translateY(-1px);
    }

    &.view {
        color: #475569;

        &:hover {
            color: #0f172a;
            background: #f1f5f9;
            border-color: #cbd5e1;
        }
    }

    &.copy {
        color: #0284c7;
        background: #f0f9ff;
        border-color: #bae6fd;

        &:hover {
            color: #0369a1;
            background: #e0f2fe;
            border-color: #7dd3fc;
        }
    }

    &.edit {
        color: #2563eb;
        background: #eff6ff;
        border-color: #bfdbfe;

        &:hover {
            color: #1d4ed8;
            background: #dbeafe;
            border-color: #93c5fd;
        }
    }

    &.delete {
        color: #dc2626;
        background: #fef2f2;
        border-color: #fecaca;

        &:hover {
            color: #b91c1c;
            background: #fee2e2;
            border-color: #fca5a5;
        }
    }
}

.detail-dialog {
    .el-dialog__body {
        min-width: 300px;
        max-width: 600px;
        overflow: auto;
    }
}

.dialog-content {
    width: 100%;
}

.detail-container {
    display: flex;
    gap: 18px;

    .detail-column {
        padding: 20px;
        border: 1px solid #e2e8f0;
        border-radius: 10px;
        background: #f8fafc;

        &.left-column {
            flex: 0 0 320px;

            .detail-item {
                margin-bottom: 14px;

                &:last-child {
                    margin-bottom: 0;
                }
            }
        }

        &.right-column {
            flex: 1;
            min-width: 0;
            background: #ffffff;

            .detail-item {
                margin-bottom: 14px;
                padding-bottom: 12px;
                border-bottom: 1px dashed #e2e8f0;

                &:last-child {
                    border-bottom: none;
                    margin-bottom: 0;
                    padding-bottom: 0;
                }
            }
        }

        .detail-item {
            line-height: 1.65;

            &.full-width {
                display: flex;
                flex-direction: column;
                gap: 6px;
            }

            .label {
                color: #64748b;
                font-weight: 600;
                font-size: 13px;
                margin-right: 8px;
                white-space: nowrap;
            }

            .value {
                color: #0f172a;
                font-size: 13.5px;
                word-break: break-all;
                white-space: pre-wrap;

                :deep(br) {
                    display: block;
                    margin: 8px 0;
                }
            }
        }
    }
}

:deep(.el-dialog) {
    .el-dialog__header {
        margin-right: 0;
        border-bottom: 1px solid #e2e8f0;
        padding-bottom: 15px;
    }

    .el-dialog__body {
        padding: 24px;
        max-height: 72vh;
        overflow-y: auto;
    }
}

.table-container {
    width: 100%;
    background: #fff;
}

.copy-btn {
    margin-left: 16px;
    padding: 4px 10px;
    font-size: 13px;

    .el-icon {
        margin-right: 4px;
    }
}
</style>

<style>
/* 使用更强的选择器来覆盖默认样式 */
.el-popper.is-dark.custom-tooltip,
.el-popper.is-light.custom-tooltip {
    max-width: 300px !important;
    /* 设置最大宽度 */
    width: fit-content !important;
    white-space: normal !important;
    word-break: break-word;
    line-height: 1.5;
}

.el-popper.custom-tooltip .el-popper__content {
    width: fit-content !important;
    max-width: 300px !important;
}
</style>