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
    } else {
      ElMessage.error(`Error: ${error.message}`);
    }
    return Promise.reject(error);
  },
);
