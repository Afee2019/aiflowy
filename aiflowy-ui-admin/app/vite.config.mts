import { defineConfig } from '@aiflowy/vite-config';

import ElementPlus from 'unplugin-element-plus/vite';

export default defineConfig(async () => {
  return {
    application: {},
    vite: {
      plugins: [
        ElementPlus({
          format: 'esm',
        }),
      ],
      server: {
        port: 8212,
        proxy: {
          '/api': {
            changeOrigin: true,
            // 后端API地址 (不需要 rewrite，直接转发)
            target: 'http://localhost:8211',
            ws: true,
          },
        },
      },
    },
  };
});
