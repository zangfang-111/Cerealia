// @flow

import React from 'react'
import { Form, Select } from 'antd'
import { observer } from 'mobx-react-lite'
import appStore from '../../../stores/app-store'

const Option = Select.Option

export default observer(() => {
  const theme = appStore.appTheme === 'theme-light' ? 'Light Theme' : 'Dark Theme'
  return (
    <Form className={'profile-page'} >
      <div className={'content-title'}>Preferences</div>
      <div className={'preference-row'}>
        <div className={'preference-label'}>App Theme</div>
        <Select
          className={'select-input'}
          size={'large'}
          dropdownClassName={'select-dropdown'}
          value={theme}
          onChange={(value) => appStore.setTheme(value)}>
          <Option key={1} value={'theme-light'}>{'Light Theme'}</Option>
          <Option key={2} value={'theme-dark'}>{'Dark Theme'}</Option>
        </Select>
      </div>
    </Form>
  )
})
