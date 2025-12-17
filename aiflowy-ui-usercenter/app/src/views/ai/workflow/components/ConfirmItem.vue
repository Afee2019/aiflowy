<script setup lang="ts">
import { Download } from '@element-plus/icons-vue';
import { ElIcon, ElText } from 'element-plus';

import confirmFile from '#/assets/ai/workflow/confirm-file.png';
// 导入你的图片资源
import confirmOther from '#/assets/ai/workflow/confirm-other.png';

// 定义 Props
const props = defineProps({
  // v-model 绑定值
  modelValue: {
    type: [String, Number, Object],
    default: null,
  },
  // 数据类型: text, image, video, audio, other, file
  selectionDataType: {
    type: String,
    default: 'text',
  },
  // 数据列表
  selectionData: {
    type: Array as () => any[],
    default: () => [],
  },
});

// 定义 Emits
const emit = defineEmits(['update:modelValue', 'change']);

// 判断是否选中
const isSelected = (item: any) => {
  return props.modelValue === item;
};

// 切换选中状态
const changeValue = (item: any) => {
  if (props.modelValue === item) {
    // 如果点击已选中的，则取消选中
    emit('update:modelValue', null);
    emit('change', null); // 触发 Element Plus 表单验证
  } else {
    emit('update:modelValue', item);
    emit('change', item); // 触发 Element Plus 表单验证
  }
};

// 获取图标
const getIcon = (type: string) => {
  return type === 'other' ? confirmOther : confirmFile;
};

// 下载处理
const handleDownload = (url: string) => {
  window.open(url, '_blank');
};
</script>

<template>
  <div class="custom-radio-group">
    <template v-for="(item, index) in selectionData" :key="index">
      <!-- 类型: Text -->
      <div
        v-if="selectionDataType === 'text'"
        class="custom-radio-option"
        :class="{ selected: isSelected(item) }"
        style="width: 100%; flex-shrink: 0"
        @click="changeValue(item)"
      >
        {{ item }}
      </div>

      <!-- 类型: Image -->
      <div
        v-else-if="selectionDataType === 'image'"
        class="custom-radio-option"
        :class="{ selected: isSelected(item) }"
        style="padding: 0"
        @click="changeValue(item)"
      >
        <img
          :src="item"
          alt=""
          style="width: 80px; height: 80px; border-radius: 8px; display: block"
        />
      </div>

      <!-- 类型: Video -->
      <div
        v-else-if="selectionDataType === 'video'"
        class="custom-radio-option"
        :class="{ selected: isSelected(item) }"
        @click="changeValue(item)"
      >
        <video controls :src="item" style="width: 162px; height: 141px"></video>
      </div>

      <!-- 类型: Audio -->
      <div
        v-else-if="selectionDataType === 'audio'"
        class="custom-radio-option"
        :class="{ selected: isSelected(item) }"
        style="width: 100%; flex-shrink: 0"
        @click="changeValue(item)"
      >
        <audio
          controls
          :src="item"
          style="width: 100%; height: 44px; margin-top: 8px"
        ></audio>
      </div>

      <!-- 类型: File 或 Other -->
      <div
        v-else-if="
          selectionDataType === 'other' || selectionDataType === 'file'
        "
        class="custom-radio-option"
        :class="{ selected: isSelected(item) }"
        style="width: 100%; flex-shrink: 0"
        @click="changeValue(item)"
      >
        <div
          style="
            display: flex;
            justify-content: space-between;
            align-items: center;
          "
        >
          <div style="width: 92%; display: flex; align-items: center">
            <img
              style="width: 20px; height: 20px; margin-right: 8px"
              alt=""
              :src="getIcon(selectionDataType)"
            />
            <!-- 使用 Element Plus 的 Text 组件处理省略号，如果没有安装 Element Plus，可以用普通的 span + css -->
            <ElText truncated>
              {{ item }}
            </ElText>
          </div>
          <div class="download-icon-btn" @click.stop="handleDownload(item)">
            <ElIcon><Download /></ElIcon>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.custom-radio-group {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.custom-radio-option {
  background-color: var(--el-bg-color);
  padding: 8px;
  border-radius: 8px;
  cursor: pointer;
  position: relative;
  box-shadow: 0 0 0 1px var(--el-border-color);
  transition: all 0.2s;
  box-sizing: border-box; /* 确保 padding 不会撑大宽度 */
}

.custom-radio-option:hover {
  box-shadow: 0 0 0 1px var(--el-color-primary-light-5);
}

.custom-radio-option.selected {
  box-shadow: 0 0 0 1px var(--el-color-primary-light-3);
  padding: 8px;
  background: var(--el-color-primary-light-9);
}

.custom-radio-option.selected::after {
  content: '';
  position: absolute;
  right: 0;
  bottom: 0;
  width: 16px;
  height: 16px;
  background-color: var(--el-color-primary);
  border-radius: 6px 2px;
  box-sizing: border-box;
}

.custom-radio-option.selected::before {
  content: '';
  position: absolute;
  right: 3px;
  bottom: 7px;
  width: 9px;
  height: 4px;
  border-left: 1px solid white;
  border-bottom: 1px solid white;
  transform: rotate(-45deg);
  z-index: 1;
}

.download-icon-btn {
  font-size: 18px;
  cursor: pointer;
  margin-right: 10px;
  display: flex; /* 为了对齐图标 */
  align-items: center;
}
</style>
