// @flow

import React, { useState } from 'react'
import { Collapse } from 'react-collapse'
import Button from '../../Button/Button'
import { keyVerifyMessage } from '../../../../constants/errors'

type Props = {
  isOpenCollapse: boolean,
  publicKey: string,
  onCloseModal: Function,
  confirm: Function,
  keyVerified: boolean
}

export default function (props: Props) {
  const [error, setError] = useState(false)
  const onConfirm = () => {
    if (!props.keyVerified) {
      setError(true)
    } else {
      setError(false)
      props.confirm()
    }
  }
  const { isOpenCollapse, publicKey, onCloseModal } = props
  return (
    <Collapse isOpened={isOpenCollapse} className={'public-key-confirmation'}>
      <div className={'collapse-title'}>Do you really want to sign using
        <span className={'highlight-text'}> {publicKey}</span>
      </div>
      {
        error && <p className={'error'}>{keyVerifyMessage}</p>
      }
      <div className={'btn-group'}>
        <Button type={'primary'} text={'Sign and Submit'} onClick={onConfirm} />
        <Button type={''} text={'Cancel'} onClick={onCloseModal} />
      </div>
    </Collapse>
  )
}
