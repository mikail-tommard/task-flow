package httpapi

import (
	"encoding/json"
	"net/http"
)

func writeError(w http.ResponseWriter, statusCode int, err string) {
	writeJSON(w, statusCode, map[string]string{"error": err})
}

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}
