<script setup lang="ts">
import type { Ref } from 'vue';

import { defineEmits, nextTick, ref } from 'vue';

import { $t } from '@aiflowy/locales';

import { ElButton, ElDialog, ElForm, ElFormItem, ElInput } from 'element-plus';

interface BasicFormItem {
  key: string;
  description: string;
}

const emit = defineEmits(['success']);
const dialogVisible = ref(false);
const openDialog = (data: BasicFormItem[]) => {
  nextTick(() => {
    basicFormRef.value?.resetFields();
  });
  if (data.length === 0) {
    basicForm.value = [
      { key: '1', description: '' },
      { key: '2', description: '' },
      { key: '3', description: '' },
      { key: '4', description: '' },
      { key: '5', description: '' },
    ];
  } else {
    basicForm.value = data;
  }

  dialogVisible.value = true;
};

const basicForm: Ref<BasicFormItem[]> = ref([
  { key: '1', description: '' },
  { key: '2', description: '' },
  { key: '3', description: '' },
  { key: '4', description: '' },
  { key: '5', description: '' },
]);

const basicFormRef = ref();

defineExpose({
  openDialog(data: BasicFormItem[]) {
    openDialog(data);
  },
});

const handleConfirm = () => {
  basicFormRef.value?.validate((valid: boolean) => {
    if (valid) {
      emit('success', basicForm.value);
      dialogVisible.value = false;
    }
  });
};
</script>

<template>
  <ElDialog
    v-model="dialogVisible"
    :title="$t('button.add')"
    width="700"
    align-center
  >
    <ElForm
      ref="basicFormRef"
      style="width: 100%; margin-top: 20px"
      :model="basicForm"
      label-width="auto"
    >
      <template v-for="(item, index) in basicForm" :key="item.key">
        <ElFormItem
          :label="`${$t('bot.problemPresupposition')}${item.key}`"
          :prop="`${index}.description`"
          label-position="right"
        >
          <ElInput v-model="item.description" />
        </ElFormItem>
      </template>
    </ElForm>

    <template #footer>
      <div class="dialog-footer">
        <ElButton @click="dialogVisible = false">
          {{ $t('button.cancel') }}
        </ElButton>
        <ElButton type="primary" @click="handleConfirm">
          {{ $t('button.confirm') }}
        </ElButton>
      </div>
    </template>
  </ElDialog>
</template>

<style scoped></style>
