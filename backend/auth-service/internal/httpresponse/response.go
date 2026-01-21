package httpresponse

import (
	"auth-service/internal/dto"
	"auth-service/internal/logger"
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	js, err := json.Marshal(data)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Msg("error encoding response")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
}

func WriteError(w http.ResponseWriter, message string, status int) {
	WriteJSON(w, status, dto.ErrorResponse{
		Message: message,
		Code:    status,
	})
}
