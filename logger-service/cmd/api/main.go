package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"

	"github.com/amirrmonfared/SimpleMicroService/logger-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = ":80"
	mongoURL = "mongodb://mongo:27017"
	rpcPort  = "5001"
	gRpcPort = "50001"
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

	go server.rpcListen()
	

	// start web server
	log.Println("Starting service on port", webPort)
	err = server.Start(webPort)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}

func (server *Server) rpcListen() error {
	log.Println("Starting RPC server on port ", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
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
