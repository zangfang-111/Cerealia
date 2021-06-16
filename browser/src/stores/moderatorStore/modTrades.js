// @flow

import { action, computed, observable, runInAction } from 'mobx'
import { getAdminTradeData } from '../../graphql/trades'
import type { TradeType } from '../../model/flowType'
import { GqlClient } from '../../services/cerealia'
import tradeViewStore from '../tradeViewStore'

class ModTrades {
  @observable trades: Array<TradeType> = []
  @observable triedToFetchTrades: boolean = false
  @observable selectedTab: number = 0

  gqlClient: Object
  constructor () {
    this.gqlClient = GqlClient
  }

  @action async fetchTrades () {
    this.triedToFetchTrades = false
    let response = await this.gqlClient.query({
      query: getAdminTradeData,
      fetchPolicy: 'network-only'
    })
    runInAction('fetchSuccess', () => {
      this.trades = response.data.adminTrades
      this.triedToFetchTrades = true
    })
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

export default new ModTrades()
