<script lang="js" setup>
import { computed, watch, onMounted, ref, inject } from 'vue';
import { useFormItem } from 'element-plus';
import { first } from 'lodash-es';

const props = defineProps({
  modelValue: {
    type: [String, Number, Boolean, Object, Array],
    default: '',
  },
  options: {
    type: Array,
    default: () => [],
  },
  desc: {
    type: String,
    default: '',
  },
  allowCreate: {
    type: Boolean,
    default: false,
  },
  filterable: {
    type: Boolean,
    default: false,
  },
  loadData: {
    type: Function,
    default: null,
  },
  fieldConfig: {
    type: Object,
    default: null,
  },
});

const emit = defineEmits(['update:modelValue']);

const { formItem } = useFormItem();
const dynamicOptions = ref([]);

// 尝试从父组件注入LogicFlow实例和当前节点模型
const lf = inject('logicFlow', null);
const currentNodeModel = inject('currentNodeModel', null);
const fieldConfig = inject('fieldConfig', null);

// 计算最终的选项列表
const finalOptions = computed(() => {
  // 如果有loadData函数，优先使用动态选项
  const hasLoadDataFunc = props.loadData || (props.fieldConfig?.componentProps?.loadData);
  if (hasLoadDataFunc) {
    return dynamicOptions.value;
  }
  // 否则使用静态选项
  return props.options;
});

const modelValue = computed({
  get: () => {
    return props.modelValue;
  },
  set: (value) => {
    emit('update:modelValue', value);
  },
});

function initValue() {
  const firstValue = first(finalOptions.value)?.value;
  const hasModelValue = finalOptions.value.some(
    (item) => item.value === props.modelValue,
  );
  
  console.log('initValue 调用:', {
    finalOptionsLength: finalOptions.value.length,
    firstValue,
    hasModelValue,
    currentModelValue: props.modelValue,
    allowCreate: props.allowCreate
  });
  
  // 如果没有选项，不设置任何值
  if (finalOptions.value.length === 0) {
    return;
  }
  
  // 如果允许创建自定义值，且当前值不为空，则保持当前值
  if (props.allowCreate && props.modelValue) {
    return;
  }
  
  // 否则，如果当前值在选项中存在，保持当前值；否则设置为第一个选项的值
  if (hasModelValue) {
    modelValue.value = props.modelValue;
  } else if (firstValue !== undefined) {
    modelValue.value = firstValue;
  }
}

// 调用loadData函数加载动态选项
async function loadDynamicOptions() {
  const loadDataFunc = props.loadData || (props.fieldConfig?.componentProps?.loadData);
  console.log('loadDynamicOptions 调用:', {
    hasLoadDataFunc: !!loadDataFunc,
    hasLf: !!lf,
    hasCurrentNodeModel: !!currentNodeModel,
    fieldConfig: props.fieldConfig,
    injectedFieldConfig: fieldConfig
  });
  
  if (loadDataFunc && typeof loadDataFunc === 'function') {
    try {
      const configToUse = props.fieldConfig || fieldConfig;
      // 创建一个field对象，包含component属性
      const fieldForLoadData = {
        ...configToUse,
        component: {
          options: [],
          // 保留原有的component属性（如果是对象的话）
          ...(typeof configToUse?.component === 'object' ? configToUse.component : {})
        }
      };
      console.log('调用loadData函数:', { lf: !!lf, currentNodeModel: !!currentNodeModel, fieldForLoadData });
      console.log('loadDataFunc 函数详情:', loadDataFunc);
      console.log('传递的参数:', {
        param1: lf,
        param2: currentNodeModel,
        param3: fieldForLoadData,
        param3Type: typeof fieldForLoadData,
        param3Keys: fieldForLoadData ? Object.keys(fieldForLoadData) : null
      });
      const options = await loadDataFunc(lf, currentNodeModel, fieldForLoadData);
      console.log('loadData返回的选项:', options);
      if (Array.isArray(options)) {
        dynamicOptions.value = options;
      }
    } catch (error) {
      console.error('加载动态选项失败:', error);
    }
  }
}

watch(
  () => props.modelValue,
  () => {
    formItem?.validate?.('change');
  },
);

watch(
  () => finalOptions.value,
  () => {
    initValue();
  },
);

// 监听LogicFlow和currentNodeModel的变化，重新加载选项
watch(
  [() => lf, () => currentNodeModel],
  async () => {
    if (lf && currentNodeModel) {
      console.log('LogicFlow或currentNodeModel发生变化，重新加载选项');
      await loadDynamicOptions();
      initValue();
    }
  },
  { deep: true }
);

onMounted(async () => {
  console.log('select-form-item mounted, props:', props);
  console.log('注入的值:', {
    lf: !!lf,
    currentNodeModel: !!currentNodeModel,
    fieldConfig: !!fieldConfig,
    lfType: typeof lf,
    currentNodeModelType: typeof currentNodeModel
  });
  console.log(`=== loadData 函数详细检查 [${props.fieldConfig?.name || 'unknown'}] ===`);
  console.log('props.loadData:', props.loadData);
  console.log('props.fieldConfig:', props.fieldConfig);
  if (props.fieldConfig) {
    console.log('fieldConfig keys:', Object.keys(props.fieldConfig));
    console.log('fieldConfig.componentProps:', props.fieldConfig.componentProps);
    if (props.fieldConfig.componentProps) {
      console.log('componentProps keys:', Object.keys(props.fieldConfig.componentProps));
      console.log('componentProps.loadData:', props.fieldConfig.componentProps.loadData);
    }
  }
  console.log(`=== 检查结束 [${props.fieldConfig?.name || 'unknown'}] ===`);
  
  await loadDynamicOptions();
  initValue();
  console.log('select-form-item 初始化完成, finalOptions:', finalOptions.value);
});
</script>

<template>
  <div class="relative w-full">
    <div>
      <el-select 
        v-model="modelValue" 
        :allow-create="allowCreate"
        :filterable="filterable"
        v-bind="$attrs"
      >
        <el-option
          v-for="item in finalOptions"
          :key="item.value"
          :label="item.label"
          :value="item.value"
        />
      </el-select>
    </div>
    <div v-if="props.desc">
      <el-text size="small" type="info">
        {{ props.desc }}
      </el-text>
    </div>
  </div>
</template>
