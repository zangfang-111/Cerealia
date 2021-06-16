// @flow

import React, { useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import useReactRouter from 'use-react-router'
import { Form, Spin } from 'antd'
import UserSignupForm from './form'
import type { NewUserInputType } from '../../../model/flowType'
import { addNotificationHelper } from '../../../lib/helper'
import spinnerStore from '../../../stores/spinner-store'
import usersStore from '../../../stores/user-store'

async function sync () {
  try {
    if (usersStore.organizations.length === 0) {
      await usersStore.fetchOrganizations()
    }
  } catch (err) {
    addNotificationHelper(err, 'err')
  }
}

export default observer(() => {
  const { history } = useReactRouter()
  useEffect(() => { sync() })

  async function onSignup (input: NewUserInputType) {
    spinnerStore.showSpinner('please wait...')
    try {
      await usersStore.signup(input)
      history.push('/login')
      addNotificationHelper('You have registered successfully', 'success')
    } catch (err) {
      addNotificationHelper(err, 'error')
    }
    spinnerStore.hideSpinner()
  }

  const SignupForm = Form.create()(UserSignupForm)
  return (
    <div className={'plain-container'}>
      <div className={'userForm-content'}>
        <p className={`circle ${status}`}>
          <i className='fas fa-user' />
        </p>
        <p className={'close-icon'}>
          <i className='fas fa-times-circle hoverable' onClick={history.goBack} />
        </p>
        <div className={'confirm'}>
          <Spin spinning={spinnerStore.loading} size='large' tip={spinnerStore.spinnerTip}>
            <SignupForm onSignup={onSignup} />
          </Spin>
        </div>
      </div>
    </div>
  )
})
