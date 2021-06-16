// @flow

import React, { useState } from 'react'
import { observer } from 'mobx-react-lite'
import type { TradeStageDocType } from '../../../model/flowType'
import { approveStatus, warningExpireTime, locationMap } from '../../../constants/tradeConst'
import { Tooltip } from 'antd'
import { toLocalTime } from '../../../lib/helper'
import RejectedDocumentModal from '../../Common/Modals/ContractModal/RejectedDocumentModal'
import ApproveTradeStageDocModal from '../../Common/Modals/ContractModal/ApproveTradeStageDocModal'
import RejectTradeStageDocModal from '../../Common/Modals/ContractModal/RejectTradeStageDocModal'
import currentUser from '../../../stores/current-user'
import appStore from '../../../stores/app-store'
import * as moment from 'moment-timezone'
import { openDocFile } from '../../../lib/downloader'

type Props = {
  stageDoc: TradeStageDocType,
  stageIdx: number
}

export default observer((props: Props) => {
  const [displayApproveStageDocModal, setDisplayApproveStageDocModal] = useState(false)
  const [displayRejectStageDocModal, setDisplayRejectStageDocModal] = useState(false)
  const [displayRejectedDocumentModal, setDisplayRejectedDocumentModal] = useState(false)
  const { stageDoc, stageIdx } = props

  function renderReason () {
    switch (stageDoc.status) {
      case approveStatus.pending:
        return '--'
      case approveStatus.approved:
        return '--'
      case approveStatus.rejected:
        return (
          <p className={'show-reason hoverable'}
            onClick={() => setDisplayRejectedDocumentModal(true)}>
            <i className='fas fa-exclamation-circle' />
            Show reason
          </p>
        )
      default:
        return null
    }
  }

  function renderAction () {
    if (stageDoc.status !== approveStatus.pending ||
      currentUser.user.id === stageDoc.doc.createdBy.id ||
      appStore.adminMode) {
      return '--'
    }
    return (
      <React.Fragment>
        <p className={'approve-action hoverable'}
          onClick={() => setDisplayApproveStageDocModal(true)}>
          <i className='fas fa-check-circle' />
          approve
        </p>
        <p className={'reject-action hoverable'}
          onClick={() => setDisplayRejectStageDocModal(true)}>
          <i className='fas fa-times-circle' />
          reject
        </p>
      </React.Fragment>
    )
  }

  return (
    <div className={'doc-content'}>
      <div className={'activated'}>
        {
          stageDoc.status === approveStatus.pending &&
          <i className={'fas fa-bell doc-activated'} />
        }
      </div>
      <div className={'doc-title'}>
        <Tooltip placement={'topLeft'} title={stageDoc.doc.name}>
          <p className={'doc-name'}>{stageDoc.doc.name}</p>
        </Tooltip>
        <span className={'doc-download'}>
          <i className='fas fa-file-contract' />
          <a onClick={() => openDocFile(stageDoc.doc.id, stageDoc.doc.name, locationMap.trade)}>
              Download
          </a>
        </span>
      </div>
      <div className={'doc-submitted'}>
        <p>{`${stageDoc.doc.createdBy.firstName} ${stageDoc.doc.createdBy.lastName}`}</p>
        <p>{toLocalTime(stageDoc.doc.createdAt)}</p>
      </div>
      <div className={`doc-status ${stageDoc.status}`}>
        {renderDocStatus(stageDoc)}
      </div>
      <div className={'doc-action'}>
        {renderAction()}
      </div>
      <div className={'doc-reason'}>
        {renderReason()}
      </div>
      <div className={'doc-expire'}>
        {renderDocExpired(stageDoc)}
      </div>
      <ApproveTradeStageDocModal
        visible={displayApproveStageDocModal}
        onCloseModal={() => setDisplayApproveStageDocModal(false)}
        stageDoc={stageDoc}
        stageIdx={stageIdx} />
      <RejectTradeStageDocModal
        visible={displayRejectStageDocModal}
        onCloseModal={() => setDisplayRejectStageDocModal(false)}
        stageDoc={stageDoc}
        stageIdx={stageIdx} />
      {
        stageDoc.status === approveStatus.rejected &&
        <RejectedDocumentModal
          visible={displayRejectedDocumentModal}
          onCloseModal={() => setDisplayRejectedDocumentModal(false)}
          stageDoc={stageDoc}
          stageIdx={stageIdx}
        />
      }
    </div>
  )
})

function renderDocStatus (stageDoc: TradeStageDocType) {
  switch (stageDoc.status) {
    case approveStatus.pending:
      return <p>
        <i className='fas fa-check-circle' />Waiting for approval
      </p>
    case approveStatus.approved:
      return (
        <p className={'approved'}>
          <i className='fas fa-check-circle' />Approved by
          { stageDoc.approvedBy
            ? `${stageDoc.approvedBy.firstName} ${stageDoc.approvedBy.lastName}` : 'me' }
          {toLocalTime(stageDoc.approvedAt)}
        </p>
      )
    case approveStatus.submitted:
      return (
        <p className={'approved'}>
          <i className='fas fa-check-circle' />Submitted
        </p>
      )
    case approveStatus.rejected:
      return (
        <p className={'rejected'}>
          <i className='fas fa-times-circle' />Rejected by<br />
          { stageDoc.approvedBy
            ? `${stageDoc.approvedBy.firstName} ${stageDoc.approvedBy.lastName}` : 'me' }
          {toLocalTime(stageDoc.approvedAt)}
        </p>
      )
    case approveStatus.expired:
      return (
        <p>
          <i className='fas fa-times-circle'> Expired </i>
          {toLocalTime(stageDoc.expiresAt)}
        </p>
      )
    default:
      return null
  }
}

function renderDocExpired (stageDoc: TradeStageDocType) {
  if (stageDoc.status !== approveStatus.pending) {
    return '--'
  }
  let utcExpiresAt = moment.utc(stageDoc.expiresAt)
  return utcExpiresAt.diff(moment.utc(), 'hours') < warningExpireTime
    ? <p className='warn-doc-expired'>{utcExpiresAt.local().format()}</p>
    : <p className={'expire-time'}>{utcExpiresAt.local().format()}</p>
}
