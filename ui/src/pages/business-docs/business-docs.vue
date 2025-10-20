<script lang="js" setup>
import { ref, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import * as Api from '@src/api';
import DocEditorModal from './doc-editor-modal.vue';
import DocViewer from './doc-viewer.vue';

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
            <el-table-column prop="createTime" label="创建时间" />
            <el-table-column label="操作">
              <template #default="scope">
                <el-button size="small" @click="handleView(scope.row.id)">查看</el-button>
                <el-button size="small" @click="handleEdit(scope.row.id)">编辑</el-button>
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
</template>