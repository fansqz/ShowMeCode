import request from '@/utils/request';
import { toFormData } from '@/utils/format';

enum API {
  //获取用户列表
  LIST_USER_URL = '/manage/user/list',
  // 用户相关url
  USER_URL = '/manage/user',
  // 用户-角色控制
  USER_ROLE_URL = '/manage/user/role',
  // 获取简单的角色列表
  USER_ROLE_LIST_URL = '/manage/user/simpleRole/list',
}

// 根据用户id获取用户信息
export const reqGetUser = (userID: string): Promise<any> => {
  return request.get(API.USER_URL + `/${userID}`);
};

// 获取用户列表
export const reqUserList = (data: any): Promise<any> => {
  return request.get(API.LIST_USER_URL, {
    params: data,
  });
};

// 创建用户
export const reqInsertUser = (data: any): Promise<any> => {
  return request.post(API.USER_URL, toFormData(data), {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

// 修改用户
export const reqUpdateUser = (data: any): Promise<any> => {
  return request.put(API.USER_URL, toFormData(data), {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

// 删除角色
export const reqDeleteUser = (id: string): Promise<any> => {
  return request.delete(API.USER_URL + `/${id}`);
};

// 获取用户的角色id
export const reqUserRole = (userID: string): Promise<any> => {
  return request.get(API.USER_ROLE_URL + `/${userID}`);
};

// 更新角色可访问api
export const reqUpdateUserRole = (userID: string, roleIDs: string[]): Promise<any> => {
  return request.put(API.USER_ROLE_URL, { userID: userID, roleIDs: roleIDs });
};

// 获取简单的角色列表
export const reqSimpleRoleList = (): Promise<any> => {
  return request.get(API.USER_ROLE_LIST_URL);
};
