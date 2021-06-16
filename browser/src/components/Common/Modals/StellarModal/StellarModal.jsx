// @flow

import React from 'react'
import { invalidKeyMessage, misMatchKeyMessage } from '../../../../constants/errors'
import { addNotificationHelper } from '../../../../lib/helper'
import stellarStore from '../../../../stores/stellarStore'
import currentUser from '../../../../stores/current-user'
import Button from '../../Button/Button'
import BaseModal from '../ContractModal/BaseModal'
import type { ModalPropsType } from '../ContractModal/types'

type State = {
  error: boolean,
  errMessage: string,
  secret: string
}

class StellarModalContainer extends React.Component<ModalPropsType, State> {
  child: Object
  constructor (props: ModalPropsType) {
    super(props)

    this.state = { error: false, secret: '', errMessage: '' }
  }

  onChangeStellarKey = (e: SyntheticInputEvent<HTMLInputElement>) => {
    this.setState({
      error: !stellarStore.validateStellarSecretKey(e.currentTarget.value),
      secret: e.currentTarget.value,
      errMessage: invalidKeyMessage
    })
  }

  onConfirmKey = () => {
    const { onCloseModal } = this.props
    if (this.state.error) {
      addNotificationHelper(invalidKeyMessage, 'error')
    } else {
      if (stellarStore.validateAndSetUserKey(this.state.secret, currentUser.user.pubKey)) {
        this.setState({ secret: '' })
        onCloseModal()
      } else {
        this.setState({ error: true, errMessage: misMatchKeyMessage })
        addNotificationHelper(misMatchKeyMessage, 'error')
      }
    }
  }

  onKeyPress = (e: SyntheticKeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      this.onConfirmKey()
    }
  }

  render () {
    const { visible } = this.props

    return (
      <BaseModal
        visible={visible}
        onCloseModal={this.props.onCloseModal}
        status={'stellar_key'}
      >
        <div className={'stellar-modal-container'}>
          <p className={'modal-title'}>Stellar Secret Key</p>
          <div className={'text-field'}>
            {
              currentUser.user.pubKey
                ? <p>
                  <b>PublicKey: </b>
                  <span className={'highlight-text'}>{currentUser.user.pubKey}</span>
                </p>
                : <p className={'invalid'}>You don’t have any public key assigned</p>
            }
          </div>
          {
            stellarStore.keyVerified
              ? <div>
                <p className={'success'}>Secret key stored <i className={'fas fa-check green-tick'} /></p>
              </div>
              : <div>
                <input type={'text'} className={this.state.error ? 'invalid' : 'valid'}
                  value={this.state.secret} onChange={this.onChangeStellarKey}
                  placeholder={'Your SecretKey:'} onKeyPress={this.onKeyPress} />
                {
                  this.state.error && <p className={'invalid'}>{this.state.errMessage}</p>
                }
                <Button type={'primary'} text={'save locally'} onClick={this.onConfirmKey} />
              </div>
          }
        </div>
      </BaseModal>
    )
  }
}

export default StellarModalContainer
