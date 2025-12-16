<script setup lang="ts">
import type { UploadProps } from 'element-plus';

import { ref } from 'vue';

import { useAppConfig } from '@aiflowy/hooks';
import { useAccessStore } from '@aiflowy/stores';

import { UploadFilled } from '@element-plus/icons-vue';
import { ElIcon, ElUpload } from 'element-plus';

const props = defineProps({
  action: {
    type: String,
    default: '/api/v1/commons/upload',
  },
  visible: {
    type: Boolean,
    default: true,
  },
});

const emit = defineEmits(['success', 'onChange']);
const accessStore = useAccessStore();
const headers = ref({
  'aiflowy-token': accessStore.accessToken,
});
const { apiURL } = useAppConfig(import.meta.env, import.meta.env.PROD);

// 核心：获取ElUpload组件实例
const uploadRef = ref<InstanceType<typeof ElUpload>>();

// 上传成功回调
const handleSuccess: UploadProps['onSuccess'] = (response) => {
  emit('success', response.data.path);
};

// 文件状态变化回调
const handleChange: UploadProps['onChange'] = (file, fileList) => {
  emit('onChange', file, fileList);
};

// 暴露给父组件的方法：手动触发文件选择
const triggerFileSelect = () => {
  if (uploadRef.value) {
    // 调用ElUpload内部的上传按钮点击事件
    const uploadInput = uploadRef.value.$el.querySelector('input[type="file"]');
    if (uploadInput) {
      uploadInput.click(); // 触发原生文件选择框
    }
  }
};

// 对外暴露方法（父组件可通过ref调用）
defineExpose({
  triggerFileSelect,
});
</script>

<template>
  <!-- 给ElUpload添加ref引用 -->
  <ElUpload
    ref="uploadRef"
    class="upload-demo"
    drag
    :headers="headers"
    :action="`${apiURL}${props.action}`"
    :on-success="handleSuccess"
    :on-change="handleChange"
    multiple
    :style="{ display: props.visible ? 'block' : 'none' }"
  >
    <ElIcon class="el-icon--upload"><UploadFilled /></ElIcon>
    <div class="el-upload__text">
      Drop file here or <em>click to upload</em>
    </div>
  </ElUpload>
</template>
