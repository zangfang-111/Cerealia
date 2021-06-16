// Package render defines function for response rendering
package render

import (
	"encoding/json"

	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/robert-zaremba/errstack"
)

// Resp renders the response of request in proper Json style.
func Resp(v interface{}, c *routing.Context) errstack.E {
	response, err := json.Marshal(v)
	if err != nil {
		return errstack.WrapAsInf(err, "Can not marshal the data")
	}
	err = c.Write(response)
	return errstack.WrapAsInf(err, "Can not send the response")
}
