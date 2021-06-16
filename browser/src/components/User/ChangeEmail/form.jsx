// @flow

import React from 'react'
import { Form, Input } from 'antd'
import { withRouter } from 'react-router-dom'
import currentUser from '../../../stores/current-user'

import Button from '../../Common/Button/Button'

type Props = {
  form: Object,
  history: Object,
  onSaveEmail: Function,
}

type State = {
  emails: Array<string>
}

const FormItem = Form.Item

@withRouter
class EmailForm extends React.Component<Props, State> {
  constructor (props: Props) {
    super(props)
    this.state = {
      emails: []
    }
  }

  componentDidMount = async () => {
    await currentUser.authenticate()
    this.setState({ emails: currentUser.user.emails })
  }

  handleSubmit = (e: SyntheticEvent<HTMLButtonElement>) => {
    this.props.form.validateFieldsAndScroll((err, values) => {
      if (!err) {
        this.props.onSaveEmail(Object.values(values))
      }
    })
  }

  renderEmailField = (email: string, index: number) => {
    index++
    const { getFieldDecorator } = this.props.form
    return (
      <FormItem key={index}>
        <p className={'profile-label'}>{'Email' + index.toString()}
          <a className={'profile-edit'} onClick={() => this.deleteEmailField(index - 1)}>| Delete</a>
        </p>
        {getFieldDecorator('email' + index.toString(), {
          initialValue: email,
          rules: [{
            type: 'email', message: 'This is not valid E-mail!'
          }, {
            required: true, message: 'E-mail is required'
          }]
        })(
          <Input placeholder={'Email' + index.toString()} spellCheck={false} />
        )}
      </FormItem>
    )
  }

  addEmailField = () => {
    let emails = this.state.emails
    emails.push('')
    this.setState({ emails: emails })
  }

  deleteEmailField = (index: number) => {
    let emails = this.state.emails
    emails.splice(index, 1)
    this.setState({ emails: emails })
  }

  render () {
    return (
      <Form className={'profile-form'} >
        <div className={'content-title'}>Edit Your Emails</div>
        <p className={'add-email'}>
          <a onClick={this.addEmailField}>
            <i className={'fas fa-plus-circle'} /> Add New Email
          </a>
        </p>
        {
          this.state.emails.length &&
          this.state.emails.map(this.renderEmailField)
        }
        <div className={'btn-group'}>
          <Button type={'primary'} text={'Save Changes'} onClick={this.handleSubmit} />
          <Button type={'default'} text={'Cancel'} onClick={this.props.history.goBack} />
        </div>
      </Form>
    )
  }
}
export default EmailForm
