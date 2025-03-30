package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	u "net/url"
	"os"
)

// login
type signInResp struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

func login() {
	if len(os.Args) < 4 {
		log.Fatal("Insert both credintials or too many arguments were passed!")
	}

	username := os.Args[2]
	password := os.Args[3]

	formData := u.Values{
		"username": {username},
		"password": {password},
	}

	body, err := makeRequest("POST", URL_LOGIN, bytes.NewBufferString(formData.Encode()), "3")
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	if !bytes.Contains(body, []byte("success")) {
		log.Fatal("Login failed: Invalid credentials!")
	}

	var response signInResp
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatalf("Error parsing response: %v", err)
	}

	if err := os.WriteFile("./token", []byte(response.Data), 0644); err != nil {
		log.Fatalf("Error writing auth token to file: %v", err)
	}

	fmt.Println("Login successful!")

}

// logout
func logout() {
	jsonData := []byte(`{"key": "value"}`)
	body, err := makeRequest("POST", URL_LOGOUT, bytes.NewBuffer(jsonData), "1")
	if err != nil {
		logErr(err)
		return
	}

	if bytes.Contains(body, []byte("success")) {
		fmt.Println("Logged out successfully!")
		_ = os.Remove("token")
	} else {
		log.Println("Logout failed: You must be logged in to do this!")
	}
}
