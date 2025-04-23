<template>
    <el-button size="small" icon="Refresh" circle @click="updateRefsh"></el-button>
    <el-button size="small" icon="FullScreen" circle @click="fullScreen" style="margin-right: 15px;"></el-button>

    <!-- 下拉菜单 -->
    <el-dropdown>
        <span class="el-dropdown-link">
            <span style="vertical-align: middle;">{{ userStore.loginUser.name }}</span>
            <el-icon class="el-icon--right" style="vertical-align: middle;">
                <arrow-down />
            </el-icon>
        </span>
        <template #dropdown>
            <el-dropdown-menu>
                <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
        </template>
    </el-dropdown>
</template>

<script lang="ts" setup>
import { useRouter, useRoute } from 'vue-router';
import useUserStore from '@/store/modules/user';
import useLayOutSettingStore from '@/store/modules/setting';
import { ElMessage } from 'element-plus';
// 获取骨架的小仓库
let LayOutSettingStore = useLayOutSettingStore()
// 获取用户相关的小仓库
let userStore = useUserStore()


const router = useRouter()

// 刷新按钮点击的回调
const updateRefsh = () => {
    LayOutSettingStore.refsh = !LayOutSettingStore.refsh
}

// 全屏按钮点击的回调
const fullScreen = () => {
    // DOM对象的一个属性，可以用来判断当前是不是全屏模式，如果是全屏，返回true，否则返回null
    let full = document.fullscreenElement;
    // 切换到全屏模式
    if (!full) {
        // 利用文档根节点方法requestFullscreen，实现全屏模式
        document.documentElement.requestFullscreen();
    } else {
        // 变为非全屏模式，退出全屏模式
        document.exitFullscreen();
    }
}

// 退出登录点击回调
// const logout = async () => {
//     // 1、调用后端接口（退出登录）
//     // 2、清空数据（token、username、等用户信息）
//     // 3、跳转到登录页面
//     await userStore.userLogout()
//     // 跳转到登录页面
//     $router.push('/login')  // 回首页
//     // $router.push({ path: '/login', query: { redirect: $route.path } })
// }


const handleLogout = async () => {
    console.log('退出登录')
    await userStore.userLogout()
    try {
        await userStore.userLogout()
        ElMessage.success('退出登录成功')
        // 跳转到登录页
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

<style scoped>
.el-dropdown-link {
    display: flex;
    align-items: center;
    cursor: pointer;
}

.el-icon--right {
    margin-left: 4px;
}
</style>