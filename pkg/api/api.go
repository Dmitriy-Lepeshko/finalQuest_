package api

import "net/http"

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDateHandler)
	mux.HandleFunc("/api/task", authMiddleware(taskHandler))
	mux.HandleFunc("/api/tasks", authMiddleware(tasksHandler))
	mux.HandleFunc("/api/task/done", authMiddleware(doneHandler))
	mux.HandleFunc("/api/signin", signInHandler)
}
