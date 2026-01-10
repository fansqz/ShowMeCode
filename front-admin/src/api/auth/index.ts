import request from '@/utils/request';
import { toFormData } from 'axios';

enum API {
  LoginURL = '/auth/login',
  UserInfoURL = '/auth/get/info',
  CodeSendURL = '/auth/code/send',
  RegisterURL = '/auth/register',
}

// 登录
export const reqLogin = (data: LoginRequest): Promise<LoginResponse> => {
  return request.post(API.LoginURL, toFormData(data), {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

// 发送验证码
export const reqSendCode = (data: SendCodeRequest): Promise<SendCodeResponse> => {
  return request.post(API.CodeSendURL, toFormData(data), {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

// 注册
export const reqRegister = (data: RegisterRequest): Promise<RegisterReponse> => {
  return request.post(API.RegisterURL, toFormData(data), {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};
