<template>
    <div class="login-container">
        <div class="login-box">
            <div class="login-title">期刊管理系统</div>
            <el-form ref="loginFormRef" :model="loginForm" :rules="loginRules" class="login-form">
                <el-form-item prop="username">
                    <el-input v-model="loginForm.username" placeholder="请输入用户名" :prefix-icon="User" />
                </el-form-item>
                <el-form-item prop="password">
                    <el-input v-model="loginForm.password" type="password" placeholder="请输入密码" :prefix-icon="Lock"
                        @keyup.enter="handleLogin" />
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" :loading="loading" class="login-button" @click="handleLogin">
                        登录
                    </el-button>
                </el-form-item>
            </el-form>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElNotification } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { getTime } from '@/utils/date'
import useUserStore from '@/store/modules/user'



const userStore = useUserStore()

// 路由实例
const router = useRouter()
const route = useRoute()

// 表单引用
const loginFormRef = ref()

// 加载状态
const loading = ref(false)

// 登录表单数据
const loginForm = reactive({
    username: '',
    password: ''
})

// 表单验证规则
const loginRules = {
    username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
        { min: 3, max: 20, message: '用户名长度应在 3 到 20 个字符之间', trigger: 'blur' }
    ],
    password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, max: 20, message: '密码长度应在 6 到 20 个字符之间', trigger: 'blur' }
    ]
}

// 登录处理
const handleLogin = () => {
    loginFormRef.value?.validate(async (valid: boolean) => {
        if (!valid) {
            return
        }

        try {
            loading.value = true
            // 这里调用登录 API
            // const result = await login(loginForm)

            // 模拟登录成功
            await userStore.userLogin(loginForm)
            loading.value = false
            // router.push({ path: '/' })
            // 获取重定向地址
            const redirect = route.query.redirect as string
            router.push(redirect || '/')
            // router.push('/')
            ElNotification({
                title: '登录成功',
                message: `${getTime()}好! 欢迎你`,
                type: 'success'
            })
        } catch (error: any) {
            ElNotification({
                title: '登录失败',
                message: error,
                type: 'error'
            })
            loading.value = false
        } finally {
            loading.value = false
        }
    })
}
</script>

<style scoped>
.login-container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
    background-color: #f5f5f5;
}

.login-box {
    width: 400px;
    padding: 40px;
    background-color: white;
    border-radius: 8px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.login-title {
    text-align: center;
    font-size: 24px;
    font-weight: bold;
    color: #333;
    margin-bottom: 30px;
}

.login-form {
    .el-form-item {
        margin-bottom: 25px;
    }

    .login-button {
        width: 100%;
    }
}

/* 输入框样式 */
:deep(.el-input__wrapper) {
    padding: 1px 11px;
}

:deep(.el-input__inner) {
    height: 40px;
}
</style>
