import type { ServerSentEventMessage } from 'fetch-event-stream';

/**
 * 该文件可自行根据业务逻辑进行调整
 */
import type { RequestClientOptions } from '@aiflowy/request';

import { useAppConfig } from '@aiflowy/hooks';
import { preferences } from '@aiflowy/preferences';
import {
  authenticateResponseInterceptor,
  defaultResponseInterceptor,
  errorMessageResponseInterceptor,
  RequestClient,
} from '@aiflowy/request';
import { useAccessStore } from '@aiflowy/stores';

import { ElMessage } from 'element-plus';
import { events } from 'fetch-event-stream';

import { useAuthStore } from '#/store';

import { refreshTokenApi } from './core';

const { apiURL } = useAppConfig(import.meta.env, import.meta.env.PROD);

function createRequestClient(baseURL: string, options?: RequestClientOptions) {
  const client = new RequestClient({
    ...options,
    baseURL,
  });

  /**
   * 重新认证逻辑
   */
  async function doReAuthenticate() {
    console.warn('Access token or refresh token is invalid or expired. ');
    const accessStore = useAccessStore();
    const authStore = useAuthStore();
    accessStore.setAccessToken(null);
    if (
      preferences.app.loginExpiredMode === 'modal' &&
      accessStore.isAccessChecked
    ) {
      accessStore.setLoginExpired(true);
    } else {
      await authStore.logout();
    }
  }

  /**
   * 刷新token逻辑
   */
  async function doRefreshToken() {
    const accessStore = useAccessStore();
    const resp = await refreshTokenApi();
    const newToken = resp.data;
    accessStore.setAccessToken(newToken);
    return newToken;
  }

  function formatToken(token: null | string) {
    return token ? `${token}` : null;
  }

  // 请求头处理
  client.addRequestInterceptor({
    fulfilled: async (config) => {
      const accessStore = useAccessStore();

      config.headers['aiflowy-token'] = formatToken(accessStore.accessToken);
      config.headers['Accept-Language'] = preferences.app.locale;
      return config;
    },
  });

  // 处理返回的响应数据格式
  client.addResponseInterceptor(
    defaultResponseInterceptor({
      codeField: 'errorCode',
      dataField: 'data',
      showErrorMessage: (message) => {
        ElMessage.error(message);
      },
      successCode: 0,
    }),
  );

  // token过期的处理
  client.addResponseInterceptor(
    authenticateResponseInterceptor({
      client,
      doReAuthenticate,
      doRefreshToken,
      enableRefreshToken: preferences.app.enableRefreshToken,
      formatToken,
    }),
  );

  // 通用的错误处理,如果没有进入上面的错误处理逻辑，就会进入这里
  client.addResponseInterceptor(
    errorMessageResponseInterceptor((msg: string, error) => {
      // 这里可以根据业务进行定制,你可以拿到 error 内的信息进行定制化处理，根据不同的 code 做不同的提示，而不是直接使用 message.error 提示 msg
      // 当前mock接口返回的错误字段是 error 或者 message
      const responseData = error?.response?.data ?? {};
      const errorMessage = responseData?.error ?? responseData?.message ?? '';
      // 如果没有错误信息，则会根据状态码进行提示
      ElMessage.error(errorMessage || msg);
    }),
  );

  return client;
}

export const requestClient = createRequestClient(apiURL, {
  responseReturn: 'data',
});

export const api = createRequestClient(apiURL, {
  responseReturn: 'body',
});

export const baseRequestClient = new RequestClient({ baseURL: apiURL });

export interface SseOptions {
  onMessage?: (message: ServerSentEventMessage) => void;
  onError?: (err: any) => void;
  onFinished?: () => void;
}
export class SseClient {
  private controller: AbortController | null = null;
  private currentRequestId = 0;

  abort(): void {
    if (this.controller) {
      this.controller.abort();
      this.controller = null;
    }
  }

  isActive(): boolean {
    return this.controller !== null;
  }

  async post(url: string, data?: any, options?: SseOptions): Promise<void> {
    // 生成唯一的请求ID
    const requestId = ++this.currentRequestId;
    const currentRequestId = requestId;

    // 如果已有请求，先取消
    this.abort();

    // 创建新的控制器
    const controller = new AbortController();
    this.controller = controller;

    // 保存信号的引用到局部变量
    const signal = controller.signal;

    try {
      const res = await fetch(apiURL + url, {
        method: 'POST',
        signal, // 使用局部变量 signal
        headers: this.getHeaders(),
        body: JSON.stringify(data),
      });

      if (!res.ok) {
        const error = new Error(`HTTP ${res.status}: ${res.statusText}`);
        options?.onError?.(error);
        return;
      }

      // 在开始事件流之前检查是否还是同一个请求
      if (this.currentRequestId !== currentRequestId) {
        return;
      }

      const msgEvents = events(res, signal);

      try {
        for await (const event of msgEvents) {
          // 每次迭代都检查是否还是同一个请求
          if (this.currentRequestId !== currentRequestId) {
            break;
          }
          options?.onMessage?.(event);
        }
      } catch (innerError) {
        options?.onError?.(innerError);
      }

      // 只有在还是同一个请求的情况下才调用 onFinished
      if (this.currentRequestId === currentRequestId) {
        options?.onFinished?.();
      }
    } catch (error) {
      if (this.currentRequestId !== currentRequestId) {
        return;
      }
      console.error('SSE错误:', error);
      options?.onError?.(error);
    } finally {
      // 只有当还是当前请求时才清除 controller
      if (this.currentRequestId === currentRequestId) {
        this.controller = null;
      }
    }
  }

  private getHeaders() {
    const accessStore = useAccessStore();
    return {
      Accept: 'text/event-stream',
      'Content-Type': 'application/json',
      'aiflowy-token': accessStore.accessToken || '',
    };
  }
}

export const sseClient = new SseClient();
