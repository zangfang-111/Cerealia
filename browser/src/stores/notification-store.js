// @flow

import { action, observable, runInAction } from 'mobx'
import { notifications, notificationsTrade, notificationDismiss } from '../graphql/trades'
import type { Notification } from '../model/flowType'
import { GqlClient } from '../services/cerealia'

class NotificationStore {
  @observable notifications: Array<Notification> = []
  @observable triedToFetchNotification: boolean = false

  gqlClient: Object
  constructor () {
    this.gqlClient = GqlClient
  }

  @action async fetchAllNotifications (i: number) {
    this.triedToFetchNotification = false
    let response = await this.gqlClient.query({
      query: notifications,
      variables: { 'from': i }
    })
    runInAction('fetchSuccess', () => {
      this.notifications = response.data.notifications
      this.triedToFetchNotification = true
    })
  }

  async fetchTradeNotifications (id: string) {
    let response = await this.gqlClient.query({
      query: notificationsTrade,
      variables: { 'id': id }
    })
    return response.data.notificationsTrade
  }

  dismiss (id: string) {
    this.gqlClient.mutate({
      mutation: notificationDismiss,
      variables: { 'id': id }
    })
  }
}

export default new NotificationStore()
