// @flow

import { Switch } from 'antd'
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
  isOpenCollapse: boolean,
  isApproving: boolean,
  reason: string,
  isApprove: boolean
}

class ConfirmClosingStageModal extends React.Component<StageModalPropsType, State> {
  child: Object
  constructor (props: StageModalPropsType) {
    super(props)

    this.state = {
      error: true,
      isOpenCollapse: false,
      isApproving: false,
      reason: '',
      isApprove: true
    }
  }

  confirmPublicKey = (value: string) => {
    this.setState({ isOpenCollapse: true }, () => {
      if (value === 'reject') {
        return true
      }
      this.setState({ isApproving: true })
    })
  }

  onChangeTextArea = (reasonText: string, err: boolean) => {
    this.setState({
      reason: reasonText,
      error: err
    })
  }

  toggleConfirmReject = () => {
    this.setState({ isApprove: !this.state.isApprove })
  }

  onCancel = () => this.setState({ isOpenCollapse: false },
    () => this.props.onCloseModal())

  approveClosingStage = async () => {
    const { stageIdx, onCloseModal } = this.props
    const inputValue = createTradeStagePath(tradeViewStore.id, stageIdx)
    spinnerStore.showSpinner()
    try {
      const signedTx = await stellarStore.signStageCloseTx(inputValue, approveStatus.approved)
      await tradeViewStore.stageCloseReqApprove(inputValue, signedTx)
      onCloseModal()
      addNotificationHelper('Stage closed successfully!', 'success')
    } catch (err) {
      addNotificationHelper(err, 'error')
    }
    spinnerStore.hideSpinner()
  }

  rejectClosingStage = async () => {
    if (this.state.error) return
    const { stageIdx, onCloseModal } = this.props
    const inputValue = createTradeStagePath(tradeViewStore.id, stageIdx)
    spinnerStore.showSpinner()
    try {
      const signedTx = await stellarStore.signStageCloseTx(inputValue, approveStatus.rejected)
      await tradeViewStore.stageCloseReqReject(inputValue, this.state.reason, signedTx)
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
        <div className={'confirm-closing-stage'}>
          {
            this.state.isApprove
              ? <p className={'modal-title'}>Approve Close Stage</p>
              : <p className={'modal-title'}>Reject Close Stage</p>
          }
          <div className={'switching-section'}>
            <Switch
              className={'show-switch'}
              checkedChildren='Reject'
              unCheckedChildren='Approve'
              onChange={this.toggleConfirmReject}
            />
          </div>
          <div className={'content-title modal-content-text'}>
            {
              this.state.isApprove ? 'Approve ' : 'Reject '
            }
            the completion of the stage <span className={'highlight-text'}>
            "{tradeViewStore.stages[stageIdx] && tradeViewStore.stages[stageIdx].name}"
            </span>
          </div>
          {
            !this.state.isApprove &&
            <ReasonField onChangeTextArea={this.onChangeTextArea} />
          }
          <PublicKeyConfirmation
            isOpenCollapse={this.state.isOpenCollapse}
            publicKey={currentUser.pubKey}
            onCloseModal={this.onCancel}
            keyVerified={stellarStore.keyVerified}
            confirm={
              this.state.isApproving ? this.approveClosingStage : this.rejectClosingStage
            }
          />
          {
            !this.state.isOpenCollapse &&
            <div className={'confirm-btn'}>
              {
                this.state.isApprove
                  ? <Button type={'primary'} text={'Approve'}
                    onClick={() => this.confirmPublicKey('approve')} />
                  : <Button type={''} text={'Reject'}
                    onClick={() => this.confirmPublicKey('reject')} />
              }
            </div>
          }
        </div>
      </BaseModal>
    )
  }
}

export default ConfirmClosingStageModal
