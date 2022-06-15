package validation

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator is an interface for validation purposes.
type Validator interface {
	Bind(ctx context.Context, data interface{}) error
	DecodeError(w http.ResponseWriter, err error)
}

// message is a struct for validation error messages.
type message struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
}

// errors type is a struct for multiple error messages.
type errors struct {
	Errors []*message `json:"errors"`
}

// errorMap improves error messages.
func (v message) errorMap(err validator.FieldError) *message {
	errMap := map[string]string{
		"required": "is required",
		"email":    "is not valid",
		"min":      "minimum length is " + err.Param(),
		"gte":      "minimum length is " + err.Param(),
	}

	return &message{
		Message: "the " + strings.ToLower(err.Field()) + " field " + errMap[err.Tag()],
	}
}

// Bind checks if the given struct is valid.
func (v message) Bind(ctx context.Context, data interface{}) error {
	if err := validator.New().StructCtx(ctx, data); err != nil {
		return err
	}
	return nil
}

// DecodeError returns validation error messages.
func (v message) DecodeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)

	message := &errors{}

	for _, err := range err.(validator.ValidationErrors) {
		message.Errors = append(message.Errors, v.errorMap(err))
	}

	if err := json.NewEncoder(w).Encode(message); err != nil {
		if _, err := w.Write([]byte("could not encode the payload")); err != nil {
			return
		}
		return
	}
}

// New creates a new Validator.
func New() Validator {
	return message{}
}
