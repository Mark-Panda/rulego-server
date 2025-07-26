<script lang="js" setup>
const props = defineProps({
  name: {
    type: String,
    default: '',
  },
  description: {
    type: String,
    default: '',
  },
  disabled: {
    type: Boolean,
  },
});

const emit = defineEmits([
  'delete',
  'manage',
  'design',
  'deployment',
  'export',
]);

function deleteHandler() {
  emit('delete');
}

function manageHandler() {
  emit('manage');
}

function designHandler() {
  emit('design');
}

function deploymentHandler() {
  emit('deployment');
}
function exportHandler() {
  emit('export');
}
</script>

<template>
  <div class="h-48 w-full flex-none pb-4 pr-4 md:w-1/2 xl:w-1/4">
    <div
      :class="[
        'h-full w-full rounded-xl border border-[var(--el-border-color)] p-4 transition duration-300 hover:border-primary',
        props.disabled
          ? 'cursor-not-allowed bg-[#f5f7fa] dark:bg-gray-800'
          : 'hover:cursor-pointer hover:shadow-lg',
      ]"
    >
      <div class="flex h-full flex-col">
        <div class="flex flex-none items-center">
          <div class="relative h-10 w-10 flex-none rounded-lg bg-gradient-to-br from-blue-100 to-blue-200 dark:from-blue-900 dark:to-blue-800">
            <div
              class="absolute bottom-[-4px] right-[-4px] flex h-6 w-6 items-center justify-center rounded-full border border-[var(--el-border-color)] bg-white shadow-sm dark:bg-gray-800"
            >
              <el-icon class="text-primary"><el-icon-cpu /></el-icon>
            </div>
          </div>
          <div class="flex-grow overflow-hidden pl-4">
            <div class="truncate font-medium text-gray-800 dark:text-gray-200">{{ props.name }}</div>
            <div class="text-xs text-gray-500 dark:text-gray-400">
              {{ props.disabled ? '未部署' : '已部署' }}
            </div>
          </div>
          <div @click.stop>
            <el-dropdown
              trigger="click"
              :teleported="false"
              :show-arrow="false"
            >
              <el-button icon="el-icon-more-filled" :text="true" class="!text-gray-500 hover:!text-primary"> </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="manageHandler">
                    <el-icon class="mr-1"><el-icon-setting /></el-icon>管理
                  </el-dropdown-item>
                  <el-dropdown-item @click="designHandler">
                    <el-icon class="mr-1"><el-icon-edit /></el-icon>设计
                  </el-dropdown-item>
                  <el-dropdown-item :divided="true" @click="deploymentHandler">
                    <el-icon class="mr-1">
                      <el-icon-video-play v-if="props.disabled" />
                      <el-icon-video-pause v-else />
                    </el-icon>
                    {{ props.disabled ? '部署' : '下线' }}
                  </el-dropdown-item>
                  <el-dropdown-item :divided="true"> 
                    <el-icon class="mr-1"><el-icon-document-copy /></el-icon>复制
                  </el-dropdown-item>
                  <el-dropdown-item @click="exportHandler">
                    <el-icon class="mr-1"><el-icon-download /></el-icon>导出
                  </el-dropdown-item>
                  <el-dropdown-item :divided="true" @click="deleteHandler">
                    <el-icon class="mr-1 text-danger"><el-icon-delete /></el-icon>
                    <span class="text-danger">删除</span>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
        <div class="mt-3 flex-grow overflow-hidden">
          <div class="line-clamp-3 text-sm text-gray-500 dark:text-gray-400">
            {{ props.description || '暂无描述' }}
          </div>
        </div>
        <div class="mt-2 flex-none pt-2 border-t border-gray-100 dark:border-gray-700">
          <div class="flex items-center justify-between text-xs text-gray-500 dark:text-gray-400">
            <div>
              <el-tag size="small" :type="props.disabled ? 'info' : 'success'" effect="light" class="text-xs">
                {{ props.disabled ? '未部署' : '已部署' }}
              </el-tag>
            </div>
            <div class="flex space-x-2">
              <el-button type="primary" size="small" text @click.stop="designHandler">
                <el-icon><el-icon-edit /></el-icon>
                <span>设计</span>
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
