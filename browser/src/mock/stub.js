// @flow

/* stub object builder
   example:
     let x = stub()
               .of('function1')
               .of('function2')
*/
export default function stub () {
  return {
    of: function (name: string, callback: Function, returnValue: any) {
      this[name] = function () {
        let args = Array.prototype.slice.call(arguments)
        this[name].calls.push(args)
        let ret = null
        if (callback) {
          ret = callback.apply(this, args)
        }
        if (returnValue) {
          return returnValue
        }
        return ret
      }
      this[name].calls = []
      return this
    }
  }
}
