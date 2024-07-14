package utils

import (
	"encoding/json"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

    jsonResponse, err := json.MarshalIndent(payload, "", "  ")
    if err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }

    if _, err := w.Write(jsonResponse); err != nil {
        http.Error(w, "Failed to write response", http.StatusInternalServerError)
    }
}
