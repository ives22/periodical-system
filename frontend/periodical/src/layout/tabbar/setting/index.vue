<template>
    <div class="setting-actions">
        <el-tooltip content="刷新当前视图" placement="bottom">
            <button class="icon-action-btn" type="button" @click="updateRefsh">
                <el-icon :size="16"><Refresh /></el-icon>
            </button>
        </el-tooltip>
        <el-tooltip content="切换全屏显示" placement="bottom">
            <button class="icon-action-btn" type="button" @click="fullScreen">
                <el-icon :size="16"><FullScreen /></el-icon>
            </button>
        </el-tooltip>

        <span class="action-divider"></span>

        <!-- 用户身份胶囊下拉菜单 -->
        <el-dropdown trigger="click">
            <div class="user-profile-pill">
                <div class="avatar-circle">
                    {{ (userStore.loginUser.name || 'U').slice(0, 1).toUpperCase() }}
                </div>
                <div class="user-meta">
                    <span class="user-name">{{ userStore.loginUser.name }}</span>
                    <span v-if="userStore.roleNames" class="user-role">{{ userStore.roleNames }}</span>
                </div>
                <el-icon class="dropdown-caret">
                    <arrow-down />
                </el-icon>
            </div>
            <template #dropdown>
                <el-dropdown-menu>
                    <el-dropdown-item icon="Key" @click="openPwdDialog">修改个人密码</el-dropdown-item>
                    <el-dropdown-item divided icon="SwitchButton" @click="handleLogout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
            </template>
        </el-dropdown>

        <!-- 修改个人密码弹窗 -->
        <el-dialog v-model="pwdDialogVisible" title="修改个人密码" width="440px" append-to-body>
            <el-form ref="pwdFormRef" :model="pwdForm" :rules="pwdRules" label-width="95px">
                <el-form-item label="当前账号">
                    <el-input :value="userStore.loginUser.name" disabled />
                </el-form-item>
                <el-form-item label="旧密码" prop="old_password">
                    <el-input v-model="pwdForm.old_password" type="password" show-password placeholder="请输入旧密码" />
                </el-form-item>
                <el-form-item label="新密码" prop="new_password">
                    <el-input v-model="pwdForm.new_password" type="password" show-password placeholder="请输入新密码（6-20位）" />
                </el-form-item>
                <el-form-item label="确认新密码" prop="confirm_password">
                    <el-input v-model="pwdForm.confirm_password" type="password" show-password placeholder="请再次输入新密码" />
                </el-form-item>
            </el-form>
            <template #footer>
                <el-button @click="pwdDialogVisible = false">取消</el-button>
                <el-button type="primary" :loading="pwdSubmitting" @click="submitUpdatePwd">确定修改</el-button>
            </template>
        </el-dialog>
    </div>
</template>

<script lang="ts" setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import useUserStore from '@/store/modules/user'
import useLayOutSettingStore from '@/store/modules/setting'
import { reqUpdatePassword } from '@/api/user'
import { ElMessage } from 'element-plus'

let LayOutSettingStore = useLayOutSettingStore()
let userStore = useUserStore()
const router = useRouter()

const updateRefsh = () => {
    LayOutSettingStore.refsh = !LayOutSettingStore.refsh
}

const fullScreen = () => {
    let full = document.fullscreenElement
    if (!full) {
        document.documentElement.requestFullscreen()
    } else {
        document.exitFullscreen()
    }
}

const pwdDialogVisible = ref(false)
const pwdSubmitting = ref(false)
const pwdFormRef = ref()
const pwdForm = reactive({
    old_password: '',
    new_password: '',
    confirm_password: ''
})

const pwdRules = {
    old_password: [
        { required: true, message: '请输入旧密码', trigger: 'blur' }
    ],
    new_password: [
        { required: true, message: '请输入新密码', trigger: 'blur' },
        { min: 6, max: 20, message: '密码长度需在 6 到 20 个字符之间', trigger: 'blur' }
    ],
    confirm_password: [
        { required: true, message: '请再次输入新密码', trigger: 'blur' },
        {
            validator: (_rule: any, value: string, callback: any) => {
                if (value !== pwdForm.new_password) {
                    callback(new Error('两次输入的新密码不一致'))
                } else {
                    callback()
                }
            },
            trigger: 'blur'
        }
    ]
}

const openPwdDialog = () => {
    pwdForm.old_password = ''
    pwdForm.new_password = ''
    pwdForm.confirm_password = ''
    pwdDialogVisible.value = true
}

const submitUpdatePwd = () => {
    pwdFormRef.value?.validate(async (valid: boolean) => {
        if (!valid) return
        pwdSubmitting.value = true
        try {
            const res = await reqUpdatePassword({
                username: userStore.loginUser.name,
                old_password: pwdForm.old_password,
                new_password: pwdForm.new_password
            })
            if (res.code === 200) {
                ElMessage.success('密码修改成功，请重新登录')
                pwdDialogVisible.value = false
                userStore.clearUserInfo()
                router.push('/login')
            } else {
                ElMessage.error(res.msg || '密码修改失败')
            }
        } catch (err) {
            ElMessage.error('密码修改请求失败')
        } finally {
            pwdSubmitting.value = false
        }
    })
}

const handleLogout = async () => {
    try {
        await userStore.userLogout()
        ElMessage.success('退出登录成功')
        router.push('/login')
    } catch (error) {
        console.error('退出登录失败：', error)
        ElMessage.error('退出登录失败')
    }
}
</script>

<script lang="ts">
export default {
    name: 'Setting'
}
</script>

<style scoped lang="scss">
.setting-actions {
    display: flex;
    align-items: center;
    gap: 8px;

    .icon-action-btn {
        width: 34px;
        height: 34px;
        border-radius: 8px;
        border: 1px solid #e2e8f0;
        background: #ffffff;
        color: #475569;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        transition: all 0.2s ease;

        &:hover {
            color: #0f172a;
            background: #f8fafc;
            border-color: #cbd5e1;
        }
    }

    .action-divider {
        width: 1px;
        height: 20px;
        background: #e2e8f0;
        margin: 0 4px;
    }

    .user-profile-pill {
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 4px 10px 4px 6px;
        border-radius: 999px;
        border: 1px solid #e2e8f0;
        background: #ffffff;
        cursor: pointer;
        transition: all 0.2s ease;

        &:hover {
            background: #f8fafc;
            border-color: #cbd5e1;
        }

        .avatar-circle {
            width: 28px;
            height: 28px;
            border-radius: 50%;
            background: linear-gradient(135deg, #0f172a 0%, #2563eb 100%);
            color: #ffffff;
            font-size: 13px;
            font-weight: 600;
            display: flex;
            align-items: center;
            justify-content: center;
        }

        .user-meta {
            display: flex;
            flex-direction: column;
            line-height: 1.15;

            .user-name {
                font-size: 13px;
                font-weight: 600;
                color: #0f172a;
            }

            .user-role {
                font-size: 11px;
                color: #64748b;
            }
        }

        .dropdown-caret {
            color: #64748b;
            font-size: 12px;
        }
    }
}
</style>
