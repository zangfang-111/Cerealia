// @flow

import { action, observable } from 'mobx'
import moment from 'moment'
import { getStageCloseStatus, getStageDeleteStatus } from '../lib/helper'
import { approveStatus, canMakeReqStatus, tradeActorMap } from '../constants/tradeConst'
import {
  addStageApprove,
  addStageReject,
  createStage,
  tradeCloseReq,
  tradeCloseReqApprove,
  tradeCloseReqReject,
  tradeStageCloseReq,
  tradeStageCloseReqApprove,
  tradeStageCloseReqReject,
  tradeStageDelReq,
  tradeStageDelReqApprove,
  tradeStageDelReqReject,
  tradeStageDocApprove,
  tradeStageDocReject,
  tradeStageSetExpireTime
} from '../graphql/trades'
import type {
  ApproveReqType,
  CreateNewDocType,
  CreateTradeStageInputType,
  TradeStageAddReqType,
  TradeStageDocPathType,
  TradeStageDocType,
  TradeStagePathType,
  TradeStageType,
  TradeType,
  UserType
} from '../model/flowType'
import currentUserStore from './current-user'
import appStore from './app-store'
import { mkEmptyDoc, mkEmptyUser } from '../services/generators'
import { GqlClient } from '../services/cerealia'

class TradeViewStore {
  @observable id: string
  @observable name: string
  @observable description: string
  @observable buyer: UserType
  @observable seller: UserType
  @observable tradeCloseStatus: string
  @observable stages: Array<TradeStageType>
  @observable stageAddReqs: Array<TradeStageAddReqType>
  @observable createdBy: UserType
  @observable createdAt: Date
  @observable closeReqs: Array<ApproveReqType>
  @observable loading: boolean
  @observable scAddr: string
  @observable currentRole: string = ''
  @observable currentStageIdx: number = -1

  gqlClient: Object
  constructor () {
    this.gqlClient = GqlClient
  }

  @action updateTrade (trade: TradeType, gqlClient: Object): void {
    // HACK: gqlClient for tests
    // This should be refactored, so test can customize the store in their own way
    this.gqlClient = gqlClient
    if (!trade.id) {
      return
    }
    this.id = trade.id
    this.name = trade.name
    this.description = trade.description
    this.buyer = trade.buyer
    this.seller = trade.seller
    this.stages = trade.stages
    this.stageAddReqs = trade.stageAddReqs
    this.createdBy = trade.createdBy
    this.createdAt = trade.createdAt
    this.scAddr = trade.scAddr
    this.closeReqs = trade.closeReqs
    this.setDocIndex()
    this.resetStageStatus()
    this.currentRole = this.buyer.id === currentUserStore.user.id
      ? tradeActorMap.b : tradeActorMap.s
  }

  @action createStageDocCallback = (params: CreateNewDocType, user: UserType): void => {
    let stageDoc: TradeStageDocType = mkEmptyDoc()
    stageDoc.doc.id = params.docID
    stageDoc.doc.name = params.fileName
    stageDoc.doc.note = params.note
    stageDoc.doc.createdBy = user
    stageDoc.doc.createdAt = moment.utc().local().format()
    stageDoc.doc.hash = params.hash
    stageDoc.status = params.withApproval ? 'pending' : 'submitted'
    stageDoc.expiresAt = params.expiresAt
    if (this.stages[params.stageIdx].docs) {
      stageDoc.index = this.stages[params.stageIdx].docs.length
      this.stages[params.stageIdx].docs.push(stageDoc)
    } else {
      stageDoc.index = 0
      this.stages[params.stageIdx].docs = [stageDoc]
    }
    this.resetStageStatus()
  }

  @action createStageCallback = (addStageReq: TradeStageAddReqType,
    withApproval: boolean): void => {
    this.stageAddReqs.push(addStageReq)
    const modUser = appStore.adminMode ? currentUserStore.user : mkEmptyUser()
    const moderator = {
      user: modUser,
      createdAt: new Date()
    }
    if (!withApproval) {
      let newStage = {
        name: addStageReq.name,
        description: addStageReq.description,
        addReqIdx: 0,
        owner: addStageReq.owner,
        docs: [],
        expiresAt: '',
        delReqs: [],
        closeReqs: [],
        stageDeleteStatus: canMakeReqStatus.can,
        stageCloseStatus: canMakeReqStatus.no,
        moderator: moderator
      }
      this.stages.push(newStage)
    }
  }

  @action addStageApproveCallback = (response: Object, stageIdx: number): void => {
    let newStage: TradeStageType = response.data.tradeStageAddReqApprove
    this.stageAddReqs[stageIdx].status = approveStatus.approved
    this.stages.push(newStage)
    this.resetStageStatus()
  }

  @action addStageRejectCallback = (stageIdx: number, reason: string): void => {
    this.stageAddReqs[stageIdx].status = approveStatus.rejected
    this.stageAddReqs[stageIdx].rejectReason = reason
  }

  @action stageDelReqCallback = (response: Object, stageIdx: number): void => {
    this.stages[stageIdx].delReqs.push(response.data.tradeStageDelReq)
    this.stages[stageIdx].stageDeleteStatus = canMakeReqStatus.pending
  }

  @action stageDelReqApproveCallback = (response: Object, stageIdx: number): void => {
    let lenDelReq = this.stages[stageIdx].delReqs.length
    this.stages[stageIdx].delReqs[lenDelReq - 1].status = approveStatus.approved
    this.stages[stageIdx].stageDeleteStatus = canMakeReqStatus.approved
    this.resetStageStatus()
  }

  @action stageDelReqRejectCallback = (response: Object, stageIdx: number): void => {
    let lenDelReq = this.stages[stageIdx].delReqs.length
    this.stages[stageIdx].delReqs[lenDelReq - 1].status = approveStatus.rejected
    this.stages[stageIdx].stageDeleteStatus = canMakeReqStatus.can
  }

  @action stageCloseReqCallback = (response: Object, stageIdx: number): void => {
    this.stages[stageIdx].closeReqs.push(response.data.tradeStageCloseReq)
    this.stages[stageIdx].stageCloseStatus = canMakeReqStatus.pending
  }

  @action stageCloseReqApproveCallback = (stageIdx: number): void => {
    let lenCloseReq = this.stages[stageIdx].closeReqs.length
    this.stages[stageIdx].closeReqs[lenCloseReq - 1].status = approveStatus.approved
    this.stages[stageIdx].stageCloseStatus = canMakeReqStatus.approved
    this.resetStageStatus()
  }

  @action stageCloseReqRejectCallback = (stageIdx: number, reason: string): void => {
    let lenCloseReq = this.stages[stageIdx].closeReqs.length
    this.stages[stageIdx].closeReqs[lenCloseReq - 1].status = approveStatus.rejected
    this.stages[stageIdx].closeReqs[lenCloseReq - 1].approvedBy = currentUserStore.user
    this.stages[stageIdx].closeReqs[lenCloseReq - 1].approvedAt = new Date()
    this.stages[stageIdx].closeReqs[lenCloseReq - 1].rejectReason = reason
    this.stages[stageIdx].stageCloseStatus = canMakeReqStatus.can
  }

  @action tradeCloseReqCallback = (response: Object): void => {
    this.closeReqs.push(response.data.tradeCloseReq)
    this.tradeCloseStatus = canMakeReqStatus.pending
  }

  @action tradeCloseReqApproveCallback = (): void => {
    let lenCloseReq = this.closeReqs.length
    this.closeReqs[lenCloseReq - 1].status = approveStatus.approved
    this.tradeCloseStatus = canMakeReqStatus.approved
  }

  @action tradeCloseReqRejectCallback = (): void => {
    let lenCloseReq = this.closeReqs.length
    this.closeReqs[lenCloseReq - 1].status = approveStatus.rejected
    this.tradeCloseStatus = canMakeReqStatus.can
  }

  @action stageDocApproveCallback = (id: TradeStageDocPathType): void => {
    this.stages[id.stageIdx].docs[id.stageDocIdx].status =
      approveStatus.approved
    this.resetStageStatus()
  }

  @action stageDocRejectCallback = (id: TradeStageDocPathType,
    reason: string, user: UserType): void => {
    let doc: TradeStageDocType = this.stages[id.stageIdx].docs[id.stageDocIdx]
    doc.status = approveStatus.rejected
    doc.approvedBy = user
    doc.rejectReason = reason
    doc.approvedAt = moment()
    this.resetStageStatus()
  }

  @action setStageExpireTimeUpdate = (id: TradeStagePathType, expiresAt: string): void => {
    this.stages[id.stageIdx].expiresAt = expiresAt
  }

  @action resetStageStatus = () => {
    let foundCurrentStage = false
    let tradeCloseStatus = canMakeReqStatus.can
    this.stages && this.stages.map((stage, index) => {
      stage.stageDeleteStatus = getStageDeleteStatus(stage)
      stage.stageCloseStatus = getStageCloseStatus(stage)
      // check stage status to find current stage
      if (stage.stageDeleteStatus !== canMakeReqStatus.approved &&
        stage.stageCloseStatus !== canMakeReqStatus.approved &&
        !foundCurrentStage) {
        this.currentStageIdx = index
        foundCurrentStage = true
      }
      // check stage status to find possibility of trade close
      if (stage.stageDeleteStatus !== canMakeReqStatus.approved &&
        stage.stageCloseStatus !== canMakeReqStatus.approved) {
        tradeCloseStatus = canMakeReqStatus.no
      }
    })
    // set tradeCloseStatus
    let n = this.closeReqs ? this.closeReqs.length : 0
    if (n === 0) {
      this.tradeCloseStatus = tradeCloseStatus
      return
    }
    if (this.closeReqs[n - 1].status === approveStatus.approved) {
      this.tradeCloseStatus = canMakeReqStatus.approved
    } else if (this.closeReqs[n - 1].status === approveStatus.pending) {
      this.tradeCloseStatus = canMakeReqStatus.pending
    }
  }

  @action setDocIndex = (): void => {
    this.stages && this.stages.map((stage) => {
      if (stage.docs && stage.docs.length > 0) {
        stage.docs.map((doc, i) => {
          doc.index = i
          if (moment().isAfter(moment(doc.expiresAt)) &&
            doc.status === approveStatus.pending) {
            doc.status = approveStatus.expired
          }
        })
      }
    })
  }

  @action async createStage (input: CreateTradeStageInputType,
    signedTx: string, withApproval: boolean) {
    let response = await this.gqlClient.mutate({
      mutation: createStage,
      variables: { input, signedTx, withApproval }
    })
    this.createStageCallback(response.data.tradeStageAddReq, withApproval)
  }

  @action async addStageApprove (id: TradeStagePathType, signedTx: string) {
    let response = await this.gqlClient.mutate({
      mutation: addStageApprove,
      variables: { 'id': id, 'signedTx': signedTx }
    })
    this.addStageApproveCallback(response, id.stageIdx)
  }

  @action async addStageReject (id: TradeStagePathType, signedTx: string, reason: string) {
    await this.gqlClient.mutate({
      mutation: addStageReject,
      variables: { 'id': id, 'signedTx': signedTx, 'reason': reason }
    })
    this.addStageRejectCallback(id.stageIdx, reason)
  }

  @action async stageDelReq (id: TradeStagePathType, reason: string) {
    let response = await this.gqlClient.mutate({
      mutation: tradeStageDelReq,
      variables: { 'id': id, 'reason': reason }
    })
    this.stageDelReqCallback(response, id.stageIdx)
  }

  @action async stageDelReqApprove (id: TradeStagePathType) {
    let response = await this.gqlClient.mutate({
      mutation: tradeStageDelReqApprove,
      variables: { 'id': id }
    })
    this.stageDelReqApproveCallback(response, id.stageIdx)
  }

  @action async stageDelReqReject (id: TradeStagePathType, reason: string) {
    let response = await this.gqlClient.mutate({
      mutation: tradeStageDelReqReject,
      variables: { 'id': id, 'reason': reason }
    })
    this.stageDelReqRejectCallback(response, id.stageIdx)
  }

  @action async stageCloseReq (id: TradeStagePathType, reason: string, signedTx: string) {
    let response = await this.gqlClient.mutate({
      mutation: tradeStageCloseReq,
      variables: { 'id': id, 'reason': reason, 'signedTx': signedTx }
    })
    this.stageCloseReqCallback(response, id.stageIdx)
  }

  @action async stageCloseReqApprove (id: TradeStagePathType, signedTx: string) {
    await this.gqlClient.mutate({
      mutation: tradeStageCloseReqApprove,
      variables: { 'id': id, 'signedTx': signedTx }
    })
    this.stageCloseReqApproveCallback(id.stageIdx)
  }

  @action async stageCloseReqReject (id: TradeStagePathType,
    reason: string, signedTx: string) {
    await this.gqlClient.mutate({
      mutation: tradeStageCloseReqReject,
      variables: { 'id': id, 'reason': reason, 'signedTx': signedTx }
    })
    this.stageCloseReqRejectCallback(id.stageIdx, reason)
  }

  @action async tradeCloseReq (reason: string, signedTx: string) {
    let response = await this.gqlClient.mutate({
      mutation: tradeCloseReq,
      variables: { 'id': this.id, 'reason': reason, 'signedTx': signedTx }
    })
    this.tradeCloseReqCallback(response)
  }

  @action async tradeCloseReqApprove (signedTx: string) {
    await this.gqlClient.mutate({
      mutation: tradeCloseReqApprove,
      variables: { 'id': this.id, 'signedTx': signedTx }
    })
    this.tradeCloseReqApproveCallback()
  }

  @action async tradeCloseReqReject (reason: string, signedTx: string) {
    await this.gqlClient.mutate({
      mutation: tradeCloseReqReject,
      variables: { 'id': this.id, 'reason': reason, 'signedTx': signedTx }
    })
    this.tradeCloseReqRejectCallback()
  }

  @action async stageDocApprove (id: TradeStageDocPathType, signedTx: string) {
    await this.gqlClient.mutate({
      mutation: tradeStageDocApprove,
      variables: { 'id': id, 'signedTx': signedTx }
    })
    this.stageDocApproveCallback(id)
  }

  @action async stageDocReject (id: TradeStageDocPathType, reason: string,
    user: UserType, signedTx: string) {
    await this.gqlClient.mutate({
      mutation: tradeStageDocReject,
      variables: { 'id': id, 'reason': reason, 'signedTx': signedTx }
    })
    this.stageDocRejectCallback(id, reason, user)
  }

  @action async setStageExpireTime (id: TradeStagePathType, expiresAt: string) {
    await this.gqlClient.mutate({
      mutation: tradeStageSetExpireTime,
      variables: { 'id': id, 'expiresAt': expiresAt }
    })
    this.setStageExpireTimeUpdate(id, expiresAt)
  }
}

export default new TradeViewStore()
