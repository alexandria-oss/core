package httputil

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexandria-oss/core/exception"
	"net/http"
)

// ResponseErrJSON writes the required error's HTTP status and message using io.Writer
// from the HTTP handler
func ResponseErrJSON(w http.ResponseWriter, err error) {
	if err != nil {
		code := ErrorToCode(err)

		w.WriteHeader(code)
		errJSON := json.NewEncoder(w).Encode(&GenericResponse{
			Message: exception.GetErrorDescription(err),
			Code:    string(code),
		})

		if errJSON != nil {
			// Print application/text if not working
			w.Header().Add("Content-Type", "text/plain")
			_, _ = fmt.Fprintf(w, `%v`, &GenericResponse{
				Message: exception.GetErrorDescription(err),
				Code:    string(code),
			})
		}
	}
}

// ErrorToCode returns HTTP status code depending on the error given
func ErrorToCode(err error) int {
	switch {
	case errors.Is(err, exception.EntityNotFound) || errors.Is(err, exception.EntitiesNotFound):
		return http.StatusNotFound
	case errors.Is(err, exception.RequiredField) || errors.Is(err, exception.InvalidFieldFormat) ||
		errors.Is(err, exception.InvalidID) || errors.Is(err, exception.EmptyBody) ||
		errors.Is(err, exception.InvalidFieldRange):
		return http.StatusBadRequest
	case errors.Is(err, exception.EntityExists):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
