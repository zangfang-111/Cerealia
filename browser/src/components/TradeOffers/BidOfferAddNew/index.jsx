// @flow

import React from 'react'
import axios from 'axios'
import { observer } from 'mobx-react-lite'
import useReactRouter from 'use-react-router'
import { Form, Spin } from 'antd'
import NewTradeOfferForm from './BidOfferAddNew'
import { addNotificationHelper, mkFormDataHeaders } from '../../../lib/helper'
import spinnerStore from '../../../stores/spinner-store'
import tradeOffersStore from '../../../stores/tradeoffer-store'
import type { TradeOfferInput } from '../../../model/flowType'
import { mkLink } from '../../../services/cerealia'

export default observer(() => {
  const { history } = useReactRouter()

  async function onSubmit (input: TradeOfferInput, file: Object) {
    spinnerStore.showSpinner('please wait...')
    try {
      if (file) {
        let data = new FormData()
        data.append('formfile', file)
        let response = await axios.post(
          mkLink('v1/trade-offers'),
          data,
          {
            headers: mkFormDataHeaders()
          })
        if (!response.data) {
          addNotificationHelper('Attached file upload failed.', 'error')
          return
        }
        input.docID = response.data.toString()
      }
      await tradeOffersStore.createTradeOffer(input)
      history.push(`/trade-offer/${tradeOffersStore.offerStatus ? 'sell' : 'buy'}`)
      addNotificationHelper('Trade Offer created successfully', 'success')
    } catch (err) {
      addNotificationHelper(err, 'error')
    }
    spinnerStore.hideSpinner()
  }

  const OfferForm = Form.create()(NewTradeOfferForm)
  return (
    <Spin spinning={spinnerStore.loading} size='large' tip={spinnerStore.spinnerTip}>
      <OfferForm offerType={tradeOffersStore.offerStatus} onSubmit={onSubmit} />
    </Spin>
  )
})
