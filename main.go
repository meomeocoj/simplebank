package main

import (
	"database/sql"
	"log"

	"github.com/meomeocoj/simplebank/api"
	db "github.com/meomeocoj/simplebank/db/sqlc"
	"github.com/meomeocoj/simplebank/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("load config fail:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("connection fail:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("server start fail:", err)
	} else {
		log.Println("server started successfully at", config.ServerAddress)
	}

}

func NewStore(db *sql.DB) {
	panic("unimplemented")
}
