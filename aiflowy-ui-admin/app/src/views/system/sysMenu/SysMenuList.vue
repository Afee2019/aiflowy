<script setup lang="ts">
import { onMounted, ref } from 'vue';

import { createIconifyIcon } from '@aiflowy/icons';

import { Delete, Edit, Plus } from '@element-plus/icons-vue';
import {
  ElButton,
  ElIcon,
  ElMessage,
  ElMessageBox,
  ElTable,
  ElTableColumn,
} from 'element-plus';

import { api } from '#/api/request';
import { $t } from '#/locales';

import SysMenuModal from './SysMenuModal.vue';

onMounted(() => {
  getTree();
});

const saveDialog = ref();
const treeData = ref([]);
const loading = ref(false);
function reset() {
  getTree();
}
function showDialog(row: any) {
  saveDialog.value.openDialog({ ...row });
}
function remove(row: any) {
  ElMessageBox.confirm('确定删除吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
    beforeClose: (action, instance, done) => {
      if (action === 'confirm') {
        instance.confirmButtonLoading = true;
        api
          .post('/api/v1/sysMenu/remove', { id: row.id })
          .then((res) => {
            instance.confirmButtonLoading = false;
            if (res.errorCode === 0) {
              ElMessage.success(res.message);
              reset();
              done();
            }
          })
          .catch(() => {
            instance.confirmButtonLoading = false;
          });
      } else {
        done();
      }
    },
  }).catch(() => {});
}
function getTree() {
  loading.value = true;
  api
    .get('/api/v1/sysMenu/list', {
      params: {
        asTree: true,
      },
    })
    .then((res) => {
      loading.value = false;
      treeData.value = res.data;
    });
}
</script>

<template>
  <div class="page-container">
    <SysMenuModal ref="saveDialog" @reload="reset" />
    <div class="handle-div">
      <ElButton @click="showDialog({})" type="primary">
        <ElIcon class="mr-1">
          <Plus />
        </ElIcon>
        {{ $t('button.add') }}
      </ElButton>
    </div>
    <ElTable :data="treeData" border row-key="id" v-loading="loading">
      <ElTableColumn width="50" />
      <ElTableColumn prop="menuType" label="菜单类型">
        <template #default="{ row }">
          {{ row.menuType }}
        </template>
      </ElTableColumn>
      <ElTableColumn prop="menuTitle" label="菜单标题">
        <template #default="{ row }">
          {{ $t(row.menuTitle) }}
        </template>
      </ElTableColumn>
      <ElTableColumn prop="menuUrl" label="菜单url">
        <template #default="{ row }">
          {{ row.menuUrl }}
        </template>
      </ElTableColumn>
      <ElTableColumn prop="component" label="组件路径">
        <template #default="{ row }">
          {{ row.component }}
        </template>
      </ElTableColumn>
      <ElTableColumn prop="menuIcon" label="图标">
        <template #default="{ row }">
          <component class="size-5" :is="createIconifyIcon(row.menuIcon)" />
        </template>
      </ElTableColumn>
      <ElTableColumn prop="isShow" label="是否显示">
        <template #default="{ row }">
          {{ row.isShow }}
        </template>
      </ElTableColumn>
      <ElTableColumn prop="permissionTag" label="权限标识">
        <template #default="{ row }">
          {{ row.permissionTag }}
        </template>
      </ElTableColumn>
      <ElTableColumn prop="sortNo" label="排序">
        <template #default="{ row }">
          {{ row.sortNo }}
        </template>
      </ElTableColumn>
      <ElTableColumn prop="created" label="创建时间">
        <template #default="{ row }">
          {{ row.created }}
        </template>
      </ElTableColumn>
      <ElTableColumn label="操作" width="150">
        <template #default="{ row }">
          <ElButton @click="showDialog(row)" link type="primary">
            <ElIcon class="mr-1">
              <Edit />
            </ElIcon>
            {{ $t('button.edit') }}
          </ElButton>
          <ElButton @click="remove(row)" link type="danger">
            <ElIcon class="mr-1">
              <Delete />
            </ElIcon>
            {{ $t('button.delete') }}
          </ElButton>
        </template>
      </ElTableColumn>
    </ElTable>
  </div>
</template>

<style scoped></style>
