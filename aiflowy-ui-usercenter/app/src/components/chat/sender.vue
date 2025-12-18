<script setup lang="ts">
import { ref } from 'vue';
import { Sender } from 'vue-element-plus-x';

import { Promotion } from '@element-plus/icons-vue';
import { ElButton, ElIcon } from 'element-plus';

import { sseClient } from '#/api/request';

interface Props {
  sessionId: string | undefined;
  bot: any;
  addMessage?: (message: any) => void;
}
const props = defineProps<Props>();
const senderValue = ref('');
const btnLoading = ref(false);
// {
//   key: 0,
//     role: 'user',
//   placement: 'end',
//   content: '哈哈哈，让我试试',
//   typing: true,
// }
const msgKey = ref();
function sendMessage() {
  const data = {
    sessionId: props.sessionId,
    prompt: senderValue.value,
    botId: props.bot.id,
  };
  btnLoading.value = true;
  props.addMessage?.({
    key: Date.now(),
    role: 'user',
    placement: 'end',
    content: senderValue.value,
    typing: true,
  });
  msgKey.value = Date.now();
  sseClient.post('/userCenter/aiBot/chat', data, {
    onMessage(res) {
      const msg = {
        key: msgKey,
        role: 'assistant',
        placement: 'start',
        content: res.data,
        typing: true,
        loading: res.event !== 'finish',
      };
      props.addMessage?.(msg);
      if (res.event === 'finish') {
        btnLoading.value = false;
      }
    },
    onError(err) {
      console.error(err);
      btnLoading.value = false;
    },
    onFinished() {
      senderValue.value = '';
      btnLoading.value = false;
    },
  });
}
function getDisabled() {
  return !senderValue.value || !props.sessionId || btnLoading.value;
}
</script>

<template>
  <Sender
    v-model="senderValue"
    variant="updown"
    :auto-size="{ minRows: 2, maxRows: 5 }"
    clearable
    allow-speech
    placeholder="发送消息"
  >
    <!-- 自定义 prefix 前缀 -->
    <!-- <template #prefix>
    </template> -->

    <template #action-list>
      <div style="display: flex; align-items: center; gap: 8px">
        <ElButton
          :disabled="getDisabled()"
          @click="sendMessage"
          round
          color="#626aef"
        >
          <ElIcon><Promotion /></ElIcon>
        </ElButton>
      </div>
    </template>
  </Sender>
</template>
