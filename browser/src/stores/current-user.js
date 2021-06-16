// @flow

import { action, observable, runInAction, computed } from 'mobx'
import {
  userLogin,
  getCurrentUser,
  changePassword,
  changeEmail,
  updateUserProfile
} from '../graphql/trades'
import type {
  UserType,
  LoginInputType,
  ChangePasswordType,
  UserProfileInputType
} from '../model/flowType'
import { userRoleMap } from '../constants/tradeConst'
import { GqlClient } from '../services/cerealia'
import { mkEmptyUser } from '../services/generators'

class CurrentUser {
  @observable user: UserType = mkEmptyUser()
  gqlClient: Object
  constructor () {
    this.gqlClient = GqlClient
  }

  updateClient = (client: Object) => {
    this.gqlClient = client
  }

  @action async authenticate () {
    if (this.user.id) {
      return
    }
    try {
      let response = await this.gqlClient.query({ query: getCurrentUser })
      runInAction('fetchSuccess', () => {
        if (response.data.user) {
          this.user = response.data.user
        }
      })
    } catch (e) {
      console.warn('Invalid authentication', e)
    }
  }

  @action async login (input: LoginInputType) {
    let response = await this.gqlClient.mutate({
      mutation: userLogin,
      variables: { 'input': input }
    })
    await runInAction('fetchSuccess', async () => {
      localStorage.setItem('auth_token', response.data.userLogin.token)
    })
    // Login request doesn't have the session loaded so it can't fetch any non-public data
    await this.authenticate()
  }

  async changePassword (input: ChangePasswordType) {
    await this.gqlClient.mutate({
      mutation: changePassword,
      variables: { 'input': input }
    })
  }

  async changeEmail (input: Array<string>) {
    await this.gqlClient.mutate({
      mutation: changeEmail,
      variables: { 'input': input }
    })
    this.user.emails = input
  }

  @action async updateUserProfile (input: UserProfileInputType) {
    let response = await this.gqlClient.mutate({
      mutation: updateUserProfile,
      variables: { 'input': input }
    })
    runInAction('fetchSuccess', () => {
      this.user = response.data.userProfileUpdate
    })
  }

  @computed get userName (): string {
    let u = this.user
    return u.firstName ? u.firstName.concat(' ', u.lastName) : ''
  }

  @computed get isAuthenticated (): boolean {
    return !!this.user.id
  }

  @computed get hasModeratorRole (): boolean {
    return this.user.roles.indexOf(userRoleMap.moderator) >= 0
  }

  @computed get pubKey (): string {
    return this.user.pubKey
  }
}

export default new CurrentUser()
