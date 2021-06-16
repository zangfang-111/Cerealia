// @flow

import { action, observable } from 'mobx'
import type { TradeOfferType } from '../model/flowType'
import { GqlClient } from '../services/cerealia'

class CurTradeOffer {
  @observable tradeOffer: TradeOfferType = {}

  gqlClient: Object
  constructor () {
    this.gqlClient = GqlClient
  }

  @action setCurTradeOffer (tradeOffer: TradeOfferType) {
    this.tradeOffer = tradeOffer
  }
}

export default new CurTradeOffer()
