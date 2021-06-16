// @flow

import { AutoComplete, Button, Col, Form, Input, Row } from 'antd'
import React from 'react'
import type { OrganizationType } from '../../../model/flowType'
import stellarStore from '../../../stores/stellarStore'
import usersStore from '../../../stores/user-store'
import CreateNewOrgModal from '../CreateOrgModal/CreateNewOrgModal'

const FormItem = Form.Item
const TextArea = Input.TextArea
const defaultAvatar = ''
const pwdPattern = new RegExp(`^(?=.*[a-z])(?=.*[A-Z])(?=.*[!@#$%^&*,.;'"()])`)

type Props = {
  form: Object,
  onSignup: Function
}

type State = {
  confirmDirty: boolean,
  displayPassErrMessage: boolean,
  keyError: boolean,
  selectedOrg: ?OrganizationType,
  noMatchOrgErr: boolean,
  visibleAddOrgModal: boolean,
  curOrgName: string
}

type ruleType = {
  validator: Function,
  field: string,
  fullField: string,
  type: string
}

class SignupForm extends React.Component<Props, State> {
  constructor (props: Props) {
    super(props)

    this.state = {
      confirmDirty: false,
      displayPassErrMessage: false,
      keyError: false,
      selectedOrg: null,
      noMatchOrgErr: false,
      visibleAddOrgModal: false,
      curOrgName: ''
    }
  }

  onOpenAddOrgModal = () => this.setState({ visibleAddOrgModal: true })

  onCloseAddOrgModal = () => this.setState({ visibleAddOrgModal: false })

  handleSubmit = (e: SyntheticEvent<HTMLButtonElement>) => {
    e.preventDefault()
    this.props.form.validateFieldsAndScroll((err: any, values: any) => {
      if (!err && !this.state.keyError && !this.state.noMatchOrgErr) {
        values.avatar = defaultAvatar
        delete values.confirm
        values.orgID = this.state.selectedOrg && this.state.selectedOrg.id
        this.props.onSignup(values)
      }
    })
  }

  validateKey = (e: SyntheticInputEvent<HTMLInputElement>) => {
    this.setState({ keyError: !stellarStore.validateStellarPublicKey(e.currentTarget.value) })
  }

  handleConfirmBlur = (e: SyntheticEvent<HTMLInputElement>) => {
    this.setState({ confirmDirty: this.state.confirmDirty || !!e.currentTarget.value })
  }

  compareToFirstPassword = (rule: ruleType, value: string, callback: Function) => {
    const form = this.props.form
    if (value && value !== form.getFieldValue('password')) {
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

  searchOrganizations = (value: string) => {
    let found = false
    for (let o of usersStore.organizations) {
      if (o.name.toLowerCase().indexOf(value.toLowerCase()) >= 0) {
        found = true
        break
      }
    }
    this.setState({
      curOrgName: value,
      noMatchOrgErr: value !== '' && !found
    })
  }

  selectOrganization = (value: string) => {
    for (let o of usersStore.organizations) {
      if (o.name === value) {
        this.setState({ selectedOrg: o })
        return
      }
    }
    this.setState({ selectedOrg: undefined })
  }

  addNewOrg = (org: OrganizationType) => {
    this.setState({ selectedOrg: org, noMatchOrgErr: false })
    this.onCloseAddOrgModal()
  }

  render () {
    const { getFieldDecorator } = this.props.form
    return (
      <Form className={'signup-form'} onSubmit={this.handleSubmit} >
        <div className={'content-title form-title'}>Sign up</div>
        <Row>
          <Col sm={12} xs={24} key={1}>
            <FormItem>
              {getFieldDecorator('firstName', {
                rules: [{ required: true, message: 'Firstname is required', whitespace: true }]
              })(
                <Input placeholder={'Firstname'} spellCheck={false} />
              )}
            </FormItem>
          </Col>
          <Col sm={12} xs={24} key={2}>
            <FormItem>
              {getFieldDecorator('lastName', {
                rules: [{ required: true, message: 'Lastname is required', whitespace: true }]
              })(
                <Input placeholder={'Lastname'} spellCheck={false} />
              )}
            </FormItem>
          </Col>
        </Row>
        <Row>
          <Col sm={12} xs={24} key={3}>
            <FormItem>
              {getFieldDecorator('email', {
                rules: [{
                  type: 'email', message: 'This is not valid E-mail!'
                }, {
                  required: true, message: 'E-mail is required'
                }]
              })(
                <Input placeholder={'Email'} spellCheck={false} />
              )}
            </FormItem>
          </Col>
          <Col sm={12} xs={24} key={4}>
            <FormItem>
              {getFieldDecorator('publicKey', {
                rules: [{ required: true, message: 'Public key is required', whitespace: false }]
              })(
                <Input placeholder={'Public Key'} className={this.state.keyError ? 'invalid' : ''}
                  onChange={this.validateKey} spellCheck={false} />
              )}
            </FormItem>
          </Col>
        </Row>
        <Row>
          <Col sm={12} xs={24} key={5}>
            <FormItem>
              {getFieldDecorator('password', {
                rules: [{
                  required: true,
                  message: 'Password is required'
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
                <Input type='password' placeholder={'Password'} />
              )}
            </FormItem>
          </Col>
          <Col sm={12} xs={24} key={6}>
            <FormItem>
              {getFieldDecorator('confirm', {
                rules: [{
                  required: true, message: 'Please confirm your password!'
                }, {
                  validator: this.compareToFirstPassword
                }]
              })(
                <Input type='password' placeholder={'Confirm Password'}
                  onBlur={this.handleConfirmBlur} />
              )}
              {
                this.state.displayPassErrMessage &&
                  <p className={'error ant-form-explain'}>Two passwords that you enter is inconsistent!</p>
              }
            </FormItem>
          </Col>
        </Row>
        <Row>
          <Col sm={12} xs={24} key={7}>
            <FormItem>
              {getFieldDecorator('orgID', {
                rules: [{ required: true, message: 'Company is required', whitespace: true }]
              })(
                <AutoComplete
                  dataSource={usersStore.organizations.map(org => org.name)}
                  onSelect={this.selectOrganization}
                  onSearch={this.searchOrganizations}
                  placeholder={'Company'}
                  filterOption={(input, option) => option.props.children.toLowerCase().indexOf(
                    input.toLowerCase()) !== -1
                  }
                />
              )}
              {
                this.state.noMatchOrgErr &&
                  <p className={'error ant-form-explain'}>
                    No matches! <a onClick={this.onOpenAddOrgModal}>Create</a></p>
              }
            </FormItem>
          </Col>
          <Col sm={12} xs={24} key={8}>
            <FormItem>
              {getFieldDecorator('biography', {
                rules: [{ required: false, whitespace: true }]
              })(
                <TextArea rows={4} placeholder={'Biography'} />
              )}
            </FormItem>
          </Col>
          <Col sm={12} xs={24} key={9}>
            <FormItem className={'role-input'}>
              {getFieldDecorator('orgRole', {
                rules: [{ required: true, message: 'Company role is required', whitespace: true }]
              })(
                <Input placeholder={'Role'} />
              )}
            </FormItem>
          </Col>
        </Row>
        <Row className={'singup-bottom'}>
          <p>By signing up you are agreeing to
            Cerealia's <a>Terms of use</a> and <a>Privacy Policy</a>
          </p>
          <FormItem className={'button-container'}>
            <Button htmlType={'submit'} className={'violet-button'}>Sign Up
            </Button>
          </FormItem>
        </Row>
        <CreateNewOrgModal
          onCloseModal={this.onCloseAddOrgModal}
          onAddNewOrg={this.addNewOrg}
          visible={this.state.visibleAddOrgModal}
          curOrgName={this.state.curOrgName}
        />
      </Form>
    )
  }
}
export default SignupForm
