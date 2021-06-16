// @flow

import React from 'react'
import { withRouter, Link } from 'react-router-dom'
import { Form, Icon, Input } from 'antd'
import Button from '../../../components/Common/Button/Button'

const FormItem = Form.Item
type Props = {
  form: Object,
  onLogin: Function,
  history: Object
}

type State = {
  error: boolean,
  email: string,
  password: string
}

@withRouter
class PasswordLoginForm extends React.Component<Props, State> {
  constructor (props: Props) {
    super(props)
    this.state = {
      error: false,
      email: '',
      password: ''
    }
  }

  handleChangeEmail = (e: SyntheticInputEvent<HTMLInputElement>): void => {
    this.setState({ email: e.currentTarget.value })
  }

  handleChangePassword = (e: SyntheticInputEvent<HTMLInputElement>): void => {
    this.setState({ password: e.currentTarget.value })
  }

  onLogin = (): void => {
    this.props.form.validateFieldsAndScroll((err, values) => {
      if (err) {
        this.setState({ error: true })
      } else {
        this.props.onLogin(this.state.email, this.state.password)
      }
    })
  }

  render () {
    const { getFieldDecorator } = this.props.form
    return (
      <Form className={'login-form'}>
        <div className={'content-title form-title'}>Login</div>
        <FormItem>
          {getFieldDecorator('email', {
            rules: [{
              type: 'email', message: 'This is not a valid emai!'
            }, {
              required: true, message: 'E-mail required'
            }]
          })(
            <Input prefix={<Icon type={'user'} className={'login-prefix'} />}
              placeholder={'Email'} spellCheck={false} onChange={this.handleChangeEmail} onPressEnter={this.onLogin} />
          )}
        </FormItem>
        <FormItem>
          {getFieldDecorator('password', {
            rules: [{ required: true, message: 'Password required' }]
          })(
            <Input prefix={<Icon type={'lock'} className={'login-prefix'} />}
              type={'password'} placeholder={'Password'} onChange={this.handleChangePassword} onPressEnter={this.onLogin} />
          )}
        </FormItem>
        <div className={'forgot-password'}>
          <Link to={'/forgot-password'}>Forgot password?</Link>
        </div>
        <div className={'btn-group'}>
          <Button type={'primary'} text={'Login'} onClick={this.onLogin} />
          <Button type={''} text={'Create New Account'} onClick={() => this.props.history.push('/signup')} />
        </div>
      </Form>
    )
  }
}

export default PasswordLoginForm
