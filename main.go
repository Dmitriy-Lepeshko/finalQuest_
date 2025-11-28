package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"firstIteration/pkg/api"
	"firstIteration/pkg/db"
	"firstIteration/pkg/server"
)

const (
	dbFile = "scheduler.db"
	webDir = "./web"
)

func main() {

	godotenv.Load()

	dbFile := "scheduler.db"
	if envDB := os.Getenv("TODO_DBFILE"); envDB != "" {
		dbFile = envDB
	}

	dbConn, err := db.Init(dbFile)
	if err != nil {
		fmt.Printf("Ошибка инициализации БД: %v\n", err)
		os.Exit(1)
	}
	defer dbConn.Close()
	fmt.Printf("База данных '%s' готова.\n", dbFile)

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	api.Init(mux)

	server.Run(mux)
}
