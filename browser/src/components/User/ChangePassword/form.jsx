// @flow

import React from 'react'
import { Form, Input } from 'antd'
import Button from '../../Common/Button/Button'
import { withRouter } from 'react-router-dom'
const pwdPattern = new RegExp(`^(?=.*[a-z])(?=.*[A-Z])(?=.*[!@#$%^&*,.;'"()])`)

const FormItem = Form.Item

type Props = {
  form: Object,
  onChangePassword: Function,
  history: Object
}

type State = {
  confirmDirty: boolean,
  displayPassErrMessage: boolean,
}

type ruleType = {
  validator: Function,
  field: string,
  fullField: string,
  type: string
}

@withRouter
class ProfileForm extends React.Component<Props, State> {
  constructor (props: Props) {
    super(props)
    this.state = {
      confirmDirty: false,
      displayPassErrMessage: false
    }
  }

  handleSubmit = (e: SyntheticEvent<HTMLButtonElement>) => {
    this.props.form.validateFieldsAndScroll((err, values) => {
      if (!err) {
        delete values.confirm
        this.props.onChangePassword(values)
      }
    })
  }

  handleConfirmBlur = (e: SyntheticEvent<HTMLInputElement>) => {
    this.setState({ confirmDirty: this.state.confirmDirty || !!e.currentTarget.value })
  }

  compareToPrevPassword = (rule: ruleType, value: string, callback: Function) => {
    const form = this.props.form
    if (value && value !== form.getFieldValue('newPassword')) {
      this.setState({ displayPassErrMessage: true })
    } else {
      this.setState({ displayPassErrMessage: false })
      callback()
    }
  }

  validateToNextPassword = (rule: ruleType, value: string, callback: Function) => {
    const form = this.props.form
    if (value && this.state.confirmDirty) {
      form.validateFields(['confirm'], { force: true })
    }
    callback()
  }

  render () {
    const { getFieldDecorator } = this.props.form
    return (
      <Form className={'profile-form'} >
        <div className={'content-title'}>Change Your Password</div>
        <br />
        <FormItem>
          <p className={'profile-label'}>Old Password
            <a className={'profile-edit'} >
              <strong>Forgot your password</strong>
            </a>
          </p>
          {getFieldDecorator('oldPassword', {
            rules: [{ required: true, message: 'Please input your old password!', whitespace: true }]
          })(
            <Input type='password' spellCheck={false} placeholder={'Old Password'} />
          )}
        </FormItem>
        <FormItem>
          <p className={'profile-label'}>New Password</p>
          {getFieldDecorator('newPassword', {
            rules: [{
              required: true,
              message: 'New Password is required'
            }, {
              min: 8,
              message: 'Password length must be at least 8 characters '
            }, {
              pattern: pwdPattern,
              message: 'Password must include at least one small and one capital letter, and one sign (eg: $)'
            }, {
              validator: this.validateToNextPassword
            }]
          })(
            <Input type='password' placeholder={'New Password'} />
          )}
        </FormItem>
        <FormItem>
          <p className={'profile-label'}>Confirm New Password</p>
          {getFieldDecorator('confirm', {
            rules: [{
              required: true, message: 'Please confirm your password!'
            }, {
              validator: this.compareToPrevPassword
            }]
          })(
            <Input type='password' placeholder={'Confirm Password'}
              onBlur={this.handleConfirmBlur} />
          )}
          {
            this.state.displayPassErrMessage &&
              <p className={'error'}>Both new passwords should match!</p>
          }
        </FormItem>
        <div className={'btn-group'}>
          <Button type={'primary'} text={'Save Changes'} onClick={this.handleSubmit} />
          <Button type={'default'} text={'Cancel'} onClick={this.props.history.goBack} />
        </div>
      </Form>
    )
  }
}
export default ProfileForm
