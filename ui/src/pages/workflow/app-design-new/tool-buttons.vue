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
  isMiniMapVisible: {
    type: Boolean,
    default: true,
  },
  isAiAssistantVisible: {
    type: Boolean,
    default: false,
  },
})
const emit = defineEmits(['update:modelValue', 'optimize-layout', 'toggle-minimap', 'toggle-ai-assistant'])

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

/**
 * 切换小地图显示状态
 */
function toggleMiniMap() {
  console.log('工具按钮：toggleMiniMap 被点击');
  emit('toggle-minimap');
}

/**
 * 切换AI助手显示状态
 */
function toggleAiAssistant() {
  console.log('工具按钮：toggleAiAssistant 被点击');
  emit('toggle-ai-assistant');
}
</script>
<template>
  <div class="view-toggle-container">
    <!-- 视图切换按钮组 -->
    <div class="button-group">
      <div class="toggle-btn-wrapper">
        <el-tooltip content="流程图视图" placement="top" :show-after="300">
          <el-button 
            :type="modelValue.isFlowVisible && !modelValue.isSourceCodeVisible ? 'primary' : ''"
            :class="{ 'is-active': modelValue.isFlowVisible && !modelValue.isSourceCodeVisible }"
            size="small"
            @click="setViewMode('flow')"
            class="toggle-btn view-btn"
            circle
          >
            <el-icon :size="16">
              <NodesGraph />
            </el-icon>
          </el-button>
        </el-tooltip>
      </div>
      
      <div class="toggle-btn-wrapper">
        <el-tooltip content="源码视图" placement="top" :show-after="300">
          <el-button 
            :type="!modelValue.isFlowVisible && modelValue.isSourceCodeVisible ? 'primary' : ''"
            :class="{ 'is-active': !modelValue.isFlowVisible && modelValue.isSourceCodeVisible }"
            size="small"
            @click="setViewMode('code')"
            class="toggle-btn view-btn"
            circle
          >
            <el-icon :size="16">
              <SplitScreen />
            </el-icon>
          </el-button>
        </el-tooltip>
      </div>
      
      <div class="toggle-btn-wrapper">
        <el-tooltip content="分屏视图" placement="top" :show-after="300">
          <el-button 
            :type="modelValue.isFlowVisible && modelValue.isSourceCodeVisible ? 'primary' : ''"
            :class="{ 'is-active': modelValue.isFlowVisible && modelValue.isSourceCodeVisible }"
            size="small"
            @click="setViewMode('split')"
            class="toggle-btn view-btn"
            circle
          >
            <el-icon :size="16">
              <el-icon-grid />
            </el-icon>
          </el-button>
        </el-tooltip>
      </div>
    </div>
    
    <!-- 分隔线 -->
    <div class="divider"></div>
    
    <!-- 功能按钮组 -->
    <div class="button-group">
      <div class="toggle-btn-wrapper">
        <el-tooltip content="智能布局" placement="top" :show-after="300">
          <el-button 
            size="small"
            @click="optimizeLayout"
            class="toggle-btn function-btn optimize-btn"
            circle
          >
            <el-icon :size="16">
              <el-icon-magic-stick />
            </el-icon>
          </el-button>
        </el-tooltip>
      </div>
      
      <div class="toggle-btn-wrapper">
        <el-tooltip :content="isMiniMapVisible ? '隐藏小地图' : '显示小地图'" placement="top" :show-after="300">
          <el-button 
            :type="isMiniMapVisible ? 'primary' : ''"
            :class="{ 'is-active': isMiniMapVisible }"
            size="small"
            @click="toggleMiniMap"
            class="toggle-btn function-btn"
            circle
          >
            <el-icon :size="16">
              <el-icon-location />
            </el-icon>
          </el-button>
        </el-tooltip>
      </div>
      
      <div class="toggle-btn-wrapper">
        <el-tooltip :content="isAiAssistantVisible ? '隐藏AI助手' : '显示AI助手'" placement="top" :show-after="300">
          <el-button 
            :type="isAiAssistantVisible ? 'primary' : ''"
            :class="{ 'is-active': isAiAssistantVisible }"
            size="small"
            @click="toggleAiAssistant"
            class="toggle-btn function-btn ai-btn"
            circle
          >
            <el-icon :size="16">
              <el-icon-chat-line-round />
            </el-icon>
          </el-button>
        </el-tooltip>
      </div>
    </div>
  </div>
</template>

<style scoped>
.view-toggle-container {
  display: flex;
  flex-direction: row;
  gap: 8px;
  background: rgba(255, 255, 255, 0.98);
  border-radius: 16px;
  padding: 8px 12px;
  box-shadow: 
    0 8px 32px rgba(0, 0, 0, 0.12),
    0 4px 16px rgba(0, 0, 0, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  min-height: 56px;
  align-items: center;
  transition: box-shadow 0.3s cubic-bezier(0.4, 0, 0.2, 1); /* 只对阴影做动画 */
}

.view-toggle-container:hover {
  box-shadow: 
    0 12px 40px rgba(0, 0, 0, 0.15),
    0 6px 20px rgba(0, 0, 0, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.4);
  /* 移除transform变换，防止摆动 */
}

.button-group {
  display: flex;
  gap: 6px;
  align-items: center;
}

.toggle-btn-wrapper {
  display: flex;
  justify-content: center;
  position: relative;
}

.divider {
  width: 1px;
  height: 28px;
  background: linear-gradient(to bottom, transparent, #d1d5db, transparent);
  margin: 0 6px;
  opacity: 0.8;
}

/* 基础按钮样式 */
.toggle-btn.el-button {
  width: 40px !important;
  height: 40px !important;
  padding: 0 !important;
  border-radius: 12px !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
  border: 1px solid rgba(226, 232, 240, 0.8) !important;
  position: relative !important;
  overflow: hidden !important;
  font-weight: 500 !important;
}

/* 视图按钮样式 */
.view-btn.el-button:not(.el-button--primary):not(.is-active) {
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%) !important;
  color: #64748b !important;
  border-color: rgba(148, 163, 184, 0.3) !important;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05) !important;
}

.view-btn.el-button:not(.el-button--primary):not(.is-active):hover {
  background: linear-gradient(135deg, #e0f2fe 0%, #bae6fd 100%) !important;
  color: #0ea5e9 !important;
  border-color: rgba(14, 165, 233, 0.3) !important;
  /* 只保留缩放效果，移除垂直位移 */
  transform: scale(1.02) !important;
  box-shadow: 0 6px 16px rgba(14, 165, 233, 0.15) !important;
}

/* 功能按钮样式 */
.function-btn.el-button:not(.el-button--primary):not(.is-active) {
  background: linear-gradient(135deg, #fefefe 0%, #f9fafb 100%) !important;
  color: #6b7280 !important;
  border-color: rgba(156, 163, 175, 0.3) !important;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05) !important;
}

.function-btn.el-button:not(.el-button--primary):not(.is-active):hover {
  background: linear-gradient(135deg, #f3f4f6 0%, #e5e7eb 100%) !important;
  color: #4b5563 !important;
  border-color: rgba(75, 85, 99, 0.4) !important;
  /* 只保留轻微缩放效果 */
  transform: scale(1.02) !important;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.1) !important;
}

/* 优化布局按钮特殊样式 */
.optimize-btn.el-button:hover {
  background: linear-gradient(135deg, #fef3c7 0%, #fbbf24 100%) !important;
  color: #92400e !important;
  border-color: rgba(251, 191, 36, 0.5) !important;
  box-shadow: 0 6px 16px rgba(251, 191, 36, 0.3) !important;
}

/* AI按钮特殊样式 */
.ai-btn.el-button:not(.el-button--primary):not(.is-active):hover {
  background: linear-gradient(135deg, #d1fae5 0%, #6ee7b7 100%) !important;
  color: #059669 !important;
  border-color: rgba(5, 150, 105, 0.4) !important;
  /* 减少变换效果，只保留轻微缩放 */
  transform: scale(1.02) !important;
  box-shadow: 0 6px 16px rgba(5, 150, 105, 0.2) !important;
}

.ai-btn.el-button--primary {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%) !important;
  border-color: #10b981 !important;
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.4) !important;
}

.ai-btn.el-button--primary:hover {
  background: linear-gradient(135deg, #34d399 0%, #10b981 100%) !important;
  border-color: #34d399 !important;
  /* 减少变换，保持位置稳定 */
  transform: scale(1.02) !important;
  box-shadow: 0 6px 20px rgba(52, 211, 153, 0.5) !important;
}

/* 激活状态通用样式 */
.toggle-btn.el-button.el-button--primary:not(.ai-btn),
.toggle-btn.el-button.is-active:not(.ai-btn) {
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%) !important;
  border-color: #3b82f6 !important;
  color: white !important;
  transform: translateY(-2px) !important;
  box-shadow: 0 6px 20px rgba(59, 130, 246, 0.4) !important;
  position: relative !important;
}

.toggle-btn.el-button.el-button--primary:hover:not(.ai-btn),
.toggle-btn.el-button.is-active:hover:not(.ai-btn) {
  background: linear-gradient(135deg, #60a5fa 0%, #3b82f6 100%) !important;
  border-color: #60a5fa !important;
  /* 简化transform，只保留轻微缩放 */
  transform: translateY(-2px) scale(1.02) !important;
  box-shadow: 0 8px 25px rgba(96, 165, 250, 0.5) !important;
}

/* 激活状态指示器 */
.toggle-btn.is-active::before {
  content: '';
  position: absolute;
  top: -3px;
  left: -3px;
  right: -3px;
  bottom: -3px;
  border: 2px solid rgba(59, 130, 246, 0.6);
  border-radius: 15px;
  opacity: 1;
  animation: activeIndicator 2s ease-in-out infinite;
  z-index: -1;
}

.ai-btn.is-active::before {
  border-color: rgba(16, 185, 129, 0.6);
  animation: aiActiveIndicator 2s ease-in-out infinite;
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

@keyframes aiActiveIndicator {
  0%, 100% {
    transform: scale(1);
    opacity: 0.8;
    border-color: rgba(16, 185, 129, 0.6);
  }
  50% {
    transform: scale(1.15);
    opacity: 0.4;
    border-color: rgba(52, 211, 153, 0.8);
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .view-toggle-container {
    padding: 6px 8px;
    gap: 6px;
  }
  
  .button-group {
    gap: 4px;
  }
  
  .toggle-btn.el-button {
    width: 36px !important;
    height: 36px !important;
  }
  
  .divider {
    height: 24px;
    margin: 0 4px;
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
  /* 移除位移变换，只保留缩放 */
  transform: scale(1.02) !important;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.2) !important;
}

/* 确保激活状态在深色模式下也清晰可见 */
.dark .view-toggle-container .toggle-btn.el-button.el-button--primary,
.dark .view-toggle-container .toggle-btn.el-button.is-active {
  background: var(--el-color-primary) !important;
  border-color: var(--el-color-primary) !important;
  color: white !important;
  /* 移除位移变换 */
  transform: none !important;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.4) !important;
}

.dark .view-toggle-container .toggle-btn.el-button.el-button--primary:hover,
.dark .view-toggle-container .toggle-btn.el-button.is-active:hover {
  background: var(--el-color-primary-light-3) !important;
  border-color: var(--el-color-primary-light-3) !important;
  /* 简化变换效果 */
  transform: scale(1.02) !important;
  box-shadow: 0 4px 16px rgba(64, 158, 255, 0.5) !important;
}
</style>
