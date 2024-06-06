package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	client *sql.DB
}

const file string = "alquiler-scrapping.db"

const create string = `
CREATE TABLE IF NOT EXISTS pisos (
	id INTEGER NOT NULL PRIMARY KEY,
	url TEXT NOT NULL UNIQUE,
	price INT,
	title TEXT,
	sent INTEGER
);`

func NewDatabase() (*Database, error) {
	client, err := sql.Open("sqlite3", os.Getenv("SQLITE_DB_FILENAME"))
	if err != nil {
		return nil, err
	}
	if _, err := client.Exec(create); err != nil {
		return nil, err
	}
	return &Database{
		client: client,
	}, nil
}

func (db *Database) Insert(entry Entry) (int, error) {
	res, err := db.client.Exec("INSERT OR IGNORE INTO pisos VALUES(NULL,?,?,?,0);", entry.Url, entry.Price, entry.Title)
	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func (db *Database) ListNotSent() (Entries, error) {
	rows, err := db.client.Query(`SELECT url, price, title FROM pisos WHERE sent=?`, 0)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := []Entry{}
	for rows.Next() {
		e := Entry{}
		err = rows.Scan(&e.Url, &e.Price, &e.Title)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func (db *Database) MarkAsSent(entry Entry) (int, error) {
	res, err := db.client.Exec("UPDATE pisos SET sent = 1 WHERE url = ?;", entry.Url)
	if err != nil {
		return 0, err
	}

	var rowsAffected int64
	if rowsAffected, err = res.RowsAffected(); err != nil {
		return 0, err
	}
	return int(rowsAffected), nil
}
