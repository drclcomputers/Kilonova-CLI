package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	u "net/url"
	"os"
	"strings"
)

// login
type signin struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

func login() {
	url := URL_LOGIN

	var username, password string

	if len(os.Args) == 4 {
		username = os.Args[2]
		password = os.Args[3]
	} else {
		fmt.Println("Insert both credintials!")
		os.Exit(1)
	}

	formData := u.Values{}
	formData.Set("username", username)
	formData.Set("password", password)

	body, err := makeRequest("POST", url, bytes.NewBufferString(formData.Encode()), "3")
	if err != nil {
		logErr(err)
	}

	if strings.Contains(string(body), "success") {
		fmt.Println("Login successful!")
	} else {
		fmt.Println("Login failed. Invalid credentials!")
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
	url := URL_LOGOUT

	jsonData := []byte(`{"key": "value"}`)
	body, err := makeRequest("POST", url, bytes.NewBuffer(jsonData), "1")
	if err != nil {
		logErr(err)
	}

	//fmt.Println(string(body))

	if strings.Contains(string(body), "success") {
		fmt.Println("Logged out succesfully!")
		os.WriteFile("token", []byte(""), 0644)
	} else {
		fmt.Println("Logged out failed! You have to be logged in to do this!")
	}
}
