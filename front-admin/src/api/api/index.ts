import request from '@/utils/request';
import { toFormData } from '@/utils/format';

// 地址
enum API {
  GET_API_TREE = '/manage/api/tree',
  API_URL = '/manage/api',
}

// 获取api
export const reqGetApi = (id: number): Promise<any> => {
  return request.get<any, any>(API.API_URL + `/+${id}`);
};

// 添加api
export const reqInsertApi = (api: any): Promise<any> => {
  return request.post(API.API_URL, toFormData(api), {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

// 修改api
export const reqUpdateApi = (api: any): Promise<any> => {
  return request.put(API.API_URL, toFormData(api), {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

// 删除api
export const reqDeleteApi = (id: number): Promise<any> => {
  return request.delete<any, any>(API.API_URL + '/' + id);
};

// 获取api树
export const reqGetApiTree = (): Promise<any> => {
  return request.get<any, any>(API.GET_API_TREE);
};
