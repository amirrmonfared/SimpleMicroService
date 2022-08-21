package main

import (
	"database/sql"
	"log"
	"os"

	db "github.com/amirrmonfared/SimpleMicroService/authentication-service/db/sqlc"
	_ "github.com/lib/pq"
)

const webPort = ":8080"

func main() {
	log.Println("starting authentication service")
	log.Println("--------------------------------")

	dsn := os.Getenv("DSN")
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Panic("cannot connect to database")
	}
	log.Println("connected to database")
	log.Println("--------------------------------")

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
