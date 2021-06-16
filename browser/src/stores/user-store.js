// @flow

import { action, observable, runInAction } from 'mobx'
import {
  getAllUsers,
  userSignup,
  changePassword,
  createOrganization,
  getAllOrganizations
} from '../graphql/trades'
import type {
  UserType,
  OrgInputType,
  NewUserInputType,
  ChangePasswordType,
  OrganizationType
} from '../model/flowType'
import { GqlClient } from '../services/cerealia'

class UsersStore {
  @observable.ref users: Array<UserType> = []
  @observable organizations: Array<OrganizationType> = []
  gqlClient: Object
  constructor () {
    this.gqlClient = GqlClient
  }

  @action async fetchUsers () {
    let response = await this.gqlClient.query({ query: getAllUsers })
    runInAction('fetchSuccess', () => {
      this.users = response.data.users
    })
  }

  @action async fetchOrganizations () {
    let response = await this.gqlClient.query({ query: getAllOrganizations })
    runInAction('fetchSuccess', () => {
      this.organizations = response.data.organizations
    })
  }

  async signup (input: NewUserInputType) {
    await this.gqlClient.mutate({
      mutation: userSignup,
      variables: { 'input': input }
    })
  }

  async changePassword (input: ChangePasswordType) {
    await this.gqlClient.mutate({
      mutation: changePassword,
      variables: { 'input': input }
    })
  }

  async createOrganization (input: OrgInputType) {
    let response = await this.gqlClient.mutate({
      mutation: createOrganization,
      variables: { 'input': input }
    })
    return response.data.organizationCreate
  }

  getUserByID (userID: string) {
    return this.users.find(user => user.id === userID) || undefined
  }
}

export default new UsersStore()
