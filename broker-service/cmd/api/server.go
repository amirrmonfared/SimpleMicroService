package main

import "github.com/gin-gonic/gin"

//Server serves HTTP requests for our scraper service.
type Server struct {
	router *gin.Engine
}

func NewServer() (*Server, error) {
	server := &Server{}
	router := gin.Default()

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
