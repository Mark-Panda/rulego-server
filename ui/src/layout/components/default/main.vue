<script lang="js" setup>
import MenuLayout from '@src/layout/components/default/menu.vue';
import { ref } from 'vue';

const collapsed = ref(false);
const toggleCollapse = () => {
  collapsed.value = !collapsed.value;
};
</script>

<template>
  <div class="flex h-screen w-screen overflow-hidden bg-gray-50 dark:bg-gray-900">
    <div :class="[
      'transition-all duration-300 flex-none border-r border-[var(--el-border-color)] bg-white dark:bg-gray-800',
      collapsed ? 'w-16' : 'w-56'
    ]">
      <menu-layout :collapsed="collapsed"></menu-layout>
    </div>
    <div class="flex flex-grow flex-col overflow-hidden">
      <div class="flex h-14 flex-none items-center justify-between border-b border-[var(--el-border-color)] bg-white px-4 dark:bg-gray-800">
        <div class="flex items-center">
          <el-button type="text" @click="toggleCollapse">
            <el-icon>
              <el-icon-fold v-if="!collapsed" />
              <el-icon-expand v-else />
            </el-icon>
          </el-button>
          <div class="ml-4 text-lg font-medium text-gray-800 dark:text-gray-200">RuleGo 工作台</div>
        </div>
        <div class="flex items-center space-x-2">
          <el-tooltip content="帮助文档" placement="bottom">
            <el-button type="text">
              <el-icon><el-icon-question-filled /></el-icon>
            </el-button>
          </el-tooltip>
          <el-tooltip content="通知" placement="bottom">
            <el-button type="text">
              <el-icon><el-icon-bell /></el-icon>
            </el-button>
          </el-tooltip>
          <el-tooltip content="设置" placement="bottom">
            <el-button type="text">
              <el-icon><el-icon-setting /></el-icon>
            </el-button>
          </el-tooltip>
        </div>
      </div>
      <div class="flex-grow overflow-auto p-4">
        <slot></slot>
      </div>
    </div>
  </div>
</template>
