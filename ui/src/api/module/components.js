import { request } from '@src/utils/request';

/**
 * 获取所有组件
 * @returns {Promise} 返回所有组件列表
 */
export function getComponents() {
  return request.get('/components');
}

/**
 * 获取已安装组件列表
 * @param {Object} params 查询参数
 * @returns {Promise} 返回已安装组件列表
 */
export function getInstalledComponents(params) {
  return request.get('/components', { params });
}

/**
 * 获取组件市场列表
 * @param {Object} params 查询参数
 * @returns {Promise} 返回组件市场列表
 */
export function getMarketComponents(params) {
  return request.get('/marketplace/components', { params });
}

/**
 * 安装组件
 * @param {String} id 组件ID
 * @param {Object} data 组件数据
 * @returns {Promise} 返回安装结果
 */
export function installComponent(id, data) {
  return request.post(`/dynamic-components/${id}`, data);
}

/**
 * 卸载组件
 * @param {String} id 组件ID
 * @returns {Promise} 返回卸载结果
 */
export function uninstallComponent(id) {
  return request.delete(`/dynamic-components/${id}`);
}

/**
 * 获取组件详情
 * @param {String} id 组件ID
 * @returns {Promise} 返回组件详情
 */
export function getComponentDetail(id) {
  return request.get(`/components/${id}`);
}

/**
 * 获取自定义组件详情
 * @param {String} id 组件ID
 * @returns {Promise} 返回自定义组件详情
 */
export function getDynamicComponentDetail(id) {
  return request.get(`/dynamic-components/${id}`);
}

/**
 * 获取组件使用规则分页列表
 * @param {Object} params 查询参数
 * @param {Number} params.page 页码
 * @param {Number} params.size 每页大小
 * @param {String} params.componentType 组件类型
 * @param {Boolean} params.disabled 是否禁用
 * @param {String} params.keywords 关键词
 * @returns {Promise} 返回组件使用规则分页列表
 */
export function getComponentUseRulePage(params) {
  return request.get('/componentUseRule/page', { params });
}

/**
 * 创建组件使用规则
 * @param {Object} data 规则数据
 * @param {String} data.componentName 组件名称
 * @param {String} data.componentType 组件类型
 * @param {Boolean} data.disabled 是否禁用
 * @param {String} data.useDesc 使用描述
 * @param {String} data.useRuleDesc 使用规则描述
 * @returns {Promise} 返回创建结果
 */
export function createComponentUseRule(data) {
  return request.post('/componentUseRule/create', data);
}

/**
 * 更新组件使用规则
 * @param {Object} data 规则数据
 * @param {String|Number} data.id 规则ID
 * @param {String} data.componentName 组件名称
 * @param {String} data.componentType 组件类型
 * @param {Boolean} data.disabled 是否禁用
 * @param {String} data.useDesc 使用描述
 * @param {String} data.useRuleDesc 使用规则描述
 * @returns {Promise} 返回更新结果
 */
export function updateComponentUseRule(data) {
  return request.post('/componentUseRule/update', data);
}

/**
 * 删除组件使用规则
 * @param {Object} params 删除参数
 * @param {String|Number} params.id 规则ID
 * @returns {Promise} 返回删除结果
 */
export function deleteComponentUseRule(params) {
  return request.delete('/componentUseRule/delete', { params });
}
