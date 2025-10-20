<script lang="js" setup>
import { ref, watch } from 'vue';

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

const emit = defineEmits(['update:visible']);

// 更完善的Markdown转HTML函数
const markdownToHtml = (markdown) => {
  if (!markdown) return '';
  
  // 先处理代码块
  markdown = markdown.replace(/```([\s\S]*?)```/g, '<pre><code>$1</code></pre>');
  
  return markdown
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')  // 粗体
    .replace(/\*(.*?)\*/g, '<em>$1</em>')              // 斜体
    .replace(/~~(.*?)~~/g, '<del>$1</del>')            // 删除线
    .replace(/`(.*?)`/g, '<code>$1</code>')            // 行内代码
    .replace(/^# (.*$)/gm, '<h1>$1</h1>')              // 标题1
    .replace(/^## (.*$)/gm, '<h2>$1</h2>')             // 标题2
    .replace(/^### (.*$)/gm, '<h3>$1</h3>')            // 标题3
    .replace(/^- (.*$)/gm, '<li>$1</li>')              // 无序列表项
    .replace(/^\d+\. (.*$)/gm, '<li>$1</li>')          // 有序列表项
    .replace(/<li>(.*)<\/li>/g, '<ul>$&</ul>')         // 包装列表
    .replace(/\n/g, '<br>');                           // 换行
};

const handleClose = () => {
  emit('update:visible', false);
};
</script>

<template>
  <el-dialog
    :model-value="visible"
    :title="docData.name || '文档详情'"
    width="1200px"
    top="50px"
    @close="handleClose"
  >
    <div class="doc-viewer">
      <div class="doc-meta">
        <el-descriptions :column="1" size="small" border>
          <el-descriptions-item label="文档名称">{{ docData.name }}</el-descriptions-item>
          <el-descriptions-item label="文档描述">{{ docData.description }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ docData.createTime }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ docData.updateTime }}</el-descriptions-item>
        </el-descriptions>
      </div>
      
      <div class="doc-content">
        <h3>文档内容</h3>
        <div class="content" v-html="markdownToHtml(docData.content)"></div>
      </div>
    </div>
    
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">关闭</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<style scoped>
.doc-viewer {
  max-height: 60vh;
  overflow-y: auto;
}

.doc-meta {
  margin-bottom: 20px;
}

.doc-content h3 {
  margin: 20px 0 10px 0;
  font-size: 18px;
  font-weight: bold;
}

.content {
  padding: 16px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  background-color: #fff;
  min-height: 300px;
  line-height: 1.1;
}

.content :deep(h1) {
  font-size: 28px;
  font-weight: bold;
  margin: 6px 0 2px 0;
  border-bottom: 1px solid #eee;
  padding-bottom: 1px;
  line-height: 1.0;
}

.content :deep(h2) {
  font-size: 24px;
  font-weight: bold;
  margin: 4px 0 1px 0;
  border-bottom: 1px solid #eee;
  padding-bottom: 1px;
  line-height: 1.0;
}

.content :deep(h3) {
  font-size: 20px;
  font-weight: bold;
  margin: 2px 0 1px 0;
  line-height: 1.0;
}

.content :deep(p) {
  margin: 0 0 2px 0;
  line-height: 1.1;
}

.content :deep(strong) {
  font-weight: bold;
}

.content :deep(em) {
  font-style: italic;
}

.content :deep(del) {
  text-decoration: line-through;
}

.content :deep(code) {
  background-color: #f5f5f5;
  padding: 0px 0px;
  border-radius: 1px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 8px;
}

.content :deep(pre) {
  background-color: #f8f8f8;
  border: 1px solid #ddd;
  border-radius: 1px;
  padding: 4px;
  overflow-x: auto;
  margin: 2px 0;
  line-height: 1.0;
}

.content :deep(pre code) {
  background-color: transparent;
  padding: 0;
  border-radius: 0;
  font-size: 7px;
}

.content :deep(ul) {
  padding-left: 16px;
  margin: 1px 0;
}

.content :deep(ol) {
  padding-left: 16px;
  margin: 1px 0;
}

.content :deep(li) {
  list-style-type: disc;
  margin: 0px 0;
}

.content :deep(li)::marker {
  font-size: 0.65em;
}

.content :deep(br) {
  display: inline;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

:deep(.el-dialog__body) {
  padding: 10px 20px;
}
</style>