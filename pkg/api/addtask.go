package api

import (
	"encoding/json"
	"net/http"

	"firstIteration/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": "invalid JSON"})
		return
	}

	if err := validateTask(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "title is required"})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "failed to save task"})
		return
	}

	writeJson(w, map[string]any{"id": id})
}
