// @flow

import React from 'react'
import { addNotificationHelper } from '../../../../../lib/helper'
import { approveStatus } from '../../../../../constants/tradeConst'
import type { TradeStageDocPathType, TradeStageDocType } from '../../../../../model/flowType'
import spinnerStore from '../../../../../stores/spinner-store'
import stellarStore from '../../../../../stores/stellarStore'
import tradeViewStore from '../../../../../stores/tradeViewStore'
import currentUser from '../../../../../stores/current-user'
import Button from '../../../Button/Button'
import PublicKeyConfirmation from '../../../Collapse/PublicKeyConfirmation/PublicKeyConfirmation'
import BaseModal from '../BaseModal'
import type { StageModalPropsType } from '../types'

type Props = StageModalPropsType & {
  stageDoc: TradeStageDocType
}

type State = {
  error: boolean,
  isOpenCollapse: boolean
}

class ApproveTradeStageDocModal extends React.Component<Props, State> {
  child: Object
  constructor (props: Props) {
    super(props)

    this.state = {
      error: false,
      isOpenCollapse: false
    }
  }

  confirmPublicKey = () => this.setState({ isOpenCollapse: true })

  onCancel = () => this.setState({ isOpenCollapse: false },
    () => this.props.onCloseModal())

  approveTradeStageDoc = async () => {
    const { onCloseModal, stageIdx, stageDoc } = this.props
    let inputValue: TradeStageDocPathType = {
      tid: tradeViewStore.id,
      stageIdx: stageIdx,
      stageDocIdx: stageDoc.index,
      stageDocHash: stageDoc.doc.hash
    }
    spinnerStore.showSpinner()
    try {
      const signedTx = await stellarStore.signDocApprovalTx(inputValue, approveStatus.approved)
      await tradeViewStore.stageDocApprove(inputValue, signedTx)
      onCloseModal()
    } catch (err) {
      addNotificationHelper(err, 'error')
    }
    spinnerStore.hideSpinner()
  }

  render () {
    const { onCloseModal, stageDoc, visible } = this.props
    return (
      <BaseModal visible={visible} onCloseModal={onCloseModal} status={'approve_additional_doc'}
        spinnerStore={spinnerStore}>
        <div className={'approve-additional-doc'}>
          <p className={'modal-title'}>Document approval</p>
          <div className={'text-field'}>
            <p>
              Do you want to approve <br />
              <span className={'highlight-text'}>
                  " {stageDoc.doc && stageDoc.doc.name}"
              </span> ?
            </p>
            <p>
              Document Note: "{stageDoc.doc && stageDoc.doc.note}"
            </p>
          </div>
          <PublicKeyConfirmation
            isOpenCollapse={this.state.isOpenCollapse}
            publicKey={currentUser.pubKey}
            onCloseModal={this.onCancel}
            confirm={this.approveTradeStageDoc}
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

export default ApproveTradeStageDocModal
