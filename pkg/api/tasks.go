package api

import (
	"net/http"
	"v1/pkg/db"
)

type tasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

const tasksLimit = 50

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tasks, err := db.Tasks(tasksLimit)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}
	writeJSON(w, tasksResp{Tasks: tasks})
}
