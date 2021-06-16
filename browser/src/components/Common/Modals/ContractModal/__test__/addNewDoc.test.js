import React from 'react'
import { shallow } from 'enzyme'
import moment from 'moment'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import tradeViewStore from '../../../../../stores/tradeViewStore'
import stellarStore from '../../../../../stores/stellarStore'
import currentUserStore from '../../../../../stores/current-user'
import '../../../../../setupTests'
import
clientMock,
{
  sampleTrade1,
  sampleUserKeyPair1
} from '../../../../../graphql/test/client-mock'
import { approveStatus } from '../../../../../constants/tradeConst'
import AddNewDocumentModal from '../AddNewDocumentModal'
import { mkLink } from '../../../../../services/cerealia'
import RejectTradeStageDocModal from '../RejectTradeStageDocModal'
import ApproveTradeStageDocModal from '../ApproveTradeStageDocModal'

describe('Test to add trade new stage', () => {
  const addDocResponse = { docID: '1' }
  const uploadURL = mkLink('v1/trades/stage-docs')
  const testStageData = {
    selectedFileName: 'pdfFile.pdf',
    expiresAt: moment('2019-09-25T12:23:25.833791Z'),
    selectedFile: new File(['(⌐□_□)'], 'pdfFile.pdf'),
    timeZone: Intl.DateTimeFormat().resolvedOptions().timeZone,
    fileName: 'pdfFile_Official',
    hash: 'aaaaaaaaaaa'
  }

  beforeAll(async function () {
    currentUserStore.updateClient(clientMock)
    stellarStore.updateClient(clientMock)
    await stellarStore.initializeStellarNetwork()
    stellarStore.validateAndSetUserKey(sampleUserKeyPair1.secKey, sampleUserKeyPair1.pubKey)
    const mock = new MockAdapter(axios)
    mock
      .onPost(uploadURL)
      .reply(200, addDocResponse)
  })
  beforeEach(async function () {
    await currentUserStore.login({ email: 'ss@ss.ss', password: 'birthday' })
  })

  it('Add document without approval', async () => {
    tradeViewStore.updateTrade(sampleTrade1, clientMock)
    const addNewDocumentModalWrapper = shallow(
      <AddNewDocumentModal
        onCloseModal={() => {
        }}
        visible
        stageIdx={0} />
    ).instance()
    addNewDocumentModalWrapper.setState({
      ...testStageData,
      withApproval: false
    })
    await addNewDocumentModalWrapper.onAddNewDocument()

    expect(tradeViewStore.stages[0].docs.length).toEqual(1)
    expect(tradeViewStore.stages[0].docs[0].status).toEqual(approveStatus.submitted)
    expect(tradeViewStore.stages[0].docs[0].doc.id).toEqual('1')
    expect(tradeViewStore.stages[0].docs[0].doc.name).toEqual('pdfFile_Official.pdf')
  })

  it('Add document with approval', async () => {
    tradeViewStore.updateTrade(sampleTrade1, clientMock)
    const addNewDocumentModalWrapper = shallow(
      <AddNewDocumentModal
        onCloseModal={() => {
        }}
        visible
        stageIdx={0} />
    ).instance()
    addNewDocumentModalWrapper.setState({
      ...testStageData,
      withApproval: true
    })
    await addNewDocumentModalWrapper.onAddNewDocument()
    expect(tradeViewStore.stages[0].docs.length).toEqual(1)
    expect(tradeViewStore.stages[0].docs[0].status).toEqual(approveStatus.pending)
    expect(tradeViewStore.stages[0].docs[0].doc.id).toEqual('1')
    expect(tradeViewStore.stages[0].docs[0].doc.name).toEqual('pdfFile_Official.pdf')

    // reject test for new document
    let stageDoc = tradeViewStore.stages[0].docs[0]
    const rejectStageDocModalWrapper = shallow(
      <RejectTradeStageDocModal
        visible
        onCloseModal={() => {}}
        stageIdx={0}
        stageDoc={stageDoc} />
    ).instance()
    rejectStageDocModalWrapper.setState({ error: false })
    await rejectStageDocModalWrapper.rejectTradeStageDoc()
    expect(tradeViewStore.stages[0].docs.length).toEqual(1)
    expect(tradeViewStore.stages[0].docs[0].status).toEqual(approveStatus.rejected)
    expect(tradeViewStore.stages[0].docs[0].doc.id).toEqual('1')

    // add new document again for approve test
    addNewDocumentModalWrapper.setState({
      ...testStageData,
      withApproval: true
    })
    await addNewDocumentModalWrapper.onAddNewDocument()
    expect(tradeViewStore.stages[0].docs.length).toEqual(2)
    expect(tradeViewStore.stages[0].docs[1].status).toEqual(approveStatus.pending)

    // approve test for new document
    stageDoc = tradeViewStore.stages[0].docs[1]
    const approveStageDocModalWrapper = shallow(
      <ApproveTradeStageDocModal
        visible
        onCloseModal={() => {}}
        stageIdx={0}
        stageDoc={stageDoc} />
    ).instance()
    await approveStageDocModalWrapper.approveTradeStageDoc()
    expect(tradeViewStore.stages[0].docs.length).toEqual(2)
    expect(tradeViewStore.stages[0].docs[1].status).toEqual(approveStatus.approved)
    expect(tradeViewStore.stages[0].docs[1].doc.id).toEqual('1')
  })
})
