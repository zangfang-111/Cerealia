// @flow

import React, { useState } from 'react'
import { observer } from 'mobx-react-lite'
import { Button } from 'antd'
import { Collapse } from 'react-collapse'
import type { TradeStageDocType, TradeStageStatusType,
  TradeStageType } from '../../../model/flowType'
import AddNewDocumentModal from '../../Common/Modals/ContractModal/AddNewDocumentModal'
import StageDocView from '../Doc/Doc'
import StageHeader from './StageHeader'
import appStore from '../../../stores/app-store'
import { approveStatus, canMakeReqStatus, stageStatusMap, tradeActorMap } from '../../../constants/tradeConst'
import tradeViewStore from '../../../stores/tradeViewStore'

type Props = {
  stage: TradeStageType,
  stageIdx: number,
  stageType: TradeStageStatusType,
  showDeleteStages: boolean
}

export default observer((props: Props) => {
  const [isOpenCollapse, setIsOpenCollapse] = useState(false)
  const [displayAddNewDocumentModal, setDisplayAddNewDocumentModal] = useState(false)
  const stageStatus = computeStageStatus(props.stage, props.stageIdx)

  let fade = ''
  if (stageStatus === 'deleted') {
    fade = props.showDeleteStages ? 'fadein' : 'fadeout'
  }
  if (stageStatus === 'completed') {
    fade = props.stageType === 'drag_drop hide_completed' ? 'fadeout' : 'fadein'
  }

  return (
    <div className={`trade-stage ${stageStatus} ${props.stageType} ${fade}`} >
      <StageHeader {...props} stageStatus={stageStatus} />
      <Collapse isOpened={isOpenCollapse} className={'document-detail'}>
        {
          !isStageFinished(props.stage) &&
          !appStore.adminMode &&
          (tradeViewStore.currentRole === tradeActorMap[props.stage.owner] || props.stage.owner === 'n') &&
          <p className={'add-new-document'}
            onClick={() => setDisplayAddNewDocumentModal(true)}>
            <i className='fas fa-plus-circle' />
            Submit new document
          </p>
        }
        <div className={'doc-header'}>
          <div className={'activated'} />
          <div className={'doc-title'}>Title</div>
          <div className={'doc-submitted'}>Submitted By</div>
          <div className={'doc-status'}>Status</div>
          <div className={'doc-action'}>Action</div>
          <div className={'doc-reason'}>Reason</div>
          <div className={'doc-expire'}>Expiring time</div>
        </div>
        {
          props.stage.docs && sortStageDoc(props.stage).map((doc) => (
            <StageDocView {...props} stageDoc={doc} key={doc.index} />
          ))
        }
      </Collapse>
      <Button className={'stage-header-expand'}
        onClick={() => setIsOpenCollapse(!isOpenCollapse)}>
        {
          isOpenCollapse
            ? <i className='fas fa-angle-double-up hoverable'> Show less</i>
            : <i className='fas fa-angle-double-down hoverable'> Show more</i>
        }
      </Button>
      {
        displayAddNewDocumentModal &&
        <AddNewDocumentModal
          visible={displayAddNewDocumentModal}
          onCloseModal={() => setDisplayAddNewDocumentModal(false)}
          {...props}
        />
      }
    </div>
  )
})

function sortStageDoc (stage: TradeStageType): Array<TradeStageDocType> {
  let front = []
  let tail = []
  for (let d of stage.docs) {
    if (d.status === approveStatus.pending) {
      front.push(d)
    } else {
      tail.push(d)
    }
  }
  return front.concat(tail)
}

function computeStageStatus (stage: TradeStageType, stageIdx: number) {
  if (stage.stageDeleteStatus === canMakeReqStatus.approved) {
    return stageStatusMap.rejected
  }
  if (stage.stageCloseStatus === canMakeReqStatus.approved) {
    return stageStatusMap.approved
  }
  if (tradeViewStore.currentStageIdx === stageIdx) {
    return stageStatusMap.current
  }
  return stageStatusMap.pending
}

function isStageFinished (stage: TradeStageType): boolean {
  return stage.stageCloseStatus === canMakeReqStatus.approved ||
    stage.stageDeleteStatus === canMakeReqStatus.approved
}
