package main

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = ":80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	rabbitConn, err := amqp.Dial("amqp://guest:guest@rabbitmq")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	app := Config{
		Rabbit: rabbitConn,
	}

	log.Printf("Starting broker service on port %s\n", webPort)

	server, err := NewServer(app)
	if err != nil {
		log.Println("cannot connect to server", err)
	}

	err = server.Start(webPort)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
