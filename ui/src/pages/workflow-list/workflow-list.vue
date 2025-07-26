<script lang="js" setup>
import { ref, onMounted, toRaw } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessageBox } from 'element-plus';
import * as Api from '@src/api';
import { saveAs } from 'file-saver';

import { WORKFLOW_MENU_KEY } from '@src/constant/workflow';
import AppCard from '@src/pages/workflow-list/app-card.vue';
import CreateAppModal from '@src/pages/workflow-list/create-app-modal.vue';
import ImportDialog from '@src/pages/workflow-list/import-dialog.vue';

const router = useRouter();

const createAppModalRef = ref();
const importDialogRef = ref();

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

async function refreshData() {
  const res = await Api.getRules({
    page: paginationState.value.page,
    size: paginationState.value.size,
    keywords: formState.value.keywords,
  });
  data.value = res.items;
  paginationState.value.total = res.total;

  console.log('data.value', data);
}

function paginationChangeHandler(currentPage, pageSize) {
  paginationState.value.page = currentPage;
  paginationState.value.size = pageSize;
  refreshData();
}

function openDetail(id, tab = WORKFLOW_MENU_KEY.APP_MANAGE) {
  router.push({ path: '/workflow', query: { id, tab } });
}

function openCreateAppModalHandler() {
  createAppModalRef.value.open();
}
//导入
function importHandler() {
  importDialogRef.value.open();
}

function createAppSuccessHandler(id) {
  openDetail(id);
}
function importSuccessHandler(id) {
  refreshData();
}
function deleteHandler(id, name) {
  ElMessageBox.confirm(`删除后将无法恢复 [${name}]，是否继续删除？`, '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
    callback: async (action) => {
      if (action === 'confirm') {
        await Api.deleteRules(id);
        refreshData();
      }
    },
  });
}

function deploymentHandler(id, name, disabled) {
  const actionName = disabled ? '部署' : '下线';

  ElMessageBox.confirm(`确定${actionName}[${name}]吗？`, '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
    callback: async (action) => {
      if (action === 'confirm') {
        await Api.deploymentRules(id, disabled);
        refreshData();
      }
    },
  });
}
//导出成json
const handleDownloadData = (id, data) => {
  const jsonString = JSON.stringify(data, null, 2);
  let blob = new Blob([jsonString], { type: 'text/plain;charset=utf-8' });
  saveAs(blob, id + '.json');
};
onMounted(() => {
  refreshData();
});
</script>

<template>
  <div class="flex h-full flex-col rounded-lg bg-white dark:bg-gray-800 shadow-sm">
    <!-- 页面标题和操作按钮 -->
    <div class="flex flex-none items-center justify-between border-b border-gray-100 dark:border-gray-700 p-4">
      <div class="flex items-center">
        <h1 class="text-xl font-medium text-gray-800 dark:text-gray-200">我的应用</h1>
        <el-tag type="info" effect="plain" class="ml-2">{{ paginationState.total }} 个应用</el-tag>
      </div>
      <div class="flex items-center space-x-2">
        <el-input
          placeholder="搜索应用名称"
          :clearable="true"
          class="w-[220px]"
          v-model="formState.keywords"
          @keyup.enter="refreshData"
        >
          <template #prefix>
            <el-icon><el-icon-search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" @click="openCreateAppModalHandler">
          <el-icon class="mr-1"><el-icon-plus /></el-icon>
          <span>创建应用</span>
        </el-button>
        <el-button @click="importHandler">
          <el-icon class="mr-1"><el-icon-upload /></el-icon>
          <span>导入应用</span>
        </el-button>
      </div>
    </div>
    
    <!-- 应用列表 -->
    <div class="flex-grow overflow-auto p-4">
      <el-scrollbar class="h-full w-full">
        <div v-if="data.length === 0" class="flex h-64 w-full items-center justify-center">
          <div class="text-center">
            <el-icon class="text-4xl text-gray-300 dark:text-gray-600"><el-icon-box /></el-icon>
            <p class="mt-2 text-gray-500 dark:text-gray-400">暂无应用，点击"创建应用"开始</p>
            <el-button class="mt-4" type="primary" @click="openCreateAppModalHandler">
              <el-icon class="mr-1"><el-icon-plus /></el-icon>
              <span>创建应用</span>
            </el-button>
          </div>
        </div>
        <div v-else class="flex flex-wrap content-start">
          <app-card
            v-for="item in data"
            :key="item.ruleChain.id"
            :name="item.ruleChain.name"
            :description="item.ruleChain.additionalInfo.description"
            :disabled="item.ruleChain.disabled"
            @click.prevent="
              item.ruleChain.disabled ? null : openDetail(item.ruleChain.id)
            "
            @delete="deleteHandler(item.ruleChain.id, item.ruleChain.name)"
            @manage="
              openDetail(item.ruleChain.id, WORKFLOW_MENU_KEY.APP_MANAGE)
            "
            @design="
              openDetail(item.ruleChain.id, WORKFLOW_MENU_KEY.APP_DESIGN)
            "
            @export="handleDownloadData(item.ruleChain.id, item)"
            @deployment="
              deploymentHandler(
                item.ruleChain.id,
                item.ruleChain.name,
                item.ruleChain.disabled,
              )
            "
          ></app-card>
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
    
    <!-- 模态框 -->
    <create-app-modal
      ref="createAppModalRef"
      @success="createAppSuccessHandler"
    ></create-app-modal>

    <import-dialog
      ref="importDialogRef"
      @success="importSuccessHandler"
    ></import-dialog>
  </div>
</template>
