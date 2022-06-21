package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
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
	Log_ID  any    `json:"log_id,omitempty"`
}

func TestBrokerIsWork(t *testing.T) {
	jsonData, _ := json.MarshalIndent("empty post request", "", "\t")

	resp, _ := http.Post("http://broker-service/", "", bytes.NewBuffer(jsonData))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d. Got %d.", http.StatusOK, resp.StatusCode)
	}
}

func TestUserLogin(t *testing.T) {

	resp := userLogin()
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("Expected status code %d. Got %d.", http.StatusAccepted, resp.StatusCode)
	}
	defer resp.Body.Close()

	respLog := getLog(resp, t)
	if respLog.StatusCode != http.StatusAccepted {
		t.Fatalf("Expected status code %d. Got %d.", http.StatusAccepted, resp.StatusCode)
	}
	defer respLog.Body.Close()

}

func userLogin() *http.Response {
	authPayload := AuthPayload{
		Email:    "admin@example.com",
		Password: "verysecret",
	}

	payload := RequestPayload{
		Action: "auth",
		Auth:   authPayload,
	}

	jsonData, _ := json.MarshalIndent(payload, "", "\t")

	resp, err := http.Post("http://broker-service/handle", "", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil
	}

	return resp
}

func getLog(resp *http.Response, t *testing.T) *http.Response {

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("cannot read body")
		return nil
	}

	var response jsonResponse
	json.Unmarshal(body, &response)

	logPayload := LogPayload{
		Data: fmt.Sprintf("%v", response.Log_ID),
	}

	payload := RequestPayload{
		Action: "getlog",
		Log:    logPayload,
	}

	jsonDataForLog, _ := json.MarshalIndent(payload, "", "\t")

	respLog, _ := http.Post("http://broker-service/handle", "", bytes.NewBuffer(jsonDataForLog))

	bodyLog, err := ioutil.ReadAll(respLog.Body)
	if err != nil {
		t.Fatal("cannot read resp", err)
	}

	var respLogger jsonResponse
	json.Unmarshal(bodyLog, &respLogger)

	if respLogger.Log_ID != logPayload.Data {
		t.Fatalf("Expected log ID %s. Got %s.", logPayload.Data, respLogger.Log_ID)
	}

	return respLog
}
