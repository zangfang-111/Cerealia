// @flow

import React, { useEffect } from 'react'
import { Route } from 'react-router-dom'
import useReactRouter from 'use-react-router'
import { observer } from 'mobx-react-lite'
import usersStore from '../stores/user-store'
import currentUser from '../stores/current-user'
import Header from '../components/Header/Header'
import SideMenu from '../components/Common/SideMenu'
import tradesStore from '../stores/trade-store'
import stellarStore from '../stores/stellarStore'
import tradeOffersStore from '../stores/tradeoffer-store'
import { addNotificationHelper } from '../lib/helper'
import notificationStore from '../stores/notification-store'
import { Row, Col } from 'antd'

async function sync (history: Object) {
  await currentUser.authenticate()
  if (!currentUser.isAuthenticated) {
    history.push('/login')
    return
  }
  try {
    const fetchTrades = !tradesStore.triedToFetchTrades &&
      tradesStore.fetchTrades()
    const fetchTradeOffer = !tradeOffersStore.triedToFetchTradeOffers &&
      tradeOffersStore.fetchTradeOffers()
    const fetchUsers = usersStore.users.length === 0 &&
      usersStore.fetchUsers()
    const initializeStellarNetwork = stellarStore.initializeStellarNetwork()
    const fetchAllNotifications = !notificationStore.triedToFetchNotification &&
      notificationStore.fetchAllNotifications(0)
    await Promise.all([
      fetchTrades,
      fetchTradeOffer,
      fetchUsers,
      initializeStellarNetwork,
      fetchAllNotifications
    ])
  } catch (err) {
    addNotificationHelper(err, 'error')
  }
}

export default observer((props: any) => {
  const { history } = useReactRouter()
  const isAdmin = false

  useEffect(() => { sync(history) }) // useEffect doesn't support async functions correctly

  return (
    <React.Fragment>
      <Header isAdmin={isAdmin} />
      <Row type='flex' className={'view-container'}>
        <Col span={6} className={'side-menu'} >
          <SideMenu isAdmin={isAdmin} />
        </Col>
        <Col span={18} >
          <div className={'view-right'}>
            <Route {...props} />
          </div>
        </Col>
      </Row>
    </React.Fragment>
  )
})
