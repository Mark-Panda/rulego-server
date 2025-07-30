<script setup>
import { onMounted, onBeforeUnmount, ref, watch, nextTick, reactive } from 'vue'
import AppDesign from '@src/pages/workflow/app-design/app-design.vue';
import AppSourceCode from '@src/pages/workflow/app-source-code/app-source-code.vue';
import { mapFlowDataModelToRuleGoModel } from '@src/pages/workflow/app-design/utils';
import { cloneDeep, } from 'lodash-es';
import EventBus from '@src/utils/event-bus';
import ChatListView from "@src/pages/workflow/chat-list-view/chat-list-view.vue";
import ToolButtons from '@src/pages/workflow/app-design-new/tool-buttons.vue';
//import Assistant from '@src/assets/assistant.svg';
import AiChatDrawer from '@src/pages/workflow/chat-list-view/ai-chat-drawer.vue';

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
const isMiniMapVisible = ref(true); // 小地图显示状态

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
 * 处理图更新
 * @param {string} newVal 
 */
function handleDesignUpdate(newVal) {
  try {
    val.value = newVal;
    emit("update:modelValue", val.value);
  } catch (error) {
    console.error('Error in handleDesignUpdate:', error);
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
  <div class="flex flex-row h-full relative">
    <div v-if="toolButtons.isFlowVisible" class="flex-1 overflow-hidden">
      <app-design :model-value="val" @update:model-value="handleDesignUpdate" ref="appDesignRef" />
    </div>
    <div v-if="toolButtons.isSourceCodeVisible" class="flex-1 overflow-hidden">
      <app-source-code :model-value="val" @update:model-value="handelSourceCodeUpdate" ref="appSourceCodeRef" />
    </div>
    <div v-if="false" class="flex flex-1 overflow-hidden">
      <chat-list-view class="flex-1" />
    </div>
    <tool-buttons 
      v-model="toolButtons" 
      :is-mini-map-visible="isMiniMapVisible"
      @optimize-layout="handleOptimizeLayout" 
      @toggle-minimap="handleToggleMiniMap"
      class="absolute bottom-4 left-1/2 transform -translate-x-1/2 z-50" 
    />
    <!--20250507暂时注释掉ai助手-->
    <!-- <div @click="isAiVisible = true" class="absolute bottom-24 right-10 size-[50px] flex justify-center items-center rounded-full bg-white cursor-pointer border">
      <el-icon :size="28">
        <Assistant />
      </el-icon>
    </div> -->
    <ai-chat-drawer v-model="isAiVisible" />
  </div>
</template>
<style scoped>
.app-split-screen {
  display: flex;
  flex-direction: column;
}
</style>