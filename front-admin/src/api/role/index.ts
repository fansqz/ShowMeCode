import request from '@/utils/request';
import { toFormData } from '@/utils/format';

enum API {
  //获取角色列表
  LIST_ROLE_URL = '/manage/role/list',
  // 角色相关url
  ROLE_URL = '/manage/role',
  // 角色-apik控制
  ROLE_API = '/manage/role/api',
  // 角色-menu控制
  ROLE_MENU = '/manage/role/menu',
}

// 获取角色列表
export const reqRoleList = (data: any): Promise<any> => {
  return request.get(API.LIST_ROLE_URL, {
    params: data,
  });
};

// 获取角色
export const reqGetRole = (roleID: string): Promise<any> => {
  return request.get(API.ROLE_URL + `/${roleID}`);
};

// 创建角色
export const reqInsertRole = (data: any): Promise<any> => {
  return request.post(API.ROLE_URL, toFormData(data), {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

// 修改角色
export const reqUpdateRole = (data: any): Promise<any> => {
  return request.put(API.ROLE_URL, toFormData(data), {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

// 删除角色
export const reqDeleteRole = (id: string): Promise<any> => {
  return request.delete(API.ROLE_URL + '/' + id);
};

// 获取角色可访问api
export const reqRoleApi = (roleID: string): Promise<any> => {
  return request.get(API.ROLE_API + `/${roleID}`);
};

// 获取角色可访问的menu
export const reqRoleMenu = (roleID: string): Promise<any> => {
  return request.get(API.ROLE_MENU + `/${roleID}`);
};

// 更新角色可访问api
export const reqUpdateRoleApi = (data: any): Promise<any> => {
  return request.put(API.ROLE_API, data);
};

// 更新角色可访问的menu
export const reqUpdateRoleMenu = (data: any): Promise<any> => {
  return request.put(API.ROLE_MENU, data);
};
