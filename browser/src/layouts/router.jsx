// @flow

import React from 'react'
import { IntlProvider, addLocaleData } from 'react-intl'
import en from 'react-intl/locale-data/en'
import zh from 'react-intl/locale-data/zh'
import { Route, Switch } from 'react-router'
import { BrowserRouter } from 'react-router-dom'
import messagesEn from '../assets/locales/en.json'
import messagesZh from '../assets/locales/zh.json'
import Home from '../components/Trade/Trade'
import TradeList from '../components/TradeList/TradeList'
import BidOfferAddNew from '../components/TradeOffers/BidOfferAddNew'
import BidOfferDetails from '../components/TradeOffers/BidOfferDetails/BidOfferDetails'
import BidOfferList from '../components/TradeOffers/BidOfferList/BidOfferList'
import Login from '../components/User/Login'
import SignUp from '../components/User/Signup'
import Profile from '../components/User/Profile'
import PasswordPage from '../components/User/ChangePassword'
import EmailPage from '../components/User/ChangeEmail'
import Preferences from '../components/Common/Preferences/index'
import Logout from '../components/User/Logout'
import Landing from '../components/Landing/Landing'
import MainApp from './MainApp'
import AdminApp from './AdminApp'
import { BrowserSupport } from '../components/BrowserSupport/BrowserSupport'

addLocaleData([...en, ...zh])

const messages = {
  'en': messagesEn,
  'zh': messagesZh
}

export default function RouteComponent () {
  return (
    <IntlProvider locale='en' messages={messages['en']}>
      <BrowserRouter basename='/view'>
        <div>
          <Route path={'/browser-support'} component={BrowserSupport} />
          <Switch>
            <Route exact path={'/'} component={Landing} />
            <Route path={'/login'} component={Login} />
            <Route path={'/signup'} component={SignUp} />
            <Route path={'/logout'} component={Logout} />
            <MainApp path={'/home'} component={Home} />
            <MainApp path={'/trades'} component={TradeList} />
            <MainApp path={'/trade-offer/new'} component={BidOfferAddNew} />
            <MainApp path={'/trade-offer/details'} component={BidOfferDetails} />
            <MainApp path={'/trade-offer/:offer(buy|sell)'} component={BidOfferList} />
            <MainApp path={'/settings/profile'} component={Profile} />
            <MainApp path={'/settings/password'} component={PasswordPage} />
            <MainApp path={'/settings/email'} component={EmailPage} />
            <MainApp path={'/settings/preferences'} component={Preferences} />
            <AdminApp path={'/admin/trades'} component={TradeList} />
            <AdminApp path={'/admin/home'} component={Home} />
            <Route path={'*'} render={() => (<div className={'empty-field'}>Page Not Found</div>)} />
          </Switch>
        </div>
      </BrowserRouter>
    </IntlProvider>
  )
}
