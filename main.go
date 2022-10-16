package main

import (
	"database/sql"
	"log"

	"github.com/bank/simple-bank/api"
	db "github.com/bank/simple-bank/db/sqlc"
	"github.com/bank/simple-bank/util"
	_ "github.com/lib/pq"
)

/*
- Pass "." here, which means the current folder, cuz our config file is in the same location with this main.go file
*/
func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal("cannot create server:", err) 

	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot connect to server")
	}

}
