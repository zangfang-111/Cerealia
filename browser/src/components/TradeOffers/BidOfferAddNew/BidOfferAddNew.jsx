// @flow

import { Button, DatePicker, Form, Icon, Input, Radio, Select, Switch } from 'antd'
import React from 'react'
import { commodities, commodityTypes, incoterms } from '../../../constants/selectOptions'
import { currencies, dateFormat, timeFormat } from '../../../constants/tradeConst'
import type { OrgMapType } from '../../../model/flowType'
import currentUser from '../../../stores/current-user'
import { RenderDropZone } from '../../Common/RenderDropZone/RenderDropZone'

const { RangePicker } = DatePicker
const Option = Select.Option
const FormItem = Form.Item
const TextArea = Input.TextArea

type Props = {
  onSubmit: Function,
  form: Object,
  offerStatus: boolean
}

type State = {
  selectedFileName: string,
  selectedFile: Object,
  anonymousStatus: boolean,
  offerStatus: boolean
}

class NewTradeOfferForm extends React.Component<Props, State> {
  child: Object
  constructor (props: Props) {
    super(props)

    this.state = {
      selectedFileName: '',
      selectedFile: null,
      anonymousStatus: false,
      offerStatus: props.offerStatus
    }
  }

  handleSubmit = (e: SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault()
    this.props.form.validateFieldsAndScroll((err: any, values: any) => {
      if (!err) {
        values.isSell = values.isSell === 'OFFER'
        values.isAnonymous = this.state.anonymousStatus
        values.note = values.note || ''
        this.props.onSubmit(values, this.state.selectedFile)
      }
    })
  }

  onDropHandler = (accepted: Array<any>, rejected: Array<any>) => {
    if (accepted.length !== 1) {
      this.child.showNotification('User can select only one file')
    } else {
      this.setState({
        selectedFile: accepted[0],
        selectedFileName: accepted[0].name
      })
    }
  }

  onToggleConfirmation = () => {
    this.setState({ anonymousStatus: !this.state.anonymousStatus })
  }

  onSetOfferType = (e: Object) => {
    this.setState({ offerStatus: e.target.value === 'OFFER' })
  }

  goodsSelectOptions = (x: [string, string]) => {
    return (<Option key={x[0]} value={x[0]}>{x[1]}</Option>)
  }

  companySelectOptions = (x: OrgMapType, k: number) => {
    return (<Option key={k} value={x.org.id}>{x.org.name}</Option>)
  }

  currencyIncortemSelectOptions = (c: string) => {
    return (<Option key={c} value={c}>{c}</Option>)
  }

  render () {
    const { getFieldDecorator } = this.props.form
    const rangeConfig = {
      rules: [{ type: 'array', required: true, message: 'Please Shipment time!' }]
    }

    return (
      <Form className={'bid-add-new'}>
        <p className={'bid-offer-title'}>
          {this.state.offerStatus ? 'Offers' : 'Bids'} / Add new
        </p>
        <div className={'page-title'}>Create new</div>
        <div className={'order-type'}>
          <FormItem className={'select-name right-margin'}>
            <p>Keep Anonymous</p>
            <Switch
              className={'switch-btn'}
              checkedChildren={<Icon type={'check'} />}
              unCheckedChildren={<Icon type={'close'} />}
              onChange={this.onToggleConfirmation}
              value={this.state.anonymousStatus}
            />
          </FormItem>
          <FormItem className={'select-name'}>
            <p>Type</p>
            {getFieldDecorator('isSell', {
              rules: [{ required: true, message: 'Offer type is required' }]
            })(
              <Radio.Group size={'large'} required
                onChange={this.onSetOfferType}>
                <Radio value={'BID'}>Bid</Radio>
                <Radio value={'OFFER'}>Offer</Radio>
              </Radio.Group>
            )}
          </FormItem>
        </div>
        <div className={'goods-container'}>
          <div className={'content-title'}>Goods</div>
          <div className={'name-type-qty'}>
            <FormItem className={'select-name'}>
              <p>Commodity</p>
              {getFieldDecorator('commodity', {
                rules: [{ required: true, message: 'Commodity is required' }]
              })(
                <Select
                  className={'select-input'}
                  size={'large'}
                  dropdownClassName={'select-dropdown'}
                  placeholder={'Type something'}
                >
                  {commodities.map(this.goodsSelectOptions)}
                </Select>
              )}
            </FormItem>
            <FormItem className={'select-type'}>
              <p>Type</p>
              {getFieldDecorator('comType', {
                rules: [{ required: true, message: 'Commodity type is required' }]
              })(
                <Select
                  mode='tags'
                  className={'select-input'}
                  size={'large'}
                  placeholder={'Type something'}
                >
                  {commodityTypes.map(this.goodsSelectOptions)}
                </Select>
              )}
            </FormItem>
            <FormItem className={'input-quality'}>
              <p>Quality</p>
              {getFieldDecorator('quality', {
                rules: [{ required: true, message: 'Quality is required' }]
              })(
                <Input placeholder={'Quality'} className={'select-input'}
                  size='large' spellCheck={false} />
              )}
            </FormItem>
          </div>
          <div className={'origin-cpy-expire'}>
            <FormItem className={'select-cover'}>
              <p>Origin</p>
              {getFieldDecorator('origin', {
                rules: [{ required: true, message: 'Origin is required' }]
              })(
                <Input placeholder={'Origin'} className={'select-input'}
                  size='large' spellCheck={false} />
              )}
            </FormItem>
            <FormItem className={'select-cover'}>
              <p>Expires</p>
              {getFieldDecorator('expiresAt')(
                <DatePicker
                  className={'select-input'}
                  size='large'
                  format={timeFormat}
                />
              )}
            </FormItem>
            <FormItem className={'select-cover'}>
              <p>Company</p>
              {getFieldDecorator('orgID', {
                rules: [{ required: true, message: 'Company is required' }]
              })(
                <Select
                  className={'select-input'}
                  size={'large'}
                  dropdownClassName={'select-dropdown'}
                  placeholder={'Type something'}
                >
                  { currentUser.user.orgMap.map(this.companySelectOptions) }
                </Select>
              )}
            </FormItem>
          </div>
        </div>
        <div className={'payment-container'}>
          <div className={'content-title'}>Payment</div>
          <div className={'payment-cover'}>
            <FormItem className={'price'}>
              <p>Price</p>
              {getFieldDecorator('price', {
                rules: [{ required: true, message: 'Price is required' }]
              })(
                <Input
                  className={'select-input'}
                  size='large'
                  type='number'
                  min='1'
                  step='0.01'
                  required
                />
              )}
            </FormItem>
            <FormItem className={'price right-margin'}>
              <p>Currency</p>
              {getFieldDecorator('currency', {
                rules: [{ required: true, message: 'Currency is required' }]
              })(
                <Select
                  className={'select-input'}
                  size={'large'}
                  dropdownClassName={'select-dropdown'}
                >
                  {currencies.map(this.currencyIncortemSelectOptions)}
                </Select>
              )}
            </FormItem>
            <FormItem className={'status-cover'}>
              <p>Price Type</p>
              {getFieldDecorator('priceType', {
                rules: [{ required: true, message: 'Status is required', whitespace: true }]
              })(
                <Radio.Group size={'large'} required >
                  <Radio value={'firm'}>Firm</Radio>
                  <Radio value={'quote'}>Quote</Radio>
                </Radio.Group>
              )}
            </FormItem>
          </div>
          <FormItem className={'payment-terms'}>
            <p>Payment terms</p>
            {getFieldDecorator('note')(
              <TextArea rows={6} />
            )}
          </FormItem>
        </div>
        <div className={'basis-container'}>
          <div className={'content-title'}>Basis</div>
          <div className={'select-section'}>
            <FormItem className={'select-cover'}>
              <p>Incoterm</p>
              {getFieldDecorator('incoterm', {
                rules: [{ required: true, message: 'Incoterm is required' }]
              })(
                <Select
                  className={'select-input'}
                  size={'large'}
                  dropdownClassName={'select-dropdown'}
                >
                  {incoterms.map(this.currencyIncortemSelectOptions)}
                </Select>
              )}
            </FormItem>
            <FormItem className={'select-cover'}>
              <p>Location</p>
              {getFieldDecorator('marketLoc', {
                rules: [{ required: true, message: 'Location is required' }]
              })(
                <Input placeholder={'Location'} className={'select-input'}
                  size='large' spellCheck={false} />
              )}
            </FormItem>
          </div>
        </div>
        <div className={'qty-smt-status'}>
          <FormItem className={'quantity-cover'}>
            <div className={'content-title'}>Quantity</div>
            <div className={'quantity'}>
              {getFieldDecorator('vol', {
                rules: [{ required: true, message: 'Quantity is required', whitespace: true }]
              })(
                <Input placeholder={'Quantity'} className={'quanty-input'}
                  size='large' type='number' spellCheck={false} />
              )}
              <p className={'quantity-unit'}>t</p>
            </div>
          </FormItem>
          <Form.Item className={'shipment-cover'} >
            <div className={'content-title'}>Shipment</div>
            <Form.Item
              style={{ display: 'inline-block' }}
            >
              {getFieldDecorator('shipment', rangeConfig)(
                <RangePicker
                  size='large'
                  format={dateFormat}
                />
              )}
            </Form.Item>
          </Form.Item>
        </div>
        <FormItem className={'terms-conditions'}>
          <div className={'content-title'}>Terms & Conditions</div>
          {RenderDropZone(this.onDropHandler, this.state.selectedFileName)}
        </FormItem>
        <Button
          className={'submit-btn'}
          type={'primary'}
          onClick={this.handleSubmit}
        >
          {`Add an ${this.state.offerStatus ? 'Offer' : 'Bid'}`}
        </Button>
      </Form>
    )
  }
}

export default NewTradeOfferForm
