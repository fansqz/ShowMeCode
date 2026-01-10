import request from '@/utils/request'

import {
  StructVisualRequest,
  ArrayVisualRequest,
  StructVisualResponse,
  ArrayVisualResponse,
  VisaulDocumentResponse,
  VisualDescriptionResponse,
  Array2DVisualRequest,
  Array2DVisualResponse,
} from './type'

enum API {
  VisaulURL = '/visual/debug',
  DocumentURL = '/visual/document',
}

export const reqGetVisualDescription = (debugID: string): Promise<VisualDescriptionResponse> => {
  return request.get(`${API.VisaulURL}/description/${debugID}`)
}

export const reqStructVisual = (data: StructVisualRequest): Promise<StructVisualResponse> => {
  return request.post(`${API.VisaulURL}/struct`, data)
}

export const reqArrayVisual = (data: ArrayVisualRequest): Promise<ArrayVisualResponse> => {
  return request.post(`${API.VisaulURL}/array`, data)
}

export const reqArray2DVisual = (data: Array2DVisualRequest): Promise<Array2DVisualResponse> => {
  return request.post(`${API.VisaulURL}/array2d`, data)
}

export const reqVisaulDocument = (): Promise<VisaulDocumentResponse> => {
  return request.get(`${API.DocumentURL}`)
}
