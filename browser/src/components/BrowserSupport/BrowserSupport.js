// @flow

import React from 'react'

export const BrowserSupport = () => {
  return (
    <div className={'browser-support'}>
      <p className={'title'}>Your browser is not supported. Please run the app in one of the following browsers:</p>
      <ul className={'browser-list'}>
        <li>Chrome (or Chrome based), version &ge; 64.0</li>
        <li>Firefox, version &ge;  58</li>
        <li>Safari, version &ge; 10.1</li>
        <li>Microsoft Edge, version &ge; 16</li>
        <li>We don't support Opera and Internet Explorer</li>
      </ul>
    </div>
  )
}
