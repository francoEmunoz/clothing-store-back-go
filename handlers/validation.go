package handlers

import (
	"encoding/json"
	"net/http"
)

func sendValidationError(rw http.ResponseWriter, err error) {
	rw.WriteHeader(http.StatusBadRequest)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
