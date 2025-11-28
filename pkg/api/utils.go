package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeJson(w http.ResponseWriter, data any) {
	writeJsonWithCode(w, data, http.StatusOK)
}

func writeJsonWithCode(w http.ResponseWriter, data any, code int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Ошибка сериализации JSON: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
