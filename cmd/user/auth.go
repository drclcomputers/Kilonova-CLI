// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	utility "kilocli/cmd/utility"
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
		utility.LogError(err)
		return utility.ERROR, utility.ERROR
	}

	return username, password
}

func createLoginToken(response utility.KilonovaResponse) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		utility.LogError(err)
		return
	}
	configDir := filepath.Join(homedir, ".config", "kn-cli")
	err = os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		utility.LogError(err)
		return
	}
	tokenFile := filepath.Join(configDir, "token")

	file, err := os.Create(tokenFile)
	if err != nil {
		utility.LogError(fmt.Errorf("error creating file: %v", err))
		return
	}
	defer file.Close()

	encryptedToken, err := utility.Encrypt(response.Data)
	if err != nil {
		utility.LogError(fmt.Errorf("error encrypting token: %v", err))
		return
	}

	if err := os.WriteFile(tokenFile, []byte(encryptedToken), 0644); err != nil {
		utility.LogError(fmt.Errorf("error writing auth token to file: %v", err))
		return
	}
}

func login(username, password string) {
	formData := u.Values{
		"username": {username},
		"password": {password},
	}

	ResponseBody, err := utility.MakePostRequest(utility.URL_LOGIN, bytes.NewBufferString(formData.Encode()), utility.RequestFormGuest)
	if err != nil {
		utility.LogError(fmt.Errorf("login failed: %v", err))
		return
	}

	if !bytes.Contains(ResponseBody, []byte(utility.SUCCESS)) {
		utility.LogError(fmt.Errorf("login failed: invalid credentials"))
		return
	}

	var response utility.KilonovaResponse
	if err := json.Unmarshal(ResponseBody, &response); err != nil {
		utility.LogError(fmt.Errorf("error parsing response: %v", err))
		return
	}

	createLoginToken(response)

	fmt.Println("Login successful!")
}

// logout

func removeTokenFile() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		utility.LogError(err)
		return
	}
	configDir := filepath.Join(homedir, ".config", "kn-cli")
	tokenFile := filepath.Join(configDir, "token")
	_ = os.Remove(tokenFile)
}

func logout() {
	JSONData := []byte(`{"key": "value"}`)
	ResponseBody, err := utility.MakePostRequest(utility.URL_LOGOUT, bytes.NewBuffer(JSONData), utility.RequestFormAuth)
	if err != nil {
		utility.LogError(err)
		return
	}

	if bytes.Contains(ResponseBody, []byte(utility.SUCCESS)) {
		fmt.Println("Logged out successfully!")
		removeTokenFile()
	}
	log.Println("Logout failed: You must be logged in to do this!")
}

func extendSession() {
	ResponseBody, err := utility.MakePostRequest(utility.URL_EXTEND_SESSION, nil, utility.RequestFormAuth)
	if err != nil {
		utility.LogError(err)
		return
	}

	var resp utility.KilonovaResponse
	if err := json.Unmarshal(ResponseBody, &resp); err != nil {
		utility.LogError(fmt.Errorf("error unmarshalling response: %s", err))
		return
	}

	if resp.Status == utility.SUCCESS {
		formattedTime, err := utility.ParseTime(resp.Data)
		if err != nil {
			utility.LogError(fmt.Errorf("error parsing time: %s", err))
			return
		}
		fmt.Println("Your session has been extended until ", formattedTime)
	}
	fmt.Println(resp.Data)

}
