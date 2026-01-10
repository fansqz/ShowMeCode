import request from '@/utils/request';

enum API {
  GET_URL = '/common/getURL',
}

export const reqGetURL = (path: string): Promise<any> => {
  return request.get(API.GET_URL, {
    params: {
      path: path,
    },
  });
};
