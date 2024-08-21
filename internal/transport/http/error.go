package http

import (
	"encoding/json"
	"net/http"
)

type JSONError struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

// ErrorJSON writes error in JSON format and status code to response.
func ErrorJSON(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	e := json.NewEncoder(w)
	e.Encode(JSONError{
		Status: http.StatusText(code),
		Error:  error,
	})
}
