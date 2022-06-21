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
	Log_ID  any    `json:"log_id,omitempty"`
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

	id, err := server.models.LogEntry.Insert(event)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
		Log_ID:  id,
	}

	ctx.JSON(http.StatusAccepted, resp)
}

func (server *Server) GetLog(ctx *gin.Context) {
	var requestPayload JSONPayload

	if err := ctx.ShouldBindJSON(&requestPayload); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	data := data.LogEntry{
		ID: requestPayload.Data,
	}

	id := fmt.Sprintf("%v", data.ID)
	
	message, err := server.models.LogEntry.GetOne(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
		Data:    "message.Data",
		Log_ID:  message.ID,
	}

	ctx.JSON(http.StatusAccepted, resp)
}
