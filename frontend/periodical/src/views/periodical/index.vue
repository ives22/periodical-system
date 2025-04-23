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
                    <el-form>
                        <el-row :gutter="20">
                            <el-col :span="12">
                                <span class="form-label">出刊月份</span>
                                <el-form-item>
                                    <el-date-picker v-model="searchForm.publish_month" type="month" placeholder="选择月份"
                                        format="YYYY-MM" value-format="YYYY-MM" :clearable="true"
                                        @change="handleMonthChange" class="month-picker" />
                                </el-form-item>
                            </el-col>
                            <el-col :span="12">
                                <span class="form-label">出刊时间</span>
                                <el-form-item>
                                    <el-select v-model="searchForm.publish_period" placeholder="请选择"
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
                            <div class="direction-title">收稿方向</div>
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
                            <div class="direction-title">网站收录</div>
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
                            <div class="direction-title">期刊级别</div>
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
                            <div class="direction-title">出版周期</div>
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
                            <div class="direction-title">学术期刊</div>
                        </div>
                        <el-checkbox-group v-model="checkedAcademicJournals" @change="handleAcademicJournalsChange"
                            class="direction-group">
                            <el-checkbox :value="1" class="direction-item">第一批学术期刊</el-checkbox>
                            <el-checkbox :value="2" class="direction-item">第二批学术期刊</el-checkbox>
                        </el-checkbox-group>
                    </div>

                    <!-- 发票选项 -->
                    <div class="direction-wrapper">
                        <div class="direction-header">
                            <div class="direction-title">发票</div>
                        </div>
                        <el-checkbox-group v-model="canIssueInvoice" @change="handleInvoiceChange">
                            <el-checkbox :value="1" class="direction-item">可开杂志社发票</el-checkbox>
                        </el-checkbox-group>
                    </div>

                    <!-- 可否发布作品 -->
                    <div class="direction-wrapper">
                        <div class="direction-header">
                            <div class="direction-title">作品</div>
                        </div>
                        <el-radio-group v-model="canPublishWorks" class="direction-group"
                            @change="handlePublishWorksChange">
                            <el-radio :value="true">可</el-radio>
                            <el-radio :value="false">否</el-radio>
                        </el-radio-group>
                    </div>

                    <!-- 全包筛选 -->
                    <div class="direction-wrapper">
                        <div class="direction-header">
                            <div class="direction-title">全包</div>
                        </div>
                        <el-radio-group v-model="isWrap" class="direction-group" @change="handleWrapChange">
                            <el-radio :value="true">是</el-radio>
                            <el-radio :value="false">否</el-radio>
                        </el-radio-group>
                    </div>

                    <!-- 作品和全包筛选 -->
                    <!-- <div class="direction-wrapper">
                        <div class="filter-row">
                            <div class="filter-item">
                                <div class="direction-title">作品</div>
                                <el-radio-group v-model="canPublishWorks" class="direction-group"
                                    @change="handlePublishWorksChange">
                                    <el-radio :value="true">可</el-radio>
                                    <el-radio :value="false">否</el-radio>
                                </el-radio-group>
                            </div>
                            <div class="filter-item">
                                <div class="direction-title">全包</div>
                                <el-radio-group v-model="isWrap" class="direction-group" @change="handleWrapChange">
                                    <el-radio :value="true">是</el-radio>
                                    <el-radio :value="false">否</el-radio>
                                </el-radio-group>
                            </div>
                        </div>
                    </div> -->

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

                </div>

                <div class="column-right">
                    <!-- 搜索表单 -->
                    <el-form :model="searchForm" ref="searchFormRef" class="search-form" :inline="true"
                        @submit.prevent="handleSearch">
                        <el-row :gutter="20">
                            <el-col :span="6">
                                <el-form-item label="期刊名称" prop="name">
                                    <el-input v-model="searchForm.name" placeholder="请输入期刊名称" clearable
                                        @keyup.enter="handleSearch" />
                                </el-form-item>
                            </el-col>
                            <el-col :span="6">
                                <el-form-item label="主办单位" prop="sponsor">
                                    <el-input v-model="searchForm.organization" placeholder="请输入主办单位" clearable
                                        @keyup.enter="handleSearch" />
                                </el-form-item>
                            </el-col>
                            <el-col :span="6">
                                <el-form-item label="主管单位" prop="supervisor">
                                    <el-input v-model="searchForm.competent_unit" placeholder="请输入主管单位" clearable
                                        @keyup.enter="handleSearch" />
                                </el-form-item>
                            </el-col>

                            <el-col :span="6">
                                <el-form-item>
                                    <el-button type="primary" @click="handleSearch" plain>搜索</el-button>
                                    <el-button @click="handleReset" plain>重置</el-button>
                                    <el-button link @click="isExpand = !isExpand">
                                        {{ isExpand ? '收起' : '展开' }}
                                        <el-icon>
                                            <component :is="isExpand ? 'ArrowUp' : 'ArrowDown'" />
                                        </el-icon>
                                    </el-button>
                                </el-form-item>
                            </el-col>
                        </el-row>

                        <!-- 展开/收起的表单项 -->
                        <template v-if="isExpand">
                            <el-row :gutter="20">
                                <el-col :span="6">
                                    <el-form-item label="国内刊号" prop="domestic_number">
                                        <el-input v-model="searchForm.domestic_number" placeholder="请输入国内刊号" clearable
                                            @keyup.enter="handleSearch" />
                                    </el-form-item>
                                </el-col>
                                <el-col :span="6">
                                    <el-form-item label="国际刊号" prop="international_number">
                                        <el-input v-model="searchForm.international_number" placeholder="请输入国际刊号"
                                            clearable @keyup.enter="handleSearch" />
                                    </el-form-item>
                                </el-col>
                                <el-col :span="6">
                                    <el-form-item label="栏目设置" prop="column_setting">
                                        <el-input v-model="searchForm.column_setting" placeholder="请输入栏目设置" clearable
                                            @keyup.enter="handleSearch" />
                                    </el-form-item>
                                </el-col>
                            </el-row>
                        </template>
                    </el-form>

                    <!-- 新增操作栏 -->
                    <div class="operation-bar" v-if="userStore.loginUser.role === 1">
                        <el-button type="primary" plain @click="handleAdd">
                            <el-icon>
                                <Plus />
                            </el-icon> 新增
                        </el-button>
                    </div>

                    <!-- 期刊列表表格 -->
                    <div class="table-container">
                        <el-table :data="periodicalArr" border @sort-change="sortChange" style="width: 100%"
                            v-loading="tableLoading">
                            <el-table-column prop="name" label="期刊名称" show-overflow-tooltip>
                                <template #default="{ row }">
                                    <el-tooltip :content="row.periodical_description || '暂无期刊描述'" placement="bottom"
                                        :effect="row.periodical_description ? 'dark' : 'dark'"
                                        popper-class="custom-tooltip">
                                        <el-link type="primary" @click="handlePeriodicalClick(row)">
                                            {{ row.name }}
                                        </el-link>
                                    </el-tooltip>
                                </template>
                            </el-table-column>

                            <el-table-column label="期刊级别" show-overflow-tooltip>
                                <template #default="{ row }">
                                    {{ getPeriodicalLevel(row.categorizes) }}
                                </template>
                            </el-table-column>
                            <el-table-column label="期刊页码" prop="periodical_page" sortable="custom"
                                show-overflow-tooltip>
                                <template #default="{ row }">
                                    <span v-if="row.periodical_page">{{ row.periodical_page }}页</span>
                                </template>
                            </el-table-column>
                            <!-- <el-table-column prop="publish_time" sortable="custom" label="出刊时间" show-overflow-tooltip>
                                <template #default="{ row }">
                                    {{ formatDate(row.publication_time) }}
                                </template>
                            </el-table-column> -->
                            <el-table-column prop="publish_time" sortable="custom" label="出刊时间" show-overflow-tooltip>
                                <template #default="{ row }">
                                    {{ formatDate(row.publication_time) }}
                                </template>
                            </el-table-column>
                            <el-table-column prop="price" label="期刊发表费" sortable="custom" show-overflow-tooltip>
                                <template #default="{ row }">
                                    <!-- ¥{{ row.price.toFixed(2) }} -->
                                    ¥{{ row.price }}
                                </template>
                            </el-table-column>

                            <el-table-column prop="updated_at" label="更新时间" sortable="custom" show-overflow-tooltip>
                                <template #default="{ row }">
                                    {{ formatTimestamp(row.updated_at) }}
                                </template>
                            </el-table-column>
                            <el-table-column label="操作" width="120">
                                <template #default="{ row }">
                                    <div class="operation-wrapper">
                                        <el-tooltip content="查看详情" placement="top">
                                            <el-icon class="operation-icon" @click="handleView(row)">
                                                <View />
                                            </el-icon>
                                        </el-tooltip>
                                        <span class="operation-divider"></span>
                                        <el-tooltip content="复制期刊描述" placement="top">
                                            <el-icon class="operation-icon" @click="handleCopy(row)">
                                                <DocumentCopy />
                                            </el-icon>
                                        </el-tooltip>
                                        <!-- 管理员才能看到的操作按钮 -->
                                        <template v-if="userStore.loginUser.role === 1">
                                            <span class="operation-divider"></span>
                                            <el-tooltip content="删除" placement="top">
                                                <el-icon class="operation-icon" @click="handleDelete(row)">
                                                    <Delete />
                                                </el-icon>
                                            </el-tooltip>
                                            <span class="operation-divider"></span>
                                            <el-tooltip content="编辑期刊" placement="top">
                                                <el-icon class="operation-icon" @click="handleEdit(row)">
                                                    <Edit />
                                                </el-icon>
                                            </el-tooltip>
                                        </template>
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
                                <el-col :span="8">
                                    <el-form-item label="社内发票" prop="invoice_receipt">
                                        <el-radio-group v-model="addForm.invoice_receipt">
                                            <el-radio :value=1>可开取</el-radio>
                                            <el-radio :value=0>不可开取</el-radio>
                                        </el-radio-group>
                                    </el-form-item>
                                </el-col>
                                <el-col :span="8">
                                    <el-form-item label="是否全包" prop="is_warp">
                                        <el-radio-group v-model="addForm.is_warp">
                                            <el-radio :value=1>是</el-radio>
                                            <el-radio :value=0>否</el-radio>
                                        </el-radio-group>
                                    </el-form-item>
                                </el-col>
                                <el-col :span="8">
                                    <el-form-item label="可发作品" prop="works">
                                        <el-radio-group v-model="addForm.works">
                                            <el-radio :value=1>是</el-radio>
                                            <el-radio :value=0>否</el-radio>
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
                                <span class="label">期刊周期：</span>
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
    </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { reqCategorizeList, reqPeriodicalList, reqPeriodicalDetail, reqCreatePeriodical, reqDeletePeriodical, reqUpdatePeriodical } from '@/api/periodical';
import { ElMessage, ElMessageBox } from 'element-plus';
import { View, DocumentCopy, CopyDocument, Plus } from '@element-plus/icons-vue'
import dayjs from 'dayjs' // 需要先安装 dayjs: npm install dayjs
import useUserStore from '@/store/modules/user';
import router from '@/router';

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

    // 重置查询参数
    periodicalListParams.has_base_type_list = []
    periodicalListParams.works = null
    periodicalListParams.is_warp = null
    periodicalListParams.invoice = null
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

// 统一的复制方法
const copyToClipboard = (text: string) => {
    if (!text) {
        ElMessage.warning('暂无内容可复制')
        return
    }

    try {
        // 创建临时文本区域
        const textarea = document.createElement('textarea')
        textarea.value = text
        // 设置样式使其不可见
        textarea.style.position = 'fixed'
        textarea.style.opacity = '0'
        document.body.appendChild(textarea)
        // 选择并复制
        textarea.select()
        document.execCommand('copy')
        // 清理
        document.body.removeChild(textarea)
        ElMessage.success('复制成功')
    } catch (err) {
        console.error('复制失败:', err)
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

// 组件挂载完毕钩子
onMounted(() => {
    getCategorizeList()  // 页面挂载时候获取分类列表
    getPeriodicalList()
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
        works: '', // 可发作品
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

</script>

<style scoped lang="scss">
.periodical-page {
    width: 100%;

}

.column {
    width: 100%;
    height: 100%;
    display: flex;

    .column-left {
        width: 275px;
        height: 100%;
        padding: 10px;
        background: #fff;
        border-right: 1px solid #e4e7ed;

        .form-label {
            display: block;
            margin-bottom: 8px;
            color: #606266;
            font-size: 14px;
        }

        .el-form {
            .el-row {
                margin-bottom: 20px;
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
            margin-top: 16px;

            .direction-header {
                display: flex;
                align-items: center;
                margin-bottom: 16px;
                margin-right: 30px;
                justify-content: space-between;

                .direction-title {
                    font-size: 16px;
                    font-weight: 500;
                    color: #333;
                    margin-right: 12px;
                }

                .check-all {
                    margin-bottom: 0;
                }
            }

            .direction-group {
                display: flex;
                flex-wrap: wrap;
                gap: 12px;

                .direction-item {
                    width: calc(50% - 6px);
                    margin-right: 0;

                    :deep(.el-checkbox__label) {
                        width: 100%;
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
                    height: 32px;
                    margin-right: 0;
                }
            }
        }
    }

    .column-right {
        flex-grow: 1;
        padding: 20px;
        width: calc(100% - 275px);
        background: #fff;

        .search-form {
            padding: 20px;
            background: #f5f7fa;
            border-radius: 4px;
            margin-bottom: 20px;
        }

        .pagination-container {
            margin-top: 20px;
            display: flex;
            justify-content: flex-end;
        }

        .operation-bar {
            margin: 16px 0;
            padding: 0 20px;
        }

        .drawer-form {
            padding: 20px;

            :deep(.el-form-item) {
                margin-bottom: 18px;

                .el-form-item__label {
                    pointer-events: none;
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
            padding: 10px 20px;
            text-align: right;
        }
    }
}

.operation-wrapper {
    display: flex;
    align-items: center;
    gap: 4px;
    /* 调整图标之间的间距 */
}

.operation-icon {
    cursor: pointer;
    font-size: 16px;
}

.operation-divider {
    margin: 0 2px;
    /* 减小分隔符的间距 */
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
        border: 1px solid #e4e7ed;
        border-radius: 4px;

        &.left-column {
            flex: 0 0 320px; // 固定左侧宽度

            .detail-item {
                margin-bottom: 16px;

                &:last-child {
                    margin-bottom: 0;
                }
            }
        }

        &.right-column {
            flex: 1;
            min-width: 0; // 防止内容溢出

            .detail-item {
                margin-bottom: 12px;
            }
        }

        .detail-item {
            line-height: 1.6;

            &.full-width {
                display: flex;
                flex-direction: column;
                gap: 8px;
            }

            .label {
                color: #909399; // 更浅的文字颜色
                font-weight: normal;
                margin-right: 8px;
                white-space: nowrap;
            }

            .value {
                color: #303133; // 更深的文字颜色
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
        border-bottom: 1px solid #dcdfe6;
        padding-bottom: 15px;
    }

    .el-dialog__body {
        padding: 24px;
        max-height: 70vh; // 只添加最大高度
        overflow-y: auto; // 只添加滚动条
    }
}

.table-container {
    width: 100%;
    background: #fff;
}

.detail-container {
    .detail-column {
        .detail-item {
            .value {
                white-space: pre-wrap;
                line-height: 1.6;

                :deep(br) {
                    display: block;
                    margin: 8px 0;
                }
            }
        }
    }
}

.copy-btn {
    margin-left: 16px; // 文本后的间距
    padding: 4px 8px;
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