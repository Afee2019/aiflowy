<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';

import {
  ElCard,
  ElCheckbox,
  ElCol,
  ElImage,
  ElRadio,
  ElRow,
  ElText,
  ElTooltip,
} from 'element-plus';

import Tag from '#/components/tag/Tag.vue';
import { $t } from '#/locales';
import { useDictStore } from '#/store';
import {
  getResourceOriginColor,
  getResourceTypeColor,
  getSrc,
} from '#/utils/resource';
import PreviewModal from '#/views/ai/resource/PreviewModal.vue';

export interface ResourceCardProps {
  data: any[];
  multiple?: boolean;
  valueProp?: string;
}
const props = withDefaults(defineProps<ResourceCardProps>(), {
  multiple: false,
  valueProp: 'id',
});
const emit = defineEmits(['update:modelValue']);
onMounted(() => {
  initDict();
});
const dictStore = useDictStore();
function initDict() {
  dictStore.fetchDictionary('resourceType');
  dictStore.fetchDictionary('resourceOriginType');
}
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
function preview(row: any) {
  previewDialog.value.openDialog({ ...row });
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
              <Tag
                size="small"
                :background-color="`${getResourceOriginColor(item)}15`"
                :text-color="getResourceOriginColor(item)"
                :text="
                  dictStore.getDictLabel('resourceOriginType', item.origin)
                "
              />
              <Tag
                size="small"
                :background-color="`${getResourceTypeColor(item)}15`"
                :text-color="getResourceTypeColor(item)"
                :text="
                  dictStore.getDictLabel('resourceType', item.resourceType)
                "
              />
            </div>
          </div>
        </ElCard>
      </ElCol>
    </ElRow>
  </div>
</template>

<style scoped></style>
