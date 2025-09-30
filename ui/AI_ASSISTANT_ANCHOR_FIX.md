# AI智能助手节点输出端点显示问题修复

## 问题描述

用户反馈："AI智能助手修改schema后画布中节点的输出端只显示一个端点，只能通过点击节点详情选择输出端点"

## 问题分析

### 根本原因
AI智能助手在修改schema后，LogicFlow的渲染流程没有正确更新节点的锚点配置，特别是动态生成的输出锚点。具体问题链路：

1. **AI更新流程**：`handleAiSchemaUpdate` → `handelJsonToDesign` → `appDesignRef.value.render()`
2. **render方法**：调用 `rerenderFlowData()` 和 `updateAllNodePropertiesHeight()`
3. **缺失环节**：没有调用关键的锚点更新函数，导致节点的输出端点配置未正确同步

### 影响的节点类型
- **switch节点**：多个case分支的输出端点
- **msg-type-switch节点**：多个消息类型路由的输出端点  
- **endpoint节点**：多个路由路径的输出端点

## 修复方案

### 1. 增强handleAiSchemaUpdate函数
**文件**: `/src/pages/workflow/app-design-new/app-design-new.vue`

```javascript
function handleAiSchemaUpdate(newSchema) {
  // ... 原有逻辑 ...
  
  // 延迟更新画布，确保数据同步并强制更新所有节点锚点
  setTimeout(() => {
    handelJsonToDesign();
    // 额外延迟确保渲染完成后更新锚点
    setTimeout(() => {
      forceUpdateAllNodeAnchors();
    }, 300);
  }, 200);
}
```

### 2. 新增forceUpdateAllNodeAnchors函数
专门处理AI修改schema后的节点锚点强制更新：

```javascript
function forceUpdateAllNodeAnchors() {
  // 遍历所有节点，根据节点类型和配置重新计算锚点
  nodes.forEach((nodeData) => {
    if (nodeData.type === 'switch') {
      // 重新生成switch节点的case输出锚点
    } else if (nodeData.type === 'msg-type-switch') {
      // 重新生成消息类型路由的输出锚点
    } else if (nodeData.type === 'endpoint') {
      // 重新生成endpoint节点的路由输出锚点
    }
  });
}
```

### 3. 优化render方法
**文件**: `/src/pages/workflow/app-design/app-design.vue`

```javascript
function render() {
  if (!flowViewRef.value) return;
  rerenderFlowData();
  flowViewRef.value.updateAllNodePropertiesHeight();
  
  // 延迟更新所有节点的锚点配置，确保渲染完成后正确显示输出端点
  setTimeout(() => {
    // 更新锚点配置和位置
  }, 200);
}
```

### 4. 暴露必要的方法
**文件**: `/src/pages/workflow/app-design/flow-view.vue`

暴露锚点更新相关方法供外部调用：
- `updateNodePropertiesAnchorsById`
- `updateNodePropertiesAnchorsYById` 
- `updateNodePropertiesHeightById`

## 修复效果

### 修复前
- AI修改schema后，节点只显示一个输出端点
- 必须手动点击节点详情才能看到所有输出选项
- 用户体验差，影响工作流程配置效率

### 修复后
- AI修改schema后，节点正确显示所有输出端点
- 动态节点的输出锚点根据配置自动更新
- 用户可以直接在画布上看到并连接所有输出端点
- 实现了真正的"所见即所得"体验

## 技术要点

### 1. 异步渲染处理
使用多层setTimeout确保渲染顺序：
1. 数据更新 (200ms)
2. 基础渲染完成 (300ms)  
3. 锚点强制更新

### 2. 节点类型识别
根据不同节点类型采用不同的锚点更新策略：
- switch: 基于cases配置
- msg-type-switch: 基于routers配置
- endpoint: 基于routers.path配置

### 3. 错误处理
添加完善的try-catch包围，确保单个节点更新失败不影响整体流程

## 相关文件

- `/src/pages/workflow/app-design-new/app-design-new.vue` - 主要修复逻辑
- `/src/pages/workflow/app-design/app-design.vue` - render方法优化
- `/src/pages/workflow/app-design/flow-view.vue` - 方法暴露
- `/src/utils/anchor-update-manager.js` - 新增锚点更新管理器
- `/src/utils/__tests__/anchor-update-manager.test.js` - 锚点更新管理器测试文件

## 测试建议

1. **Switch节点测试**：
   - 创建switch节点，配置多个case分支
   - 让AI修改分支数量或条件
   - 验证输出端点数量正确显示

2. **Endpoint节点测试**：
   - 创建endpoint节点，配置多个路由
   - 让AI添加/删除路由配置
   - 验证路由端点正确更新

3. **复杂流程测试**：
   - 创建包含多种动态节点的复杂流程
   - 让AI进行批量修改
   - 验证所有节点端点显示正常

## 后续优化

1. **封装独立的锚点更新工具类**：已创建[AnchorUpdateManager](file:///Users/yangcong1/learnProject/rulego-server/ui/src/utils/anchor-update-manager.js#L12-L388)统一管理锚点更新逻辑
2. **添加更多节点类型的锚点更新支持**：已支持switch、msg-type-switch、endpoint、for、fork等多种节点类型
3. **优化更新性能，减少不必要的DOM操作**：通过缓存机制和批量更新策略优化性能