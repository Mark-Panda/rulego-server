<script lang="js" setup>
import { onMounted, ref, reactive } from 'vue';
import dayjs from 'dayjs';
import beautify from 'js-beautify';
import * as Api from '@src/api';
import { ElMessage } from 'element-plus';
import JsonEditor from '@src/components/json-editor/json-editor.vue';

const tableData = ref([]);
const paginationState = ref({
  page: 1,
  size: 10,
  pageSizes: [10, 20, 30, 40, 50],
  total: 0,
});
const dialogTitle = ref('');
const isDialogOpen = ref(false);
const dialogData = ref('');
const loading = ref(false);
const searchForm = ref({
  chainId: '',
  startTime: '',
  endTime: '',
});

async function refreshTableData() {
  loading.value = true;
  try {
    const params = {
      size: paginationState.value.size,
      page: paginationState.value.page,
    };
    
    if (searchForm.value.chainId) {
      params.chainId = searchForm.value.chainId;
    }
    
    if (searchForm.value.startTime) {
      params.startTime = dayjs(searchForm.value.startTime).format('YYYY-MM-DD HH:mm:ss');
    }
    
    if (searchForm.value.endTime) {
      params.endTime = dayjs(searchForm.value.endTime).format('YYYY-MM-DD HH:mm:ss');
    }
    
    console.log('请求参数:', params);
    
    try {
      // 使用API函数获取真实数据
      const res = await Api.getRuntimeLogs(params);
      console.log('API返回数据:', res);
      
      // 处理返回的数据结构
      if (res && res.items && Array.isArray(res.items)) {
        console.log('找到items数组，处理日志数据');
        
        // 提取所有日志项
        const processedData = [];
        
        // 遍历每个工作流执行记录
        res.items.forEach(item => {
          if (item.logs && Array.isArray(item.logs)) {
            // 遍历每个节点的日志
            item.logs.forEach(log => {
              // 构建日志记录
              processedData.push({
                ts: log.startTs,
                chainId: item.ruleChain?.id || '',
                chainName: item.ruleChain?.name || '',
                nodeId: log.nodeId,
                flowType: log.relationType,
                msg: log.inMsg || {},
                outMsg: log.outMsg || {},
                relationType: log.relationType,
                err: log.err,
                startTs: log.startTs,
                endTs: log.endTs
              });
            });
          }
        });
        
        console.log('处理后的日志数据:', processedData);
        tableData.value = processedData;
        paginationState.value.total = res.total || 0;
      } else {
        console.warn('未识别的数据结构:', res);
        tableData.value = [];
        paginationState.value.total = 0;
      }
      
      // 如果数据为空，显示提示
      if (tableData.value.length === 0) {
        ElMessage.info('没有找到符合条件的日志数据');
      }
    } catch (apiError) {
      console.error('API调用失败:', apiError);
      ElMessage.error('获取日志数据失败，请检查网络连接或API配置');
      tableData.value = [];
      paginationState.value.total = 0;
    }
  } catch (error) {
    console.error('获取运行日志失败:', error);
    tableData.value = [];
    paginationState.value.total = 0;
    ElMessage.error('处理日志数据时发生错误');
  } finally {
    loading.value = false;
  }
}

function paginationChangeHandler(currentPage, pageSize) {
  paginationState.value.page = currentPage;
  paginationState.value.size = pageSize;
  refreshTableData();
}

function tsFormatter(row, column) {
  return dayjs(row.ts).format('YYYY-MM-DD HH:mm:ss');
}

function openDialog() {
  isDialogOpen.value = true;
}

function closeDialog() {
  isDialogOpen.value = false;
}

function showDataHandler(msg) {
  dialogTitle.value = '数据';
  try {
    // 检查数据是否为字符串，如果是则尝试解析为JSON
    let dataContent = msg.data;
    if (typeof dataContent === 'string') {
      try {
        dataContent = JSON.parse(dataContent);
      } catch (e) {
        // 如果解析失败，保持原始字符串
      }
    }
    
    // 美化JSON输出
    if (dataContent && typeof dataContent === 'object') {
      dialogData.value = beautify.js(JSON.stringify(dataContent), { indent_size: 2 });
    } else {
      dialogData.value = String(dataContent || '');
    }
    
    console.log('显示数据:', dialogData.value);
  } catch (error) {
    console.error('格式化数据失败:', error);
    dialogData.value = String(msg.data || '');
  }
  openDialog();
}

function showMetadataHandler(msg) {
  dialogTitle.value = '元数据';
  try {
    // 检查元数据是否为字符串，如果是则尝试解析为JSON
    let metadataContent = msg.metadata;
    if (typeof metadataContent === 'string') {
      try {
        metadataContent = JSON.parse(metadataContent);
      } catch (e) {
        // 如果解析失败，保持原始字符串
      }
    }
    
    // 美化JSON输出
    if (metadataContent && typeof metadataContent === 'object') {
      dialogData.value = beautify.js(JSON.stringify(metadataContent), { indent_size: 2 });
    } else {
      dialogData.value = String(metadataContent || '');
    }
    
    console.log('显示元数据:', dialogData.value);
  } catch (error) {
    console.error('格式化元数据失败:', error);
    dialogData.value = String(msg.metadata || '');
  }
  openDialog();
}

function showOutputDataHandler(msg) {
  dialogTitle.value = '输出数据';
  try {
    // 检查数据是否为字符串，如果是则尝试解析为JSON
    let dataContent = msg.data;
    if (typeof dataContent === 'string') {
      try {
        dataContent = JSON.parse(dataContent);
      } catch (e) {
        // 如果解析失败，保持原始字符串
      }
    }
    
    // 美化JSON输出
    if (dataContent && typeof dataContent === 'object') {
      dialogData.value = beautify.js(JSON.stringify(dataContent), { indent_size: 2 });
    } else {
      dialogData.value = String(dataContent || '');
    }
    
    console.log('显示输出数据:', dialogData.value);
  } catch (error) {
    console.error('格式化输出数据失败:', error);
    dialogData.value = String(msg.data || '');
  }
  openDialog();
}

function showErrorHandler(err) {
  dialogTitle.value = '错误';
  dialogData.value = err || '';
  openDialog();
}

function handleSearch() {
  paginationState.value.page = 1;
  refreshTableData();
}

function handleReset() {
  searchForm.value = {
    chainId: '',
    startTime: '',
    endTime: '',
  };
  paginationState.value.page = 1;
  refreshTableData();
}

onMounted(() => {
  console.log('运行日志组件已挂载');
  refreshTableData();
  
  // 检查表格数据是否正确加载
  setTimeout(() => {
    console.log('表格数据检查:', tableData.value);
    if (tableData.value.length > 0) {
      console.log('表格第一行数据:', tableData.value[0]);
    } else {
      console.warn('表格数据为空');
    }
  }, 1000);
});
</script>

<template>
  <div class="p-4">
    <div class="mb-6">
      <h2 class="text-2xl font-bold mb-2">运行日志</h2>
      <p class="text-gray-500 dark:text-gray-400">查看工作流运行日志，监控工作流执行情况</p>
    </div>
    
    <div class="bg-white dark:bg-gray-800 p-4 rounded-md shadow mb-6">
      <el-form :model="searchForm" inline>
        <el-form-item label="工作流ID">
          <el-input v-model="searchForm.chainId" placeholder="请输入工作流ID" clearable></el-input>
        </el-form-item>
        <el-form-item label="开始时间">
          <el-date-picker
            v-model="searchForm.startTime"
            type="datetime"
            placeholder="选择开始时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            clearable
          ></el-date-picker>
        </el-form-item>
        <el-form-item label="结束时间">
          <el-date-picker
            v-model="searchForm.endTime"
            type="datetime"
            placeholder="选择结束时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            clearable
          ></el-date-picker>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <div class="bg-white dark:bg-gray-800 rounded-md shadow">
      <div class="flex justify-between items-center p-4 border-b border-gray-200 dark:border-gray-700">
        <div class="font-bold text-lg">运行日志列表</div>
        <el-button @click="refreshTableData">
          <el-icon class="mr-1"><el-icon-refresh-right /></el-icon>
          刷新
        </el-button>
      </div>
      
      <el-table
        v-loading="loading"
        height="calc(100vh - 350px)"
        size="small"
        :data="tableData"
        :border="true"
        stripe
        @row-click="(row) => console.log('点击行:', row)"
      >
        <el-table-column
          prop="ts"
          label="事件时间"
          width="160"
          :formatter="tsFormatter"
        />
        <el-table-column label="工作流" width="160">
          <template #default="scope">
            <el-tooltip
              v-if="scope.row.chainId"
              effect="dark"
              :content="scope.row.chainName || scope.row.chainId"
              placement="top"
            >
              <span>
                {{
                  scope.row.chainName || 
                  (scope.row.chainId && scope.row.chainId.length > 10
                    ? scope.row.chainId.substring(0, 10) + '...'
                    : scope.row.chainId)
                }}
              </span>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column prop="nodeId" label="节点ID" width="120">
          <template #default="scope">
            <el-tooltip
              v-if="scope.row.nodeId"
              effect="dark"
              :content="scope.row.nodeId"
              placement="top"
            >
              <span>
                {{
                  scope.row.nodeId && scope.row.nodeId.length > 10
                    ? scope.row.nodeId.substring(0, 10) + '...'
                    : scope.row.nodeId
                }}
              </span>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column prop="relationType" label="关系类型" width="90" />
        <el-table-column label="消息ID" width="130">
          <template #default="scope">
            <el-tooltip
              v-if="scope.row.msg && scope.row.msg.id"
              effect="dark"
              :content="scope.row.msg.id"
              placement="top"
            >
              <span>
                {{
                  scope.row.msg && scope.row.msg.id && scope.row.msg.id.length > 14
                    ? scope.row.msg.id.substring(0, 14) + '...'
                    : (scope.row.msg ? scope.row.msg.id : '')
                }}
              </span>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column label="消息类型" width="130">
          <template #default="scope">
            <el-tooltip
              v-if="scope.row.msg && scope.row.msg.type"
              effect="dark"
              :content="scope.row.msg.type"
              placement="top"
            >
              <span>
                {{
                  scope.row.msg && scope.row.msg.type && scope.row.msg.type.length > 10
                    ? scope.row.msg.type.substring(0, 10) + '...'
                    : (scope.row.msg ? scope.row.msg.type : '')
                }}
              </span>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column label="输入数据" width="60" align="center">
          <template #default="scope">
            <el-tooltip 
              v-if="scope.row.msg && scope.row.msg.data"
              effect="dark" 
              content="查看输入数据" 
              placement="top"
            >
              <el-button
                @click="showDataHandler(scope.row.msg)"
                :link="true"
              >
                <el-icon><el-icon-more-filled /></el-icon>
              </el-button>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column label="输出数据" width="60" align="center">
          <template #default="scope">
            <el-tooltip 
              v-if="scope.row.outMsg && scope.row.outMsg.data"
              effect="dark" 
              content="查看输出数据" 
              placement="top"
            >
              <el-button
                @click="showDataHandler(scope.row.outMsg)"
                :link="true"
              >
                <el-icon><el-icon-more-filled /></el-icon>
              </el-button>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column
          label="元数据"
          width="60"
          align="center"
        >
          <template #default="scope">
            <el-tooltip 
              v-if="scope.row.msg && scope.row.msg.metadata"
              effect="dark" 
              content="查看元数据" 
              placement="top"
            >
              <el-button
                @click="showMetadataHandler(scope.row.msg)"
                :link="true"
              >
                <el-icon><el-icon-more-filled /></el-icon>
              </el-button>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column label="错误" width="80">
          <template #default="scope">
            <el-tooltip
              v-if="scope.row.err"
              effect="dark"
              content="查看错误"
              placement="top"
            >
              <el-button
                @click="showErrorHandler(scope.row.err)"
                :link="true"
                type="danger"
                :disabled="!scope.row.err"
              >
                <el-icon><el-icon-warning /></el-icon>
              </el-button>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column label="执行时间(ms)" width="100" align="right">
          <template #default="scope">
            <span>{{ scope.row.endTs && scope.row.startTs ? scope.row.endTs - scope.row.startTs : '-' }}</span>
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
          @change="paginationChangeHandler"
        ></el-pagination>
      </div>
    </div>
    
    <el-dialog
      :append-to-body="false"
      :destroy-on-close="true"
      :close-on-click-modal="true"
      :close-on-press-escape="true"
      :draggable="true"
      top="10px"
      :before-close="closeDialog"
      v-model="isDialogOpen"
      :title="dialogTitle"
      width="60%"
    >
      <json-editor v-model="dialogData"></json-editor>
      <template #footer>
        <div class="flex justify-end">
          <el-button @click="closeDialog">关闭</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>
