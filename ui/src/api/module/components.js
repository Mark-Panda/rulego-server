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
  return request.get('marketplace/components', { params });
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
