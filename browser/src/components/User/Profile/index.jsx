// @flow

import React, { useEffect } from 'react'
import { Form, Spin } from 'antd'
import { observer } from 'mobx-react-lite'
import useReactRouter from 'use-react-router'
import axios from 'axios'
import type { UserProfileInputType } from '../../../model/flowType'
import { addNotificationHelper, mkFormDataHeaders } from '../../../lib/helper'
import UserProfileForm from './form'
import { mkLink } from '../../../services/cerealia'
import spinnerStore from '../../../stores/spinner-store'
import usersStore from '../../../stores/user-store'
import currentUser from '../../../stores/current-user'

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

  async function onSave (input: UserProfileInputType, file: Object) {
    spinnerStore.showSpinner('please wait...')
    try {
      await uploadAvatar(file)
      await currentUser.updateUserProfile(input)
      history.push('/home')
      addNotificationHelper('User profile is updated successfully!', 'success')
    } catch (err) {
      addNotificationHelper(err, 'error')
    }
    spinnerStore.hideSpinner()
  }

  async function uploadAvatar (file: Object) {
    if (!file.name) {
      return
    }
    let data = new FormData()
    data.append('avatarFile', file)
    try {
      await axios.post(
        mkLink('v1/users/avatar'),
        data,
        {
          headers: mkFormDataHeaders()
        })
    } catch (err) {
      throw err
    }
  }

  const ProfileForm = Form.create()(UserProfileForm)
  return (
    <div className={'profile-page'}>
      <Spin spinning={spinnerStore.loading} size='large' tip={spinnerStore.spinnerTip}>
        <ProfileForm onSave={onSave} />
      </Spin>
    </div>
  )
})
