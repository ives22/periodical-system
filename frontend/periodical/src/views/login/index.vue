<template>
    <div class="login-page">
        <div class="login-shell">
            <!-- 左侧：品牌与专业学术视觉区 -->
            <div class="brand-panel">
                <div class="brand-header">
                    <div class="brand-badge">P</div>
                    <span class="brand-name">PERIODICAL SYSTEM</span>
                </div>

                <div class="brand-hero">
                    <p class="eyebrow">ACADEMIC &amp; PUBLISHING OPERATIONS</p>
                    <h1>学术期刊数字化<br />综合管理平台</h1>
                    <p class="hero-desc">
                        集多维收稿方向检索、复合影响因子分析与刊期批次追踪于一体的专业出版工作台。
                    </p>
                </div>

                <div class="brand-footer">
                    <span>© {{ new Date().getFullYear() }} Periodical Management System · v2.0</span>
                </div>
            </div>

            <!-- 右侧：极净登录表单区 -->
            <div class="form-panel">
                <div class="form-box">
                    <div class="form-header">
                        <h2>欢迎登录</h2>
                        <p>请输入您的系统账号与密码以访问工作台</p>
                    </div>

                    <el-form
                        ref="loginFormRef"
                        :model="loginForm"
                        :rules="loginRules"
                        label-position="top"
                        class="login-form"
                    >
                        <el-form-item label="登录账号" prop="username">
                            <el-input
                                v-model="loginForm.username"
                                placeholder="请输入用户名"
                                :prefix-icon="User"
                                size="large"
                                autocomplete="username"
                            />
                        </el-form-item>

                        <el-form-item label="登录密码" prop="password">
                            <el-input
                                v-model="loginForm.password"
                                type="password"
                                show-password
                                placeholder="请输入登录密码"
                                :prefix-icon="Lock"
                                size="large"
                                autocomplete="current-password"
                                @keyup.enter="handleLogin"
                            />
                        </el-form-item>

                        <el-form-item class="submit-item">
                            <el-button
                                type="primary"
                                size="large"
                                :loading="loading"
                                class="login-button"
                                @click="handleLogin"
                            >
                                进入工作台
                            </el-button>
                        </el-form-item>
                    </el-form>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElNotification } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { getTime } from '@/utils/date'
import useUserStore from '@/store/modules/user'

const userStore = useUserStore()
const router = useRouter()
const route = useRoute()

const loginFormRef = ref()
const loading = ref(false)

const loginForm = reactive({
    username: '',
    password: ''
})

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

const handleLogin = () => {
    loginFormRef.value?.validate(async (valid: boolean) => {
        if (!valid) {
            return
        }

        try {
            loading.value = true
            await userStore.userLogin(loginForm)
            const redirect = route.query.redirect as string
            router.push(redirect || '/')
            ElNotification({
                title: '登录成功',
                message: `${getTime()}好，${loginForm.username}！欢迎回到期刊管理系统`,
                type: 'success'
            })
        } catch (error: any) {
            ElNotification({
                title: '登录失败',
                message: typeof error === 'string' ? error : '用户名或密码错误',
                type: 'error'
            })
        } finally {
            loading.value = false
        }
    })
}
</script>

<style scoped lang="scss">
.login-page {
    min-height: 100vh;
    width: 100%;
    background-color: #f1f5f9;
    background-image: radial-gradient(#cbd5e1 1px, transparent 1px);
    background-size: 24px 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 24px;
}

.login-shell {
    width: 100%;
    max-width: 920px;
    min-height: 480px;
    background: #ffffff;
    border-radius: 16px;
    border: 1px solid #e2e8f0;
    box-shadow: 0 25px 50px -12px rgba(15, 23, 42, 0.12);
    display: grid;
    grid-template-columns: 1.05fr 1fr;
    overflow: hidden;
}

.brand-panel {
    background: linear-gradient(155deg, #0f172a 0%, #1e293b 100%);
    color: #f8fafc;
    padding: 44px 42px;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    position: relative;
    overflow: hidden;

    &::after {
        content: '';
        position: absolute;
        right: -80px;
        bottom: -80px;
        width: 260px;
        height: 260px;
        border-radius: 50%;
        background: radial-gradient(circle, rgba(37, 99, 235, 0.28) 0%, transparent 70%);
        pointer-events: none;
    }

    .brand-header {
        display: flex;
        align-items: center;
        gap: 12px;

        .brand-badge {
            width: 34px;
            height: 34px;
            border-radius: 8px;
            background: #2563eb;
            color: #ffffff;
            font-weight: 700;
            font-size: 16px;
            display: flex;
            align-items: center;
            justify-content: center;
            box-shadow: 0 4px 12px rgba(37, 99, 235, 0.4);
        }

        .brand-name {
            font-size: 13px;
            font-weight: 600;
            letter-spacing: 0.12em;
            color: #cbd5e1;
        }
    }

    .brand-hero {
        margin: 28px 0;

        .eyebrow {
            font-size: 11px;
            font-weight: 600;
            letter-spacing: 0.14em;
            color: #38bdf8;
            margin-bottom: 12px;
        }

        h1 {
            font-size: 32px;
            font-weight: 700;
            line-height: 1.28;
            color: #ffffff;
            margin-bottom: 16px;
            letter-spacing: -0.01em;
        }

        .hero-desc {
            font-size: 14px;
            line-height: 1.7;
            color: #94a3b8;
        }
    }

    .feature-grid {
        display: flex;
        flex-direction: column;
        gap: 16px;
        border-top: 1px solid rgba(255, 255, 255, 0.08);
        padding-top: 24px;

        .feature-item {
            display: flex;
            align-items: flex-start;
            gap: 14px;

            .feature-index {
                font-size: 12px;
                font-weight: 700;
                color: #38bdf8;
                background: rgba(56, 189, 248, 0.12);
                padding: 4px 8px;
                border-radius: 6px;
                font-family: monospace;
            }

            h4 {
                font-size: 14px;
                font-weight: 600;
                color: #f1f5f9;
                margin-bottom: 4px;
            }

            p {
                font-size: 12.5px;
                color: #94a3b8;
                line-height: 1.5;
            }
        }
    }

    .brand-footer {
        font-size: 12px;
        color: #64748b;
        margin-top: 20px;
    }
}

.form-panel {
    padding: 48px 44px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #ffffff;

    .form-box {
        width: 100%;
        max-width: 360px;

        .form-header {
            margin-bottom: 32px;

            h2 {
                font-size: 24px;
                font-weight: 700;
                color: #0f172a;
                margin-bottom: 8px;
            }

            p {
                font-size: 13.5px;
                color: #64748b;
            }
        }

        .login-form {
            :deep(.el-form-item__label) {
                font-size: 13px;
                font-weight: 600;
                color: #334155;
                padding-bottom: 6px;
            }

            .submit-item {
                margin-top: 28px;
            }

            .login-button {
                width: 100%;
                height: 44px;
                font-size: 15px;
                font-weight: 600;
                letter-spacing: 0.02em;
                border-radius: 8px;
                background: #0f172a;
                border-color: #0f172a;

                &:hover {
                    background: #1e293b;
                    border-color: #1e293b;
                }
            }
        }

        .form-footnote {
            margin-top: 24px;
            padding-top: 18px;
            border-top: 1px solid #f1f5f9;
            display: flex;
            align-items: center;
            gap: 8px;
            font-size: 12px;
            color: #64748b;

            .security-dot {
                width: 7px;
                height: 7px;
                border-radius: 50%;
                background: #10b981;
                flex-shrink: 0;
            }
        }
    }
}

@media (max-width: 820px) {
    .login-shell {
        grid-template-columns: 1fr;
        max-width: 440px;
    }

    .brand-panel {
        display: none;
    }
}
</style>
