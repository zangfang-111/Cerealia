// @flow

import React from 'react'
import { Modal, Spin } from 'antd'
import { observer } from 'mobx-react'
import spinnerStore from '../../../../stores/spinner-store'

const icons = {
  approve: <i className='fas fa-check-circle' />,
  decline: <i className='fas fa-times-circle' />,
  new_stage: <i className='fas fa-plus-circle' />,
  new_document: <i className='fas fa-file-invoice' />,
  rejected_reason: <i className='fas fa-exclamation-circle' />,
  create_trade: <i className='fas fa-file-contract' />,
  approve_additional_stage: <i className='fas fa-check-circle' />,
  reject_additional_stage: <i className='fas fa-exclamation-triangle' />,
  approve_additional_doc: <i className='fas fa-check-circle' />,
  reject_additional_doc: <i className='fas fa-times-circle' />,
  delete_stage: <i className='fas fa-trash' />,
  close_stage: <i className='fas fa-check-circle' />,
  stellar_key: <i className='fas fa-key' />,
  user: <i className='fas fa-user' />,
  history: <i className='fas fa-history' />
}

type Props = {
  children: any,
  visible: boolean,
  onCloseModal: Function,
  status: string,
}

@observer
class BaseModalContainer extends React.Component<Props, {}> {
  render () {
    const { visible, onCloseModal, status } = this.props
    return (
      <Modal
        closable={false}
        title='Basic Modal'
        visible={visible}
        className={'modal-container'}
        wrapClassName='vertical-center-modal'
        align={{}}
      >
        <div className={'modal-content'}>
          <p className={`circle ${status}`}>
            {icons[status]}
          </p>
          <p className={'close-icon'}>
            <i className='fas fa-times-circle hoverable' onClick={onCloseModal} />
          </p>
          <div className={'confirm'}>
            <Spin spinning={spinnerStore.loading} size='large' tip={spinnerStore.spinnerTip}>
              {this.props.children}
            </Spin>
          </div>
        </div>
      </Modal>
    )
  }
}

export default BaseModalContainer
