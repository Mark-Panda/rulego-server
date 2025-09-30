<template>
  <div class="smart-ai-assistant h-full flex flex-col bg-gradient-to-br from-slate-50 to-blue-50">
    <!-- 功能选项卡 -->
    <div class="tabs-section">
      <el-tabs v-model="activeTab" class="ai-tabs" @tab-change="handleTabChange">
        <el-tab-pane label="💬 AI对话" name="chat">
          <!-- 聊天界面内容 -->
        </el-tab-pane>
        <el-tab-pane label="📋 历史记录" name="history">
          <!-- 历史记录面板 -->
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 聊天界面 -->
    <div v-show="activeTab === 'chat'" class="chat-tab-content">
      <!-- 配置区域 -->
      <div class="config-section">
        <div class="config-header" @click="configCollapsed = !configCollapsed">
          <el-icon :size="16" class="text-blue-600">
            <el-icon-setting />
          </el-icon>
          <span class="text-sm font-medium text-gray-700">AI服务配置</span>
          <el-button 
            :icon="configCollapsed ? 'el-icon-arrow-down' : 'el-icon-arrow-up'" 
            size="small" 
            text 
            class="ml-auto"
          />
        </div>
        
        <el-collapse-transition>
          <div v-show="!configCollapsed">
            <el-form 
              :model="gptConfig" 
              :rules="rules" 
              ref="formRef" 
              label-width="80px" 
              class="config-form"
              size="small"
            >
              <el-form-item label="Base URL" prop="url">
                <el-input 
                  placeholder="https://api.openai.com" 
                  v-model="gptConfig.url" 
                  clearable 
                  class="config-input"
                />
              </el-form-item>
              <el-form-item label="API Key" prop="apiKey">
                <el-input 
                  type="password" 
                  placeholder="sk-..." 
                  v-model="gptConfig.apiKey" 
                  show-password 
                  clearable 
                  class="config-input"
                />
              </el-form-item>
              <el-form-item label="Model" prop="model">
                <el-input 
                  placeholder="gpt-3.5-turbo" 
                  v-model="gptConfig.model" 
                  clearable 
                  class="config-input"
                />
              </el-form-item>
            </el-form>
          </div>
        </el-collapse-transition>
      </div>
      
      <!-- 聊天界面 -->
      <div class="chat-container">
        <div class="chat-messages">
          <BubbleList 
            :list="list" 
            max-height="auto" 
            class="bubble-list"
          />
        </div>
        
        <!-- 输入区域 -->
        <div class="input-section">
          <div class="input-wrapper">
            <Sender 
              ref="senderRef" 
              v-model="inputMessage" 
              @submit="handleSend" 
              placeholder="💬 告诉AI如何修改你的流程图..."
              class="message-sender"
            />
          </div>
          <div class="input-tips">
            <el-text size="small" type="info">
              💡 试试说："添加一个HTTP请求节点"或"优化流程布局"
            </el-text>
          </div>
        </div>
      </div>
    </div>

    <!-- 历史记录界面 -->
    <div v-show="activeTab === 'history'" class="history-tab-content">
      <schema-history-panel 
        :history-manager="props.schemaHistoryManager"
        @schema-change="handleHistorySchemaChange"
        @undo="handleHistoryUndo"
        @redo="handleHistoryRedo"
      />
    </div>
  </div>
</template>

<script setup>
import OpenAI from 'openai';
import { onMounted, reactive, ref, watch } from 'vue';
import { BubbleList, Sender } from 'vue-element-plus-x';
import { ElMessage } from 'element-plus';
import SchemaHistoryPanel from '@src/components/schema-history/schema-history-panel.vue';

const props = defineProps({
  currentSchema: {
    type: Object,
    default: () => ({})
  },
  onSchemaUpdate: {
    type: Function,
    default: null
  },
  schemaHistoryManager: {
    type: Object,
    default: null
  }
});

const gptConfig = reactive({
  url: 'https://api.deepseek.com',
  apiKey: '',
  model: 'deepseek-chat',
});

const rules = {
  url: [{ required: true, message: 'Base URL is required', trigger: 'blur' }],
  apiKey: [{ required: true, message: 'API Key is required', trigger: 'blur' }],
  model: [{ required: true, message: 'Model is required', trigger: 'blur' }],
};

const formRef = ref();
const inputMessage = ref('');
let openai;
const senderRef = ref();
const configCollapsed = ref(true);
const activeTab = ref('chat');

const list = ref([
  {
    key: 1,
    role: 'ai',
    placement: 'start',
    content: '你好！我是RuleGo智能助手，我可以帮你修改流程图的配置。我可以帮你修改节点配置和属性、添加新的节点和连线、删除不需要的节点、优化流程结构、解释当前流程配置。新功能：点击"历史记录"选项卡查看所有修改记录、支持撤销/重做操作、自动保存你的修改历史。请告诉我你想要如何修改你的流程图！',
    loading: false,
    shape: 'corner',
    variant: 'filled',
    isMarkdown: true,
    typing: false,
    avatar: 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png',
    avatarSize: '24px',
  }
]);

// 初始化 OpenAI 实例
function initOpenAi() {
  if (!gptConfig.apiKey || !gptConfig.url) return;
  
  openai = new OpenAI({
    baseURL: gptConfig.url,
    apiKey: gptConfig.apiKey,
    dangerouslyAllowBrowser: true,
  });
}

// 生成AI系统提示词
function generateSystemPrompt() {
  return `你是RuleGo工作流设计助手，专门帮助用户修改和优化流程图配置。

当前流程图配置:
\`\`\`json
${JSON.stringify(props.currentSchema || {}, null, 2)}
\`\`\`

请根据用户的需求，提供以下格式的响应：
{"schema":{修改后的完整JSON配置}}

注意：
- 保持现有配置的完整性
- 只修改用户明确要求的部分
- 确保JSON格式正确
- 保持节点ID的唯一性
- 维护节点间的连接关系

用户的问题：`;
}

// 处理发送消息
async function handleSend(val) {
  const isValid = await new Promise(resolve => {
    formRef.value.validate((valid) => {
      resolve(valid);
    });
  });
  
  if (!isValid) {
    ElMessage.warning('请先完整配置AI服务信息');
    return;
  }

  senderRef.value?.clear();
  
  // 记录用户输入到历史
  if (props.schemaHistoryManager) {
    props.schemaHistoryManager.addRecord(
      props.currentSchema,
      'user-input',
      `用户询问: ${val}`,
      { userInput: val }
    );
  }
  
  // 添加用户消息到列表
  list.value.push({
    key: list.value.length + 1,
    role: 'user',
    placement: 'end',
    content: val,
    loading: false,
    shape: 'corner',
    variant: 'outlined',
    isMarkdown: true,
    typing: false,
    avatar: 'https://avatars.githubusercontent.com/u/76239030?v=4',
    avatarSize: '24px',
  });

  // 添加加载状态的AI回复
  const loadingMessageKey = list.value.length + 1;
  list.value.push({
    key: loadingMessageKey,
    role: 'ai',
    placement: 'start',
    content: '🤖 正在分析你的需求并生成方案...',
    loading: true,
    shape: 'corner',
    variant: 'filled',
    isMarkdown: true,
    typing: true,
    avatar: 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png',
    avatarSize: '24px',
  });

  try {
    const prompt = generateSystemPrompt();
    console.log('prompt',prompt);
    console.log("用户问题", val)
    const response = await openai.chat.completions.create({
      messages: [
        {
          role: 'system',
          content: prompt,
        },
        {
          role: 'user',
          content: val,
        },
      ],
      model: gptConfig.model,
      temperature: 0.3,
    });

    const aiResponse = response.choices[0].message.content || '抱歉，我无法理解你的需求。';
    console.log('AI Response:', aiResponse);
    // 更新AI回复
    const loadingMessageIndex = list.value.findIndex(item => item.key === loadingMessageKey);
    if (loadingMessageIndex !== -1) {
      list.value[loadingMessageIndex] = {
        ...list.value[loadingMessageIndex],
        content: aiResponse,
        loading: false,
        typing: false,
      };
    }

    // 检查是否包含schema更新
    const schemaUpdateMatch = aiResponse;
    if (schemaUpdateMatch && props.onSchemaUpdate) {
      try {
        const schemaUpdate = JSON.parse(schemaUpdateMatch);
        const newSchema = schemaUpdate["schema"];
        console.log("newSchema", newSchema);
        
        // 记录AI修改到历史
        let historyRecord = null;
        if (props.schemaHistoryManager) {
          historyRecord = props.schemaHistoryManager.addRecord(
            newSchema,
            'ai-modification',
            `AI修改: ${val}`,
            { 
              aiResponse: aiResponse,
              userInput: val,
              timestamp: Date.now()
            }
          );
        }
        
        // 应用schema更新
        props.onSchemaUpdate(newSchema);
        
        // 添加成功提示消息
        list.value.push({
          key: list.value.length + 1,
          role: 'ai',
          placement: 'start',
          content: `✅ 流程图配置已成功更新！

🔄 修改记录ID: ${historyRecord?.id || '未记录'}
📝 版本: ${historyRecord?.version || '未知'}
⏰ 时间: ${new Date().toLocaleString()}

你可以在"历史记录"选项卡中查看详细信息，或使用撤销/重做功能。`,
          loading: false,
          shape: 'corner',
          variant: 'filled',
          isMarkdown: true,
          typing: false,
          avatar: 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png',
          avatarSize: '24px',
        });
        
        ElMessage.success('AI已成功修改流程图配置');
      } catch (error) {
        console.error('解析AI返回的schema失败:', error);
        ElMessage.error('AI返回的配置格式有误，无法应用更改');
      }
    }

  } catch (error) {
    console.error('AI调用失败:', error);
    
    const loadingMessageIndex = list.value.findIndex(item => item.key === loadingMessageKey);
    if (loadingMessageIndex !== -1) {
      list.value[loadingMessageIndex] = {
        ...list.value[loadingMessageIndex],
        content: '❌ 抱歉，AI服务暂时不可用。请检查网络连接和API配置。',
        loading: false,
        typing: false,
      };
    }
    
    ElMessage.error('AI调用失败，请检查配置和网络连接');
  }
}

// 处理选项卡切换
function handleTabChange(tabName) {
  activeTab.value = tabName;
  if (tabName === 'history') {
    if (props.schemaHistoryManager) {
      props.schemaHistoryManager.stopAutoSave();
    }
  } else if (tabName === 'chat') {
    startAutoSave();
  }
}

// 处理历史记录的schema变更
function handleHistorySchemaChange(newSchema) {
  if (props.onSchemaUpdate) {
    props.onSchemaUpdate(newSchema);
  }
}

// 处理历史记录撤销
function handleHistoryUndo(result) {
  list.value.push({
    key: list.value.length + 1,
    role: 'ai',
    placement: 'start',
    content: `🔄 撤销操作完成\n\n已回退到: ${result.record.description}\n时间: ${new Date(result.record.timestamp).toLocaleString()}`,
    loading: false,
    shape: 'corner',
    variant: 'filled',
    isMarkdown: true,
    typing: false,
    avatar: 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png',
    avatarSize: '24px',
  });
}

// 处理历史记录重做
function handleHistoryRedo(result) {
  list.value.push({
    key: list.value.length + 1,
    role: 'ai',
    placement: 'start',
    content: `🔄 重做操作完成\n\n已前进到: ${result.record.description}\n时间: ${new Date(result.record.timestamp).toLocaleString()}`,
    loading: false,
    shape: 'corner',
    variant: 'filled',
    isMarkdown: true,
    typing: false,
    avatar: 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png',
    avatarSize: '24px',
  });
}

// 开始自动保存
function startAutoSave() {
  if (props.schemaHistoryManager) {
    props.schemaHistoryManager.startAutoSave(() => props.currentSchema);
  }
}

// 监听配置变化
watch(
  () => [gptConfig.apiKey, gptConfig.url, gptConfig.model],
  () => {
    initOpenAi();
    localStorage.setItem('smartAiConfig', JSON.stringify(gptConfig));
  }
);

// 监听当前schema变化
watch(
  () => props.currentSchema,
  (newSchema, oldSchema) => {
    if (newSchema && (!oldSchema || JSON.stringify(newSchema) !== JSON.stringify(oldSchema))) {
      if (props.schemaHistoryManager && props.schemaHistoryManager.history.length === 0) {
        props.schemaHistoryManager.addRecord(
          newSchema,
          'initial',
          '初始流程图状态',
          { isInitial: true }
        );
      }
    }
  },
  { deep: true, immediate: true }
);

// 初始化
onMounted(() => {
  const storedConfig = localStorage.getItem('smartAiConfig');
  if (storedConfig) {
    try {
      const parsedConfig = JSON.parse(storedConfig);
      Object.assign(gptConfig, parsedConfig);
    } catch (error) {
      console.warn('无法解析存储的AI配置:', error);
    }
  }
  initOpenAi();
  
  if (props.schemaHistoryManager) {
    props.schemaHistoryManager.loadFromLocalStorage();
    startAutoSave();
  }
});
</script>

<style scoped>
.smart-ai-assistant {
  height: 100%;
  background: linear-gradient(135deg, #f8fafc 0%, #e0f2fe 100%);
  border-radius: 0;
}

/* 选项卡样式 */
.tabs-section {
  border-bottom: 1px solid #e2e8f0;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

.ai-tabs {
  padding: 0 16px;
}

.ai-tabs .el-tabs__header {
  margin: 0;
}

.ai-tabs .el-tabs__nav-wrap::after {
  display: none;
}

.ai-tabs .el-tabs__item {
  font-size: 13px;
  font-weight: 500;
  color: #6b7280;
  transition: all 0.2s ease;
}

.ai-tabs .el-tabs__item.is-active {
  color: #3b82f6;
  font-weight: 600;
}

/* 选项卡内容区域 */
.chat-tab-content,
.history-tab-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.history-tab-content {
  padding: 0;
}

/* 配置区域样式 */
.config-section {
  background: rgba(255, 255, 255, 0.95);
  border-bottom: 1px solid #e2e8f0;
  backdrop-filter: blur(10px);
}

.config-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.config-header:hover {
  background: rgba(59, 130, 246, 0.05);
}

.config-form {
  padding: 0 16px 16px;
}

.config-input .el-input__wrapper {
  border-radius: 8px;
  transition: all 0.2s ease;
}

.config-input .el-input__wrapper:hover {
  box-shadow: 0 0 0 1px #3b82f6;
}

/* 聊天容器样式 */
.chat-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.chat-messages {
  flex: 1;
  overflow: hidden;
  padding: 16px;
}

.bubble-list {
  height: 100%;
}

/* 输入区域样式 */
.input-section {
  background: rgba(255, 255, 255, 0.98);
  border-top: 1px solid #e2e8f0;
  padding: 16px;
  backdrop-filter: blur(10px);
}

.input-wrapper {
  margin-bottom: 8px;
}

.message-sender {
  border-radius: 12px;
}

.input-tips {
  text-align: center;
  opacity: 0.7;
}

/* 响应式设计 */
@media (max-width: 480px) {
  .config-form {
    padding: 0 12px 12px;
  }
  
  .chat-messages {
    padding: 12px;
  }
  
  .input-section {
    padding: 12px;
  }
}

/* 深色模式适配 */
.dark .smart-ai-assistant {
  background: linear-gradient(135deg, #1e293b 0%, #0f172a 100%);
}

.dark .config-section {
  background: rgba(30, 41, 59, 0.95);
  border-bottom-color: #475569;
}

.dark .config-header:hover {
  background: rgba(59, 130, 246, 0.1);
}

.dark .input-section {
  background: rgba(30, 41, 59, 0.98);
  border-top-color: #475569;
}
</style>