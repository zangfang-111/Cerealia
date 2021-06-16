import React from 'react'
import { shallow } from 'enzyme'
import tradeViewStore from '../../../../../stores/tradeViewStore'
import stellarStore from '../../../../../stores/stellarStore'
import '../../../../../setupTests'
import
clientMock,
{
  sampleTrade1,
  sampleUserKeyPair1
} from '../../../../../graphql/test/client-mock'
import AddNewStageModal from '../AddNewStageModal'
import { approveStatus, buyerActor } from '../../../../../constants/tradeConst'
import ApproveTradeStageReqModal from '../ApproveTradeStageReqModal'
import RejectTradeStageReqModal from '../RejectTradeStageReqModal'

describe('Test to add trade new stage', () => {
  beforeAll(async function () {
    stellarStore.updateClient(clientMock)
    await stellarStore.initializeStellarNetwork()
    stellarStore.validateAndSetUserKey(sampleUserKeyPair1.secKey, sampleUserKeyPair1.pubKey)
  })
  it('Add stage without approval', async () => {
    tradeViewStore.updateTrade(sampleTrade1, clientMock)
    const addStageModalWrapper = shallow(
      <AddNewStageModal
        onCloseModal={() => {}}
        visible
      />
    ).instance()
    await addStageModalWrapper.createNewStage({
      tid: tradeViewStore.id,
      name: 'new Stage2',
      description: 'This is test stage description',
      reason: 'Require vessel info',
      owner: buyerActor
    })
    expect(tradeViewStore.stageAddReqs.length).toEqual(1)
    expect(tradeViewStore.stageAddReqs[0].status).toEqual(approveStatus.approved)
    expect(tradeViewStore.stageAddReqs[0].owner).toEqual(buyerActor)
    expect(tradeViewStore.stageAddReqs[0].name).toEqual('new Stage2')
    expect(tradeViewStore.stageAddReqs[0].reqBy.id).toEqual('1')

    expect(tradeViewStore.stages.length).toEqual(3)
    expect(tradeViewStore.stages[2].name).toEqual('new Stage2')
    expect(tradeViewStore.stages[0].owner).toEqual(buyerActor)
  })

  it('Add stage request with rejection', async () => {
    tradeViewStore.updateTrade(sampleTrade1, clientMock)
    const addStageModalWrapper = shallow(
      <AddNewStageModal
        onCloseModal={() => {}}
        visible
      />
    ).instance()
    addStageModalWrapper.onToggleConfirmation()
    await addStageModalWrapper.createNewStage({
      tid: tradeViewStore.id,
      name: 'new Stage2',
      description: 'This is test stage description',
      reason: 'Require vessel info',
      owner: buyerActor
    })
    expect(tradeViewStore.stageAddReqs.length).toEqual(1)
    expect(tradeViewStore.stageAddReqs[0].status).toEqual(approveStatus.pending)
    expect(tradeViewStore.stageAddReqs[0].owner).toEqual(buyerActor)
    expect(tradeViewStore.stageAddReqs[0].name).toEqual('new Stage2')

    const addStageRejectModalWrapper = shallow(
      <RejectTradeStageReqModal
        onCloseModal={() => {}}
        stageIdx={0}
        visible
      />
    ).instance()
    let rejectReason = 'I can not agree this stage addition'
    addStageRejectModalWrapper.onChangeTextArea(rejectReason, false)
    await addStageRejectModalWrapper.rejectTradeStageReq()
    expect(tradeViewStore.stageAddReqs[0].status).toEqual(approveStatus.rejected)
    expect(tradeViewStore.stageAddReqs[0].rejectReason).toEqual(rejectReason)
    expect(tradeViewStore.stages.length).toEqual(2)
  })

  it('Add stage with approval test', async () => {
    tradeViewStore.updateTrade(sampleTrade1, clientMock)
    const addStageModalWrapper = shallow(
      <AddNewStageModal
        onCloseModal={() => {}}
        visible
      />
    ).instance()
    addStageModalWrapper.onToggleConfirmation()
    await addStageModalWrapper.createNewStage({
      tid: tradeViewStore.id,
      name: 'new Stage2',
      description: 'This is test stage description',
      reason: 'Require vessel info',
      owner: buyerActor
    })
    expect(tradeViewStore.stageAddReqs.length).toEqual(1)
    expect(tradeViewStore.stageAddReqs[0].status).toEqual(approveStatus.pending)
    expect(tradeViewStore.stageAddReqs[0].owner).toEqual(buyerActor)
    expect(tradeViewStore.stageAddReqs[0].name).toEqual('new Stage2')

    const addStageApproveModalWrapper = shallow(
      <ApproveTradeStageReqModal
        onCloseModal={() => {}}
        stageIdx={0}
        visible
      />
    ).instance()
    await addStageApproveModalWrapper.approveTradeStageReq()

    expect(tradeViewStore.stageAddReqs[0].status).toEqual(approveStatus.approved)

    expect(tradeViewStore.stages.length).toEqual(3)
    expect(tradeViewStore.stages[2].name).toEqual('new Stage2')
    expect(tradeViewStore.stages[0].owner).toEqual(buyerActor)
  })
})
