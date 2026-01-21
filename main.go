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

	port := ":7540"
	webDir := "./web"
	http.HandleFunc("/api/nextdate", api.NextDateHandler)
	http.HandleFunc("/api/task", api.TaskHandler)
	http.HandleFunc("/api/task/done", api.TaskDoneHandler)
	http.HandleFunc("/api/tasks", api.TasksHandler)
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Println("Server started on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
