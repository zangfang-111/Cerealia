// @flow

import React, { useState } from 'react'
import { observer } from 'mobx-react-lite'
import { Popover, Tooltip } from 'antd'
import * as moment from 'moment-timezone'
import {
  approveStatus,
  canMakeReqStatus, stageStatusMap,
  tradeActorMap, tradeContractName,
  warningExpireTime
} from '../../../constants/tradeConst'
import type { TradeStageType } from '../../../model/flowType'
import CloseStageModal from '../../Common/Modals/ContractModal/CloseStageModal'
import DeleteStageModal from '../../Common/Modals/ContractModal/DeleteStageModal'
import ConfirmDeletingStageModal from '../../Common/Modals/ContractModal/ConfirmDeletingStageModal'
import ConfirmClosingStageModal from '../../Common/Modals/ContractModal/ConfirmClosingStageModal'
import StageHistoryModal from '../../Common/Modals/ContractModal/StageHistoryModal'
import ExpireTime from './ExpireTime'
import appStore from '../../../stores/app-store'
import tradeViewStore from '../../../stores/tradeViewStore'

type Props = {
  stage: TradeStageType,
  stageIdx: number,
  stageStatus: string,
}

export default observer((props: Props) => {
  const [displayHistoryModal, setDisplayHistoryModal] = useState(false)
  const { stage, stageIdx, stageStatus } = props
  return (
    <div className={'stage-header'}>
      <div className={'activated'}>
        {
          stage.moderator.user && stage.moderator.user.id &&
          <Tooltip placement='topLeft' title={`This stage was created by Cerealia moderator: ${stage.moderator.user.firstName}`}>
            <i className='fas fa-user-shield' />
          </Tooltip>
        }
        {
          hasNotification(stage) && <i className='fas fa-bell' />
        }
      </div>
      <div className={'title'}>
        Title<br />
        <span className={'header-title'}>{stage.name}</span><br />
        <span>{stage.description}</span>
      </div>
      <div className={'owner'}>
        Owner<br />
        <span className={'header-title'}>
          {
            stage.owner === 'n' ? '---' : tradeActorMap[stage.owner]
          }
        </span>
      </div>
      <div className={'status'}>
        Status<br />
        { showStageStatus(stage, stageStatus) }
        <a onClick={() => setDisplayHistoryModal(true)}>
          <i className={'fas fa-history'} />history
        </a>
        <StageHistoryModal
          visible={displayHistoryModal}
          onCloseModal={() => setDisplayHistoryModal(false)}
          stageIdx={stageIdx}
        />
      </div>
      <div className={'action'}>
        Action<br />
        <div>
          {
            appStore.adminMode
              ? renderAdminStageAction()
              : useRenderStageAction(stage, stageIdx)
          }
        </div>
      </div>
      <div className={'expiring'} >
        Expiring time
        {
          useRenderStageExpireTime(stage, stageIdx)
        }
      </div>
    </div>
  )
})

function hasNotification (stage: TradeStageType) {
  if (stage.stageCloseStatus === canMakeReqStatus.pending ||
    stage.stageDeleteStatus === canMakeReqStatus.pending) {
    return true
  }
  for (let doc of stage.docs) {
    if (doc.status === approveStatus.pending) {
      return true
    }
  }
  return false
}

function showStagePendingStatus (stage: TradeStageType) {
  if (stage.stageDeleteStatus === canMakeReqStatus.pending) {
    return 'Waiting for delete approval'
  }
  if (stage.stageCloseStatus === canMakeReqStatus.pending) {
    return 'Waiting for close approval'
  }
  return 'Completion stage'
}

function showStageStatus (stage: TradeStageType, status: string) {
  switch (status) {
    case stageStatusMap.pending:
    case stageStatusMap.current:
      return <p>
        <i className='fas fa-spinner' />
        <span className={'header-title'}>{showStagePendingStatus(stage)}</span>
      </p>
    case stageStatusMap.approved:
      return <p>
        <i className='fas fa-check-circle' />
        <span className={'header-title'}>{stageStatusMap.approved}</span>
      </p>
    case stageStatusMap.rejected:
      return <p>
        <i className='fas fa-trash' />
        <span className={'header-title'}>{stageStatusMap.rejected}</span>
      </p>
    default:
      return null
  }
}

function renderAdminStageAction () {
  return (
    <div className={'stage-header-action'}>
      <a>
        <i className='fas fa-comments hoverable' />Open chat
      </a>
    </div>
  )
}

function useRenderStageExpireTime (stage: TradeStageType, stageIdx: number) {
  const [displayDatePicker, setDisplayDatePicker] = useState(false)
  let renderItems = []
  // if the user is the owner of the stage and close request is within
  // pending status, then show the set expire time button.
  if (stage.stageCloseStatus === canMakeReqStatus.approved) {
    return (
      <div key={0} className={'disabled'}>--</div>
    )
  }
  if (tradeViewStore.currentRole !== tradeActorMap[stage.owner]) {
    renderItems.push((
      <Popover
        placement='bottom'
        title={'Set Stage Expiring Time'}
        content={<ExpireTime stageIdx={stageIdx}
          onCloseDatePicker={() => setDisplayDatePicker(false)} />}
        trigger='click'
        visible={displayDatePicker}
        onVisibleChange={setDisplayDatePicker}
        key={1}
      >
        <div className={'set-exp-time'}>
          <i className={'fas fa-clock'} />set time
        </div>
      </Popover>
    ))
  }
  let utcExpiresAt = moment.utc(stage.expiresAt)
  if (utcExpiresAt > moment()) {
    let expTime = utcExpiresAt.diff(moment.utc(), 'hours') < warningExpireTime
      ? <p key={2} className='warn-doc-expired'>{utcExpiresAt.local().format()}</p>
      : <p key={2} className={'expire-time'}>{utcExpiresAt.local().format()}</p>
    renderItems.push(expTime)
  } else {
    renderItems.push((<p key={3} >---</p>))
  }
  return (
    <React.Fragment>
      {renderItems}
    </React.Fragment>
  )
}

function useRenderStageAction (stage: TradeStageType, stageIdx: number) {
  const [displayCloseStageModal, setDisplayCloseStageModal] = useState(false)
  const [displayDeleteStageModal, setDisplayDeleteStageModal] = useState(false)
  const [displayConfirmDeletingModal, setDisplayConfirmDeletingModal] = useState(false)
  const [displayConfirmClosingModal, setDisplayConfirmClosingModal] = useState(false)

  let menuItems = []
  menuItems.push((
    <div key={0}>
      <a onClick={() => {}}>
        <i className='fas fa-comments hoverable' />Open chat
      </a>
    </div>
  ))
  // check the stage deleted approval status
  if (stage.stageDeleteStatus === canMakeReqStatus.pending) {
    if (tradeViewStore.currentRole !== tradeActorMap[stage.owner]) {
      menuItems.push((
        <div key={1}>
          <a onClick={() => setDisplayConfirmDeletingModal(true)} className={'action-button'}>
            <i className='fas fa-trash hoverable' />Confirm deleting
          </a>
        </div>
      ))
    }
  }
  // check the stage completed approval status
  if (stage.stageCloseStatus === canMakeReqStatus.pending) {
    if (tradeViewStore.currentRole !== tradeActorMap[stage.owner]) {
      menuItems.push((
        <div key={2}>
          <a onClick={() => setDisplayConfirmClosingModal(true)} className={'action-button'}>
            <i className='fas fa-check-circle hoverable' />Confirm closing
          </a>
        </div>
      ))
    }
  }
  // check the stage action status
  if (tradeViewStore.currentRole === tradeActorMap[stage.owner] || stage.owner === 'n') {
    if (stage.stageDeleteStatus === canMakeReqStatus.can &&
      stage.name !== tradeContractName) {
      menuItems.push((
        <div key={3}>
          <a onClick={() => setDisplayDeleteStageModal(true)} className={'action-button'}>
            <i className='fas fa-trash hoverable' />Delete stage
          </a>
        </div>
      ))
    } else if (stage.stageCloseStatus === canMakeReqStatus.can) {
      menuItems.push((
        <div key={4}>
          <a onClick={() => setDisplayCloseStageModal(true)} className={'action-button'}>
            <i className='fas fa-check-circle hoverable' />Complete stage
          </a>
        </div>
      ))
    }
  }
  return (
    <div className={'stage-header-action'}>
      {menuItems}
      <ConfirmDeletingStageModal
        visible={displayConfirmDeletingModal}
        onCloseModal={() => setDisplayConfirmDeletingModal(false)}
        stageIdx={stageIdx}
      />
      <ConfirmClosingStageModal
        visible={displayConfirmClosingModal}
        onCloseModal={() => setDisplayConfirmClosingModal(false)}
        stageIdx={stageIdx}
      />
      <DeleteStageModal
        onCloseModal={() => setDisplayDeleteStageModal(false)}
        visible={displayDeleteStageModal}
        stageIdx={stageIdx}
      />
      <CloseStageModal
        onCloseModal={() => setDisplayCloseStageModal(false)}
        visible={displayCloseStageModal}
        stageIdx={stageIdx}
      />
    </div>
  )
}
