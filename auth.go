package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	u "net/url"
	"os"
)

// login
type signin struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

func login() {
	url := "https://kilonova.ro/api/auth/login"

	var username, password string

	if len(os.Args) >= 4 {
		username = os.Args[2]
		password = os.Args[3]
	} else {
		fmt.Println("Insert both credintials!")
		os.Exit(1)
	}

	formData := u.Values{}
	formData.Set("username", username)
	formData.Set("password", password)

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	req.Header.Set("User-Agent", "KilonovaCLIClient/1.0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "guest")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s", err)
		os.Exit(1)
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Login successful!")
	} else {
		fmt.Println("Login failed.")
		os.Exit(1)
	}

	var signin signin
	if err := json.Unmarshal(body, &signin); err != nil {
		fmt.Printf("error unmarshalling response: %s", err)
		os.Exit(1)
	}

	token := signin.Data

	err = os.WriteFile("./token", []byte(token), 0644)
	if err != nil {
		fmt.Println("error writing auth token to file! Err: ", err)
		os.Exit(1)
	}

}

// logout
func logout() {
	url := "https://kilonova.ro/api/auth/logout"

	token, err := os.ReadFile("token")
	if err != nil {
		fmt.Println("Could not read session ID from file. Make sure you are logged in!")
		os.Exit(1)
	}

	jsonData := []byte(`{"key": "value"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	req.Header.Set("User-Agent", "KilonovaCLIClient/1.0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", string(token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Logout successful!")
		os.WriteFile("token", []byte(""), 0644)
	} else {
		fmt.Println("Logout failed. You must be authenticated to do this.")
		os.Exit(1)
	}
}
