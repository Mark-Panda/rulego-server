<script lang="js" setup>
import { ref, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { QuillEditor } from '@vueup/vue-quill';
import '@vueup/vue-quill/dist/vue-quill.snow.css';

const tableData = ref([]);
const loading = ref(false);
const searchForm = ref({
  name: '',
  type: '',
  status: ''
});

const paginationState = ref({
  page: 1,
  size: 10,
  total: 0
});

// 新增规则表单相关
const dialogVisible = ref(false);
const formRef = ref();
const ruleForm = ref({
  name: '',
  description: '',
  usageRules: '',
  type: '',
  status: 'active'
});

// 富文本编辑器内容
const editorContent = ref('');

// Quill编辑器配置
const editorOptions = {
  theme: 'snow',
  modules: {
    toolbar: [
      [{ 'size': ['small', false, 'large', 'huge'] }],
      [{ 'color': [] }, { 'background': [] }],
      ['bold', 'italic', 'underline', 'strike'],
      [{ 'list': 'ordered'}, { 'list': 'bullet' }],
      ['link'],
      ['clean']
    ]
  },
  placeholder: '请输入组件使用规则...'
};

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入组件名称', trigger: 'blur' },
    { min: 2, max: 50, message: '组件名称长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入组件描述', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择组件类型', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择组件状态', trigger: 'change' }
  ]
};

// 模拟数据
const mockData = [
  {
    id: 1,
    name: 'HTTP请求验证规则',
    type: 'endpoint',
    description: '用于验证HTTP请求的格式和参数',
    status: 'active',
    createTime: '2024-01-15 10:30:00',
    updateTime: '2024-01-15 10:30:00'
  },
  {
    id: 2,
    name: 'MQTT消息过滤规则',
    type: 'filter',
    description: '过滤MQTT消息中的无效数据',
    status: 'inactive',
    createTime: '2024-01-14 15:20:00',
    updateTime: '2024-01-14 15:20:00'
  },
  {
    id: 3,
    name: '数据转换规则',
    type: 'transform',
    description: '将输入数据转换为指定格式',
    status: 'active',
    createTime: '2024-01-13 09:15:00',
    updateTime: '2024-01-13 09:15:00'
  }
];

const typeOptions = [
  { label: '全部', value: '' },
  { label: '接入端', value: 'endpoint' },
  { label: '过滤器', value: 'filter' },
  { label: '转换器', value: 'transform' },
  { label: '动作', value: 'action' }
];

const statusOptions = [
  { label: '全部', value: '' },
  { label: '启用', value: 'active' },
  { label: '禁用', value: 'inactive' }
];

// 获取组件规则列表
const fetchRulesList = async () => {
  loading.value = true;
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500));
    
    let filteredData = [...mockData];
    
    // 应用搜索过滤
    if (searchForm.value.name) {
      filteredData = filteredData.filter(item => 
        item.name.toLowerCase().includes(searchForm.value.name.toLowerCase())
      );
    }
    
    if (searchForm.value.type) {
      filteredData = filteredData.filter(item => item.type === searchForm.value.type);
    }
    
    if (searchForm.value.status) {
      filteredData = filteredData.filter(item => item.status === searchForm.value.status);
    }
    
    // 分页处理
    const start = (paginationState.value.page - 1) * paginationState.value.size;
    const end = start + paginationState.value.size;
    
    tableData.value = filteredData.slice(start, end);
    paginationState.value.total = filteredData.length;
    
  } catch (error) {
    console.error('获取组件规则列表失败:', error);
    ElMessage.error('获取组件规则列表失败');
  } finally {
    loading.value = false;
  }
};

// 搜索
const handleSearch = () => {
  paginationState.value.page = 1;
  fetchRulesList();
};

// 重置搜索
const handleReset = () => {
  searchForm.value = {
    name: '',
    type: '',
    status: ''
  };
  paginationState.value.page = 1;
  fetchRulesList();
};

// 分页变化
const handlePageChange = (page) => {
  paginationState.value.page = page;
  fetchRulesList();
};

const handleSizeChange = (size) => {
  paginationState.value.size = size;
  paginationState.value.page = 1;
  fetchRulesList();
};

// 启用/禁用规则
const toggleRuleStatus = async (row) => {
  try {
    const action = row.status === 'active' ? '禁用' : '启用';
    await ElMessageBox.confirm(
      `确定要${action}规则 "${row.name}" 吗？`,
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    );
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 300));
    
    row.status = row.status === 'active' ? 'inactive' : 'active';
    row.updateTime = new Date().toLocaleString('zh-CN');
    
    ElMessage.success(`${action}成功`);
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败');
    }
  }
};

// 编辑规则
const editRule = (row) => {
  ElMessage.info('编辑功能开发中...');
};

// 删除规则
const deleteRule = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除规则 "${row.name}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    );
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 300));
    
    const index = tableData.value.findIndex(item => item.id === row.id);
    if (index > -1) {
      tableData.value.splice(index, 1);
      paginationState.value.total--;
    }
    
    ElMessage.success('删除成功');
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败');
    }
  }
};

// 新增规则
const addRule = () => {
  dialogVisible.value = true;
  // 重置表单
  ruleForm.value = {
    name: '',
    description: '',
    usageRules: '',
    type: '',
    status: 'enabled'
  };
  // 重置富文本编辑器内容
  editorContent.value = '';
};

// 富文本编辑器内容变化处理
const onEditorChange = (content) => {
  ruleForm.value.usageRules = content;
};

// 关闭弹窗
const closeDialog = () => {
  dialogVisible.value = false;
  if (formRef.value) {
    formRef.value.resetFields();
  }
  ruleForm.value = {
    name: '',
    description: '',
    usageRules: '',
    type: '',
    status: 'active'
  };
  editorContent.value = '';
};

// 提交表单
const submitForm = async () => {
  if (!formRef.value) return;
  
  try {
    await formRef.value.validate();
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500));
    
    // 创建新规则对象
    const newRule = {
      id: Date.now(),
      name: ruleForm.value.name,
      description: ruleForm.value.description,
      usageRules: ruleForm.value.usageRules,
      type: ruleForm.value.type,
      status: ruleForm.value.status,
      createTime: new Date().toLocaleString('zh-CN'),
      updateTime: new Date().toLocaleString('zh-CN')
    };
    
    // 添加到列表开头
    tableData.value.unshift(newRule);
    paginationState.value.total++;
    
    ElMessage.success('新增规则成功');
    closeDialog();
    
  } catch (error) {
    console.error('表单验证失败:', error);
  }
};

// 获取状态标签类型
const getStatusTagType = (status) => {
  return status === 'active' ? 'success' : 'info';
};

// 获取状态文本
const getStatusText = (status) => {
  return status === 'active' ? '启用' : '禁用';
};

// 获取类型文本
const getTypeText = (type) => {
  const typeMap = {
    endpoint: '接入端',
    filter: '过滤器',
    transform: '转换器',
    action: '动作'
  };
  return typeMap[type] || type;
};
// 页面初始化
onMounted(() => {
  fetchRulesList();
});
</script>

<template>
  <div class="p-6">
    <!-- 页面标题 -->
    <div class="mb-6">
      <h2 class="text-2xl font-bold mb-2">组件规则</h2>
      <p class="text-gray-500 dark:text-gray-400">管理组件的业务规则和验证逻辑</p>
    </div>
    
    <!-- 搜索和操作区域 -->
    <div class="bg-white dark:bg-gray-800 p-4 rounded-md shadow mb-6">
      <div class="flex justify-between items-center">
        <div class="flex items-center space-x-4">
          <el-input
            v-model="searchForm.name"
            placeholder="搜索规则名称"
            style="width: 220px;"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><el-icon-search /></el-icon>
            </template>
          </el-input>
          
          <el-select 
            v-model="searchForm.type" 
            placeholder="规则类型" 
            style="width: 150px;"
            clearable
          >
            <el-option
              v-for="option in typeOptions"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>
          
          <el-select 
            v-model="searchForm.status" 
            placeholder="状态" 
            style="width: 120px;"
            clearable
          >
            <el-option
              v-for="option in statusOptions"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>
          
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </div>
        
        <div>
          <el-button type="primary" @click="addRule">
            <el-icon><el-icon-plus /></el-icon>
            新增规则
          </el-button>
        </div>
      </div>
    </div>
    
    <!-- 规则列表 -->
    <div class="bg-white dark:bg-gray-800 rounded-md shadow">
      <el-table 
        :data="tableData" 
        v-loading="loading"
        style="width: 100%"
        stripe
      >
        <el-table-column prop="name" label="规则名称" min-width="200">
          <template #default="{ row }">
            <div class="font-medium">{{ row.name }}</div>
          </template>
        </el-table-column>
        
        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ getTypeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="description" label="描述" min-width="250">
          <template #default="{ row }">
            <span class="text-gray-600 dark:text-gray-300">{{ row.description }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="updateTime" label="更新时间" width="180">
          <template #default="{ row }">
            <span class="text-gray-500">{{ row.updateTime }}</span>
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <div class="flex space-x-2">
              <el-button 
                type="primary" 
                size="small" 
                link
                @click="editRule(row)"
              >
                编辑
              </el-button>
              
              <el-button 
                :type="row.status === 'active' ? 'warning' : 'success'" 
                size="small" 
                link
                @click="toggleRuleStatus(row)"
              >
                {{ row.status === 'active' ? '禁用' : '启用' }}
              </el-button>
              
              <el-button 
                type="danger" 
                size="small" 
                link
                @click="deleteRule(row)"
              >
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="flex justify-center p-4">
        <el-pagination
          v-model:current-page="paginationState.page"
          v-model:page-size="paginationState.size"
          :page-sizes="[10, 20, 50, 100]"
          :total="paginationState.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </div>
    
    <!-- 新增规则弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      title="新增组件规则"
      width="800px"
      top="8vh"
      :before-close="closeDialog"
    >
      <el-form
        ref="formRef"
        :model="ruleForm"
        :rules="formRules"
        label-width="120px"
        label-position="left"
      >
        <el-form-item label="组件名称" prop="name">
          <el-input
            v-model="ruleForm.name"
            placeholder="请输入组件名称"
            clearable
          />
        </el-form-item>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="组件类型" prop="type">
              <el-select
                v-model="ruleForm.type"
                placeholder="请选择组件类型"
                style="width: 100%"
              >
                <el-option label="接入端" value="endpoint" />
                <el-option label="过滤器" value="filter" />
                <el-option label="转换器" value="transform" />
                <el-option label="动作" value="action" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="组件状态" prop="status">
              <el-switch
                v-model="ruleForm.status"
                active-text="启用"
                inactive-text="禁用"
                active-value="active"
                inactive-value="inactive"
                style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
              />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-form-item label="组件描述" prop="description">
          <el-input 
            v-model="ruleForm.description" 
            type="textarea" 
            :rows="3"
            placeholder="请输入组件描述" 
          />
        </el-form-item>
        
        <el-form-item label="组件使用规则" prop="usageRules">
          <div class="quill-editor-wrapper">
            <QuillEditor
              v-model:content="editorContent"
              :options="editorOptions"
              @update:content="onEditorChange"
              content-type="html"
            />
          </div>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取消</el-button>
          <el-button type="primary" @click="submitForm">提交</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
/* 自定义样式 */
.el-table {
  --el-table-border-color: var(--el-border-color-lighter);
}

/* 暗黑模式下的表格样式 */
:deep(.dark .el-table) {
  --el-table-bg-color: var(--el-bg-color);
  --el-table-tr-bg-color: var(--el-bg-color);
  --el-table-expanded-cell-bg-color: var(--el-bg-color);
}

:deep(.dark .el-table__row) {
  background-color: var(--el-bg-color);
}

:deep(.dark .el-table__row:hover > td) {
  background-color: var(--el-fill-color-light);
}



/* VueQuill编辑器样式 */
.quill-editor-wrapper {
  border: 1px solid var(--el-border-color);
  border-radius: 6px;
  overflow: hidden;
  width: 100%;
  box-sizing: border-box;
}

:deep(.ql-toolbar) {
  border: none;
  border-bottom: 1px solid var(--el-border-color-lighter);
  background: var(--el-fill-color-extra-light);
  padding: 8px 12px;
}

:deep(.ql-container) {
  border: none;
  font-size: 14px;
  line-height: 1.6;
  min-height: 250px;
  width: 100%;
}

:deep(.ql-editor) {
  padding: 12px 16px;
  min-height: 250px;
  color: var(--el-text-color-primary);
  width: 100%;
  box-sizing: border-box;
}

:deep(.ql-editor.ql-blank::before) {
  color: var(--el-text-color-placeholder);
  font-style: normal;
}

:deep(.ql-snow .ql-tooltip) {
  z-index: 2000;
}

/* 暗黑模式适配 */
:deep(.dark .quill-editor-wrapper) {
  border-color: var(--el-border-color);
}

:deep(.dark .ql-toolbar) {
  background: var(--el-fill-color);
  border-bottom-color: var(--el-border-color);
}

:deep(.dark .ql-editor) {
  background: var(--el-bg-color);
  color: var(--el-text-color-primary);
}

/* 弹窗样式 */
.dialog-footer {
  text-align: right;
}

:deep(.el-dialog__header) {
  padding: 20px 20px 10px;
}

:deep(.el-dialog__body) {
  padding: 10px 20px 20px;
}
</style>