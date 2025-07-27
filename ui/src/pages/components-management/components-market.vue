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
});

// 组件类型选项
const componentTypeOptions = [
  { label: '全部类型', value: '' },
  { label: '过滤器', value: 'filter' },
  { label: '转换器', value: 'transform' },
  { label: '动作', value: 'action' },
  { label: '连接器', value: 'connector' }
];

// 组件分类选项
const componentCategoryOptions = [
  { label: '全部分类', value: '' },
  { label: '基础组件', value: 'basic' },
  { label: '高级组件', value: 'advanced' },
  { label: '第三方集成', value: 'integration' },
  { label: '自定义组件', value: 'custom' }
];

async function refreshTableData() {
  loading.value = true;
  try {
    const params = {
      page: paginationState.value.page,
      size: paginationState.value.size,
      keywords: searchForm.value.name,
      checkMy: 'true'
    };
    
    // 调用API获取组件市场列表
    const res = await Api.getMarketComponents(params);
    console.log('组件市场API返回数据:', res);
    
    // 处理API返回的数据结构
    if (res && res.items && Array.isArray(res.items)) {
      // 处理数据，提取需要的字段
      const processedData = res.items.map(item => {
        const ruleChain = item.ruleChain || {};
        const additionalInfo = ruleChain.additionalInfo || {};
        
        return {
          id: ruleChain.id,
          name: ruleChain.name || '未命名组件',
          root: ruleChain.root ? '是' : '否',
          description: additionalInfo.description || '-',
          author: additionalInfo.username || '-',
          updateTime: additionalInfo.updateTime || '-',
          _original: item // 保存原始数据
        };
      });
      
      tableData.value = processedData;
      paginationState.value.total = res.total || processedData.length;
    } else {
      // 其他情况，确保是空数组
      tableData.value = [];
      paginationState.value.total = 0;
      console.warn('无法解析的数据结构:', res);
    }
    
    console.log('处理后的表格数据:', tableData.value);
  } catch (error) {
    console.error('获取组件市场列表失败:', error);
    ElMessage.error('获取组件市场列表失败');
    tableData.value = [];
    paginationState.value.total = 0;
  } finally {
    loading.value = false;
  }
}

function handleSizeChange(val) {
  paginationState.value.size = val;
  refreshTableData();
}

function handleCurrentChange(val) {
  paginationState.value.page = val;
  refreshTableData();
}

function handleSearch() {
  paginationState.value.page = 1;
  refreshTableData();
}

function handleReset() {
  searchForm.value.name = '';
  paginationState.value.page = 1;
  refreshTableData();
}

function handleInstall(row) {
  ElMessageBox.confirm(`确定要安装组件 "${row.name}" 吗？安装后可在工作流中使用该组件。`, '安装确认', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'info',
  }).then(async () => {
    try {
      await Api.installComponent(row.id);
      ElMessage.success('安装组件成功');
      refreshTableData();
    } catch (error) {
      console.error('安装组件失败:', error);
      ElMessage.error('安装组件失败');
    }
  }).catch(() => {
    // 用户取消安装
  });
}

function handleViewDetail(row) {
  // 查看组件详情 - 显示原始数据
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
    case 'connector':
      return 'info';
    default:
      return 'info'; // 使用Element Plus支持的类型值
  }
}

// 获取组件分类的显示文本
function getComponentCategoryText(category) {
  const option = componentCategoryOptions.find(opt => opt.value === category);
  return option ? option.label : category;
}

onMounted(() => {
  refreshTableData();
});
</script>

<template>
  <div class="p-4">
    <div class="mb-6">
      <h2 class="text-2xl font-bold mb-2">组件市场</h2>
      <p class="text-gray-500 dark:text-gray-400">浏览和安装可用的组件</p>
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
          
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </div>
        
        <div>
          <el-button type="primary" @click="$router.push('/components-installed')">
            <el-icon class="mr-1"><el-icon-suitcase /></el-icon>
            已安装组件
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
        <el-table-column prop="root" label="根组件" width="100" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="author" label="作者" min-width="120" />
        <el-table-column prop="updateTime" label="更新时间" min-width="180" />
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
              type="success"
              @click="handleInstall(scope.row)"
              :disabled="scope.row.installed"
            >
              {{ scope.row.installed ? '已安装' : '安装' }}
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