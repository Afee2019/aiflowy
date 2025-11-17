<script setup lang="ts">
import { ref } from 'vue';

import { ElButton } from 'element-plus';

import PageData from '#/components/page/PageData.vue';
import { $t } from '#/locales';

const pageDataRef = ref();
const search = () => {
  pageDataRef.value.setQuery({ loginName: 'test' });
};
const reset = () => {
  pageDataRef.value.setQuery({});
};
</script>

<template>
  <div class="page-container">
    <PageData
      ref="pageDataRef"
      page-url="/api/v1/sysAccount/page"
      :page-size="10"
      :init-query-params="{ status: 1 }"
    >
      <template #default="{ pageList }">
        <div v-for="item in pageList" :key="item.id">
          {{ item.loginName }}
        </div>
      </template>
    </PageData>

    <ElButton @click="search" type="primary">
      {{ $t('button.search.submit') }}
    </ElButton>
    <ElButton @click="reset">{{ $t('button.search.reset') }}</ElButton>
  </div>
</template>

<style scoped></style>
