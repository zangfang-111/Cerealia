// @flow

import { action, observable } from 'mobx'

const defaultSpinText = `Blockchain operation submitted. It will take few seconds
         to process and store it in the blockchain`

class SpinnerStore {
  @observable loading: boolean = false
  @observable spinnerTip: string = ''

  @action showSpinner = (spinTip: string = defaultSpinText): void => {
    this.loading = true
    this.spinnerTip = spinTip
  }

  @action hideSpinner = (): void => {
    this.loading = false
  }
}

export default new SpinnerStore()
