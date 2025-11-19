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

import SysLogModal from './SysLogModal.vue';

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
          .post('/api/v1/sysLog/remove', { id: row.id })
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
    <SysLogModal ref="saveDialog" @reload="reset" />
    <ElForm ref="formRef" :inline="true" :model="formInline">
      <ElFormItem :label="$t('sysLog.id')" prop="id">
        <ElInput v-model="formInline.id" :placeholder="$t('sysLog.id')" />
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
    <PageData ref="pageDataRef" page-url="/api/v1/sysLog/page" :page-size="10">
      <template #default="{ pageList }">
        <ElTable :data="pageList" border>
          <ElTableColumn prop="accountId" :label="$t('sysLog.accountId')">
            <template #default="{ row }">
              {{ row.accountId }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="actionName" :label="$t('sysLog.actionName')">
            <template #default="{ row }">
              {{ row.actionName }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="actionType" :label="$t('sysLog.actionType')">
            <template #default="{ row }">
              {{ row.actionType }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="actionClass" :label="$t('sysLog.actionClass')">
            <template #default="{ row }">
              {{ row.actionClass }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="actionMethod" :label="$t('sysLog.actionMethod')">
            <template #default="{ row }">
              {{ row.actionMethod }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="actionUrl" :label="$t('sysLog.actionUrl')">
            <template #default="{ row }">
              {{ row.actionUrl }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="actionIp" :label="$t('sysLog.actionIp')">
            <template #default="{ row }">
              {{ row.actionIp }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="actionParams" :label="$t('sysLog.actionParams')">
            <template #default="{ row }">
              {{ row.actionParams }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="actionBody" :label="$t('sysLog.actionBody')">
            <template #default="{ row }">
              {{ row.actionBody }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="status" :label="$t('sysLog.status')">
            <template #default="{ row }">
              {{ row.status }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="created" :label="$t('sysLog.created')">
            <template #default="{ row }">
              {{ row.created }}
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
