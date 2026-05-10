// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

// Package user implements user authentication and account management commands
// for the Kilonova CLI application.
// It handles login, logout, session management, and user information retrieval.
package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kncli/cmd/database"
	"kncli/internal"
	u "net/url"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
)

// loginForm creates and displays a form for collecting username and password
// credentials from the user. The password input is masked for security.
//
// Returns:
//   - username: The username entered by the user
//   - password: The password entered by the user
//   - If an error occurs during form processing, returns internal.ERROR for both values
func loginForm() (string, string, bool) {
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
		return "", "", false
	}

	return username, password, true
}

// createLoginToken securely stores the authentication token received from the server.
// The token is encrypted before being written to disk for security.
//
// Parameters:
//   - response: The KilonovaResponse containing the authentication token data
//
// The token is stored in the configuration directory with the filename defined
// by internal.TOKENFILENAME. File permissions are set to 0644 (readable by owner
// and group, readable by others).
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

// login authenticates a user with the Kilonova platform using username and password.
// It sends a POST request to the login endpoint and processes the response.
//
// On successful login:
//   - Creates and stores an encrypted authentication token
//   - Initializes the local database if it doesn't exist
//   - Prints "Login successful!" to the console
//
// On failure:
//   - Logs the error using internal.LogError
//   - Returns without further action
//
// Parameters:
//   - username: The user's username or email
//   - password: The user's password
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

	var response internal.KilonovaResponse
	if err := json.Unmarshal(ResponseBody, &response); err != nil {
		internal.LogError(fmt.Errorf("error parsing response: %v", err))
		return
	}

	if response.Status != internal.SUCCESS {
		internal.LogError(fmt.Errorf("login failed: invalid credentials"))
		return
	}

	createLoginToken(response)

	fmt.Println("Login successful!")

	if !internal.DBExists() {
		database.CreateDB()
	}
}

// logout ends the current user session by sending a logout request to the server
// and removing the local authentication token.
//
// On successful logout:
//   - Prints "Logged out successfully!" to the console
//   - Removes the local token file
//
// On failure:
//   - Logs the error using internal.LogError
//   - Returns without further action
func logout() {
	ResponseBody, err := internal.MakePostRequest(internal.URL_LOGOUT, nil, internal.RequestFormAuth)
	if err != nil {
		internal.LogError(err)
		return
	}

	var response internal.KilonovaResponse
	if err := json.Unmarshal(ResponseBody, &response); err != nil {
		internal.LogError(fmt.Errorf("error parsing response: %v", err))
		return
	}

	if response.Status == internal.SUCCESS {
		fmt.Println("Logged out successfully!")
		removeTokenFile()
	} else {
		internal.LogError(fmt.Errorf("logout failed: %v", response.Data))
	}
}

// removeTokenFile deletes the local authentication token file from the configuration directory.
// This function is used during logout to clear stored credentials.
//
// The token file location is determined by combining the configuration directory
// (from internal.GetConfigDir()) with the token filename constant (internal.TOKENFILENAME).
// Any error during file removal is ignored as indicated by the blank identifier.
func removeTokenFile() {
	tokenFile := filepath.Join(internal.GetConfigDir(), internal.TOKENFILENAME)
	_ = os.Remove(tokenFile)
}

// extendSession sends a request to the server to prolong the current user session.
// It queries the session extension endpoint and processes the response to display
// the new session expiration time.
//
// On success:
//   - Parses and formats the expiration time from the response
//   - Prints the extended session expiration time to the console
//   - Also prints the raw response data
//
// On failure:
//   - Logs the error using internal.LogError
//   - Returns without further action
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
	} else {
		fmt.Println(resp.Data)
	}
}
