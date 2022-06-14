package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (server *Server) Broker(ctx *gin.Context) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hi! I'm from broker",
	}

	ctx.JSON(http.StatusOK, payload)
}

func (server *Server) HandleSubmission(ctx *gin.Context) {
	var requestPayload RequestPayload

	if err := ctx.ShouldBindJSON(&requestPayload); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var err error
	var errPayload jsonResponse
	errPayload.Error = true
	errPayload.Message = err.Error()

	switch requestPayload.Action {
	case "auth":
		server.authenticate(ctx, requestPayload.Auth)
	case "log":
		server.logItem(ctx, requestPayload.Log)
	default:
		ctx.JSON(http.StatusBadRequest, errPayload)
	}
}

// authenticate calls the authentication microservice and sends back the appropriate response
func (server *Server) authenticate(ctx *gin.Context, a AuthPayload) {
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		ctx.JSON(http.StatusBadRequest, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		ctx.JSON(http.StatusBadRequest, errors.New("error calling auth service"))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if jsonFromService.Error {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	ctx.JSON(http.StatusAccepted, payload)
}

func (server *Server) logItem(ctx *gin.Context, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	ctx.JSON(http.StatusAccepted, payload)
}
