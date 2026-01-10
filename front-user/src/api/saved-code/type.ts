import { BaseResponse, PageQuery } from '../type'

// 用户保存的代码实体类型
export interface SavedCode {
  id: number
  documentID?: number
  language: string
  code: string
  remark?: string
  sortOrder: number
  createdAt: string
  updatedAt: string
}

// 创建用户保存代码的请求参数
export interface CreateSavedCodeRequest {
  documentID?: number
  language: string
  code: string
  remark?: string
}

// 更新用户保存代码的请求参数
export interface UpdateSavedCodeRequest {
  id: number
  documentID?: number
  language: string
  code: string
  remark?: string
}

// 获取代码列表的查询参数
export interface GetSavedCodeListQuery extends PageQuery {
  documentID?: number
  language?: string
}

// 创建用户保存代码的响应
export type CreateSavedCodeResponse = BaseResponse & {
  data: SavedCode
}

// 更新用户保存代码的响应
export type UpdateSavedCodeResponse = BaseResponse & {
  data: SavedCode
}

// 获取代码列表的响应
export type GetSavedCodeListResponse = BaseResponse & {
  data: {
    list: SavedCode[]
    total: number
    page: number
    size: number
  }
}

// 获取代码详情的响应
export type GetSavedCodeDetailResponse = BaseResponse & {
  data: SavedCode
}

// 删除代码的响应
export type DeleteSavedCodeResponse = BaseResponse & {
  data: null
}
