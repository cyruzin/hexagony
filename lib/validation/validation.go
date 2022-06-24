package validation

import (
	"context"
	"encoding/json"
	"hexagony/lib/clog"
	"net/http"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

// Validator is an interface for validation purposes.
type Validator interface {
	BindStruct(ctx context.Context, data interface{}) error
	BindField(ctx context.Context, data interface{}, tag string) error
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

// single instance for caching
var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

// New creates a new Validator.
func New() Validator {
	en := en.New()
	uni = ut.New(en, en)

	trans, _ = uni.GetTranslator("en")

	validate = validator.New()
	if err := enTranslations.RegisterDefaultTranslations(validate, trans); err != nil {
		clog.Error(err, "failed to register default translations")
	}

	return message{}
}

// BindStruct checks if the given struct is valid.
func (v message) BindStruct(ctx context.Context, data interface{}) error {
	if err := validate.StructCtx(ctx, data); err != nil {
		return err
	}
	return nil
}

// BindField checks if the given field is valid.
func (v message) BindField(ctx context.Context, data interface{}, tag string) error {
	if err := validate.VarCtx(ctx, data, tag); err != nil {
		return err
	}
	return nil
}

// DecodeError returns validation error messages.
func (v message) DecodeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	payload := &errors{}

	for _, err := range err.(validator.ValidationErrors) {
		payload.Errors = append(payload.Errors, &message{Message: err.Translate(trans)})
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		if _, err := w.Write([]byte("could not encode the payload")); err != nil {
			return
		}
		return
	}
}
