// Schema修改历史记录管理器
import { ref, reactive } from 'vue';
import { nanoid } from 'nanoid';
import { cloneDeep } from 'lodash-es';

export class SchemaHistoryManager {
  constructor() {
    // 历史记录数组
    this.history = reactive([]);
    // 当前位置指针
    this.currentIndex = ref(-1);
    // 最大历史记录数量
    this.maxHistorySize = 50;
    // 自动保存间隔（毫秒）
    this.autoSaveInterval = 3000;
    // 自动保存定时器
    this.autoSaveTimer = null;
  }

  /**
   * 添加新的历史记录
   * @param {Object} schema - schema数据
   * @param {string} operation - 操作类型
   * @param {string} description - 操作描述
   * @param {Object} metadata - 元数据
   */
  addRecord(schema, operation = 'manual', description = '', metadata = {}) {
    const record = {
      id: nanoid(),
      timestamp: Date.now(),
      schema: cloneDeep(schema),
      operation,
      description,
      metadata: {
        ...metadata,
        nodeCount: this._getNodeCount(schema),
        edgeCount: this._getEdgeCount(schema),
      },
      version: this.history.length + 1
    };

    // 如果当前不在最新位置，删除当前位置之后的所有记录
    if (this.currentIndex.value < this.history.length - 1) {
      this.history.splice(this.currentIndex.value + 1);
    }

    // 添加新记录
    this.history.push(record);
    this.currentIndex.value = this.history.length - 1;

    // 限制历史记录数量
    if (this.history.length > this.maxHistorySize) {
      this.history.shift();
      this.currentIndex.value--;
    }

    // 保存到本地存储
    this._saveToLocalStorage();

    return record;
  }

  /**
   * 撤销操作
   * @returns {Object|null} 上一个schema或null
   */
  undo() {
    if (!this.canUndo()) {
      return null;
    }

    this.currentIndex.value--;
    const record = this.history[this.currentIndex.value];
    
    return {
      schema: cloneDeep(record.schema),
      record
    };
  }

  /**
   * 重做操作
   * @returns {Object|null} 下一个schema或null
   */
  redo() {
    if (!this.canRedo()) {
      return null;
    }

    this.currentIndex.value++;
    const record = this.history[this.currentIndex.value];
    
    return {
      schema: cloneDeep(record.schema),
      record
    };
  }

  /**
   * 检查是否可以撤销
   */
  canUndo() {
    return this.currentIndex.value > 0;
  }

  /**
   * 检查是否可以重做
   */
  canRedo() {
    return this.currentIndex.value < this.history.length - 1;
  }

  /**
   * 获取当前记录
   */
  getCurrentRecord() {
    if (this.currentIndex.value >= 0 && this.currentIndex.value < this.history.length) {
      return this.history[this.currentIndex.value];
    }
    return null;
  }

  /**
   * 获取历史记录列表（用于UI显示）
   */
  getHistoryList() {
    return this.history.map((record, index) => ({
      ...record,
      isCurrent: index === this.currentIndex.value,
      isAI: record.operation === 'ai-modification'
    }));
  }

  /**
   * 跳转到指定历史记录
   * @param {string} recordId - 记录ID
   */
  jumpToRecord(recordId) {
    const index = this.history.findIndex(record => record.id === recordId);
    if (index !== -1) {
      this.currentIndex.value = index;
      const record = this.history[index];
      return {
        schema: cloneDeep(record.schema),
        record
      };
    }
    return null;
  }

  /**
   * 清空历史记录
   */
  clearHistory() {
    this.history.splice(0);
    this.currentIndex.value = -1;
    this._saveToLocalStorage();
  }

  /**
   * 开始自动保存
   * @param {Function} getCurrentSchema - 获取当前schema的函数
   */
  startAutoSave(getCurrentSchema) {
    this.stopAutoSave();
    
    this.autoSaveTimer = setInterval(() => {
      try {
        const currentSchema = getCurrentSchema();
        if (currentSchema && this._hasChanges(currentSchema)) {
          this.addRecord(
            currentSchema, 
            'auto-save', 
            '自动保存',
            { autoSave: true }
          );
        }
      } catch (error) {
        console.warn('自动保存失败:', error);
      }
    }, this.autoSaveInterval);
  }

  /**
   * 停止自动保存
   */
  stopAutoSave() {
    if (this.autoSaveTimer) {
      clearInterval(this.autoSaveTimer);
      this.autoSaveTimer = null;
    }
  }

  /**
   * 导出历史记录
   */
  exportHistory() {
    return {
      version: '1.0',
      exportTime: Date.now(),
      history: this.history,
      currentIndex: this.currentIndex.value
    };
  }

  /**
   * 导入历史记录
   * @param {Object} data - 导入的数据
   */
  importHistory(data) {
    if (data && data.history && Array.isArray(data.history)) {
      this.history.splice(0, this.history.length, ...data.history);
      this.currentIndex.value = Math.min(
        data.currentIndex || 0, 
        this.history.length - 1
      );
      this._saveToLocalStorage();
      return true;
    }
    return false;
  }

  /**
   * 获取统计信息
   */
  getStatistics() {
    const aiRecords = this.history.filter(r => r.operation === 'ai-modification');
    const manualRecords = this.history.filter(r => r.operation === 'manual');
    const autoSaveRecords = this.history.filter(r => r.operation === 'auto-save');

    return {
      totalRecords: this.history.length,
      aiModifications: aiRecords.length,
      manualModifications: manualRecords.length,
      autoSaves: autoSaveRecords.length,
      currentPosition: this.currentIndex.value + 1,
      canUndo: this.canUndo(),
      canRedo: this.canRedo(),
      oldestRecord: this.history[0],
      newestRecord: this.history[this.history.length - 1]
    };
  }

  // 私有方法

  /**
   * 获取节点数量
   */
  _getNodeCount(schema) {
    return schema?.ruleChain?.nodes?.length || 0;
  }

  /**
   * 获取连线数量
   */
  _getEdgeCount(schema) {
    return schema?.ruleChain?.connections?.length || 0;
  }

  /**
   * 检查是否有变化
   */
  _hasChanges(currentSchema) {
    const currentRecord = this.getCurrentRecord();
    if (!currentRecord) return true;

    // 简单比较节点和连线数量
    const currentNodeCount = this._getNodeCount(currentSchema);
    const currentEdgeCount = this._getEdgeCount(currentSchema);
    const lastNodeCount = currentRecord.metadata.nodeCount;
    const lastEdgeCount = currentRecord.metadata.edgeCount;

    return currentNodeCount !== lastNodeCount || currentEdgeCount !== lastEdgeCount;
  }

  /**
   * 保存到本地存储
   */
  _saveToLocalStorage() {
    try {
      const data = {
        history: this.history.slice(-20), // 只保存最近20条记录到本地存储
        currentIndex: Math.max(0, this.currentIndex.value - (this.history.length - 20))
      };
      localStorage.setItem('schema-history', JSON.stringify(data));
    } catch (error) {
      console.warn('保存历史记录到本地存储失败:', error);
    }
  }

  /**
   * 从本地存储加载
   */
  loadFromLocalStorage() {
    try {
      const data = localStorage.getItem('schema-history');
      if (data) {
        const parsed = JSON.parse(data);
        if (parsed.history && Array.isArray(parsed.history)) {
          this.history.splice(0, this.history.length, ...parsed.history);
          this.currentIndex.value = parsed.currentIndex || 0;
          return true;
        }
      }
    } catch (error) {
      console.warn('从本地存储加载历史记录失败:', error);
    }
    return false;
  }
}

// 创建全局实例
export const schemaHistoryManager = new SchemaHistoryManager();