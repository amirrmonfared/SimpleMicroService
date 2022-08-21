package main

import (
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

//Server serves HTTP requests for our scraper service.
type Server struct {
	router *gin.Engine
	config Config
}

func NewServer(config Config) (*Server, error) {
	server := &Server{
		config: config,
	}

	// Initialize a new Gin router
	router := gin.New()

	// Apply the middleware to the router (works with groups too)
	router.Use(cors.Middleware(cors.Config{
		Origins:        "*",
		Methods:        "GET, PUT, POST, DELETE",
		RequestHeaders: "Origin, Authorization, Content-Type",
		ExposedHeaders: "",
		MaxAge:         50 * time.Second,
	}))

	router.POST("/", server.Broker)
	router.POST("/handle", server.HandleSubmission)

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
