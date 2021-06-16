// @flow

import { TradeStageStatusMap } from '../constants/tradeConst'

export type reqActor = 's' | 'b' | 'n'
export type Approval = 'approved' | 'rejected' | 'pending' | 'expired' | 'submitted'
export type CanMakeReq = 'no' | 'can' | 'pending' | 'approved'
export type TradeStageStatusType = $Values<typeof TradeStageStatusMap>
export type DoneStatus = 'nil' | 'doing' | 'done'
export type OfferPriceType = 'firm' | 'quote'
export type Incoterm = 'CAF' | 'CFR' | 'CFRC' | 'CIF' | 'CIFFO' | 'CNF' | 'CNFFO' | 'DEL' | 'FOB'
export type Currency = 'USD' | 'EUR' | 'RUB' | 'TRY' | 'UAH' | 'CNY'
export type NotifType = 'action' | 'alert'
export type LocationType = 'trade' | 'tradeOffer'
export type OrganizationType = {
  id: string,
  name: string,
  telephone: string,
  email: string,
  address: string,
}

export type OrgMapType = {
  org: OrganizationType,
  role: string
}

export type UserType = {
  id: string,
  firstName: string,
  lastName: string,
  emails: Array<string>,
  orgMap: Array<OrgMapType>,
  avatar: string,
  pubKey: string,
  roles: Array<string>,
  createdAt: Date,
  biography: string
}

export type StageModerator = {
  user: UserType,
  createdAt: Date
}

export type DocType = {
  id: string,
  name: string,
  hash: string,
  note: string,
  url: string,
  type: string,
  createdBy: UserType,
  createdAt: Date
}

export type TradeStageDocType = {
  doc: DocType,
  status: Approval,
  approvedBy: UserType,
  approvedAt: Date,
  rejectReason: string,
  expiresAt: Date,
  reqTx: string,
  index: number
}

export type ApproveReqType = {
  reqBy: UserType,
  reqAt: Date,
  reqActor: reqActor,
  reqReason: string,
  rejectReason: string,
  approvedBy: UserType,
  approvedAt: Date,
  status: Approval,
  reqTx: string,
  approvedTx: string
}

export type TradeStageAddReqType = {
  name: string,
  description: string,
  owner: reqActor,
  status: Approval,
  reqBy: UserType,
  reqAt: Date,
  reqActor: reqActor,
  reqReason: string,
  approvedBy: UserType,
  approvedAt: Date,
  rejectReason: string
}

export type TradeStageType = {
  name: string,
  description: string,
  addReqIdx: number,
  docs: Array<TradeStageDocType>,
  owner: reqActor,
  expiresAt: string,
  delReqs: Array<ApproveReqType>,
  closeReqs: Array<ApproveReqType>,
  moderator: StageModerator,
  stageDeleteStatus: CanMakeReq,
  stageCloseStatus: CanMakeReq
}

export type TradeStageTemplateType = {
  name: string,
  description: string,
  owner: reqActor
}

export type TradeTemplateType = {
  id: string,
  name: string,
  description: string,
  stages: Array<TradeStageTemplateType>
}

export type TradeType = {
  id: string,
  name: string,
  description: string,
  template: TradeTemplateType,
  buyer: UserType,
  seller: UserType,
  tradeCloseStatus: CanMakeReq,
  stages: Array<TradeStageType>,
  stageAddReqs: Array<TradeStageAddReqType>,
  createdBy: UserType,
  createdAt: Date,
  closeReqs: Array<ApproveReqType>,
  scAddr: string
}

export type TradeOfferType = {
  id: string,
  price: number,
  isSell: Boolean,
  priceType: OfferPriceType,
  currency: Currency,
  createdBy: UserType,
  createdAt: Date,
  expiresAt: Date,
  closedAt: Date,
  org: OrganizationType,
  isAnonymous: Boolean,
  commodity: string,
  comType: Array<string>,
  quality: string,
  origin: string,
  incoterm: Incoterm,
  marketLoc: string,
  vol: number,
  shipment: Array<Date>,
  note: string,
  terms: DocType
}

export type LoginInputType = {
  email: string,
  password: string
}

export type OrgMapInputType = {
  id: string,
  role: string
}

export type NewUserInputType = {
  firstName: string,
  lastName: string,
  email: string,
  password: string,
  avatar: string,
  publicKey: string,
  biography: string,
  orgMap: OrgMapInputType
}

export type NewTradeInputType = {
  templateID: string,
  name: string,
  description: string,
  sellerID: string,
  buyerID: string,
  tradeOfferID: string
}

export type NewStageInputType = {
  tid: string,
  name: string,
  description: string,
  reason: string,
  owner: string
}

export type TradeStagePathType = {
  tid: string,
  stageIdx: number
}

export type TradeStageDocPathType = {
  tid: string,
  stageIdx: number,
  stageDocIdx: number,
  stageDocHash: string
}

export type CreateNewDocType = {
  docID: string,
  tid: string,
  note: string,
  fileName: string,
  stageIdx: number,
  expiresAt: Date,
  hash: string,
  signedTx: string,
  withApproval: boolean
}

export type CreateTradeStageInputType = {
  tid: string,
  name: string,
  description: string,
  reason: string,
  owner: string
}

export type ChangePasswordType = {
  oldPassword: string,
  newPassword: string
}

export type UserProfileInputType = {
  firstName: string,
  lastName: string,
  orgMap: Array<OrgMapInputType>,
  biography: string
}

export type OrgInputType = {
  name: string,
  address: string,
  telephone: string,
  email: string
}

export type TradeOfferInput = {
  price: number,
  priceType: OfferPriceType,
  isSell: boolean,
  currency: Currency,
  expiresAt: Date,
  commodity: string,
  comType: Array<string>,
  quality: string,
  orgID: string,
  origin: string,
  incoterm: Incoterm,
  marketLoc: string,
  vol: number,
  shipment: Array<Date>,
  note: string,
  docID: string
}

export type Notification = {
  id: string,
  createdAt: Date,
  triggeredBy: UserType,
  receiver: Array<string>,
  type: NotifType,
  dismissed: Array<string>,
  entityID: string,
  msg: string,
  action: Approval
}
