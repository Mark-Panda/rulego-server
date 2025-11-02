<script lang="js" setup>
import { ref, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import * as Api from '@src/api';
import DocEditorModal from './doc-editor-modal.vue';
import DocViewer from './doc-viewer.vue';
import router from '@src/router';

const paginationState = ref({
  page: 1,
  size: 10,
  pageSizes: [10, 20, 30, 40, 50],
  total: 0,
});
const formState = ref({
  keywords: '',
});
const data = ref([]);
const loading = ref(false);

// 控制模态框显示
const editorModalVisible = ref(false);
const viewerModalVisible = ref(false);

// 当前编辑或查看的文档
const currentDoc = ref({});

// 生成工作流弹窗状态
const generateDialogVisible = ref(false);
const generateLoading = ref(false);
const selectedChainId = ref('');
const workflowOptions = ref([]);
let generateSourceDoc = null;

async function refreshData() {
  try {
    loading.value = true;
    const params = {
      page: paginationState.value.page,
      size: paginationState.value.size,
      keywords: formState.value.keywords
    };
    
    const res = await Api.getDocs(params);
    data.value = res.items || [];
    paginationState.value.total = res.total || 0;
  } catch (error) {
    console.error('获取文档列表失败:', error);
    ElMessage.error('获取文档列表失败');
  } finally {
    loading.value = false;
  }
}

function paginationChangeHandler(currentPage, pageSize) {
  paginationState.value.page = currentPage;
  paginationState.value.size = pageSize;
  refreshData();
}

// 处理新建文档
function handleCreate() {
  currentDoc.value = {};
  editorModalVisible.value = true;
}

// 处理编辑文档
function handleEdit(id) {
  const doc = data.value.find(item => item.id === id);
  if (doc) {
    currentDoc.value = { ...doc };
    editorModalVisible.value = true;
  }
}

// 处理删除文档
function handleDelete(id, name) {
  ElMessageBox.confirm(`确定要删除文档 "${name}" 吗？此操作不可恢复。`, '删除确认', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await Api.deleteDoc(id);
      ElMessage.success('删除成功');
      refreshData();
    } catch (error) {
      console.error('删除文档失败:', error);
      ElMessage.error('删除文档失败');
    }
  }).catch(() => {
    // 用户取消删除
  });
}

// 处理查看文档详情
function handleView(id) {
  const doc = data.value.find(item => item.id === id);
  if (doc) {
    currentDoc.value = { ...doc };
    viewerModalVisible.value = true;
  }
}

// 保存文档
async function handleSaveDoc(docData) {
  try {
    if (docData.id) {
      // 更新文档
      await Api.updateDoc(docData.id, docData);
      ElMessage.success('文档更新成功');
    } else {
      // 创建文档
      await Api.createDoc(docData);
      ElMessage.success('文档创建成功');
    }
    refreshData();
  } catch (error) {
    console.error('保存文档失败:', error);
    ElMessage.error('保存文档失败');
  }
}

// 查看关联工作流
function handleOpenWorkflow(chainId) {
  if (!chainId) {
    ElMessage.warning('该文档未关联工作流');
    return;
  }
  router.push({ path: '/workflow', query: { id: chainId } });
}

// 打开生成工作流弹窗
async function handleGenerateWorkflow(id) {
  const doc = data.value.find(item => item.id === id);
  if (!doc) return;
  generateSourceDoc = doc;
  selectedChainId.value = '';
  await loadDeployedWorkflows();
  generateDialogVisible.value = true;
}

// 加载已部署的工作流作为下拉选项
async function loadDeployedWorkflows() {
  try {
    const res = await Api.getRules({ page: 1, size: 200 });
    const items = Array.isArray(res?.items) ? res.items : [];
    // 过滤已部署（未禁用）
    const deployed = items.filter(it => it?.ruleChain && it.ruleChain.disabled === false);
    workflowOptions.value = deployed.map(it => ({
      label: it.ruleChain.name,
      value: it.ruleChain.id,
    }));
  } catch (e) {
    workflowOptions.value = [];
  }
}

// 提交生成工作流执行
async function submitGenerateWorkflow() {
  if (!selectedChainId.value) {
    ElMessage.warning('请选择工作流');
    return;
  }
  try {
    generateLoading.value = true;
    let body = generateSourceDoc?.content || '';
    let dataBody = body;
    let headers = { 'Content-Type': 'application/json' };
    // 尝试将content解析为JSON
    try {
      dataBody = body ? JSON.parse(body) : {};
    } catch (err) {
      // 解析失败则作为纯文本发送
      headers = { 'Content-Type': 'text/plain' };
      dataBody = body || '';
    }
    await Api.executeRules({ id: selectedChainId.value, msgType: 'json', data: dataBody, headers });
    ElMessage.success('已触发工作流执行');
    generateDialogVisible.value = false;
  } catch (error) {
    ElMessage.error('触发工作流执行失败');
  } finally {
    generateLoading.value = false;
  }
}

onMounted(() => {
  refreshData();
});
</script>

<template>
  <div class="flex h-full flex-col rounded-lg bg-white dark:bg-gray-800 shadow-sm">
    <!-- 页面标题和操作按钮 -->
    <div class="flex flex-none items-center justify-between border-b border-gray-100 dark:border-gray-700 p-4">
      <div class="flex items-center">
        <h1 class="text-xl font-medium text-gray-800 dark:text-gray-200">业务文档管理</h1>
        <el-tag type="info" effect="plain" class="ml-2">{{ paginationState.total }} 个文档</el-tag>
      </div>
      <div class="flex items-center space-x-4">
        <div class="flex items-center space-x-3">
          <el-input
            placeholder="搜索文档名称"
            :clearable="true"
            style="width: 220px;"
            v-model="formState.keywords"
            @keyup.enter="refreshData"
          >
            <template #prefix>
              <el-icon><el-icon-search /></el-icon>
            </template>
          </el-input>
        </div>
        <div class="flex items-center space-x-3">
          <el-button type="primary" @click="handleCreate">
            <el-icon class="mr-1"><el-icon-plus /></el-icon>
            <span>新建文档</span>
          </el-button>
        </div>
      </div>
    </div>
    
    <!-- 文档列表 -->
    <div class="flex-grow overflow-auto p-4">
      <el-scrollbar class="h-full w-full">
        <div v-if="data.length === 0 && !loading" class="flex h-64 w-full items-center justify-center">
          <div class="text-center">
            <el-icon class="text-4xl text-gray-300 dark:text-gray-600"><el-icon-document /></el-icon>
            <p class="mt-2 text-gray-500 dark:text-gray-400">暂无业务文档，点击"新建文档"开始</p>
            <el-button class="mt-4" type="primary" @click="handleCreate">
              <el-icon class="mr-1"><el-icon-plus /></el-icon>
              <span>新建文档</span>
            </el-button>
          </div>
        </div>
        <div v-else class="flex flex-wrap content-start">
          <!-- 文档表格 -->
          <el-table :data="data" style="width: 100%" v-loading="loading">
            <el-table-column prop="name" label="文档名称" />
            <el-table-column prop="description" label="描述" />
            <el-table-column prop="chainName" label="关联工作流" />
            <el-table-column prop="createTime" label="创建时间" />
            <el-table-column label="操作">
              <template #default="scope">
                <el-button size="small" @click="handleView(scope.row.id)">查看</el-button>
                <el-button size="small" @click="handleEdit(scope.row.id)">编辑</el-button>
                <el-button 
                  v-if="scope.row.chainId"
                  size="small" 
                  type="primary" 
                  @click="handleOpenWorkflow(scope.row.chainId)"
                >查看工作流</el-button>
                <el-button 
                  size="small" 
                  type="success" 
                  @click="handleGenerateWorkflow(scope.row.id)"
                >生成工作流</el-button>
                <el-button size="small" type="danger" @click="handleDelete(scope.row.id, scope.row.name)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-scrollbar>
    </div>
    
    <!-- 分页 -->
    <div class="flex flex-none justify-end border-t border-gray-100 dark:border-gray-700 p-4">
      <el-pagination
        v-model:current-page="paginationState.page"
        v-model:page-size="paginationState.size"
        size="small"
        layout="total, sizes, prev, pager, next, jumper"
        :page-sizes="paginationState.pageSizes"
        :total="paginationState.total"
        :background="true"
        @change="paginationChangeHandler"
      ></el-pagination>
    </div>
  </div>
  
  <!-- 文档编辑模态框 -->
  <doc-editor-modal
    v-model:visible="editorModalVisible"
    :doc-data="currentDoc"
    @save="handleSaveDoc"
  />
  
  <!-- 文档查看模态框 -->
  <doc-viewer
    v-model:visible="viewerModalVisible"
    :doc-data="currentDoc"
  />

  <!-- 生成工作流弹窗 -->
  <el-dialog 
    v-model="generateDialogVisible" 
    title="生成工作流" 
    width="520px" 
    :close-on-click-modal="false"
  >
    <div>
      <el-form label-width="100px">
        <el-form-item label="选择工作流" required>
          <el-select 
            v-model="selectedChainId" 
            placeholder="请选择已部署工作流" 
            filterable 
            style="width:100%"
          >
            <el-option 
              v-for="opt in workflowOptions" 
              :key="opt.value" 
              :label="opt.label" 
              :value="opt.value" 
            />
          </el-select>
        </el-form-item>
      </el-form>
    </div>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="generateDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="generateLoading" @click="submitGenerateWorkflow">确定</el-button>
      </span>
    </template>
  </el-dialog>
</template>