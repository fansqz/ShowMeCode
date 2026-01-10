import request from '@/utils/request';
import { toFormData } from '@/utils/format';

// 地址
enum API {
  GET_MENU_TREE = '/manage/menu/tree',
  MENU_URL = '/manage/menu',
}

// 获取menu
export const reqGetMenu = (id: number): Promise<any> => {
  return request.get(API.MENU_URL + `/+${id}`);
};

// 添加menu
export const reqInsertMenu = (menu: any): Promise<any> => {
  return request.post(API.MENU_URL, toFormData(menu), {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

// 修改menu
export const reqUpdateMenu = (menu: any): Promise<any> => {
  return request.put(API.MENU_URL, toFormData(menu), {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

// 删除menu
export const reqDeleteMenu = (id: number): Promise<any> => {
  return request.delete(API.MENU_URL + '/' + id);
};

// 获取menu树
export const reqGetMenuTree = (): Promise<any> => {
  return request.get(API.GET_MENU_TREE);
};
