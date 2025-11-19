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

import SysAccountModal from './SysAccountModal.vue';

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
          .post('/api/v1/sysAccount/remove', { id: row.id })
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
    <SysAccountModal ref="saveDialog" @reload="reset" />
    <ElForm ref="formRef" :inline="true" :model="formInline">
      <ElFormItem :label="$t('sysAccount.id')" prop="id">
        <ElInput
          v-model="formInline.id"
          :placeholder="`${$t('sysAccount.id')}`"
        />
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
    <PageData
      ref="pageDataRef"
      page-url="/api/v1/sysAccount/page"
      :page-size="10"
    >
      <template #default="{ pageList }">
        <ElTable :data="pageList" border>
          <ElTableColumn prop="deptId" :label="$t('sysAccount.deptId')">
            <template #default="{ row }">
              {{ row.deptId }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="loginName" :label="$t('sysAccount.loginName')">
            <template #default="{ row }">
              {{ row.loginName }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="password" :label="$t('sysAccount.password')">
            <template #default="{ row }">
              {{ row.password }}
            </template>
          </ElTableColumn>
          <ElTableColumn
            prop="accountType"
            :label="$t('sysAccount.accountType')"
          >
            <template #default="{ row }">
              {{ row.accountType }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="nickname" :label="$t('sysAccount.nickname')">
            <template #default="{ row }">
              {{ row.nickname }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="mobile" :label="$t('sysAccount.mobile')">
            <template #default="{ row }">
              {{ row.mobile }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="email" :label="$t('sysAccount.email')">
            <template #default="{ row }">
              {{ row.email }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="avatar" :label="$t('sysAccount.avatar')">
            <template #default="{ row }">
              {{ row.avatar }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="dataScope" :label="$t('sysAccount.dataScope')">
            <template #default="{ row }">
              {{ row.dataScope }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="deptIdList" :label="$t('sysAccount.deptIdList')">
            <template #default="{ row }">
              {{ row.deptIdList }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="status" :label="$t('sysAccount.status')">
            <template #default="{ row }">
              {{ row.status }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="created" :label="$t('sysAccount.created')">
            <template #default="{ row }">
              {{ row.created }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="remark" :label="$t('sysAccount.remark')">
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
