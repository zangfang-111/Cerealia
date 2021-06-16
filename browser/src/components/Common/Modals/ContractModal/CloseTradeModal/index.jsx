// @flow

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
  reason: string
}

class CloseTradeModal extends React.Component<ModalPropsType, State> {
  child: Object
  constructor (props: ModalPropsType) {
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

  closeTrade = async () => {
    if (this.state.error) return
    const { onCloseModal } = this.props
    spinnerStore.showSpinner()
    try {
      const signedTx = await stellarStore.signTradeCloseTx(
        tradeViewStore.id,
        approveStatus.pending)
      await tradeViewStore.tradeCloseReq(this.state.reason, signedTx)
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
        <div className={'close-stage'}>
          <p className={'modal-title'}>Close Trade</p>
          <p>
            Closing the trade <span className={'highlight-text'}>
            "{tradeViewStore.name}" </span>
          </p>
          <ReasonField onChangeTextArea={this.onChangeTextArea} />
          <PublicKeyConfirmation
            isOpenCollapse={this.state.isOpenCollapse}
            publicKey={currentUser.pubKey}
            onCloseModal={this.onCancel}
            confirm={this.closeTrade}
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

export default CloseTradeModal
