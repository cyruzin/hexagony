package rest

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
)

var json = jsoniter.ConfigFastest

// APIMessage is a struct for generic JSON response.
type APIMessage struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
}

// InvalidRequest handles API errors.
func InvalidRequest(
	w http.ResponseWriter,
	r *http.Request,
	err error,
	message string,
	httpCode int,
) {
	log.Error().
		Err(err).
		Str("method", r.Method).
		Str("url", r.URL.String()).
		Str("agent", r.UserAgent()).
		Str("referer", r.Referer()).
		Str("proto", r.Proto).
		Str("remote_address", r.RemoteAddr).
		Int("status", httpCode).
		Msg(message)

	w.WriteHeader(httpCode)

	apiError := &APIMessage{message, httpCode}
	if err := json.NewEncoder(w).Encode(apiError); err != nil {
		return
	}
}

// ToJSON returns a JSON response.
func ToJSON(
	w http.ResponseWriter,
	httpCode int,
	dest interface{},
) {
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(dest)
}
