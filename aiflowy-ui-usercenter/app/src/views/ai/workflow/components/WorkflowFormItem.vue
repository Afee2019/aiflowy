<script setup lang="ts">
import {
  ElAlert,
  ElCheckboxGroup,
  ElFormItem,
  ElInput,
  ElRadioGroup,
  ElSelect,
} from 'element-plus';

import { $t } from '#/locales';
import ChooseResource from '#/views/ai/resource/ChooseResource.vue';

const props = defineProps({
  parameters: {
    type: Array<any>,
    required: true,
  },
  runParams: {
    type: Object,
    required: true,
  },
  propPrefix: {
    type: String,
    default: '',
  },
});
const emit = defineEmits(['update:runParams']);
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
function updateParam(name: string, value: any) {
  const newValue = { ...props.runParams, [name]: value };
  emit('update:runParams', newValue);
}
function choose(data: any, propName: string) {
  updateParam(propName, data.resourceUrl);
}
</script>

<template>
  <ElFormItem
    v-for="(item, idx) in parameters"
    :prop="`${propPrefix}${item.name}`"
    :key="idx"
    :label="item.formLabel || item.name"
    :rules="
      item.required ? [{ required: true, message: $t('message.required') }] : []
    "
  >
    <template v-if="getContentType(item) === 'text'">
      <ElInput
        v-if="item.formType === 'input' || !item.formType"
        :model-value="runParams[item.name]"
        @update:model-value="(val) => updateParam(item.name, val)"
        :placeholder="item.formPlaceholder"
      />
      <ElSelect
        v-if="item.formType === 'select'"
        :model-value="runParams[item.name]"
        @update:model-value="(val) => updateParam(item.name, val)"
        :placeholder="item.formPlaceholder"
        :options="getCheckboxOptions(item)"
        clearable
      />
      <ElInput
        v-if="item.formType === 'textarea'"
        :model-value="runParams[item.name]"
        @update:model-value="(val) => updateParam(item.name, val)"
        :placeholder="item.formPlaceholder"
        :rows="3"
        type="textarea"
      />
      <ElRadioGroup
        v-if="item.formType === 'radio'"
        :model-value="runParams[item.name]"
        @update:model-value="(val) => updateParam(item.name, val)"
        :options="getCheckboxOptions(item)"
      />
      <ElCheckboxGroup
        v-if="item.formType === 'checkbox'"
        :model-value="runParams[item.name]"
        @update:model-value="(val) => updateParam(item.name, val)"
        :options="getCheckboxOptions(item)"
      />
    </template>
    <template v-if="getContentType(item) === 'other'">
      <ElInput
        :model-value="runParams[item.name]"
        @update:model-value="(val) => updateParam(item.name, val)"
        :placeholder="item.formPlaceholder"
      />
    </template>
    <template v-if="isResource(getContentType(item))">
      <ElInput
        :model-value="runParams[item.name]"
        @update:model-value="(val) => updateParam(item.name, val)"
        :placeholder="item.formPlaceholder"
      />
      <ChooseResource :attr-name="item.name" @choose="choose" />
    </template>
    <ElAlert v-if="item.formDescription" type="info" style="margin-top: 5px">
      {{ item.formDescription }}
    </ElAlert>
  </ElFormItem>
</template>

<style scoped></style>
