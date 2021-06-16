// @flow

import React from 'react'
import { addNotificationHelper, createTradeStagePath, hasPendingDoc } from '../../../../../lib/helper'
import { approveStatus } from '../../../../../constants/tradeConst'
import spinnerStore from '../../../../../stores/spinner-store'
import stellarStore from '../../../../../stores/stellarStore'
import tradeViewStore from '../../../../../stores/tradeViewStore'
import currentUser from '../../../../../stores/current-user'
import Button from '../../../Button/Button'
import PublicKeyConfirmation from '../../../Collapse/PublicKeyConfirmation/PublicKeyConfirmation'
import ReasonField from '../../../ReasonField'
import BaseModal from '../BaseModal'
import type { StageModalPropsType } from '../types'

type State = {
  error: boolean,
  isOpenCollapse: boolean,
  reason: string
}

class CloseStageModal extends React.Component<StageModalPropsType, State> {
  child: Object
  constructor (props: StageModalPropsType) {
    super(props)
    this.state = {
      error: true,
      isOpenCollapse: false,
      reason: ''
    }
  }

  confirmPublicKey = () => this.setState({ isOpenCollapse: true })

  onCancel = () => this.setState({ isOpenCollapse: false },
    () => this.props.onCloseModal())

  onChangeTextArea = (reasonText: string, err: boolean) => {
    this.setState({
      reason: reasonText,
      error: err
    })
  }

  closeStage = async () => {
    if (this.state.error) return
    const { onCloseModal, stageIdx } = this.props
    const inputValue = createTradeStagePath(
      tradeViewStore.id,
      stageIdx
    )
    spinnerStore.showSpinner()
    if (hasPendingDoc(tradeViewStore.stages[stageIdx])) {
      addNotificationHelper('This stage has a document with pending status', 'warning')
    }
    try {
      const signedTx = await stellarStore.signStageCloseTx(inputValue, approveStatus.pending)
      await tradeViewStore.stageCloseReq(inputValue, this.state.reason, signedTx)
      onCloseModal()
    } catch (err) {
      addNotificationHelper(err, 'error')
    }
    spinnerStore.hideSpinner()
  }

  render () {
    const { onCloseModal, stageIdx, visible } = this.props
    return (
      <BaseModal visible={visible} onCloseModal={onCloseModal} status={'close_stage'}
        spinnerStore={spinnerStore}>
        <div className={'close-stage'}>
          <p className={'modal-title'}>Close Stage</p>
          <p>
            Closing the
            <span className={'highlight-text'}>
            "{tradeViewStore.stages[stageIdx] &&
            tradeViewStore.stages[stageIdx].name}"
            </span> stage close request
          </p>
          <ReasonField onChangeTextArea={this.onChangeTextArea} />
          <PublicKeyConfirmation
            isOpenCollapse={this.state.isOpenCollapse}
            publicKey={currentUser.pubKey}
            onCloseModal={this.onCancel}
            confirm={this.closeStage}
            keyVerified={stellarStore.keyVerified}
          />
          {
            !this.state.isOpenCollapse &&
            <div className={'btn-group'}>
              <Button type={'primary'} text={'Yes, close'} onClick={this.confirmPublicKey} />
              <Button type={''} text={'cancel'} onClick={onCloseModal} />
            </div>
          }
        </div>
      </BaseModal>
    )
  }
}

export default CloseStageModal
