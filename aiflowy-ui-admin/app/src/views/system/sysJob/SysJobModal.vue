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
  deptId: '',
  jobName: '',
  jobType: '',
  jobParams: '',
  cronExpression: '',
  allowConcurrent: '',
  misfirePolicy: '',
  options: '',
  status: '',
  remark: '',
});
const btnLoading = ref(false);
const rules = ref({
  deptId: [
    { required: true, message: $t('message.required'), trigger: 'blur' },
  ],
  jobName: [
    { required: true, message: $t('message.required'), trigger: 'blur' },
  ],
  jobType: [
    { required: true, message: $t('message.required'), trigger: 'blur' },
  ],
  cronExpression: [
    { required: true, message: $t('message.required'), trigger: 'blur' },
  ],
  allowConcurrent: [
    { required: true, message: $t('message.required'), trigger: 'blur' },
  ],
  misfirePolicy: [
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
          isAdd.value ? 'api/v1/sysJob/save' : 'api/v1/sysJob/update',
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
      <ElFormItem prop="deptId" :label="$t('sysJob.deptId')">
        <ElInput v-model.trim="entity.deptId" />
      </ElFormItem>
      <ElFormItem prop="jobName" :label="$t('sysJob.jobName')">
        <ElInput v-model.trim="entity.jobName" />
      </ElFormItem>
      <ElFormItem prop="jobType" :label="$t('sysJob.jobType')">
        <ElInput v-model.trim="entity.jobType" />
      </ElFormItem>
      <ElFormItem prop="jobParams" :label="$t('sysJob.jobParams')">
        <ElInput v-model.trim="entity.jobParams" />
      </ElFormItem>
      <ElFormItem prop="cronExpression" :label="$t('sysJob.cronExpression')">
        <ElInput v-model.trim="entity.cronExpression" />
      </ElFormItem>
      <ElFormItem prop="allowConcurrent" :label="$t('sysJob.allowConcurrent')">
        <ElInput v-model.trim="entity.allowConcurrent" />
      </ElFormItem>
      <ElFormItem prop="misfirePolicy" :label="$t('sysJob.misfirePolicy')">
        <ElInput v-model.trim="entity.misfirePolicy" />
      </ElFormItem>
      <ElFormItem prop="options" :label="$t('sysJob.options')">
        <ElInput v-model.trim="entity.options" />
      </ElFormItem>
      <ElFormItem prop="status" :label="$t('sysJob.status')">
        <ElInput v-model.trim="entity.status" />
      </ElFormItem>
      <ElFormItem prop="remark" :label="$t('sysJob.remark')">
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
