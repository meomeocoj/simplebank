package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dataDriver = "postgres"
	dataSource = "postgres://root:secret@localhost:6432/simplebank?sslmode=disable"
)

var testingQueries *Queries

var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open(dataDriver, dataSource)
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	defer db.Close()
	testingQueries = New(db)

	os.Exit(m.Run())
}
