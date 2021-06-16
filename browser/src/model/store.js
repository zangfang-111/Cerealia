// @flow

import type {
  UserType,
  TradeType,
  TradeTemplateType,
  TradeStageType,
  TradeStageAddReqType,
  ApproveReqType, OrganizationType
} from './flowType'

export type TradeStoreType = {
  trades: Array<TradeType>,
  tradeTemplates: Array<TradeTemplateType>,
  selectedTab: number,
  triedToFetchTrades: boolean,
  getCurTrade: TradeType,
  fetchTrades: Function,
  initializeTrades: Function,
  createTrade: Function,
  fetchTradeTemplates: Function,
  setSelectedTab: Function
}

export type TradeViewType = {
  currentRole: string,
  id: string,
  name: string,
  description: string,
  buyer: UserType,
  seller: UserType,
  status: boolean,
  stages: Array<TradeStageType>,
  stageAddReqs: Array<TradeStageAddReqType>,
  createdBy: UserType,
  createdAt: Date,
  loading: boolean,
  scAddr: string,
  currentStageIdx: number,
  closeReqs: Array<ApproveReqType>,
  tradeCloseStatus: string,
  updateTrade: Function,
  createStage: Function,
  addStageApprove: Function,
  addStageReject: Function,
  stageDelReq: Function,
  stageDelReqApprove: Function,
  stageDelReqReject: Function,
  stageCloseReq: Function,
  stageCloseReqApprove: Function,
  stageCloseReqReject: Function,
  stageDocApprove: Function,
  stageDocReject: Function,
  stages: Array<TradeStageType>,
  stageAddReqs: Array<TradeStageAddReqType>,
  createStageDocUpdate: Function,
  tradeCloseReq: Function,
  tradeCloseReqApprove: Function,
  tradeCloseReqReject: Function,
  setStageExpireTimeUpdate: Function,
  setStageExpireTime: Function
}

export type UserStoreType = {
  user: UserType,
  users: Array<UserType>,
  organizations: Array<OrganizationType>,
  fetchUsers: Function,
  fetchOrganizations: Function,
  login: Function,
  getUserByID: Function,
  userName: Function,
  authenticate: Function,
  signup: Function,
  changePassword: Function,
  changeEmail: Function,
  updateUserProfile: Function,
  isAuthenticated: boolean,
  hasModeratorRole: Function,
  pubKey: Function
}

export type AppStoreType = {
  adminMode: boolean,
  appTheme: string,
  getCookie: Function,
  setAdminMode: Function,
  setTheme: Function
}

export type StellarNetwork = {
  name: string,
  url: string,
  passphrase: string
}

export type StellarStoreType = {
  keyVerified: boolean,
  network: StellarNetwork,
  updateClient: Function,
  initializeStellarNetwork: Function,
  validateStellarSecretKey: Function,
  validateStellarPublicKey: Function,
  validateAndSetUserKey: Function,
  signRawTX: Function,
  sign: Function,
  mkTradeStageDocTx: Function,
  mkTradeStageCloseTx: Function,
  mkTradeStageAddTx: Function,
  mkTradeCloseTx: Function,
  signDocApprovalTx: Function,
  signStageCloseTx: Function,
  signStageAddTx: Function,
  signTradeCloseTx: Function
}

export type SpinnerStoreType = {
  loading: boolean,
  spinnerTip: string,
  showSpinner: Function,
  hideSpinner: Function,
}

export type ModeratorTradeStoreType = {
  trades: Array<TradeType>,
  triedToFetchTrades: boolean,
  selectedTab: number,
  fetchTrades: Function,
  setSelectedTab: Function,
  getCurTrade: Function
}
