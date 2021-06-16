// @flow

import { Icon, Radio, Switch } from 'antd'
import React from 'react'
import { addNotificationHelper, createTradeStagePath } from '../../../../../lib/helper'
import {
  approveStatus,
  buyerActor,
  sellerActor,
  tradeActorMap
} from '../../../../../constants/tradeConst'
import type { CreateTradeStageInputType } from '../../../../../model/flowType'
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
  selectedUserType: string,
  reason: string,
  confirm: boolean,
  isOpenCollapse: boolean
}

class AddNewStageModal extends React.Component<ModalPropsType, State> {
  stageName: string
  description: string
  child: Object
  constructor (props: ModalPropsType) {
    super(props)

    this.stageName = ''
    this.description = ''

    this.state = {
      error: true,
      selectedUserType: tradeViewStore.currentRole === tradeActorMap.b
        ? buyerActor : sellerActor,
      reason: '',
      isOpenCollapse: false,
      confirm: false
    }
  }

  confirmPublicKey = () => this.setState({ isOpenCollapse: true })

  onCancel = () => this.setState({ isOpenCollapse: false }, () => this.props.onCloseModal())

  onCreateNewStage = (): void => {
    if (this.state.error) return
    const inputValue = {
      tid: tradeViewStore.id,
      name: this.stageName,
      description: this.description,
      reason: this.state.reason,
      owner: this.state.selectedUserType
    }
    this.createNewStage(inputValue)
  }

  createNewStage = async (value: CreateTradeStageInputType) => {
    const { onCloseModal } = this.props
    const tradeStagePath = createTradeStagePath(value.tid, tradeViewStore.stageAddReqs.length)
    const operation = this.state.confirm ? approveStatus.pending : approveStatus.approved
    spinnerStore.showSpinner()
    try {
      const signedTx = await stellarStore.signStageAddTx(tradeStagePath, operation)
      await tradeViewStore.createStage(value, signedTx, this.state.confirm)
      this.stageName = ''
      this.description = ''
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

  selectUserType = (e: SyntheticInputEvent<HTMLInputElement>) => {
    // currentType dosesn't exists for a group element
    this.setState({ selectedUserType: e.target.value })
  }

  onToggleConfirmation = () => {
    this.setState({ confirm: !this.state.confirm })
  }

  render () {
    const { visible, onCloseModal } = this.props
    const RadioGroup = Radio.Group
    return (
      <BaseModal visible={visible} onCloseModal={onCloseModal}
        status={'new_stage'} spinnerStore={spinnerStore}>
        <div className={`add-new-stage ${this.state.error ? `error` : ''}`}>
          <p className={'modal-title'}>Add New Stage</p>
          <p>
            Select an owner and name of the stage
          </p>
          <RadioGroup
            defaultValue={this.state.selectedUserType}
            onChange={this.selectUserType}
          >
            <Radio value={buyerActor}>Buyer</Radio>
            <Radio value={sellerActor}>Seller</Radio>
          </RadioGroup>
          <input type={'text'} size={'large'} placeholder={'stage name'} ref={input => (this.stageName = input ? input.value : '')} />
          <textarea className={'description'} placeholder={'Type some description...'} ref={input => (this.description = input ? input.value : '')} />
          <ReasonField onChangeTextArea={this.onChangeTextArea} />
          <div className={'confirm-approve'}>
            <Switch className={'switch-btn'} checkedChildren={<Icon type={'check'} />}
              unCheckedChildren={<Icon type={'close'} />} onChange={this.onToggleConfirmation} />
            <span>Counterparty approve required</span>
          </div>
          <PublicKeyConfirmation
            isOpenCollapse={this.state.isOpenCollapse}
            publicKey={currentUser.pubKey}
            onCloseModal={this.onCancel}
            confirm={this.onCreateNewStage}
            keyVerified={stellarStore.keyVerified}
          />
          {
            !this.state.isOpenCollapse &&
            <div className={'btn-group'}>
              <Button type={'primary'} text={'Add'} onClick={this.confirmPublicKey} />
            </div>
          }
        </div>
      </BaseModal>
    )
  }
}

export default AddNewStageModal
