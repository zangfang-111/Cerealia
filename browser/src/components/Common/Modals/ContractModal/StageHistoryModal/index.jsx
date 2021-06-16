// @flow

import { Table } from 'antd'
import React from 'react'
import { toLocalTime } from '../../../../../lib/helper'
import { approveStatus } from '../../../../../constants/tradeConst'
import tradeViewStore from '../../../../../stores/tradeViewStore'
import Button from '../../../Button/Button'
import BaseModal from '../BaseModal'
import type { StageModalPropsType } from '../types'

class StageHistoryModal extends React.Component<StageModalPropsType> {
  renderStageHistory = () => {
    const { stageIdx } = this.props
    const columns = [{
      title: 'No',
      dataIndex: 'index',
      key: 'index'
    }, {
      title: 'By',
      dataIndex: 'name',
      key: 'name'
    }, {
      title: 'Date',
      dataIndex: 'date',
      key: 'date'
    }, {
      title: 'Reason',
      dataIndex: 'reason',
      key: 'reason',
      width: 232
    }]
    let closeReqs = tradeViewStore.stages[stageIdx].closeReqs
    let closeReqItems = []
    closeReqs.length > 0 && closeReqs.map((req, index) => {
      if (req.status === approveStatus.rejected) {
        closeReqItems.push({
          key: index,
          index: index + 1,
          name: req.approvedBy.firstName + ' ' + req.approvedBy.lastName,
          date: toLocalTime(req.approvedAt),
          reason: req.rejectReason
        })
      }
    })
    return (
      <div>
        <p>Rejected stage "<span className={'highlight-text'}>
          {tradeViewStore.stages[stageIdx].name}</span>" close requests</p>
        <Table columns={columns} dataSource={closeReqItems} />
        {
          closeReqItems.length === 0 && <p>There is no rejected close request</p>
        }
      </div>
    )
  }
  render () {
    const { onCloseModal, visible } = this.props
    return (
      <BaseModal visible={visible} onCloseModal={onCloseModal} status={'history'}>
        <div className={'confirm-closing-stage'}>
          <p className={'modal-title'}>Stage Close History</p>
          {
            this.renderStageHistory()
          }
          <div className={'bottom-right-button'}>
            <Button type={'primary'} text={'Close'} onClick={onCloseModal} />
          </div>
        </div>
      </BaseModal>
    )
  }
}

export default StageHistoryModal
