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
import ReasonField from '../../../ReasonField'
import BaseModal from '../BaseModal'
import type { StageModalPropsType } from '../types'

type Props = StageModalPropsType & {
  stageDoc: TradeStageDocType
}

type State = {
  error: boolean,
  isOpenCollapse: boolean,
  reason: string,
  loading: boolean
}

class RejectTradeStageDocModal extends React.Component<Props, State> {
  child: Object
  constructor (props: Props) {
    super(props)

    this.state = {
      error: true,
      isOpenCollapse: false,
      reason: '',
      loading: false
    }
  }

  confirmPublicKey = () => this.setState({ isOpenCollapse: true })

  onCancel = () => this.setState({ isOpenCollapse: false },
    () => this.props.onCloseModal())

  rejectTradeStageDoc = async () => {
    const { onCloseModal, stageIdx, stageDoc } = this.props
    if (this.state.error) return
    let inputValue: TradeStageDocPathType = {
      tid: tradeViewStore.id,
      stageIdx: stageIdx,
      stageDocIdx: stageDoc.index,
      stageDocHash: stageDoc.doc.hash
    }
    this.setState({ loading: true })
    try {
      spinnerStore.showSpinner()
      const signedTx = await stellarStore.signDocApprovalTx(inputValue, approveStatus.rejected)
      await tradeViewStore.stageDocReject(inputValue, this.state.reason,
        currentUser.user, signedTx)
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
    const { onCloseModal, visible, stageDoc } = this.props
    return (
      <BaseModal visible={visible} onCloseModal={onCloseModal} status={'reject_additional_doc'}
        spinnerStore={spinnerStore}>
        <div className={'reject-additional-doc'}>
          <p className={'modal-title'}>Document Rejection</p>
          <div className={'text-field'}>
            <p>
              Do you want to reject <br />
              <span className={'highlight-text'}>
                "{stageDoc.doc && stageDoc.doc.name}"
              </span> ?
            </p>
            <p>
              Document Note: "{stageDoc.doc && stageDoc.doc.note}"
            </p>
          </div>
          <ReasonField onChangeTextArea={this.onChangeTextArea} />
          <PublicKeyConfirmation
            isOpenCollapse={this.state.isOpenCollapse}
            publicKey={currentUser.user.pubKey}
            onCloseModal={this.onCancel}
            confirm={this.rejectTradeStageDoc}
            keyVerified={stellarStore.keyVerified}
          />
          {
            !this.state.isOpenCollapse &&
            <p className={'btn-group'}>
              <Button type={'primary'} text={'Yes, reject'} onClick={this.confirmPublicKey} />
              <Button type={''} text={'cancel'} onClick={onCloseModal} />
            </p>
          }
        </div>
      </BaseModal>
    )
  }
}
export default RejectTradeStageDocModal
