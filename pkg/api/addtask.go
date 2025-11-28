package api

import (
	"encoding/json"
	"log"
	"net/http"

	"firstIteration/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJsonWithCode(w, map[string]string{"error": "invalid JSON"}, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJsonWithCode(w, map[string]string{"error": "title is required"}, http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		writeJsonWithCode(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		log.Printf("Ошибка сохранения задачи: %v", err)
		writeJsonWithCode(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]any{"id": id})
}
