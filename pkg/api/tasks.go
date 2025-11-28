package api

import (
	"net/http"
	"strconv"
	"time"

	"firstIteration/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if lStr := r.FormValue("limit"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = l
		}
	}

	search := r.FormValue("search")

	var tasks []*db.Task
	var err error

	if search != "" {
		if len(search) == 10 && search[2] == '.' && search[5] == '.' {
			if parsed, err2 := time.Parse("02.01.2006", search); err2 == nil {
				dateStr := parsed.Format("20060102")
				tasks, err = db.SearchByDate(dateStr, limit)
			} else {
				tasks, err = db.SearchByKeyword(search, limit)
			}
		} else {
			tasks, err = db.SearchByKeyword(search, limit)
		}
	} else {
		tasks, err = db.Tasks(limit)
	}

	if err != nil {
		writeJson(w, map[string]string{"error": "failed to fetch tasks"})
		return
	}

	writeJson(w, TasksResp{Tasks: tasks})
}
