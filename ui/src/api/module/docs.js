import { request } from '@src/utils/request';

// 获取业务文档列表
export function getDocs(params) {
  return request({
    url: '/api/docs',
    method: 'get',
    params
  });
}

// 获取业务文档详情
export function getDoc(id) {
  return request({
    url: `/api/docs/${id}`,
    method: 'get'
  });
}

// 创建业务文档
export function createDoc(data) {
  return request({
    url: '/api/docs',
    method: 'post',
    data
  });
}

// 更新业务文档
export function updateDoc(id, data) {
  return request({
    url: `/api/docs/${id}`,
    method: 'put',
    data
  });
}

// 删除业务文档
export function deleteDoc(id) {
  return request({
    url: `/api/docs/${id}`,
    method: 'delete'
  });
}