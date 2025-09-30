/**
 * 锚点更新管理器
 * 负责统一管理所有节点类型的锚点更新逻辑
 */
import { cloneDeep } from 'lodash-es';

export class AnchorUpdateManager {
  constructor() {
    this.updateStrategies = new Map();
    this.performanceCache = new Map();
    this.lastUpdateTime = 0;
    this.updateQueue = [];
    this.isUpdating = false;
    this.lfInstance = null;
    
    // 注册默认的更新策略
    this.registerDefaultStrategies();
    
    // 绑定事件监听器
    this.bindEventListeners();
  }

  /**
   * 绑定LogicFlow事件监听器
   */
  bindEventListeners() {
    // 这里将在setLogicFlowInstance中绑定实际的事件监听器
  }

  /**
   * 设置LogicFlow实例并绑定事件监听器
   * @param {object} lf LogicFlow实例
   */
  setLogicFlowInstance(lf) {
    if (!lf) return;
    
    this.lfInstance = lf;
    
    // 绑定节点新增事件
    lf.on('node:add', this.handleNodeAdd.bind(this));
    
    // 绑定节点删除事件
    lf.on('node:delete', this.handleNodeDelete.bind(this));
    
    // 绑定节点数据更新事件
    lf.on('node:data-change', this.handleNodeDataChange.bind(this));
  }

  /**
   * 处理节点新增事件
   * @param {object} event 事件对象
   */
  handleNodeAdd(event) {
    console.log('AnchorUpdateManager: 检测到节点新增', event.data);
    
    // 延迟更新锚点，确保节点完全渲染
    setTimeout(() => {
      this.updateSingleNodeAnchor(this.lfInstance, event.data);
      
      // 再次延迟更新锚点位置
      setTimeout(() => {
        try {
          const nodeModel = this.lfInstance.getNodeModelById(event.data.id);
          if (nodeModel && typeof nodeModel.updateAnchorPosition === 'function') {
            nodeModel.updateAnchorPosition();
          }
        } catch (error) {
          console.warn('更新节点锚点位置时出错:', error);
        }
      }, 50);
    }, 100);
  }

  /**
   * 处理节点删除事件
   * @param {object} event 事件对象
   */
  handleNodeDelete(event) {
    console.log('AnchorUpdateManager: 检测到节点删除', event.data);
    
    // 从缓存中移除该节点
    this.clearCache(event.data.id);
  }

  /**
   * 处理节点数据变更事件
   * @param {object} event 事件对象
   */
  handleNodeDataChange(event) {
    console.log('AnchorUpdateManager: 检测到节点数据变更', event.data);
    
    // 检查是否需要更新锚点
    if (this.needsUpdate(event.data.id, event.data.properties)) {
      setTimeout(() => {
        this.updateSingleNodeAnchor(this.lfInstance, event.data);
      }, 50);
    }
  }

  /**
   * 注册默认的节点类型更新策略
   */
  registerDefaultStrategies() {
    // Switch节点更新策略
    this.registerStrategy('switch', (lf, nodeId, properties) => {
      const cases = properties?.formData?.cases || [];
      const staticAnchors = properties?.anchors?.filter(anchor => anchor.isStatic) || [];
      const newAnchors = [...staticAnchors];
      
      cases.forEach((caseItem, index) => {
        const anchorId = caseItem.label || `CASE ${index + 1}`;
        newAnchors.push({
          id: anchorId,
          isStatic: false,
          name: anchorId,
          top: 0,
          type: 'output'
        });
      });
      
      return newAnchors;
    });

    // MsgTypeSwitch节点更新策略
    this.registerStrategy('msg-type-switch', (lf, nodeId, properties) => {
      const routers = properties?.formData?.routers || [];
      const staticAnchors = properties?.anchors?.filter(anchor => anchor.isStatic) || [];
      const newAnchors = [...staticAnchors];
      
      routers.forEach((router) => {
        newAnchors.push({
          id: router,
          isStatic: false,
          name: router,
          top: 0,
          type: 'output'
        });
      });
      
      return newAnchors;
    });

    // Endpoint节点更新策略
    this.registerStrategy('endpoint', (lf, nodeId, properties) => {
      const routers = properties?.formData?.routers || [];
      const staticAnchors = properties?.anchors?.filter(anchor => anchor.isStatic) || [];
      const newAnchors = [...staticAnchors];
      
      routers.forEach((router) => {
        newAnchors.push({
          id: router.path,
          isStatic: false,
          name: router.path,
          top: 0,
          type: 'output'
        });
      });
      
      return newAnchors;
    });

    // For循环节点更新策略
    this.registerStrategy('for', (lf, nodeId, properties) => {
      const staticAnchors = properties?.anchors?.filter(anchor => anchor.isStatic) || [];
      const newAnchors = [...staticAnchors];
      
      // For节点通常有True和False两个输出
      if (!newAnchors.find(anchor => anchor.id === 'True')) {
        newAnchors.push({
          id: 'True',
          isStatic: false,
          name: 'True',
          top: 0,
          type: 'output'
        });
      }
      
      if (!newAnchors.find(anchor => anchor.id === 'False')) {
        newAnchors.push({
          id: 'False',
          isStatic: false,
          name: 'False',
          top: 0,
          type: 'output'
        });
      }
      
      return newAnchors;
    });

    // Fork分支节点更新策略
    this.registerStrategy('fork', (lf, nodeId, properties) => {
      const branches = properties?.formData?.branches || [];
      const staticAnchors = properties?.anchors?.filter(anchor => anchor.isStatic) || [];
      const newAnchors = [...staticAnchors];
      
      branches.forEach((branch, index) => {
        const anchorId = branch.name || `Branch ${index + 1}`;
        newAnchors.push({
          id: anchorId,
          isStatic: false,
          name: anchorId,
          top: 0,
          type: 'output'
        });
      });
      
      return newAnchors;
    });
  }

  /**
   * 注册节点类型的更新策略
   * @param {string} nodeType 节点类型
   * @param {function} strategy 更新策略函数
   */
  registerStrategy(nodeType, strategy) {
    this.updateStrategies.set(nodeType, strategy);
  }

  /**
   * 获取节点类型的更新策略
   * @param {string} nodeType 节点类型
   * @returns {function|null} 更新策略函数
   */
  getStrategy(nodeType) {
    return this.updateStrategies.get(nodeType) || null;
  }

  /**
   * 生成缓存键
   * @param {string} nodeId 节点ID
   * @param {object} properties 节点属性
   * @returns {string} 缓存键
   */
  generateCacheKey(nodeId, properties) {
    const relevantData = {
      type: properties?.type,
      cases: properties?.formData?.cases,
      routers: properties?.formData?.routers,
      branches: properties?.formData?.branches
    };
    return `${nodeId}_${JSON.stringify(relevantData)}`;
  }

  /**
   * 检查是否需要更新（基于缓存）
   * @param {string} nodeId 节点ID
   * @param {object} properties 节点属性
   * @returns {boolean} 是否需要更新
   */
  needsUpdate(nodeId, properties) {
    const cacheKey = this.generateCacheKey(nodeId, properties);
    const cachedKey = this.performanceCache.get(nodeId);
    
    if (cachedKey === cacheKey) {
      return false; // 数据未变化，无需更新
    }
    
    this.performanceCache.set(nodeId, cacheKey);
    return true;
  }

  /**
   * 批量更新节点锚点（性能优化版本）
   * @param {object} lf LogicFlow实例
   * @param {array} nodes 节点数组
   * @param {object} options 选项
   */
  async batchUpdateAnchors(lf, nodes, options = {}) {
    const {
      maxBatchSize = 5, // 每批最大处理节点数
      batchDelay = 50,   // 批次间延迟
      forceUpdate = false // 是否强制更新（忽略缓存）
    } = options;

    if (!lf || !nodes || nodes.length === 0) {
      console.warn('AnchorUpdateManager: 无效的参数');
      return;
    }

    // 防止并发更新
    if (this.isUpdating) {
      console.log('AnchorUpdateManager: 正在更新中，跳过本次请求');
      return;
    }

    this.isUpdating = true;
    const startTime = Date.now();
    
    try {
      console.log(`AnchorUpdateManager: 开始批量更新 ${nodes.length} 个节点的锚点`);
      
      // 过滤需要更新的节点
      const nodesToUpdate = nodes.filter(nodeData => {
        const nodeModel = lf.getNodeModelById(nodeData.id);
        if (!nodeModel) return false;
        
        const properties = nodeModel.getProperties();
        const strategy = this.getStrategy(nodeData.type);
        
        return strategy && (forceUpdate || this.needsUpdate(nodeData.id, properties));
      });

      console.log(`AnchorUpdateManager: 过滤后需要更新 ${nodesToUpdate.length} 个节点`);

      // 分批处理
      for (let i = 0; i < nodesToUpdate.length; i += maxBatchSize) {
        const batch = nodesToUpdate.slice(i, i + maxBatchSize);
        
        await this.processBatch(lf, batch);
        
        // 批次间延迟，避免阻塞UI
        if (i + maxBatchSize < nodesToUpdate.length) {
          await new Promise(resolve => setTimeout(resolve, batchDelay));
        }
      }

      const endTime = Date.now();
      console.log(`AnchorUpdateManager: 批量更新完成，耗时 ${endTime - startTime}ms`);
      
    } catch (error) {
      console.error('AnchorUpdateManager: 批量更新失败:', error);
    } finally {
      this.isUpdating = false;
    }
  }

  /**
   * 处理一批节点
   * @param {object} lf LogicFlow实例
   * @param {array} batch 节点批次
   */
  async processBatch(lf, batch) {
    const updatePromises = batch.map(nodeData => {
      return new Promise((resolve) => {
        try {
          const success = this.updateSingleNodeAnchor(lf, nodeData);
          resolve({ nodeId: nodeData.id, success });
        } catch (error) {
          console.error(`更新节点 ${nodeData.id} 锚点失败:`, error);
          resolve({ nodeId: nodeData.id, success: false, error });
        }
      });
    });

    const results = await Promise.all(updatePromises);
    
    // 统计结果
    const successCount = results.filter(r => r.success).length;
    const failCount = results.length - successCount;
    
    if (failCount > 0) {
      console.warn(`批次处理完成: ${successCount} 成功, ${failCount} 失败`);
    }
  }

  /**
   * 更新单个节点的锚点
   * @param {object} lf LogicFlow实例
   * @param {object} nodeData 节点数据
   * @returns {boolean} 是否更新成功
   */
  updateSingleNodeAnchor(lf, nodeData) {
    try {
      const nodeModel = lf.getNodeModelById(nodeData.id);
      if (!nodeModel) {
        console.warn(`节点 ${nodeData.id} 模型不存在`);
        return false;
      }

      const properties = nodeModel.getProperties();
      const strategy = this.getStrategy(nodeData.type);
      
      if (!strategy) {
        // 无更新策略的节点跳过
        return true;
      }

      // 执行更新策略
      const newAnchors = strategy(lf, nodeData.id, properties);
      
      if (!newAnchors || !Array.isArray(newAnchors)) {
        console.warn(`节点 ${nodeData.id} 更新策略返回无效的锚点配置`);
        return false;
      }

      // 检查锚点是否实际发生变化
      const currentAnchors = properties?.anchors || [];
      if (this.anchorsEqual(currentAnchors, newAnchors)) {
        return true; // 无变化，无需更新
      }

      // 更新锚点配置
      lf.setProperties(nodeData.id, {
        ...properties,
        anchors: newAnchors
      });

      console.log(`节点 ${nodeData.id} (${nodeData.type}) 锚点更新成功，锚点数量: ${newAnchors.length}`);
      return true;
      
    } catch (error) {
      console.error(`更新节点 ${nodeData.id} 锚点时出错:`, error);
      return false;
    }
  }

  /**
   * 比较两个锚点数组是否相等
   * @param {array} anchors1 锚点数组1
   * @param {array} anchors2 锚点数组2
   * @returns {boolean} 是否相等
   */
  anchorsEqual(anchors1, anchors2) {
    if (!anchors1 || !anchors2) return false;
    if (anchors1.length !== anchors2.length) return false;
    
    for (let i = 0; i < anchors1.length; i++) {
      const a1 = anchors1[i];
      const a2 = anchors2[i];
      
      if (a1.id !== a2.id || a1.type !== a2.type || a1.name !== a2.name) {
        return false;
      }
    }
    
    return true;
  }

  /**
   * 清理缓存
   * @param {string} nodeId 可选，指定节点ID清理特定缓存
   */
  clearCache(nodeId = null) {
    if (nodeId) {
      this.performanceCache.delete(nodeId);
    } else {
      this.performanceCache.clear();
    }
  }

  /**
   * 获取缓存统计
   * @returns {object} 缓存统计信息
   */
  getCacheStats() {
    return {
      size: this.performanceCache.size,
      supportedNodeTypes: Array.from(this.updateStrategies.keys()),
      isUpdating: this.isUpdating
    };
  }

  /**
   * 销毁管理器，清理资源
   */
  destroy() {
    // 清理LogicFlow事件监听器
    if (this.lfInstance) {
      this.lfInstance.off('node:add', this.handleNodeAdd.bind(this));
      this.lfInstance.off('node:delete', this.handleNodeDelete.bind(this));
      this.lfInstance.off('node:data-change', this.handleNodeDataChange.bind(this));
    }
    
    this.updateStrategies.clear();
    this.performanceCache.clear();
    this.updateQueue = [];
    this.isUpdating = false;
    this.lfInstance = null;
  }
}

// 创建单例实例
export const anchorUpdateManager = new AnchorUpdateManager();

// 默认导出
export default anchorUpdateManager;