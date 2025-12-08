<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';

import { sortNodes } from '@aiflowy/utils';

import { ElCard, ElCol, ElRow } from 'element-plus';

import { api } from '#/api/request';
import { $t } from '#/locales';
import ExecResult from '#/views/ai/workflow/components/ExecResult.vue';
import WorkflowForm from '#/views/ai/workflow/components/WorkflowForm.vue';
import WorkflowSteps from '#/views/ai/workflow/components/WorkflowSteps.vue';

onMounted(async () => {
  pageLoading.value = true;
  await Promise.all([getWorkflowInfo(workflowId.value), getRunningParams()]);
  pageLoading.value = false;
});
const pageLoading = ref(false);
const route = useRoute();
const workflowId = ref(route.query.id);
const workflowInfo = ref<any>({});
const runParams = ref<any>(null);
const initState = ref(false);
const tinyFlowData = ref<any>(null);
const workflowForm = ref();
async function getWorkflowInfo(workflowId: any) {
  api.get(`/api/v1/aiWorkflow/detail?id=${workflowId}`).then((res) => {
    workflowInfo.value = res.data;
    tinyFlowData.value = workflowInfo.value.content
      ? JSON.parse(workflowInfo.value.content)
      : {};
  });
}
async function getRunningParams() {
  api
    .get(`/api/v1/aiWorkflow/getRunningParameters?id=${workflowId.value}`)
    .then((res) => {
      runParams.value = res.data;
    });
}
function onSubmit() {
  initState.value = !initState.value;
}
function resumeChain(data: any) {
  workflowForm.value?.resume(data);
}
const chainInfo = ref<any>(null);
function onAsyncExecute(info: any) {
  chainInfo.value = info;
}
</script>

<template>
  <div v-loading="pageLoading" class="container p-2.5">
    <ElRow :gutter="10">
      <ElCol :span="10">
        <ElCard
          shadow="never"
          class="mb-2.5"
          body-style="height: 300px;overflow-y: auto"
        >
          <div class="mb-2.5 font-semibold">
            {{ $t('aiWorkflow.params') }}：
          </div>
          <WorkflowForm
            v-if="runParams && tinyFlowData"
            ref="workflowForm"
            :workflow-id="workflowId"
            :workflow-params="runParams"
            :on-submit="onSubmit"
            :on-async-execute="onAsyncExecute"
            :tiny-flow-data="tinyFlowData"
          />
        </ElCard>
        <ElCard shadow="never" body-style="height: 300px;overflow-y: auto">
          <div class="mb-2.5 font-semibold">{{ $t('aiWorkflow.steps') }}：</div>
          <WorkflowSteps
            v-if="tinyFlowData"
            :workflow-id="workflowId"
            :node-json="sortNodes(tinyFlowData)"
            :init-signal="initState"
            :polling-data="chainInfo"
            @resume="resumeChain"
          />
        </ElCard>
      </ElCol>
      <ElCol :span="14">
        <ElCard shadow="never" body-style="height: 612px;overflow-y: auto">
          <div class="mb-2.5 mt-2.5 font-semibold">
            {{ $t('aiWorkflow.result') }}：
          </div>
          <ExecResult
            v-if="tinyFlowData"
            :workflow-id="workflowId"
            :node-json="sortNodes(tinyFlowData)"
            :init-signal="initState"
            :polling-data="chainInfo"
          />
        </ElCard>
      </ElCol>
    </ElRow>
  </div>
</template>

<style scoped>
.container {
  background: var(--el-bg-color-page);
}
</style>
