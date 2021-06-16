// @flow

import React from 'react'
import { Spin } from 'antd'
import _ from '../../polyfills/underscore.js'
import { observer } from 'mobx-react-lite'
import useReactRouter from 'use-react-router'
import EmptyTrades from '../Trade/EmptyTrade'
import TradeCard from './TradeCard'
import tradesStore from '../../stores/trade-store'
import modTradesStore from '../../stores/moderatorStore/modTrades'

type Props = {
  location: Object
}

export default observer((props: Props) => {
  const { history } = useReactRouter()
  const trades = tradesStore.trades
  const spinning = !tradesStore.triedToFetchTrades
  const auth = props.location.auth === 'admin'
  const selectTrade = auth ? modTradesStore.setSelectedTab : tradesStore.setSelectedTab
  const detailPath = auth ? '/admin/home' : '/home'

  return (
    <div className={'trades-list'}>
      <Spin spinning={spinning} tip={'Loading trade data...'}>
        <p className={'title'}>Trade List</p>
        <div className={'trades-list-content'}>
          {
            _.isEmpty(trades)
              ? <EmptyTrades />
              : trades.map((trade, i) => (
                <TradeCard
                  key={i}
                  tradeNum={i}
                  trade={trade}
                  selectTrade={selectTrade}
                  detailPath={detailPath}
                  history={history}
                />)
              )
          }
        </div>
      </Spin>
    </div>
  )
})
