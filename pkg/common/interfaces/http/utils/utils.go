package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload) // Игнорируем ошибку энкодинга, не критично
}

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request format",
			"code":  "VALIDATION_ERROR",
		})
		return false
	}
	return true
}
