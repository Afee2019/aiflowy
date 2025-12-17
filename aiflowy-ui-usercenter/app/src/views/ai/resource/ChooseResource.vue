<script setup lang="ts">
import { ref, watch } from 'vue';

import { ElButton, ElDialog } from 'element-plus';

import PageData from '#/components/page/PageData.vue';
import { $t } from '#/locales';
import ResourceCardList from '#/views/ai/resource/ResourceCardList.vue';

const props = defineProps({
  attrName: {
    type: String,
    required: true,
  },
});

const emit = defineEmits(['choose']);

const pageDataRef = ref();
const dialogVisible = ref(false);
const chooseResources = ref([]);
const currentChoose = ref<any>({});
function openDialog() {
  dialogVisible.value = true;
}
function closeDialog() {
  dialogVisible.value = false;
}
function confirm() {
  emit('choose', currentChoose.value, props.attrName);
  closeDialog();
}
watch(
  () => chooseResources.value,
  (newValue) => {
    currentChoose.value = newValue.length > 0 ? newValue[0] : {};
  },
);
</script>

<template>
  <div>
    <ElDialog
      v-model="dialogVisible"
      draggable
      :title="$t('aiResource.choose')"
      :before-close="closeDialog"
      :close-on-click-modal="false"
      width="80%"
      destroy-on-close
    >
      <PageData
        ref="pageDataRef"
        page-url="/userCenter/aiResource/page"
        :page-size="8"
        :page-sizes="[8, 12, 16, 20]"
      >
        <template #default="{ pageList }">
          <ResourceCardList v-model="chooseResources" :data="pageList" />
        </template>
      </PageData>
      <template #footer>
        <ElButton @click="closeDialog">
          {{ $t('button.cancel') }}
        </ElButton>
        <ElButton type="primary" @click="confirm">
          {{ $t('button.confirm') }}
        </ElButton>
      </template>
    </ElDialog>
    <ElButton @click="openDialog()">
      {{ $t('button.choose') }}
    </ElButton>
  </div>
</template>

<style scoped></style>
