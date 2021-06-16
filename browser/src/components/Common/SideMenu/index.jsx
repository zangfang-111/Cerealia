// @flow

import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'

const appEntries = [
  ['Home', '/home', 'fa-home', []],
  ['Trade List', '/trades', 'fa-handshake', []],
  ['Bids & Offers', '/trade-offer/', 'fa-file-alt', [
    ['Bids', '/trade-offer/buy'],
    ['Offers', '/trade-offer/sell']
  ]],
  ['Chat', '/chat', 'fa-comments', []],
  ['FAQ', '/faq', 'fa-question', []],
  ['Settings', '/settings/', 'fa-cogs', [
    ['Profile edit', '/settings/profile'],
    ['Password change', '/settings/password'], // TODO - link update
    ['Email change', '/settings/email'],
    ['Preferences', '/settings/preferences']
  ]]
]

const adminEntries = [
  ['Main app', '/home', 'fa-home', []],
  ['All Trades', '/admin/trades', 'fa-handshake', []]
]

function renderSubmenuEntry (e: Array<string>) {
  return (<Link key={e[1]} to={e[1]} className={location.pathname === e[1] ? 'active' : ''}>
    {e[0]}</Link>)
}

function menuComponents (entries: Array<any>, submenu: string, setSubmenu: (string) => void) {
  let menu = []
  for (let e of entries) {
    let [name, link, icon, subentries] = e
    icon = 'fas ' + icon
    let active = location.pathname.startsWith('/view' + link) ? 'active' : ''
    let c
    let drilldown = submenu === link ? <i className='fas fa-caret-up' /> : <i className='fas fa-caret-down' />
    if (subentries.length !== 0) {
      c = <div key={link} className={active} >
        <a onClick={() => setSubmenu(link)}><i className={icon} /> {name} {drilldown}</a>
        {submenu === link &&
        <div className={'submenu'} >{subentries.map(renderSubmenuEntry)}</div>}
      </div>
    } else {
      c = <div key={link} className={active} >
        <Link to={link} onClick={() => setSubmenu(link)} ><i className={icon} /> {name}</Link>
      </div>
    }
    menu.push(c)
  }
  return menu
}

export default function (props: {isAdmin: boolean}) {
  const [submenu, setSubmenu] = useState('')
  const [isMobile, setIsMobile] = useState(window.innerWidth < 1150)
  const updateWindowDimensions = () => setIsMobile(window.innerWidth < 1150)

  useEffect(() => {
    window.addEventListener('resize', updateWindowDimensions)
    return () => window.removeEventListener('resize', updateWindowDimensions)
  })

  let entries = props.isAdmin ? adminEntries : appEntries
  return (
    <div className={props.isAdmin ? 'admin-panel' : 'main-panel'}>
      {isMobile &&
      <div><i className={'fas fa-bars'} /></div> }
      {menuComponents(entries, submenu, setSubmenu)}
      {process.env.NODE_ENV !== 'production' &&
      <div className='side-menu-item'>
        Branch: {process.env.REACT_APP_GIT_BRANCH}@{process.env.REACT_APP_GIT_SHA}
      </div>}
    </div>
  )
}
