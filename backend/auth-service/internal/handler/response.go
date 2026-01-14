package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	js, err := json.Marshal(data)
	if err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	//var buf bytes.Buffer
	//
	//if err := json.NewEncoder(&buf).Encode(data); err != nil {
	//	log.Printf("error encoding response: %v", err)
	//	http.Error(w, "internal server error", http.StatusInternalServerError)
	//	return
	//}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
}

func writeError(w http.ResponseWriter, message string, status int) {
	writeJSON(w, status, ErrorResponse{Message: message, Code: status})
}
