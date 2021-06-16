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
import { approveStatus, canMakeReqStatus } from '../../../../../constants/tradeConst'
import AddNewDocumentModal from '../AddNewDocumentModal'
import { mkLink } from '../../../../../services/cerealia'
import CloseStageModal from '../CloseStageModal'
import ConfirmClosingStageModal from '../ConfirmClosingStageModal'

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

  it('Stage close request approval test', async () => {
    tradeViewStore.updateTrade(sampleTrade1, clientMock)

    // add doc in the first stage 'trade contract'
    const addNewDocumentModalWrapper = shallow(
      <AddNewDocumentModal
        onCloseModal={() => {}}
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

    // test for stage close request
    const stageCloseModalWrapper = shallow(
      <CloseStageModal
        visible
        onCloseModal={() => {}}
        stageIdx={0} />
    ).instance()
    stageCloseModalWrapper.setState({ error: false })
    await stageCloseModalWrapper.closeStage()
    expect(tradeViewStore.stages[0].closeReqs.length).toEqual(1)
    expect(tradeViewStore.stages[0].closeReqs[0].status).toEqual(approveStatus.pending)
    expect(tradeViewStore.stages[0].stageCloseStatus).toEqual(canMakeReqStatus.pending)
    // test for reject close stage request
    const confirmCloseStageModalWrapper = shallow(
      <ConfirmClosingStageModal
        visible
        onCloseModal={() => {}}
        stageIdx={0} />
    ).instance()
    confirmCloseStageModalWrapper.setState({ error: false })
    await confirmCloseStageModalWrapper.rejectClosingStage()
    expect(tradeViewStore.stages[0].closeReqs.length).toEqual(1)
    expect(tradeViewStore.stages[0].closeReqs[0].status).toEqual(approveStatus.rejected)
    expect(tradeViewStore.stages[0].stageCloseStatus).toEqual(canMakeReqStatus.can)

    // test for approve close stage request
    await stageCloseModalWrapper.closeStage()
    expect(tradeViewStore.stages[0].closeReqs.length).toEqual(2)
    expect(tradeViewStore.stages[0].closeReqs[1].status).toEqual(approveStatus.pending)
    expect(tradeViewStore.stages[0].stageCloseStatus).toEqual(canMakeReqStatus.pending)

    await confirmCloseStageModalWrapper.approveClosingStage()
    expect(tradeViewStore.stages[0].closeReqs.length).toEqual(2)
    expect(tradeViewStore.stages[0].closeReqs[1].status).toEqual(approveStatus.approved)
    expect(tradeViewStore.stages[0].stageCloseStatus).toEqual(canMakeReqStatus.approved)
  })
})
