import { toFormData } from 'axios'

import request from '@/utils/request'

import {
  AccountInfoResponse,
  UpdateAccountRequest,
  UpdateAccountResponse,
  UploadAvatarRequest,
  UploadAvatarResponse,
} from './type'

enum API {
  AccountInfoURL = '/account/info',
  AccountURL = '/account',
  PasswordResetURL = '/account/password/reseet',
  PasswordURL = '/account/password',
  AvatarURL = '/avatar',
}

// 获取账号信息
export const reqAccountInfo = (): Promise<AccountInfoResponse> => {
  return request.get(API.AccountInfoURL)
}

// 更新账号信息
export const reqUpdateAccount = (data: UpdateAccountRequest): Promise<UpdateAccountResponse> => {
  return request.put(API.AccountURL, toFormData(data), {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
}

// 上传头像
export const reqUploadAvatar = (data: UploadAvatarRequest): Promise<UploadAvatarResponse> => {
  return request.post(API.AvatarURL, toFormData(data), {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
}
