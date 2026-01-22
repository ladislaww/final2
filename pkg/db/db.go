package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR NOT NULL,
	comment TEXT,
	"repeat" VARCHAR(128)
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`


var db *sql.DB

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := false
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		install = true
	}

	conn, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	if err := conn.Ping(); err != nil {
		conn.Close()
		return err
	}
	if install {
		if _, err := conn.Exec(schema); err != nil {
			conn.Close()
			return err
		}
	}
	db = conn
	return nil
}


func Conn() *sql.DB {
	return db
}

// Close closes the opened database connection.
func Close() error {
	if db == nil {
		return nil
	}
	return db.Close()
}
