package main

import (
	"fmt"
	"net/http"

	"github.com/amirrmonfared/testMicroServices/logger-service/data"
	"github.com/gin-gonic/gin"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (server *Server) WriteLog(ctx *gin.Context) {
	var requestPayload JSONPayload

	if err := ctx.ShouldBindJSON(&requestPayload); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := server.models.LogEntry.Insert(event)
	fmt.Println("hi")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	ctx.JSON(http.StatusAccepted, resp)
}
