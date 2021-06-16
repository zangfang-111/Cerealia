// @flow

import React from 'react'
import useReactRouter from 'use-react-router'
import { observer } from 'mobx-react-lite'
import { Form, Spin } from 'antd'
import PasswordLoginForm from './form'
import { addNotificationHelper } from '../../../lib/helper'
import currentUser from '../../../stores/current-user'
import spinnerStore from '../../../stores/spinner-store'

export default observer(() => {
  const { history } = useReactRouter()
  async function onLogin (email: string, password: string) {
    let inputValue = {
      email: email,
      password: password
    }
    spinnerStore.showSpinner('please wait...')
    try {
      await currentUser.login(inputValue)
      history.push('/home')
    } catch (err) {
      addNotificationHelper(err, 'error')
    }
    spinnerStore.hideSpinner()
  }

  const LoginForm = Form.create()(PasswordLoginForm)
  return (
    <div className={'plain-container'}>
      <div className={'userForm-content'}>
        <p className={`circle ${status}`}>
          <i className='fas fa-user' />
        </p>
        <div className={'confirm'}>
          <Spin spinning={spinnerStore.loading} size='large' tip={spinnerStore.spinnerTip}>
            <LoginForm onLogin={onLogin} />
          </Spin>
        </div>
      </div>
    </div>
  )
})
