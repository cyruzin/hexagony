package rest

import (
	"encoding/json"
	"log"
	"net/http"
)

// APIMessage is a struct for generic JSON response.
type APIMessage struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
}

// InvalidRequest handles API errors.
func InvalidRequest(
	w http.ResponseWriter,
	err error,
	message string,
	httpCode int,
) {
	log.Printf("Error: %s", err.Error())
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
