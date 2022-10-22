package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/meomeocoj/simplebank/utils"
)

var testingQueries *Queries

var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	db, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	defer db.Close()
	testingQueries = New(db)

	os.Exit(m.Run())
}
