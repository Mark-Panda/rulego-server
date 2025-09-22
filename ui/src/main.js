import { createApp } from 'vue';
import ElementPlus from 'element-plus';
import * as ElementPlusIconsVue from '@element-plus/icons-vue';
import App from '@src/app.vue';
import router from '@src/router';
import store from '@src/store';
import { initTheme } from '@src/utils/theme';
import useUserStore from '@src/store/module/user';

import '@src/style/tailwind.css';
import 'element-plus/dist/index.css';
import '@logicflow/core/lib/style/index.css';
import '@logicflow/extension/lib/style/index.css';
import '@src/style/logic-flow.css';
import '@src/style/element-plus.css';

// 初始化主题
initTheme();

const app = createApp(App);

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(`ElIcon${key}`, component);
}

app.use(ElementPlus).use(store).use(router);

// 在应用启动后检查token有效性
app.mount('#app');

// 在应用挂载后初始化用户store
const userStore = useUserStore();
userStore.checkTokenExpiry();
