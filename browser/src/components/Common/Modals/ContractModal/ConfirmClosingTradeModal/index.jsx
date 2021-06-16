// @flow

import { Switch } from 'antd'
import React from 'react'
import { addNotificationHelper } from '../../../../../lib/helper'
import { approveStatus } from '../../../../../constants/tradeConst'
import spinnerStore from '../../../../../stores/spinner-store'
import stellarStore from '../../../../../stores/stellarStore'
import tradeViewStore from '../../../../../stores/tradeViewStore'
import currentUser from '../../../../../stores/current-user'
import Button from '../../../Button/Button'
import PublicKeyConfirmation from '../../../Collapse/PublicKeyConfirmation/PublicKeyConfirmation'
import ReasonField from '../../../ReasonField'
import BaseModal from '../BaseModal'
import type { ModalPropsType } from '../types'

type State = {
  error: boolean,
  isOpenCollapse: boolean,
  isApproving: boolean,
  reason: string,
  withApproval: boolean
}

class ConfirmClosingTradeModal extends React.Component<ModalPropsType, State> {
  child: Object
  constructor (props: ModalPropsType) {
    super(props)

    this.state = {
      error: true,
      isOpenCollapse: false,
      isApproving: false,
      reason: '',
      withApproval: true
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

  onCancel = () => this.setState({ isOpenCollapse: false },
    () => this.props.onCloseModal())

  approveClosingTrade = async () => {
    const { onCloseModal } = this.props
    spinnerStore.showSpinner()
    try {
      const signedTx = await stellarStore.signTradeCloseTx(
        tradeViewStore.id,
        approveStatus.approved)
      await tradeViewStore.tradeCloseReqApprove(signedTx)
      onCloseModal()
    } catch (err) {
      addNotificationHelper(err, 'error')
    }
    spinnerStore.hideSpinner()
  }

  toggleConfirmReject = () => {
    this.setState({ withApproval: !this.state.withApproval })
  }

  rejectClosingTrade = async () => {
    if (this.state.error) return
    const { onCloseModal } = this.props
    spinnerStore.showSpinner()
    try {
      const signedTx = await stellarStore.signTradeCloseTx(
        tradeViewStore.id,
        approveStatus.rejected)
      await tradeViewStore.tradeCloseReqReject(this.state.reason, signedTx)
      onCloseModal()
    } catch (err) {
      addNotificationHelper(err, 'error')
    }
    spinnerStore.hideSpinner()
  }

  render () {
    const { onCloseModal, visible } = this.props
    return (
      <BaseModal visible={visible} onCloseModal={onCloseModal} status={'close_stage'}
        spinnerStore={spinnerStore}>
        <div className={'confirm-closing-stage'}>
          {
            this.state.withApproval
              ? <p className={'modal-title'}>Approve Close Trade</p>
              : <p className={'modal-title'}>Reject Close Trade</p>
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
              this.state.withApproval ? 'Approve ' : 'Reject '
            }
            the completion of the trade <span className={'highlight-text'}>
            "{tradeViewStore.name}"
            </span>
          </div>
          {
            !this.state.withApproval &&
            <ReasonField onChangeTextArea={this.onChangeTextArea} />
          }
          <PublicKeyConfirmation
            isOpenCollapse={this.state.isOpenCollapse}
            publicKey={currentUser.pubKey}
            onCloseModal={this.onCancel}
            keyVerified={stellarStore.keyVerified}
            confirm={
              this.state.isApproving ? this.approveClosingTrade : this.rejectClosingTrade
            }
          />
          {
            !this.state.isOpenCollapse &&
            <div className={'confirm-btn'}>
              {
                this.state.withApproval
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

export default ConfirmClosingTradeModal
