type LoginPassportRequest = {
  type: string;
  account: string;
  password: string;
};

type LoginEmaiCodeRequest = {
  type: string;
  email: string;
  code: string;
};

type LoginRequest = LoginEmaiCodeRequest | LoginPassportRequest;

type LoginResponse = BaseResponse & {
  data: string;
};

type SendCodeRequest = {
  email: string;
  type: string;
};

type SendCodeResponse = BaseResponse & {};

type RegisterRequest = {
  username: string;
  email: string;
  password: string;
  code: string;
};

type RegisterReponse = BaseResponse & {};
