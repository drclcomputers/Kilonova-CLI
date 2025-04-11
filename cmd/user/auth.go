// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kncli/cmd/database"
	"kncli/internal"
	"log"
	u "net/url"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
)

// login
func loginForm() (string, string) {
	var username, password string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Username:").
				Value(&username),
			huh.NewInput().
				Title("Password:").
				Value(&password).
				EchoMode(huh.EchoModePassword),
		),
	)

	if err := form.Run(); err != nil {
		internal.LogError(err)
		return internal.ERROR, internal.ERROR
	}

	return username, password
}

func createLoginToken(response internal.KilonovaResponse) {
	tokenFile := filepath.Join(internal.GetConfigDir(), internal.TOKENFILENAME)

	file, err := os.Create(tokenFile)
	if err != nil {
		internal.LogError(fmt.Errorf("error creating file: %v", err))
		return
	}
	defer file.Close()

	encryptedToken, err := internal.Encrypt(response.Data)
	if err != nil {
		internal.LogError(fmt.Errorf("error encrypting token: %v", err))
		return
	}

	if err := os.WriteFile(tokenFile, []byte(encryptedToken), 0644); err != nil {
		internal.LogError(fmt.Errorf("error writing auth token to file: %v", err))
		return
	}
}

func login(username, password string) {
	formData := u.Values{
		"username": {username},
		"password": {password},
	}

	ResponseBody, err := internal.MakePostRequest(internal.URL_LOGIN, bytes.NewBufferString(formData.Encode()), internal.RequestFormGuest)
	if err != nil {
		internal.LogError(fmt.Errorf("login failed: %v", err))
		return
	}

	if !bytes.Contains(ResponseBody, []byte(internal.SUCCESS)) {
		internal.LogError(fmt.Errorf("login failed: invalid credentials"))
		return
	}

	var response internal.KilonovaResponse
	if err := json.Unmarshal(ResponseBody, &response); err != nil {
		internal.LogError(fmt.Errorf("error parsing response: %v", err))
		return
	}

	createLoginToken(response)

	fmt.Println("Login successful!")

	if !internal.DBExists() {
		database.CreateDB()
	}
}

// logout

func removeTokenFile() {
	tokenFile := filepath.Join(internal.GetConfigDir(), internal.TOKENFILENAME)
	_ = os.Remove(tokenFile)
}

func logout() {
	JSONData := []byte(`{"key": "value"}`)
	ResponseBody, err := internal.MakePostRequest(internal.URL_LOGOUT, bytes.NewBuffer(JSONData), internal.RequestFormAuth)
	if err != nil {
		internal.LogError(err)
		return
	}

	if bytes.Contains(ResponseBody, []byte(internal.SUCCESS)) {
		fmt.Println("Logged out successfully!")
		removeTokenFile()
	}
	log.Println("Logout failed: You must be logged in to do this!")
}

func extendSession() {
	ResponseBody, err := internal.MakePostRequest(internal.URL_EXTEND_SESSION, nil, internal.RequestFormAuth)
	if err != nil {
		internal.LogError(err)
		return
	}

	var resp internal.KilonovaResponse
	if err := json.Unmarshal(ResponseBody, &resp); err != nil {
		internal.LogError(fmt.Errorf("error unmarshalling response: %s", err))
		return
	}

	if resp.Status == internal.SUCCESS {
		formattedTime, err := internal.ParseTime(resp.Data)
		if err != nil {
			internal.LogError(fmt.Errorf("error parsing time: %s", err))
			return
		}
		fmt.Println("Your session has been extended until ", formattedTime)
	}
	fmt.Println(resp.Data)
}
