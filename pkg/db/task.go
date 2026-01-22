package db

import (
	"database/sql"
	"errors"
	"strconv"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func Tasks(limit int) ([]*Task, error) {
	conn := Conn()

	rows, err := conn.Query(
		`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var (
			id      int64
			date    string
			title   string
			comment sql.NullString
			repeat  sql.NullString
		)
		if err := rows.Scan(&id, &date, &title, &comment, &repeat); err != nil {
			return nil, err
		}
		task := &Task{
			ID:      strconv.FormatInt(id, 10),
			Date:    date,
			Title:   title,
			Comment: comment.String,
			Repeat:  repeat.String,
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func AddTask(task *Task) (int64, error) {
	if task == nil {
		return 0, errors.New("task is nil")
	}
	conn := Conn()
	if conn == nil {
		return 0, sql.ErrConnDone
	}
	res, err := conn.Exec(
		`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetTask(id string) (*Task, error) {
	conn := Conn()
	if conn == nil {
		return nil, sql.ErrConnDone
	}

	var (
		tid     int64
		date    string
		title   string
		comment sql.NullString
		repeat  sql.NullString
	)
	err := conn.QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`,
		id,
	).Scan(&tid, &date, &title, &comment, &repeat)
	if err != nil {
		return nil, err
	}

	return &Task{
		ID:      strconv.FormatInt(tid, 10),
		Date:    date,
		Title:   title,
		Comment: comment.String,
		Repeat:  repeat.String,
	}, nil
}


func UpdateTask(task *Task) error {
	if task == nil {
		return errors.New("task is nil")
	}
	conn := Conn()
	if conn == nil {
		return sql.ErrConnDone
	}
	res, err := conn.Exec(
		`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func UpdateDate(next string, id string) error {
	conn := Conn()
	if conn == nil {
		return sql.ErrConnDone
	}
	res, err := conn.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, next, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func DeleteTask(id string) error {
	conn := Conn()
	if conn == nil {
		return sql.ErrConnDone
	}
	res, err := conn.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
