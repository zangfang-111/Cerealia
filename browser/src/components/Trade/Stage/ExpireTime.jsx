// @flow

import React, { useState } from 'react'
import { observer } from 'mobx-react-lite'
import { Button, DatePicker, Spin } from 'antd'
import TimezonePicker from 'react-timezone'
import * as moment from 'moment-timezone'
import { addNotificationHelper,
  createTradeStagePath } from '../../../lib/helper'
import { timeFormat } from '../../../constants/tradeConst'
import tradeViewStore from '../../../stores/tradeViewStore'
import spinnerStore from '../../../stores/spinner-store'

type Props = {
  stageIdx: number,
  onCloseDatePicker: Function
}

export default observer((props: Props) => {
  const [expiresAt, setExpiresAt] = useState(moment())
  const [timeZone, setTimeZone] = useState('')
  const [timeZoneError, setTimeZoneError] = useState(false)

  function selectDate (e: Function) {
    // e is moment function of antd DatePicker
    setExpiresAt(e)
    setTimeZoneError(moment(e) < moment())
  }

  function selectTimeZone (timeZone: string) {
    setTimeZone(timeZone)
    setExpiresAt(expiresAt.tz(timeZone))
  }

  async function submitExpireTime () {
    setTimeZoneError(expiresAt < moment())
    if (timeZoneError) {
      return
    }
    const et = expiresAt.tz(Intl.DateTimeFormat().resolvedOptions().timeZone)
    const inputValue = createTradeStagePath(tradeViewStore.id, props.stageIdx)
    spinnerStore.showSpinner()
    try {
      await tradeViewStore.setStageExpireTime(inputValue, moment.utc(et).format())
      props.onCloseDatePicker()
    } catch (e) {
      addNotificationHelper(e, 'error')
    }
    spinnerStore.hideSpinner()
  }

  return (
    <div className={'stage-expire-time'} >
      <Spin spinning={spinnerStore.loading} size={'large'} tip={'setting expire time...'}>
        <div className={'stage-datepicker'}>
          select expiring date and time:
          <DatePicker
            showTime
            format={timeFormat}
            onChange={selectDate}
            value={expiresAt}
          />
          check your timezone:
          <TimezonePicker
            value={timeZone}
            onChange={selectTimeZone}
            inputProps={{
              placeholder: 'Select Timezone...',
              name: 'timezone'
            }}
            className={'timezone-picker'}
          />
          {
            timeZoneError &&
            <p className={'error'}>Expire time can't be in past </p>
          }
        </div>
        <div className={'datepicker-btngroup'}>
          <Button type={'primary'} onClick={submitExpireTime}>Submit</Button>
          <Button type={'default'} onClick={props.onCloseDatePicker}>Cancel</Button>
        </div>
      </Spin>
    </div>
  )
})
