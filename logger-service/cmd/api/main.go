package main

import (
	"context"
	"log"
	"time"

	"github.com/amirrmonfared/SimpleMicroService/logger-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = ":80"
	mongoURL = "mongodb://mongo:27017"
)

var client *mongo.Client

func main() {

	log.Println("starting logger service")
	log.Println("--------------------------------")
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	log.Println("connected to database")
	log.Println("--------------------------------")

	client = mongoClient
	models := data.New(client)

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	server, err := NewServer(models)
	if err != nil {
		log.Println("cannot connect to server", err)
	}

	err = server.Start(webPort)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}

func connectToMongo() (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	log.Println("connected to mongo")

	return c, nil
}
