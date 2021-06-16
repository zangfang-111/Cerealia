// @flow

import { useEffect } from 'react'
import useReactRouter from 'use-react-router'

export default function () {
  const { history } = useReactRouter()
  // linting: complex expression that is used in the dependency array.
  const localAuthToken = localStorage.getItem('auth_token')
  useEffect(() => {
    if (localAuthToken) {
      history.push('/home')
    } else {
      history.push('/login')
    }
  }, [localAuthToken, history])
  return null
}
