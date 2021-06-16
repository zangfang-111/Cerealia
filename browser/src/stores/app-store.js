// @flow

import { action, observable } from 'mobx'

class AppStore {
  @observable adminMode: boolean = false
  @observable appTheme: String = 'theme-light'

  constructor () {
    let curTheme = 'theme-light'
    // Set browser cookie for theme
    if (document.cookie.indexOf('cr_app_theme') === -1) {
      // Create theme cookie
      document.cookie = 'cr_app_theme=' + curTheme
    } else {
      curTheme = this.getCookie('cr_app_theme')
      let rootElement = document.getElementById('root')
      rootElement.classList.add(curTheme)
    }
    this.appTheme = curTheme
  }

  // Found logic for fetch cookie on stack overflow
  // https://stackoverflow.com/questions/10730362/get-cookie-by-name
  getCookie = (name: string) => {
    var value = '; ' + document.cookie
    var parts = value.split('; ' + name + '=')
    if (parts.length === 2) return parts.pop().split(';').shift()
  }

  @action setAdminMode = (adminMode: boolean) => {
    this.adminMode = adminMode
  }

  @action setTheme = (addTheme: string) => {
    let removeTheme = addTheme === 'theme-light' ? 'theme-dark' : 'theme-light'
    let rootElement = document.getElementById('root')
    rootElement.classList.remove(removeTheme)
    rootElement.classList.add(addTheme)
    this.appTheme = addTheme
    // update cookie value
    document.cookie = 'cr_app_theme=' + addTheme
  }
}

export default new AppStore()
