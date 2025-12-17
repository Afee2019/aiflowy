<script setup lang="ts">
import type { FormInstance } from 'element-plus';

import { reactive, ref } from 'vue';

import { useAppConfig } from '@aiflowy/hooks';
import { useAccessStore, useUserStore } from '@aiflowy/stores';

import {
  ElAvatar,
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElSpace,
  ElUpload,
} from 'element-plus';
import { tryit } from 'radash';

import { api } from '#/api/request';
import { useAuthStore } from '#/store';

const { apiURL } = useAppConfig(import.meta.env, import.meta.env.PROD);
const userStore = useUserStore();
const useAuth = useAuthStore();
const formRef = ref<FormInstance>();
const formData = reactive({
  nickname: userStore?.userInfo?.nickname,
  avatar: userStore?.userInfo?.avatar,
  mobile: userStore?.userInfo?.mobile,
});

const editing = reactive({
  nickname: false,
  avatar: false,
  mobile: false,
});
const rules = {
  nickname: [{ required: true, message: '请输入用户名' }],
  mobile: [
    { required: true, message: '请输入手机号' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号' },
  ],
};

const handleFieldChange = async (field: keyof typeof formData) => {
  if (!formRef.value) return;
  const [err] = await tryit(formRef.value.validateField)([field]);

  if (!err) {
    submit(field);
  }
};
const handleCancelEdit = async (field: keyof typeof formData) => {
  editing[field] = false;
  formRef.value?.resetFields([field]);
};

const handleUploadChange = (response: any) => {
  formData.avatar = response?.data.path;
  submit('avatar');
};
const accessStore = useAccessStore();
const btnLoading = ref(false);
function submit(field: keyof typeof formData) {
  btnLoading.value = true;
  api.post('/userCenter/sysAccount/updateProfile', formData).then((res) => {
    btnLoading.value = false;
    if (res.errorCode === 0) {
      ElMessage.success('修改成功');
      editing[field] = false;
      useAuth.fetchUserInfo();
    }
  });
}
function logout() {
  useAuth.logout();
}
</script>

<template>
  <div
    class="h-full w-full bg-[linear-gradient(153deg,#FFFFFF,#EFF8FF)] pt-[156px]"
  >
    <ElForm
      ref="formRef"
      class="mx-auto"
      style="max-width: 619px; width: 100%"
      label-position="top"
      :model="formData"
      :rules="rules"
      hide-required-asterisk
    >
      <ElFormItem label="用户名：" prop="nickname">
        <div class="flex w-full justify-between gap-5">
          <template v-if="editing.nickname">
            <ElInput v-model="formData.nickname" />
            <ElSpace>
              <ElButton
                :loading="btnLoading"
                type="primary"
                @click="handleFieldChange('nickname')"
              >
                确定
              </ElButton>
              <ElButton @click="handleCancelEdit('nickname')"> 取消 </ElButton>
            </ElSpace>
          </template>
          <template v-else>
            <span>{{ formData.nickname }}</span>
            <ElButton @click="editing.nickname = true">编辑用户名</ElButton>
          </template>
        </div>
      </ElFormItem>
      <ElFormItem label="头像：">
        <ElSpace>
          <ElUpload
            :action="`${apiURL}/userCenter/commons/upload`"
            :headers="{
              'aiflowy-token': accessStore.accessToken,
            }"
            :show-file-list="false"
            :on-success="handleUploadChange"
          >
            <ElAvatar :src="formData.avatar" :size="46" />
          </ElUpload>
          <span>支持 2M 以内的 JPG 或 PNG 图片</span>
        </ElSpace>
      </ElFormItem>
      <ElFormItem label="手机号：" prop="mobile">
        <div class="flex w-full justify-between gap-5">
          <template v-if="editing.mobile">
            <ElInput v-model="formData.mobile" />
            <ElSpace>
              <ElButton
                :loading="btnLoading"
                type="primary"
                @click="handleFieldChange('mobile')"
              >
                确定
              </ElButton>
              <ElButton @click="handleCancelEdit('mobile')">取消</ElButton>
            </ElSpace>
          </template>
          <template v-else>
            <span>{{ formData.mobile }}</span>
            <ElButton @click="editing.mobile = true">修改手机号</ElButton>
          </template>
        </div>
      </ElFormItem>
      <ElFormItem>
        <div class="mt-20 flex w-full justify-center">
          <ElButton @click="logout" type="primary" class="!h-11 w-[333px]">退出登录</ElButton>
        </div>
      </ElFormItem>
    </ElForm>
  </div>
</template>
