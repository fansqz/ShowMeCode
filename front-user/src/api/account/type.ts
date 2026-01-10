import { BaseResponse } from '../type'

export type AccountInfo = {
  avatar: string
  loginName: string
  username: string
  email: string
  phone: string
  introduction: string
  sex: string
  birthDay: string
  codingAge: number
}

export type AccountInfoResponse = BaseResponse & {
  data: AccountInfo
}

export type UpdateAccountRequest = {
  avatar: string
  username: string
  introduction: string
  sex: string
  birthDay: string
}

export type UpdateAccountResponse = BaseResponse & {}

export type UploadAvatarRequest = {
  avatar: File
}

export type UploadAvatarResponse = BaseResponse & {
  data: string
}
