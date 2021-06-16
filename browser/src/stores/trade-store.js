// @flow

import { action, observable, computed, runInAction } from 'mobx'
import { getTemplates,
  createTrade,
  getTradeData
} from '../graphql/trades'
import type { NewTradeInputType, TradeTemplateType, TradeType } from '../model/flowType'
import tradeViewStore from './tradeViewStore'
import { GqlClient } from '../services/cerealia'

class TradesStore {
  @observable trades: Array<TradeType> = []
  @observable tradeTemplates: Array<TradeTemplateType> = []
  @observable selectedTab: number = 0
  @observable triedToFetchTrades: boolean = false

  gqlClient: Object
  constructor () {
    this.gqlClient = GqlClient
  }

  @action async fetchTrades () {
    this.triedToFetchTrades = false
    let response = await this.gqlClient.query({
      query: getTradeData,
      fetchPolicy: 'network-only'
    })
    runInAction('fetchSuccess', () => {
      this.trades = response.data.trades
      this.triedToFetchTrades = true
      tradeViewStore.updateTrade(this.getCurTrade, this.gqlClient)
    })
  }

  @action initializeTrades () {
    this.trades = []
    this.selectedTab = 0
    this.triedToFetchTrades = false
  }

  @action async createTrade (input: NewTradeInputType) {
    let response = await this.gqlClient.mutate({
      mutation: createTrade,
      variables: { 'input': input }
    })
    runInAction('fetchSuccess', () => {
      this.trades.push(response.data.tradeCreate)
      this.trades.length && this.setSelectedTab(this.trades.length - 1)
      tradeViewStore.updateTrade(this.getCurTrade, this.gqlClient)
    })
  }

  @action async fetchTradeTemplates () {
    try {
      let response = await this.gqlClient.query({
        query: getTemplates,
        variables: {}
      })
      runInAction('fetchSuccess', () => {
        this.tradeTemplates = response.data.tradeTemplates
      })
    } catch (error) {
      console.error('Get trade Templates error:', error)
    }
  }

  @action setSelectedTab = (index: number) => {
    this.selectedTab = index
    tradeViewStore.updateTrade(this.getCurTrade, this.gqlClient)
  }

  @computed get getCurTrade (): any {
    if (this.trades.length === 0) {
      return {}
    }
    return this.trades[this.selectedTab]
  }
}

export default new TradesStore()
