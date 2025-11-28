package api

import (
	"net/http"
	"time"

	"firstIteration/pkg/db"
)

func doneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "id is required"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
	} else {
		now := time.Now()
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
		if err := db.UpdateDate(next, id); err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
	}

	writeJson(w, map[string]any{})
}
