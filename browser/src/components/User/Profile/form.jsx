// @flow

import { AutoComplete, Col, Form, Input, Row, Tag } from 'antd'
import React from 'react'
import DropZone from 'react-dropzone'
import { withRouter } from 'react-router-dom'
import Button from '../../../components/Common/Button/Button'
import { addNotificationHelper } from '../../../lib/helper'
import type { OrgMapType, OrganizationType, UserProfileInputType } from '../../../model/flowType'
import { mkPublicLink } from '../../../services/cerealia'
import usersStore from '../../../stores/user-store'
import currentUser from '../../../stores/current-user'
import CreateNewOrgModal from '../CreateOrgModal/CreateNewOrgModal'

const FormItem = Form.Item
const TextArea = Input.TextArea

type Props = {
  form: Object,
  onSave: Function,
  history: Object
}

type CompanyType = {
  orgItem: OrgMapType,
  noMatchOrgErr: boolean,
  curOrgName: string
}

type State = {
  keyError: boolean,
  avatarFile: Object,
  avatar: string,
  newPassword: string,
  changeAvatar: boolean,
  visibleAddOrgModal: boolean,
  curOrgIdx: number,
  companies: Array<CompanyType>
}

@withRouter
class ProfileForm extends React.Component<Props, State> {
  constructor (props: Props) {
    super(props)
    this.state = {
      keyError: false,
      avatarFile: {},
      avatar: '',
      newPassword: '',
      changeAvatar: false,
      visibleAddOrgModal: false,
      curOrgIdx: 0,
      companies: []
    }
  }

  onOpenAddOrgModal = (index: number) => this.setState({
    curOrgIdx: index,
    visibleAddOrgModal: true
  })

  onCloseAddOrgModal = () => this.setState({ visibleAddOrgModal: false })

  componentDidMount = () => {
    let companies = []
    for (let o of currentUser.user.orgMap) {
      companies.push({
        orgItem: o,
        noMatchOrgErr: false,
        curOrgName: o.org.name
      })
    }
    this.setState({
      avatar: mkPublicLink('/user/avatars/' + currentUser.user.avatar),
      companies
    })
  }

  handleSubmit = (e: SyntheticEvent<HTMLButtonElement>) => {
    if (this.state.companies.find(company => company.noMatchOrgErr === true)) {
      return
    }
    this.props.form.validateFieldsAndScroll((err, values) => {
      if (!err) {
        this.props.onSave(this.mkChangeData(values), this.state.avatarFile)
      }
    })
  }

  mkChangeData = (values: Object): UserProfileInputType => {
    let newUserData = {}
    let orgMap = []
    newUserData.firstName = values.firstName
    newUserData.lastName = values.lastName
    newUserData.biography = values.biography
    for (let company of this.state.companies) {
      orgMap.push({
        id: company.orgItem.org.id,
        role: company.orgItem.role
      })
    }
    newUserData.orgMap = orgMap
    return newUserData
  }

  showAvatar = () => {
    if (this.state.avatarFile.name || currentUser.user.avatar) {
      return React.createElement('img', {
        src: this.state.avatar,
        className: 'card-image'
      })
    }
    return React.createElement('i', { className: 'fas fa-user-circle avatar-icon' })
  }

  handleDrop = (files: Array<Object>) => {
    if (!files.length) {
      return
    }
    this.setState({
      avatarFile: files[0],
      avatar: files[0].preview,
      changeAvatar: true
    })
  }

  onCancelChangeAvatar = () => {
    this.setState({
      avatarFile: {},
      avatar: mkPublicLink('/user/avatars/' + currentUser.user.avatar),
      changeAvatar: false
    })
  }

  addCompany = () => {
    let companies = this.state.companies
    companies.push({
      orgItem: {
        org: { name: '', id: '', email: '', telephone: '', address: '' },
        role: ''
      },
      noMatchOrgErr: false,
      curOrgName: ''
    })
    this.setState({ companies })
  }

  deleteOrgMap = (index: number) => {
    let companies = this.state.companies
    companies.splice(index, 1)
    this.setState({ companies })
  }

  onChangeRole = (e: SyntheticInputEvent<HTMLInputElement>, index: number) => {
    let companies = this.state.companies
    companies[index].orgItem.role = e.currentTarget.value
    this.setState({ companies })
  }

  searchOrganizations = (value: string, index: number) => {
    if (!value) {
      return
    }
    let found = false
    for (let o of usersStore.organizations) {
      if (o.name.toLowerCase().indexOf(value.toLowerCase()) >= 0) {
        found = true
        break
      }
    }
    let companies = this.state.companies
    companies[index].noMatchOrgErr = !found
    companies[index].curOrgName = value
    this.setState({ companies })
  }

  selectOrganization = (value: string, index: number) => {
    for (let o of usersStore.organizations) {
      if (o.name === value) {
        let companies = this.state.companies
        companies[index].orgItem.org = o
        this.setState({ companies })
        break
      }
    }
  }

  addNewOrg = (org: OrganizationType) => {
    const sameOrg = usersStore.organizations.find(
      item => item.name === org.name
    )
    if (sameOrg) {
      addNotificationHelper('Same company name already exists', 'error')
      return
    }
    let companies = this.state.companies
    companies[this.state.curOrgIdx].orgItem.org = org
    companies[this.state.curOrgIdx].noMatchOrgErr = false
    this.setState({ companies })
    this.onCloseAddOrgModal()
  }

  renderOrgMap = (orgMapItem: CompanyType, index: number) => {
    const { getFieldDecorator } = this.props.form
    return (
      <Row key={index}>
        <Col span={11} >
          <FormItem>
            <p className={'profile-label'}>{'Company' + (index + 1).toString()}</p>
            {getFieldDecorator('company' + index.toString(), {
              initialValue: orgMapItem.orgItem.org.name,
              rules: [{ required: true, message: 'Company is required', whitespace: true }]
            })(
              <AutoComplete
                dataSource={usersStore.organizations.map(org => org.name)}
                onSelect={(value) => this.selectOrganization(value, index)}
                onSearch={(value) => this.searchOrganizations(value, index)}
                placeholder={'Company'}
                filterOption={(input, option) => option.props.children.toLowerCase().indexOf(
                  input.toLowerCase()) !== -1
                }
              />
            )}
            {
              this.state.companies[index].noMatchOrgErr &&
              <p className={'error ant-form-explain'}>
                No matches! <a onClick={() => this.onOpenAddOrgModal(index)}>Create</a></p>
            }
          </FormItem>
        </Col>
        <Col span={11} offset={2} >
          <FormItem>
            <p className={'profile-label'}>{'Role' + (index + 1).toString()}
              <a className={'profile-edit'} onClick={() => this.deleteOrgMap(index)}>| Delete</a>
            </p>
            {getFieldDecorator('role' + index.toString(), {
              initialValue: orgMapItem.orgItem.role,
              rules: [{
                required: true, message: 'Role is required'
              }]
            })(
              <Input onChange={(e) => this.onChangeRole(e, index)}
                placeholder={'Role'} spellCheck={false} />
            )}
          </FormItem>
        </Col>
      </Row>
    )
  }

  render () {
    const { getFieldDecorator } = this.props.form
    const user = currentUser.user
    return (
      <Form className={'profile-form'} >
        <div className={'content-title'}>Edit Your Profile</div>
        <div className={'profile-avatar'}>
          { this.showAvatar() }
          <DropZone accept={'image/jpeg,image/jpg,image/png,image/gif'}
            maxSize={2000000}
            onDrop={this.handleDrop}
            multiple={false}
            className={'avatar-text'}>
              Change your avatar
          </DropZone>
          {
            this.state.changeAvatar &&
            <a className={'profile-cancel'} onClick={this.onCancelChangeAvatar}>| Cancel</a>
          }
        </div>
        <FormItem>
          <p className={'profile-label'}>First Name</p>
          {getFieldDecorator('firstName', {
            initialValue: user.firstName,
            rules: [{ required: true, message: 'Please input your firstname!', whitespace: true }]
          })(
            <Input spellCheck={false} />
          )}
        </FormItem>
        <FormItem>
          <p className={'profile-label'}>Last Name</p>
          {getFieldDecorator('lastName', {
            initialValue: user.lastName,
            rules: [{ required: true, message: 'Please input your lastname!', whitespace: true }]
          })(
            <Input spellCheck={false} />
          )}
        </FormItem>
        <div className={'page-title'}>Organizations</div>
        <br />
        {
          this.state.companies.length > 0 &&
          this.state.companies.map(this.renderOrgMap)
        }
        <div className={'align-right'}>
          <Tag color={'#20a040'} onClick={this.addCompany} >Add new company</Tag>
        </div>
        <hr />
        <br />
        <FormItem>
          <p className={'profile-label'}>PublicKeys</p>
          {getFieldDecorator('publicKey', {
            initialValue: user.pubKey,
            rules: [{ required: true, message: 'Please input your public key!', whitespace: false }]
          })(
            <Input spellCheck={false} disabled />
          )}
          {
            this.state.keyError && <p className={'invalid-key'}>this is not valid stellar public key</p>
          }
        </FormItem>
        <FormItem>
          <p className={'profile-label'}>Biography</p>
          {getFieldDecorator('biography', {
            initialValue: user.biography,
            rules: [{ required: false, whitespace: true }]
          })(
            <TextArea rows={4} placeholder={'Biography'} />
          )}
        </FormItem>
        <div className={'btn-group'}>
          <Button type={'primary'} text={'Save Changes'} onClick={this.handleSubmit} />
          <Button type={'default'} text={'Cancel'} onClick={() => this.props.history.goBack()} />
        </div>
        {
          this.state.visibleAddOrgModal &&
          <CreateNewOrgModal
            onCloseModal={this.onCloseAddOrgModal}
            onAddNewOrg={this.addNewOrg}
            visible={this.state.visibleAddOrgModal}
            curOrgName={this.state.companies[this.state.curOrgIdx].curOrgName}
            usersStore={usersStore}
          />
        }
      </Form>
    )
  }
}
export default ProfileForm
