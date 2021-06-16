// @flow

import { mkLink } from '../services/cerealia'
import { mkFormDataHeaders, addNotificationHelper } from './helper'
import Axios from 'axios'
import type { LocationType } from '../model/flowType'
import { locationMap } from '../constants/tradeConst'

export async function openDocFile (docID: string, filename: string, location: LocationType) {
  let objUrl
  try {
    let url = location === locationMap.trade ? 'v1/trades/stage-docs/' : 'v1/trade-offers/docs/'
    let response = await downloadFile(url + docID)
    objUrl = URL.createObjectURL(new Blob([response.data]))
    let a = document.createElement('a')
    a.style.display = 'none'
    const body = document.body
    if (body) {
      body.appendChild(a)
    }
    a.download = filename
    a.href = objUrl
    a.click()
  } catch (e) {
    addNotificationHelper(e, 'error')
  } finally {
    window.URL.revokeObjectURL(objUrl)
  }
}

export function downloadFile (url: string) {
  return Axios.get(
    mkLink(url),
    {
      headers: mkFormDataHeaders(),
      responseType: 'arraybuffer'
    })
}
