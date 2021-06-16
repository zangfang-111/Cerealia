// @flow

import React from 'react'
import { addNotificationHelper, createTradeStagePath } from '../../../../../lib/helper'
import tradeViewStore from '../../../../../stores/tradeViewStore'
import Button from '../../../Button/Button'
import ReasonField from '../../../ReasonField'
import BaseModal from '../BaseModal'
import type { StageModalPropsType } from '../types'

type State = {
  error: boolean,
  reason: string
}

class ConfirmDeletingStageModal extends React.Component<StageModalPropsType, State> {
  child: Object
  constructor (props: StageModalPropsType) {
    super(props)

    this.state = {
      error: true,
      reason: ''
    }
  }

  approveDeletingStage = async () => {
    const { stageIdx, onCloseModal } = this.props
    const inputValue = createTradeStagePath(tradeViewStore.id, stageIdx)
    try {
      await tradeViewStore.stageDelReqApprove(inputValue)
      onCloseModal()
      addNotificationHelper('Stage deleted successfully!', 'success')
    } catch (err) {
      addNotificationHelper(err, 'error')
    }
  }

  rejectDeletingStage = async () => {
    if (this.state.error) return
    const { stageIdx, onCloseModal } = this.props
    const inputValue = createTradeStagePath(tradeViewStore.id, stageIdx)
    try {
      await tradeViewStore.stageDelReqReject(inputValue, this.state.reason)
      onCloseModal()
    } catch (err) {
      addNotificationHelper(err, 'error')
    }
  }

  onChangeTextArea = (reasonText: string, err: boolean) => {
    this.setState({
      reason: reasonText,
      error: err
    })
  }

  render () {
    const { onCloseModal, stageIdx, visible } = this.props
    return (
      <BaseModal visible={visible} onCloseModal={onCloseModal} status={'delete_stage'}>
        <div className={'confirm-deleting-stage'}>
          <p className={'modal-title'}>Confirm Delete Stage</p>
          <p>
            Deleting the
            <span className={'highlight-text'}>
            "{tradeViewStore.stages[stageIdx] &&
            tradeViewStore.stages[stageIdx].name}"
            </span> stage add request
          </p>
          <ReasonField onChangeTextArea={this.onChangeTextArea} />
          <div className={'btn-group'}>
            <Button type={'primary'} text={'Confirm'} onClick={this.approveDeletingStage} />
            <Button type={''} text={'Reject'} onClick={this.rejectDeletingStage} />
          </div>
        </div>
      </BaseModal>
    )
  }
}

export default ConfirmDeletingStageModal
