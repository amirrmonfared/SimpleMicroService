package main

import (
	"log"
	"os"
	"strconv"
)

const webPort = ":8080"

func main() {
	log.Println("starting mail service")
	log.Println("--------------------------------")

	mail := createMail()
	server, err := NewServer(mail)
	if err != nil {
		log.Println("cannot connect to server", err)
	}

	err = server.Start(webPort)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
	}

	return m
}
