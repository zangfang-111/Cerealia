// @flow

import React, { useState } from 'react'
import { observer } from 'mobx-react-lite'
import type { TradeStageAddReqType } from '../../../model/flowType'
import ApproveTradeStageReqModal from '../../Common/Modals/ContractModal/ApproveTradeStageReqModal'
import RejectTradeStageReqModal from '../../Common/Modals/ContractModal/RejectTradeStageReqModal'
import currentUser from '../../../stores/current-user'
import tradeViewStore from '../../../stores/tradeViewStore'
import { tradeActorMap } from '../../../constants/tradeConst'

export default observer(() => {
  const [displayApprove, setDisplayApprove] = useState(false)
  const [displayReject, setDisplayReject] = useState(false)
  const [stageIdx, setStageIdx] = useState(0)

  function onOpenApproveModal (index: number) {
    setDisplayApprove(true)
    setStageIdx(index)
  }

  function onOpenRejectModal (index: number) {
    setDisplayReject(true)
    setStageIdx(index)
  }

  function renderStageAddAction (stage: TradeStageAddReqType, index: number) {
    if (stage.reqBy.id === currentUser.user.id) {
      return `Waiting for ${stage.reqActor === 's' ? 'buyer' : 'seller'}'s approval`
    }
    return (
      <p>
        <i className='fas fa-check-circle hoverable'
          onClick={() => onOpenApproveModal(index)}>Approve</i>
        <i className='fas fa-times-circle hoverable'
          onClick={() => onOpenRejectModal(index)}>Reject</i>
      </p>
    )
  }

  const stageAddReqs = tradeViewStore.stageAddReqs
  return (
    <div className={'additional-stage'}>
      { stageAddReqs &&
        stageAddReqs.map((stageAddReq, index) =>
          renderAdditionalStage(stageAddReq, index, renderStageAddAction))
      }
      { displayApprove && (
        <ApproveTradeStageReqModal
          onCloseModal={() => setDisplayApprove(false)}
          stageIdx={stageIdx}
          visible={displayApprove}
        />
      )}
      { displayReject && (
        <RejectTradeStageReqModal
          onCloseModal={() => setDisplayReject(false)}
          stageIdx={stageIdx}
          visible={displayReject}
        />
      )}
    </div>
  )
})

const renderAdditionalStage = (stage: TradeStageAddReqType, index: number,
  renderStageAddAction: Function) => {
  return (
    <div className={`stage-content ${stage.status}`} key={index}>
      <table>
        <thead>
          <tr>
            <td>Title</td>
            <td>Stage Owner</td>
            <td>reason</td>
            <td>created By</td>
            <td>Action</td>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td className={'title'}>
              <span>{stage.name}</span>
              <span>{stage.description}</span>
            </td>
            <td className={'owner'}>
              <span>{tradeActorMap[stage.owner]}</span>
            </td>
            <td className={'reason'}>
              <span>{stage.reqReason}</span>
            </td>
            <td className={'created-by'}>
              <span>{`${stage.reqBy.firstName} ${stage.reqBy.lastName}`}</span>
            </td>
            <td className={'action'}>
              {
                renderStageAddAction(stage, index)
              }
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  )
}
