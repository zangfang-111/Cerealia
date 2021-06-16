// @flow

export const sellerActor = 's'
export const buyerActor = 'b'
export const tradeContractName = 'trade contract'
export const approveStatus = {
  nil: 'nil',
  pending: 'pending',
  approved: 'approved',
  rejected: 'rejected',
  expired: 'expired',
  submitted: 'submitted'
}
export const canMakeReqStatus = {
  no: 'no',
  can: 'can',
  pending: 'pending',
  approved: 'approved'
}
export const tradeActorMap = {
  s: 'seller',
  b: 'buyer',
  n: 'noOwner'
}

export const stageStatusMap = {
  pending: 'pending',
  approved: 'completed',
  rejected: 'deleted',
  current: 'current-stage'
}

export const userRoleMap = {
  trader: 'trader',
  moderator: 'moderator'
}

export const locationMap = {
  trade: 'trade',
  tradeOffer: 'tradeOffer'
}

export const timeFormat = 'YYYY-MM-DD HH:mm:ss'
export const RegPhone = new RegExp(/^[+]?[(]?[0-9]{1,4}[)]?[-\s0-9]*$/)

// within 12h of expires time left, the warning displays
export const warningExpireTime = 24

export const canCloseTooltipText = 'Trade can be completed only when all stages have been completed'

export const TradeStageStatusMap = Object.freeze({
  pending: 'drag_drop',
  completed: 'drag_drop hide_completed',
  deleted: 'no_drag_drop'
})

export const currencies = ['USD', 'EUR', 'RUB', 'TRY', 'UAH', 'CNY']

export const dateFormat = 'YYYY-MM-DD'
