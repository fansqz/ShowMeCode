type BaseResponse = {
  code: number;
  message: string;
};

type PageInfo = {
  size: number;
  total: number;
};

type PageQuery = {
  page: number;
  pageSize: number;
  sortProperty?: string;
  sortRule?: string;
};
