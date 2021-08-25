package data

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func OpenDatabase() error {
	var err error

	db, err = sql.Open("sqlite3", "./main.db")
	if err != nil {
		return err
	}

	return db.Ping()
}

func CreateTable() {
	createTableSQL := ` CREATE TABLE IF NOT EXISTS stubby (
		"idNote" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"word" TEXT,
		"definition" TEXT,
		"category" TEXT
	);`

	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	statement.Exec()
	log.Println("stubby table created")

}

func InsertNote(word string, definition string, category string) {
	insertNoteSQL := `INSERT INTO stubby (word, definition, category)
	VALUES (?, ?, ?)`

	statement, err := db.Prepare(insertNoteSQL)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = statement.Exec(word, definition, category)
	if err != nil {
		log.Fatalln(err)
	}

}

func ListAllNotes() {
	row, err := db.Query("SELECT * FROM stubby")
	if err != nil {
		log.Fatalln(err)
	}

	defer row.Close()

	for row.Next() {
		var idNote int
		var word string
		var definition string
		var category string

		row.Scan(&idNote, &word, &definition, &category)
		fmt.Printf("[%s] %s - %s\n", category, word, definition)
	}

}
