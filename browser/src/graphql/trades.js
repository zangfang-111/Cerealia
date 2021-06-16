// @flow

import gql from 'graphql-tag'

const orgFragment = gql`
  fragment organization on Organization {
    id
    name
    telephone
    email
    address
  }
`

const userFragment = gql`
  fragment user on User {
    id
    firstName
    lastName
    emails
    orgMap{
      org{
        ...organization
      }
      role
    }
    roles
    biography
    avatar
    pubKey
    createdAt
  }
  ${orgFragment}
`

const authUserFragment = gql`
  fragment authUser on AuthUser {
    token
  }
`

const docFragment = gql`
  fragment doc on Doc {
    id
    name
    hash
    note
    url
    type
    createdBy{
      ...user
    }
    createdAt
  }
  ${userFragment}
`

const stageDocFragment = gql`
  fragment tradeStageDoc on TradeStageDoc {
    doc {
      ...doc
    }
    status
    approvedBy{
      ...user
    }
    approvedAt
    rejectReason
    expiresAt
  }
  ${userFragment}
  ${docFragment}
`

const approveReqFragment = gql`
  fragment approveReq on ApproveReq {
    reqBy{
      ...user
    }
    reqAt
    reqActor
    reqReason
    rejectReason
    approvedBy{
      ...user
    }
    approvedAt
    status
  }
  ${userFragment}
`

const stageAddReqFragment = gql`
  fragment tradeStageAddReq on TradeStageAddReq {
    name
    description
    owner
    reqBy{
      ...user
    }
    reqAt
    reqActor
    reqReason
    rejectReason
    approvedBy{
      ...user
    }
    approvedAt
    status
  }
  ${userFragment}
`

const moderatorFragment = gql`
  fragment stageModerator on StageModerator {
    user{
      ...user
    }
    createdAt
  }
  ${userFragment}
`

const stageFragment = gql`
  fragment tradeStage on TradeStage {
    name
    description
    addReqIdx
    docs{
      ...tradeStageDoc
    }
    owner
    expiresAt
    delReqs{
      ...approveReq
    }
    closeReqs{
      ...approveReq
    }
    moderator{
      ...stageModerator
    }
  }
  ${stageDocFragment}
  ${approveReqFragment}
  ${moderatorFragment}
`

const tradeOfferFragment = gql`
  fragment tradeOffer on TradeOffer {
    id
    price
    priceType
    isSell
    currency
    createdBy{
      ...user
    }
    createdAt
    expiresAt
    closedAt
    org{
      ...organization
    }
    isAnonymous
    commodity
    comType
    quality
    origin
    incoterm
    marketLoc
    vol
    shipment
    note
    terms{
      ...doc
    }
  }
  ${userFragment}
  ${orgFragment}
  ${docFragment}
`

const tradeFragment = gql`
  fragment trade on Trade {
    id
    name
    description
    buyer{
      ...user
    }
    seller{
      ...user
    }
    stages{
      ...tradeStage
    }
    stageAddReqs{
      ...tradeStageAddReq
    }
    createdBy{
      ...user
    }
    closeReqs {
      ...approveReq
    }
    createdAt
    scAddr
    moderating
    tradeOffer{
      ...tradeOffer
    }
  }
  ${userFragment}
  ${stageFragment}
  ${stageAddReqFragment}
  ${tradeOfferFragment}
`

const notificationFragment = gql`
  fragment notification on Notification {
    id
    createdAt
    triggeredBy {
      ...user
    }
    receiver
    type
    dismissed
    entityID
    msg
    action
  }
  ${userFragment}
`

export const getTradeData = gql`
  query trades {
    trades{
      ...trade
    }
  }
  ${tradeFragment}
`

export const getTradeOffers = gql`
  query tradeOffers {
    tradeOffers{
      ...tradeOffer
    }
  }
  ${tradeOfferFragment}
`

export const getAdminTradeData = gql`
  query adminTrades {
    adminTrades{
      ...trade
    }
  }
  ${tradeFragment}
`

export const getTemplates = gql`
  query getTemplates {
    tradeTemplates {
      id
      name
      description
    }
  }
`

export const createTrade = gql`
  mutation createTrade($input: NewTradeInput!){
    tradeCreate(input: $input){
        ...trade
    }
  }
  ${tradeFragment}
`

export const createStage = gql`
  mutation createStage($input: NewStageInput!, $signedTx: String!, $withApproval: Boolean!) {
    tradeStageAddReq(input: $input, signedTx: $signedTx, withApproval: $withApproval){
      ...tradeStageAddReq
    }
  }
  ${stageAddReqFragment}
`

export const addStageApprove = gql`
  mutation addStageApprove($id: TradeStagePath!, $signedTx: String!) {
    tradeStageAddReqApprove(id: $id, signedTx: $signedTx){
      ...tradeStage
    }
  }
  ${stageFragment}
`

export const addStageReject = gql`
  mutation addStageReject($id: TradeStagePath!, $signedTx: String!, $reason: String!) {
    tradeStageAddReqReject(id: $id, signedTx: $signedTx, reason: $reason)
  }
`

export const tradeStageDelReq = gql`
  mutation tradeStageDelReq($id: TradeStagePath!, $reason: String!) {
    tradeStageDelReq(id: $id, reason: $reason){
      ...approveReq
    }
  }
  ${approveReqFragment}
`

export const tradeStageDelReqApprove = gql`
  mutation tradeStageDelReqApprove($id: TradeStagePath!) {
    tradeStageDelReqApprove(id: $id){
      ...approveReq
    }
  }
  ${approveReqFragment}
`

export const tradeStageDelReqReject = gql`
  mutation tradeStageDelReqReject($id: TradeStagePath!, $reason: String!) {
    tradeStageDelReqReject(id: $id, reason: $reason)
  }
`

export const tradeStageCloseReq = gql`
  mutation tradeStageCloseReq($id: TradeStagePath!, $reason: String!, $signedTx: String!) {
    tradeStageCloseReq(id: $id, reason: $reason, signedTx: $signedTx){
      ...approveReq
    }
  }
  ${approveReqFragment}
`

export const tradeStageCloseReqApprove = gql`
  mutation tradeStageCloseReqApprove($id: TradeStagePath!, $signedTx: String!) {
    tradeStageCloseReqApprove(id: $id, signedTx: $signedTx){
      ...approveReq
    }
  }
  ${approveReqFragment}
`

export const tradeStageCloseReqReject = gql`
  mutation
  tradeStageCloseReqReject($id: TradeStagePath!, $reason: String!, $signedTx: String!) {
    tradeStageCloseReqReject(id: $id, reason: $reason, signedTx: $signedTx)
  }
`

export const tradeCloseReq = gql`
  mutation tradeCloseReq($id: String!, $reason: String!, $signedTx: String!) {
    tradeCloseReq(id: $id, reason: $reason, signedTx: $signedTx){
      ...approveReq
    }
  }
  ${approveReqFragment}
`

export const tradeCloseReqApprove = gql`
  mutation tradeCloseReqApprove($id: String!, $signedTx: String!) {
    tradeCloseReqApprove(id: $id, signedTx: $signedTx){
      ...approveReq
    }
  }
  ${approveReqFragment}
`

export const tradeCloseReqReject = gql`
  mutation
  tradeCloseReqReject($id: String!, $reason: String!, $signedTx: String!) {
    tradeCloseReqReject(id: $id, reason: $reason, signedTx: $signedTx)
  }
`

export const tradeStageDocApprove = gql`
  mutation tradeStageDocApprove($id: TradeStageDocPath!, $signedTx: String!) {
    tradeStageDocApprove(id: $id, signedTx: $signedTx){
      ...tradeStageDoc
    }
  }
  ${stageDocFragment}
`

export const tradeStageDocReject = gql`
  mutation tradeStageDocReject($id: TradeStageDocPath!, $signedTx: String!, $reason: String!) {
    tradeStageDocReject(id: $id, signedTx: $signedTx, reason: $reason)
  }
`

export const tradeStageSetExpireTime = gql`
  mutation tradeStageSetExpireTime($id: TradeStagePath!, $expiresAt: String!){
    tradeStageSetExpireTime(id: $id, expiresAt: $expiresAt)
  }
`

export const notificationDismiss = gql`
  mutation notificationDismiss($id: String!){
    notificationDismiss(id: $id)
  }
`

export const mkTradeStageDocTx = gql`
  mutation mkTradeStageDocTx(
    $id: TradeStageDocPath!,
    $operationType: Approval!,
    $expiresAt: Time){
    mkTradeStageDocTx(id: $id, operationType: $operationType, expiresAt: $expiresAt)
  }
`

export const mkTradeStageCloseTx = gql`
  mutation mkTradeStageCloseTx($id: TradeStagePath!, $operationType: Approval!){
    mkTradeStageCloseTx(id: $id, operationType: $operationType)
  }
`

export const mkTradeStageAddTx = gql`
  mutation mkTradeStageAddTx($id: TradeStagePath!, $operationType: Approval!){
    mkTradeStageAddTx(id: $id, operationType: $operationType)
  }
`

export const mkTradeCloseTx = gql`
  mutation mkTradeCloseTx($id: String!, $operationType: Approval!){
    mkTradeCloseTx(id: $id, operationType: $operationType)
  }
`

export const getCurrentUser = gql`
  query user{
    user{
      ...user
    }
  }
  ${userFragment}
`

export const createTradeOffer = gql`
  mutation tradeOfferCreate($input: TradeOfferInput!){
    tradeOfferCreate(input: $input) {
      ...tradeOffer
    }
  }
  ${tradeOfferFragment}
`

export const createOrganization = gql`
  mutation organizationCreate($input: OrgInput!){
    organizationCreate(input: $input) {
      ...organization
    }
  }
  ${orgFragment}
`

export const userSignup = gql`
  mutation userSignup($input: NewUserInput) {
    userSignup(input: $input)
  }
`

export const userLogin = gql`
  mutation userLogin($input: UserLoginInput!){
    userLogin(input: $input){
      ...authUser
    }
  }
  ${authUserFragment}
`

export const getStellarInfo = gql`
  query stellarNet{
    stellarNet{
      url
      passphrase
      name
    }
  }
`

export const getAllUsers = gql`
  query allUsers{
    users{
      ...user
    }
  }
  ${userFragment}
`

export const getAllOrganizations = gql`
  query allOrganizations{
    organizations{
      ...organization
    }
  }
  ${orgFragment}
`

export const changePassword = gql`
  mutation userPasswordChange($input: ChangePasswordInput!){
    userPasswordChange(input: $input)
  }
`

export const changeEmail = gql`
  mutation userEmailChange($input: [Email!]!){
    userEmailChange(input: $input)
  }
`

export const updateUserProfile = gql`
  mutation userProfileUpdate($input: UserProfileInput!){
    userProfileUpdate(input: $input){
      ...user
    }
  }
  ${userFragment}
`

export const notifications = gql`
  query notifications($from: Uint!){
    notifications(from: $from){
      ...notification
    }
  }
  ${notificationFragment}
`

export const notificationsTrade = gql`
  query notificationsTrade($id: String!){
    notificationsTrade(id: $id){
      ...notification
    }
  }
  ${notificationFragment}
`
