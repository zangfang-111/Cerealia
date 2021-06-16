import { canMakeReqStatus } from '../constants/tradeConst'

export const mkEmptyUser = () => {
  return {
    id: '',
    firstName: '',
    lastName: '',
    emails: [],
    orgMap: [{
      org: {
        id: '',
        name: '',
        telephone: '',
        email: '',
        address: ''
      },
      role: ''
    }],
    roles: [],
    pubKey: '',
    createdAt: new Date(),
    avatar: '',
    biography: ''
  }
}

export const mkEmptyDoc = () => {
  return {
    doc: {
      id: '',
      name: '',
      hash: '',
      note: '',
      url: '',
      type: '',
      createdBy: mkEmptyUser(),
      createdAt: new Date()
    },
    tradeCloseStatus: canMakeReqStatus.can,
    approvedBy: mkEmptyUser(),
    approvedAt: new Date(),
    rejectReason: '',
    expiresAt: new Date(),
    closeReqs: [],
    status: 'pending',
    reqTx: '',
    index: 0
  }
}

export const mkEmptyTemplate = () => {
  return {
    id: '',
    name: '',
    description: '',
    stages: []
  }
}

export const mkEmptyTrade = () => {
  return {
    id: '',
    name: '',
    description: '',
    template: mkEmptyTemplate(),
    buyer: mkEmptyUser(),
    seller: mkEmptyUser(),
    tradeCloseStatus: 'no',
    scAddr: '',
    stages: [],
    stageAddReqs: [],
    closeReqs: [],
    createdAt: new Date(),
    createdBy: mkEmptyUser()
  }
}
