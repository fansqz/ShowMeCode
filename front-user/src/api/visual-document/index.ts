import request from '@/utils/request'

import { VisualDocumentDirectoryResponse, VisualDocumentResponse } from './type.ts'

enum API {
  //获取用户列表
  VISUAL_DOCUMENT_URL = '/learn/visual/document',
  VISUAL_DOCUMENT_DIRECTORY_URL = '/learn/visual/document/directory',
  VISUAL_DOCUMENT_BANK_URL = '/learn/visual/document/bank/all',
}

// 获取可视化文档目录
export const reqVisualDocumentDirectory = (
  bankID: number,
): Promise<VisualDocumentDirectoryResponse> => {
  return request.get(API.VISUAL_DOCUMENT_DIRECTORY_URL + '/' + bankID)
}

// 获取可视化文档
export const reqVisualDocument = (id: number): Promise<VisualDocumentResponse> => {
  return request.get(`${API.VISUAL_DOCUMENT_URL}/${id}`)
}
