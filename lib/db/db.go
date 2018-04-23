package db

import (
	"bufio"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/util"
	"github.com/syou6162/go-active-learning/lib/util/file"
)

func CreateDBConnection() (*sql.DB, error) {
	host := util.GetEnv("POSTGRES_HOST", "localhost")
	dbUser := util.GetEnv("DB_USER", "nobody")
	dbPassword := util.GetEnv("DB_PASSWORD", "nobody")
	dbName := util.GetEnv("DB_NAME", "go-active-learning")
	return sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, dbUser, dbPassword, dbName))
}

func InsertOrUpdateExample(db *sql.DB, e *example.Example) (sql.Result, error) {
	var url string
	now := time.Now()

	err := db.QueryRow(`SELECT url FROM example WHERE url = $1`, e.Url).Scan(&url)
	switch {
	case err == sql.ErrNoRows:
		return db.Exec(`INSERT INTO example (url, label, created_at, updated_at) VALUES ($1, $2, $3, $4)`, e.Url, e.Label, now, now)
	case err != nil:
		return nil, err
	default:
		return db.Exec(`UPDATE example SET label = $2, updated_at = $3 WHERE url = $1 `, e.Url, e.Label, now)
	}
}

func InsertExampleFromScanner(db *sql.DB, scanner *bufio.Scanner) (*example.Example, error) {
	line := scanner.Text()
	e, err := file.ParseLine(line)
	if err != nil {
		return nil, err
	}
	_, err = InsertOrUpdateExample(db, e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func ReadExamples(db *sql.DB) ([]*example.Example, error) {
	rows, err := db.Query(`SELECT url, label FROM example`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var examples example.Examples

	for rows.Next() {
		var label example.LabelType
		var url string
		if err := rows.Scan(&url, &label); err != nil {
			return nil, err
		}
		e := example.Example{Url: url, Label: label}
		examples = append(examples, &e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return examples, nil
}

func DeleteAllExamples(db *sql.DB) (sql.Result, error) {
	return db.Exec(`DELETE FROM example`)
}
