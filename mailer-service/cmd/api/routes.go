package main

import (
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

//Server serves HTTP requests for our scraper service.
type Server struct {
	router *gin.Engine
	Mailer Mail
}

func NewServer(mailer Mail) (*Server, error) {
	server := &Server{
		Mailer: mailer,
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

	router.POST("/send", server.SendMail)

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
