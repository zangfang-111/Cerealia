// @flow

import React, { useState } from 'react'
import { Badge, Button, Dropdown, Icon, Input, Menu, Select, Tooltip, Spin } from 'antd'
import { observer } from 'mobx-react-lite'
import _ from '../../polyfills/underscore.js'
import { Link } from 'react-router-dom'
import { logoImageLight, logoImageDark } from '../../assets/images'
import StellarKeyInputModal from '../Common/Modals/StellarModal/StellarModal'
import CreateNewTradeModal from '../Common/Modals/ContractModal/CreateNewTradeModal'
import type { TradeType } from '../../model/flowType'
import { addNotificationHelper } from '../../lib/helper'
import currentUser from '../../stores/current-user'
import tradesStore from '../../stores/trade-store'
import modTradesStore from '../../stores/moderatorStore/modTrades'
import appStore from '../../stores/app-store'
import tradeOfferStore from '../../stores/tradeoffer-store'

const Option = Select.Option

export default observer((props: { isAdmin: boolean }) => {
  const [displayStellarKeyModal, setDisplayStellarKeyModal] = useState(false)
  const [spinning, setSpinning] = useState(false)

  async function onRefreshTrades () {
    setSpinning(true)
    try {
      if (props.isAdmin) {
        await modTradesStore.fetchTrades()
      } else {
        await tradesStore.fetchTrades()
      }
      await tradeOfferStore.fetchTradeOffers()
      setSpinning(false)
    } catch (err) {
      addNotificationHelper(err, 'error')
    }
  }

  // on combining headers, still created 2 menus to avoid excessive use of ?:
  const adminUserMenu = (
    <Menu>
      {
        appStore.adminMode &&
        <Menu.Item key='0'>
          <Link to={'/'} onClick={() => appStore.setAdminMode(false)}>
            <i className={'fas fa-users'} />Main App
          </Link>
        </Menu.Item>
      }
      <Menu.Item key='1'>
        <Link to={'/logout'}>
          <i className={'fas fa-sign-out-alt'} />Log out
          <span className={'user-logout-name'}>{currentUser.user.firstName}</span>
        </Link>
      </Menu.Item>
    </Menu>
  )

  const appUserMenu = currentUser.isAuthenticated
    ? (
      <Menu>
        {
          currentUser.hasModeratorRole &&
          <Menu.Item key='0'>
            <Link
              to={{
                pathname: '/admin/trades',
                auth: 'admin'
              }}
              onClick={() => appStore.setAdminMode(true)}
            >
              <i className={'fas fa-user-secret'} />Admin Panel
            </Link>
          </Menu.Item>
        }
        <Menu.Item key='1'>
          <Link to={'/logout'}>
            <i className={'fas fa-sign-out-alt'} />Log out
            <span className={'user-logout-name'}>{currentUser.user.firstName}</span>
          </Link>
        </Menu.Item>
      </Menu>
    )
    : (
      <Menu>
        <Menu.Item key='0'>
          <Link to={'/login'}><i className={'fas fa-sign-in-alt'} />Log in</Link>
        </Menu.Item>
        <Menu.Item key='1'>
          <Link to={'/signup'}><i className={'fas fa-user-plus'} />Sign up</Link>
        </Menu.Item>
      </Menu>
    )

  const spinIcon = <Icon type='loading' style={{ fontSize: 18 }} spin />
  const logoImage = appStore.appTheme === 'theme-light' ? logoImageLight : logoImageDark

  return (
    <div className={'header'}>
      <div className={'logo-image'}>
        <Link to={props.isAdmin ? '/' : '/home'} >
          <img className={'logo-brand'} src={logoImage} alt='logo brand' />
        </Link>
        {process.env.NODE_ENV !== 'production' &&
        <div className='dev-env'>{process.env.NODE_ENV}</div>}
      </div>
      { !props.isAdmin ? <TestUtils /> : null }
      <div className={'header-right'}>
        { spinning
          ? <Spin indicator={spinIcon} spinning={spinning} />
          : <Badge>
            <i className='fas fa-sync-alt' title={'refresh'} onClick={onRefreshTrades} />
          </Badge>
        }
        <Badge count={50} dot>
          <i className='fas fa-envelope' />
        </Badge>
        <Badge count={20} dot>
          <i className='fas fa-bell' />
        </Badge>
        <Badge>
          <i className='fas fa-key' title={'key verify'}
            onClick={() => setDisplayStellarKeyModal(true)} />
        </Badge>
        <Dropdown overlay={props.isAdmin ? adminUserMenu : appUserMenu}>
          <a className='ant-dropdown-link'>
            <i className='fas fa-user' />
            {
              currentUser.userName
            }
            <Icon type='down' />
          </a>
        </Dropdown>
      </div>
      { !props.isAdmin ? <StellarKeyInputModal
        onCloseModal={() => setDisplayStellarKeyModal(false)}
        visible={displayStellarKeyModal} /> : null }
    </div>
  )
})

const TestUtils = observer(() => {
  const [expandSearch, setExpandSearch] = useState(false)
  const [displayCreateNewTradeModal, setDisplayCreateNewTradeModal] = useState(false)
  const Search = Input.Search

  return (
    <div className={'trade-utils'}>
      <Tooltip placement={'bottomLeft'} title={'Start New Contract'}>
        <Button className={'btn-new-trade'}
          onClick={() => setDisplayCreateNewTradeModal(true)}>
          <Icon type='plus' />
        </Button>
      </Tooltip>
      <Select
        className={'select-trades'}
        size={'large'}
        dropdownClassName={'trade-list'}
        placeholder={<i className={'fas fa-handshake select-icon'}>Trades</i>}
        value={
          !_.isEmpty(tradesStore.trades)
            ? tradesStore.trades[tradesStore.selectedTab].name
            : ''
        }
      >
        {
          tradesStore.trades.map(mkTradeListEntry)
        }
      </Select>
      <Search
        size={'large'}
        onSearch={value => console.log(value)}
        onFocus={() => setExpandSearch(true)}
        onBlur={() => setExpandSearch(false)}
        className={`${expandSearch ? 'expand-search' : 'search-trades'}`}
      />
      <CreateNewTradeModal
        visible={displayCreateNewTradeModal}
        onCloseModal={() => setDisplayCreateNewTradeModal(false)}
      />
    </div>
  )
})

function mkTradeListEntry (t: TradeType, i: number) {
  return (
    <Option
      value={t.id}
      key={i}
      onClick={() => tradesStore.setSelectedTab(i)}
    >
      <Link to='/'>{t.name}</Link>
    </Option>
  )
}
