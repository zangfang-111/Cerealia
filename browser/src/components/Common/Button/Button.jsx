// @flow

import { Button, Tooltip } from 'antd'
import React from 'react'

type Props = {
  type: string,
  text: string,
  onClick?: Function,
  disabled?: boolean,
  tooltipText?: string
}

export default function (props: Props) {
  const { type, text, onClick, disabled, tooltipText } = props
  let button = <Button disabled={disabled} type={type} className={`${type}`} onClick={onClick}>
    {text}
  </Button>
  return (
    <div className={'button-container'}>
      {
        disabled
          ? <Tooltip title={tooltipText}>
            {button}
          </Tooltip>
          : button
      }
    </div>
  )
}
