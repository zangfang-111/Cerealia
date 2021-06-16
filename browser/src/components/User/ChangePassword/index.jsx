// @flow

import React from 'react'
import { Form, Spin } from 'antd'
import { observer } from 'mobx-react-lite'
import useReactRouter from 'use-react-router'
import PasswordForm from './form'
import type { ChangePasswordType } from '../../../model/flowType'
import { addNotificationHelper } from '../../../lib/helper'
import spinnerStore from '../../../stores/spinner-store'
import usersStore from '../../../stores/user-store'

export default observer(() => {
  const { history } = useReactRouter()
  async function onChangePassword (values: ChangePasswordType) {
    spinnerStore.showSpinner('please wait...')
    try {
      await usersStore.changePassword(values)
      history.push('/home')
      addNotificationHelper('Password is changed successfully!', 'error')
    } catch (err) {
      addNotificationHelper(err, 'error')
    }
    spinnerStore.hideSpinner()
  }

  const NewPasswordForm = Form.create()(PasswordForm)
  return (
    <div className={'profile-page'}>
      <Spin spinning={spinnerStore.loading} size='large' tip={spinnerStore.spinnerTip}>
        <NewPasswordForm onChangePassword={onChangePassword} />
      </Spin>
    </div>
  )
})
