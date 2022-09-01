package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/rpc"

	"github.com/amirrmonfared/SimpleMicroService/broker-service/event"
	"github.com/gin-gonic/gin"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	LogID   interface{} `json:"log_id,omitempty"`
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

	var errPayload jsonResponse
	errPayload.Error = true

	switch requestPayload.Action {
	case "auth":
		server.authenticate(ctx, requestPayload.Auth)
	case "log":
		server.logItemViaRPC(ctx, requestPayload.Log)
	case "getlog":
		server.getLog(ctx, requestPayload.Log)
	case "mail":
		server.SendMail(ctx, requestPayload.Mail)
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
	payload.LogID = jsonFromService.LogID

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

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"
	payload.LogID = jsonFromService.LogID

	ctx.JSON(http.StatusAccepted, payload)
}

func (server *Server) getLog(ctx *gin.Context, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/get-log"

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

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.LogID = jsonFromService.LogID

	ctx.JSON(http.StatusAccepted, payload)
}

func (server *Server) SendMail(ctx *gin.Context, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	// call the mail service
	mailServiceURL := "http://mailer-service/send"

	// post to mail service
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
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

	// make sure we get back the right status code
	if response.StatusCode != http.StatusAccepted {
		ctx.JSON(http.StatusBadRequest, errors.New("error calling mail service"))
		return
	}

	// send back json
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To

	ctx.JSON(http.StatusAccepted, payload)
}

// logEventViaRabbit logs an event using the logger-service. It makes the call by pushing the data to RabbitMQ.
func (server *Server) logEventViaRabbit(ctx *gin.Context, l LogPayload) {
	err := server.pushToQueue(l.Name, l.Data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	ctx.JSON(http.StatusAccepted, payload)
}

// pushToQueue pushes a message into RabbitMQ
func (server *Server) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(server.config.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}

type RPCPayload struct {
	Name string
	Data string
}

func (server *Server) logItemViaRPC(ctx *gin.Context, l LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}

	var result string
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: result,
	}

	ctx.JSON(http.StatusAccepted, payload)
}
