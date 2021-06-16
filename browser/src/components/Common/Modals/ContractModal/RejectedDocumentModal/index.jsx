// @flow

import React from 'react'
import { toLocalTime } from '../../../../../lib/helper'
import type { TradeStageDocType } from '../../../../../model/flowType'
import tradeViewStore from '../../../../../stores/tradeViewStore'
import usersStore from '../../../../../stores/user-store'
import BaseModal from '../BaseModal'
import type { StageModalPropsType } from '../types'

type Props = StageModalPropsType & {
  stageDoc: TradeStageDocType
}

class RejectedDocumentModal extends React.Component<Props> {
  getRejector = () => {
    const { stageDoc } = this.props
    let rejectorID = stageDoc.approvedBy.id
    let user = usersStore.getUserByID(rejectorID)
    return user ? user.firstName.concat(' ', user.lastName) : ''
  }

  render () {
    const { visible, onCloseModal, stageIdx, stageDoc } = this.props
    return (
      <BaseModal visible={visible} onCloseModal={onCloseModal}
        status={'reject_additional_doc'}>
        {
          stageDoc.doc &&
          <div className={'rejected-document-reason'}>
            <p className={'modal-title'}>Rejected Document</p>
            <p>This document has been rejected by {this.getRejector()}.</p>
            <p>Reason:</p>
            <p>
              {
                tradeViewStore.stages[stageIdx].docs[stageDoc.index].rejectReason
              }
            </p>
            <p>Rejected: {
              toLocalTime(tradeViewStore.stages[stageIdx].docs[stageDoc.index].approvedAt)
            }</p>
          </div>
        }
      </BaseModal>
    )
  }
}

export default RejectedDocumentModal
