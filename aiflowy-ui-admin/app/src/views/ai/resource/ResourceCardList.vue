<script setup lang="ts">
import {onMounted, ref, watch} from 'vue';

import {
  ElCard,
  ElCheckbox,
  ElCol,
  ElImage,
  ElRadio,
  ElRow,
  ElText,
  ElTooltip,
  ElTag,
} from 'element-plus';

import audioIcon from '#/assets/ai/resource/audio-icon.png';
import docIcon from '#/assets/ai/resource/doc-icon.png';
import otherIcon from '#/assets/ai/resource/other-icon.png';
import videoIcon from '#/assets/ai/resource/video-icon.png';
import { $t } from '#/locales';
import PreviewModal from '#/views/ai/resource/PreviewModal.vue';
import {useDictStore} from "#/store";

export interface ResourceCardProps {
  data: any[];
  multiple?: boolean;
  valueProp?: string;
}
const props = withDefaults(defineProps<ResourceCardProps>(), {
  multiple: false,
  valueProp: 'id',
});
onMounted(() => {
  initDict();
});
const dictStore = useDictStore();
function initDict() {
  dictStore.fetchDictionary('resourceType');
  dictStore.fetchDictionary('resourceOriginType');
}
const emit = defineEmits(['update:modelValue']);
const previewDialog = ref();
const radioValue = ref('');
const checkAll = ref(false);
function choose() {
  const arr = [];
  if (props.multiple) {
    for (const data of props.data) {
      if (data.checkboxValue) {
        arr.push(data);
      }
    }
  } else {
    if (radioValue.value) {
      for (const data of props.data) {
        if (data[props.valueProp] === radioValue.value) {
          arr.push(data);
        }
      }
    }
  }
  emit('update:modelValue', arr);
}
function handleCheckAllChange(val: any) {
  if (val) {
    for (const data of props.data) {
      data.checkboxValue = data[props.valueProp];
    }
  } else {
    for (const data of props.data) {
      data.checkboxValue = '';
    }
  }
}
function getSrc(item: any) {
  switch (item.resourceType) {
    case 0: {
      return item.resourceUrl;
    }
    case 1: {
      return audioIcon;
    }
    case 2: {
      return videoIcon;
    }
    case 3: {
      return docIcon;
    }
    default: {
      return otherIcon;
    }
  }
}
function preview(row: any) {
  previewDialog.value.openDialog({ ...row });
}
function getResourceTypeColor(item: any) {
  switch (item.resourceType) {
    case 0: {
      return '#0066FF';
    }
    case 1: {
      return '#FFA200';
    }
    case 2: {
      return '#5600FF';
    }
    case 3: {
      return '#0099CC';
    }
    default: {
      return '#757575';
    }
  }
}
function getResourceOriginColor(item: any) {
  switch (item.origin) {
    case 0: {
      return '#00B8A9';
    }
    case 1: {
      return '#0066FF';
    }
    default: {
      return '';
    }
  }
}
watch(
  [() => radioValue.value, () => props.data],
  () => {
    choose();
  },
  { deep: true },
);
</script>

<template>
  <div>
    <PreviewModal ref="previewDialog" />
    <ElCheckbox
      v-if="multiple"
      :label="$t('button.selectAll')"
      v-model="checkAll"
      @change="handleCheckAllChange"
    />
    <ElRow :gutter="20">
      <ElCol :span="6" v-for="item in data" :key="item.id" class="mb-5">
        <ElCard
          :body-style="{ padding: '12px', height: '285px' }"
          shadow="hover"
        >
          <div>
            <div>
              <ElCheckbox
                v-if="multiple"
                v-model="item.checkboxValue"
                :true-value="item[valueProp]"
                false-value=""
              />
              <ElRadio v-else v-model="radioValue" :value="item[valueProp]" />
            </div>
            <div>
              <ElImage
                @click="preview(item)"
                :src="getSrc(item)"
                style="width: 100%; height: 150px; cursor: pointer"
              />
            </div>
            <div>
              <ElTooltip
                :content="`${item.resourceName}.${item.suffix}`"
                placement="top"
              >
                <ElText truncated>
                  {{ item.resourceName }}.{{ item.suffix }}
                </ElText>
              </ElTooltip>
            </div>
            <div class="flex gap-1.5">
              <ElTag :color="getResourceOriginColor(item)" effect="dark">
                {{ dictStore.getDictLabel('resourceOriginType', item.origin) }}
              </ElTag>
              <ElTag :color="getResourceTypeColor(item)" effect="dark">
                {{ dictStore.getDictLabel('resourceType', item.resourceType) }}
              </ElTag>
            </div>
          </div>
        </ElCard>
      </ElCol>
    </ElRow>
  </div>
</template>

<style scoped></style>
