// @flow

import './styles/index.scss'

import { onError } from 'mobx-react'
import { configure } from 'mobx'
import React from 'react'
import ReactDOM from 'react-dom'
import { configureDevtool } from 'mobx-react-devtools'
import App from './App.jsx'
import registerServiceWorker from './registerServiceWorker'
import usersStore from './stores/user-store'
import tradesStore from './stores/trade-store'
import tradeViewStore from './stores/tradeViewStore'
import stellarStore from './stores/stellarStore'
import spinnerStore from './stores/spinner-store'
import modTradesStore from './stores/moderatorStore/modTrades'
import appStore from './stores/app-store'

onError(error => console.log(error))

let env = (process.env.NODE_ENV || '').toUpperCase()

if (env.indexOf('PROD') === -1) {
  // Any configurations are optional
  configureDevtool({
    // Turn on logging changes button programmatically:
    logEnabled: true,
    // Turn off displaying components updates button programmatically:
    updatesEnabled: false,
    logFilter: change => change.type === 'reaction'
  })
}

configure({
  enforceActions: 'observed'
})

const stores = {
  usersStore: usersStore,
  tradesStore: tradesStore,
  tradeViewStore: tradeViewStore,
  stellarStore: stellarStore,
  spinnerStore: spinnerStore,
  modTradesStore: modTradesStore,
  appStore: appStore
}

// For easier debugging
window._____APP_STATE_____ = stores

ReactDOM.render(
  <App />,
  document.getElementById('root')
)
registerServiceWorker()
