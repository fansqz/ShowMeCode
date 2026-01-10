import request from '@/utils/request'

import { URLResponse } from './type'

enum API {
  GetURL = '/common/getURL',
}

export const reqGetURL = (path: string): Promise<URLResponse> => {
  return request.get(API.GetURL, {
    params: {
      path: path,
    },
  })
}
