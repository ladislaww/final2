package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"v1/pkg/db"
)

type taskRequest struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetTask(w, r)
		return
	case http.MethodPost:
		handleAddTask(w, r)
		return
	case http.MethodPut:
		handleUpdateTask(w, r)
		return
	case http.MethodDelete:
		handleDeleteTask(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func writeError(w http.ResponseWriter, err error) {
	writeJSON(w, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}

func handleGetTask(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeError(w, errors.New("Не указан идентификатор"))
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, errors.New("Задача не найдена"))
			return
		}
		writeError(w, err)
		return
	}
	writeJSON(w, task)
}

func handleAddTask(w http.ResponseWriter, r *http.Request) {
	var req taskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err)
		return
	}
	task, err := validateTask(req, false)
	if err != nil {
		writeError(w, err)
		return
	}

	conn := db.Conn()
	if conn == nil {
		writeError(w, errors.New("База данных не инициализирована"))
		return
	}

	res, err := conn.Exec(`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		writeError(w, err)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	var req taskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err)
		return
	}
	task, err := validateTask(req, true)
	if err != nil {
		writeError(w, err)
		return
	}

	if err := db.UpdateTask(task); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, errors.New("Задача не найдена"))
			return
		}
		writeError(w, err)
		return
	}
	writeJSON(w, map[string]string{})
}

func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeError(w, errors.New("Не указан идентификатор"))
		return
	}
	if err := db.DeleteTask(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, errors.New("Задача не найдена"))
			return
		}
		writeError(w, err)
		return
	}
	writeJSON(w, map[string]string{})
}

func validateTask(req taskRequest, requireID bool) (*db.Task, error) {
	id := strings.TrimSpace(req.ID)
	if requireID && id == "" {
		return nil, errors.New("Не указан идентификатор")
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, errors.New("Не указан заголовок задачи")
	}

	today := time.Now().Format(dateLayout)
	dateStr := strings.TrimSpace(req.Date)
	if dateStr == "" {
		dateStr = today
	} else if _, err := time.Parse(dateLayout, dateStr); err != nil {
		return nil, errors.New("Неверный формат даты")
	}

	repeat := strings.TrimSpace(req.Repeat)
	if repeat != "" {
		next, err := NextDate(today, dateStr, repeat)
		if err != nil {
			return nil, errors.New("Неверный формат правила повторения")
		}
		if dateStr < today {
			dateStr = next
		}
	} else if dateStr < today {
		dateStr = today
	}

	return &db.Task{
		ID:      id,
		Date:    dateStr,
		Title:   title,
		Comment: req.Comment,
		Repeat:  repeat,
	}, nil
}
