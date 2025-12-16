<!-- src/components/MarkdownRenderer.vue -->
<script setup lang="ts">
import { computed } from 'vue';

import DOMPurify from 'dompurify';
import hljs from 'highlight.js';
import markdownit from 'markdown-it';

import 'highlight.js/styles/github-dark.css';

// 定义 props
const props = defineProps<{
  content: string;
  removePrefix?: string;
}>();

const md = markdownit({
  html: true,
  linkify: true,
  typographer: true,
  breaks: true,
  highlight: (str: string, lang: string) => {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return `<pre class="hljs"><code>${hljs.highlight(str, { language: lang, ignoreIllegals: true }).value}</code></pre>`;
      } catch (error) {
        console.error('代码高亮失败:', error);
      }
    }
    return `<pre class="hljs"><code>${hljs.highlightAuto(str).value}</code></pre>`;
  },
});

// 解析 Markdown 为安全的 HTML
const renderedHtml = computed(() => {
  if (!props.content) return '';
  // 移除指定前缀
  const pureContent = props.removePrefix
    ? props.content.replace(new RegExp(`^${props.removePrefix}\\s*`, 'i'), '')
    : props.content;
  // 解析 Markdown
  let html = md.render(pureContent);
  // XSS 过滤（关键）
  html = DOMPurify.sanitize(html);
  return html;
});
</script>

<template>
  <div class="markdown-renderer" v-html="renderedHtml"></div>
</template>

<style scoped>
/* Markdown 基础样式 */
.markdown-renderer {
  line-height: 1.8;
  color: #333;
  font-size: 14px;
}

/* 代码块样式 */
.markdown-renderer :deep(pre) {
  padding: 12px;
  border-radius: 8px;
  overflow-x: auto;
  background: #1e293b !important;
  margin: 16px 0;
}

.markdown-renderer :deep(code) {
  padding: 2px 4px;
  border-radius: 4px;
  background: #f1f5f9;
  color: #e11d48;
  font-size: 13px;
}

/* 标题样式 */
.markdown-renderer :deep(h1) {
  font-size: 24px;
  font-weight: 600;
  margin: 20px 0 12px;
}
.markdown-renderer :deep(h2) {
  font-size: 20px;
  font-weight: 600;
  margin: 16px 0 10px;
}
.markdown-renderer :deep(h3) {
  font-size: 18px;
  font-weight: 600;
  margin: 14px 0 8px;
}

/* 列表/引用块 */
.markdown-renderer :deep(ul),
.markdown-renderer :deep(ol) {
  padding-left: 24px;
  margin: 8px 0;
}
.markdown-renderer :deep(blockquote) {
  border-left: 4px solid #3b82f6;
  padding: 8px 16px;
  background: #f0f7ff;
  margin: 12px 0;
}

/* 链接样式 */
.markdown-renderer :deep(a) {
  color: #3b82f6;
  text-decoration: underline;
}
.markdown-renderer :deep(a:hover) {
  color: #2563eb;
}

/* 表格样式（可选） */
.markdown-renderer :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin: 16px 0;
}
.markdown-renderer :deep(th),
.markdown-renderer :deep(td) {
  border: 1px solid #e2e8f0;
  padding: 8px 12px;
}
.markdown-renderer :deep(th) {
  background: #f8fafc;
  font-weight: 600;
}
</style>
