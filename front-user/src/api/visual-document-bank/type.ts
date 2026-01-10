import { BaseResponse } from '../type'

export type VisualDocumentBank = {
  id: number
  name: string
  description: string
  creatorID?: number
  creatorName?: string
  enable: boolean
}

export type AllVisualDocumentBankyResponse = BaseResponse & {
  data: VisualDocumentBank[]
}

export type VisualDocumentResponse = BaseResponse & {
  data: VisualDocumentBank
}
