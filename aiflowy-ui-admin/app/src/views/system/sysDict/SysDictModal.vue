<script setup lang="ts">
import type { FormInstance } from 'element-plus';

import { onMounted, ref } from 'vue';

import {
  ElButton,
  ElDialog,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
} from 'element-plus';

import { api } from '#/api/request';
import { $t } from '#/locales';

const emit = defineEmits(['reload']);
// vue
onMounted(() => {});
defineExpose({
  openDialog,
});
const saveForm = ref<FormInstance>();
// variables
const dialogVisible = ref(false);
const isAdd = ref(true);
const entity = ref<any>({
  name: '',
  code: '',
  description: '',
  dictType: '',
  sortNo: '',
  status: '',
  options: '',
});
const btnLoading = ref(false);
const rules = ref({
  code: [{ required: true, message: $t('message.required'), trigger: 'blur' }],
});
// functions
function openDialog(row: any) {
  if (row.id) {
    isAdd.value = false;
  }
  entity.value = row;
  dialogVisible.value = true;
}
function save() {
  saveForm.value?.validate((valid) => {
    if (valid) {
      btnLoading.value = true;
      api
        .post(
          isAdd.value ? 'api/v1/sysDict/save' : 'api/v1/sysDict/update',
          entity.value,
        )
        .then((res) => {
          btnLoading.value = false;
          if (res.errorCode === 0) {
            ElMessage.success(res.message);
            emit('reload');
            closeDialog();
          }
        })
        .catch(() => {
          btnLoading.value = false;
        });
    }
  });
}
function closeDialog() {
  saveForm.value?.resetFields();
  isAdd.value = true;
  entity.value = {};
  dialogVisible.value = false;
}
</script>

<template>
  <ElDialog
    v-model="dialogVisible"
    draggable
    :title="isAdd ? $t('button.add') : $t('button.edit')"
    :before-close="closeDialog"
    :close-on-click-modal="false"
  >
    <ElForm
      label-width="120px"
      ref="saveForm"
      :model="entity"
      status-icon
      :rules="rules"
    >
      <ElFormItem prop="name" :label="$t('sysDict.name')">
        <ElInput v-model.trim="entity.name" />
      </ElFormItem>
      <ElFormItem prop="code" :label="$t('sysDict.code')">
        <ElInput v-model.trim="entity.code" />
      </ElFormItem>
      <ElFormItem prop="description" :label="$t('sysDict.description')">
        <ElInput v-model.trim="entity.description" />
      </ElFormItem>
      <ElFormItem prop="dictType" :label="$t('sysDict.dictType')">
        <ElInput v-model.trim="entity.dictType" />
      </ElFormItem>
      <ElFormItem prop="sortNo" :label="$t('sysDict.sortNo')">
        <ElInput v-model.trim="entity.sortNo" />
      </ElFormItem>
      <ElFormItem prop="status" :label="$t('sysDict.status')">
        <ElInput v-model.trim="entity.status" />
      </ElFormItem>
      <ElFormItem prop="options" :label="$t('sysDict.options')">
        <ElInput v-model.trim="entity.options" />
      </ElFormItem>
    </ElForm>
    <template #footer>
      <ElButton @click="closeDialog">
        {{ $t('button.cancel') }}
      </ElButton>
      <ElButton
        type="primary"
        @click="save"
        :loading="btnLoading"
        :disabled="btnLoading"
      >
        {{ $t('button.save') }}
      </ElButton>
    </template>
  </ElDialog>
</template>

<style scoped></style>
