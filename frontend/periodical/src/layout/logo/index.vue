<template>
    <div class="logo-container" :class="{ 'is-fold': layOutSettingStore.fold }" v-if="setting.logoHidden">
        <div class="logo-badge">
            <img :src="setting.logo" alt="logo" />
        </div>
        <div class="logo-text" v-show="!layOutSettingStore.fold">
            <p class="brand-title">{{ setting.title }}</p>
            <span class="brand-sub">EDITORIAL &amp; PUBLISHING</span>
        </div>
    </div>
</template>

<script setup lang="ts">
import setting from '@/setting'
import useLayOutSettingStore from '@/store/modules/setting'

const layOutSettingStore = useLayOutSettingStore()
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
    color: #ffffff;
    display: flex;
    align-items: center;
    padding: 0 18px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.07);
    overflow: hidden;
    transition: padding 0.25s cubic-bezier(0.22, 1, 0.36, 1);

    &.is-fold {
        padding: 0;
        justify-content: center;
    }

    .logo-badge {
        width: 36px;
        height: 36px;
        border-radius: 9px;
        background: linear-gradient(135deg, #2563eb 0%, #0284c7 100%);
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
        box-shadow: 0 2px 8px rgba(37, 99, 235, 0.35);

        img {
            width: 20px;
            height: 20px;
            display: block;
            filter: brightness(1.15);
        }
    }

    .logo-text {
        margin-left: 10px;
        display: flex;
        flex-direction: column;
        white-space: nowrap;

        .brand-title {
            font-size: 14.5px;
            font-weight: 600;
            letter-spacing: 0.01em;
            color: #f8fafc;
            line-height: 1.2;
        }

        .brand-sub {
            font-size: 9.5px;
            font-weight: 500;
            letter-spacing: 0.04em;
            color: #64748b;
            margin-top: 2px;
        }
    }
}
</style>
