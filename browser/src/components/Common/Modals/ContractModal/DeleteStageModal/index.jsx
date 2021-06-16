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

class DeleteStageModal extends React.Component<StageModalPropsType, State> {
  child: Object
  constructor (props: StageModalPropsType) {
    super(props)
    this.state = {
      error: true,
      reason: ''
    }
  }

  deleteStage = async () => {
    if (this.state.error) return
    const { stageIdx, onCloseModal } = this.props
    const inputValue = createTradeStagePath(tradeViewStore.id, stageIdx)
    try {
      await tradeViewStore.stageDelReq(inputValue, this.state.reason)
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
        <div className={'delete-stage'}>
          <p className={'modal-title'}>Delete Stage</p>
          <p>
            Deleting the
            <span className={'highlight-text'}>
            "{tradeViewStore.stages[stageIdx] &&
            tradeViewStore.stages[stageIdx].name}"
            </span> stage add request
          </p>
          <ReasonField onChangeTextArea={this.onChangeTextArea} />
          <div className={'btn-group'}>
            <Button type={'primary'} text={'Yes, delete'} onClick={this.deleteStage} />
            <Button type={''} text={'cancel'} onClick={onCloseModal} />
          </div>
        </div>
      </BaseModal>
    )
  }
}

export default DeleteStageModal
