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
  parentId: '',
  ancestors: '',
  deptName: '',
  deptCode: '',
  sortNo: '',
  status: '',
  remark: '',
});
const btnLoading = ref(false);
const rules = ref({
  parentId: [
    { required: true, message: $t('message.required'), trigger: 'blur' },
  ],
  deptName: [
    { required: true, message: $t('message.required'), trigger: 'blur' },
  ],
  status: [
    { required: true, message: $t('message.required'), trigger: 'blur' },
  ],
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
          isAdd.value ? 'api/v1/sysDept/save' : 'api/v1/sysDept/update',
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
      <ElFormItem prop="parentId" :label="$t('sysDept.parentId')">
        <ElInput v-model.trim="entity.parentId" />
      </ElFormItem>
      <ElFormItem prop="ancestors" :label="$t('sysDept.ancestors')">
        <ElInput v-model.trim="entity.ancestors" />
      </ElFormItem>
      <ElFormItem prop="deptName" :label="$t('sysDept.deptName')">
        <ElInput v-model.trim="entity.deptName" />
      </ElFormItem>
      <ElFormItem prop="deptCode" :label="$t('sysDept.deptCode')">
        <ElInput v-model.trim="entity.deptCode" />
      </ElFormItem>
      <ElFormItem prop="sortNo" :label="$t('sysDept.sortNo')">
        <ElInput v-model.trim="entity.sortNo" />
      </ElFormItem>
      <ElFormItem prop="status" :label="$t('sysDept.status')">
        <ElInput v-model.trim="entity.status" />
      </ElFormItem>
      <ElFormItem prop="remark" :label="$t('sysDept.remark')">
        <ElInput v-model.trim="entity.remark" />
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
