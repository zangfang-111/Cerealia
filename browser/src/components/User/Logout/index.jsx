// @flow

import { useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import useReactRouter from 'use-react-router'
import currentUser from '../../../stores/current-user'

export default observer(() => {
  const { history } = useReactRouter()
  const isAuth = currentUser.isAuthenticated
  useEffect(() => {
    localStorage.setItem('auth_token', '')
    if (isAuth) {
      window.location.reload()
    } else {
      history.push('/login')
    }
  }, [isAuth, history])
  return null
})
