package validation

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validation struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
}

// APIErrors type is a struct for multiple error messages.
type APIErrors struct {
	Errors []*Validation `json:"errors"`
}

func validationMap(err validator.FieldError) *Validation {
	errMap := map[string]string{
		"required": "is required",
		"email":    "is not valid",
		"min":      "minimum length is " + err.Param(),
		"gte":      "minimum length is " + err.Param(),
	}

	return &Validation{
		Message: "the " + strings.ToLower(err.Field()) + " field " + errMap[err.Tag()],
	}
}

// Message handles validation error messages.
func Message(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)

	apiValidator := &APIErrors{}

	for _, err := range err.(validator.ValidationErrors) {
		apiValidator.Errors = append(apiValidator.Errors, validationMap(err))
	}

	if err := json.NewEncoder(w).Encode(apiValidator); err != nil {
		if _, err := w.Write([]byte("could not encode the payload")); err != nil {
			return
		}
		return
	}
}
