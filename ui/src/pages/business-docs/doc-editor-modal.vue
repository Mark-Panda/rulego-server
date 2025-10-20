<script lang="js" setup>
import { ref, reactive, watch } from 'vue';
import MarkdownEditor from '@src/components/markdown-editor/markdown-editor.vue';

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  docData: {
    type: Object,
    default: () => ({})
  }
});

const emit = defineEmits(['update:visible', 'save']);

const formRef = ref(null);
const formState = reactive({
  id: '',
  name: '',
  content: '',
  description: ''
});

const rules = {
  name: [
    { required: true, message: '请输入文档名称', trigger: 'blur' }
  ]
};

watch(() => props.visible, (val) => {
  if (val) {
    // 如果有传入的文档数据，则初始化表单
    if (props.docData && props.docData.id) {
      formState.id = props.docData.id;
      formState.name = props.docData.name;
      formState.content = props.docData.content || '';
      formState.description = props.docData.description || '';
    } else {
      // 否则重置表单
      formState.id = '';
      formState.name = '';
      formState.content = '';
      formState.description = '';
    }
  }
});

const handleClose = () => {
  emit('update:visible', false);
};

const handleSubmit = async () => {
  try {
    await formRef.value.validate();
    // 发送保存事件
    emit('save', { ...formState });
    handleClose();
  } catch (error) {
    console.error('表单验证失败:', error);
  }
};
</script>

<template>
  <el-dialog
    :model-value="visible"
    :title="formState.id ? '编辑文档' : '新建文档'"
    width="1200px"
    top="50px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formState"
      :rules="rules"
      label-width="80px"
    >
      <el-form-item label="文档名称" prop="name">
        <el-input
          v-model="formState.name"
          placeholder="请输入文档名称"
        />
      </el-form-item>
      
      <el-form-item label="文档描述" prop="description">
        <el-input
          v-model="formState.description"
          type="textarea"
          :rows="2"
          placeholder="请输入文档描述"
        />
      </el-form-item>
      
      <el-form-item label="文档内容" prop="content">
        <markdown-editor
          v-model="formState.content"
          placeholder="请输入文档内容，支持Markdown语法"
          height="400px"
          width="100%"
          :split-view="true"
        />
      </el-form-item>
    </el-form>
    
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="handleSubmit">保存</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<style scoped>
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

:deep(.el-dialog__body) {
  padding: 10px 20px;
}
</style>