// @flow

import { action, observable, runInAction } from 'mobx'
import {
  createTradeOffer,
  getTradeOffers
} from '../graphql/trades'
import type { TradeOfferInput, TradeOfferType } from '../model/flowType'
import { GqlClient } from '../services/cerealia'

class TradeOffersStore {
  @observable tradeOffers: Array<TradeOfferType> = []
  @observable triedToFetchTradeOffers: boolean = false
  @observable offerStatus: boolean = false

  gqlClient: Object
  constructor () {
    this.gqlClient = GqlClient
  }

  @action async fetchTradeOffers () {
    this.triedToFetchTradeOffers = false
    let response = await this.gqlClient.query({
      query: getTradeOffers,
      fetchPolicy: 'network-only'
    })
    runInAction('fetchSuccess', () => {
      this.tradeOffers = response.data.tradeOffers
      this.triedToFetchTradeOffers = true
    })
  }

  @action async createTradeOffer (input: TradeOfferInput) {
    let response = await this.gqlClient.mutate({
      mutation: createTradeOffer,
      variables: { 'input': input }
    })
    runInAction('fetchSuccess', () => {
      this.tradeOffers.push(response.data.tradeOfferCreate)
    })
  }

  @action setOfferStatus (status: boolean) {
    this.offerStatus = status
  }
}

export default new TradeOffersStore()
