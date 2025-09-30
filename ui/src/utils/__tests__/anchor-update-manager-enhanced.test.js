/**
 * 锚点更新管理器增强测试文件
 */
import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { AnchorUpdateManager } from '@/utils/anchor-update-manager.js';

// 模拟LogicFlow实例
const mockLf = {
  on: vi.fn(),
  off: vi.fn(),
  getNodeModelById: vi.fn(),
  setProperties: vi.fn(),
  getGraphData: vi.fn()
};

describe('AnchorUpdateManager Enhanced', () => {
  let anchorUpdateManager;

  beforeEach(() => {
    anchorUpdateManager = new AnchorUpdateManager();
    vi.clearAllMocks();
  });

  afterEach(() => {
    anchorUpdateManager.destroy();
  });

  it('应该能够设置LogicFlow实例并绑定事件监听器', () => {
    anchorUpdateManager.setLogicFlowInstance(mockLf);
    
    expect(anchorUpdateManager.lfInstance).toBe(mockLf);
    expect(mockLf.on).toHaveBeenCalledWith('node:add', expect.any(Function));
    expect(mockLf.on).toHaveBeenCalledWith('node:delete', expect.any(Function));
    expect(mockLf.on).toHaveBeenCalledWith('node:data-change', expect.any(Function));
  });

  it('应该能够处理节点新增事件', () => {
    // 模拟锚点更新管理器的方法
    anchorUpdateManager.updateSingleNodeAnchor = vi.fn();
    
    anchorUpdateManager.setLogicFlowInstance(mockLf);
    
    const mockNodeData = {
      id: 'test-node',
      type: 'switch',
      properties: {
        formData: {
          cases: [{ label: 'CASE 1' }]
        }
      }
    };
    
    // 模拟getNodeModelById返回节点模型
    mockLf.getNodeModelById.mockReturnValue({
      updateAnchorPosition: vi.fn()
    });
    
    // 触发节点新增事件处理
    anchorUpdateManager.handleNodeAdd({ data: mockNodeData });
    
    // 等待异步操作完成
    setTimeout(() => {
      expect(mockLf.getNodeModelById).toHaveBeenCalledWith('test-node');
    }, 150);
  });

  it('应该能够处理节点删除事件', () => {
    anchorUpdateManager.setLogicFlowInstance(mockLf);
    
    const mockNodeData = {
      id: 'test-node'
    };
    
    // 先添加到缓存
    const properties = { type: 'switch', formData: { cases: [] } };
    anchorUpdateManager.needsUpdate('test-node', properties);
    expect(anchorUpdateManager.performanceCache.size).toBe(1);
    
    // 触发节点删除事件处理
    anchorUpdateManager.handleNodeDelete({ data: mockNodeData });
    
    // 检查缓存是否被清理
    expect(anchorUpdateManager.performanceCache.size).toBe(0);
  });

  it('应该能够处理节点数据变更事件', () => {
    // 模拟锚点更新管理器的方法
    anchorUpdateManager.updateSingleNodeAnchor = vi.fn();
    anchorUpdateManager.needsUpdate = vi.fn().mockReturnValue(true);
    
    anchorUpdateManager.setLogicFlowInstance(mockLf);
    
    const mockNodeData = {
      id: 'test-node',
      properties: {
        type: 'switch',
        formData: {
          cases: [{ label: 'CASE 1' }, { label: 'CASE 2' }]
        }
      }
    };
    
    // 触发节点数据变更事件处理
    anchorUpdateManager.handleNodeDataChange({ data: mockNodeData });
    
    // 检查是否正确处理了数据变更
    setTimeout(() => {
      // 应该调用更新锚点的方法
      expect(anchorUpdateManager.updateSingleNodeAnchor).toHaveBeenCalled();
    }, 100);
  });

  it('应该在销毁时清理LogicFlow事件监听器', () => {
    anchorUpdateManager.setLogicFlowInstance(mockLf);
    anchorUpdateManager.destroy();
    
    expect(mockLf.off).toHaveBeenCalledWith('node:add', expect.any(Function));
    expect(mockLf.off).toHaveBeenCalledWith('node:delete', expect.any(Function));
    expect(mockLf.off).toHaveBeenCalledWith('node:data-change', expect.any(Function));
    expect(anchorUpdateManager.lfInstance).toBeNull();
  });
});