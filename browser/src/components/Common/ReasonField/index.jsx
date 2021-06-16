// @flow

import React, { useState } from 'react'
import { Input } from 'antd'
import { invalidReasonMessage } from '../../../constants/errors'

type Props = {
  onChangeTextArea: Function
}

export default function (props: Props) {
  const [error, setError] = useState(false)

  const handleChangeTextArea = (e: SyntheticInputEvent<HTMLInputElement>) => {
    setError(e.currentTarget.value.length < 10)
    props.onChangeTextArea(e.currentTarget.value, e.currentTarget.value.length < 10)
  }

  return (
    <div >
      <Input.TextArea placeholder={'Type the reason here...'}
        rows={3} onChange={handleChangeTextArea} />
      {
        error && <p className={'error'}>{invalidReasonMessage}</p>
      }
    </div>
  )
}
