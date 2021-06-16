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
import BaseModal from '../BaseModal'
import type { StageModalPropsType } from '../types'

type State = {
  isOpenCollapse: boolean
}

class ApproveTradeStageReqModal extends React.Component<StageModalPropsType, State> {
  child: Object
  constructor (props: StageModalPropsType) {
    super(props)

    this.state = {
      isOpenCollapse: false
    }
  }

  confirmPublicKey = () => this.setState({ isOpenCollapse: true })

  onCancel = () => this.setState({ isOpenCollapse: false },
    () => this.props.onCloseModal())

  approveTradeStageReq = async () => {
    const { onCloseModal, stageIdx } = this.props
    const value = createTradeStagePath(tradeViewStore.id, stageIdx)
    spinnerStore.showSpinner()
    try {
      const signedTx = await stellarStore.signStageAddTx(value, approveStatus.approved)
      await tradeViewStore.addStageApprove(value, signedTx)
      onCloseModal()
    } catch (err) {
      console.error('error', err)
      addNotificationHelper(err, 'error')
    }
    spinnerStore.hideSpinner()
  }

  render () {
    const { onCloseModal, stageIdx, visible } = this.props

    return (
      <BaseModal visible={visible} onCloseModal={onCloseModal} spinnerStore={spinnerStore}
        status={'approve_additional_stage'} >
        <div className={'approve-additional-stage'}>
          <p className={'modal-title'}>Approve New Stage</p>
          <p>
            Approving the
            <span className={'highlight-text'}>
            " {tradeViewStore.stageAddReqs[stageIdx] &&
            tradeViewStore.stageAddReqs[stageIdx].name}"
            </span> stage add request
          </p>
          <PublicKeyConfirmation
            isOpenCollapse={this.state.isOpenCollapse}
            publicKey={currentUser.pubKey}
            onCloseModal={this.onCancel}
            confirm={this.approveTradeStageReq}
            keyVerified={stellarStore.keyVerified}
          />
          {
            !this.state.isOpenCollapse &&
            <div className={'btn-group'}>
              <Button type={'primary'} text={'Yes, approve'} onClick={this.confirmPublicKey} />
              <Button type={''} text={'cancel'} onClick={onCloseModal} />
            </div>
          }
        </div>
      </BaseModal>
    )
  }
}

export default ApproveTradeStageReqModal
