// @flow

import React from 'react'
import useReactRouter from 'use-react-router'
import { observer } from 'mobx-react-lite'
import { Form, Spin } from 'antd'
import PasswordForm from './form'
import { addNotificationHelper } from '../../../lib/helper'
import spinnerStore from '../../../stores/spinner-store'
import currentUser from '../../../stores/current-user'

export default observer(() => {
  const { history } = useReactRouter()
  async function onSaveEmail (values: Array<string>) {
    if (!values || values.length === 0) {
      addNotificationHelper('A user should have at least one email', 'error')
      return
    }
    spinnerStore.showSpinner('please wait...')
    try {
      await currentUser.changeEmail(values)
      history.push('/home')
      addNotificationHelper('Email is changed successfully!', 'error')
    } catch (err) {
      addNotificationHelper(err, 'error')
    }
    spinnerStore.hideSpinner()
  }

  const NewPasswordForm = Form.create()(PasswordForm)
  return (
    <div className={'profile-page'}>
      <Spin spinning={spinnerStore.loading} size='large' tip={spinnerStore.spinnerTip}>
        <NewPasswordForm onSaveEmail={onSaveEmail} />
      </Spin>
    </div>
  )
})
