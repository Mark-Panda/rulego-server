<script lang="js" setup>
import { ref, onMounted } from 'vue';
import * as Api from '@src/api';
import { ElMessage, ElMessageBox } from 'element-plus';

const tableData = ref([]);
const loading = ref(false);
const paginationState = ref({
  page: 1,
  size: 10,
  pageSizes: [10, 20, 30, 40, 50],
  total: 0,
});
const searchForm = ref({
  name: '',
  type: '',
});

// 组件类型选项
const componentTypeOptions = [
  { label: '全部类型', value: '' },
  { label: '过滤器', value: 'filter' },
  { label: '转换器', value: 'transform' },
  { label: '动作', value: 'action' },
  { label: '外部组件', value: 'external' },
  { label: '子规则链', value: 'flow' },
  { label: '输入端', value: 'endpoint' }
];

// 存储所有组件数据
const allComponentsData = ref([]);

async function refreshTableData() {
  loading.value = true;
  try {
    // 调用API获取已安装组件列表，不传分页参数
    const res = await Api.getInstalledComponents();
    console.log('已安装组件API返回数据:', res);
    
    // 处理API返回的数据结构
    const processedData = [];
    
    // 处理endpoints数据（输入端组件）
    if (res && res.endpoints && Array.isArray(res.endpoints)) {
      res.endpoints.forEach(item => {
        processedData.push({
          name: item.name || item.type,
          type: 'endpoint',
          category: 'endpoint',
          version: item.version || '1.0.0',
          description: item.desc || '输入端组件',
          author: item.author || '-',
          id: item.id || item.type,
          _original: item // 保存原始组件信息
        });
      });
    }
    
    // 处理nodes数据（过滤器、转换器、动作等）
    if (res && res.nodes && Array.isArray(res.nodes)) {
      res.nodes.forEach(item => {
        // 处理category，如果以external/开头，则归类为external
        let category = item.category;
        if (category && category.startsWith('external/')) {
          category = 'external';
        }
        
        processedData.push({
          name: item.type || item.name,
          type: category,
          category: category,
          originalCategory: item.category, // 保存原始category
          version: item.version || '1.0.0',
          description: item.desc || '-',
          author: item.author || '-',
          id: item.id || item.type,
          _original: item // 保存原始组件信息
        });
      });
    }
    
    // 保存所有组件数据
    allComponentsData.value = processedData;
    
    // 应用过滤和分页
    applyFilterAndPagination();
    
  } catch (error) {
    ElMessage.error('获取已安装组件列表失败');
    allComponentsData.value = [];
    tableData.value = [];
    paginationState.value.total = 0;
  } finally {
    loading.value = false;
  }
}

// 应用过滤和分页
function applyFilterAndPagination() {
  // 根据搜索条件过滤数据
  let filteredData = allComponentsData.value;
  if (searchForm.value.type) {
    filteredData = filteredData.filter(item => item.category === searchForm.value.type);
  }
  if (searchForm.value.name) {
    const keyword = searchForm.value.name.toLowerCase();
    filteredData = filteredData.filter(item => 
      (item.name && item.name.toLowerCase().includes(keyword)) || 
      (item.description && item.description.toLowerCase().includes(keyword))
    );
  }
  
  // 更新总数
  paginationState.value.total = filteredData.length;
  
  // 应用分页
  const start = (paginationState.value.page - 1) * paginationState.value.size;
  const end = start + paginationState.value.size;
  tableData.value = filteredData.slice(start, end);
}

function handleSizeChange(val) {
  paginationState.value.size = val;
  applyFilterAndPagination();
}

function handleCurrentChange(val) {
  paginationState.value.page = val;
  applyFilterAndPagination();
}

function handleSearch() {
  paginationState.value.page = 1;
  applyFilterAndPagination();
}

function handleReset() {
  searchForm.value.name = '';
  searchForm.value.type = '';
  paginationState.value.page = 1;
  applyFilterAndPagination();
}

function handleUninstall(row) {
  ElMessageBox.confirm(`确定要卸载组件 "${row.name}" 吗？卸载后将无法在工作流中使用该组件。`, '卸载确认', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    try {
      await Api.uninstallComponent(row.id);
      ElMessage.success('卸载组件成功');
      refreshTableData();
    } catch (error) {
      console.error('卸载组件失败:', error);
      ElMessage.error('卸载组件失败');
    }
  }).catch(() => {
    // 用户取消卸载
  });
}

// 获取组件类型的显示文本
function getComponentTypeText(type) {
  const option = componentTypeOptions.find(opt => opt.value === type);
  return option ? option.label : type;
}

// 获取组件类型的标签类型
function getComponentTypeTagType(type) {
  switch (type) {
    case 'filter':
      return 'success';
    case 'transform':
      return 'warning';
    case 'action':
      return 'danger';
    case 'external':
      return 'info';
    case 'flow':
      return 'primary';
    case 'endpoint':
      return 'info';
    default:
      return 'info'; // 使用Element Plus支持的类型值
  }
}

// 查看组件详情
function handleViewDetail(row) {
  // 直接展示原始组件的JSON信息
  const originalData = row._original || row;
  
  // 使用简单的方式显示JSON
  const jsonContent = JSON.stringify(originalData, null, 2);
  const content = `<pre style="white-space: pre-wrap; word-break: break-all;">${jsonContent}</pre>`;
  
  ElMessageBox.alert(content, '组件详情', {
    dangerouslyUseHTMLString: true,
    confirmButtonText: '确定',
    customClass: 'component-detail-dialog'
  });
}

onMounted(() => {
  refreshTableData();
});
</script>

<template>
  <div class="p-4">
    <div class="mb-6">
      <h2 class="text-2xl font-bold mb-2">已安装组件</h2>
      <p class="text-gray-500 dark:text-gray-400">管理系统中已安装的组件</p>
    </div>
    
    <div class="bg-white dark:bg-gray-800 p-4 rounded-md shadow mb-6">
      <div class="flex justify-between items-center">
        <div class="flex items-center space-x-4">
          <el-input
            v-model="searchForm.name"
            placeholder="搜索组件名称"
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
            placeholder="组件类型" 
            style="width: 140px;"
            clearable
          >
            <el-option
              v-for="item in componentTypeOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
          
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </div>
        
        <div>
          <el-button type="primary" @click="$router.push('/components-market')">
            <el-icon class="mr-1"><el-icon-shopping-cart /></el-icon>
            组件市场
          </el-button>
        </div>
      </div>
    </div>
    
    <div class="bg-white dark:bg-gray-800 rounded-md shadow">
      <el-table
        v-loading="loading"
        :data="tableData"
        style="width: 100%"
        border
        stripe
      >
        <el-table-column prop="name" label="组件名称" min-width="150" />
        <el-table-column label="组件类型" min-width="120">
          <template #default="scope">
            <el-tag v-if="scope.row.type" :type="getComponentTypeTagType(scope.row.type)">
              {{ getComponentTypeText(scope.row.type) }}
            </el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button
              size="small"
              type="primary"
              @click="handleViewDetail(scope.row)"
            >
              详情
            </el-button>
            <el-button
              size="small"
              type="danger"
              @click="handleUninstall(scope.row)"
            >
              卸载
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="flex justify-end p-4">
        <el-pagination
          v-model:current-page="paginationState.page"
          v-model:page-size="paginationState.size"
          size="small"
          layout="total, sizes, prev, pager, next, jumper"
          :page-sizes="paginationState.pageSizes"
          :total="paginationState.total"
          :background="true"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        ></el-pagination>
      </div>
    </div>
  </div>
</template>