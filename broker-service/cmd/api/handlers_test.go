package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestBroker(t *testing.T) {
	jsonData, _ := json.MarshalIndent("empty post request", "", "\t")

	resp, _ := http.Post("http://localhost:8080/", "", bytes.NewBuffer(jsonData))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d. Got %d.", http.StatusOK, resp.StatusCode)
	}
}

func TestAuthenticate(t *testing.T) {
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
	resp, _ := http.Post("http://localhost:8080/handle", "", bytes.NewBuffer(jsonData))
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("Expected status code %d. Got %d.", http.StatusAccepted, resp.StatusCode)
	}
	defer resp.Body.Close()

}
