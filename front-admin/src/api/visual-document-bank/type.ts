export type AllVisualDocumentBankyResponse = BaseResponse & {
  data: VisualDocumentBank[];
};

export type VisualDocumentResponse = BaseResponse & {
  data: VisualDocumentBank;
};

export type InsertVisualDocumentBankResponse = BaseResponse & {
  data: number;
};

export type UpdateVisualDocumentBankResponse = BaseResponse & {
  data: string;
};

export type DeleteVisualDocumentBankResponse = BaseResponse & {
  data: string;
};

export type VisualDocumentBank = {
  id: number;
  name: string;
  description: string;
  creatorID?: number;
  creatorName?: string;
  enable: boolean;
};
