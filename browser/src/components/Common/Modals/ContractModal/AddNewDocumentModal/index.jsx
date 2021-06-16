// @flow

import 'moment-timezone'

import { DatePicker, Icon, Switch } from 'antd'
import axios from 'axios'
import moment from 'moment'
import React from 'react'
import TimezonePicker from 'react-timezone'
import { invalidDateMessage } from '../../../../../constants/errors'
import { addNotificationHelper, fileHash, mkFormDataHeaders } from '../../../../../lib/helper'
import { approveStatus, timeFormat } from '../../../../../constants/tradeConst'
import { mkLink } from '../../../../../services/cerealia'
import spinnerStore from '../../../../../stores/spinner-store'
import stellarStore from '../../../../../stores/stellarStore'
import tradeViewStore from '../../../../../stores/tradeViewStore'
import currentUser from '../../../../../stores/current-user'
import Button from '../../../Button/Button'
import PublicKeyConfirmation from '../../../Collapse/PublicKeyConfirmation/PublicKeyConfirmation'
import { RenderDropZone } from '../../../RenderDropZone/RenderDropZone'
import BaseModal from '../BaseModal'
import type { StageModalPropsType } from '../types'

const path = require('path')

type State = {
  error: boolean,
  selectedFileName: string,
  expiresAt: moment,
  selectedFile: Object,
  hash: string,
  loading: boolean,
  isOpenCollapse: boolean,
  fileName: string,
  timeZone: ?string,
  withApproval: boolean
}

class AddNewDocumentModal extends React.Component<StageModalPropsType, State> {
  documentDescription: string
  child: Object
  constructor (props: StageModalPropsType) {
    super(props)

    this.documentDescription = ''
    this.state = {
      error: false,
      selectedFileName: '',
      expiresAt: moment(),
      selectedFile: {},
      timeZone: Intl.DateTimeFormat().resolvedOptions().timeZone,
      fileName: '',
      hash: '',
      isOpenCollapse: false,
      loading: false,
      withApproval: true
    }
  }

  confirmPublicKey = () => this.setState({ isOpenCollapse: true })

  onCancel = () => this.setState({ isOpenCollapse: false },
    () => this.props.onCloseModal())

  modalclear = () => {
    this.setState({
      selectedFileName: '',
      selectedFile: {},
      fileName: '',
      expiresAt: moment(),
      timeZone: Intl.DateTimeFormat().resolvedOptions().timeZone
    })
    this.documentDescription = ''
  }

  onDropHandler = (accepted: Array<any>, rejected: Array<any>) => {
    if (accepted.length !== 1) {
      addNotificationHelper('User can select only one file', 'error')
    } else {
      fileHash(accepted[0], (x) => {
        this.setState({
          selectedFileName: accepted[0].name,
          selectedFile: accepted[0],
          fileName: path.basename(accepted[0].name, path.extname(accepted[0].name)),
          hash: x
        })
      })
    }
  }

  selectTimeZone = (timeZone: string) =>
    this.setState((state) => {
      return {
        timeZone: timeZone,
        expiresAt: state.expiresAt.tz(timeZone) }
    })

  checkDocName = (officialDocName: string) => {
    let errMessage = ''
    if (!officialDocName) {
      errMessage = 'Official FileName cannot be empty'
    }
    if (officialDocName.indexOf('.') === 0) {
      errMessage = 'Official FileName must not be started from "."'
    }
    if (officialDocName.indexOf('/') >= 0) {
      errMessage = 'Official FileName must not include "/"'
    }
    return errMessage
  }

  selectDate = (e: Date) => this.setState({ expiresAt: e })

  onFilenameChange = (e: SyntheticInputEvent<HTMLInputElement>) => {
    this.setState({ fileName: e.currentTarget.value })
  }

  calculateNextDocumentIDx = (stageIdx: number) => {
    return tradeViewStore.stages[stageIdx].docs
      ? tradeViewStore.stages[stageIdx].docs.length
      : 0
  }

  createNewDocInputParams = async (fullOfficialDocName: string, expiresAt: Date) => {
    const { stageIdx } = this.props
    const { withApproval, hash } = this.state
    let inputValue = {
      tid: tradeViewStore.id,
      stageIdx: stageIdx,
      stageDocIdx: this.calculateNextDocumentIDx(stageIdx),
      stageDocHash: hash
    }
    let opType = withApproval ? approveStatus.pending : approveStatus.submitted
    const signedTx = await stellarStore.signDocApprovalTx(
      inputValue,
      opType,
      moment.utc(expiresAt).format())
    return {
      tid: tradeViewStore.id,
      note: this.documentDescription,
      fileName: fullOfficialDocName,
      stageIdx: stageIdx,
      expiresAt: moment.utc(expiresAt).format(),
      hash: hash,
      signedTx: signedTx,
      docID: '',
      withApproval: withApproval
    }
  }

  onAddNewDocument = async () => {
    let officialDocName = this.state.fileName.replace(/\s/g, '_')
    let fileExt = path.extname(this.state.selectedFileName)
    let fullOfficialDocName = officialDocName.concat(fileExt)
    let docNameError = this.checkDocName(officialDocName)
    if (docNameError) {
      addNotificationHelper(docNameError, 'error')
      return
    }
    const expiresAt = this.state.expiresAt.tz(Intl.DateTimeFormat().resolvedOptions().timeZone)
    if (moment(this.state.expiresAt) < moment() && this.state.withApproval) {
      this.setState({ error: true })
      return
    }
    spinnerStore.showSpinner()
    let inputParams
    try {
      inputParams = await this.createNewDocInputParams(fullOfficialDocName, expiresAt)
    } catch (err) {
      spinnerStore.hideSpinner()
      addNotificationHelper(err, 'error')
      return
    }
    this.setState({ error: false })
    let data = new FormData()
    data.append('formfile', this.state.selectedFile, fullOfficialDocName)
    try {
      let response = await axios.post(
        mkLink('v1/trades/stage-docs'),
        data,
        {
          headers: mkFormDataHeaders(),
          params: inputParams
        })
      inputParams.docID = response.data.docID
      tradeViewStore.createStageDocCallback(inputParams, currentUser.user)
      this.props.onCloseModal()
      this.modalclear()
    } catch (err) {
      addNotificationHelper(err, 'error')
    }
    spinnerStore.hideSpinner()
  }

  onToggleConfirmation = () => {
    this.setState({ withApproval: !this.state.withApproval })
  }

  render () {
    return (
      <BaseModal visible={this.props.visible}
        onCloseModal={this.props.onCloseModal} status={'new_document'}
        spinnerStore={spinnerStore}>
        <div className={'add-new-document'}>
          <p className={'modal-title'}>Submit New Document</p>
          {RenderDropZone(this.onDropHandler, this.state.selectedFileName)}
          <input type={'text'} placeholder={'Rename proper document name...'}
            size='large' value={this.state.fileName} onChange={this.onFilenameChange} />
          <p>
            This filename will be used officially
            in the stage view. Please consider using the right file names
          </p>
          <textarea className={'description'} placeholder={'Type some description...'}
            ref={input => (this.documentDescription = input ? input.value : '')} />

          <div className={'confirm-approve'}>
            <Switch
              className={'switch-btn'}
              checkedChildren={<Icon type={'check'} />}
              unCheckedChildren={<Icon type={'close'} />}
              onChange={this.onToggleConfirmation}
              defaultChecked
            />
            <span>Counterparty approve required</span>
          </div>
          { this.state.withApproval &&
            <div className={'date-section'}>
              <p className={'first-part'}>
                <span>Approval expire time</span>
              </p>
              <div className={'second-part'}>
                <DatePicker
                  showTime
                  format={timeFormat}
                  onChange={this.selectDate}
                  value={this.state.expiresAt}
                />
                <TimezonePicker
                  value={this.state.timeZone}
                  onChange={this.selectTimeZone}
                  inputProps={{
                    placeholder: 'Select Timezone...',
                    name: 'timezone'
                  }}
                />
              </div>
              {
                this.state.error && <p className={'error'}>{invalidDateMessage}</p>
              }
            </div>
          }
          <PublicKeyConfirmation
            isOpenCollapse={this.state.isOpenCollapse}
            publicKey={currentUser.pubKey}
            onCloseModal={this.onCancel}
            confirm={this.onAddNewDocument}
            keyVerified={stellarStore.keyVerified}
          />
          {
            !this.state.isOpenCollapse &&
            <Button type={'primary'} text={'Submit'} onClick={this.confirmPublicKey} />
          }
        </div>
      </BaseModal>
    )
  }
}

export default AddNewDocumentModal
