// @flow

export function isEmpty (obj: any) {
  obj = obj || obj
  return [Object, Array].includes(obj.constructor) && !Object.entries(obj).length
}

export default {
  isEmpty
}
