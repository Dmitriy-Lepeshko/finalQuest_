package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

const defaultPort = 7540

func Run(mux *http.ServeMux) {
	port := defaultPort
	if p := os.Getenv("TODO_PORT"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			port = parsed
		}
	}
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Сервер запущен на http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Printf("Сервер остановлен: %v\n", err)
		os.Exit(1)
	}
}
