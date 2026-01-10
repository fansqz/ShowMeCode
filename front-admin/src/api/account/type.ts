type AccountInfo = {
  avatar: string;
  loginName: string;
  username: string;
  email: string;
  phone: string;
  introduction: string;
  sex: string;
  birthDay: string;
  codingAge: number;
  roles: string[];
  menus: string[];
};

type AccountInfoResponse = BaseResponse & {
  data: AccountInfo;
};

type UpdateAccountRequest = {
  avatar: string;
  username: string;
  introduction: string;
  sex: string;
  birthDay: string;
};

type UpdateAccountResponse = BaseResponse & {};

type UploadAvatarRequest = {
  avatar: File;
};

type UploadAvatarResponse = BaseResponse & {
  data: string;
};
