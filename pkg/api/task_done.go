package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"v1/pkg/db"
)

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeError(w, errors.New("Не указан идентификатор"), http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, errors.New("Задача не найдена"), http.StatusNotFound)
			return
		}
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	if strings.TrimSpace(task.Repeat) == "" {
		if err := db.DeleteTask(id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeError(w, errors.New("Задача не найдена"), http.StatusNotFound)
				return
			}
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]string{})
		return
	}

	today := time.Now().Format(dateLayout)
	next, err := NextDate(today, task.Date, task.Repeat)
	if err != nil {
		writeError(w, errors.New("Неверный формат правила повторения"), http.StatusBadRequest)
		return
	}
	if err := db.UpdateDate(next, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, errors.New("Задача не найдена"), http.StatusNotFound)
			return
		}
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{})
}
