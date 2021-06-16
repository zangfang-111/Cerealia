package middleware

import (
	"encoding/json"
	"net/http"
)

// WriteHTTPAppError writes app error into the HTTP response
func WriteHTTPAppError(w http.ResponseWriter, handlerError error, message string, code int) {
	errObj := struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}{
		handlerError.Error(),
		message,
	}
	// TODO: check if this is needed. Might be resolved by
	// content.TypeNegotiator in router.go
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	j, err := json.Marshal(errObj)
	if err != nil {
		logger.Error("Can't serialize error object", err)
		_, err = w.Write(j)
	} else {
		_, err = w.Write([]byte("Serialization error"))
	}
	if err != nil {
		logger.Error("Can't write the response", err)
	}
}
