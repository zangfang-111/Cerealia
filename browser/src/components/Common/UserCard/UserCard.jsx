// @flow

import React from 'react'
import { mkPublicLink } from '../../../services/cerealia'

type Props = {
  actor: string,
  name?: string,
  orgName: string,
  role: string,
  avatar: string,
};

export default function (props: Props) {
  const { actor, name, orgName, role, avatar } = props

  const showAvatar = (avatar) => {
    if (avatar) {
      return React.createElement('img', {
        src: mkPublicLink('/user/avatars/' + avatar),
        className: 'card-image'
      })
    }
    return React.createElement('i', { className: 'fas fa-user-circle avatar-icon' })
  }

  return (
    <div className={'user-card'}>
      <div className={'card-header'}>
        <div>{ showAvatar(avatar) }</div>
        <div className={'user-name'}>
          <p>{actor}</p>
          <p>{name}</p>
          <p>{orgName} - {role}</p>
        </div>
      </div>
    </div>
  )
}
