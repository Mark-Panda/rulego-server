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
  // 新增：是否允许添加新选项
  allowAddOption: {
    type: Boolean,
    default: false,
  },
  // 新增：添加选项的回调函数
  onAddOption: {
    type: Function,
    default: null,
  },
});

const emit = defineEmits(['update:modelValue', 'add-option']);

const { formItem } = useFormItem();
const dynamicOptions = ref([]);
const userAddedOptions = ref([]); // 用户添加的选项
const isAddingOption = ref(false); // 是否正在添加选项
const newOptionLabel = ref(''); // 新选项的标签

// 尝试从父组件注入LogicFlow实例和当前节点模型
const lf = inject('logicFlow', null);
const currentNodeModel = inject('currentNodeModel', null);
const fieldConfig = inject('fieldConfig', null);

// 计算最终的选项列表
const finalOptions = computed(() => {
  let baseOptions = [];
  
  // 如果有loadData函数，优先使用动态选项
  const hasLoadDataFunc = props.loadData || (props.fieldConfig?.componentProps?.loadData);
  if (hasLoadDataFunc) {
    baseOptions = dynamicOptions.value;
  } else {
    // 否则使用静态选项
    baseOptions = props.options;
  }
  
  // 合并用户添加的选项
  const allOptions = [...baseOptions, ...userAddedOptions.value];
  
  // 如果允许添加选项，在列表末尾添加"添加新选项"的特殊选项
  if (props.allowAddOption && !isAddingOption.value) {
    allOptions.push({
      value: '__ADD_NEW_OPTION__',
      label: '+ 添加新选项',
      isAddOption: true
    });
  }
  
  return allOptions;
});

const modelValue = computed({
  get: () => {
    return props.modelValue;
  },
  set: (value) => {
    // 如果选择的是"添加新选项"，则开始添加流程
    if (value === '__ADD_NEW_OPTION__') {
      startAddingOption();
      return;
    }
    emit('update:modelValue', value);
  },
});

// 开始添加新选项
function startAddingOption() {
  isAddingOption.value = true;
  newOptionLabel.value = '';
  // 重置选择值，避免显示"添加新选项"
  emit('update:modelValue', props.modelValue);
}

// 确认添加新选项
function confirmAddOption() {
  if (!newOptionLabel.value.trim()) {
    return;
  }
  
  const newOption = {
    value: newOptionLabel.value.trim(),
    label: newOptionLabel.value.trim(),
    isUserAdded: true
  };
  
  // 检查是否已存在相同的选项
  const exists = finalOptions.value.some(option => 
    option.value === newOption.value && !option.isAddOption
  );
  
  if (exists) {
    // 可以显示提示信息
    console.warn('选项已存在:', newOption.value);
    cancelAddOption();
    return;
  }
  
  // 添加到用户选项列表
  userAddedOptions.value.push(newOption);
  
  // 设置为当前选中值
  emit('update:modelValue', newOption.value);
  
  // 通知父组件
  if (props.onAddOption) {
    props.onAddOption(newOption);
  }
  emit('add-option', newOption);
  
  // 结束添加状态
  isAddingOption.value = false;
  newOptionLabel.value = '';
}

// 取消添加新选项
function cancelAddOption() {
  isAddingOption.value = false;
  newOptionLabel.value = '';
}

// 处理输入框的键盘事件
function handleInputKeydown(event) {
  if (event.key === 'Enter') {
    event.preventDefault();
    confirmAddOption();
  } else if (event.key === 'Escape') {
    event.preventDefault();
    cancelAddOption();
  }
}

function initValue() {
  const firstValue = first(finalOptions.value.filter(item => !item.isAddOption))?.value;
  
  // 检查当前值是否在选项中存在（支持多选）
  let hasModelValue = false;
  if (Array.isArray(props.modelValue)) {
    // 多选情况：检查数组中的每个值是否都在选项中
    hasModelValue = props.modelValue.length > 0 && props.modelValue.every(value => 
      finalOptions.value.some(item => item.value === value && !item.isAddOption)
    );
  } else {
    // 单选情况
    hasModelValue = finalOptions.value.some(
      (item) => item.value === props.modelValue && !item.isAddOption,
    );
  }
  
  console.log('initValue 调用:', {
    finalOptionsLength: finalOptions.value.length,
    firstValue,
    hasModelValue,
    currentModelValue: props.modelValue,
    isArray: Array.isArray(props.modelValue),
    allowCreate: props.allowCreate
  });
  
  // 如果没有选项，不设置任何值
  if (finalOptions.value.filter(item => !item.isAddOption).length === 0) {
    return;
  }
  
  // 如果允许创建自定义值，且当前值不为空，则保持当前值
  if (props.allowCreate && props.modelValue) {
    return;
  }
  
  // 如果当前值已经存在且有效，保持当前值
  if (hasModelValue) {
    return; // 不需要重新设置，保持原值
  }
  
  // 只有在当前值无效时才设置默认值
  if (firstValue !== undefined) {
    // 对于多选组件，如果当前值是空数组或无效，不设置默认值
    // 对于单选组件，设置第一个选项为默认值
    if (!Array.isArray(props.modelValue)) {
      modelValue.value = firstValue;
    }
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
  
  // 强制调用 loadDynamicOptions，即使 lf 或 currentNodeModel 为空
  console.log('强制调用 loadDynamicOptions...');
  await loadDynamicOptions();
  initValue();
  console.log('select-form-item 初始化完成, finalOptions:', finalOptions.value);
  
  // 如果初始时没有 lf 或 currentNodeModel，设置一个定时器重试
  if (!lf || !currentNodeModel) {
    console.log('lf 或 currentNodeModel 为空，设置定时器重试...');
    setTimeout(async () => {
      console.log('定时器重试，当前注入值:', {
        lf: !!lf,
        currentNodeModel: !!currentNodeModel
      });
      if (lf && currentNodeModel) {
        console.log('重试调用 loadDynamicOptions...');
        await loadDynamicOptions();
        initValue();
      }
    }, 1000);
  }
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
          :class="{ 'add-option-item': item.isAddOption }"
        />
        
        <!-- 添加新选项的输入框 -->
        <div v-if="isAddingOption" class="px-3 py-2 border-t border-gray-200">
          <div class="flex items-center space-x-2">
            <el-input
              v-model="newOptionLabel"
              placeholder="输入新选项名称"
              size="small"
              @keydown="handleInputKeydown"
              ref="newOptionInput"
              class="flex-1"
            />
            <el-button 
              type="primary" 
              size="small" 
              @click="confirmAddOption"
              :disabled="!newOptionLabel.trim()"
            >
              确认
            </el-button>
            <el-button 
              size="small" 
              @click="cancelAddOption"
            >
              取消
            </el-button>
          </div>
        </div>
      </el-select>
    </div>
    <div v-if="props.desc">
      <el-text size="small" type="info">
        {{ props.desc }}
      </el-text>
    </div>
  </div>
</template>

<style scoped>
.add-option-item {
  color: #409eff;
  font-style: italic;
}
</style>
