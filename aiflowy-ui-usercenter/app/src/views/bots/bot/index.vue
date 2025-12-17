<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { cn, sortNodes } from '@aiflowy/utils';

import { ArrowLeft } from '@element-plus/icons-vue';
import {
  ElAside,
  ElButton,
  ElContainer,
  ElHeader,
  ElIcon,
  ElMain,
  ElSpace,
} from 'element-plus';

import { api } from '#/api/request';
import defaultBotAvatar from '#/assets/defaultBotAvatar.png';
import {
  Card,
  CardAvatar,
  CardContent,
  CardDescription,
  CardTitle,
} from '#/components/card';
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
const router = useRouter();
const workflowId = ref(route.params.id);
const workflowInfo = ref<any>({});
const runParams = ref<any>(null);
const initState = ref(false);
const tinyFlowData = ref<any>(null);
const workflowForm = ref();
async function getWorkflowInfo(workflowId: any) {
  api.get(`/userCenter/aiWorkflow/detail?id=${workflowId}`).then((res) => {
    workflowInfo.value = res.data;
    tinyFlowData.value = workflowInfo.value.content
      ? JSON.parse(workflowInfo.value.content)
      : {};
  });
}
async function getRunningParams() {
  api
    .get(`/userCenter/aiWorkflow/getRunningParameters?id=${workflowId.value}`)
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
  <ElContainer class="h-full">
    <ElHeader class="!px-8 !py-4" height="fit-content">
      <div class="flex flex-col gap-6">
        <ElSpace class="cursor-pointer" :size="10" @click="router.back()">
          <ElIcon color="#969799" size="24"><ArrowLeft /></ElIcon>
          <h1 class="text-2xl font-medium text-[#333333]">
            {{ workflowInfo?.title }}
          </h1>
        </ElSpace>
        <div
          class="flex items-center justify-between rounded-lg bg-[linear-gradient(153deg,#D3E3FD,#CBDEFE)] px-10 py-9"
        >
          <Card class="max-w-none cursor-auto items-center gap-7">
            <CardAvatar
              :size="72"
              :src="workflowInfo?.icon"
              :default-avatar="defaultBotAvatar"
            />
            <CardContent class="gap-3">
              <CardTitle class="text-3xl font-medium text-[#1A1A1A]">
                {{ workflowInfo?.title }}
              </CardTitle>
              <CardDescription class="text-base text-[#5E6673]">
                {{ workflowInfo?.description }}
              </CardDescription>
            </CardContent>
          </Card>
          <RouterLink to="">
            <ElButton color="#0066FF" size="large" round plain>
              执行记录
            </ElButton>
          </RouterLink>
        </div>
      </div>
    </ElHeader>
    <ElMain class="!px-8 !pb-4 !pt-0">
      <ElContainer class="h-full gap-4">
        <ElAside
          width="366px"
          class="flex flex-col gap-6 rounded-lg bg-white p-5"
        >
          <h1 class="text-base font-medium text-[#1A1A1A]">输入参数</h1>
          <WorkflowForm
            v-if="runParams && tinyFlowData"
            ref="workflowForm"
            :workflow-id="workflowId"
            :workflow-params="runParams"
            :on-submit="onSubmit"
            :on-async-execute="onAsyncExecute"
            :tiny-flow-data="tinyFlowData"
          />
        </ElAside>
        <ElAside width="366px">
          <div
            :class="cn('flex h-full flex-col gap-6 rounded-lg bg-white p-5')"
          >
            <h1 class="text-base font-medium text-[#1A1A1A]">执行步骤</h1>
            <WorkflowSteps
              v-if="tinyFlowData"
              :workflow-id="workflowId"
              :node-json="sortNodes(tinyFlowData)"
              :init-signal="initState"
              :polling-data="chainInfo"
              @resume="resumeChain"
            />
          </div>
        </ElAside>
        <div :class="cn('flex flex-1 flex-col gap-6 rounded-lg bg-white p-5')">
          <h1 class="text-base font-medium text-[#1A1A1A]">运行结果</h1>
          <div
            class="flex-1 rounded-lg border border-[#F0F0F0] bg-[#F7F7F7] p-4"
          >
            <ExecResult
              v-if="tinyFlowData"
              :workflow-id="workflowId"
              :node-json="sortNodes(tinyFlowData)"
              :init-signal="initState"
              :polling-data="chainInfo"
            />
          </div>
        </div>
      </ElContainer>
    </ElMain>
  </ElContainer>
</template>
