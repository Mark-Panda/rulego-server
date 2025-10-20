<script lang="js" setup>
import { ref, onMounted, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { ElMessage } from 'element-plus';

const props = defineProps({
  collapsed: {
    type: Boolean,
    default: false
  }
});

const router = useRouter();
const activeIndex = ref('workflow-list');

function handleOpen() {}
function handleClose() {}
function handleSelect(val) {
  if (val === activeIndex.value) return;
  activeIndex.value = val;
  router.push({ name: val });
}



const menuItems = [
  {
    title: '工作台',
    items: [
      { 
        index: 'workflow-list', 
        icon: 'el-icon-menu', 
        label: '工作流管理',
        implemented: true
      },
      { 
        index: 'business-docs', 
        icon: 'el-icon-document', 
        label: '业务文档管理',
        implemented: true
      },
      { 
        index: 'runtime-logs', 
        icon: 'el-icon-video-play', 
        label: '运行日志',
        implemented: true
      },
      { 
        index: 'components-installed', 
        icon: 'el-icon-suitcase', 
        label: '组件管理',
        implemented: true,
        children: [
          {
            index: 'components-installed',
            label: '已安装组件',
            implemented: true
          },
          {
            index: 'components-market',
            label: '组件市场',
            implemented: true
          },
          {
            index: 'components-rules',
            label: '组件规则',
            implemented: true
          }
        ]
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
            
            <template v-for="item in section.items" :key="item.index">
              <!-- 有子菜单的项目 -->
              <el-sub-menu v-if="item.children && item.children.length > 0" :index="item.index">
                <template #title>
                  <el-icon><component :is="item.icon" /></el-icon>
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
                
                <el-menu-item 
                  v-for="child in item.children" 
                  :key="child.index" 
                  :index="child.index"
                  @click="handleItemClick(child)"
                >
                  <template #title>
                    <span>{{ child.label }}</span>
                    <el-tag 
                      v-if="!child.implemented" 
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
              
              <!-- 没有子菜单的项目 -->
              <el-menu-item 
                v-else
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
            </template>
          </el-sub-menu>
        </template>
      </el-menu>
    </div>
    
    <!-- Bottom Actions -->
    <!-- 主题切换功能已禁用 -->
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
