// @flow

import { notification } from 'antd'
import type {
  TradeStageType,
  TradeStagePathType,
  TradeType
} from '../model/flowType'
import { approveStatus, canMakeReqStatus } from '../constants/tradeConst'
import * as moment from 'moment-timezone'

export const addNotificationHelper =
  (content: Object, notificationType: string): void => {
    notification[notificationType]({
      message: notificationType,
      description: content.toString()
    })
    if (notificationType === 'error') {
      console.error(content)
    }
  }

export const mkJWTHeader = (): Object => {
  const token = localStorage.getItem('auth_token')
  let header = {
    Authorization: `Bearer ${token || ''}`
  }
  return header
}

export const mkFormDataHeaders = () => {
  let header = {
    Accept: 'multipart/form-data',
    'Content-Type': 'multipart/form-data'
  }
  return Object.assign(header, mkJWTHeader())
}

export const createTradeStagePath = (
  tradeID: string,
  stageIndex: number): TradeStagePathType => {
  return {
    tid: tradeID,
    stageIdx: stageIndex
  }
}

export const toLocalTime = (utcTime: Date): string => {
  return moment.utc(utcTime).local().format('YYYY-MM-DD HH:mm:ss')
}

export const getStageDeleteStatus = (stage: TradeStageType) => {
  if (stage.docs && stage.docs.length > 0) {
    return canMakeReqStatus.no
  }
  if (stage.delReqs && stage.delReqs.length > 0) {
    let lastReq = stage.delReqs[stage.delReqs.length - 1]
    if (lastReq.status === approveStatus.approved) {
      return canMakeReqStatus.approved
    } else if (lastReq.status === approveStatus.pending) {
      return canMakeReqStatus.pending
    }
  }
  return canMakeReqStatus.can
}

export const getStageCloseStatus = (stage: TradeStageType) => {
  if (stage.closeReqs && stage.closeReqs.length > 0) {
    let lastReq = stage.closeReqs[stage.closeReqs.length - 1]
    if (lastReq.status === approveStatus.approved) {
      return canMakeReqStatus.approved
    } else if (lastReq.status === approveStatus.pending) {
      return canMakeReqStatus.pending
    }
  }
  for (let doc of stage.docs) {
    if (doc.status !== approveStatus.rejected) {
      return canMakeReqStatus.can
    }
  }
  return canMakeReqStatus.no
}

export const hasPendingDoc = (stage: TradeStageType) => {
  for (let doc of stage.docs) {
    if (doc.status === approveStatus.pending) {
      return true
    }
  }
  return false
}

export const getTradeStatus = (trade: TradeType) => {
  let tradeStatus = canMakeReqStatus.can
  trade.stages && trade.stages.map((stage) => {
    let stageDeleteStatus = getStageDeleteStatus(stage)
    let stageCloseStatus = getStageCloseStatus(stage)
    if (stageDeleteStatus !== canMakeReqStatus.approved &&
      stageCloseStatus !== canMakeReqStatus.approved) {
      tradeStatus = canMakeReqStatus.no
    }
  })
  // set tradeCloseStatus
  let n = trade.closeReqs ? trade.closeReqs.length : 0
  if (n === 0) {
    return tradeStatus
  }
  if (trade.closeReqs[n - 1].status === approveStatus.approved) {
    return canMakeReqStatus.approved
  } else if (trade.closeReqs[n - 1].status === approveStatus.pending) {
    return canMakeReqStatus.pending
  } else {
    return tradeStatus
  }
}

export const fileHash = (file: Object, callback: Function) => {
  let BLAKE2s = require('blake2s')
  let reader = new FileReader()
  let h = new BLAKE2s(32)
  let fileByteArray = []
  reader.onload = function () {
    let arrayBuffer = this.result
    let array = new Uint8Array(arrayBuffer)
    for (let i = 0; i < array.length; i++) {
      fileByteArray.push(array[i])
    }
    h.update(fileByteArray)
    callback(h.digest('hex'))
  }
  reader.readAsArrayBuffer(file)
}

let browsers = [ // list of [agent name, min version]
  ['chrome', 64],
  ['chromium', 64],
  ['firefox.*', 58],
  ['edge.*', 16],
  ['AppleWebKit[a-z0-9.(), /]+Version', 10.1]
]

export const isSupportedBrowser = () => {
  const agent = navigator.userAgent
  for (let b of browsers) {
    let [browserName, expectedVersion] = b
    let output = new RegExp(browserName + '/([0-9.]+)', 'i').exec(agent)
    if (!output) {
      continue
    }
    let currentVersion = parseFloat(output[1])
    if (currentVersion >= expectedVersion) {
      return true
    }
  }
  return false
}

export const findInPairs = (data: Array<any>, key: string) => {
  for (let commodity of data) {
    if (commodity[0] === key) {
      return commodity[1]
    }
  }
  console.warn(`can't find`, key, `in`, data)
  return ''
}
