<script lang="js" setup>
import { computed, watch, ref } from 'vue';
import { useFormItem } from 'element-plus';

const props = defineProps({
  modelValue: {
    type: [String, Object],
    default: '',
  },
  desc: {
    type: String,
    default: '',
  },
  placeholder: {
    type: String,
    default: '请输入JSON格式内容...',
  },
  rows: {
    type: Number,
    default: 6,
  },
});

const emit = defineEmits(['update:modelValue']);

const { formItem } = useFormItem();
const jsonError = ref('');

const modelValue = computed({
  get: () => {
    if (typeof props.modelValue === 'object') {
      return JSON.stringify(props.modelValue, null, 2);
    }
    return props.modelValue || '';
  },
  set: (value) => {
    // 验证JSON格式
    if (value.trim()) {
      try {
        JSON.parse(value);
        jsonError.value = '';
        emit('update:modelValue', value);
      } catch (error) {
        jsonError.value = 'JSON格式错误: ' + error.message;
        emit('update:modelValue', value);
      }
    } else {
      jsonError.value = '';
      emit('update:modelValue', value);
    }
  },
});

// 格式化JSON
const formatJson = () => {
  try {
    const parsed = JSON.parse(modelValue.value);
    modelValue.value = JSON.stringify(parsed, null, 2);
    jsonError.value = '';
  } catch (error) {
    jsonError.value = 'JSON格式错误，无法格式化';
  }
};

watch(
  () => props.modelValue,
  () => {
    formItem?.validate?.('change');
  },
);
</script>

<template>
  <div class="relative w-full">
    <div class="relative">
      <el-input 
        v-model="modelValue" 
        type="textarea" 
        :rows="rows"
        :placeholder="placeholder"
        v-bind="$attrs"
        :class="{ 'is-error': jsonError }"
      ></el-input>
      <div class="absolute top-2 right-2">
        <el-button 
          size="small" 
          type="primary" 
          link 
          @click="formatJson"
          title="格式化JSON"
        >
          格式化
        </el-button>
      </div>
    </div>
    <div v-if="jsonError" class="mt-1">
      <el-text size="small" type="danger">
        {{ jsonError }}
      </el-text>
    </div>
    <div v-if="props.desc" class="mt-1">
      <el-text size="small" type="info">
        {{ props.desc }}
      </el-text>
    </div>
  </div>
</template>

<style scoped>
.is-error :deep(.el-textarea__inner) {
  border-color: var(--el-color-danger);
}
</style>