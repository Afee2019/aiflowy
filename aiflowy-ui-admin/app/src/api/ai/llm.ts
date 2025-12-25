import { api } from '#/api/request.js';

// 获取LLM供应商
export async function getLlmProviderList() {
  return api.get('/api/v1/modelProvider/list');
}

// 保存LLM
export async function saveLlm(data: string) {
  return api.post('/api/v1/model/save', data);
}

// 删除LLM
export async function deleteLlm(data: any) {
  return api.post(`/api/v1/model/remove`, data);
}

// 修改LLM
export async function updateLlm(data: any) {
  return api.post(`/api/v1/model/update`, data);
}

// 一键添加LLM
export async function quickAddLlm(data: any) {
  return api.post(`/api/v1/model/quickAdd`, data);
}

export interface llmType {
  id: string;
  title: string;
  provider: string;
  llmModel: string;
  icon: string;
  description: string;
  modelType: string;
  groupName: string;
  added: boolean;
  aiLlmProvider: any;
}
