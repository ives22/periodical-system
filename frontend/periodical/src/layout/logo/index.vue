<template>
    <div class="logo-container" v-if="setting.logoHidden">
        <img :src="setting.logo" alt="logo">
        <p>{{ setting.title }}</p>
    </div>
</template>

<script setup lang="ts">
// 引入配置
import setting from '@/setting';
</script>

<script lang="ts">
export default {
    name: 'Logo'
}
</script>

<style scoped lang="scss">
.logo-container {
    width: 100%;
    height: $baseMenuLogoHeight;
    color: white;
    display: flex;
    align-items: center;
    padding: 20px;

    img {
        width: 25px;
        height: 25px;
    }

    p {
        font-size: $baseMenuLogoFontSize;
        margin-left: 10px;
    }
}
</style>