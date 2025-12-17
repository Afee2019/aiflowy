<script setup lang="ts">
import { computed } from 'vue';

interface Props {
  /** 背景颜色 */
  backgroundColor?: string;
  /** 文字颜色 */
  textColor?: string;
  /** 标签文本内容 */
  text: number | string;
  /** 标签尺寸 */
  size?: 'large' | 'medium' | 'small';
  /** 是否可关闭 */
  closable?: boolean;
  /** 是否为圆形标签 */
  round?: boolean;
  /** 边框类型 */
  border?: 'dashed' | 'none' | 'solid';
  /** 边框颜色 */
  borderColor?: string;
  /** 主题类型 */
  type?: 'danger' | 'info' | 'primary' | 'success' | 'warning';
}

const props = withDefaults(defineProps<Props>(), {
  backgroundColor: 'rgb(236, 245, 255)',
  textColor: '#409eff',
  size: 'medium',
  closable: false,
  round: false,
  border: 'solid',
  borderColor: '',
  type: 'primary',
});

const emit = defineEmits<{
  close: [text: number | string];
}>();

// 根据类型自动设置颜色
const tagStyle = computed(() => {
  const style: Record<string, string> = {};

  switch (props.type) {
    case 'danger': {
      style.backgroundColor = 'rgb(254, 240, 240)';
      style.color = '#f56c6c';
      style.borderColor = props.borderColor || '#f56c6c';

      break;
    }
    case 'info': {
      style.backgroundColor = 'rgb(244, 244, 245)';
      style.color = '#909399';
      style.borderColor = props.borderColor || '#909399';

      break;
    }
    case 'primary': {
      style.backgroundColor = 'rgb(236, 245, 255)';
      style.color = '#409eff';
      style.borderColor = props.borderColor || '#409eff';

      break;
    }
    case 'success': {
      style.backgroundColor = 'rgb(240, 249, 235)';
      style.color = '#67c23a';
      style.borderColor = props.borderColor || '#67c23a';

      break;
    }
    case 'warning': {
      style.backgroundColor = 'rgb(253, 246, 236)';
      style.color = '#e6a23c';
      style.borderColor = props.borderColor || '#e6a23c';

      break;
    }
    // No default
  }

  // 自定义颜色优先级高于类型预设
  if (props.backgroundColor) {
    style.backgroundColor = props.backgroundColor;
  }
  if (props.textColor) {
    style.color = props.textColor;
  }
  if (props.borderColor) {
    style.borderColor = props.borderColor;
  }

  return style;
});

const handleClose = () => {
  emit('close', props.text);
};

// 尺寸映射
const sizeMap = {
  small: {
    fontSize: '12px',
    padding: '0 8px',
    height: '24px',
  },
  medium: {
    fontSize: '14px',
    padding: '0 12px',
    height: '32px',
  },
  large: {
    fontSize: '16px',
    padding: '0 16px',
    height: '40px',
  },
};
</script>

<template>
  <div
    class="tag"
    :class="[`tag--${size}`, { 'tag--round': round }, `tag--border-${border}`]"
    :style="[
      tagStyle,
      {
        fontSize: sizeMap[size].fontSize,
        padding: sizeMap[size].padding,
        height: sizeMap[size].height,
        lineHeight: sizeMap[size].height,
      },
    ]"
  >
    <span class="tag__content">
      <slot>{{ text }}</slot>
    </span>

    <span v-if="closable" class="tag__close" @click.stop="handleClose">
      <slot name="close-icon"> × </slot>
    </span>
  </div>
</template>

<style scoped>
.tag {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  box-sizing: border-box;
  white-space: nowrap;
  user-select: none;
  transition: all 0.2s ease-in-out;
  font-family:
    -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue',
    Arial, sans-serif;
}

.tag--small {
  font-size: 12px;
}

.tag--medium {
  font-size: 14px;
}

.tag--large {
  font-size: 16px;
}

.tag--round {
  border-radius: 9999px;
}

.tag--border-solid {
  border-style: solid;
}

.tag--border-dashed {
  border-style: dashed;
}

.tag--border-none {
  border: none;
}

.tag__content {
  display: inline-flex;
  align-items: center;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tag__close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  margin-left: 4px;
  cursor: pointer;
  width: 16px;
  height: 16px;
  font-size: 12px;
  line-height: 1;
  border-radius: 50%;
  transition: background-color 0.2s;
}

.tag__close:hover {
  background-color: rgba(0, 0, 0, 0.1);
}

/* 为可关闭标签调整内边距 */
.tag:has(.tag__close) {
  padding-right: 8px;
}

.tag:has(.tag__close).tag--small {
  padding-right: 6px;
}

.tag:has(.tag__close).tag--large {
  padding-right: 12px;
}
</style>
