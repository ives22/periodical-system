<template>
    <div class="home-container">
        <div class="welcome-section">
            <h1>欢迎使用期刊管理系统</h1>
            <div class="time-display">
                <p>{{ currentTime }}</p>
                <p>{{ welcomeMessage }}</p>
            </div>
        </div>

        <el-row :gutter="20" class="info-cards">
            <el-col :span="8">
                <el-card shadow="hover">
                    <template #header>
                        <div class="card-header">
                            <span>系统信息</span>
                        </div>
                    </template>
                    <div class="card-content">
                        <p>当前用户：{{ userStore.loginUser.name }}</p>
                        <!-- <p>用户角色：{{ userStore.loginUser.role === 1 ? '管理员' : '普通用户' }}</p> -->
                    </div>
                </el-card>
            </el-col>

            <el-col :span="8">
                <el-card shadow="hover">
                    <template #header>
                        <div class="card-header">
                            <span>快捷导航</span>
                        </div>
                    </template>
                    <div class="card-content">
                        <el-button type="primary" @click="router.push('/periodical')">期刊管理</el-button>
                    </div>
                </el-card>
            </el-col>

            <el-col :span="8">
                <el-card shadow="hover">
                    <template #header>
                        <div class="card-header">
                            <span>系统公告</span>
                        </div>
                    </template>
                    <div class="card-content">
                        <p>欢迎使用期刊管理系统！</p>
                    </div>
                </el-card>
            </el-col>
        </el-row>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import useUserStore from '@/store/modules/user'

const router = useRouter()
const userStore = useUserStore()
const currentTime = ref('')
const welcomeMessage = ref('')

// 获取当前时间和欢迎语
const updateTime = () => {
    const now = new Date()
    const hours = now.getHours()
    currentTime.value = now.toLocaleString()

    if (hours < 12) {
        welcomeMessage.value = '早上好！'
    } else if (hours < 18) {
        welcomeMessage.value = '下午好！'
    } else {
        welcomeMessage.value = '晚上好！'
    }
}

onMounted(() => {
    updateTime()
    // 每秒更新时间
    setInterval(updateTime, 1000)
})
</script>

<style scoped lang="scss">
.home-container {
    padding: 20px;

    .welcome-section {
        text-align: center;
        margin-bottom: 40px;

        h1 {
            font-size: 28px;
            color: #303133;
            margin-bottom: 20px;
        }

        .time-display {
            font-size: 16px;
            color: #606266;

            p {
                margin: 10px 0;
            }
        }
    }

    .info-cards {
        .el-card {
            height: 200px;

            .card-header {
                font-weight: bold;
                color: #303133;
            }

            .card-content {
                padding: 10px 0;

                p {
                    margin: 10px 0;
                    color: #606266;
                }

                .el-button {
                    margin: 10px 0;
                }
            }
        }
    }
}
</style>
