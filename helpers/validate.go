package helpers

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
)

func JSONValidate(err error) string {
	var ute *json.UnmarshalTypeError
	var se *json.SyntaxError

	switch {
	case errors.Is(err, io.EOF), errors.Is(err, io.ErrUnexpectedEOF), errors.As(err, &se):
		return "Malformed JSON payload"

	case errors.As(err, &ute):
		return "Unexpected field value"

	case strings.HasPrefix(err.Error(), "json: unknown field "):
		return "Unexpected field name"

	default:
		return "Invalid request payload"

	}
}
