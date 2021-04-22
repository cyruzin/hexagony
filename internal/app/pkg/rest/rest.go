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

func DecodeError(w http.ResponseWriter, r *http.Request, err error, httpCode int) {
	log.Error().
		Err(err).
		Stack().
		Str("method", r.Method).
		Str("url", r.URL.String()).
		Str("agent", r.UserAgent()).
		Str("referer", r.Referer()).
		Str("proto", r.Proto).
		Str("remote_address", r.RemoteAddr).
		Int("status", httpCode).
		Msg(err.Error())

	w.WriteHeader(httpCode)

	apiError := &APIMessage{err.Error(), httpCode}
	if err := json.NewEncoder(w).Encode(apiError); err != nil {
		return
	}
}

func EncodeJSON(w http.ResponseWriter, httpCode int, dest interface{}) {
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(dest)
}
