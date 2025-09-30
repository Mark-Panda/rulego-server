/**
 * 锚点更新管理器测试文件
 */
import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { AnchorUpdateManager } from '@/utils/anchor-update-manager.js';

// 模拟LogicFlow实例
const mockLf = {
  getNodeModelById: vi.fn(),
  setProperties: vi.fn(),
  getGraphData: vi.fn()
};

// 测试数据
const mockNodes = [
  {
    id: 'node1',
    type: 'switch',
    properties: {
      formData: {
        cases: [
          { label: 'CASE 1' },
          { label: 'CASE 2' }
        ]
      },
      anchors: [
        { id: 'input', type: 'input', isStatic: true }
      ]
    }
  },
  {
    id: 'node2',
    type: 'msg-type-switch',
    properties: {
      formData: {
        routers: ['type1', 'type2']
      },
      anchors: [
        { id: 'input', type: 'input', isStatic: true }
      ]
    }
  }
];

describe('AnchorUpdateManager', () => {
  let anchorUpdateManager;

  beforeEach(() => {
    anchorUpdateManager = new AnchorUpdateManager();
    vi.clearAllMocks();
  });

  afterEach(() => {
    anchorUpdateManager.destroy();
  });

  it('应该正确创建实例', () => {
    expect(anchorUpdateManager).toBeInstanceOf(AnchorUpdateManager);
    expect(anchorUpdateManager.updateStrategies.size).toBeGreaterThan(0);
  });

  it('应该注册默认的更新策略', () => {
    expect(anchorUpdateManager.getStrategy('switch')).toBeTypeOf('function');
    expect(anchorUpdateManager.getStrategy('msg-type-switch')).toBeTypeOf('function');
    expect(anchorUpdateManager.getStrategy('endpoint')).toBeTypeOf('function');
    expect(anchorUpdateManager.getStrategy('for')).toBeTypeOf('function');
    expect(anchorUpdateManager.getStrategy('fork')).toBeTypeOf('function');
  });

  it('应该能够注册新的更新策略', () => {
    const customStrategy = vi.fn();
    anchorUpdateManager.registerStrategy('custom-node', customStrategy);
    
    const strategy = anchorUpdateManager.getStrategy('custom-node');
    expect(strategy).toBe(customStrategy);
  });

  it('应该能够生成缓存键', () => {
    const properties = {
      type: 'switch',
      formData: {
        cases: [{ label: 'Test' }]
      }
    };
    
    const cacheKey = anchorUpdateManager.generateCacheKey('node1', properties);
    expect(cacheKey).toContain('node1');
    expect(cacheKey).toContain('switch');
  });

  it('应该能够检查是否需要更新', () => {
    const properties = {
      type: 'switch',
      formData: {
        cases: [{ label: 'Test' }]
      }
    };
    
    // 第一次检查应该需要更新
    expect(anchorUpdateManager.needsUpdate('node1', properties)).toBe(true);
    
    // 第二次检查相同数据应该不需要更新
    expect(anchorUpdateManager.needsUpdate('node1', properties)).toBe(false);
    
    // 修改数据后应该需要更新
    const newProperties = {
      type: 'switch',
      formData: {
        cases: [{ label: 'Test2' }]
      }
    };
    expect(anchorUpdateManager.needsUpdate('node1', newProperties)).toBe(true);
  });

  it('应该能够比较锚点数组是否相等', () => {
    const anchors1 = [
      { id: 'input', type: 'input', name: 'Input' },
      { id: 'output1', type: 'output', name: 'Output1' }
    ];
    
    const anchors2 = [
      { id: 'input', type: 'input', name: 'Input' },
      { id: 'output1', type: 'output', name: 'Output1' }
    ];
    
    const anchors3 = [
      { id: 'input', type: 'input', name: 'Input' },
      { id: 'output2', type: 'output', name: 'Output2' }
    ];
    
    expect(anchorUpdateManager.anchorsEqual(anchors1, anchors2)).toBe(true);
    expect(anchorUpdateManager.anchorsEqual(anchors1, anchors3)).toBe(false);
  });

  it('应该能够清理缓存', () => {
    const properties = {
      type: 'switch',
      formData: {
        cases: [{ label: 'Test' }]
      }
    };
    
    // 触发缓存
    anchorUpdateManager.needsUpdate('node1', properties);
    expect(anchorUpdateManager.performanceCache.size).toBe(1);
    
    // 清理特定节点缓存
    anchorUpdateManager.clearCache('node1');
    expect(anchorUpdateManager.performanceCache.size).toBe(0);
  });

  it('应该能够获取缓存统计', () => {
    const stats = anchorUpdateManager.getCacheStats();
    expect(stats).toHaveProperty('size');
    expect(stats).toHaveProperty('supportedNodeTypes');
    expect(stats).toHaveProperty('isUpdating');
  });
});