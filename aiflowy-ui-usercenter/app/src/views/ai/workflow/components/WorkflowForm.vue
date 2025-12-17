<script setup lang="ts">
import type { FormInstance } from 'element-plus';

import { computed, onUnmounted, ref } from 'vue';

import { Position } from '@element-plus/icons-vue';
import { ElButton, ElForm, ElFormItem } from 'element-plus';

import { api } from '#/api/request';
import { $t } from '#/locales';
import WorkflowFormItem from '#/views/ai/workflow/components/WorkflowFormItem.vue';

export type WorkflowFormProps = {
  onAsyncExecute?: (values: any) => void;
  onSubmit?: (values: any) => void;
  tinyFlowData: any;
  workflowId: any;
  workflowParams: any;
};
const props = withDefaults(defineProps<WorkflowFormProps>(), {
  onExecuting: () => {
    console.warn('no execute method');
  },
  onSubmit: () => {
    console.warn('no submit method');
  },
  onAsyncExecute: () => {
    console.warn('no async execute method');
  },
});
defineExpose({
  resume,
});
const runForm = ref<FormInstance>();
const runParams = ref<any>({});
const submitLoading = ref(false);
const parameters = computed(() => {
  return props.workflowParams.parameters;
});
const executeId = ref('');
function resume(data: any) {
  data.executeId = executeId.value;
  submitLoading.value = true;
  api.post('/userCenter/aiWorkflow/resume', data).then((res) => {
    if (res.errorCode === 0) {
      startPolling(executeId.value);
    }
  });
}
function submitV2() {
  runForm.value?.validate((valid) => {
    if (valid) {
      const data = {
        id: props.workflowId,
        variables: {
          ...runParams.value,
        },
      };
      props.onSubmit?.(runParams.value);
      submitLoading.value = true;
      api.post('/userCenter/aiWorkflow/runAsync', data).then((res) => {
        if (res.errorCode === 0 && res.data) {
          // executeId
          executeId.value = res.data;
          startPolling(res.data);
        }
      });
    }
  });
}
const timer = ref();
const nodes = ref(
  props.tinyFlowData.nodes.map((node: any) => ({
    nodeId: node.id,
    nodeName: node.data.title,
  })),
);
// 轮询执行结果
function startPolling(executeId: any) {
  if (timer.value) return;
  timer.value = setInterval(() => executePolling(executeId), 1000);
}
function executePolling(executeId: any) {
  api
    .post('/userCenter/aiWorkflow/getChainStatus', {
      executeId,
      nodes: nodes.value,
    })
    .then((res) => {
      // 5 是挂起状态
      if (res.data.status !== 1 || res.data.status === 5) {
        stopPolling();
      }
      props.onAsyncExecute?.(res.data);
    });
}
function stopPolling() {
  submitLoading.value = false;
  if (timer.value) {
    clearInterval(timer.value);
    timer.value = null;
  }
}
onUnmounted(() => {
  stopPolling();
});
</script>

<template>
  <div>
    <ElForm label-position="top" ref="runForm" :model="runParams">
      <WorkflowFormItem
        v-model:run-params="runParams"
        :parameters="parameters"
      />
      <ElFormItem>
        <ElButton
          type="primary"
          @click="submitV2"
          :loading="submitLoading"
          :icon="Position"
        >
          {{ $t('button.run') }}
        </ElButton>
      </ElFormItem>
    </ElForm>
  </div>
</template>

<style scoped></style>
