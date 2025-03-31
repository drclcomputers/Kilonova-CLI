package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	u "net/url"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type signInResp struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

type extendResp struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

var signinCmd = &cobra.Command{
	Use:   "signin [username] [password]",
	Short: "Sign in to your account",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		login(args[0], args[1])
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of your account",
	Run: func(cmd *cobra.Command, args []string) {
		logout()
	},
}

var userGetDetailsCmd = &cobra.Command{
	Use:   "user [User ID or nf (get self ID)]",
	Short: "Get details about a user.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userGetDetails(args[0])
	},
}

var extendSessionCmd = &cobra.Command{
	Use:   "extendsession",
	Short: "Extend the current session for 30 days more.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		extendSession()
	},
}

func init() {
	rootCmd.AddCommand(signinCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(userGetDetailsCmd)
	rootCmd.AddCommand(extendSessionCmd)
}

// login
func login(username, password string) {

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

func userGetDetails(user_id string) {
	body, err := makeRequest("GET", URL_SELF, nil, "1")
	if err != nil {
		logErr(err)
		return
	}

	fmt.Println(string(body))
}

// extend session
func extendSession() {
	body, err := makeRequest("POST", URL_EXTEND_SESSION, nil, "1")
	if err != nil {
		logErr(err)
		return
	}

	var resp extendResp
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Printf("error unmarshalling response: %s", err)
		return
	}

	if resp.Status == "success" {
		parsedTime, err := time.Parse(time.RFC3339Nano, resp.Data)
		if err != nil {
			logErr(err)
			return
		}
		formattedTime := parsedTime.Format("2006-01-02 15:04:05")
		fmt.Println("Your session has been extended until ", formattedTime)
	} else {
		fmt.Println(resp.Data)
	}

}
