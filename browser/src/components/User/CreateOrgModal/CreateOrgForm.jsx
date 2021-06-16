// @flow

import React from 'react'
import { Input, Form } from 'antd'
import { RegPhone } from '../../../constants/tradeConst'
import Button from '../../Common/Button/Button'

const FormItem = Form.Item

type Props = {
  curOrgName: string,
  onCloseModal: Function,
  createOrganization: Function,
  form: Object
}

export default (props: Props) => {
  function handleSubmit (e: SyntheticEvent<HTMLButtonElement>) {
    props.form.validateFieldsAndScroll((err, values) => {
      if (!err) {
        props.createOrganization(values)
      }
    })
  }

  const { getFieldDecorator } = props.form
  const formItemLayout = {
    labelCol: { span: 7 },
    wrapperCol: { span: 17 }
  }

  return (
    <Form className={'Add-Org-Form'}>
      <div className={'content-title modal-title'}>Add New Company</div>
      <FormItem {...formItemLayout} label={'Name'}>
        {getFieldDecorator('name', {
          initialValue: props.curOrgName,
          rules: [{ required: true, message: 'Company name is required', whitespace: true }]
        })(
          <Input disabled />
        )}
      </FormItem>
      <FormItem {...formItemLayout} label={'Address'}>
        {getFieldDecorator('address', {
          rules: [{ required: true, message: 'Company address is required', whitespace: true }]
        })(
          <Input spellCheck={false} placeholder={'Address'} />
        )}
      </FormItem>
      <FormItem {...formItemLayout} label={'Email'}>
        {getFieldDecorator('email', {
          rules: [
            { type: 'email', message: 'This is not valid E-mail!' },
            { required: true, message: 'Company Email is required' }
          ]
        })(
          <Input placeholder={'Email'} spellCheck={false} />
        )}
      </FormItem>
      <Form.Item {...formItemLayout} label={'Phone Number'}>
        {getFieldDecorator('telephone', {
          rules: [
            { required: true, message: 'Company phone number is required' },
            {
              pattern: RegPhone,
              message: 'This is not valid phone number'
            }]
        })(
          <Input placeholder={'Phone Number'} />
        )}
      </Form.Item>
      <div className={'btn-group'}>
        <Button type={'primary'} text={'Add'} onClick={handleSubmit} />
        <Button type={'default'} text={'Cancel'} onClick={props.onCloseModal} />
      </div>
    </Form>
  )
}
