import { descriptions } from '@/constants/description';

export type VisualDocumentDirectoryResponse = BaseResponse & {
  data: VisualDocumentDirectory[];
};

export type VisualDocumentDirectory = {
  // 可视化节点
  id: number;
  bankID: number;
  parentID: number;
  title: string;
  enable: boolean;
  order: number;
  chidren: VisualDocumentDirectory[];
};

export type VisualDocumentResponse = BaseResponse & {
  data: VisualDocument;
};

export type VisualDocument = {
  id: number;
  bankID: number;
  parentID: number;
  title: string;
  content: string;
  creatorID: number;
  enable: boolean;
  codeList: VisualDocumentCode[];
};

export type InsertVisualDocumentRequest = {
  bankID: number;
  parentID: number;
  title: string;
  content: string;
  enable: boolean;
};

export type InsertVisualDocumentResponse = BaseResponse & {
  data: number;
};

export type UpdateVisualDocumentRequest = {
  id: number;
  parentID: number;
  title: string;
  content: string;
  enable: boolean;
  codeList: VisualDocumentCode[];
};

export type VisualDocumentCode = {
  code: string;
  language: string;
  breakpoints: number[];
};

export type UpdateVisualDocumentResponse = BaseResponse & {
  data: string;
};

export type UpdateVisualDocumentParentResponse = BaseResponse & {
  data: string;
};

export type DeleteVisualDocumentResponse = BaseResponse & {
  data: string;
};