<script lang="js" setup>
import { ref, onMounted, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { ElMessage } from 'element-plus';
import { getThemeMode, toggleThemeMode } from '@src/utils/theme';

const props = defineProps({
  collapsed: {
    type: Boolean,
    default: false
  }
});

const router = useRouter();
const activeIndex = ref('workflow-list');
const themeMode = ref(getThemeMode());

function handleOpen() {}
function handleClose() {}
function handleSelect(val) {
  if (val === activeIndex.value) return;
  activeIndex.value = val;
  router.push({ name: val });
}

// 切换主题
const switchTheme = () => {
  const newMode = toggleThemeMode();
  themeMode.value = newMode;
  ElMessage.success(`已切换到${newMode === 'dark' ? '暗黑' : '明亮'}模式`);
};

const menuItems = [
  {
    title: '工作台',
    items: [
      { 
        index: 'workflow-list', 
        icon: 'el-icon-menu', 
        label: '我的应用',
        implemented: true
      },
      { 
        index: 'runs', 
        icon: 'el-icon-video-play', 
        label: '运行日志',
        implemented: false
      }
    ]
  },
  {
    title: '应用市场',
    items: [
      { 
        index: 'app-market', 
        icon: 'el-icon-box', 
        label: '应用市场',
        implemented: false
      }
    ]
  },
  {
    title: '管理',
    items: [
      { 
        index: 'share-node-list', 
        icon: 'el-icon-connection', 
        label: '授权配置',
        implemented: true
      },
      { 
        index: 'node-manage', 
        icon: 'el-icon-suitcase', 
        label: '组件管理',
        implemented: false
      },
      { 
        index: 'system-manage', 
        icon: 'el-icon-setting', 
        label: '系统管理',
        implemented: false
      }
    ]
  }
];

onMounted(() => {
  const { name } = useRoute();
  activeIndex.value = name;
});

const handleItemClick = (item) => {
  if (!item.implemented) {
    ElMessage.error('功能暂未实现');
    return;
  }
  handleSelect(item.index);
};

const logoStyle = computed(() => {
  return props.collapsed 
    ? 'justify-center text-xl py-4' 
    : 'justify-start text-xl py-4 px-4';
});
</script>

<template>
  <div class="flex h-full flex-col overflow-hidden">
    <!-- Logo -->
    <div :class="['flex items-center flex-none border-b border-[var(--el-border-color)]', logoStyle]">
      <div class="flex items-center">
        <el-icon class="text-primary text-xl"><el-icon-connection /></el-icon>
        <span v-if="!collapsed" class="ml-2 font-bold text-primary">RuleGo</span>
      </div>
    </div>
    
    <!-- Menu -->
    <div class="flex-grow overflow-auto">
      <el-menu
        :default-active="activeIndex"
        class="default-layout-menu h-full"
        :collapse="collapsed"
        @open="handleOpen"
        @close="handleClose"
        @select="handleSelect"
      >
        <template v-for="(section, sectionIndex) in menuItems" :key="sectionIndex">
          <el-sub-menu :index="'section-' + sectionIndex">
            <template #title>
              <span>{{ section.title }}</span>
            </template>
            
            <el-menu-item 
              v-for="item in section.items" 
              :key="item.index" 
              :index="item.index"
              @click="handleItemClick(item)"
            >
              <el-icon><component :is="item.icon" /></el-icon>
              <template #title>
                <span>{{ item.label }}</span>
                <el-tag 
                  v-if="!item.implemented" 
                  size="small" 
                  type="info" 
                  effect="plain"
                  class="ml-2"
                >
                  开发中
                </el-tag>
              </template>
            </el-menu-item>
          </el-sub-menu>
        </template>
      </el-menu>
    </div>
    
    <!-- Bottom Actions -->
    <div class="flex-none border-t border-[var(--el-border-color)] p-2">
      <el-tooltip :content="collapsed ? '切换主题' : ''" placement="right">
        <el-button type="text" class="w-full justify-center" @click="switchTheme">
          <el-icon>
            <el-icon-moon v-if="themeMode === 'light'" />
            <el-icon-sunny v-else />
          </el-icon>
          <span v-if="!collapsed" class="ml-2">
            {{ themeMode === 'light' ? '切换到暗黑模式' : '切换到明亮模式' }}
          </span>
        </el-button>
      </el-tooltip>
    </div>
  </div>
</template>

<style scoped>
.default-layout-menu {
  border-right: none;
}

.default-layout-menu :deep(.el-sub-menu__title) {
  font-weight: 500;
}

.default-layout-menu :deep(.el-menu-item.is-active) {
  background-color: rgba(64, 158, 255, 0.1);
}

/* 暗黑模式下的菜单样式 */
:deep(.dark .el-menu) {
  --el-menu-bg-color: var(--el-bg-color);
  --el-menu-text-color: var(--el-text-color-primary);
  --el-menu-hover-bg-color: var(--el-fill-color-light);
  --el-menu-hover-text-color: var(--el-color-primary);
  --el-menu-active-color: var(--el-color-primary);
}

:deep(.dark .el-menu-item.is-active) {
  background-color: rgba(64, 158, 255, 0.2);
}
</style>
