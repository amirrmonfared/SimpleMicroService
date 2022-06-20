package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func TestBroker(t *testing.T) {
	jsonData, _ := json.MarshalIndent("empty post request", "", "\t")

	resp, _ := http.Post("http://broker-service", "", bytes.NewBuffer(jsonData))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d. Got %d.", http.StatusOK, resp.StatusCode)
	}
}

func TestUserLogin(t *testing.T) {

	authPayload := AuthPayload{
		Email:    "admin@example.com",
		Password: "verysecret",
	}

	payload := RequestPayload{
		Action: "auth",
		Auth:   authPayload,
	}

	jsonData, _ := json.MarshalIndent(payload, "", "\t")

	// check user login request accepted
	resp, _ := http.Post("http://broker-service/handle", "", bytes.NewBuffer(jsonData))
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("Expected status code %d. Got %d.", http.StatusAccepted, resp.StatusCode)
	}
	defer resp.Body.Close()

	fmt.Println("hi1")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("cannot read resp", err)
	}

	var respones jsonResponse
	json.Unmarshal(body, &respones)

	// check user login request logged into log database
	logID := respones.Log_ID

	jsonDataForLog, _ := json.MarshalIndent(logID, "", "\t")

	respLog, _ := http.Post("http://logger-service/get-log", "", bytes.NewBuffer(jsonDataForLog))
	if respLog.StatusCode != http.StatusAccepted {
		t.Fatalf("Expected status code %d. Got %d.", http.StatusAccepted, respLog.StatusCode)
	}

	fmt.Println("hi2")

	defer resp.Body.Close()

	bodyLog, err := ioutil.ReadAll(respLog.Body)
	if err != nil {
		t.Fatal("cannot read resp", err)
	}

	var responeslog jsonResponse
	json.Unmarshal(bodyLog, &responeslog)

	fmt.Println(resp)
}
