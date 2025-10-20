<script setup>
import { onMounted, onBeforeUnmount, ref, watch, nextTick, reactive } from 'vue'
import AppDesign from '@src/pages/workflow/app-design/app-design.vue';
import AppSourceCode from '@src/pages/workflow/app-source-code/app-source-code.vue';
import { mapFlowDataModelToRuleGoModel } from '@src/pages/workflow/app-design/utils';
import { cloneDeep, } from 'lodash-es';
import EventBus from '@src/utils/event-bus';
import ChatListView from "@src/pages/workflow/chat-list-view/chat-list-view.vue";
import SmartAiAssistant from "@src/pages/workflow/chat-list-view/smart-ai-assistant.vue";
import ToolButtons from '@src/pages/workflow/app-design-new/tool-buttons.vue';
//import Assistant from '@src/assets/assistant.svg';
import AiChatDrawer from '@src/pages/workflow/chat-list-view/ai-chat-drawer.vue';
import { SchemaHistoryManager } from '@src/utils/schema-history-manager.js';
import { anchorUpdateManager } from '@src/utils/anchor-update-manager.js';

const props = defineProps({
  modelValue: {
    type: [String, Number, Object],
    default: '',
  },
});

const emit = defineEmits(["update:modelValue"]);
const logicflowNodeMouseUp = EventBus.logicflowNodeMouseUp();


const val = ref(cloneDeep(props.modelValue));
const appDesignRef = ref(null);
const appSourceCodeRef = ref(null);
const toolButtons = ref({
  isFlowVisible: true,
  isSourceCodeVisible: false,
});
const isAiVisible = ref(false);
const isAiAssistantVisible = ref(false); // 新增智能AI助手状态
const isMiniMapVisible = ref(true); // 小地图显示状态
const isMobile = ref(false); // 移动端检测

// 初始化schema历史记录管理器
const schemaHistoryManager = new SchemaHistoryManager();

// AI助手侧边栏宽度相关状态
const aiSidebarWidth = ref(420); // 默认宽度420px
const minSidebarWidth = 320; // 最小宽度
const maxSidebarWidth = 800; // 最大宽度
const isResizing = ref(false);
const resizeStartX = ref(0);
const resizeStartWidth = ref(0);

function handelDesignToJson() {
  const flowData = appDesignRef.value.getData();
  const ruleGoModel = mapFlowDataModelToRuleGoModel(flowData, val.value);
  val.value = ruleGoModel;
  emit("update:modelValue", val.value);
}

let jsonToDocTimer = 0;
async function handelJsonToDesign() {
  try {
    //连续输入code不要频繁更新图。
    if (jsonToDocTimer) {
      clearTimeout(jsonToDocTimer);
    }
    jsonToDocTimer = setTimeout(async () => {
      // 等待下一个tick确保DOM完全挂载
      await nextTick();
      
      // 添加安全检查，确保组件已挂载且render方法存在
      if (appDesignRef.value && typeof appDesignRef.value.render === 'function') {
        try {
          // 再次确保组件已完全初始化
          setTimeout(() => {
            if (appDesignRef.value && typeof appDesignRef.value.render === 'function') {
              appDesignRef.value.render();
            }
          }, 100);
        } catch (error) {
          console.warn('渲染失败，组件可能未完全初始化:', error);
        }
      }
    }, 500);
  } catch (error) {
    console.error('Error in handelJsonToDesign:', error);
  }
}

function getData() {
  try {
    return val.value || { nodes: [], edges: [] };
  } catch (error) {
    console.error('Error in getData:', error);
    return { nodes: [], edges: [] };
  }
}

/**
 * 处理代码更新
 * @param {string} newVal 
 */
function handelSourceCodeUpdate(newVal) {
  try {
    val.value = newVal;
    emit("update:modelValue", val.value);
    handelJsonToDesign();
  } catch (error) {
    console.error('Error in handelSourceCodeUpdate:', error);
  }
}

/**
 * 处理设计器更新
 */
function handleDesignUpdate(newVal) {
  try {
    // 记录设计器的手动修改
    if (val.value && JSON.stringify(newVal) !== JSON.stringify(val.value)) {
      schemaHistoryManager.addRecord(
        newVal,
        'manual',
        '手动修改流程图',
        { 
          manual: true,
          timestamp: Date.now()
        }
      );
    }
    
    val.value = newVal;
    emit("update:modelValue", val.value);
  } catch (error) {
    console.error('Error in handleDesignUpdate:', error);
  }
}

/**
 * 处理手动操作的schema更新
 */
function handleManualSchemaUpdate(newVal) {
  try {
    // 记录手动修改
    if (val.value && JSON.stringify(newVal) !== JSON.stringify(val.value)) {
      schemaHistoryManager.addRecord(
        newVal,
        'manual',
        '手动修改流程图',
        { 
          manual: true,
          timestamp: Date.now()
        }
      );
    }
    
    val.value = newVal;
    emit("update:modelValue", val.value);
  } catch (error) {
    console.error('Error in handleManualSchemaUpdate:', error);
  }
}

/**
 * 处理优化布局
 */
async function handleOptimizeLayout() {
  try {
    // 等待DOM更新
    await nextTick();
    
    if (appDesignRef.value) {
      const lf = appDesignRef.value.getLf();
      if (lf) {
        try {
          // 延迟执行以确保LogicFlow完全初始化
          setTimeout(() => {
            if (appDesignRef.value?.getLf()) {
              // 执行布局优化
              optimizeNodeLayout(appDesignRef.value.getLf());
              
              // 触发数据更新
              handelDesignToJson();
            }
          }, 200);
          
        } catch (error) {
          console.error('布局优化失败:', error);
        }
      }
    }
  } catch (error) {
    console.error('Error in handleOptimizeLayout:', error);
  }
}

/**
 * 同步小地图状态
 */
async function syncMiniMapState() {
  console.log('app-design-new：syncMiniMapState 被调用');
  
  if (appDesignRef.value && appDesignRef.value.getMiniMapVisible) {
    try {
      const currentState = await appDesignRef.value.getMiniMapVisible();
      console.log('app-design-new：获取到的小地图状态:', currentState);
      isMiniMapVisible.value = currentState;
      console.log('app-design-new：同步后的 isMiniMapVisible:', isMiniMapVisible.value);
    } catch (error) {
      console.error('app-design-new：同步小地图状态时出错:', error);
    }
  } else {
    console.warn('app-design-new：appDesignRef.value 或 getMiniMapVisible 方法不存在');
  }
}

/**
 * 处理小地图切换
 */
async function handleToggleMiniMap() {
  console.log('app-design-new：handleToggleMiniMap 被调用');
  
  if (appDesignRef.value && appDesignRef.value.toggleMiniMap) {
    console.log('app-design-new：调用 app-design 的 toggleMiniMap 方法');
    try {
      const result = await appDesignRef.value.toggleMiniMap();
      console.log('app-design-new：toggleMiniMap 返回状态:', result);
      
      // 同步状态
      isMiniMapVisible.value = result;
      console.log('app-design-new：更新后的 isMiniMapVisible:', isMiniMapVisible.value);
    } catch (error) {
      console.error('app-design-new：调用 toggleMiniMap 时出错:', error);
    }
  } else {
    console.error('app-design-new：appDesignRef.value 或 toggleMiniMap 方法不存在');
  }
}

/**
 * 处理AI助手切换
 */
function handleToggleAiAssistant() {
  console.log('app-design-new：handleToggleAiAssistant 被调用');
  isAiAssistantVisible.value = !isAiAssistantVisible.value;
  console.log('app-design-new：AI助手状态更新为:', isAiAssistantVisible.value);
}

/**
 * 开始拖拽调整侧边栏宽度
 */
function handleResizeStart(event) {
  if (isMobile.value) return; // 移动端不支持拖拽调整
  
  isResizing.value = true;
  resizeStartX.value = event.clientX;
  resizeStartWidth.value = aiSidebarWidth.value;
  
  // 添加全局事件监听
  document.addEventListener('mousemove', handleResizeMove);
  document.addEventListener('mouseup', handleResizeEnd);
  
  // 防止文本选中
  document.body.style.userSelect = 'none';
  document.body.style.cursor = 'col-resize';
}

/**
 * 拖拽调整过程
 */
function handleResizeMove(event) {
  if (!isResizing.value) return;
  
  const deltaX = resizeStartX.value - event.clientX; // 向左拖拽为正值
  const newWidth = resizeStartWidth.value + deltaX;
  
  // 限制宽度范围
  aiSidebarWidth.value = Math.max(minSidebarWidth, Math.min(maxSidebarWidth, newWidth));
}

/**
 * 结束拖拽调整
 */
function handleResizeEnd() {
  isResizing.value = false;
  
  // 移除全局事件监听
  document.removeEventListener('mousemove', handleResizeMove);
  document.removeEventListener('mouseup', handleResizeEnd);
  
  // 恢复样式
  document.body.style.userSelect = '';
  document.body.style.cursor = '';
  
  // 保存宽度到本地存储
  localStorage.setItem('aiSidebarWidth', aiSidebarWidth.value.toString());
}

/**
 * 处理AI更新schema
 */
function handleAiSchemaUpdate(newSchema) {
  console.log('app-design-new：收到AI的schema更新:', newSchema);
  try {
    // 记录AI修改
    schemaHistoryManager.addRecord(
      newSchema,
      'ai',
      'AI智能助手修改',
      { 
        ai: true,
        timestamp: Date.now(),
        previousSchema: JSON.parse(JSON.stringify(val.value))
      }
    );
    
    val.value = newSchema;
    emit("update:modelValue", val.value);
    
    // 延迟更新画布，确保数据同步并强制更新所有节点锚点
    setTimeout(() => {
      handelJsonToDesign();
      // 额外延迟确保渲染完成后更新锚点
      setTimeout(() => {
        forceUpdateAllNodeAnchors();
      }, 300);
    }, 200);
    
    console.log('app-design-new：schema更新成功，已保存到历史记录');
  } catch (error) {
    console.error('app-design-new：工作流AI schema更新失败:', error);
  }
}

/**
 * 更新所有边的位置
 * @param {Object} lf LogicFlow 实例
 */
function updateAllEdgesPosition(lf) {
  const graphData = lf.getGraphData();
  const nodes = graphData.nodes;
  
  // 更新所有节点的出边
  nodes.forEach(node => {
    const nodeModel = lf.getNodeModelById(node.id);
    if (nodeModel && nodeModel.outgoing && nodeModel.outgoing.edges) {
      nodeModel.outgoing.edges.forEach(edge => {
        if (edge.updatePathByAnchor) {
          edge.updatePathByAnchor();
        }
      });
    }
  });
}

/**
 * 强制更新所有节点的锚点
 * 解决AI修改schema后节点输出端点显示问题
 * 使用优化后的锚点更新管理器
 */
function forceUpdateAllNodeAnchors() {
  console.log('开始强制更新所有节点的锚点');
  
  if (!appDesignRef.value) {
    console.warn('appDesignRef不存在，无法更新锚点');
    return;
  }
  
  try {
    const lf = appDesignRef.value.getLf();
    if (!lf) {
      console.warn('LogicFlow实例不存在，无法更新锚点');
      return;
    }
    
    const graphData = lf.getGraphData();
    const nodes = graphData.nodes || [];
    
    console.log(`找到 ${nodes.length} 个节点，开始使用锚点管理器更新`);
    
    // 使用锚点更新管理器批量更新
    anchorUpdateManager.batchUpdateAnchors(lf, nodes, {
      maxBatchSize: 10,
      batchDelay: 30,
      forceUpdate: true
    }).then(() => {
      console.log('锚点管理器批量更新完成');
      
      // 更新所有连线位置
      setTimeout(() => {
        try {
          updateAllEdgesPosition(lf);
          console.log('所有节点锚点和连线更新完成');
        } catch (error) {
          console.error('更新连线位置时出错:', error);
        }
      }, 100);
    });
    
  } catch (error) {
    console.error('强制更新节点锚点时出错:', error);
  }
}

/**
 * 优化节点布局
 * @param {Object} lf LogicFlow 实例
 */
function optimizeNodeLayout(lf) {
  const graphData = lf.getGraphData();
  const nodes = graphData.nodes;
  const edges = graphData.edges;
  
  if (nodes.length === 0) return;
  
  // 找到开始节点
  const startNode = nodes.find(node => node.type === 'start');
  if (!startNode) {
    // 如果没有开始节点，只执行 fitView
    lf.fitView();
    return;
  }
  
  // 布局参数 - 优化间距避免连线遮挡
  const HORIZONTAL_SPACING = 300; // 增加水平层级间距
  const VERTICAL_SPACING = 150;   // 增加同层节点垂直间距
  const START_X = 150;
  const START_Y = 300;
  
  // 构建图的邻接表和入度表
  const adjacencyList = new Map();
  const inDegree = new Map();
  const outDegree = new Map();
  
  nodes.forEach(node => {
    adjacencyList.set(node.id, []);
    inDegree.set(node.id, 0);
    outDegree.set(node.id, 0);
  });
  
  edges.forEach(edge => {
    if (adjacencyList.has(edge.sourceNodeId)) {
      adjacencyList.get(edge.sourceNodeId).push(edge.targetNodeId);
      outDegree.set(edge.sourceNodeId, outDegree.get(edge.sourceNodeId) + 1);
    }
    if (inDegree.has(edge.targetNodeId)) {
      inDegree.set(edge.targetNodeId, inDegree.get(edge.targetNodeId) + 1);
    }
  });
  
  // 使用改进的拓扑排序进行层次分配
  const levels = [];
  const nodeLevel = new Map();
  const queue = [];
  
  // 找到所有入度为0的节点作为起始点
  nodes.forEach(node => {
    if (inDegree.get(node.id) === 0) {
      queue.push({ nodeId: node.id, level: 0 });
      nodeLevel.set(node.id, 0);
    }
  });
  
  // 如果没有入度为0的节点，从start节点开始
  if (queue.length === 0) {
    queue.push({ nodeId: startNode.id, level: 0 });
    nodeLevel.set(startNode.id, 0);
  }
  
  // 拓扑排序分层
  while (queue.length > 0) {
    const { nodeId, level } = queue.shift();
    
    // 确保层级数组存在
    if (!levels[level]) {
      levels[level] = [];
    }
    levels[level].push(nodeId);
    
    // 处理邻居节点
    const neighbors = adjacencyList.get(nodeId) || [];
    neighbors.forEach(neighborId => {
      if (!nodeLevel.has(neighborId)) {
        const newLevel = level + 1;
        nodeLevel.set(neighborId, newLevel);
        queue.push({ nodeId: neighborId, level: newLevel });
      }
    });
  }
  
  // 处理未分配层级的节点
  nodes.forEach(node => {
    if (!nodeLevel.has(node.id)) {
      const lastLevel = levels.length;
      if (!levels[lastLevel]) {
        levels[lastLevel] = [];
      }
      levels[lastLevel].push(node.id);
      nodeLevel.set(node.id, lastLevel);
    }
  });
  
  // 优化每层节点的排序，减少连线交叉
  levels.forEach((levelNodes, levelIndex) => {
    if (levelNodes.length <= 1) return;
    
    // 根据节点的连接关系进行排序
    levelNodes.sort((a, b) => {
      // 获取前一层连接到当前节点的节点
      const aIncomingNodes = edges
        .filter(edge => edge.targetNodeId === a && nodeLevel.get(edge.sourceNodeId) === levelIndex - 1)
        .map(edge => edge.sourceNodeId);
      const bIncomingNodes = edges
        .filter(edge => edge.targetNodeId === b && nodeLevel.get(edge.sourceNodeId) === levelIndex - 1)
        .map(edge => edge.sourceNodeId);
      
      if (aIncomingNodes.length === 0 && bIncomingNodes.length === 0) return 0;
      if (aIncomingNodes.length === 0) return 1;
      if (bIncomingNodes.length === 0) return -1;
      
      // 根据前一层节点的位置进行排序
      const aAvgPos = aIncomingNodes.reduce((sum, nodeId) => {
        const prevLevel = levels[levelIndex - 1];
        return sum + prevLevel.indexOf(nodeId);
      }, 0) / aIncomingNodes.length;
      
      const bAvgPos = bIncomingNodes.reduce((sum, nodeId) => {
        const prevLevel = levels[levelIndex - 1];
        return sum + prevLevel.indexOf(nodeId);
      }, 0) / bIncomingNodes.length;
      
      return aAvgPos - bAvgPos;
    });
  });
  
  // 计算每层节点的位置并移动 - 水平布局
  levels.forEach((levelNodes, levelIndex) => {
    if (levelNodes.length === 0) return;
    
    // 水平布局：X坐标按层级递增，Y坐标按同层节点分布
    const levelX = START_X + levelIndex * HORIZONTAL_SPACING;
    
    // 计算该层节点的总高度和起始Y坐标
    const totalHeight = Math.max(0, (levelNodes.length - 1) * VERTICAL_SPACING);
    const startY = START_Y - totalHeight / 2;
    
    levelNodes.forEach((nodeId, nodeIndex) => {
      const y = levelNodes.length === 1 ? START_Y : startY + nodeIndex * VERTICAL_SPACING;
      const nodeModel = lf.getNodeModelById(nodeId);
      if (nodeModel) {
        nodeModel.moveTo(levelX, y);
      }
    });
  });
  
  // 分步骤更新连线，确保连线不丢失
  setTimeout(() => {
    // 第一步：更新所有边的路径
    updateAllEdgesPosition(lf);
    
    // 第二步：延迟执行 fitView
    setTimeout(() => {
      lf.fitView();
    }, 150);
  }, 150);
}

watch(
  () => props.modelValue,
  (newVal) => {
    val.value = newVal; // 当父组件更新 modelValue 时，同步更新 val
  },
  { deep: true }
);

onMounted(async () => {
  try {
    // 检测移动端
    const checkMobile = () => {
      isMobile.value = window.innerWidth <= 768;
    };
    checkMobile();
    window.addEventListener('resize', checkMobile);
    
    // 恢复保存的AI助手侧边栏宽度
    const savedWidth = localStorage.getItem('aiSidebarWidth');
    if (savedWidth) {
      const width = parseInt(savedWidth, 10);
      if (width >= minSidebarWidth && width <= maxSidebarWidth) {
        aiSidebarWidth.value = width;
      }
    }
    
    // 初始化schema历史记录管理器
    if (val.value) {
      schemaHistoryManager.addRecord(
        val.value,
        'initial',
        '初始化画布状态',
        { 
          initialization: true,
          timestamp: Date.now()
        }
      );
    }
    
    // 等待DOM完全挂载
    await nextTick();
    
    // 等待一个短暂的时间确保所有子组件都已挂载
    setTimeout(() => {
      try {
        logicflowNodeMouseUp.on(handelDesignToJson);
        
        // 同步小地图状态
        setTimeout(() => {
          syncMiniMapState();
        }, 500); // 等待LogicFlow完全初始化
      } catch (error) {
        console.error('事件监听器注册失败:', error);
      }
    }, 100);
  } catch (error) {
    console.error('组件挂载失败:', error);
  }
});

onBeforeUnmount(() => {
  try {
    // 清理事件监听器
    if (logicflowNodeMouseUp && typeof logicflowNodeMouseUp.off === 'function') {
      logicflowNodeMouseUp.off(handelDesignToJson);
    }
    
    // 清理定时器
    if (jsonToDocTimer) {
      clearTimeout(jsonToDocTimer);
      jsonToDocTimer = 0;
    }
    
    // 清理拖拽相关事件监听器
    document.removeEventListener('mousemove', handleResizeMove);
    document.removeEventListener('mouseup', handleResizeEnd);
    
    // 清理组件引用
    appDesignRef.value = null;
    appSourceCodeRef.value = null;
  } catch (error) {
    console.error('组件销毁时出错:', error);
  }
});

defineExpose({
  getData,
  handelDesignToJson,
  handelJsonToDesign
});

</script>
<template>
  <!-- 当AI助手打开时，主内容区域自适应调整 -->
  <div class="app-design-container" :class="{ 'ai-assistant-open': isAiAssistantVisible }">
    <!-- 主内容区域 -->
    <div class="main-content">
      <div v-if="toolButtons.isFlowVisible" class="content-panel flow-panel">
        <app-design :model-value="val" @update:model-value="handleDesignUpdate" ref="appDesignRef" />
      </div>
      <div v-if="toolButtons.isSourceCodeVisible" class="content-panel code-panel">
        <app-source-code :model-value="val" @update:model-value="handelSourceCodeUpdate" ref="appSourceCodeRef" />
      </div>
      <div v-if="false" class="content-panel chat-panel">
        <chat-list-view class="flex-1" />
      </div>
    </div>
    
    <!-- 工具栏 -->
    <tool-buttons 
      v-model="toolButtons" 
      :is-mini-map-visible="isMiniMapVisible"
      :is-ai-assistant-visible="isAiAssistantVisible"
      @optimize-layout="handleOptimizeLayout" 
      @toggle-minimap="handleToggleMiniMap"
      @toggle-ai-assistant="handleToggleAiAssistant"
      class="floating-toolbar" 
    />
    
    <!-- AI助手侧边栏 -->
    <transition name="sidebar-slide">
      <div 
        v-if="isAiAssistantVisible" 
        class="ai-sidebar"
        :class="{ 'sidebar-mobile': isMobile }"
        :style="{ width: isMobile ? '100vw' : aiSidebarWidth + 'px' }"
      >
        <!-- 拖拽手柄 -->
        <div 
          v-if="!isMobile"
          class="resize-handle"
          @mousedown="handleResizeStart"
          :class="{ 'resizing': isResizing }"
        >
          <div class="resize-indicator"></div>
        </div>
        <!-- 侧边栏头部 -->
        <div class="sidebar-header">
          <div class="header-content">
            <div class="header-icon">
              <el-icon :size="22" class="text-emerald-600">
                <el-icon-chat-line-round />
              </el-icon>
            </div>
            <div class="header-text">
              <h3 class="header-title">AI智能助手</h3>
              <p class="header-subtitle">实时修改流程图配置</p>
            </div>
          </div>
          <el-button 
            size="small" 
            circle 
            @click="isAiAssistantVisible = false"
            class="close-btn"
          >
            <el-icon><el-icon-close /></el-icon>
          </el-button>
        </div>
        
        <!-- AI助手组件 -->
        <div class="sidebar-body">
          <smart-ai-assistant 
            :current-schema="val" 
            :on-schema-update="handleAiSchemaUpdate"
            :schema-history-manager="schemaHistoryManager"
          />
        </div>
      </div>
    </transition>
    
    <!-- 移动端遮罩 -->
    <transition name="overlay-fade">
      <div 
        v-if="isAiAssistantVisible && isMobile" 
        class="mobile-overlay"
        @click="isAiAssistantVisible = false"
      ></div>
    </transition>
    
    <!-- 原有的AI聊天抽屉（保留兼容性） -->
    <ai-chat-drawer v-model="isAiVisible" />
  </div>
</template>
<style scoped>
.app-design-container {
  height: 100%;
  position: relative;
  display: flex;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 主内容区域 */
.main-content {
  flex: 1;
  display: flex;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.ai-assistant-open .main-content {
  margin-right: v-bind('aiSidebarWidth + "px"');
}

.content-panel {
  flex: 1;
  overflow: hidden;
  position: relative;
}

/* 工具栏样式 */
.floating-toolbar {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 60;
  transition: all 0.3s ease;
}

.ai-assistant-open .floating-toolbar {
  left: v-bind('"calc(50% - " + (aiSidebarWidth / 2) + "px)"'); /* 向左偏移侧边栏宽度的一半 */
}

/* AI侧边栏样式 */
.ai-sidebar {
  position: absolute;
  right: 0;
  top: 0;
  height: 100%;
  min-width: v-bind('minSidebarWidth + "px"');
  max-width: v-bind('maxSidebarWidth + "px"');
  background: linear-gradient(135deg, #ffffff 0%, #f8fafc 100%);
  border-left: 1px solid #e2e8f0;
  box-shadow: 
    -8px 0 32px rgba(0, 0, 0, 0.12),
    -4px 0 16px rgba(0, 0, 0, 0.08);
  z-index: 50;
  display: flex;
  flex-direction: column;
  backdrop-filter: blur(20px);
}

/* 侧边栏头部 */
.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  background: linear-gradient(135deg, #f0fdfa 0%, #ecfdf5 100%);
  border-bottom: 1px solid #d1fae5;
  backdrop-filter: blur(10px);
}

.header-content {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);
  color: white;
}

.header-text {
  flex: 1;
}

.header-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #065f46;
  line-height: 1.2;
}

.header-subtitle {
  margin: 2px 0 0;
  font-size: 12px;
  color: #047857;
  opacity: 0.8;
  line-height: 1.2;
}

.close-btn {
  color: #6b7280;
  border: 1px solid #d1d5db;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(4px);
  transition: all 0.2s ease;
}

.close-btn:hover {
  color: #ef4444;
  border-color: #fca5a5;
  background: #fef2f2;
  transform: scale(1.05);
}

/* 侧边栏主体 */
.sidebar-body {
  flex: 1;
  overflow: hidden;
}

/* 拖拽手柄样式 */
.resize-handle {
  position: absolute;
  left: -5px;
  top: 0;
  width: 10px;
  height: 100%;
  cursor: col-resize;
  z-index: 60;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.resize-handle:hover .resize-indicator,
.resize-handle.resizing .resize-indicator {
  background: #3b82f6;
  box-shadow: 0 0 8px rgba(59, 130, 246, 0.4);
}

.resize-indicator {
  width: 2px;
  height: 60px;
  background: #d1d5db;
  border-radius: 1px;
  transition: all 0.2s ease;
}

.resize-handle:hover .resize-indicator {
  height: 80px;
}

.resize-handle.resizing .resize-indicator {
  height: 100px;
  background: #1d4ed8;
}

/* 拖拽状态样式 */
.resize-handle.resizing {
  background: rgba(59, 130, 246, 0.1);
}

/* 全局拖拽状态 */
body.resizing {
  cursor: col-resize !important;
  user-select: none !important;
}

/* 移动端适配 */
.sidebar-mobile {
  width: 100vw;
  box-shadow: none;
  border-left: none;
}

.mobile-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 40;
  backdrop-filter: blur(4px);
}

/* 过渡动画 */
.sidebar-slide-enter-active,
.sidebar-slide-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.sidebar-slide-enter-from {
  transform: translateX(100%);
  opacity: 0;
}

.sidebar-slide-leave-to {
  transform: translateX(100%);
  opacity: 0;
}

.overlay-fade-enter-active,
.overlay-fade-leave-active {
  transition: all 0.2s ease;
}

.overlay-fade-enter-from,
.overlay-fade-leave-to {
  opacity: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .ai-assistant-open .main-content {
    margin-right: 0;
  }
  
  .ai-assistant-open .floating-toolbar {
    left: 50%;
  }
  
  .sidebar-header {
    padding: 16px 20px;
  }
  
  .header-icon {
    width: 40px;
    height: 40px;
  }
  
  .header-title {
    font-size: 15px;
  }
  
  .header-subtitle {
    font-size: 11px;
  }
}

@media (max-width: 480px) {
  .sidebar-header {
    padding: 14px 16px;
  }
  
  .header-content {
    gap: 10px;
  }
  
  .header-icon {
    width: 36px;
    height: 36px;
  }
  
  .header-title {
    font-size: 14px;
  }
}

/* 深色模式适配 */
@media (prefers-color-scheme: dark) {
  .ai-sidebar {
    background: linear-gradient(135deg, #1e293b 0%, #0f172a 100%);
    border-left-color: #374151;
  }
  
  .sidebar-header {
    background: linear-gradient(135deg, #064e3b 0%, #022c22 100%);
    border-bottom-color: #065f46;
  }
  
  .header-title {
    color: #d1fae5;
  }
  
  .header-subtitle {
    color: #a7f3d0;
  }
  
  .close-btn {
    color: #9ca3af;
    border-color: #4b5563;
    background: rgba(31, 41, 55, 0.8);
  }
  
  .close-btn:hover {
    color: #f87171;
    border-color: #ef4444;
    background: #7f1d1d;
  }
  
  .mobile-overlay {
    background: rgba(0, 0, 0, 0.7);
  }
}
</style>