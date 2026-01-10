import { BaseResponse } from '../type'

export type LoginPassportRequest = {
  type: string
  account: string
  password: string
}

export type LoginEmaiCodeRequest = {
  type: string
  email: string
  code: string
}

export type LoginRequest = LoginEmaiCodeRequest | LoginPassportRequest

export type LoginResponse = BaseResponse & {
  data: string
}

export type SendCodeRequest = {
  email: string
  type: string
}

export type SendCodeResponse = BaseResponse & {}

export type RegisterRequest = {
  username: string
  email: string
  password: string
  code: string
}

export type RegisterReponse = BaseResponse & {}
