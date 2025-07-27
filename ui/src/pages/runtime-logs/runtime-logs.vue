<script lang="js" setup>
import { onMounted, ref, reactive } from 'vue';
import dayjs from 'dayjs';
import beautify from 'js-beautify';
import * as Api from '@src/api';
import { ElMessage } from 'element-plus';
import JsonEditor from '@src/components/json-editor/json-editor.vue';

const tableData = ref([]);
const expandedRows = ref(new Set());
const detailData = ref({});
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
      current: paginationState.value.page,
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
        
        // 清空详情数据
        detailData.value = {};
        
        // 处理主表数据
        const mainData = res.items.map(item => {
          // 计算总执行时间
          const executionTime = item.endTs && item.startTs ? item.endTs - item.startTs : 0;
          
          // 计算节点数量
          const nodeCount = item.logs && Array.isArray(item.logs) ? item.logs.length : 0;
          
          // 检查是否有错误
          let hasError = false;
          if (item.logs && Array.isArray(item.logs)) {
            hasError = item.logs.some(log => log.err && log.err.length > 0);
          }
          
          // 构建主表记录
          return {
            id: item.id,
            ts: item.startTs,
            chainId: item.ruleChain?.id || '',
            chainName: item.ruleChain?.name || '',
            isRoot: item.ruleChain?.root || false,
            executionTime: executionTime,
            nodeCount: nodeCount,
            hasError: hasError,
            startTs: item.startTs,
            endTs: item.endTs
          };
        });
        
        // 处理详情数据
        res.items.forEach(item => {
          if (item.logs && Array.isArray(item.logs)) {
            // 为每个主记录创建详情数据数组
            detailData.value[item.id] = item.logs.map(log => ({
              nodeId: log.nodeId,
              relationType: log.relationType,
              msg: log.inMsg || {},
              outMsg: log.outMsg || {},
              err: log.err,
              startTs: log.startTs,
              endTs: log.endTs,
              executionTime: log.endTs && log.startTs ? log.endTs - log.startTs : 0
            }));
          }
        });
        
        console.log('处理后的主表数据:', mainData);
        console.log('处理后的详情数据:', detailData.value);
        
        tableData.value = mainData;
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

function toggleRowExpansion(row) {
  const rowId = row.id;
  if (expandedRows.value.has(rowId)) {
    expandedRows.value.delete(rowId);
  } else {
    expandedRows.value.add(rowId);
  }
}

function isRowExpanded(row) {
  return expandedRows.value.has(row.id);
}

function handleSizeChange(val) {
  paginationState.value.size = val;
  refreshTableData();
}

function handleCurrentChange(val) {
  paginationState.value.page = val;
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
      
      <!-- 主表 -->
      <el-table
        v-loading="loading"
        height="calc(100vh - 350px)"
        size="small"
        :data="tableData"
        :border="true"
        stripe
        row-key="id"
      >
        <!-- 展开按钮列 -->
        <el-table-column type="expand" width="50">
          <template #default="props">
            <!-- 子表 -->
            <div class="p-4 bg-gray-50">
              <h4 class="text-lg font-medium mb-3">节点执行详情</h4>
              <el-table
                :data="detailData[props.row.id] || []"
                size="small"
                border
                stripe
                style="width: 100%"
              >
                <el-table-column prop="nodeId" label="节点ID" min-width="140" align="left">
                  <template #default="scope">
                    <el-tooltip
                      v-if="scope.row.nodeId"
                      effect="dark"
                      :content="scope.row.nodeId"
                      placement="top"
                    >
                      <span>
                        {{
                          scope.row.nodeId && scope.row.nodeId.length > 12
                            ? scope.row.nodeId.substring(0, 12) + '...'
                            : scope.row.nodeId
                        }}
                      </span>
                    </el-tooltip>
                  </template>
                </el-table-column>
                <el-table-column prop="relationType" label="关系类型" min-width="120" align="left" />
                <el-table-column label="消息ID" min-width="180" align="left">
                  <template #default="scope">
                    <el-tooltip
                      v-if="scope.row.msg && scope.row.msg.id"
                      effect="dark"
                      :content="scope.row.msg.id"
                      placement="top"
                    >
                      <span>
                        {{
                          scope.row.msg && scope.row.msg.id && scope.row.msg.id.length > 16
                            ? scope.row.msg.id.substring(0, 16) + '...'
                            : (scope.row.msg ? scope.row.msg.id : '')
                        }}
                      </span>
                    </el-tooltip>
                  </template>
                </el-table-column>
                <el-table-column label="消息类型" min-width="150" align="left">
                  <template #default="scope">
                    <el-tooltip
                      v-if="scope.row.msg && scope.row.msg.type"
                      effect="dark"
                      :content="scope.row.msg.type"
                      placement="top"
                    >
                      <span>
                        {{
                          scope.row.msg && scope.row.msg.type && scope.row.msg.type.length > 12
                            ? scope.row.msg.type.substring(0, 12) + '...'
                            : (scope.row.msg ? scope.row.msg.type : '')
                        }}
                      </span>
                    </el-tooltip>
                  </template>
                </el-table-column>
                <el-table-column label="输入数据" width="80" align="center">
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
                <el-table-column label="输出数据" width="80" align="center">
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
                  width="80"
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
                <el-table-column label="错误" width="80" align="center">
                  <template #default="scope">
                    <el-tooltip
                      v-if="scope.row.err"
                      effect="dark"
                      :content="scope.row.err"
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
                <el-table-column label="执行时间(ms)" min-width="120" align="left">
                  <template #default="scope">
                    <span>{{ scope.row.executionTime || '-' }}</span>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </template>
        </el-table-column>
        
        <!-- 主表列 -->
        <el-table-column
          prop="ts"
          label="执行时间"
          min-width="180"
          :formatter="tsFormatter"
          align="left"
        />
        <el-table-column label="工作流" min-width="180" align="left">
          <template #default="scope">
            <div class="flex items-center">
              <el-tooltip
                v-if="scope.row.chainId"
                effect="dark"
                :content="scope.row.chainName || scope.row.chainId"
                placement="top"
              >
                <span>
                  {{
                    scope.row.chainName || 
                    (scope.row.chainId && scope.row.chainId.length > 12
                      ? scope.row.chainId.substring(0, 12) + '...'
                      : scope.row.chainId)
                  }}
                </span>
              </el-tooltip>
              <el-tag
                v-if="scope.row.isRoot"
                size="small"
                type="success"
                class="ml-2"
              >
                根规则链
              </el-tag>
              <el-tag
                v-else
                size="small"
                type="info"
                class="ml-2"
              >
                子规则链
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="nodeCount" label="节点数" min-width="100" align="left" />
        <el-table-column label="状态" width="100" align="center">
          <template #default="scope">
            <el-tag :type="scope.row.hasError ? 'danger' : 'success'">
              {{ scope.row.hasError ? '失败' : '成功' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="执行时间(ms)" min-width="120" align="left">
          <template #default="scope">
            <span>{{ scope.row.executionTime || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" align="center">
          <template #default="scope">
            <el-button
              size="small"
              @click="toggleRowExpansion(scope.row)"
            >
              查看详情
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
