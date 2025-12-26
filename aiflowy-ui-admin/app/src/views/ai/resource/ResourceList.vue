<script setup lang="ts">
import type { FormInstance } from 'element-plus';

import { computed, onMounted, ref } from 'vue';

import { formatBytes, tryit } from '@aiflowy/utils';

import {
  Delete,
  Download,
  Edit,
  MoreFilled,
  Plus,
} from '@element-plus/icons-vue';
import {
  ElAvatar,
  ElButton,
  ElDialog,
  ElDropdown,
  ElDropdownItem,
  ElDropdownMenu,
  ElForm,
  ElFormItem,
  ElIcon,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElTable,
  ElTableColumn,
  ElText,
  ElTooltip,
} from 'element-plus';

import { api } from '#/api/request';
import DictSelect from '#/components/dict/DictSelect.vue';
import PageData from '#/components/page/PageData.vue';
import PageSide from '#/components/page/PageSide.vue';
import Tag from '#/components/tag/Tag.vue';
import { $t } from '#/locales';
import { useDictStore } from '#/store';
import {
  getResourceOriginColor,
  getResourceTypeColor,
  getSrc,
} from '#/utils/resource';
import PreviewModal from '#/views/ai/resource/PreviewModal.vue';

import ResourceModal from './ResourceModal.vue';

onMounted(() => {
  initDict();
  getSideList();
});
const formRef = ref<FormInstance>();
const pageDataRef = ref();
const saveDialog = ref();
const previewDialog = ref();
const formInline = ref({
  resourceName: '',
  resourceType: '',
});
const dictStore = useDictStore();
function initDict() {
  dictStore.fetchDictionary('resourceType');
  dictStore.fetchDictionary('resourceOriginType');
}
function search(formEl: FormInstance | undefined) {
  formEl?.validate((valid) => {
    if (valid) {
      pageDataRef.value.setQuery(formInline.value);
    }
  });
}
function reset(formEl: FormInstance | undefined) {
  formEl?.resetFields();
  pageDataRef.value.setQuery({});
}
function showDialog(row: any) {
  saveDialog.value.openDialog({ ...row });
}
function remove(row: any) {
  ElMessageBox.confirm($t('message.deleteAlert'), $t('message.noticeTitle'), {
    confirmButtonText: $t('message.ok'),
    cancelButtonText: $t('message.cancel'),
    type: 'warning',
    beforeClose: (action, instance, done) => {
      if (action === 'confirm') {
        instance.confirmButtonLoading = true;
        api
          .post('/api/v1/resource/remove', { id: row.id })
          .then((res) => {
            instance.confirmButtonLoading = false;
            if (res.errorCode === 0) {
              ElMessage.success(res.message);
              reset(formRef.value);
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
function preview(row: any) {
  previewDialog.value.openDialog({ ...row });
}
function download(row: any) {
  window.open(row.resourceUrl, '_blank');
}

const fieldDefinitions = ref<any[]>([
  {
    prop: 'categoryName',
    label: $t('aiWorkflowCategory.categoryName'),
    type: 'input',
    required: true,
    placeholder: $t('aiWorkflowCategory.categoryName'),
  },
  {
    prop: 'sortNo',
    label: $t('aiWorkflowCategory.sortNo'),
    type: 'number',
    required: false,
    placeholder: $t('aiWorkflowCategory.sortNo'),
  },
]);
const sideList = ref<any>([]);
const controlBtns = [
  {
    icon: Edit,
    label: $t('button.edit'),
    onClick(row: any) {
      showControlDialog(row);
    },
  },
  {
    type: 'danger',
    icon: Delete,
    label: $t('button.delete'),
    onClick(row: any) {
      removeCategory(row);
    },
  },
];
const footerButton = {
  icon: Plus,
  label: $t('button.add'),
  onClick() {
    showControlDialog({});
  },
};
const sideDialogVisible = ref(false);
const sideFormData = ref<any>({});
const sideFormRef = ref<FormInstance>();
const sideFormRules = computed(() => {
  const rules: Record<string, any[]> = {};
  fieldDefinitions.value.forEach((field) => {
    const fieldRules = [];
    if (field.required) {
      fieldRules.push({
        required: true,
        message: `${$t('message.required')}`,
        trigger: 'blur',
      });
    }
    if (fieldRules.length > 0) {
      rules[field.prop] = fieldRules;
    }
  });
  return rules;
});
const sideSaveLoading = ref(false);

function changeCategory(category: any) {
  pageDataRef.value.setQuery({ categoryId: category.id });
}
function showControlDialog(item: any) {
  sideFormRef.value?.resetFields();
  sideFormData.value = { ...item };
  sideDialogVisible.value = true;
}
function removeCategory(row: any) {
  ElMessageBox.confirm($t('message.deleteAlert'), $t('message.noticeTitle'), {
    confirmButtonText: $t('message.ok'),
    cancelButtonText: $t('message.cancel'),
    type: 'warning',
    beforeClose: (action, instance, done) => {
      if (action === 'confirm') {
        instance.confirmButtonLoading = true;
        api
          .post('/api/v1/resourceCategory/remove', { id: row.id })
          .then((res) => {
            instance.confirmButtonLoading = false;
            if (res.errorCode === 0) {
              ElMessage.success(res.message);
              done();
              getSideList();
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
function handleSideSubmit() {
  formRef.value?.validate((valid) => {
    if (valid) {
      sideSaveLoading.value = true;
      const url = sideFormData.value.id
        ? '/api/v1/resourceCategory/update'
        : '/api/v1/resourceCategory/save';
      api.post(url, sideFormData.value).then((res) => {
        sideSaveLoading.value = false;
        if (res.errorCode === 0) {
          ElMessage.success(res.message);
          sideDialogVisible.value = false;
          getSideList();
        }
      });
    }
  });
}
const getSideList = async () => {
  const [, res] = await tryit<any>(
    api.get('/api/v1/resourceCategory/list', {
      params: { sortKey: 'sortNo', sortType: 'asc' },
    }),
  );

  if (res && res.errorCode === 0) {
    sideList.value = [
      {
        id: '',
        categoryName: $t('common.allCategories'),
      },
      ...res.data,
    ];
  }
};
</script>

<template>
  <div class="flex h-full flex-col gap-1.5 p-6">
    <PreviewModal ref="previewDialog" />
    <ResourceModal ref="saveDialog" @reload="reset" />
    <div class="flex items-center justify-between">
      <ElForm ref="formRef" inline :model="formInline">
        <ElFormItem prop="resourceType" class="!mr-3">
          <DictSelect
            v-model="formInline.resourceType"
            dict-code="resourceType"
            :placeholder="$t('aiResource.resourceType')"
          />
        </ElFormItem>
        <ElFormItem prop="resourceName" class="!mr-3">
          <ElInput
            v-model="formInline.resourceName"
            :placeholder="$t('aiResource.resourceName')"
          />
        </ElFormItem>
        <ElFormItem>
          <ElButton @click="search(formRef)" type="primary">
            {{ $t('button.query') }}
          </ElButton>
          <ElButton @click="reset(formRef)">
            {{ $t('button.reset') }}
          </ElButton>
        </ElFormItem>
      </ElForm>
      <div class="handle-div">
        <ElButton
          v-access:code="'/api/v1/resource/save'"
          @click="showDialog({})"
          type="primary"
        >
          <ElIcon class="mr-1">
            <Plus />
          </ElIcon>
          {{ $t('button.add') }}
        </ElButton>
      </div>
    </div>

    <div class="flex max-h-[calc(100vh-191px)] flex-1 gap-6">
      <PageSide
        label-key="categoryName"
        value-key="id"
        :menus="sideList"
        :control-btns="controlBtns"
        :footer-button="footerButton"
        @change="changeCategory"
      />
      <div class="bg-background h-full flex-1 overflow-auto rounded-lg p-5">
        <PageData
          ref="pageDataRef"
          page-url="/api/v1/resource/page"
          :page-size="10"
        >
          <template #default="{ pageList }">
            <ElTable :data="pageList" border>
              <ElTableColumn
                prop="resourceName"
                :label="$t('aiResource.resourceName')"
                width="300"
              >
                <template #default="{ row }">
                  <div class="flex items-center gap-2.5">
                    <ElAvatar :src="getSrc(row)" shape="square" :size="32" />
                    <div class="w-[200px]">
                      <ElTooltip
                        :content="`${row.resourceName}`"
                        placement="top"
                      >
                        <ElText truncated>
                          {{ row.resourceName }}
                        </ElText>
                      </ElTooltip>
                    </div>
                  </div>
                </template>
              </ElTableColumn>
              <ElTableColumn
                align="center"
                prop="suffix"
                :label="$t('aiResource.suffix')"
                width="60"
              >
                <template #default="{ row }">
                  {{ row.suffix }}
                </template>
              </ElTableColumn>
              <ElTableColumn
                align="center"
                prop="fileSize"
                :label="$t('aiResource.fileSize')"
              >
                <template #default="{ row }">
                  {{ formatBytes(row.fileSize) }}
                </template>
              </ElTableColumn>
              <ElTableColumn
                align="center"
                prop="origin"
                :label="$t('aiResource.origin')"
              >
                <template #default="{ row }">
                  <Tag
                    size="small"
                    :background-color="`${getResourceOriginColor(row)}15`"
                    :text-color="getResourceOriginColor(row)"
                    :text="
                      dictStore.getDictLabel('resourceOriginType', row.origin)
                    "
                  />
                </template>
              </ElTableColumn>
              <ElTableColumn
                align="center"
                prop="resourceType"
                :label="$t('aiResource.resourceType')"
              >
                <template #default="{ row }">
                  <Tag
                    size="small"
                    :background-color="`${getResourceTypeColor(row)}15`"
                    :text-color="getResourceTypeColor(row)"
                    :text="
                      dictStore.getDictLabel('resourceType', row.resourceType)
                    "
                  />
                </template>
              </ElTableColumn>
              <ElTableColumn prop="created" :label="$t('aiResource.created')">
                <template #default="{ row }">
                  {{ row.created }}
                </template>
              </ElTableColumn>
              <ElTableColumn
                :label="$t('common.handle')"
                width="140"
                align="right"
              >
                <template #default="{ row }">
                  <div class="flex items-center gap-3">
                    <div class="flex items-center">
                      <ElButton link type="primary" @click="preview(row)">
                        {{ $t('button.view') }}
                      </ElButton>
                      <ElButton link type="primary" @click="showDialog(row)">
                        {{ $t('button.edit') }}
                      </ElButton>
                    </div>
                    <ElDropdown>
                      <ElButton link :icon="MoreFilled" />

                      <template #dropdown>
                        <ElDropdownMenu>
                          <ElDropdownItem @click="download(row)">
                            <ElButton :icon="Download" link>
                              {{ $t('button.download') }}
                            </ElButton>
                          </ElDropdownItem>
                          <div v-access:code="'/api/v1/resource/remove'">
                            <ElDropdownItem @click="remove(row)">
                              <ElButton type="danger" :icon="Delete" link>
                                {{ $t('button.delete') }}
                              </ElButton>
                            </ElDropdownItem>
                          </div>
                        </ElDropdownMenu>
                      </template>
                    </ElDropdown>
                  </div>
                </template>
              </ElTableColumn>
            </ElTable>
          </template>
        </PageData>
      </div>
    </div>
    <ElDialog
      v-model="sideDialogVisible"
      :title="sideFormData.id ? `${$t('button.edit')}` : `${$t('button.add')}`"
      :close-on-click-modal="false"
    >
      <ElForm
        ref="sideFormRef"
        :model="sideFormData"
        :rules="sideFormRules"
        label-width="120px"
      >
        <!-- 动态生成表单项 -->
        <ElFormItem
          v-for="field in fieldDefinitions"
          :key="field.prop"
          :label="field.label"
          :prop="field.prop"
        >
          <ElInput
            v-if="!field.type || field.type === 'input'"
            v-model="sideFormData[field.prop]"
            :placeholder="field.placeholder"
          />
          <ElInputNumber
            v-else-if="field.type === 'number'"
            v-model="sideFormData[field.prop]"
            :placeholder="field.placeholder"
            style="width: 100%"
          />
        </ElFormItem>
      </ElForm>

      <template #footer>
        <ElButton @click="sideDialogVisible = false">
          {{ $t('button.cancel') }}
        </ElButton>
        <ElButton
          type="primary"
          @click="handleSideSubmit"
          :loading="sideSaveLoading"
        >
          {{ $t('button.confirm') }}
        </ElButton>
      </template>
    </ElDialog>
  </div>
</template>
