// Helper package for miscellaneous functions for validations
package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// JSONValidator to validate json request payload 👁👁👀
func JSONValidator(err error) string {

	var (
		ute *json.UnmarshalTypeError
		se  *json.SyntaxError
	)

	switch {
	case errors.Is(err, io.EOF), errors.Is(err, io.ErrUnexpectedEOF), errors.As(err, &se):
		return "Malformed JSON payload"

	case errors.As(err, &ute):
		return "Unexpected field value"

	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fn := strings.TrimPrefix(err.Error(), "json: unknown field ")
		return fmt.Sprintf("Unexpected field name: %s", fn)

	default:
		return "Invalid request payload"

	}
}
