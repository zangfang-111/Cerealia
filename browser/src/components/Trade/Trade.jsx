// @flow

import React, { useState, useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { Spin } from 'antd'
import _ from '../../polyfills/underscore.js'
import TradeView from './TradeView/TradeView.jsx'
import StellarKeyInputModal from '../Common/Modals/StellarModal/StellarModal'
import EmptyTrades from './EmptyTrade'
import stellarStore from '../../stores/stellarStore'
import currentUser from '../../stores/current-user'
import tradesStore from '../../stores/trade-store'

export default observer(() => {
  const [displayStellarKeyModal, setDisplayStellarKeyModal] = useState(false)
  const isVerified = !stellarStore.keyVerified && currentUser.isAuthenticated
  useEffect(() => {
    if (isVerified) {
      setDisplayStellarKeyModal(true)
    }
  }, [isVerified])
  return (
    <div className={'home'}>
      <Spin spinning={!tradesStore.triedToFetchTrades} tip={'Loading trade data...'}>
        {
          _.isEmpty(tradesStore.trades)
            ? <div className={'home-empty'}>
              <p>Home</p>
              <EmptyTrades />
            </div>
            : <div className={'contract-section'}>
              <div className={'trade'}>
                <TradeView />
              </div>
            </div>
        }
        <StellarKeyInputModal
          onCloseModal={() => setDisplayStellarKeyModal(false)}
          visible={displayStellarKeyModal}
        />
      </Spin>
    </div>
  )
})
