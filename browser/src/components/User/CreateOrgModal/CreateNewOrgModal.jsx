// @flow

import React from 'react'
import { Form } from 'antd'
import BaseModal from '../../Common/Modals/ContractModal/BaseModal'
import CreateOrgForm from './CreateOrgForm'
import type { OrgInputType } from '../../../model/flowType'
import { addNotificationHelper } from '../../../lib/helper'
import usersStore from '../../../stores/user-store'

type Props = {
  visible: boolean,
  curOrgName: string,
  onCloseModal: Function,
  onAddNewOrg: Function
}
export default (props: Props) => {
  async function onCreateOrganization (input: OrgInputType) {
    try {
      let org = await usersStore.createOrganization(input)
      props.onAddNewOrg(org)
      props.onCloseModal()
    } catch (err) {
      addNotificationHelper(err, 'error')
    }
  }

  const AddNewOrgForm = Form.create()(CreateOrgForm)
  return (
    <BaseModal visible={props.visible} onCloseModal={props.onCloseModal}
      status={'new_stage'} >
      <AddNewOrgForm curOrgName={props.curOrgName} onCloseModal={props.onCloseModal}
        createOrganization={onCreateOrganization} />
    </BaseModal>
  )
}
