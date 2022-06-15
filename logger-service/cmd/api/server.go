package main

import (
	"github.com/amirrmonfared/testMicroServices/logger-service/data"
	"github.com/gin-gonic/gin"
)

//Server serves HTTP requests for our scraper service.
type Server struct {
	router *gin.Engine
	models data.Models
}

func NewServer() (*Server, error) {
	server := &Server{}
	router := gin.Default()

	router.POST("/log", server.WriteLog)

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
