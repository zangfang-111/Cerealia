// @flow

import React from 'react'
import { addNotificationHelper, createTradeStagePath } from '../../../../../lib/helper'
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
  reason: string,
  isOpenCollapse: boolean
}

class RejectTradeStageReqModal extends React.Component<StageModalPropsType, State> {
  child: Object
  constructor (props: StageModalPropsType) {
    super(props)

    this.state = {
      error: true,
      reason: '',
      isOpenCollapse: false
    }
  }

  confirmPublicKey = () => this.setState({ isOpenCollapse: true })

  onCancel = () => this.setState({ isOpenCollapse: false },
    () => this.props.onCloseModal())

  rejectTradeStageReq = async () => {
    const { onCloseModal, stageIdx } = this.props
    if (this.state.error) return
    const inputValue = createTradeStagePath(tradeViewStore.id, stageIdx)
    spinnerStore.showSpinner()
    try {
      const signedTx = await stellarStore.signStageAddTx(inputValue, approveStatus.rejected)
      await tradeViewStore.addStageReject(inputValue, signedTx, this.state.reason)
      onCloseModal()
    } catch (err) {
      addNotificationHelper(err, 'error')
    }
    spinnerStore.hideSpinner()
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
      <BaseModal visible={visible} onCloseModal={onCloseModal} status={'reject_additional_stage'}
        spinnerStore={spinnerStore}>
        <div className={'reject-additional-stage'}>
          <p className={'modal-title'}>Reject New Stage</p>
          <p>
            Rejecting the
            <span className={'highlight-text'}>
              "{tradeViewStore.stageAddReqs[stageIdx] &&
            tradeViewStore.stageAddReqs[stageIdx].name}"
            </span> stage add request
          </p>
          <ReasonField onChangeTextArea={this.onChangeTextArea} />
          <PublicKeyConfirmation
            isOpenCollapse={this.state.isOpenCollapse}
            publicKey={currentUser.pubKey}
            onCloseModal={this.onCancel}
            confirm={this.rejectTradeStageReq}
            keyVerified={stellarStore.keyVerified}
          />
          {
            !this.state.isOpenCollapse &&
            <div className={'btn-group'}>
              <Button type={'primary'} text={'Yes, reject'} onClick={this.confirmPublicKey} />
              <Button type={''} text={'cancel'} onClick={onCloseModal} />
            </div>
          }
        </div>
      </BaseModal>
    )
  }
}

export default RejectTradeStageReqModal
