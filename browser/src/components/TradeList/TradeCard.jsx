// @flow

import React from 'react'
import type { TradeType } from '../../model/flowType'
import { getTradeStatus } from '../../lib/helper'
import { canMakeReqStatus } from '../../constants/tradeConst'

type Props = {
  trade: TradeType,
  tradeNum: number,
  selectTrade: Function,
  detailPath: string,
  history: Object
}

class TradeCardContainer extends React.Component<Props, {}> {
  renderSelectedTrade = () => {
    const { selectTrade, detailPath, history, tradeNum } = this.props
    selectTrade(tradeNum)
    history.push(detailPath)
  }

  renderTradeStatus = (trade: TradeType) => {
    let status = getTradeStatus(trade)
    switch (status) {
      case canMakeReqStatus.approved:
        return (<p><i className={'fas fa-check'} />Completed</p>)
      case canMakeReqStatus.pending:
        return (<p>Waiting for close request approval</p>)
      default:
        return (<p><i className={'fas fa-spinner'} />In Progress</p>)
    }
  }

  render () {
    const { trade } = this.props
    return (
      <div className={'trade-card-body'}>
        <div className={'trade-title'}>
          <p>Title</p>
          <p onClick={this.renderSelectedTrade}>{trade.name}</p>
          <p onClick={this.renderSelectedTrade}>{trade.description}</p>
        </div>
        <div className={'trade-status'}>
          <p>Status</p>
          {
            this.renderTradeStatus(trade)
          }
        </div>
        <div className={'trade-last-action'}>
          <p>Last action</p>
          <p>
            {trade.stages[trade.stages.length - 1].name}
          </p>
          <p>
            {trade.stages[trade.stages.length - 1].description}
          </p>
        </div>
        <div className={'trade-involved'}>
          <p>Who is involved</p>
          <p>{trade.buyer.firstName + ' ' + trade.buyer.lastName}</p>
          <p>{trade.seller.firstName + ' ' + trade.seller.lastName}</p>
        </div>
      </div>
    )
  }
}

export default TradeCardContainer
