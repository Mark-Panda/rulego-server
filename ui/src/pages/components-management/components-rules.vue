<script lang="js" setup>
import { ref, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { QuillEditor } from '@vueup/vue-quill';
import '@vueup/vue-quill/dist/vue-quill.snow.css';
import * as Api from '@src/api';

const tableData = ref([]);
const loading = ref(false);
const searchForm = ref({
  keywords: '',
  componentType: '',
  disabled: ''
});

const paginationState = ref({
  page: 1,
  size: 10,
  pageSizes: [10, 20, 30, 40, 50],
  total: 0
});

// 新增/编辑/查看规则表单相关
const dialogVisible = ref(false);
const formRef = ref();
const isEditing = ref(false); // 是否为编辑模式
const isViewing = ref(false); // 是否为查看模式
const editingId = ref(null); // 编辑的规则ID
const ruleForm = ref({
  componentName: '',
  componentType: '',
  disabled: false,
  useDesc: '',
  useRuleDesc: ''
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
  componentName: [
    { required: true, message: '请输入组件名称', trigger: 'blur' },
    { min: 2, max: 50, message: '组件名称长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  componentType: [
    { required: true, message: '请选择组件类型', trigger: 'change' }
  ],
  useDesc: [
    { required: true, message: '请输入使用描述', trigger: 'blur' }
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

// 组件类型选项（从已安装组件获取）
const componentTypeOptions = ref([
  { label: '全部类型', value: '' },
  { label: '接入端', value: 'endpoint' },
  { label: '过滤器', value: 'filter' },
  { label: '转换器', value: 'transform' },
  { label: '动作', value: 'action' }
]);

// 已安装组件列表
const installedComponents = ref([]);
// 组件名称选项
const componentNameOptions = ref([]);

const disabledOptions = [
  { label: '全部状态', value: '' },
  { label: '启用', value: false },
  { label: '禁用', value: true }
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

// 获取已安装组件列表
const fetchInstalledComponents = async () => {
  try {
    const res = await Api.getInstalledComponents();
    console.log('已安装组件API返回数据:', res);
    
    const processedData = [];
    const typeSet = new Set();
    
    // 处理endpoints数据（输入端组件）
    if (res && res.endpoints && Array.isArray(res.endpoints)) {
      res.endpoints.forEach(item => {
        const component = {
          name: item.name || item.type,
          type: 'endpoint',
          category: 'endpoint',
          description: item.desc || '输入端组件',
          _original: item
        };
        processedData.push(component);
        typeSet.add('endpoint');
      });
    }
    
    // 处理nodes数据（过滤器、转换器、动作等）
    if (res && res.nodes && Array.isArray(res.nodes)) {
      res.nodes.forEach(item => {
        let category = item.category;
        if (category && category.startsWith('external/')) {
          category = 'external';
        }
        
        const component = {
          name: item.type || item.name,
          type: category,
          category: category,
          description: item.desc || '-',
          _original: item
        };
        processedData.push(component);
        typeSet.add(category);
      });
    }
    
    installedComponents.value = processedData;
    
    // 更新组件名称选项
    componentNameOptions.value = processedData.map(item => ({
      label: item.name,
      value: item.name,
      type: item.type
    }));
    
    // 更新组件类型选项
    const typeOptionsMap = {
      'endpoint': '输入端',
      'filter': '过滤器', 
      'transform': '转换器',
      'action': '动作',
      'external': '外部组件',
      'flow': '子规则链'
    };
    
    const newTypeOptions = [{ label: '全部类型', value: '' }];
    typeSet.forEach(type => {
      if (type && typeOptionsMap[type]) {
        newTypeOptions.push({
          label: typeOptionsMap[type],
          value: type
        });
      }
    });
    componentTypeOptions.value = newTypeOptions;
    
  } catch (error) {
    console.error('获取已安装组件列表失败:', error);
    ElMessage.error('获取已安装组件列表失败');
  }
};

// 获取组件规则列表
const fetchRulesList = async () => {
  loading.value = true;
  try {
    const params = {
      page: paginationState.value.page,
      size: paginationState.value.size,
      keywords: searchForm.value.keywords || undefined,
      componentType: searchForm.value.componentType || undefined,
      disabled: searchForm.value.disabled !== '' ? searchForm.value.disabled : undefined
    };
    
    // 调用API获取组件使用规则分页列表
    const res = await Api.getComponentUseRulePage(params);
    console.log('组件使用规则API返回数据:', res);
    
    // 处理API返回的数据结构
    if (res && res.list && Array.isArray(res.list)) {
      // 处理数据，转换为表格需要的格式
      tableData.value = res.list.map(item => ({
        id: item.id,
        name: item.componentName || '未知组件',
        type: item.componentType || '',
        description: item.useDesc || '暂无描述',
        status: item.disabled ? 'inactive' : 'active',
        createTime: item.createdAt ? new Date(item.createdAt).toLocaleString('zh-CN') : '-',
        updateTime: item.updatedAt ? new Date(item.updatedAt).toLocaleString('zh-CN') : '-',
        usageRules: item.useRuleDesc || '',
        _original: item // 保存原始数据
      }));
      
      paginationState.value.total = res.total || 0;
    } else {
      // 处理空数据或异常数据结构
      tableData.value = [];
      paginationState.value.total = 0;
      if (!res || res.list === null) {
        console.log('组件使用规则暂无数据');
      } else {
        console.warn('无法解析的数据结构:', res);
      }
    }
    
    console.log('处理后的表格数据:', tableData.value);
  } catch (error) {
    console.error('获取组件使用规则列表失败:', error);
    ElMessage.error('获取组件使用规则列表失败');
    tableData.value = [];
    paginationState.value.total = 0;
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
    keywords: '',
    componentType: '',
    disabled: ''
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
    
    // 准备更新数据
    const updateData = {
      id: String(row.id),
      componentName: row.name,
      componentType: row.type,
      disabled: row.status === 'active', // 如果当前是启用状态，则设置为禁用
      useDesc: row.description || '',
      useRuleDesc: row.usageRules || ''
    };
    
    // 调用与编辑相同的更新接口
    await Api.updateComponentUseRule(updateData);
    
    // 更新本地数据
    row.status = row.status === 'active' ? 'inactive' : 'active';
    row.updateTime = new Date().toLocaleString('zh-CN');
    
    ElMessage.success(`${action}成功`);
  } catch (error) {
    if (error !== 'cancel') {
      console.error('切换状态失败:', error);
      if (error.response && error.response.data && error.response.data.message) {
        ElMessage.error(`操作失败: ${error.response.data.message}`);
      } else {
        ElMessage.error('操作失败，请稍后重试');
      }
    }
  }
};

// 编辑规则
const editRule = (row) => {
  isEditing.value = true;
  isViewing.value = false;
  editingId.value = row.id;
  dialogVisible.value = true;
  
  // 填充表单数据
  ruleForm.value = {
    componentName: row.name || '',
    componentType: row.type || '',
    disabled: row.status === 'inactive',
    useDesc: row.description || '',
    useRuleDesc: row.usageRules || ''
  };
  
  // 设置富文本编辑器内容
  editorContent.value = row.usageRules || '';
  // 恢复默认placeholder
  editorOptions.placeholder = '请输入组件使用规则...';
};

// 查看规则
const viewRule = (row) => {
  isEditing.value = false;
  isViewing.value = true;
  editingId.value = row.id;
  dialogVisible.value = true;
  
  // 填充表单数据
  ruleForm.value = {
    componentName: row.name || '',
    componentType: row.type || '',
    disabled: row.status === 'inactive',
    useDesc: row.description || '',
    useRuleDesc: row.usageRules || ''
  };
  
  // 设置富文本编辑器内容
  editorContent.value = row.usageRules || '';
  
  // 如果有内容，清除placeholder
  if (row.usageRules && row.usageRules.trim()) {
    editorOptions.placeholder = '';
  } else {
    editorOptions.placeholder = '请输入组件使用规则...';
  }
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
    
    // 调用删除API
    await Api.deleteComponentUseRule({ id: row.id });
    
    // 删除成功后从表格中移除该行
    const index = tableData.value.findIndex(item => item.id === row.id);
    if (index > -1) {
      tableData.value.splice(index, 1);
      paginationState.value.total--;
    }
    
    ElMessage.success('删除成功');
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除规则失败:', error);
      ElMessage.error('删除失败');
    }
  }
};

// 新增规则
const addRule = () => {
  isEditing.value = false;
  isViewing.value = false;
  editingId.value = null;
  dialogVisible.value = true;
  // 重置表单
  ruleForm.value = {
    componentName: '',
    componentType: '',
    disabled: false,
    useDesc: '',
    useRuleDesc: ''
  };
  // 重置富文本编辑器内容
  editorContent.value = '';
  // 恢复默认placeholder
  editorOptions.placeholder = '请输入组件使用规则...';
};

// 富文本编辑器内容变化处理
const onEditorChange = (content) => {
  ruleForm.value.useRuleDesc = content;
};

// 组件名称变化处理
const onComponentNameChange = (componentName) => {
  // 根据选择的组件名称自动设置组件类型
  const selectedComponent = componentNameOptions.value.find(option => option.value === componentName);
  if (selectedComponent) {
    ruleForm.value.componentType = selectedComponent.type;
  }
};

// 关闭弹窗
const closeDialog = () => {
  dialogVisible.value = false;
  isEditing.value = false;
  isViewing.value = false;
  editingId.value = null;
  if (formRef.value) {
    formRef.value.resetFields();
  }
  ruleForm.value = {
    componentName: '',
    componentType: '',
    disabled: false,
    useDesc: '',
    useRuleDesc: ''
  };
  editorContent.value = '';
  // 恢复默认placeholder
  editorOptions.placeholder = '请输入组件使用规则...';
};

// 提交表单
const submitForm = async () => {
  if (!formRef.value) return;
  
  try {
    await formRef.value.validate();
    
    // 准备提交数据
    const submitData = {
      componentName: ruleForm.value.componentName,
      componentType: ruleForm.value.componentType,
      disabled: ruleForm.value.disabled,
      useDesc: ruleForm.value.useDesc,
      useRuleDesc: ruleForm.value.useRuleDesc
    };
    
    console.log('提交数据:', submitData);
    
    if (isEditing.value) {
      // 编辑模式：添加ID并调用更新接口
      submitData.id = String(editingId.value);
      await Api.updateComponentUseRule(submitData);
      ElMessage.success('规则更新成功');
    } else {
      // 新增模式：调用创建接口
      await Api.createComponentUseRule(submitData);
      ElMessage.success('新增规则成功');
    }
    
    closeDialog();
    
    // 刷新列表
    await fetchRulesList();
    
  } catch (error) {
    const action = isEditing.value ? '更新' : '创建';
    console.error(`${action}组件使用规则失败:`, error);
    if (error.response && error.response.data && error.response.data.message) {
      ElMessage.error(`${action}失败: ${error.response.data.message}`);
    } else {
      ElMessage.error(`${action}规则失败，请稍后重试`);
    }
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
onMounted(async () => {
  await fetchInstalledComponents();
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
            v-model="searchForm.keywords"
            placeholder="搜索关键词"
            style="width: 220px;"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><el-icon-search /></el-icon>
            </template>
          </el-input>
          
          <el-select 
            v-model="searchForm.componentType" 
            placeholder="组件类型" 
            style="width: 150px;"
            clearable
          >
            <el-option
              v-for="option in componentTypeOptions"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>
          
          <el-select 
            v-model="searchForm.disabled" 
            placeholder="状态" 
            style="width: 120px;"
            clearable
          >
            <el-option
              v-for="option in disabledOptions"
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
                type="info" 
                size="small" 
                link
                @click="viewRule(row)"
              >
                查看
              </el-button>
              
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
          :page-sizes="paginationState.pageSizes"
          :total="paginationState.total"
          layout="total, sizes, prev, pager, next, jumper"
          :background="true"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </div>
    
    <!-- 新增/编辑/查看规则弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isViewing ? '查看组件规则' : (isEditing ? '编辑组件规则' : '新增组件规则')"
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
        <el-form-item label="组件名称" prop="componentName">
          <el-select
            v-model="ruleForm.componentName"
            placeholder="请选择组件名称"
            style="width: 100%"
            :disabled="isEditing || isViewing"
            filterable
            @change="onComponentNameChange"
          >
            <el-option
              v-for="option in componentNameOptions"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>
        </el-form-item>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="组件类型" prop="componentType">
              <el-select
                v-model="ruleForm.componentType"
                placeholder="请选择组件类型"
                style="width: 100%"
                :disabled="isViewing"
              >
                <el-option
                  v-for="option in componentTypeOptions.filter(opt => opt.value !== '')"
                  :key="option.value"
                  :label="option.label"
                  :value="option.value"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="组件状态" prop="disabled">
              <el-switch
                v-model="ruleForm.disabled"
                active-text="禁用"
                inactive-text="启用"
                :active-value="true"
                :inactive-value="false"
                :disabled="isViewing"
                style="--el-switch-on-color: #ff4949; --el-switch-off-color: #13ce66"
              />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-form-item label="使用描述" prop="useDesc">
          <el-input
            v-model="ruleForm.useDesc"
            type="textarea"
            :rows="3"
            placeholder="请输入使用描述"
            :readonly="isViewing"
          />
        </el-form-item>
        
        <el-form-item label="使用规则描述" prop="useRuleDesc">
          <div class="quill-editor-wrapper">
            <QuillEditor
              v-model:content="editorContent"
              :options="editorOptions"
              @update:content="onEditorChange"
              content-type="html"
              :read-only="isViewing"
            />
          </div>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="closeDialog">
            {{ isViewing ? '关闭' : '取消' }}
          </el-button>
          <el-button 
            v-if="!isViewing"
            type="primary" 
            @click="submitForm"
          >
            {{ isEditing ? '更新' : '确定' }}
          </el-button>
        </span>
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