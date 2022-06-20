package main

import (
	"database/sql"
	"log"
	"os"

	db "github.com/amirrmonfared/testMicroServices/authentication-service/db/sqlc"
	_ "github.com/lib/pq"
)

const webPort = ":80"

var counts int64

func main() {
	log.Println("starting authentication service")

	//connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	store := db.NewStore(conn)
	server, err := NewServer(store)
	if err != nil {
		log.Println("cannot connect to server", err)
	}

	err = server.Start(webPort)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil
	}

	err = conn.Ping()
	if err != nil {
		return nil
	}

	return conn
}
