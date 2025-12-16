<script setup lang="ts">
import type { PropType } from 'vue';

import type { llmType } from '#/api';

import { Minus, Setting } from '@element-plus/icons-vue';
import { ElIcon, ElImage } from 'element-plus';

defineProps({
  llmList: {
    type: Array as PropType<llmType[]>,
    default: () => [],
  },
  icon: {
    type: String,
    default: '',
  },
});
const emit = defineEmits(['deleteLlm', 'editLlm']);
const handleDeleteLlm = (id: string) => {
  emit('deleteLlm', id);
};
const handleEditLlm = (id: string) => {
  emit('editLlm', id);
};
</script>

<template>
  <div v-for="llm in llmList" :key="llm.id" class="container">
    <div class="llm-item">
      <div class="start">
        <ElImage :src="icon" style="width: 24px; height: 24px" />
        <div>
          {{ llm.title }}
        </div>
      </div>
      <div class="end">
        <ElIcon
          size="16"
          @click="handleEditLlm(llm.id)"
          style="cursor: pointer"
        >
          <Setting />
        </ElIcon>
        <ElIcon
          size="16"
          @click="handleDeleteLlm(llm.id)"
          style="cursor: pointer"
        >
          <Minus />
        </ElIcon>
      </div>
    </div>
  </div>
</template>

<style scoped>
.llm-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 40px;
}
.container {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding-left: 18px;
  padding-right: 18px;
}
.start {
  display: flex;
  align-items: center;
  gap: 12px;
}
.end {
  display: flex;
  align-items: center;
  gap: 12px;
}
</style>
