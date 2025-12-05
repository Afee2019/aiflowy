import type { BotInfo, Message, RequestResult, Session } from '@aiflowy/types';

import { api } from '#/api/request.js';

/** 获取bot详情 */
export const getBotDetails = (id: string) => {
  return api.get<RequestResult<BotInfo>>('/api/v1/aiBot/getDetail', {
    params: { id },
  });
};

export interface GetSessionListParams {
  botId: string;
  tempUserId: string;
}
/** 获取bot对话列表 */
export const getSessionList = (params: GetSessionListParams) => {
  return api.get<RequestResult<{ cons: Session[] }>>(
    '/api/v1/conversation/externalList',
    { params },
  );
};

export interface SaveBotParams {
  icon: string;
  title: string;
  alias: string;
  description: string;
}
/** 创建Bot */
export const saveBot = (params: SaveBotParams) => {
  return api.post<RequestResult>('/api/v1/aiBot/save', { ...params });
};

export interface UpdateBotParams extends SaveBotParams {
  id: string;
}
/** 修改Bot */
export const updateBotApi = (params: UpdateBotParams) => {
  return api.post<RequestResult>('/api/v1/aiBot/update', { ...params });
};

/** 删除Bot */
export const removeBotFromId = (id: string) => {
  return api.post<RequestResult>('/api/v1/aiBot/update', { id });
};

export interface GetMessageListParams {
  sessionId: string;
  botId: string;
  isExternalMsg: number;
  tempUserId: string;
}
/** 获取单个对话的信息列表 */
export const getMessageList = (params: GetMessageListParams) => {
  return api.get<RequestResult<Message[]>>('/api/v1/aiBotMessage/messageList', {
    params,
  });
};
