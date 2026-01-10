import request from '@/utils/request';
import {
  AllVisualDocumentBankyResponse,
  InsertVisualDocumentBankResponse,
  UpdateVisualDocumentBankResponse,
  VisualDocumentResponse,
  DeleteVisualDocumentBankResponse,
  VisualDocumentBank,
} from './type.ts';

enum API {
  //知识库相关url
  BANK_URL = '/manage/visual/document/bank',
  //获取知识库列表
  ALL_BANK_URL = '/manage/visual/document/bank/all',
}

// 获取知识库列表
export const reqAllVisualDocumentBank = (): Promise<AllVisualDocumentBankyResponse> => {
  return request.get<any, any>(API.ALL_BANK_URL);
};

// 创建知识库
export const reqInsertVisualDocumentBank = (
  data: VisualDocumentBank,
): Promise<InsertVisualDocumentBankResponse> => {
  return request.post(API.BANK_URL, data);
};

// 修改知识库
export const reqUpdateVisualDocumentBank = (
  data: VisualDocumentBank,
): Promise<UpdateVisualDocumentBankResponse> => {
  return request.put(API.BANK_URL, data);
};

// 根据id获取知识库信息
export const reqVisualDocumentBank = (id: number): Promise<VisualDocumentResponse> => {
  return request.get<any, any>(API.BANK_URL + '/' + id);
};

// 删除知识库
export const reqDeleteVisualDocumentBank = (
  id: string,
): Promise<DeleteVisualDocumentBankResponse> => {
  return request.delete(API.BANK_URL + '/' + id);
};
