<script lang="js" setup>
import { ref, onMounted } from 'vue';
import zhCn from 'element-plus/es/locale/lang/zh-cn';
import { initTheme } from '@src/utils/theme';

const locale = ref(zhCn);
const currentTheme = ref('light');

onMounted(() => {
  // 初始化主题
  currentTheme.value = initTheme();
});
</script>

<template>
  <el-config-provider :locale="locale">
    <router-view></router-view>
  </el-config-provider>
</template>
