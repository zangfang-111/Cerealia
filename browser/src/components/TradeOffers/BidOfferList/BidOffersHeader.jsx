// @flow

import React from 'react'
import { Select, Icon } from 'antd'
import { observer } from 'mobx-react-lite'
import useReactRouter from 'use-react-router'
import tradeOfferStore from '../../../stores/tradeoffer-store'

const Option = Select.Option

export default observer(() => {
  const { history } = useReactRouter()
  function handleCreateNew () { history.push('/trade-offer/new') }
  return (
    <div className={'bid-header'}>
      <div className={'header-left'}>
        <div className={'btn-group'} onClick={handleCreateNew}>
          <p>{`Create New ${tradeOfferStore.offerStatus ? 'Offers' : 'Bids'}`}</p>
          <i className='fas fa-plus-circle' />
        </div>
        <div className={'page-title'}>{tradeOfferStore.offerStatus ? 'Offers' : 'Bids'}</div>
      </div>
      <div className={'header-right'}>
        <Icon type='filter' theme='filled' className={'filter-icon'} />
        <div>
          <p>Product origin</p>
          <Select
            className={'select-origin'}
            size={'large'}
            dropdownClassName={'select-dropdown'}
            placeholder={'Type something'}
          >
            <Option value={'TEST1'}>test1</Option>
            <Option value={'TEST2'}>test2</Option>
            <Option value={'TEST3'}>test3</Option>
          </Select>
        </div>
        <div>
          <p>Product commodity</p>
          <Select
            className={'select-commodity'}
            size={'large'}
            dropdownClassName={'select-dropdown'}
            placeholder={'Type something'}
          >
            <Option value={'TEST1'}>test1</Option>
            <Option value={'TEST2'}>test2</Option>
            <Option value={'TEST3'}>test3</Option>
          </Select>
        </div>
        <div>
          <p>Shipment date</p>
          <Select
            className={'select-shipment'}
            size={'large'}
            dropdownClassName={'select-dropdown'}
            placeholder={'2019-02-15'}
          >
            <Option value={'TEST1'}>test1</Option>
            <Option value={'TEST2'}>test2</Option>
            <Option value={'TEST3'}>test3</Option>
          </Select>
        </div>
      </div>
    </div>
  )
})
