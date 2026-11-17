<template>
    <div class="user-manage-container">
        <!-- 页面标题栏 -->
        <div class="page-header">
            <div class="header-meta">
                <h2>用户与多角色管理</h2>
                <p>维护系统登录账号、部门备注信息，并支持为单个用户绑定一个或多个业务角色</p>
            </div>
            <el-button
                v-permission="'user:create'"
                type="primary"
                icon="Plus"
                @click="handleAddUser"
            >
                新增用户
            </el-button>
        </div>

        <!-- 搜索与操作栏 -->
        <el-card shadow="never" class="search-card">
            <el-form :inline="true" :model="searchForm" class="search-form">
                <el-form-item label="登录账号">
                    <el-input
                        v-model="searchForm.username"
                        placeholder="请输入用户名搜索"
                        clearable
                        style="width: 220px;"
                        @keyup.enter="handleSearch"
                    />
                </el-form-item>
                <el-form-item label="所属角色">
                    <el-select v-model="searchForm.role_id" placeholder="全部角色" clearable style="width: 180px;">
                        <el-option
                            v-for="r in roleOptions"
                            :key="r.id"
                            :label="r.name"
                            :value="r.id"
                        />
                    </el-select>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" icon="Search" @click="handleSearch">搜索</el-button>
                    <el-button icon="Refresh" @click="handleReset">重置</el-button>
                </el-form-item>
            </el-form>
        </el-card>

        <!-- 用户列表表格 -->
        <el-card shadow="never" class="table-card">
            <el-table v-loading="loading" :data="userList" stripe style="width: 100%;">
                <el-table-column prop="id" label="ID" width="70" align="center" />
                <el-table-column prop="username" label="登录用户名" min-width="160">
                    <template #default="{ row }">
                        <div class="user-cell">
                            <span class="user-avatar-mini">{{ (row.username || 'U').slice(0, 1).toUpperCase() }}</span>
                            <span style="font-weight: 600; color: #0f172a;">{{ row.username }}</span>
                            <el-tag
                                v-if="row.username === userStore.loginUser.name"
                                size="small"
                                type="success"
                                effect="plain"
                            >
                                当前登录
                            </el-tag>
                        </div>
                    </template>
                </el-table-column>
                <el-table-column label="已绑定角色 (多角色并集)" min-width="230">
                    <template #default="{ row }">
                        <template v-if="row.roles && row.roles.length > 0">
                            <el-tag
                                v-for="r in row.roles"
                                :key="r.id"
                                :type="r.code === 'super_admin' ? 'danger' : r.code === 'admin' ? 'primary' : 'info'"
                                size="small"
                                effect="light"
                                style="margin-right: 6px; margin-bottom: 4px;"
                            >
                                {{ r.name }}{{ r.status === 0 ? '(已禁用)' : '' }}
                            </el-tag>
                        </template>
                        <span v-else style="color: #94a3b8;">未分配角色</span>
                    </template>
                </el-table-column>
                <el-table-column label="姓名 / 部门 / 备注" min-width="230">
                    <template #default="{ row }">
                        <div v-if="row.label && Object.keys(row.label).length > 0" class="label-tags">
                            <span v-if="row.label.real_name" class="label-item"><strong>姓名：</strong>{{ row.label.real_name }}</span>
                            <span v-if="row.label.department" class="label-item"><strong>部门：</strong>{{ row.label.department }}</span>
                            <span v-if="row.label.remark" class="label-item"><strong>备注：</strong>{{ row.label.remark }}</span>
                        </div>
                        <span v-else style="color: #cbd5e1;">-</span>
                    </template>
                </el-table-column>
                <el-table-column label="创建时间" width="175" align="center">
                    <template #default="{ row }">
                        {{ formatTimestamp(row.created_at) }}
                    </template>
                </el-table-column>
                <el-table-column label="更新时间" width="175" align="center">
                    <template #default="{ row }">
                        {{ formatTimestamp(row.updated_at || row.created_at) }}
                    </template>
                </el-table-column>
                <el-table-column label="操作" width="200" align="center" fixed="right">
                    <template #default="{ row }">
                        <el-button
                            v-permission="'user:update'"
                            type="primary"
                            link
                            icon="Edit"
                            :disabled="isRowProtected(row)"
                            @click="handleEditUser(row)"
                        >
                            编辑 / 角色
                        </el-button>
                        <el-button
                            v-permission="'user:delete'"
                            type="danger"
                            link
                            icon="Delete"
                            :disabled="row.username === userStore.loginUser.name || isRowProtected(row)"
                            @click="handleDeleteUser(row)"
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
                    @size-change="fetchUserList"
                    @current-change="fetchUserList"
                />
            </div>
        </el-card>

        <!-- 新增 / 编辑用户弹窗 -->
        <el-dialog
            v-model="dialogVisible"
            :title="isEditing ? `编辑用户 - ${userForm.username}` : '新增用户'"
            width="520px"
            destroy-on-close
        >
            <el-form ref="userFormRef" :model="userForm" :rules="userRules" label-width="100px">
                <el-form-item label="用户名" prop="username">
                    <el-input
                        v-model="userForm.username"
                        placeholder="请输入登录用户名（3-20位）"
                        :disabled="isEditing"
                        maxlength="20"
                    />
                </el-form-item>
                <el-form-item :label="isEditing ? '重置密码' : '初始密码'" prop="password">
                    <el-input
                        v-model="userForm.password"
                        type="password"
                        show-password
                        :placeholder="isEditing ? '留空则不修改原密码（如需重置请输入6-20位新密码）' : '请输入初始密码（6-20位）'"
                        maxlength="20"
                    />
                </el-form-item>
                <el-form-item label="分配角色" prop="role_ids">
                    <el-select
                        v-model="userForm.role_ids"
                        multiple
                        placeholder="请选择一个或多个角色"
                        style="width: 100%;"
                    >
                        <el-option
                            v-for="r in roleOptions"
                            :key="r.id"
                            :label="`${r.name} (${r.code})`"
                            :value="r.id"
                            :disabled="r.code === 'super_admin' && !userStore.loginUser.is_super_admin"
                        />
                    </el-select>
                </el-form-item>
                <el-form-item label="真实姓名">
                    <el-input v-model="userForm.real_name" placeholder="选填，例如：张三" maxlength="32" />
                </el-form-item>
                <el-form-item label="所属部门">
                    <el-input v-model="userForm.department" placeholder="选填，例如：编辑部" maxlength="32" />
                </el-form-item>
                <el-form-item label="备注信息">
                    <el-input
                        v-model="userForm.remark"
                        type="textarea"
                        :rows="2"
                        placeholder="选填，其他备注信息"
                        maxlength="100"
                    />
                </el-form-item>
            </el-form>
            <template #footer>
                <el-button @click="dialogVisible = false">取消</el-button>
                <el-button type="primary" :loading="submitting" @click="handleSubmitUser">保存</el-button>
            </template>
        </el-dialog>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import useUserStore from '@/store/modules/user'
import { reqUserList, reqCreateUser, reqUpdateUser, reqDeleteUser } from '@/api/user'
import { reqAllEnabledRoles } from '@/api/role'

defineOptions({
    name: 'SystemUser',
})

const userStore = useUserStore()

const loading = ref(false)
const userList = ref<any[]>([])
const total = ref(0)
const roleOptions = ref<any[]>([])

const searchForm = reactive<{ username: string; role_id: number | undefined }>({
    username: '',
    role_id: undefined
})

const pagination = reactive({
    page_number: 1,
    page_size: 20
})

const formatTimestamp = (ts?: number) => {
    if (!ts) return '-'
    return new Date(ts * 1000).toLocaleString()
}

const isRowProtected = (row: any) => {
    if (userStore.loginUser.is_super_admin) return false
    const roles = row.roles || []
    return roles.some((r: any) => r.code === 'super_admin')
}

const fetchRoleOptions = async () => {
    const res = await reqAllEnabledRoles()
    if (res.code === 200) {
        roleOptions.value = res.data || []
    }
}

const fetchUserList = async () => {
    loading.value = true
    try {
        const res = await reqUserList({
            page_number: pagination.page_number,
            page_size: pagination.page_size,
            username: searchForm.username,
            role_id: Number(searchForm.role_id) || 0
        })
        if (res.code === 200 && res.data) {
            userList.value = res.data.items || []
            total.value = res.data.total || 0
        }
    } finally {
        loading.value = false
    }
}

const handleSearch = () => {
    pagination.page_number = 1
    fetchUserList()
}

const handleReset = () => {
    searchForm.username = ''
    searchForm.role_id = undefined
    pagination.page_number = 1
    fetchUserList()
}

const dialogVisible = ref(false)
const isEditing = ref(false)
const submitting = ref(false)
const userFormRef = ref()

const userForm = reactive({
    username: '',
    password: '',
    role_ids: [] as number[],
    real_name: '',
    department: '',
    remark: ''
})

const userRules = computed(() => ({
    username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
        { min: 3, max: 20, message: '用户名长度需在 3 到 20 个字符之间', trigger: 'blur' }
    ],
    password: isEditing.value
        ? [{ min: 6, max: 20, message: '密码长度需在 6 到 20 个字符之间', trigger: 'blur' }]
        : [
              { required: true, message: '请输入初始密码', trigger: 'blur' },
              { min: 6, max: 20, message: '密码长度需在 6 到 20 个字符之间', trigger: 'blur' }
          ],
    role_ids: [
        { type: 'array', required: true, min: 1, message: '请至少选择一个角色', trigger: 'change' }
    ]
}))

const handleAddUser = () => {
    isEditing.value = false
    userForm.username = ''
    userForm.password = ''
    const defaultRole = roleOptions.value.find((r: any) => r.code === 'author')
    userForm.role_ids = defaultRole ? [defaultRole.id] : []
    userForm.real_name = ''
    userForm.department = ''
    userForm.remark = ''
    dialogVisible.value = true
}

const handleEditUser = (row: any) => {
    isEditing.value = true
    userForm.username = row.username
    userForm.password = ''
    userForm.role_ids = (row.roles || []).map((r: any) => r.id)
    userForm.real_name = row.label?.real_name || ''
    userForm.department = row.label?.department || ''
    userForm.remark = row.label?.remark || ''
    dialogVisible.value = true
}

const handleSubmitUser = () => {
    userFormRef.value?.validate(async (valid: boolean) => {
        if (!valid) return
        submitting.value = true
        try {
            const label: Record<string, string> = {}
            if (userForm.real_name.trim()) label.real_name = userForm.real_name.trim()
            if (userForm.department.trim()) label.department = userForm.department.trim()
            if (userForm.remark.trim()) label.remark = userForm.remark.trim()

            let res: any
            if (isEditing.value) {
                res = await reqUpdateUser({
                    username: userForm.username,
                    new_password: userForm.password || '',
                    role_ids: userForm.role_ids,
                    label
                })
            } else {
                res = await reqCreateUser({
                    username: userForm.username,
                    password: userForm.password,
                    role_ids: userForm.role_ids,
                    label
                })
            }

            if (res.code === 200) {
                ElMessage.success(isEditing.value ? '更新用户成功' : '创建用户成功')
                dialogVisible.value = false
                fetchUserList()
                if (isEditing.value && userForm.username === userStore.loginUser.name && !userForm.password) {
                    await userStore.fetchUserProfile()
                }
            } else {
                ElMessage.error(res.msg || '操作失败')
            }
        } finally {
            submitting.value = false
        }
    })
}

const handleDeleteUser = (row: any) => {
    ElMessageBox.confirm(`确定要永久删除用户 "${row.username}" 吗？`, '删除确认', {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(async () => {
        const res = await reqDeleteUser({ username: row.username })
        if (res.code === 200) {
            ElMessage.success('删除用户成功')
            fetchUserList()
        } else {
            ElMessage.error(res.msg || '删除用户失败')
        }
    }).catch(() => {})
}

onMounted(() => {
    fetchRoleOptions()
    fetchUserList()
})
</script>

<style scoped lang="scss">
.user-manage-container {
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

    .user-cell {
        display: flex;
        align-items: center;
        gap: 8px;

        .user-avatar-mini {
            width: 24px;
            height: 24px;
            border-radius: 6px;
            background: #eff6ff;
            color: #2563eb;
            font-size: 12px;
            font-weight: 700;
            display: inline-flex;
            align-items: center;
            justify-content: center;
            flex-shrink: 0;
        }
    }

    .label-tags {
        display: flex;
        flex-direction: column;
        gap: 2px;
        font-size: 12.5px;
        color: #475569;

        strong {
            font-weight: 600;
            color: #334155;
        }
    }

    .pagination-wrapper {
        margin-top: 16px;
        display: flex;
        justify-content: flex-end;
    }
}
</style>
