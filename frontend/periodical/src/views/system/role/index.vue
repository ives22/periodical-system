<template>
    <div class="role-manage-container">
        <!-- 页面标题栏 -->
        <div class="page-header">
            <div class="header-meta">
                <h2>角色与权限矩阵管理</h2>
                <p>动态创建业务角色，并通过可视化权限树配置菜单可见性与操作按钮/接口级权限</p>
            </div>
            <el-button
                v-permission="'role:create'"
                type="primary"
                icon="Plus"
                @click="handleAddRole"
            >
                新增角色
            </el-button>
        </div>

        <!-- 搜索与操作栏 -->
        <el-card shadow="never" class="search-card">
            <el-form :inline="true" :model="searchForm" class="search-form">
                <el-form-item label="关键字">
                    <el-input
                        v-model="searchForm.keyword"
                        placeholder="搜索角色名称 / 编码 / 描述"
                        clearable
                        style="width: 240px;"
                        @keyup.enter="handleSearch"
                    />
                </el-form-item>
                <el-form-item label="状态">
                    <el-select v-model="searchForm.status" placeholder="全部状态" clearable style="width: 140px;">
                        <el-option label="启用" :value="1" />
                        <el-option label="禁用" :value="0" />
                    </el-select>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" icon="Search" @click="handleSearch">搜索</el-button>
                    <el-button icon="Refresh" @click="handleReset">重置</el-button>
                </el-form-item>
            </el-form>
        </el-card>

        <!-- 角色数据表格 -->
        <el-card shadow="never" class="table-card">
            <el-table v-loading="loading" :data="roleList" stripe style="width: 100%;">
                <el-table-column prop="id" label="ID" width="70" align="center" />
                <el-table-column prop="name" label="角色名称" min-width="150">
                    <template #default="{ row }">
                        <span style="font-weight: 600; color: #0f172a;">{{ row.name }}</span>
                        <el-tag v-if="row.is_system === 1" size="small" type="danger" effect="plain" style="margin-left: 6px;">
                            系统内置
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="code" label="角色编码" min-width="140">
                    <template #default="{ row }">
                        <el-tag type="info" effect="light">{{ row.code }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="description" label="角色职责描述" min-width="230" show-overflow-tooltip />
                <el-table-column prop="sort" label="排序" width="80" align="center" />
                <el-table-column label="状态" width="100" align="center">
                    <template #default="{ row }">
                        <el-switch
                            v-model="row.status"
                            :active-value="1"
                            :inactive-value="0"
                            :disabled="row.code === 'super_admin' || !userStore.hasPermission('role:update')"
                            @change="handleStatusChange(row)"
                        />
                    </template>
                </el-table-column>
                <el-table-column label="更新时间" width="175" align="center">
                    <template #default="{ row }">
                        {{ formatTimestamp(row.updated_at || row.created_at) }}
                    </template>
                </el-table-column>
                <el-table-column label="操作" width="240" align="center" fixed="right">
                    <template #default="{ row }">
                        <el-button
                            v-permission="'role:update'"
                            type="primary"
                            link
                            icon="Setting"
                            @click="handleOpenPermDrawer(row)"
                        >
                            分配权限
                        </el-button>
                        <el-button
                            v-permission="'role:update'"
                            type="primary"
                            link
                            icon="Edit"
                            @click="handleEditRole(row)"
                        >
                            编辑
                        </el-button>
                        <el-button
                            v-permission="'role:delete'"
                            type="danger"
                            link
                            icon="Delete"
                            :disabled="row.is_system === 1 || row.code === 'super_admin'"
                            @click="handleDeleteRole(row)"
                        >
                            删除
                        </el-button>
                    </template>
                </el-table-column>
            </el-table>

            <div class="pagination-wrapper">
                <el-pagination
                    v-model:current-page="pagination.page_number"
                    v-model:page-size="pagination.page_size"
                    :page-sizes="[10, 20, 50]"
                    layout="total, sizes, prev, pager, next, jumper"
                    :total="total"
                    @size-change="fetchRoleList"
                    @current-change="fetchRoleList"
                />
            </div>
        </el-card>

        <!-- 新增 / 编辑角色弹窗 -->
        <el-dialog
            v-model="dialogVisible"
            :title="isEditing ? '编辑角色' : '新增角色'"
            width="500px"
            destroy-on-close
        >
            <el-form ref="roleFormRef" :model="roleForm" :rules="roleRules" label-width="95px">
                <el-form-item label="角色名称" prop="name">
                    <el-input v-model="roleForm.name" placeholder="例如：期刊编辑、审核主管" maxlength="32" />
                </el-form-item>
                <el-form-item label="角色编码" prop="code">
                    <el-input
                        v-model="roleForm.code"
                        placeholder="唯一英文标识，例如：editor"
                        maxlength="32"
                        :disabled="isEditing && roleForm.is_system === 1"
                    />
                </el-form-item>
                <el-form-item label="显示排序" prop="sort">
                    <el-input-number v-model="roleForm.sort" :min="0" :max="999" />
                </el-form-item>
                <el-form-item label="状态" prop="status">
                    <el-radio-group v-model="roleForm.status" :disabled="roleForm.code === 'super_admin'">
                        <el-radio :value="1">启用</el-radio>
                        <el-radio :value="0">禁用</el-radio>
                    </el-radio-group>
                </el-form-item>
                <el-form-item label="角色描述" prop="description">
                    <el-input
                        v-model="roleForm.description"
                        type="textarea"
                        :rows="3"
                        placeholder="请输入角色职责说明"
                        maxlength="200"
                        show-word-limit
                    />
                </el-form-item>
            </el-form>
            <template #footer>
                <el-button @click="dialogVisible = false">取消</el-button>
                <el-button type="primary" :loading="submitting" @click="handleSubmitRole">保存</el-button>
            </template>
        </el-dialog>

        <!-- 分配权限抽屉 -->
        <el-drawer
            v-model="permDrawerVisible"
            :title="`分配权限 - ${currentRole?.name || ''} (${currentRole?.code || ''})`"
            size="520px"
            destroy-on-close
        >
            <div v-loading="permLoading" class="perm-drawer-body">
                <el-alert
                    v-if="currentRole?.code === 'super_admin'"
                    title="系统内置超级管理员默认拥有全部菜单与操作权限，无需手动修改。"
                    type="warning"
                    :closable="false"
                    show-icon
                    style="margin-bottom: 15px;"
                />
                <div class="tree-toolbar">
                    <el-button size="small" @click="handleExpandAll(true)">全部展开</el-button>
                    <el-button size="small" @click="handleExpandAll(false)">全部折叠</el-button>
                    <el-button
                        size="small"
                        type="primary"
                        plain
                        :disabled="currentRole?.code === 'super_admin'"
                        @click="handleCheckAll(true)"
                    >
                        全选
                    </el-button>
                    <el-button
                        size="small"
                        plain
                        :disabled="currentRole?.code === 'super_admin'"
                        @click="handleCheckAll(false)"
                    >
                        全不选
                    </el-button>
                </div>

                <el-tree
                    ref="permTreeRef"
                    :data="permissionTree"
                    show-checkbox
                    node-key="id"
                    default-expand-all
                    :props="{ label: 'name', children: 'children' }"
                    class="perm-tree"
                >
                    <template #default="{ data }">
                        <div class="perm-tree-node">
                            <span>{{ data.name }}</span>
                            <div class="perm-node-meta">
                                <el-tag :type="data.type === 1 ? 'primary' : 'success'" size="small" effect="plain">
                                    {{ data.type === 1 ? '菜单' : '操作/接口' }}
                                </el-tag>
                                <span class="perm-code">{{ data.code }}</span>
                            </div>
                        </div>
                    </template>
                </el-tree>
            </div>
            <template #footer>
                <div style="text-align: right;">
                    <el-button @click="permDrawerVisible = false">取消</el-button>
                    <el-button
                        type="primary"
                        :loading="permSubmitting"
                        :disabled="currentRole?.code === 'super_admin'"
                        @click="handleSavePermissions"
                    >
                        保存权限配置
                    </el-button>
                </div>
            </template>
        </el-drawer>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import useUserStore from '@/store/modules/user'
import {
    reqRoleList,
    reqCreateRole,
    reqUpdateRole,
    reqDeleteRole,
    reqPermissionTree,
    reqRolePermissionIds,
    reqAssignRolePermissions
} from '@/api/role'

defineOptions({
    name: 'SystemRole',
})

const userStore = useUserStore()

const loading = ref(false)
const roleList = ref<any[]>([])
const total = ref(0)

const searchForm = reactive<{ keyword: string; status: number | undefined }>({
    keyword: '',
    status: undefined
})

const pagination = reactive({
    page_number: 1,
    page_size: 20
})

const formatTimestamp = (ts?: number) => {
    if (!ts) return '-'
    return new Date(ts * 1000).toLocaleString()
}

const fetchRoleList = async () => {
    loading.value = true
    try {
        const payload: any = {
            page_number: pagination.page_number,
            page_size: pagination.page_size,
            keyword: searchForm.keyword
        }
        if (searchForm.status !== undefined && searchForm.status !== null && (searchForm.status as any) !== '') {
            payload.status = Number(searchForm.status)
        }
        const res = await reqRoleList(payload)
        if (res.code === 200 && res.data) {
            roleList.value = res.data.items || []
            total.value = res.data.total || 0
        }
    } finally {
        loading.value = false
    }
}

const handleSearch = () => {
    pagination.page_number = 1
    fetchRoleList()
}

const handleReset = () => {
    searchForm.keyword = ''
    searchForm.status = undefined
    pagination.page_number = 1
    fetchRoleList()
}

const handleStatusChange = async (row: any) => {
    try {
        const res = await reqUpdateRole(row.id, {
            name: row.name,
            code: row.code,
            description: row.description,
            status: row.status,
            sort: row.sort
        })
        if (res.code === 200) {
            ElMessage.success(`角色 "${row.name}" 已${row.status === 1 ? '启用' : '禁用'}`)
        } else {
            row.status = row.status === 1 ? 0 : 1
            ElMessage.error(res.msg || '状态更新失败')
        }
    } catch (error) {
        row.status = row.status === 1 ? 0 : 1
    }
}

const dialogVisible = ref(false)
const isEditing = ref(false)
const submitting = ref(false)
const roleFormRef = ref()
const roleForm = reactive({
    id: 0,
    name: '',
    code: '',
    description: '',
    status: 1,
    sort: 10,
    is_system: 0
})

const roleRules = {
    name: [
        { required: true, message: '请输入角色名称', trigger: 'blur' },
        { min: 2, max: 32, message: '长度在 2 到 32 个字符之间', trigger: 'blur' }
    ],
    code: [
        { required: true, message: '请输入角色唯一编码', trigger: 'blur' },
        { pattern: /^[a-zA-Z0-9_:-]+$/, message: '编码仅支持字母、数字、下划线或短横线', trigger: 'blur' }
    ]
}

const handleAddRole = () => {
    isEditing.value = false
    roleForm.id = 0
    roleForm.name = ''
    roleForm.code = ''
    roleForm.description = ''
    roleForm.status = 1
    roleForm.sort = 10
    roleForm.is_system = 0
    dialogVisible.value = true
}

const handleEditRole = (row: any) => {
    isEditing.value = true
    roleForm.id = row.id
    roleForm.name = row.name
    roleForm.code = row.code
    roleForm.description = row.description
    roleForm.status = row.status
    roleForm.sort = row.sort
    roleForm.is_system = row.is_system
    dialogVisible.value = true
}

const handleSubmitRole = () => {
    roleFormRef.value?.validate(async (valid: boolean) => {
        if (!valid) return
        submitting.value = true
        try {
            let res: any
            if (isEditing.value) {
                res = await reqUpdateRole(roleForm.id, {
                    name: roleForm.name,
                    code: roleForm.code,
                    description: roleForm.description,
                    status: roleForm.status,
                    sort: roleForm.sort
                })
            } else {
                res = await reqCreateRole({
                    name: roleForm.name,
                    code: roleForm.code,
                    description: roleForm.description,
                    status: roleForm.status,
                    sort: roleForm.sort,
                    permission_ids: [1, 2, 3]
                })
            }
            if (res.code === 200) {
                ElMessage.success(isEditing.value ? '更新角色成功' : '创建角色成功')
                dialogVisible.value = false
                fetchRoleList()
            } else {
                ElMessage.error(res.msg || '保存失败')
            }
        } finally {
            submitting.value = false
        }
    })
}

const handleDeleteRole = (row: any) => {
    ElMessageBox.confirm(`确定要删除角色 "${row.name}" 吗？`, '删除确认', {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(async () => {
        const res = await reqDeleteRole(row.id)
        if (res.code === 200) {
            ElMessage.success('删除角色成功')
            fetchRoleList()
        } else {
            ElMessage.error(res.msg || '删除角色失败')
        }
    }).catch(() => {})
}

const permDrawerVisible = ref(false)
const permLoading = ref(false)
const permSubmitting = ref(false)
const currentRole = ref<any>(null)
const permissionTree = ref<any[]>([])
const permTreeRef = ref<any>(null)

const getLeafNodeIds = (nodes: any[]): number[] => {
    let leaves: number[] = []
    for (const node of nodes) {
        if (node.children && node.children.length > 0) {
            leaves = leaves.concat(getLeafNodeIds(node.children))
        } else {
            leaves.push(node.id)
        }
    }
    return leaves
}

const getAllNodeIds = (nodes: any[]): number[] => {
    let ids: number[] = []
    for (const node of nodes) {
        ids.push(node.id)
        if (node.children && node.children.length > 0) {
            ids = ids.concat(getAllNodeIds(node.children))
        }
    }
    return ids
}

const handleOpenPermDrawer = async (row: any) => {
    currentRole.value = row
    permDrawerVisible.value = true
    permLoading.value = true
    try {
        if (permissionTree.value.length === 0) {
            const treeRes = await reqPermissionTree()
            if (treeRes.code === 200) {
                permissionTree.value = treeRes.data || []
            }
        }
        const permRes = await reqRolePermissionIds(row.id)
        const checkedIds: number[] = permRes.code === 200 ? (permRes.data || []) : []
        const leafIds = new Set(getLeafNodeIds(permissionTree.value))

        await nextTick()
        if (permTreeRef.value) {
            permTreeRef.value.setCheckedKeys([])
            for (const id of checkedIds) {
                if (leafIds.has(id)) {
                    permTreeRef.value.setChecked(id, true, false)
                }
            }
        }
    } finally {
        permLoading.value = false
    }
}

const handleExpandAll = (expand: boolean) => {
    const nodes = permTreeRef.value?.store?.nodesMap
    if (!nodes) return
    Object.values(nodes).forEach((node: any) => {
        node.expanded = expand
    })
}

const handleCheckAll = (checked: boolean) => {
    if (!permTreeRef.value) return
    if (checked) {
        permTreeRef.value.setCheckedKeys(getAllNodeIds(permissionTree.value))
    } else {
        permTreeRef.value.setCheckedKeys([])
    }
}

const handleSavePermissions = async () => {
    if (!currentRole.value || !permTreeRef.value) return
    permSubmitting.value = true
    try {
        const checkedKeys: number[] = permTreeRef.value.getCheckedKeys() || []
        const halfCheckedKeys: number[] = permTreeRef.value.getHalfCheckedKeys() || []
        const finalIds = Array.from(new Set([...checkedKeys, ...halfCheckedKeys]))

        const res = await reqAssignRolePermissions(currentRole.value.id, {
            permission_ids: finalIds
        })
        if (res.code === 200) {
            ElMessage.success('角色权限分配已保存')
            permDrawerVisible.value = false
            fetchRoleList()
            await userStore.fetchUserProfile()
        } else {
            ElMessage.error(res.msg || '保存权限失败')
        }
    } finally {
        permSubmitting.value = false
    }
}

onMounted(() => {
    fetchRoleList()
})
</script>

<style scoped lang="scss">
.role-manage-container {
    display: flex;
    flex-direction: column;
    gap: 16px;

    .page-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        background: #ffffff;
        padding: 18px 22px;
        border-radius: 12px;
        border: 1px solid #e2e8f0;

        .header-meta {
            h2 {
                font-size: 18px;
                font-weight: 700;
                color: #0f172a;
                margin-bottom: 4px;
            }

            p {
                font-size: 13px;
                color: #64748b;
            }
        }
    }

    .search-card {
        .search-form {
            .el-form-item {
                margin-bottom: 0;
            }
        }
    }

    .pagination-wrapper {
        margin-top: 16px;
        display: flex;
        justify-content: flex-end;
    }
}

.perm-drawer-body {
    .tree-toolbar {
        margin-bottom: 15px;
        display: flex;
        gap: 8px;
    }

    .perm-tree {
        border: 1px solid #e2e8f0;
        border-radius: 8px;
        padding: 12px;
        background: #f8fafc;

        .perm-tree-node {
            flex: 1;
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding-right: 8px;
            font-size: 13.5px;

            .perm-node-meta {
                display: flex;
                align-items: center;
                gap: 8px;

                .perm-code {
                    color: #64748b;
                    font-size: 12px;
                    font-family: monospace;
                }
            }
        }
    }
}
</style>
