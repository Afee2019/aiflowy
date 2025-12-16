<script setup lang="ts">
import { reactive, ref } from 'vue';

import {
  ElButton,
  ElDialog,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElOption,
  ElSelect,
  ElTag,
} from 'element-plus';

import { api } from '#/api/request';
import { $t } from '#/locales';
import { modelTypes } from '#/views/ai/llm/modelTypes';

import providerList from './providerList.json';

const props = defineProps({
  providerId: {
    type: String,
    default: '',
  },
});

const emit = defineEmits(['reload']);

const formDataRef = ref();
defineExpose({
  openAddDialog() {
    if (formDataRef.value) {
      formDataRef.value.resetFields();
    }
    Object.assign(formData, {
      modelType: '',
      title: '',
      llmModel: '',
      groupName: '',
      provider: '',
      endPoint: '',
      providerId: '',
      supportFeatures: [] as string[],
      options: {
        llmEndpoint: '',
        chatPath: '',
        embedPath: '',
        rerankPath: '',
      },
    });
    modelAbility.value.forEach((tag) => (tag.selected = false));
    dialogVisible.value = true;
  },
  openEditDialog(item: any) {
    dialogVisible.value = true;
    isAdd.value = false;

    Object.assign(formData, {
      id: item.id,
      modelType: item.modelType || '',
      title: item.title || '',
      llmModel: item.llmModel || '',
      groupName: item.groupName || '',
      provider: item.provider || '',
      endPoint: item.endPoint || '',
      supportFeatures: Array.isArray(item.supportFeatures)
        ? [...item.supportFeatures]
        : [],
      options: {
        llmEndpoint: item.options?.llmEndpoint || '',
        chatPath: item.options?.chatPath || '',
        embedPath: item.options?.embedPath || '',
        rerankPath: item.options?.rerankPath || '',
      },
    });

    modelAbility.value.forEach((tag) => {
      tag.selected = formData.supportFeatures.includes(tag.value);
    });
  },
});

const providerOptions =
  ref<Array<{ label: string; options: any; value: string }>>(providerList);
const isAdd = ref(true);
const dialogVisible = ref(false);
const formData = reactive({
  modelType: '',
  title: '',
  llmModel: '',
  groupName: '',
  provider: '',
  endPoint: '',
  providerId: '',
  supportFeatures: [] as string[],
  options: {
    llmEndpoint: '',
    chatPath: '',
    embedPath: '',
    rerankPath: '',
  },
});
const modelAbility = ref<
  Array<{
    activeType: 'danger' | 'info' | 'primary' | 'success' | 'warning'; // 选中后的专属类型
    defaultType: 'info'; // 默认灰色类型
    label: string;
    selected: boolean; // 选中状态
    value: string;
  }>
>([
  {
    label: $t('llm.modelAbility.reasoning'),
    value: 'reasoning',
    defaultType: 'info',
    activeType: 'success',
    selected: false,
  },
  {
    label: $t('llm.modelAbility.tool'),
    value: 'tool',
    defaultType: 'info',
    activeType: 'primary',
    selected: false,
  },
  {
    label: $t('llm.modelAbility.embedding'),
    value: 'embedding',
    defaultType: 'info',
    activeType: 'warning',
    selected: false,
  },
  {
    label: $t('llm.modelAbility.rerank'),
    value: 'rerank',
    defaultType: 'info',
    activeType: 'danger',
    selected: false,
  },
]);
const handleTagClick = (item: (typeof modelAbility.value)[0]) => {
  item.selected = !item.selected;

  if (!formData.supportFeatures) {
    formData.supportFeatures = [];
  }
  if (item.selected) {
    formData.supportFeatures.push(item.value);
  } else {
    formData.supportFeatures = formData.supportFeatures.filter(
      (v) => v !== item.value,
    );
  }
};
const closeDialog = () => {
  dialogVisible.value = false;
};
const rules = {
  title: [
    {
      required: true,
      message: $t('message.required'),
      trigger: 'blur',
    },
  ],
  llmModel: [
    {
      required: true,
      message: $t('message.required'),
      trigger: 'blur',
    },
  ],
  groupName: [
    {
      required: true,
      message: $t('message.required'),
      trigger: 'blur',
    },
  ],
  modelType: [
    {
      required: true,
      message: $t('message.required'),
      trigger: 'blur',
    },
  ],
  providerName: [
    {
      required: true,
      message: $t('message.required'),
      trigger: 'blur',
    },
  ],
  provider: [
    {
      required: true,
      message: $t('message.required'),
      trigger: 'blur',
    },
  ],
};
const btnLoading = ref(false);
const save = async () => {
  btnLoading.value = true;
  try {
    await formDataRef.value.validate();
    if (isAdd.value) {
      api.post('/api/v1/aiLlm/save', formData).then((res) => {
        if (res.errorCode === 0) {
          ElMessage.success(res.message);
          emit('reload');
          closeDialog();
        }
      });
    } else {
      api.post('/api/v1/aiLlm/update', formData).then((res) => {
        if (res.errorCode === 0) {
          ElMessage.success(res.message);
          emit('reload');
          closeDialog();
        }
      });
    }
  } finally {
    btnLoading.value = false;
  }
};
const handleChangeProvider = (val: string) => {
  const tempProvider = providerList.find((item) => item.value === val);
  if (!tempProvider) {
    return;
  }
  formData.provider = tempProvider.value;
  formData.providerId = props.providerId;
  formData.options.llmEndpoint = providerOptions.value.find(
    (item) => item.value === val,
  )?.options.llmEndpoint;
  formData.options.embedPath = providerOptions.value.find(
    (item) => item.value === val,
  )?.options.embedPath;
  formData.options.chatPath = providerOptions.value.find(
    (item) => item.value === val,
  )?.options.chatPath;
  formData.options.rerankPath = providerOptions.value.find(
    (item) => item.value === val,
  )?.options.rerankPath;
};
</script>

<template>
  <ElDialog
    v-model="dialogVisible"
    draggable
    :title="isAdd ? $t('button.add') : $t('button.edit')"
    :before-close="closeDialog"
    :close-on-click-modal="false"
    align-center
    width="482"
  >
    <ElForm
      label-width="100px"
      ref="formDataRef"
      :model="formData"
      status-icon
      :rules="rules"
    >
      <ElFormItem prop="title" :label="$t('llm.title')">
        <ElInput v-model.trim="formData.title" />
      </ElFormItem>
      <ElFormItem prop="modelType" :label="$t('llm.modelType')">
        <ElSelect v-model="formData.modelType" @change="handleChangeProvider">
          <ElOption
            v-for="item in modelTypes"
            :key="item.value"
            :label="item.label"
            :value="item.value || ''"
          />
        </ElSelect>
      </ElFormItem>
      <ElFormItem prop="provider" :label="$t('llm.provider')">
        <ElSelect v-model="formData.provider" @change="handleChangeProvider">
          <ElOption
            v-for="item in providerOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value || ''"
          />
        </ElSelect>
      </ElFormItem>

      <ElFormItem prop="llmModel" :label="$t('llm.llmModel')">
        <ElInput v-model.trim="formData.llmModel" />
      </ElFormItem>
      <ElFormItem prop="endPoint" :label="$t('llmProvider.endPoint')">
        <ElInput v-model.trim="formData.endPoint" />
      </ElFormItem>
      <ElFormItem prop="chatPath" :label="$t('llmProvider.chatPath')">
        <ElInput v-model.trim="formData.options.chatPath" />
      </ElFormItem>
      <ElFormItem prop="embedPath" :label="$t('llmProvider.embedPath')">
        <ElInput v-model.trim="formData.options.embedPath" />
      </ElFormItem>
      <ElFormItem prop="groupName" :label="$t('llm.groupName')">
        <ElInput v-model.trim="formData.groupName" />
      </ElFormItem>
      <ElFormItem prop="ability" :label="$t('llm.ability')">
        <div class="model-ability">
          <ElTag
            class="model-ability-tag"
            v-for="item in modelAbility"
            :key="item.value"
            :type="item.selected ? item.activeType : item.defaultType"
            @click="handleTagClick(item)"
            :class="{ 'tag-selected': item.selected }"
          >
            {{ item.label }}
          </ElTag>
        </div>
      </ElFormItem>
    </ElForm>
    <template #footer>
      <ElButton @click="closeDialog">
        {{ $t('button.cancel') }}
      </ElButton>
      <ElButton
        type="primary"
        @click="save"
        :loading="btnLoading"
        :disabled="btnLoading"
      >
        {{ $t('button.save') }}
      </ElButton>
    </template>
  </ElDialog>
</template>

<style scoped>
.headers-container-reduce {
  align-items: center;
}
.addHeadersBtn {
  width: 100%;
  border-style: dashed;
  border-color: var(--el-color-primary);
  border-radius: 8px;
  margin-top: 8px;
}
.head-con-content {
  margin-bottom: 8px;
  align-items: center;
}
.model-ability {
  display: flex;
  flex-wrap: nowrap;
  align-items: center;
  gap: 8px;
}
.model-ability-tag {
  cursor: pointer;
}
</style>
