import request from '@/utils/request';
import {
  VisualDocumentDirectoryResponse,
  VisualDocumentResponse,
  InsertVisualDocumentRequest,
  InsertVisualDocumentResponse,
  UpdateVisualDocumentRequest,
  UpdateVisualDocumentResponse,
  UpdateVisualDocumentParentResponse,
  DeleteVisualDocumentResponse,
} from './type.ts';

enum API {
  //获取用户列表
  VISUAL_DOCUMENT_URL = '/manage/visual/document',
  VISUAL_DOCUMENT_DIRECTORY_URL = '/manage/visual/document/directory',
  VISUAL_DOCUMENT_CODE_URL = '/manage/visual/document/code',
  VISUAL_DOCUMENT_FILE_URL = '/manage/visual/document/file',
}

// 获取可视化文档目录
export const reqVisualDocumentDirectory = (
  bankID: number,
): Promise<VisualDocumentDirectoryResponse> => {
  return request.get(API.VISUAL_DOCUMENT_DIRECTORY_URL + '/' + bankID);
};

// 获取可视化文档
export const reqVisualDocument = (id: number): Promise<VisualDocumentResponse> => {
  return request.get(`${API.VISUAL_DOCUMENT_URL}/${id}`);
};

// 添加可视化文档
export const reqInsertVisualDocument = (
  req: InsertVisualDocumentRequest,
): Promise<InsertVisualDocumentResponse> => {
  return request.post(API.VISUAL_DOCUMENT_URL, req);
};

// 更新可视化文档
export const reqUpdateVisualDocument = (
  req: UpdateVisualDocumentRequest,
): Promise<UpdateVisualDocumentResponse> => {
  return request.put(API.VISUAL_DOCUMENT_URL, req);
};

// 更新可视化文档的父节点
export const reqUpdateVisualDocumentDirectory = (
  draggingNodeID: number,
  dragNodeID: number,
  type: string,
): Promise<UpdateVisualDocumentParentResponse> => {
  return request.post(API.VISUAL_DOCUMENT_URL + '/directory', {
    draggingDocumentID: draggingNodeID,
    dragDocumentID: dragNodeID,
    eventType: type,
  });
};

// 删除可视化文档
export const reqDeleteVisualDocument = (id: number): Promise<DeleteVisualDocumentResponse> => {
  return request.delete(`${API.VISUAL_DOCUMENT_URL}/${id}`);
};

// 上传可视化文档文件
export const reqUploadVisualDocumentFile = (data: { file: File }): Promise<any> => {
  const formData = new FormData();
  formData.append('file', data.file);
  return request.post(API.VISUAL_DOCUMENT_FILE_URL, formData);
};
