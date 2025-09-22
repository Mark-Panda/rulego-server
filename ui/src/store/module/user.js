import { ref, computed } from 'vue';
import { defineStore } from 'pinia';
import { SESSIONSTORAGE_KEYS, getSession, setSession, clearAllSession } from '@src/utils/sessionstorage';
import { logout as logoutApi } from '@src/api/module/login';
import router from '@src/router';

const useUserStore = defineStore('user', () => {
  const token = ref(getSession(SESSIONSTORAGE_KEYS.TOKEN));
  const isLoggedIn = computed(() => !!token.value);

  function setToken(newToken) {
    token.value = newToken;
    setSession(SESSIONSTORAGE_KEYS.TOKEN, newToken);
  }

  async function logout() {
    try {
      // 尝试调用后端登出接口
      await logoutApi();
    } catch (error) {
      // 即使后端调用失败（如CORS、网络错误等），也要清理本地数据
      console.warn('后端登出接口调用失败，可能是CORS或网络问题:', error);
      // 对于CORS错误，这是正常现象，不影响登出流程
    } finally {
      // 无论后端调用是否成功，都要清理本地数据
      token.value = null;
      clearAllSession();
      router.push('/login');
    }
  }

  function checkTokenExpiry() {
    if (!token.value) {
      logout();
    }
  }

  return { 
    token, 
    isLoggedIn, 
    setToken, 
    logout, 
    checkTokenExpiry 
  };
});

export default useUserStore;
