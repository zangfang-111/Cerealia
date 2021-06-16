package stellar

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/facebookgo/stack"
	"github.com/robert-zaremba/errstack"
	"github.com/robert-zaremba/log15"
	"github.com/stellar/go/clients/horizon"
)

var logger = log15.Root()

type herror struct {
	E          *horizon.Error
	stacktrace stack.Stack
	msg        string
}

func wrapErr(err error, msg string) errstack.E {
	if err == nil {
		return nil
	}
	if herr, ok := err.(*horizon.Error); ok {
		return herror{herr, stack.Callers(1), msg}
	}
	return errstack.WrapAsInf(err, msg)
}

func (e herror) WithMsg(msg string) errstack.E {
	return herror{e.E, e.stacktrace, fmt.Sprint(msg, " [", msg, "]")}
}

func (e herror) Error() string {
	return fmt.Sprintf("[%s]. %s", e.msg, e.E.Error())
}

// StatusCode returns the HTTP status code.
func (e herror) StatusCode() int {
	return 500 // e.E.Problem.Status
}

func (e herror) IsReq() bool {
	return e.StatusCode() < 500
}

func (e herror) Stacktrace() stack.Stack {
	return e.stacktrace
}

func (e herror) Cause() error {
	return e.E
}

// MarshalJSON implements json.Marshaler interface
func (e herror) MarshalJSON() ([]byte, error) {
	m := e.Extensions()
	m["message"] = e.msg
	return json.Marshal(m)
}

func (e herror) Extensions() map[string]interface{} {
	return map[string]interface{}{
		"result_codes": e.E.Problem.Extras["result_codes"],
		"title":        "Horizon error: " + e.E.Problem.Title,
		"type":         e.E.Problem.Type,
		"status":       e.E.Problem.Status}
}

func (e herror) Log() {
	codes, errCodes := e.E.ResultCodes()
	if errCodes != nil {
		logger.Error("Horizon: can't get error ResultCodes", errCodes)
	}

	logger.Error("Horizon error. "+e.msg,
		"title", e.E.Problem.Title,
		"status", strconv.Itoa(e.E.Problem.Status),
		"result_codes", codes,
		log15.Alone("stacktrace", skipInternalStack(e.stacktrace)),
	)
}

const pkgImport = "bitbucket.org/cerealia/apps"
const pkgVendor = pkgImport + "/vendor"

// skipInternalStack removes bottom stack frames which are not related to the
// project stack trace
func skipInternalStack(st stack.Stack) stack.Stack {
	var end = len(st) - 1
	for ; end > 0; end-- {
		f := st[end]
		if strings.HasPrefix(f.File, pkgImport) &&
			!strings.HasPrefix(f.File, pkgVendor) {
			break
		}
	}
	return st[:end+1]
}
