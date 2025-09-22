import axios from 'axios';
import { SESSIONSTORAGE_KEYS, getSession, setSession } from '@src/utils/sessionstorage';
import { ElMessage } from 'element-plus';
import router from '@src/router';

export const baseURL = window.config.baseURL + '/api/v1';

export const request = axios.create({
  baseURL,
  timeout: 1000 * 5,
});

export const rawRequest = axios.create({
  baseURL,
  timeout: 1000 * 5,
});

rawRequest.interceptors.request.use(
  (config) => {
    const token = getSession(SESSIONSTORAGE_KEYS.TOKEN);
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

request.interceptors.request.use(
  (config) => {
    const token = getSession(SESSIONSTORAGE_KEYS.TOKEN);
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

request.interceptors.response.use(
  (response) => {
    return response.data;
  },
  (error) => {
    if (error.response) {
      const { status, data } = error.response;
      
      // 处理token过期或未授权的情况
      if (status === 401 || status === 403) {
        // 清除token
        setSession(SESSIONSTORAGE_KEYS.TOKEN, null);
        // 显示提示信息
        ElMessage.error('登录已过期，请重新登录');
        // 重定向到登录页
        router.push('/login');
        return Promise.reject(error);
      }
      
      ElMessage.error(`Error ${status}: ${data.message || JSON.stringify(data)}`);
    } else if (error.code === 'ERR_NETWORK' || error.message.includes('CORS')) {
      // 对于CORS或网络错误，不显示错误信息（特别是登出接口）
      console.warn('网络请求失败（可能是CORS问题）:', error.message);
      // 对于登出接口的CORS错误，不显示用户错误信息
      if (!error.config?.url?.includes('/logout')) {
        ElMessage.error('网络连接失败，请检查网络连接');
      }
    } else {
      ElMessage.error(`Error: ${error.message}`);
    }
    return Promise.reject(error);
  },
);
