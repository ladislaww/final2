package main

import (
	"log"
	"net/http"
	"v1/pkg/api"
	"v1/pkg/db"
)

func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("Failed to close DB:", err)
		}
	}()

	port := ":7540"
	webDir := "./web"
	http.HandleFunc("/api/nextdate", api.NextDateHandler)
	http.HandleFunc("/api/task", api.TaskHandler)
	http.HandleFunc("/api/task/done", api.TaskDoneHandler)
	http.HandleFunc("/api/tasks", api.TasksHandler)
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Println("Server started on port", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Println("Server stopped with error:", err)
	}
}
