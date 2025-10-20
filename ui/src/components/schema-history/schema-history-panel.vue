<template>
  <div class="schema-history-panel">
    <!-- 历史记录头部 -->
    <div class="history-header">
      <div class="header-title">
        <el-icon :size="16" class="text-indigo-600">
          <el-icon-clock />
        </el-icon>
        <span class="title-text">修改历史</span>
        <el-badge :value="historyList.length" :max="99" class="history-count" />
      </div>
      
      <div class="header-actions">
        <el-tooltip content="撤销" placement="top">
          <el-button 
            :disabled="!canUndo" 
            size="small" 
            circle 
            @click="handleUndo"
            class="action-btn undo-btn"
          >
            <el-icon><el-icon-refresh-left /></el-icon>
          </el-button>
        </el-tooltip>
        
        <el-tooltip content="重做" placement="top">
          <el-button 
            :disabled="!canRedo" 
            size="small" 
            circle 
            @click="handleRedo"
            class="action-btn redo-btn"
          >
            <el-icon><el-icon-refresh-right /></el-icon>
          </el-button>
        </el-tooltip>
        
        <el-dropdown @command="handleMenuCommand" trigger="click">
          <el-button size="small" circle class="action-btn">
            <el-icon><el-icon-more /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="export">导出历史</el-dropdown-item>
              <el-dropdown-item command="import">导入历史</el-dropdown-item>
              <el-dropdown-item command="clear" divided>清空历史</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <!-- 统计信息 -->
    <div class="statistics-section" v-if="statistics">
      <div class="stat-item">
        <span class="stat-label">总记录</span>
        <span class="stat-value">{{ statistics.totalRecords }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">AI修改</span>
        <span class="stat-value ai-count">{{ statistics.aiModifications }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">手动修改</span>
        <span class="stat-value manual-count">{{ statistics.manualModifications }}</span>
      </div>
    </div>

    <!-- 历史记录列表 -->
    <div class="history-list">
      <div class="list-container">
        <div 
          v-for="record in displayHistoryList" 
          :key="record.id"
          :class="[
            'history-item',
            { 'is-current': record.isCurrent },
            { 'is-ai': record.isAI }
          ]"
          @click="handleJumpToRecord(record.id)"
        >
          <!-- 操作类型图标 -->
          <div class="operation-icon">
            <el-icon v-if="record.isAI" class="ai-icon">
              <el-icon-cpu />
            </el-icon>
            <el-icon v-else-if="record.operation === 'auto-save'" class="auto-save-icon">
              <el-icon-download />
            </el-icon>
            <el-icon v-else class="manual-icon">
              <el-icon-edit />
            </el-icon>
          </div>

          <!-- 记录信息 -->
          <div class="record-info">
            <div class="record-header">
              <span class="operation-type">
                {{ getOperationLabel(record.operation) }}
              </span>
              <span class="record-time">
                {{ formatTime(record.timestamp) }}
              </span>
            </div>
            
            <div class="record-description">
              {{ record.description || '无描述' }}
            </div>
            
            <div class="record-metadata">
              <span class="metadata-item">
                <el-icon :size="12"><el-icon-connection /></el-icon>
                {{ record.metadata.nodeCount }}节点
              </span>
              <span class="metadata-item">
                <el-icon :size="12"><el-icon-link /></el-icon>
                {{ record.metadata.edgeCount }}连线
              </span>
              <el-button 
                v-if="canShowDiff(record)"
                size="small" 
                type="primary" 
                text
                @click.stop="handleShowDiff(record)"
                class="diff-btn"
              >
                <el-icon :size="12"><el-icon-view /></el-icon>
                查看差异
              </el-button>
            </div>
          </div>

          <!-- 当前指示器 -->
          <div v-if="record.isCurrent" class="current-indicator">
            <el-icon><el-icon-location /></el-icon>
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-if="historyList.length === 0" class="empty-state">
      <el-icon :size="48" class="empty-icon">
        <el-icon-document />
      </el-icon>
      <p class="empty-text">暂无修改历史</p>
      <p class="empty-hint">开始修改流程图后将自动记录历史</p>
    </div>

    <!-- 隐藏的文件输入 -->
    <input 
      ref="fileInputRef" 
      type="file" 
      accept=".json" 
      style="display: none" 
      @change="handleFileImport"
    />
    
    <!-- Schema差异对比对话框 -->
    <el-dialog 
      v-model="diffDialogVisible" 
      title="Schema 差异对比" 
      width="90%"
      :before-close="handleCloseDiffDialog"
      class="diff-dialog"
    >
      <div class="diff-container">
        <!-- 对比头部 -->
        <div class="diff-header">
          <div class="diff-info">
            <div class="before-info">
              <h4>🔴 修改前</h4>
              <p class="version-info">
                <span class="version">版本: {{ beforeRecord?.version || 'N/A' }}</span>
                <span class="time">时间: {{ formatTime(beforeRecord?.timestamp) }}</span>
              </p>
            </div>
            <div class="after-info">
              <h4>🟢 修改后</h4>
              <p class="version-info">
                <span class="version">版本: {{ afterRecord?.version || 'N/A' }}</span>
                <span class="time">时间: {{ formatTime(afterRecord?.timestamp) }}</span>
              </p>
            </div>
          </div>
        </div>
        
        <!-- 差异摘要 -->
        <div class="diff-summary">
          <h5>📈 修改统计</h5>
          <div class="summary-stats">
            <span class="stat-item added">新增: {{ diffStats.added }} 项</span>
            <span class="stat-item modified">修改: {{ diffStats.modified }} 项</span>
            <span class="stat-item removed">删除: {{ diffStats.removed }} 项</span>
          </div>
        </div>
        
        <!-- 差异内容 -->
        <div class="diff-content">
          <div class="changes-list">
            <div v-for="change in detailedChanges" :key="change.path" class="change-item">
              <div class="change-header">
                <span class="change-type" :class="change.type">
                  {{ getChangeTypeLabel(change.type) }}
                </span>
                <code class="change-path">{{ change.path }}</code>
              </div>
              <div class="change-content">
                <div v-if="change.type === 'modified'" class="change-diff">
                  <div class="old-value">
                    <span class="label">原值:</span>
                    <code>{{ formatValue(change.oldValue) }}</code>
                  </div>
                  <div class="new-value">
                    <span class="label">新值:</span>
                    <code>{{ formatValue(change.newValue) }}</code>
                  </div>
                </div>
                <div v-else-if="change.type === 'added'" class="change-added">
                  <span class="label">新增内容:</span>
                  <code>{{ formatValue(change.value) }}</code>
                </div>
                <div v-else-if="change.type === 'removed'" class="change-removed">
                  <span class="label">删除内容:</span>
                  <code>{{ formatValue(change.value) }}</code>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="diffDialogVisible = false">关闭</el-button>
          <el-button type="primary" @click="handleApplyBeforeSchema">恢复到修改前</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  historyManager: {
    type: Object,
    default: null
  }
});

const emit = defineEmits(['schema-change', 'undo', 'redo']);

const fileInputRef = ref();

// 差异对比相关状态
const diffDialogVisible = ref(false);
const beforeRecord = ref(null);
const afterRecord = ref(null);
const beforeSchema = ref(null);
const afterSchema = ref(null);
const detailedChanges = ref([]);
const diffStats = ref({ added: 0, modified: 0, removed: 0 });

// 历史记录数据
const historyList = computed(() => {
  return props.historyManager ? props.historyManager.getHistoryList() : [];
});
const statistics = computed(() => {
  return props.historyManager ? props.historyManager.getStatistics() : null;
});
const canUndo = computed(() => {
  return props.historyManager ? props.historyManager.canUndo() : false;
});
const canRedo = computed(() => {
  return props.historyManager ? props.historyManager.canRedo() : false;
});

// 显示的历史记录列表（倒序显示，最新的在前面）
const displayHistoryList = computed(() => {
  return [...historyList.value].reverse();
});

// 操作类型标签映射
const operationLabels = {
  'manual': '手动修改',
  'ai-modification': 'AI智能修改',
  'auto-save': '自动保存',
  'import': '导入操作',
  'initial': '初始状态'
};

/**
 * 获取操作类型标签
 */
function getOperationLabel(operation) {
  return operationLabels[operation] || operation;
}

/**
 * 检查是否可以显示差异
 */
function canShowDiff(record) {
  if (!props.historyManager) return false;
  
  const currentIndex = historyList.value.findIndex(r => r.id === record.id);
  // 只有非第一个记录才能显示差异
  return currentIndex > 0;
}

/**
 * 显示Schema差异
 */
function handleShowDiff(record) {
  if (!props.historyManager) return;
  
  const currentIndex = historyList.value.findIndex(r => r.id === record.id);
  if (currentIndex <= 0) return;
  
  // 获取当前记录和前一个记录
  afterRecord.value = historyList.value[currentIndex];
  beforeRecord.value = historyList.value[currentIndex - 1];
  
  afterSchema.value = afterRecord.value.schema;
  beforeSchema.value = beforeRecord.value.schema;
  
  // 计算差异
  calculateDifferences();
  
  // 显示对话框
  diffDialogVisible.value = true;
}

/**
 * 计算Schema差异
 */
function calculateDifferences() {
  if (!beforeSchema.value || !afterSchema.value) return;
  
  const changes = [];
  const stats = { added: 0, modified: 0, removed: 0 };
  
  // 比较节点变化
  const beforeNodes = beforeSchema.value.metadata?.nodes || [];
  const afterNodes = afterSchema.value.metadata?.nodes || [];
  
  const beforeNodeMap = new Map(beforeNodes.map(node => [node.id, node]));
  const afterNodeMap = new Map(afterNodes.map(node => [node.id, node]));
  
  // 检查新增的节点
  afterNodes.forEach(node => {
    if (!beforeNodeMap.has(node.id)) {
      changes.push({
        type: 'added',
        path: `metadata.nodes[${node.id}]`,
        value: node,
        description: `新增节点: ${node.name || node.type}`
      });
      stats.added++;
    }
  });
  
  // 检查删除的节点
  beforeNodes.forEach(node => {
    if (!afterNodeMap.has(node.id)) {
      changes.push({
        type: 'removed',
        path: `metadata.nodes[${node.id}]`,
        value: node,
        description: `删除节点: ${node.name || node.type}`
      });
      stats.removed++;
    }
  });
  
  // 检查修改的节点
  afterNodes.forEach(afterNode => {
    const beforeNode = beforeNodeMap.get(afterNode.id);
    if (beforeNode) {
      const nodeChanges = compareObjects(beforeNode, afterNode, `metadata.nodes[${afterNode.id}]`);
      if (nodeChanges.length > 0) {
        changes.push(...nodeChanges);
        stats.modified++;
      }
    }
  });
  
  // 比较连线变化
  const beforeConnections = beforeSchema.value.metadata?.connections || [];
  const afterConnections = afterSchema.value.metadata?.connections || [];
  
  if (beforeConnections.length !== afterConnections.length) {
    changes.push({
      type: 'modified',
      path: 'metadata.connections',
      oldValue: `${beforeConnections.length} 个连线`,
      newValue: `${afterConnections.length} 个连线`,
      description: '连线数量变化'
    });
    stats.modified++;
  }
  
  // 比较规则链基本信息
  const ruleChainChanges = compareObjects(
    beforeSchema.value.ruleChain, 
    afterSchema.value.ruleChain, 
    'ruleChain'
  );
  if (ruleChainChanges.length > 0) {
    changes.push(...ruleChainChanges);
    stats.modified++;
  }
  
  detailedChanges.value = changes;
  diffStats.value = stats;
}

/**
 * 比较两个对象的差异
 */
function compareObjects(obj1, obj2, basePath = '') {
  const changes = [];
  const allKeys = new Set([...Object.keys(obj1 || {}), ...Object.keys(obj2 || {})]);
  
  for (const key of allKeys) {
    const path = basePath ? `${basePath}.${key}` : key;
    const val1 = obj1?.[key];
    const val2 = obj2?.[key];
    
    if (val1 === undefined && val2 !== undefined) {
      changes.push({
        type: 'added',
        path,
        value: val2,
        description: `新增属性: ${key}`
      });
    } else if (val1 !== undefined && val2 === undefined) {
      changes.push({
        type: 'removed',
        path,
        value: val1,
        description: `删除属性: ${key}`
      });
    } else if (val1 !== val2) {
      if (typeof val1 === 'object' && typeof val2 === 'object' && !Array.isArray(val1) && !Array.isArray(val2)) {
        // 递归比较对象
        changes.push(...compareObjects(val1, val2, path));
      } else {
        changes.push({
          type: 'modified',
          path,
          oldValue: val1,
          newValue: val2,
          description: `修改属性: ${key}`
        });
      }
    }
  }
  
  return changes;
}

/**
 * 获取变化类型标签
 */
function getChangeTypeLabel(type) {
  const labels = {
    'added': '新增',
    'modified': '修改',
    'removed': '删除'
  };
  return labels[type] || type;
}

/**
 * 格式化值显示
 */
function formatValue(value) {
  if (value === null) return 'null';
  if (value === undefined) return 'undefined';
  if (typeof value === 'string') return `"${value}"`;
  if (typeof value === 'object') return JSON.stringify(value, null, 2);
  return String(value);
}

/**
 * 关闭差异对话框
 */
function handleCloseDiffDialog() {
  diffDialogVisible.value = false;
  beforeRecord.value = null;
  afterRecord.value = null;
  beforeSchema.value = null;
  afterSchema.value = null;
  detailedChanges.value = [];
  diffStats.value = { added: 0, modified: 0, removed: 0 };
}

/**
 * 工作流修改前的Schema
 */
function handleApplyBeforeSchema() {
  if (beforeSchema.value) {
    emit('schema-change', beforeSchema.value);
    ElMessage.success('已恢复到修改前的状态');
    diffDialogVisible.value = false;
  }
}

/**
 * 格式化时间
 */
function formatTime(timestamp) {
  const date = new Date(timestamp);
  const now = new Date();
  const diff = now - date;
  
  if (diff < 60000) { // 1分钟内
    return '刚刚';
  } else if (diff < 3600000) { // 1小时内
    return `${Math.floor(diff / 60000)}分钟前`;
  } else if (diff < 86400000) { // 1天内
    return `${Math.floor(diff / 3600000)}小时前`;
  } else {
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
  }
}

/**
 * 处理撤销
 */
function handleUndo() {
  if (!props.historyManager) return;
  
  const result = props.historyManager.undo();
  if (result) {
    emit('undo', result);
    emit('schema-change', result.schema);
    ElMessage.success('已撤销到上一步');
  }
}

/**
 * 处理重做
 */
function handleRedo() {
  if (!props.historyManager) return;
  
  const result = props.historyManager.redo();
  if (result) {
    emit('redo', result);
    emit('schema-change', result.schema);
    ElMessage.success('已重做到下一步');
  }
}

/**
 * 跳转到指定记录
 */
function handleJumpToRecord(recordId) {
  if (!props.historyManager) return;
  
  const result = props.historyManager.jumpToRecord(recordId);
  if (result) {
    emit('schema-change', result.schema);
    ElMessage.success(`已跳转到: ${result.record.description || '指定版本'}`);
  }
}

/**
 * 处理菜单命令
 */
async function handleMenuCommand(command) {
  switch (command) {
    case 'export':
      handleExportHistory();
      break;
    case 'import':
      handleImportHistory();
      break;
    case 'clear':
      await handleClearHistory();
      break;
  }
}

/**
 * 导出历史记录
 */
function handleExportHistory() {
  if (!props.historyManager) return;
  
  try {
    const data = props.historyManager.exportHistory();
    const blob = new Blob([JSON.stringify(data, null, 2)], { 
      type: 'application/json' 
    });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `schema-history-${new Date().getTime()}.json`;
    link.click();
    URL.revokeObjectURL(url);
    ElMessage.success('历史记录已导出');
  } catch (error) {
    ElMessage.error('导出失败: ' + error.message);
  }
}

/**
 * 导入历史记录
 */
function handleImportHistory() {
  fileInputRef.value?.click();
}

/**
 * 处理文件导入
 */
function handleFileImport(event) {
  if (!props.historyManager) return;
  
  const file = event.target.files[0];
  if (!file) return;

  const reader = new FileReader();
  reader.onload = (e) => {
    try {
      const data = JSON.parse(e.target.result);
      const success = props.historyManager.importHistory(data);
      if (success) {
        ElMessage.success('历史记录导入成功');
      } else {
        ElMessage.error('导入失败：文件格式不正确');
      }
    } catch (error) {
      ElMessage.error('导入失败：' + error.message);
    }
  };
  reader.readAsText(file);
  
  // 清空文件输入
  event.target.value = '';
}

/**
 * 清空历史记录
 */
async function handleClearHistory() {
  if (!props.historyManager) return;
  
  try {
    await ElMessageBox.confirm(
      '确定要清空所有历史记录吗？此操作不可恢复。',
      '确认清空',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    );
    
    props.historyManager.clearHistory();
    ElMessage.success('历史记录已清空');
  } catch {
    // 用户取消
  }
}

onMounted(() => {
  // 从本地存储加载历史记录
  if (props.historyManager) {
    props.historyManager.loadFromLocalStorage();
  }
});

onUnmounted(() => {
  // 组件销毁时停止自动保存
  if (props.historyManager) {
    props.historyManager.stopAutoSave();
  }
});
</script>

<style scoped>
.schema-history-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #f8fafc 0%, #e0f2fe 100%);
}

/* 头部样式 */
.history-header {
  padding: 16px;
  border-bottom: 1px solid #e2e8f0;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.title-text {
  font-weight: 600;
  color: #1e293b;
}

.history-count {
  margin-left: 4px;
}

.header-actions {
  display: flex;
  gap: 4px;
}

.action-btn {
  border: 1px solid #d1d5db;
  background: rgba(255, 255, 255, 0.8);
  transition: all 0.2s ease;
}

.action-btn:hover:not(:disabled) {
  background: rgba(59, 130, 246, 0.1);
  border-color: #3b82f6;
  color: #3b82f6;
}

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 统计信息样式 */
.statistics-section {
  padding: 12px 16px;
  border-bottom: 1px solid #e2e8f0;
  background: rgba(255, 255, 255, 0.7);
  display: flex;
  gap: 16px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
}

.stat-label {
  font-size: 11px;
  color: #6b7280;
  margin-bottom: 2px;
}

.stat-value {
  font-size: 16px;
  font-weight: 600;
  color: #1e293b;
}

.ai-count {
  color: #10b981;
}

.manual-count {
  color: #3b82f6;
}

/* 历史记录列表样式 */
.history-list {
  flex: 1;
  overflow: hidden;
}

.list-container {
  height: 100%;
  overflow-y: auto;
  padding: 8px;
}

.history-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px;
  margin-bottom: 8px;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 8px;
  border: 1px solid transparent;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
}

.history-item:hover {
  background: rgba(255, 255, 255, 0.95);
  border-color: #d1d5db;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.history-item.is-current {
  background: linear-gradient(135deg, #dbeafe 0%, #bfdbfe 100%);
  border-color: #3b82f6;
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.2);
}

.history-item.is-ai {
  border-left: 3px solid #10b981;
}

.operation-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f1f5f9;
  flex-shrink: 0;
}

.ai-icon {
  color: #10b981;
}

.manual-icon {
  color: #3b82f6;
}

.auto-save-icon {
  color: #f59e0b;
}

.record-info {
  flex: 1;
  min-width: 0;
}

.record-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.operation-type {
  font-weight: 500;
  color: #1e293b;
  font-size: 13px;
}

.record-time {
  font-size: 11px;
  color: #6b7280;
}

.record-description {
  color: #4b5563;
  font-size: 12px;
  margin-bottom: 6px;
  line-height: 1.4;
}

.record-metadata {
  display: flex;
  gap: 12px;
}

.metadata-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  color: #6b7280;
}

.diff-btn {
  margin-left: 8px;
  font-size: 11px;
  padding: 2px 6px;
  height: auto;
  border-radius: 4px;
}

.diff-btn:hover {
  background: rgba(59, 130, 246, 0.1);
}

.current-indicator {
  position: absolute;
  right: 8px;
  top: 8px;
  color: #3b82f6;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

/* 空状态样式 */
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 32px 16px;
  text-align: center;
}

.empty-icon {
  color: #9ca3af;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 16px;
  font-weight: 500;
  color: #6b7280;
  margin-bottom: 8px;
}

.empty-hint {
  font-size: 14px;
  color: #9ca3af;
}

/* 响应式设计 */
@media (max-width: 480px) {
  .statistics-section {
    padding: 8px 12px;
    gap: 8px;
  }
  
  .history-item {
    padding: 8px;
    gap: 8px;
  }
  
  .operation-icon {
    width: 28px;
    height: 28px;
  }
}

/* 差异对比样式 */
.diff-dialog :deep(.el-dialog) {
  max-width: 95vw;
  max-height: 90vh;
}

.diff-container {
  max-height: 70vh;
  display: flex;
  flex-direction: column;
}

.diff-header {
  padding: 16px 0;
  border-bottom: 1px solid #e5e7eb;
  margin-bottom: 16px;
}

.diff-info {
  display: flex;
  gap: 32px;
}

.before-info,
.after-info {
  flex: 1;
}

.before-info h4,
.after-info h4 {
  margin: 0 0 8px 0;
  font-size: 14px;
  font-weight: 600;
}

.version-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 12px;
  color: #6b7280;
}

.version,
.time {
  margin: 0;
}

.diff-summary {
  margin-bottom: 20px;
  padding: 12px;
  background: #f9fafb;
  border-radius: 6px;
  border: 1px solid #e5e7eb;
}

.diff-summary h5 {
  margin: 0 0 8px 0;
  font-size: 14px;
  font-weight: 600;
  color: #374151;
}

.summary-stats {
  display: flex;
  gap: 16px;
}

.stat-item {
  font-size: 12px;
  font-weight: 500;
  padding: 2px 6px;
  border-radius: 4px;
}

.stat-item.added {
  color: #059669;
  background: #d1fae5;
}

.stat-item.modified {
  color: #d97706;
  background: #fef3c7;
}

.stat-item.removed {
  color: #dc2626;
  background: #fee2e2;
}

.diff-content {
  flex: 1;
  overflow-y: auto;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 16px;
}

.changes-list {
  space-y: 12px;
}

.change-item {
  margin-bottom: 12px;
  padding: 12px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  background: #ffffff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.change-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.change-type {
  font-size: 10px;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 3px;
  text-transform: uppercase;
}

.change-type.added {
  color: #059669;
  background: #d1fae5;
  border: 1px solid #10b981;
}

.change-type.modified {
  color: #d97706;
  background: #fef3c7;
  border: 1px solid #f59e0b;
}

.change-type.removed {
  color: #dc2626;
  background: #fee2e2;
  border: 1px solid #ef4444;
}

.change-path {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 11px;
  color: #6366f1;
  background: #f1f5f9;
  padding: 2px 4px;
  border-radius: 3px;
  border: 1px solid #e2e8f0;
}

.change-content {
  font-size: 12px;
}

.change-diff {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.old-value,
.new-value,
.change-added,
.change-removed {
  padding: 8px;
  border-radius: 4px;
  border: 1px solid;
}

.old-value {
  background: #fef2f2;
  border-color: #fecaca;
}

.new-value {
  background: #f0fdf4;
  border-color: #bbf7d0;
}

.change-added {
  background: #f0fdf4;
  border-color: #bbf7d0;
}

.change-removed {
  background: #fef2f2;
  border-color: #fecaca;
}

.label {
  font-weight: 600;
  color: #374151;
  display: block;
  margin-bottom: 4px;
}

.old-value code,
.new-value code,
.change-added code,
.change-removed code {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 11px;
  background: transparent;
  padding: 0;
  border: none;
  word-break: break-all;
  white-space: pre-wrap;
  color: inherit;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 16px;
}

/* 移动端优化 */
@media (max-width: 768px) {
  .diff-info {
    flex-direction: column;
    gap: 16px;
  }
  
  .change-diff {
    grid-template-columns: 1fr;
  }
  
  .summary-stats {
    flex-direction: column;
    gap: 8px;
  }
  
  .diff-dialog :deep(.el-dialog) {
    width: 95%;
    margin: 5vh auto;
  }
}
</style>