export type BaseResponse = {
  code: number
  message: string
}

export type PageInfo = {
  size: number
  total: number
}

export type PageQuery = {
  page: number
  pageSize: number
  sortProperty?: string
  sortRule?: string
}
