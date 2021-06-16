// @flow

import { Radio, Select } from 'antd'
import { observer } from 'mobx-react'
import React from 'react'
import { addNotificationHelper } from '../../../../../lib/helper'
import { buyerActor, sellerActor } from '../../../../../constants/tradeConst'
import type { TradeTemplateType, UserType } from '../../../../../model/flowType'
import _ from '../../../../../polyfills/underscore.js'
import spinnerStore from '../../../../../stores/spinner-store'
import stellarStore from '../../../../../stores/stellarStore'
import tradesStore from '../../../../../stores/trade-store'
import curTradeOfferStore from '../../../../../stores/current-tradeoffer'
import usersStore from '../../../../../stores/user-store'
import currentUser from '../../../../../stores/current-user'
import Button from '../../../Button/Button'
import PublicKeyConfirmation from '../../../Collapse/PublicKeyConfirmation/PublicKeyConfirmation'
import BaseModal from '../BaseModal'
import type { ModalPropsType } from '../types'

type State = {
  titleError: boolean,
  selectedUserType: string,
  tradeTemplateID: string,
  partnerID: string,
  isOpenCollapse: boolean
}

type Props = ModalPropsType & {
  withTradeOffer?: boolean
}

@observer
class CreateNewTradeModal extends React.Component<Props, State> {
  newTradeName: string
  newTradeDescription: string
  child: Object
  constructor (props: Props) {
    super(props)

    this.newTradeName = ''
    this.newTradeDescription = ''

    this.state = {
      titleError: false,
      selectedUserType: '',
      tradeTemplateID: '',
      partnerID: '',
      isOpenCollapse: false
    }
  }

  componentDidMount = async () => {
    try {
      if (tradesStore.tradeTemplates.length === 0) {
        await tradesStore.fetchTradeTemplates()
      }
    } catch (err) {
      addNotificationHelper(err, 'error')
    }
  }

  confirmPublicKey = () => this.setState({ isOpenCollapse: true })

  onCancel = () => this.setState({ isOpenCollapse: false },
    () => this.props.onCloseModal())

  createNewTrade = async () => {
    const { onCloseModal, withTradeOffer } = this.props
    const inputValue = {
      name: this.newTradeName,
      description: this.newTradeDescription,
      templateID: this.state.tradeTemplateID,
      sellerID: this.state.selectedUserType === sellerActor
        ? currentUser.user.id
        : this.state.partnerID,
      buyerID: this.state.selectedUserType === buyerActor
        ? currentUser.user.id
        : this.state.partnerID,
      tradeOfferID: withTradeOffer ? curTradeOfferStore.tradeOffer.id : ''
    }
    spinnerStore.showSpinner()
    try {
      await tradesStore.createTrade(inputValue)
      this.initialize()
      onCloseModal()
    } catch (err) {
      addNotificationHelper(err, 'error')
    }
    spinnerStore.hideSpinner()
  }

  initialize = () => {
    this.newTradeName = ''
    this.newTradeDescription = ''
  }

  onCreateNewTrade = () => {
    this.setState({
      titleError: false
    })
    if (!this.newTradeName) {
      this.setState({ titleError: true })
    }

    if (this.state.titleError) {
      return
    }
    this.createNewTrade()
  }

  selectUserType = (e: SyntheticInputEvent<HTMLInputElement>) => {
    // currentType dosesn't exists for a group element
    this.setState({ selectedUserType: e.target.value })
  }

  onChangeTemplateId = (value: string) => {
    this.setState({ tradeTemplateID: value })
  }
  onChangeUserId = (value: string) => {
    this.setState({ partnerID: value })
  }
  render () {
    const { visible, onCloseModal } = this.props

    const Option = Select.Option
    const RadioGroup = Radio.Group
    return (
      <BaseModal visible={visible} onCloseModal={onCloseModal}
        status={'create_trade'} spinnerStore={spinnerStore}>
        <div className={'create-new-trade'}>
          <p className={'modal-title'}>Fill initial details</p>
          <p>
            Select name and template of a trade.
          </p>
          <div className={'input-container'}>
            <RadioGroup defaultValue={this.state.selectedUserType}
              onChange={this.selectUserType}>
              <Radio value={buyerActor}>I am a Buyer</Radio>
              <Radio value={sellerActor}>I am a Seller</Radio>
            </RadioGroup>
            <input
              type={'text'}
              className={`${this.state.titleError ? 'error' : ''}`}
              placeholder={'trade name'}
              size='large'
              ref={input => (this.newTradeName = input ? input.value : '')} />
            {
              this.state.titleError && <p className={'error'}>Please input the new trade name correctly</p>
            }
            <div className={'input-selector-group'}>
              <Select size={'large'} placeholder={'Template'} onChange={this.onChangeTemplateId}>
                {
                  !_.isEmpty(tradesStore.tradeTemplates) &&
                  tradesStore.tradeTemplates.map((template: TradeTemplateType, i) => {
                    return (
                      <Option value={template.id} key={i}>{template.name}</Option>
                    )
                  })
                }
              </Select>
              <Select size={'large'} placeholder={'Counterparty'} onChange={this.onChangeUserId}>
                {
                  !_.isEmpty(usersStore.users) &&
                  usersStore.users.filter(item => item.id !== currentUser.user.id &&
                    item.roles.indexOf('trader') !== -1)
                    .map((user: UserType, i) => {
                      return (
                        <Option
                          value={user.id}
                          key={i}
                        >
                          {`${user.orgMap[0].org.name} - ${user.firstName} ${user.lastName}`}
                        </Option>
                      )
                    })
                }
              </Select>
            </div>
            <div className={'description'}>
              <p>
                Please provide more details
              </p>
              <textarea
                placeholder={'Type something...'}
                ref={input => (this.newTradeDescription = input ? input.value : '')} />
            </div>
            <PublicKeyConfirmation
              isOpenCollapse={this.state.isOpenCollapse}
              publicKey={currentUser.pubKey}
              onCloseModal={this.onCancel}
              confirm={this.onCreateNewTrade}
              keyVerified={stellarStore.keyVerified}
            />
            { !this.state.isOpenCollapse && <Button type={'primary'} text={'confirm'} onClick={this.confirmPublicKey} /> }
          </div>
        </div>
      </BaseModal>
    )
  }
}

export default CreateNewTradeModal
