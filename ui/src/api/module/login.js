import { request, rawRequest } from '@src/utils/request';

export function login(data) {
  return request.post('/login', data);
}

export function logout() {
  // 使用rawRequest避免触发响应拦截器的错误处理
  return rawRequest.post('/logout');
}
