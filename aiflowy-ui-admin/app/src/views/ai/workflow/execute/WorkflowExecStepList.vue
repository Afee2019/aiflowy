<script setup lang="ts">
import type { FormInstance } from 'element-plus';

import { onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { ArrowLeft } from '@element-plus/icons-vue';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElTable,
  ElTableColumn,
  ElTag,
} from 'element-plus';

import PageData from '#/components/page/PageData.vue';
import { $t } from '#/locales';
import { useDictStore } from '#/store';

const router = useRouter();
const $route = useRoute();
onMounted(() => {
  initDict();
});
const formRef = ref<FormInstance>();
const pageDataRef = ref();
const formInline = ref({
  nodeName: '',
});
const dictStore = useDictStore();
function initDict() {
  dictStore.fetchDictionary('dataStatus');
}
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
function getTagType(row: any) {
  switch (row.status) {
    case 1: {
      return 'primary';
    }
    case 6: {
      return 'warning';
    }
    case 10: {
      return 'danger';
    }
    case 20: {
      return 'success';
    }
    case 21: {
      return 'danger';
    }
    default: {
      return 'info';
    }
  }
}
</script>

<template>
  <div class="page-container border-border border">
    <div class="mb-3">
      <ElButton :icon="ArrowLeft" @click="router.back()">
        {{ $t('button.back') }}
      </ElButton>
    </div>
    <ElForm ref="formRef" :inline="true" :model="formInline">
      <ElFormItem class="w-full max-w-[300px]" prop="nodeName">
        <ElInput
          v-model="formInline.nodeName"
          :placeholder="$t('aiWorkflowRecordStep.nodeName')"
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
    <div class="handle-div"></div>
    <PageData
      ref="pageDataRef"
      page-url="/api/v1/workflowExecStep/page"
      :page-size="10"
      :extra-query-params="{
        recordId: $route.query.recordId,
        sortKey: 'startTime',
        sortType: 'asc',
      }"
    >
      <template #default="{ pageList }">
        <ElTable :data="pageList" border>
          <ElTableColumn
            show-overflow-tooltip
            prop="execKey"
            :label="$t('aiWorkflowRecordStep.execKey')"
          >
            <template #default="{ row }">
              {{ row.execKey }}
            </template>
          </ElTableColumn>
          <ElTableColumn
            show-overflow-tooltip
            prop="nodeId"
            :label="$t('aiWorkflowRecordStep.nodeId')"
          >
            <template #default="{ row }">
              {{ row.nodeId }}
            </template>
          </ElTableColumn>
          <ElTableColumn
            prop="nodeName"
            :label="$t('aiWorkflowRecordStep.nodeName')"
          >
            <template #default="{ row }">
              {{ row.nodeName }}
            </template>
          </ElTableColumn>
          <ElTableColumn
            show-overflow-tooltip
            prop="input"
            :label="$t('aiWorkflowRecordStep.input')"
          >
            <template #default="{ row }">
              {{ row.input }}
            </template>
          </ElTableColumn>
          <ElTableColumn
            show-overflow-tooltip
            prop="output"
            :label="$t('aiWorkflowRecordStep.output')"
          >
            <template #default="{ row }">
              {{ row.output }}
            </template>
          </ElTableColumn>
          <ElTableColumn
            prop="execTime"
            :label="$t('aiWorkflowRecordStep.execTime')"
          >
            <template #default="{ row }">
              {{ row.execTime || '-' }} ms
            </template>
          </ElTableColumn>
          <ElTableColumn
            prop="status"
            :label="$t('aiWorkflowRecordStep.status')"
          >
            <template #default="{ row }">
              <ElTag :type="getTagType(row)">
                {{ $t(`aiWorkflowRecordStep.status${row.status}`) }}
              </ElTag>
            </template>
          </ElTableColumn>
          <ElTableColumn
            show-overflow-tooltip
            prop="errorInfo"
            :label="$t('aiWorkflowRecordStep.errorInfo')"
          >
            <template #default="{ row }">
              {{ row.errorInfo }}
            </template>
          </ElTableColumn>
        </ElTable>
      </template>
    </PageData>
  </div>
</template>

<style scoped></style>
