// @flow

import React, { useState } from 'react'
import { Switch } from 'antd'
import { observer } from 'mobx-react-lite'
import DragSortableList from 'react-drag-sortable'
import type { TradeStageType } from '../../../model/flowType'
import AddStageRequests from '../Stage/AddStage'
import MainStage from '../Stage/Stage'
import UserCard from '../../Common/UserCard/UserCard'
import { canMakeReqStatus, canCloseTooltipText,
  TradeStageStatusMap } from '../../../constants/tradeConst'
import Button from '../../Common/Button/Button'
import CloseTradeModal from '../../Common/Modals/ContractModal/CloseTradeModal'
import ConfirmClosingTradeModal from '../../Common/Modals/ContractModal/ConfirmClosingTradeModal'
import AddNewStageModal from '../../Common/Modals/ContractModal/AddNewStageModal'
import stellarStore from '../../../stores/stellarStore'
import currentUser from '../../../stores/current-user'
import appStore from '../../../stores/app-store'
import tradeViewStore from '../../../stores/tradeViewStore'

export default observer(() => {
  const [showDeleteStages, setShowDeleteStages] = useState(false)
  const [showCompletedStages, setShowCompletedStages] = useState(false)
  const [displayAddNewStageModal, setDisplayAddNewStageModal] = useState(false)

  const mainStageClass = showCompletedStages
    ? TradeStageStatusMap.pending : TradeStageStatusMap.completed
  const deletedStageClass = TradeStageStatusMap.deleted
  const stagesRendered = tradeViewStore.stages &&
    tradeViewStore.stages.map((stage: TradeStageType, i: number) => {
      return { content: (
        <MainStage
          key={i}
          stage={stage}
          stageIdx={i}
          stageType={mainStageClass}
        />
      ) }
    })

  return (
    <div className={'trade-container main-trade-section'}>
      <div className={'user-information'}>
        <UserCard
          actor='Buyer'
          name={`${tradeViewStore.buyer.firstName} ${tradeViewStore.buyer.lastName}`}
          orgName={tradeViewStore.buyer.orgMap[0].org.name}
          role={tradeViewStore.buyer.orgMap[0].role}
          avatar={tradeViewStore.buyer.avatar}
        />
        <UserCard
          actor='Seller'
          name={`${tradeViewStore.seller.firstName} ${tradeViewStore.seller.lastName}`}
          orgName={tradeViewStore.seller.orgMap[0].org.name}
          role={tradeViewStore.seller.orgMap[0].role}
          avatar={tradeViewStore.seller.avatar}
        />
      </div>
      <div className={'trade-header'}>
        <div className={'left-column'}>
          <p className={'trade-name'}>{ tradeViewStore.name}</p>
          <p> { tradeViewStore.description } </p>
          <p className='smart-contract'>Smart Contract address: &nbsp;
            <a href={stellarStore.linkToHorizonOperations(tradeViewStore.scAddr)} target='_blank'>
              { tradeViewStore.scAddr }
            </a>
            <br />
            [<a href={stellarStore.linkToSteexp(tradeViewStore.scAddr)} target='_blank' >
               operations explorer
            </a>]
          </p>
        </div>
        <div className={'right-column'}>
          <div className={'close-trade-button'}>
            {
              appStore.adminMode
                ? renderAdminTradeCloseStatus()
                : useRenderTradeCloseStatus()
            }
          </div>
          <div className={'switching-body'}>
            <div className={'switching-stage-title'}>Show deleted stages</div>
            <div className={'switching-section'}>
              <Switch
                className={'show-switch'}
                checkedChildren='ON'
                unCheckedChildren='OFF'
                onChange={() => setShowDeleteStages(!showDeleteStages)}
              />
            </div>
          </div>
          <div className={'switching-body'}>
            <div className={'switching-stage-title'}>Show completed stages</div>
            <div className={'switching-section'}>
              <Switch
                className={'show-switch'}
                checkedChildren='ON'
                unCheckedChildren='OFF'
                onChange={() => setShowCompletedStages(!showCompletedStages)}
              />
            </div>
          </div>
        </div>
      </div>
      <div className={'trade-content'}>
        <AddStageRequests />
        <DragSortableList items={stagesRendered || []}
          moveTransitionDuration={0.3} onSort={() => {}} type='vertical' />
        <React.Fragment>
          <hr />
          <p className={`deleted-section-title deleted ${showDeleteStages ? 'fadein' : 'fadeout'}`}>Deleted Stages</p>
          {
            tradeViewStore.stages.map((stage: TradeStageType, i: number) => {
              return (
                <MainStage
                  key={i}
                  stage={stage}
                  stageIdx={i}
                  userId={currentUser.user.id}
                  stageType={deletedStageClass}
                  showDeleteStages={showDeleteStages}
                />
              )
            })
          }
        </React.Fragment>
        {
          (tradeViewStore.tradeCloseStatus === canMakeReqStatus.no ||
            tradeViewStore.tradeCloseStatus === canMakeReqStatus.can) &&
            <div className={'btn-group hoverable'}
              onClick={() => setDisplayAddNewStageModal(true)}>
              <p>Add new stage</p>
              <i className='fas fa-plus-circle' />
            </div>
        }
      </div>
      <AddNewStageModal
        onCloseModal={() => setDisplayAddNewStageModal(false)}
        visible={displayAddNewStageModal}
      />
    </div>
  )
})

function renderAdminTradeCloseStatus () {
  let content = (<p className={'trade-status green'}>This trade is in progress</p>)
  switch (tradeViewStore.tradeCloseStatus) {
    case canMakeReqStatus.approved:
      content = (<p className={'trade-status green'}>This trade has been completed</p>)
      break
    default:
      break
  }
  return content
}

function useRenderTradeCloseStatus () {
  const [displayCloseTradeModal, setDisplayCloseTradeModal] = useState(false)
  const [displayConfirmTradeCloseModal, setDisplayConfirmTradeCloseModal] = useState(false)

  const showTradeStatus = () => {
    switch (tradeViewStore.tradeCloseStatus) {
      case canMakeReqStatus.can:
        return (
          <Button onClick={() => setDisplayCloseTradeModal(true)}
            text={'Complete trade'} type={'primary'}
            tooltipText={canCloseTooltipText} />)
      case canMakeReqStatus.no:
        return (
          <Button text={'Complete trade'} type={'primary'}
            disabled tooltipText={canCloseTooltipText} />)
      case canMakeReqStatus.pending:
        if (tradeViewStore.closeReqs[tradeViewStore.closeReqs.length - 1].reqBy.id ===
          currentUser.user.id) {
          return (<p className={'trade-status'}>Waiting for counterParty's approval</p>)
        }
        return (
          <Button onClick={() => setDisplayConfirmTradeCloseModal(true)}
            text={'Approve Trade Completion'} type={'primary'} />)
      case canMakeReqStatus.approved:
        return (<p className={'trade-status green'}>This trade has been completed</p>)
      default:
        return null
    }
  }

  return (
    <React.Fragment>
      {showTradeStatus()}
      <CloseTradeModal
        onCloseModal={() => setDisplayCloseTradeModal(false)}
        visible={displayCloseTradeModal}
      />
      <ConfirmClosingTradeModal
        onCloseModal={() => setDisplayConfirmTradeCloseModal(false)}
        visible={displayConfirmTradeCloseModal}
      />
    </React.Fragment>
  )
}
