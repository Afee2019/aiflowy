<script setup lang="ts">
import type { CheckboxValueType } from 'element-plus';

import { onMounted, ref } from 'vue';

import { cn } from '@aiflowy/utils';

import { ElCheckbox, ElCheckboxGroup, ElSpace } from 'element-plus';

import { Card, CardContent, CardTitle } from '#/components/card';
import Tag from '#/components/tag/Tag.vue';
import { useDictStore } from '#/store';
import {
  getResourceOriginColor,
  getResourceTypeColor,
  getSrc,
} from '#/utils/resource';

interface GridProps {
  data: any[];
  onCheckedChange?: (ids: any[]) => void;
}

const props = defineProps<GridProps>();
const checkAll = ref(false);
const isIndeterminate = ref(false);
const checkedIds = ref<number[]>([]);

const handleCheckAllChange = (val: CheckboxValueType) => {
  checkedIds.value = val ? props.data.map((asset) => asset.id) : [];
  isIndeterminate.value = false;
  props.onCheckedChange?.(checkedItems());
};
const handleCheckedIdsChange = (value: CheckboxValueType[]) => {
  const checkedCount = value.length;
  checkAll.value = checkedCount === props.data.length;
  isIndeterminate.value = checkedCount > 0 && checkedCount < props.data.length;
  props.onCheckedChange?.(checkedItems());
};
onMounted(() => {
  initDict();
});
const dictStore = useDictStore();
function initDict() {
  dictStore.fetchDictionary('resourceType');
  dictStore.fetchDictionary('resourceOriginType');
}
function checkedItems() {
  return props.data.filter((asset) => checkedIds.value.includes(asset.id));
}
</script>

<template>
  <div class="flex w-full flex-col items-start gap-6">
    <ElCheckbox
      v-model="checkAll"
      :indeterminate="isIndeterminate"
      @change="handleCheckAllChange"
    >
      全选
    </ElCheckbox>
    <ElCheckboxGroup
      class="flex w-full flex-wrap items-center gap-5"
      v-model="checkedIds"
      @change="handleCheckedIdsChange"
    >
      <Card
        class="group relative max-w-[378px] flex-col gap-3 border border-[#F0F0F0] p-3 transition hover:-translate-y-2 hover:shadow-xl"
        v-for="asset in props.data"
        :key="asset.id"
      >
        <div
          class="flex h-[174px] w-full items-center justify-center rounded-lg bg-[#F9F8F8]"
        >
          <img class="shrink-0 object-cover w-[100px] h-[100px]" :src="getSrc(asset)" />
        </div>
        <CardContent class="w-full gap-3">
          <CardTitle class="font-medium text-[#1A1A1A]">
            {{ asset.resourceName }}
          </CardTitle>
          <div class="flex items-center justify-between">
            <ElSpace :size="10">
              <Tag
                size="small"
                :background-color="`${getResourceOriginColor(asset)}15`"
                :text-color="getResourceOriginColor(asset)"
                :text="
                  dictStore.getDictLabel('resourceOriginType', asset.origin)
                "
              />
              <Tag
                size="small"
                :background-color="`${getResourceTypeColor(asset)}15`"
                :text-color="getResourceTypeColor(asset)"
                :text="
                  dictStore.getDictLabel('resourceType', asset.resourceType)
                "
              />
            </ElSpace>
            <span class="text-xs text-[#969799]">{{ asset.size }}</span>
          </div>
        </CardContent>
        <div
          :class="
            cn(
              'absolute left-2.5 top-2.5 group-hover:block',
              !checkedIds.includes(asset.id) && 'hidden',
            )
          "
        >
          <ElCheckbox style="--el-checkbox-height: 1" :value="asset.id" />
        </div>
      </Card>
    </ElCheckboxGroup>
  </div>
</template>
