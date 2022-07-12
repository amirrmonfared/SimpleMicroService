package main

import (
	"time"

	db "github.com/amirrmonfared/SimpleMicroService/authentication-service/db/sqlc"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

//Server serves HTTP requests for our scraper service.
type Server struct {
	router *gin.Engine
	store  db.Store
}

func NewServer(store db.Store) (*Server, error) {
	server := &Server{
		store: store,
	}
	// Initialize a new Gin router
	router := gin.New()

	// Apply the middleware to the router (works with groups too)
	router.Use(cors.Middleware(cors.Config{
		Origins:        "*",
		Methods:        "GET, PUT, POST, DELETE, OPTIONS",
		RequestHeaders: "Accept, Authorization, Content-Type, X-CSRF-Token",
		ExposedHeaders: "Link",
		MaxAge:         50 * time.Second,
	}))

	router.POST("/authenticate", server.Authenticate)

	server.router = router
	return server, nil
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
