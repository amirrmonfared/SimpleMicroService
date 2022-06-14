package main

import (
	"log"
)

const webPort = ":80"

func main() {
	log.Printf("Starting broker service on port %s\n", webPort)

	server, err := NewServer()
	if err != nil {
		log.Println("cannot connect to server", err)
	}

	err = server.Start(webPort)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
