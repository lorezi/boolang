// Helper package for miscellaneous functions for validations
package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// JsonDecoder to check req payload size, check unknown fields and decode req
func JSONDecoder(r io.ReadCloser, w http.ResponseWriter, m interface{}) (*json.Decoder, error) {
	r = http.MaxBytesReader(w, r, 1048576)
	d := json.NewDecoder(r)
	d.DisallowUnknownFields()
	err := d.Decode(m)

	return d, err
}

// JSONValidator to validate json request payload 👁👁👀
func JSONValidator(err error) (string, int) {

	var (
		ute *json.UnmarshalTypeError
		se  *json.SyntaxError
	)

	switch {
	case errors.Is(err, io.EOF), errors.Is(err, io.ErrUnexpectedEOF), errors.As(err, &se):
		return "malformed JSON payload", http.StatusBadRequest

	case errors.As(err, &ute):
		return "unexpected field value", http.StatusBadRequest

	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fn := strings.TrimPrefix(err.Error(), "json: unknown field ")
		fn = fmt.Sprintf("unexpected field name: %s", fn)
		return fn, http.StatusBadRequest

	case err.Error() == "http: request body too large":
		return "request body must not be larger than 1mb", http.StatusRequestEntityTooLarge

	default:
		return "invalid request payload", http.StatusInternalServerError

	}

}
