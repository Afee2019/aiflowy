<script setup lang="ts">
import type { FormInstance } from 'element-plus';

import type { Ref } from 'vue';

import { nextTick, ref } from 'vue';

import { $t } from '@aiflowy/locales';

import {
  ElButton,
  ElDialog,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
} from 'element-plus';

import { api } from '#/api/request';

interface BasicFormItem {
  weChatMpAppId: string;
  weChatMpSecret: string;
  weChatMpToken: string;
  EncodingAESKey: string;
}

const dialogVisible = ref(false);
const botId = ref('');
const openDialog = (newBotId: string, options: BasicFormItem) => {
  nextTick(() => {
    basicFormRef.value?.resetFields();
  });
  botId.value = newBotId;
  basicForm.value = options;
  dialogVisible.value = true;
};

const basicForm: Ref<BasicFormItem> = ref({
  weChatMpAppId: '',
  weChatMpSecret: '',
  weChatMpToken: '',
  EncodingAESKey: '',
});
const basicFormRef = ref<FormInstance>();
defineExpose({
  openDialog(botId: string, options: BasicFormItem) {
    openDialog(botId, options);
  },
});

const rules = ref({
  weChatMpAppId: [
    {
      required: true,
      message: $t('message.required'),
      trigger: 'blur',
    },
  ],
  weChatMpSecret: [
    {
      required: true,
      message: $t('message.required'),
      trigger: 'blur',
    },
  ],
  weChatMpToken: [
    {
      required: true,
      message: $t('message.required'),
    },
  ],
});
const handleConfirm = () => {
  basicFormRef.value?.validate((valid: boolean) => {
    if (valid) {
      api
        .post('/api/v1/aiBot/updateOptions', {
          id: botId.value,
          options: {
            weChatMpAppId: basicForm.value.weChatMpAppId,
            weChatMpSecret: basicForm.value.weChatMpSecret,
            weChatMpToken: basicForm.value.weChatMpToken,
            EncodingAESKey: basicForm.value.EncodingAESKey,
          },
        })
        .then((res) => {
          if (res.errorCode === 0) {
            ElMessage.success($t('message.updateOkMessage'));
          } else {
            ElMessage.error(res.message);
          }
        });
      dialogVisible.value = false;
    }
  });
};
</script>

<template>
  <ElDialog
    v-model="dialogVisible"
    title="微信公众号配置"
    width="700"
    align-center
  >
    <ElForm
      ref="basicFormRef"
      style="width: 100%; margin-top: 20px"
      :model="basicForm"
      label-width="auto"
      :rules="rules"
    >
      <ElFormItem label="AppId" prop="weChatMpAppId" label-position="right">
        <ElInput v-model="basicForm.weChatMpAppId" />
      </ElFormItem>
      <ElFormItem label="Secret" prop="weChatMpSecret" label-position="right">
        <ElInput v-model="basicForm.weChatMpSecret" />
      </ElFormItem>
      <ElFormItem label="Token" prop="weChatMpToken" label-position="right">
        <ElInput v-model="basicForm.weChatMpToken" />
      </ElFormItem>
      <ElFormItem
        label="EncodingAESKey"
        prop="EncodingAESKey"
        label-position="right"
      >
        <ElInput v-model="basicForm.EncodingAESKey" />
      </ElFormItem>
    </ElForm>

    <template #footer>
      <div class="dialog-footer">
        <ElButton @click="dialogVisible = false">
          {{ $t('button.cancel') }}
        </ElButton>
        <ElButton type="primary" @click="handleConfirm">
          {{ $t('button.confirm') }}
        </ElButton>
      </div>
    </template>
  </ElDialog>
</template>

<style scoped></style>
