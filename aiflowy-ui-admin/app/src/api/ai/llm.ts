import { api } from '#/api/request.js';

// 获取LLM品牌列表
export async function getLlmBrandList() {
  return api.get('/api/v1/aiLlmBrand/list?asTree=true');
}

// 获取LLM供应商
export async function getLlmProviderList() {
  return api.get('/api/v1/aiLlmProvider/list');
}

// 保存LLM
export async function saveLlm(data: string) {
  return api.post('/api/v1/aiLlm/save', data);
}

// 删除LLM
export async function deleteLlm(data: any) {
  return api.post(`/api/v1/aiLlm/remove`, data);
}

// 修改LLM
export async function updateLlm(data: any) {
  return api.post(`/api/v1/aiLlm/update`, data);
}

// 一键添加LLM
export async function quickAddLlm(data: any) {
  return api.post(`/api/v1/aiLlm/quickAdd`, data);
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
}
