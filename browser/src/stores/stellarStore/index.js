// @flow

import { action, observable, runInAction } from 'mobx'
import {
  getStellarInfo,
  mkTradeCloseTx,
  mkTradeStageAddTx,
  mkTradeStageCloseTx,
  mkTradeStageDocTx
} from '../../graphql/trades'
import type {
  TradeStageDocPathType,
  TradeStagePathType
} from '../../model/flowType'
import type { StellarNetwork } from '../../model/store'
import { GqlClient } from '../../services/cerealia'

const StellarBase = process.env.REACT_APP_USE_CDN === 'no'
  ? require('stellar-base')
  : window.StellarBase
const StellarSdk = process.env.REACT_APP_USE_CDN === 'no'
  ? require('stellar-sdk')
  : window.StellarSdk

const steexpNetworks = {
  'noop': 'NOOP/',
  'horizon-main': 'https://steexp.com',
  'horizon-test': 'https://testnet.steexp.com'
}

class StellarNOOP {}

class StellarStore {
  gqlClient: Object
  server: Object
  keyPair: Object
  @observable network: StellarNetwork
  @observable keyVerified: boolean = false

  constructor () {
    this.gqlClient = GqlClient
    this.server = {}
    // This variable is a "private property"
    // The key should be encapsulated by stellar's browser wallet.
    this.keyPair = {}
  }

  updateClient = (client: Object) => {
    this.gqlClient = client
  }

  linkToSteexp (account: string): string {
    if (!this.network) {
      return 'uninitialized/'
    }
    const url = steexpNetworks[this.network.name] || 'undefined'
    return `${url}/account/${account}#operations`
  }

  linkToHorizonOperations (account: string): string {
    if (!this.network) {
      return 'uninitialized/'
    }
    return `${this.network.url}/accounts/${account}/operations`
  }

  @action async initializeStellarNetwork () {
    if (this.network) {
      return
    }
    let r = await this.gqlClient.query({ query: getStellarInfo })
    runInAction('fetchSuccess', () => {
      StellarSdk.Network.use(new StellarSdk.Network(r.data.stellarNet.passphrase))
      if (r.data.stellarNet.url.startsWith('noop')) {
        console.warn('Using NOOP Stellar Server')
        this.server = new StellarNOOP()
      } else {
        this.network = r.data.stellarNet
        console.log('Stellar Server: ', this.network.url)
        this.server = new StellarSdk.Server(this.network.url)
      }
    })
  }

  sign = (txb: StellarSdk.TransactionBuilder): StellarSdk.TransactionBuilder => {
    txb.sign(this.keyPair)
    return txb
  }

  signRawTX = (tx: string): string => {
    let txb = new StellarSdk.Transaction(tx)
    // constructs js-xdr object with signed transaction builder to send it to backend
    return this.sign(txb).toEnvelope().toXDR().toString('base64')
  }

  validateStellarSecretKey = (secretKey: string): void => {
    return StellarBase.StrKey.isValidEd25519SecretSeed(secretKey)
  }

  validateStellarPublicKey = (publicKey: string): void => {
    return StellarBase.StrKey.isValidEd25519PublicKey(publicKey)
  }

  @action validateAndSetUserKey (secretKey: string, userPubKey: string): boolean {
    if (!this.validateStellarSecretKey(secretKey)) {
      return false
    }
    const keypair = StellarSdk.Keypair.fromSecret(secretKey)
    if (keypair.publicKey() === userPubKey) {
      this.keyVerified = true
      this.keyPair = keypair
      return true
    }
    this.keyVerified = false
    return false
  }

  async mkTradeStageDocTx (id: TradeStageDocPathType, operationType: string, expiresAt?: Date) {
    const response = await this.gqlClient.mutate({
      mutation: mkTradeStageDocTx,
      variables: { 'id': id, 'operationType': operationType, 'expiresAt': expiresAt }
    })
    return response.data.mkTradeStageDocTx
  }

  async mkTradeStageCloseTx (id: TradeStagePathType, operationType: string) {
    const response = await this.gqlClient.mutate({
      mutation: mkTradeStageCloseTx,
      variables: { 'id': id, 'operationType': operationType }
    })
    return response.data.mkTradeStageCloseTx
  }

  async mkTradeStageAddTx (id: TradeStagePathType, operationType: string) {
    const response = await this.gqlClient.mutate({
      mutation: mkTradeStageAddTx,
      variables: { 'id': id, 'operationType': operationType }
    })
    return response.data.mkTradeStageAddTx
  }

  async mkTradeCloseTx (tid: string, operationType: string) {
    const response = await this.gqlClient.mutate({
      mutation: mkTradeCloseTx,
      variables: { 'id': tid, 'operationType': operationType }
    })
    return response.data.mkTradeCloseTx
  }

  async signDocApprovalTx (
    inputValue: TradeStageDocPathType,
    operationType: string,
    expiresAt?: Date) {
    const docTx = await this.mkTradeStageDocTx(inputValue, operationType, expiresAt)
    return this.signRawTX(docTx)
  }

  async signStageCloseTx (inputValue: TradeStagePathType, operationType: string) {
    const docTx = await this.mkTradeStageCloseTx(inputValue, operationType)
    return this.signRawTX(docTx)
  }

  async signStageAddTx (inputValue: TradeStagePathType, operationType: string) {
    const docTx = await this.mkTradeStageAddTx(inputValue, operationType)
    return this.signRawTX(docTx)
  }

  async signTradeCloseTx (tid: string, operationType: string) {
    const docTx = await this.mkTradeCloseTx(tid, operationType)
    return this.signRawTX(docTx)
  }
}

export default new StellarStore()
