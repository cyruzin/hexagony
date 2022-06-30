package rest

import (
	"encoding/json"
	"net/http"
)

// Message is a struct for generic JSON response.
type Message struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
}

// DecodeError returns unsuccessful JSON error message.
func DecodeError(w http.ResponseWriter, r *http.Request, err error, httpCode int) {
	w.WriteHeader(httpCode)

	errorMessage := &Message{err.Error(), httpCode}
	if err := json.NewEncoder(w).Encode(errorMessage); err != nil {
		return
	}
}

// JSON returns a successful JSON message.
func JSON(w http.ResponseWriter, httpCode int, dest interface{}) {
	w.WriteHeader(httpCode)
	if err := json.NewEncoder(w).Encode(dest); err != nil {
		return
	}
}
