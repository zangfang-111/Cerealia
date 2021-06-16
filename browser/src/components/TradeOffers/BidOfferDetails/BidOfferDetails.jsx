// @flow

import React, { useState } from 'react'
import { observer } from 'mobx-react-lite'
import CreateNewTradeModal from '../../Common/Modals/ContractModal/CreateNewTradeModal'
import tradeOffersStore from '../../../stores/tradeoffer-store'
import curTradeOfferStore from '../../../stores/current-tradeoffer'
import { findInPairs, toLocalTime } from '../../../lib/helper'
import { commodities, commodityTypes } from '../../../constants/selectOptions'
import { openDocFile } from '../../../lib/downloader'
import { locationMap } from '../../../constants/tradeConst'

function renderComType (com: string) {
  return (<div className={'type-item'}>
    {findInPairs(commodityTypes, com)}
  </div>)
}

export default observer(() => {
  // TODO
  const [displayCreateNewTradeModal, setDisplayCreateNewTradeModal] = useState(false)
  const tradeOffer = curTradeOfferStore.tradeOffer

  function handleCreateTradeModal (status) {
    setDisplayCreateNewTradeModal(status)
  }

  return (
    <div className={'bid-details'}>
      <p>
        {tradeOffersStore.offerStatus ? 'Offers' : 'Bids'} /&nbsp;
        {findInPairs(commodities, tradeOffer.commodity)}
      </p>
      <div className={'header-section'}>
        <div className={'title'}>
          <div className={'page-title'}>{findInPairs(commodities, tradeOffer.commodity)}</div>
          <p>{tradeOffer.quality}</p>
        </div>
        <div className={'type'}>
          { tradeOffer.comType.map(renderComType) }
        </div>
        <div className={'open-chat'}>
          <i className='fas fa-comments' />
          <p>Open chat</p>
        </div>
      </div>
      <div className={'body-container'}>
        <div className={'left-body'}>
          <div className={'row'}>
            <div className={'sub-title'}>Price</div>
            <p>${tradeOffer.price}</p>
          </div>
          <div className={'row'}>
            <div className={'sub-title'}>Origin</div>
            <p>{tradeOffer.origin}</p>
          </div>
          <div className={'row'}>
            <div className={'sub-title'}>Basis</div>
            <p>{tradeOffer.incoterm}</p>
          </div>
          <div className={'row'}>
            <div className={'sub-title'}>QTY</div>
            <p>{tradeOffer.vol}t</p>
          </div>
          <div className={'row'}>
            <div className={'sub-title'}>Shipment</div>
            <div className={'shipment'}>
              <p>from {toLocalTime(tradeOffer.shipment[0])}</p>
              <p> till {toLocalTime(tradeOffer.shipment[1])}</p>
            </div>
          </div>
          <div className={'row-bottom'}>
            <div className={'sub-title'}>Submitted</div>
            <p>{toLocalTime(tradeOffer.createdAt)}</p>
          </div>
        </div>
        <div className={'right-body'}>
          <div className={'row'}>
            <div className={'sub-title'}>Status</div>
            <p>{tradeOffer.priceType}</p>
          </div>
          <div className={'row'}>
            <div className={'sub-title'}>Expires</div>
            <p>{toLocalTime(tradeOffer.expiresAt)}</p>
          </div>
          <div className={'row'}>
            <div className={'sub-title'}>Company</div>
            {
              tradeOffer.org &&
              <p>{tradeOffer.org.name}</p>
            }
          </div>
          <div className={'row'}>
            <div className={'sub-title'}>Contract terms</div>
            {
              tradeOffer.terms &&
              <div className={'download-btn'}>
                <i className='fas fa-download' />
                <a onClick={() => openDocFile(
                  tradeOffer.terms.id, tradeOffer.terms.name, locationMap.tradeOffer
                )}>
                  Download
                </a>
              </div>
            }
          </div>
          <div className={'row'}>
            <div className={'sub-title'}>Payment terms</div>
            <p>{tradeOffer.note}</p>
          </div>
          <div className={'contract-btn'} onClick={() => handleCreateTradeModal(true)}>
            <p>Create a contract</p>
            <i className='fas fa-plus-circle' />
          </div>
        </div>
      </div>
      <CreateNewTradeModal
        visible={displayCreateNewTradeModal}
        onCloseModal={() => handleCreateTradeModal(false)}
        withTradeOffer
      />
    </div>
  )
})
