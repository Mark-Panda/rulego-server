<script setup>
import NodesGraph from '@src/assets/nodes-graph.svg'
import SplitScreen from '@src/assets/split-screen.svg'

const props = defineProps({
  modelValue: {
    type: Object,
    default: {
      isFlowVisible: true,
      isSourceCodeVisible: false,
    },
  },
})
const emit = defineEmits(['update:modelValue', 'optimize-layout'])

/**
 * 设置视图模式
 * @param {string} mode - 'flow', 'code', 'split'
 */
function setViewMode(mode) {
  const newValue = { ...props.modelValue }
  
  switch (mode) {
    case 'flow':
      newValue.isFlowVisible = true
      newValue.isSourceCodeVisible = false
      break
    case 'code':
      newValue.isFlowVisible = false
      newValue.isSourceCodeVisible = true
      break
    case 'split':
      newValue.isFlowVisible = true
      newValue.isSourceCodeVisible = true
      break
  }
  
  emit('update:modelValue', newValue)
}

/**
 * 优化布局
 */
function optimizeLayout() {
  emit('optimize-layout')
}
</script>
<template>
  <div class="view-toggle-container">
    <div class="toggle-btn-wrapper">
      <el-tooltip content="流程图视图" placement="top" :show-after="500">
        <el-button 
          :type="modelValue.isFlowVisible && !modelValue.isSourceCodeVisible ? 'primary' : ''"
          :class="{ 'is-active': modelValue.isFlowVisible && !modelValue.isSourceCodeVisible }"
          size="small"
          @click="setViewMode('flow')"
          class="toggle-btn"
          circle
        >
          <el-icon :size="16">
            <NodesGraph />
          </el-icon>
        </el-button>
      </el-tooltip>
    </div>
    
    <div class="toggle-btn-wrapper">
      <el-tooltip content="源码视图" placement="top" :show-after="500">
        <el-button 
          :type="!modelValue.isFlowVisible && modelValue.isSourceCodeVisible ? 'primary' : ''"
          :class="{ 'is-active': !modelValue.isFlowVisible && modelValue.isSourceCodeVisible }"
          size="small"
          @click="setViewMode('code')"
          class="toggle-btn"
          circle
        >
          <el-icon :size="16">
            <SplitScreen />
          </el-icon>
        </el-button>
      </el-tooltip>
    </div>
    
    <div class="toggle-btn-wrapper">
      <el-tooltip content="分屏视图" placement="top" :show-after="500">
        <el-button 
          :type="modelValue.isFlowVisible && modelValue.isSourceCodeVisible ? 'primary' : ''"
          :class="{ 'is-active': modelValue.isFlowVisible && modelValue.isSourceCodeVisible }"
          size="small"
          @click="setViewMode('split')"
          class="toggle-btn"
          circle
        >
          <el-icon :size="16">
            <el-icon-grid />
          </el-icon>
        </el-button>
      </el-tooltip>
    </div>
    
    <!-- 分隔线 -->
    <div class="divider"></div>
    
    <div class="toggle-btn-wrapper">
      <el-tooltip content="优化布局" placement="top" :show-after="500">
        <el-button 
          size="small"
          @click="optimizeLayout"
          class="toggle-btn"
          circle
        >
          <el-icon :size="16">
            <el-icon-magic-stick />
          </el-icon>
        </el-button>
      </el-tooltip>
    </div>
  </div>
</template>

<style scoped>
.view-toggle-container {
  display: flex;
  flex-direction: row;
  gap: 6px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  padding: 6px 10px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  min-height: 48px;
}

.toggle-btn-wrapper {
  display: flex;
  justify-content: center;
}

.divider {
  width: 1px;
  height: 24px;
  background: var(--el-border-color-lighter);
  margin: 0 4px;
  opacity: 0.6;
}

/* 基础按钮样式重置 */
.toggle-btn.el-button {
  width: 36px !important;
  height: 36px !important;
  padding: 0 !important;
  border-radius: 50% !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
  border: 1px solid var(--el-border-color-lighter) !important;
  position: relative !important;
  overflow: hidden !important;
}

/* 默认状态 - 无背景色 */
.view-toggle-container .toggle-btn.el-button:not(.el-button--primary):not(.is-active) {
  background: transparent !important;
  color: var(--el-text-color-regular) !important;
  border-color: var(--el-border-color-lighter) !important;
  box-shadow: none !important;
  transform: none !important;
}

.view-toggle-container .toggle-btn.el-button:not(.el-button--primary):not(.is-active):hover {
  background: var(--el-color-primary-light-9) !important;
  color: var(--el-color-primary) !important;
  border-color: var(--el-color-primary-light-7) !important;
  transform: translateY(-3px) scale(1.05) !important;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.2) !important;
}

/* 激活状态 - 蓝色背景 */
.view-toggle-container .toggle-btn.el-button.el-button--primary,
.view-toggle-container .toggle-btn.el-button.is-active {
  background: var(--el-color-primary) !important;
  border-color: var(--el-color-primary) !important;
  color: white !important;
  transform: translateY(-2px) !important;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.4) !important;
  position: relative !important;
}

.view-toggle-container .toggle-btn.el-button.el-button--primary:hover,
.view-toggle-container .toggle-btn.el-button.is-active:hover {
  background: var(--el-color-primary-light-3) !important;
  border-color: var(--el-color-primary-light-3) !important;
  transform: translateY(-3px) scale(1.05) !important;
  box-shadow: 0 4px 16px rgba(64, 158, 255, 0.5) !important;
}

/* 激活状态指示器 - 只对真正激活的按钮显示 */
.toggle-btn.is-active::before {
  content: '';
  position: absolute;
  top: -2px;
  left: -2px;
  right: -2px;
  bottom: -2px;
  border: 2px solid var(--el-color-primary);
  border-radius: 50%;
  opacity: 0.6;
  animation: activeIndicator 2s ease-in-out infinite;
}

/* 移除非激活按钮的伪元素 */
.toggle-btn:not(.is-active)::before {
  display: none !important;
}

@keyframes activeIndicator {
  0%, 100% {
    transform: scale(1);
    opacity: 0.6;
  }
  50% {
    transform: scale(1.1);
    opacity: 0.3;
  }
}

/* 深色模式适配 */
.dark .view-toggle-container {
  background: rgba(31, 41, 55, 0.95);
  border-color: rgba(75, 85, 99, 0.3);
}

.dark .view-toggle-container .toggle-btn.el-button:not(.el-button--primary):not(.is-active) {
  background: transparent !important;
  color: var(--el-text-color-regular) !important;
  border-color: var(--el-border-color) !important;
  box-shadow: none !important;
  transform: none !important;
}

.dark .view-toggle-container .toggle-btn.el-button:not(.el-button--primary):not(.is-active):hover {
  background: rgba(64, 158, 255, 0.1) !important;
  border-color: var(--el-color-primary-light-7) !important;
  color: var(--el-color-primary) !important;
  transform: translateX(3px) scale(1.05) !important;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.2) !important;
}

/* 确保激活状态在深色模式下也清晰可见 */
.dark .view-toggle-container .toggle-btn.el-button.el-button--primary,
.dark .view-toggle-container .toggle-btn.el-button.is-active {
  background: var(--el-color-primary) !important;
  border-color: var(--el-color-primary) !important;
  color: white !important;
  transform: translateX(2px) !important;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.4) !important;
}

.dark .view-toggle-container .toggle-btn.el-button.el-button--primary:hover,
.dark .view-toggle-container .toggle-btn.el-button.is-active:hover {
  background: var(--el-color-primary-light-3) !important;
  border-color: var(--el-color-primary-light-3) !important;
  transform: translateX(3px) scale(1.05) !important;
  box-shadow: 0 4px 16px rgba(64, 158, 255, 0.5) !important;
}
</style>
