// @flow

import { approveStatus, buyerActor } from '../../constants/tradeConst'
import { sampleUser1, stageOpTx, stageDocAddTx } from './client-mock'
import tradeViewStore from '../../stores/tradeViewStore'
import type {
  CreateTradeStageInputType,
  NewUserInputType,
  TradeStageDocPathType,
  TradeStagePathType
} from '../../model/flowType'
import { mkEmptyUser } from '../../services/generators'

type TradeStageAddInput = {
  input: CreateTradeStageInputType,
  signedTx: string,
  withApproval: boolean
}

type TradeStageReqType = {
  id: TradeStagePathType,
  signedTx: string
}

type TradeStageWithReasonType = {
  id: TradeStagePathType,
  signedTx: string,
  reason: string
}

type TradeStageDocApprove = {
  id: TradeStageDocPathType,
  signedTx: string
}

type TradeStageDocReject = {
  id: TradeStageDocPathType,
  signedTx: string,
  reason: string
}

type TradeStageOpTxType = {
  id: TradeStagePathType,
  operationType: string
}

type TradeStageDocAddTx = {
  id: TradeStageDocPathType,
  operationType: string
}

export const resolvers = {
  Query: {
    stellarNet: (parent: any) => ({
      name: 'stellar-testNet',
      url: 'https://horizon-testnet.stellar.org',
      passphrase: 'Test SDF Network ; September 2015'
    })
  },
  Mutation: {
    tradeStageAddReq: (parent: any, { input, signedTx, withApproval }: TradeStageAddInput) => ({
      name: input.name,
      description: input.description,
      owner: input.owner,
      status: withApproval ? approveStatus.pending : approveStatus.approved,
      reqActor: buyerActor,
      reqBy: {
        id: '1',
        firstName: 'Sergey',
        lastName: 'Ivanov',
        roles: ['trader'],
        avatar: '',
        createdAt: '2018-09-02T14:43:35.423969Z',
        biography: ''
      },
      reqAt: '2018-09-02T14:43:35.423969Z',
      reqReason: input.reason
    }),
    tradeStageAddReqApprove: (parent: any, { id, signedTx }: TradeStageReqType) => ({
      name: tradeViewStore.stageAddReqs[id.stageIdx].name,
      description: tradeViewStore.stageAddReqs[id.stageIdx].description,
      owner: tradeViewStore.stageAddReqs[id.stageIdx].owner,
      addReqIdx: 0,
      docs: [],
      delReqs: [],
      closeReqs: []
    }),
    tradeStageAddReqReject: (
      parent: any, { id, signedTx, reason }: TradeStageWithReasonType
    ) => 0,
    tradeStageDocApprove: (parent: any, { id, signedTx }: TradeStageDocApprove) => ({
      docID: '1',
      status: approveStatus.approved,
      expiresAt: new Date(),
      approvedAt: new Date(),
      reqTx: '',
      approvedTx: signedTx,
      approvedBy: mkEmptyUser(),
      rejectReason: ''
    }),
    tradeStageDocReject: (parent: any, { id, signedTx, reason }: TradeStageDocReject) => 0,
    tradeStageCloseReq: (parent: any, { id, signedTx, reason }: TradeStageWithReasonType) => ({
      status: approveStatus.pending,
      reqActor: 'b',
      reqBy: mkEmptyUser(),
      reqAt: new Date(),
      reqTx: signedTx,
      reqReason: reason,
      approvedBy: mkEmptyUser(),
      approvedAt: new Date(),
      approvedTx: '',
      rejectReason: ''
    }),
    tradeStageCloseReqApprove: (parent: any, { id, signedTx }: TradeStageReqType) => ({
      status: approveStatus.approved,
      reqActor: 'b',
      reqBy: mkEmptyUser(),
      reqAt: new Date(),
      reqTx: signedTx,
      reqReason: '',
      approvedBy: mkEmptyUser(),
      approvedAt: new Date(),
      approvedTx: '',
      rejectReason: ''
    }),
    tradeStageCloseReqReject: (
      parent: any, { id, reason, signedTx }: TradeStageWithReasonType
    ) => 0,
    userLogin: (parent: any, input: NewUserInputType) => ({
      user: sampleUser1,
      token: 'aaaaaaaaaa'
    }),
    mkTradeStageAddTx: (
      parent: any, { id, operationType }: TradeStageOpTxType
    ) => stageOpTx,
    mkTradeStageDocTx: (
      parent: any, { id, operationType }: TradeStageDocAddTx
    ) => stageDocAddTx,
    mkTradeStageCloseTx: (
      parent: any, { id, operationType }: TradeStageOpTxType
    ) => stageOpTx
  }
}
