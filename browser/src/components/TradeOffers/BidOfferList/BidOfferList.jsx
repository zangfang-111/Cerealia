// @flow

import React from 'react'
import { observer } from 'mobx-react-lite'
import useReactRouter from 'use-react-router'
import { Spin } from 'antd'
import moment from 'moment'
import tradeOffersStore from '../../../stores/tradeoffer-store'
import curTradeOfferStore from '../../../stores/current-tradeoffer'
import BidOffersHeader from './BidOffersHeader'
import _ from '../../../polyfills/underscore.js'
import type { TradeOfferType } from '../../../model/flowType'
import { findInPairs, toLocalTime } from '../../../lib/helper'
import { commodities, commodityTypes } from '../../../constants/selectOptions'

type Props = {
  offerType: string,
  tradeOffer: TradeOfferType,
  match: Object
}

export default observer((props: Props) => {
  const { history } = useReactRouter()

  // onClick: offer row click
  function handleBidOfferDetails (tradeoffer) {
    curTradeOfferStore.setCurTradeOffer(tradeoffer)
    history.push('/trade-offer/details')
  }

  function displayComType (companyType: Array<string>) {
    let comType = []
    companyType.map(com => {
      comType.push(findInPairs(commodityTypes, com))
    })
    return comType.join(', ')
  }

  function displayShipment (shipmentValue: Array<any>) {
    return moment(shipmentValue[0]).utc().format('D MMM') + ' - ' + moment(shipmentValue[1]).utc().format(`D MMM 'YY`)
  }

  // parameterized the route to `trade-offer/:offer` and resticted to value buy|sell
  const offerType = props.match.params.offer
  const isSell = offerType === 'sell'
  function showTradeOfferRows () {
    const tradeOffers = tradeOffersStore.tradeOffers.filter(item =>
      item.isSell === isSell)
    tradeOffersStore.setOfferStatus(isSell)
    if (_.isEmpty(tradeOffers)) {
      return <div className={'empty-field'}>There are no active bids.</div>
    }
    return tradeOffers.map((item, key) => (
      <div key={key} value={item} className={'bid-row row-body'} onClick={() => handleBidOfferDetails(item)} >
        <div className={'row-goods'}>
          <p className={'goods-title'}>{findInPairs(commodities, item.commodity)}</p>
          <p>{displayComType(item.comType)}, {item.quality}</p>
        </div>
        <p>{item.origin}</p>
        <p>{item.incoterm}, {item.marketLoc}</p>
        <p>{item.vol}</p>
        <p>{displayShipment(item.shipment)}</p>
        <p>{toLocalTime(item.createdAt)}</p>
      </div>
    ))
  }
  const spinning = !tradeOffersStore.triedToFetchTradeOffers
  const tradeOfferType = tradeOffersStore.offerStatus ? 'Offers' : 'Bids'

  return (
    <div className={'bid-list'}>
      <Spin spinning={spinning} tip={`Loading ${tradeOfferType} data...`}>
        <BidOffersHeader />
        <div className={'bid-list-content'}>
          <div className={'bid-row'}>
            <div className={'row-goods'}>Goods</div>
            <p>Origin</p>
            <p>Basis</p>
            <p>QTY(t)</p>
            <p>Shipment</p>
            <p>Submitted</p>
          </div>
          <div>
            { showTradeOfferRows() }
          </div>
        </div>
      </Spin>
    </div>
  )
})
