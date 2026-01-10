import request from '@/utils/request'

import { AllVisualDocumentBankyResponse, VisualDocumentResponse } from './type.ts'

enum API {
  ALL_VISUAL_DOCUMENT_BANK_URL = '/learn/visual/document/bank/all',
  VISUAL_DOCUMENT_BANK_URL = '/learn/visual/document/bank',
}

// 获取知识库列表
export const reqAllVisualDocumentBank = (): Promise<AllVisualDocumentBankyResponse> => {
  return request.get<any, any>(API.ALL_VISUAL_DOCUMENT_BANK_URL)
}

// 根据id获取知识库信息
export const reqVisualDocumentBank = (id: number): Promise<VisualDocumentResponse> => {
  return request.get<any, any>(API.VISUAL_DOCUMENT_BANK_URL + '/' + id)
}
