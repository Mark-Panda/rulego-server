<script lang="js" setup>
import MenuLayout from '@src/layout/components/default/menu.vue';
import { ref } from 'vue';
import useUserStore from '@src/store/module/user';
import { ElMessageBox } from 'element-plus';

const collapsed = ref(false);
const userStore = useUserStore();

const toggleCollapse = () => {
  collapsed.value = !collapsed.value;
};

const handleLogout = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要退出登录吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    );
    userStore.logout();
  } catch {
    // 用户取消操作
  }
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
          <el-dropdown @command="handleLogout">
            <el-button type="text" class="flex items-center">
              <el-avatar :size="32" class="mr-2">
                <el-icon><el-icon-user /></el-icon>
              </el-avatar>
              <span class="text-sm text-gray-600 dark:text-gray-300">admin</span>
              <el-icon class="ml-1"><el-icon-arrow-down /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">
                  <el-icon><el-icon-switch-button /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
      <div class="flex-grow overflow-auto p-4">
        <slot></slot>
      </div>
    </div>
  </div>
</template>
