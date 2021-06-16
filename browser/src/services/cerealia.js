// @flow

import { ApolloClient, ApolloLink, HttpLink, InMemoryCache } from 'apollo-boost'
import { mkJWTHeader } from '../lib/helper'

var host = process.env.REACT_APP_API_HOST || ''
if (host === '') {
  throw Error('env REACT_APP_API_HOST is not set')
}
if (!host.startsWith('http')) {
  throw Error('REACT_APP_API_HOST has wrong value. Should have the follwing form: "http[s]://hostname[:portnum]/" without tailing /')
}
if (host.endsWith('/')) {
  host = host.slice(-1)
}

/* mkLink constructs a proper link for the Cerealia backend
   @param path represents the REST endpoint.
      The leadin slash will be added if it's not present.
 */
export function mkLink (path: string) {
  if (!path.startsWith('/')) {
    path = '/' + path
  }
  return host + path
}

export function mkPublicLink (path: string) {
  if (!path.startsWith('/')) {
    path = '/' + path
  }
  return host + '/assets/public' + path
}

const httpLink = new HttpLink({
  uri: mkLink('/query') })

const authLink = new ApolloLink((operation, forward) => {
  operation.setContext({
    headers: mkJWTHeader()
  })
  return forward(operation)
})

export const GqlClient = new ApolloClient({
  link: authLink.concat(httpLink),
  cache: new InMemoryCache()
})
