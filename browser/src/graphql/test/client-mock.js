// @flow

import { SchemaLink } from 'apollo-link-schema'
import { makeExecutableSchema } from 'graphql-tools'
import { ApolloClient, InMemoryCache } from 'apollo-boost'
import fs from 'fs'
import { resolvers } from './resolver'
import type { TradeType } from '../../model/flowType'
import { mkEmptyTemplate, mkEmptyUser } from '../../services/generators'

const typeDefs = fs.readFileSync(`${__dirname}/../../../../api/websrv_schema.graphql`, 'utf8')
const cache = new InMemoryCache()
const executableSchema = makeExecutableSchema({
  typeDefs: typeDefs,
  resolvers,
  resolverValidationOptions: {
    requireResolversForResolveType: false
  }
})

export const sampleUserKeyPair1 = {
  secKey: 'SCRHEWSED4VXPCA55HBNWSZ7ESRU2NLLNJRIO532DJBPSB5X3F7UVNRM',
  pubKey: 'GDOKCE5VFBB3CCPWG6HLQXW7AL4QELDSHHLABWDBYMRSI4UGYY4BBGS3'
}

export const stageOpTx = 'AAAAADDjh70/h8oj6Ae81FL+JXPIQKzuaSjL2WI+c6avLIc0AAAAyAAAh74AAAABAAAAAAAAAAAAAAACAAAAAQAAAAAw44e9P4fKI+gHvNRS/iVzyECs7mkoy9liPnOmryyHNAAAAAoAAAADaWR4AAAAAAEAAAABMQAAAAAAAAEAAAAAMOOHvT+HyiPoB7zUUv4lc8hArO5pKMvZYj5zpq8shzQAAAAKAAAABmVudGl0eQAAAAAAAQAAAAxzdGFnZV9hZGRyZXEAAAAAAAAAAA=='
export const stageDocAddTx = 'AAAAADDjh70/h8oj6Ae81FL+JXPIQKzuaSjL2WI+c6avLIc0AAAAyAAAh74AAAABAAAAAAAAAAAAAAACAAAAAQAAAAAw44e9P4fKI+gHvNRS/iVzyECs7mkoy9liPnOmryyHNAAAAAoAAAADaWR4AAAAAAEAAAABMQAAAAAAAAEAAAAAMOOHvT+HyiPoB7zUUv4lc8hArO5pKMvZYj5zpq8shzQAAAAKAAAABmVudGl0eQAAAAAAAQAAAAxzdGFnZV9hZGRyZXEAAAAAAAAAAA=='
export const sampleUser1 = {
  id: '1',
  firstName: 'Sergey',
  lastName: 'Ivanov',
  emails: ['ss@ss.ss'],
  roles: ['trader'],
  pubKey: 'GDOKCE5VFBB3CCPWG6HLQXW7AL4QELDSHHLABWDBYMRSI4UGYY4BBGS3',
  biography: '',
  avatar: '',
  createdAt: new Date(),
  orgMap: null
}

export const sampleTrade1: TradeType = {
  id: '1993134',
  name: 'testTrade1',
  description: 'This is sample test trade data',
  template: mkEmptyTemplate(),
  buyer: {
    id: '1',
    firstName: 'Sergey',
    lastName: 'Ivanov',
    pubKey: 'GDOKCE5VFBB3CCPWG6HLQXW7AL4QELDSHHLABWDBYMRSI4UGYY4BBGS3',
    avatar: '',
    biography: '',
    createdAt: new Date(),
    emails: [],
    orgMap: [],
    roles: []
  },
  seller: {
    id: '2',
    firstName: 'Bahir',
    lastName: 'Abadi',
    pubKey: 'GCZVKTLPQY54OGK5R3EEAU22Q7COES2XP5544H2C3CRNN7GGRL74IBUM',
    avatar: '',
    biography: '',
    createdAt: new Date(),
    emails: [],
    orgMap: [],
    roles: []
  },
  tradeCloseStatus: 'no',
  scAddr: 'GCULNM3FMJM4UXISN4CRMFXVLFOIXSAETM7RMWJLJLYK5NCOGDSF6PUB',
  stages: [{
    name: 'trade contract',
    description: 'this is trade contract stage',
    addReqIdx: -1,
    owner: 'b',
    expiresAt: new Date(),
    docs: [],
    delReqs: [],
    closeReqs: [],
    moderator: null
  }, {
    name: 'stage 1',
    description: 'this is vessel info stage',
    addReqIdx: -1,
    owner: 'b',
    expiresAt: new Date(),
    docs: [],
    delReqs: [],
    closeReqs: [],
    moderator: null
  }],
  stageAddReqs: [],
  closeReqs: [],
  createdAt: new Date(),
  createdBy: mkEmptyUser()
}

export default new ApolloClient({
  link: new SchemaLink({ schema: executableSchema }),
  cache
})
