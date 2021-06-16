// @flow

import React, { useEffect } from 'react'
import { Route } from 'react-router-dom'
import useReactRouter from 'use-react-router'
import { observer } from 'mobx-react-lite'
import { addNotificationHelper } from '../lib/helper'
import currentUser from '../stores/current-user'
import modTradesStore from '../stores/moderatorStore/modTrades'
import Header from '../components/Header/Header'
import SideMenu from '../components/Common/SideMenu'
import { Row, Col } from 'antd'

async function sync (history: Object) {
  await currentUser.authenticate()
  if (!currentUser.isAuthenticated) {
    history.push('/login')
  } else if (!currentUser.hasModeratorRole) {
    history.goBack()
  } else {
    try {
      const fetchTrades = !modTradesStore.triedToFetchTrades && modTradesStore.fetchTrades()
      await Promise.all([fetchTrades])
    } catch (err) {
      addNotificationHelper(err, 'error')
    }
  }
}

export default observer((props: any) => {
  const { history } = useReactRouter()
  const isAdmin = true

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
