<script setup lang="ts">
import type { FormInstance } from 'element-plus';

import { ref } from 'vue';

import { Delete, Edit, Plus } from '@element-plus/icons-vue';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElMessageBox,
  ElTable,
  ElTableColumn,
} from 'element-plus';

import { api } from '#/api/request';
import PageData from '#/components/page/PageData.vue';
import { $t } from '#/locales';

import SysJobModal from './SysJobModal.vue';

const formRef = ref<FormInstance>();
const pageDataRef = ref();
const saveDialog = ref();
const formInline = ref({
  id: '',
});
function search(formEl: FormInstance | undefined) {
  formEl?.validate((valid) => {
    if (valid) {
      pageDataRef.value.setQuery(formInline.value);
    }
  });
}
function reset(formEl: FormInstance | undefined) {
  formEl?.resetFields();
  pageDataRef.value.setQuery({});
}
function showDialog(row: any) {
  saveDialog.value.openDialog({ ...row });
}
function remove(row: any) {
  ElMessageBox.confirm($t('message.deleteAlert'), $t('message.noticeTitle'), {
    confirmButtonText: $t('message.ok'),
    cancelButtonText: $t('message.cancel'),
    type: 'warning',
    beforeClose: (action, instance, done) => {
      if (action === 'confirm') {
        instance.confirmButtonLoading = true;
        api
          .post('/api/v1/sysJob/remove', { id: row.id })
          .then((res) => {
            instance.confirmButtonLoading = false;
            if (res.errorCode === 0) {
              ElMessage.success(res.message);
              reset(formRef.value);
              done();
            }
          })
          .catch(() => {
            instance.confirmButtonLoading = false;
          });
      } else {
        done();
      }
    },
  }).catch(() => {});
}
</script>

<template>
  <div class="page-container">
    <SysJobModal ref="saveDialog" @reload="reset" />
    <ElForm ref="formRef" :inline="true" :model="formInline">
      <ElFormItem :label="$t('sysJob.id')" prop="id">
        <ElInput v-model="formInline.id" :placeholder="$t('sysJob.id')" />
      </ElFormItem>
      <ElFormItem>
        <ElButton @click="search(formRef)" type="primary">
          {{ $t('button.query') }}
        </ElButton>
        <ElButton @click="reset(formRef)">
          {{ $t('button.reset') }}
        </ElButton>
      </ElFormItem>
    </ElForm>
    <div class="handle-div">
      <ElButton @click="showDialog({})" type="primary">
        <ElIcon class="mr-1">
          <Plus />
        </ElIcon>
        {{ $t('button.add') }}
      </ElButton>
    </div>
    <PageData ref="pageDataRef" page-url="/api/v1/sysJob/page" :page-size="10">
      <template #default="{ pageList }">
        <ElTable :data="pageList" border>
          <ElTableColumn prop="deptId" :label="$t('sysJob.deptId')">
            <template #default="{ row }">
              {{ row.deptId }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="jobName" :label="$t('sysJob.jobName')">
            <template #default="{ row }">
              {{ row.jobName }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="jobType" :label="$t('sysJob.jobType')">
            <template #default="{ row }">
              {{ row.jobType }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="jobParams" :label="$t('sysJob.jobParams')">
            <template #default="{ row }">
              {{ row.jobParams }}
            </template>
          </ElTableColumn>
          <ElTableColumn
            prop="cronExpression"
            :label="$t('sysJob.cronExpression')"
          >
            <template #default="{ row }">
              {{ row.cronExpression }}
            </template>
          </ElTableColumn>
          <ElTableColumn
            prop="allowConcurrent"
            :label="$t('sysJob.allowConcurrent')"
          >
            <template #default="{ row }">
              {{ row.allowConcurrent }}
            </template>
          </ElTableColumn>
          <ElTableColumn
            prop="misfirePolicy"
            :label="$t('sysJob.misfirePolicy')"
          >
            <template #default="{ row }">
              {{ row.misfirePolicy }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="options" :label="$t('sysJob.options')">
            <template #default="{ row }">
              {{ row.options }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="status" :label="$t('sysJob.status')">
            <template #default="{ row }">
              {{ row.status }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="created" :label="$t('sysJob.created')">
            <template #default="{ row }">
              {{ row.created }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="remark" :label="$t('sysJob.remark')">
            <template #default="{ row }">
              {{ row.remark }}
            </template>
          </ElTableColumn>
          <ElTableColumn :label="$t('common.handle')" width="150">
            <template #default="{ row }">
              <ElButton @click="showDialog(row)" link type="primary">
                <ElIcon class="mr-1">
                  <Edit />
                </ElIcon>
                {{ $t('button.edit') }}
              </ElButton>
              <ElButton @click="remove(row)" link type="danger">
                <ElIcon class="mr-1">
                  <Delete />
                </ElIcon>
                {{ $t('button.delete') }}
              </ElButton>
            </template>
          </ElTableColumn>
        </ElTable>
      </template>
    </PageData>
  </div>
</template>

<style scoped></style>
