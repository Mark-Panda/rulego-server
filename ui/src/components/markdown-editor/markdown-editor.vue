<script lang="js" setup>
import { ref, onMounted } from 'vue';

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: ''
  },
  height: {
    type: String,
    default: '300px'
  },
  width: {
    type: String,
    default: '100%'
  },
  // 是否启用双屏显示
  splitView: {
    type: Boolean,
    default: true
  },
  // 是否可调整大小
  resizable: {
    type: Boolean,
    default: true
  }
});

const emit = defineEmits(['update:modelValue']);

const textareaRef = ref(null);
const editorRef = ref(null);
const currentHeight = ref(props.height);
const currentWidth = ref(props.width);

const handleInput = (event) => {
  emit('update:modelValue', event.target.value);
};

// 简单的Markdown转HTML函数
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

// 处理拖拽调整大小
const handleResize = (e, direction) => {
  if (!props.resizable) return;
  
  e.preventDefault();
  
  const startX = e.clientX;
  const startY = e.clientY;
  const startWidth = parseInt(currentWidth.value, 10);
  const startHeight = parseInt(currentHeight.value, 10);
  
  const doDrag = (dragEvent) => {
    if (direction === 'horizontal') {
      // 调整宽度
      const newWidth = startWidth + (dragEvent.clientX - startX);
      if (newWidth > 300) { // 最小宽度300px
        currentWidth.value = newWidth + 'px';
      }
    } else {
      // 调整高度
      const newHeight = startHeight + (dragEvent.clientY - startY);
      if (newHeight > 100) { // 最小高度100px
        currentHeight.value = newHeight + 'px';
      }
    }
  };
  
  const stopDrag = () => {
    document.removeEventListener('mousemove', doDrag);
    document.removeEventListener('mouseup', stopDrag);
  };
  
  document.addEventListener('mousemove', doDrag);
  document.addEventListener('mouseup', stopDrag);
};

onMounted(() => {
  if (textareaRef.value) {
    textareaRef.value.value = props.modelValue;
  }
});
</script>

<template>
  <div 
    class="markdown-editor" 
    ref="editorRef"
    :style="{ width: currentWidth }"
  >
    <div v-if="splitView" class="editor-split-view" :style="{ height: currentHeight }">
      <div class="editor-pane">
        <textarea
          ref="textareaRef"
          :placeholder="placeholder"
          :value="modelValue"
          @input="handleInput"
          class="editor-textarea"
        ></textarea>
      </div>
      <div class="preview-pane">
        <div
          class="editor-preview"
          v-html="markdownToHtml(modelValue)"
        ></div>
      </div>
    </div>
    <div v-else class="editor-content" :style="{ height: currentHeight }">
      <textarea
        ref="textareaRef"
        :placeholder="placeholder"
        :value="modelValue"
        @input="handleInput"
        class="editor-textarea"
      ></textarea>
    </div>
    
    <!-- 调整大小的手柄 -->
    <div 
      v-if="resizable"
      class="resize-handle resize-handle-bottom" 
      @mousedown="handleResize($event, 'vertical')"
    ></div>
    
    <!-- 调整宽度的手柄 -->
    <div 
      v-if="resizable"
      class="resize-handle resize-handle-right" 
      @mousedown="handleResize($event, 'horizontal')"
    ></div>
  </div>
</template>

<style scoped>
.markdown-editor {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
  position: relative;
  min-width: 300px;
}

.editor-split-view {
  display: flex;
  position: relative;
}

.editor-pane {
  flex: 1;
  border-right: 1px solid #dcdfe6;
}

.preview-pane {
  flex: 1;
}

.editor-content {
  position: relative;
}

.editor-textarea {
  width: 100%;
  height: 100%;
  padding: 12px;
  border: none;
  outline: none;
  resize: none;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.5;
}

.editor-preview {
  padding: 12px;
  height: 100%;
  overflow-y: auto;
  background-color: #fafafa;
  line-height: 1.1;
}

.editor-preview :deep(h1) {
  font-size: 24px;
  font-weight: bold;
  margin: 4px 0 1px 0;
  line-height: 1.0;
}

.editor-preview :deep(h2) {
  font-size: 20px;
  font-weight: bold;
  margin: 2px 0 1px 0;
  line-height: 1.0;
}

.editor-preview :deep(h3) {
  font-size: 18px;
  font-weight: bold;
  margin: 1px 0 1px 0;
  line-height: 1.0;
}

.editor-preview :deep(p) {
  margin: 0 0 1px 0;
  line-height: 1.1;
}

.editor-preview :deep(strong) {
  font-weight: bold;
}

.editor-preview :deep(em) {
  font-style: italic;
}

.editor-preview :deep(del) {
  text-decoration: line-through;
}

.editor-preview :deep(code) {
  background-color: #f5f5f5;
  padding: 0px 0px;
  border-radius: 1px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 8px;
}

.editor-preview :deep(pre) {
  background-color: #f8f8f8;
  border: 1px solid #ddd;
  border-radius: 1px;
  padding: 2px;
  overflow-x: auto;
  margin: 2px 0;
  line-height: 1.0;
}

.editor-preview :deep(pre code) {
  background-color: transparent;
  padding: 0;
  border-radius: 0;
  font-size: 7px;
}

.editor-preview :deep(ul) {
  padding-left: 14px;
  margin: 1px 0;
}

.editor-preview :deep(ol) {
  padding-left: 14px;
  margin: 1px 0;
}

.editor-preview :deep(li) {
  list-style-type: disc;
  margin: 0px 0;
}

.editor-preview :deep(li)::marker {
  font-size: 0.65em;
}

.editor-preview :deep(br) {
  display: inline;
}

.resize-handle {
  position: absolute;
  background-color: #dcdfe6;
  transition: background-color 0.2s;
}

.resize-handle-bottom {
  bottom: 0;
  left: 0;
  right: 0;
  height: 5px;
  cursor: ns-resize;
}

.resize-handle-right {
  top: 0;
  bottom: 0;
  right: 0;
  width: 5px;
  cursor: ew-resize;
}

.resize-handle:hover {
  background-color: #409eff;
}
</style>