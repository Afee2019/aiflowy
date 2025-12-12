<script setup lang="ts">
import type { FilesCardProps } from 'vue-element-plus-x/types/FilesCard';

import { ref } from 'vue';
import { Attachments } from 'vue-element-plus-x';

import { ElMessage } from 'element-plus';

const senderRef = ref();
const showHeaderFlog = ref(false);

type SelfFilesCardProps = FilesCardProps & {
  id?: number | string;
};

const files = ref<SelfFilesCardProps[]>([]);
defineExpose( {
  init(firstFile: File, fileList: FileList) {
    console.log('firstFile', firstFile);
    showHeaderFlog.value = true;
    files.value = [
      {
        id: 0,
        uid: `${firstFile.name}_${firstFile.size}`,
        name: firstFile.name,
        fileSize: firstFile.size,
        imgFile: firstFile,
        showDelIcon: true,
        imgVariant: 'square',
      },
    ];
  },
});

function closeHeader() {
  showHeaderFlog.value = false;
  senderRef.value.closeHeader();
}

function handlePasteFile(firstFile: File, fileList: FileList) {
  showHeaderFlog.value = true;
  senderRef.value.openHeader();
  const fileArray = [...fileList];

  fileArray.forEach((file, index) => {
    files.value.push({
      id: index,
      uid: `${index}_${file.name}_${file.size}`,
      name: file.name,
      fileSize: file.size,
      imgFile: file,
      showDelIcon: true,
      imgVariant: 'square',
    });
  });
}

async function handleHttpRequest(options: any) {
  const formData = new FormData();
  formData.append('file', options.file);
  ElMessage.info('上传中...');

  setTimeout(() => {
    const res = {
      message: '文件上传成功',
      fileName: options.file.name,
      uid: options.file.uid,
      fileSize: options.file.size,
      imgFile: options.file,
    };
    files.value.push({
      id: files.value.length,
      uid: res.uid,
      name: res.fileName,
      fileSize: res.fileSize,
      imgFile: res.imgFile,
      showDelIcon: true,
      imgVariant: 'square',
    });

    ElMessage.success('上传成功');
  }, 1000);
}

function handleDeleteCard(item: SelfFilesCardProps) {
  files.value = files.value.filter((items: any) => items.id !== item.id);
  ElMessage.success('删除成功');
}
</script>

<template>
  <div class="header-self-wrap">
    <Attachments
      :items="files"
      :http-request="handleHttpRequest"
      @delete-card="handleDeleteCard"
    />
  </div>
</template>

<style scoped>
.header-self-wrap {
  display: flex;
  flex-direction: row;
  padding: 16px;
  width: 100%;
  overflow: auto;
  .header-self-title {
    width: 100%;
    display: flex;
    height: 30px;
    align-items: center;
    justify-content: space-between;
    padding-bottom: 8px;
  }
  .header-self-content {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 20px;
    color: #626aef;
    font-weight: 600;
  }
}
</style>
