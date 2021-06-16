// @flow

import React from 'react'
import DevTools from 'mobx-react-devtools'
import Layouts from './layouts/router.jsx'
import { BrowserSupport } from './components/BrowserSupport/BrowserSupport'
import { isSupportedBrowser } from './lib/helper'

export default function App () {
  let env = (process.env.NODE_ENV || '').toUpperCase()
  return (
    <div className='container'>
      {
        !isSupportedBrowser()
          ? <BrowserSupport />
          : <Layouts />
      }
      { env.indexOf('PROD') === -1 ? <DevTools /> : null }
    </div>
  )
}
