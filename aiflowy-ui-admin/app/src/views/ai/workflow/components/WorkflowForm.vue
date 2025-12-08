<script setup lang="ts">
import type { FormInstance } from 'element-plus';

import { computed, onUnmounted, ref } from 'vue';

import { Position } from '@element-plus/icons-vue';
import {
  ElAlert,
  ElButton,
  ElCheckboxGroup,
  ElForm,
  ElFormItem,
  ElInput,
  ElRadioGroup,
  ElSelect,
} from 'element-plus';

import { api } from '#/api/request';
import { $t } from '#/locales';
import ChooseResource from '#/views/ai/resource/ChooseResource.vue';

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
function getContentType(item: any) {
  return item.contentType || 'text';
}
function isResource(contentType: any) {
  return ['audio', 'file', 'image', 'video'].includes(contentType);
}
function getCheckboxOptions(item: any) {
  if (item.enums) {
    return (
      item.enums?.map((option: any) => ({
        label: option,
        value: option,
      })) || []
    );
  }
  return [];
}
function choose(data: any, propName: string) {
  runParams.value[propName] = data.resourceUrl;
}
function resume(data: any) {
  data.executeId = executeId.value;
  submitLoading.value = true;
  api.post('/api/v1/aiWorkflow/resume', data).then((res) => {
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
      api.post('/api/v1/aiWorkflow/runAsync', data).then((res) => {
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
  timer.value = setInterval(() => executePolling(executeId), 500);
}
function executePolling(executeId: any) {
  api
    .post('/api/v1/aiWorkflow/getChainStatus', {
      executeId,
      nodes: props.tinyFlowData.nodes.map((node: any) => ({
        nodeId: node.id,
        nodeName: nodes.value,
      })),
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
      <ElFormItem
        v-for="(item, idx) in parameters"
        :prop="item.name"
        :key="idx"
        :label="item.formLabel || item.name"
        :rules="
          item.required
            ? [{ required: true, message: $t('message.required') }]
            : []
        "
      >
        <template v-if="getContentType(item) === 'text'">
          <ElInput
            v-if="item.formType === 'input' || !item.formType"
            v-model="runParams[item.name]"
            :placeholder="item.formPlaceholder"
          />
          <ElSelect
            v-if="item.formType === 'select'"
            v-model="runParams[item.name]"
            :placeholder="item.formPlaceholder"
            :options="getCheckboxOptions(item)"
            clearable
          />
          <ElInput
            v-if="item.formType === 'textarea'"
            v-model="runParams[item.name]"
            :placeholder="item.formPlaceholder"
            :rows="3"
            type="textarea"
          />
          <ElRadioGroup
            v-if="item.formType === 'radio'"
            v-model="runParams[item.name]"
            :options="getCheckboxOptions(item)"
          />
          <ElCheckboxGroup
            v-if="item.formType === 'checkbox'"
            v-model="runParams[item.name]"
            :options="getCheckboxOptions(item)"
          />
        </template>
        <template v-if="getContentType(item) === 'other'">
          <ElInput
            v-model="runParams[item.name]"
            :placeholder="item.formPlaceholder"
          />
        </template>
        <template v-if="isResource(getContentType(item))">
          <ElInput
            v-model="runParams[item.name]"
            :placeholder="item.formPlaceholder"
          />
          <ChooseResource :attr-name="item.name" @choose="choose" />
        </template>
        <ElAlert
          v-if="item.formDescription"
          type="info"
          style="margin-top: 5px"
        >
          {{ item.formDescription }}
        </ElAlert>
      </ElFormItem>
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
