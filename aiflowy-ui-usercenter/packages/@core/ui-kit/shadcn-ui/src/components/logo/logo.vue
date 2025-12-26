<script setup lang="ts">
import { computed } from 'vue';

import { cn } from '@aiflowy-core/shared/utils';

import { AIFlowyAvatar } from '../avatar';

interface Props {
  /**
   * @zh_CN 是否收起文本
   */
  collapsed?: boolean;
  /**
   * @zh_CN Logo 图片适应方式
   */
  fit?: 'contain' | 'cover' | 'fill' | 'none' | 'scale-down';
  /**
   * @zh_CN Logo 跳转地址
   */
  href?: string;
  /**
   * @zh_CN Logo 图片大小
   */
  logoSize?: number;
  /**
   * @zh_CN Logo 图标
   */
  src?: string;
  /**
   * @zh_CN 暗色主题 Logo 图标 (可选，若不设置则使用 src)
   */
  srcDark?: string;
  /**
   * @zh_CN Logo 文本
   */
  text: string;
  /**
   * @zh_CN Logo 主题
   */
  theme?: string;
}

defineOptions({
  name: 'AIFlowyLogo',
});

const props = withDefaults(defineProps<Props>(), {
  collapsed: false,
  href: 'javascript:void 0',
  logoSize: 120,
  src: '',
  srcDark: '',
  theme: 'light',
  fit: 'cover',
});

/**
 * @zh_CN 根据主题选择合适的 logo 图标
 */
const logoSrc = computed(() => {
  if (props.collapsed) {
    return '/logoMini.svg';
  }
  // 如果是暗色主题且提供了 srcDark，则使用暗色主题的 logo
  if (props.theme === 'dark' && props.srcDark) {
    return props.srcDark;
  }
  // 否则使用默认的 src
  return props.src;
});
</script>

<template>
  <div :class="theme" class="flex h-full items-center text-lg">
    <a
      :href="href"
      :class="
        cn(
          $attrs.class as string,
          'flex h-full items-center gap-2 overflow-hidden px-6 text-lg leading-normal transition-all duration-500',
          collapsed && 'px-5',
        )
      "
    >
      <AIFlowyAvatar
        v-if="logoSrc"
        :alt="text"
        :src="logoSrc"
        :size="collapsed ? 20 : logoSize"
        :fit="fit"
        class="relative rounded-none bg-transparent"
      />
      <!-- <template v-if="!collapsed">
        <slot name="text">
          <span class="text-foreground truncate text-nowrap font-semibold">
            {{ text }}
          </span>
        </slot>
      </template> -->
    </a>
  </div>
</template>
