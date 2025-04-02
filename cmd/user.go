// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	u "net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
)

type userDetailResp struct {
	Data struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Admin       bool   `json:"admin"`
		Proposer    bool   `json:"proposer"`
		DisplayName string `json:"display_name"`
	} `json:"data"`
}

type userSolvedProblems struct {
	Data []struct {
		Problem_ID int     `json:"id"`
		Name       string  `json:"name"`
		Source     string  `json:"source_credits"`
		Score      float64 `json:"score_scale"`
	} `json:"data"`
}

var signinCmd = &cobra.Command{
	Use:   "signin [username] [password]",
	Short: "Sign in to your account",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		username, password := LoginForm()
		action := func() { login(username, password) }
		if err := spinner.New().Title("Logging in...").Action(action).Run(); err != nil {
			log.Fatal(err)
		}
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
	Use:   "user [User ID or me (get self ID)]",
	Short: "Get details about a user.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userGetDetails(args[0], "user")
	},
}

var userSolvedProblemsCmd = &cobra.Command{
	Use:   "solvedproblems [User ID or me (get self ID)]",
	Short: "Get list of solved problems by user.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userGetSolvedProblems(args[0])
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

var setBioCmd = &cobra.Command{
	Use:   "setbio [bio]",
	Short: "Set your profile's bio.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		setBio(args[0])
	},
}

var changeNameCmd = &cobra.Command{
	Use:   "changename [new name] [password]",
	Short: "Change your profile name.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		changeName(args[0], args[1])
	},
}

var changePassCmd = &cobra.Command{
	Use:   "changepass [old password] [new password]",
	Short: "Change your account password.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		changePass(args[0], args[1])
	},
}

var resetPassCmd = &cobra.Command{
	Use:   "resetpass [email]",
	Short: "Reset password via email when forgotten.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resetPass(args[0])
	},
}

var deleteUserCmd = &cobra.Command{
	Use:   "deleteuser",
	Short: "Delete your Kilonova account.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		deleteUser()
	},
}

var changeEmailCmd = &cobra.Command{
	Use:   "changemail [new email] [password]",
	Short: "Change your account email.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		changeEmail(args[0], args[1])
	},
}

var resendEmailCmd = &cobra.Command{
	Use:   "resendemail",
	Short: "Resend verification mail.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		resendEmail()
	},
}

func init() {
	rootCmd.AddCommand(signinCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(userGetDetailsCmd)
	rootCmd.AddCommand(extendSessionCmd)
	rootCmd.AddCommand(userSolvedProblemsCmd)
	rootCmd.AddCommand(setBioCmd)
	rootCmd.AddCommand(changeNameCmd)
	rootCmd.AddCommand(changePassCmd)
	rootCmd.AddCommand(changeEmailCmd)
	rootCmd.AddCommand(resetPassCmd)
	rootCmd.AddCommand(deleteUserCmd)
	rootCmd.AddCommand(resendEmailCmd)
}

func LoginForm() (string, string) {
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
		log.Fatal(err)
	}

	return username, password
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

	var response KNResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatalf("Error parsing response: %v", err)
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		logErr(err)
	}
	configDir := filepath.Join(homedir, ".config", "kn-cli")
	err = os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		logErr(err)
	}
	tokenFile := filepath.Join(configDir, "token")

	file, err := os.Create(tokenFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	if err := os.WriteFile(tokenFile, []byte(response.Data), 0644); err != nil {
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
		homedir, err := os.UserHomeDir()
		if err != nil {
			logErr(err)
		}
		configDir := filepath.Join(homedir, ".config", "kn-cli")
		tokenFile := filepath.Join(configDir, "token")
		_ = os.Remove(tokenFile)
	} else {
		log.Println("Logout failed: You must be logged in to do this!")
	}
}

func getBio(name string) string {
	res, err := http.Get(fmt.Sprintf("https://kilonova.ro/profile/%s", name))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// Parse the HTML
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the About Me section and extract text from <p>
	bio := doc.Find("div.segment-panel.reset-list.statement-content.enhance-tables p").First().Text()

	return bio
}

func userGetDetails(user_id, use_case string) bool {
	var url string
	if user_id == "me" {
		url = URL_SELF
	} else {
		url = fmt.Sprintf(URL_USER, user_id)
	}

	body, err := makeRequest("GET", url, nil, "3")
	if err != nil {
		logErr(err)
		return false
	}

	var dataUser userDetailResp

	err = json.Unmarshal(body, &dataUser)
	if err != nil {
		logErr(err)
		return false
	}

	if dataUser.Data.DisplayName == "" {
		dataUser.Data.DisplayName = "-"
	}

	if use_case == "isadmin" {
		return dataUser.Data.Admin
	} else {
		bio := getBio(dataUser.Data.Name)
		fmt.Printf("ID: %d\nName: %s\nA.K.A: %s\nBio: %s\nAdmin: %t\nProposer: %t\n\n",
			dataUser.Data.Id, dataUser.Data.Name, dataUser.Data.DisplayName, bio,
			dataUser.Data.Admin, dataUser.Data.Proposer)
		return false
	}
}

func userGetSolvedProblems(user_id string) {
	var url string
	if user_id == "me" {
		url = URL_SELF_PROBLEMS
	} else {
		url = fmt.Sprintf(URL_USER_PROBLEMS, user_id)
	}

	body, err := makeRequest("GET", url, nil, "3")
	if err != nil {
		logErr(err)
		return
	}

	var dataUser userSolvedProblems

	err = json.Unmarshal(body, &dataUser)
	if err != nil {
		logErr(err)
		return
	}

	var rows []table.Row

	for _, problem := range dataUser.Data {

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", problem.Problem_ID),
			problem.Name,
			problem.Source,
			fmt.Sprintf("%.0f", problem.Score),
		})
	}

	columns := []table.Column{
		{Title: "Problem ID", Width: 7},
		{Title: "Name", Width: 20},
		{Title: "Source", Width: 40},
		{Title: "Score", Width: 7},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	t.SetStyles(table.DefaultStyles())

	p := tea.NewProgram(model{table: t})
	if t, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %s %v", t, err)
	}
}

// extend session
func extendSession() {
	body, err := makeRequest("POST", URL_EXTEND_SESSION, nil, "1")
	if err != nil {
		logErr(err)
		return
	}

	var resp KNResponse
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

func isAdmin(user_id string) bool {
	return userGetDetails(user_id, "isadmin")
}

func setBio(bio string) {
	data := map[string]string{"bio": bio}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	body, err := makeRequest("POST", URL_SELF_SET_BIO, bytes.NewBuffer(jsonData), "2")
	if err != nil {
		logErr(err)
		return
	}

	var dataKN KNResponse
	err = json.Unmarshal(body, &dataKN)
	if err != nil {
		logErr(err)
		return
	}

	if dataKN.Status == "success" {
		fmt.Println("Success! Bio changed!")
	} else {
		fmt.Println("Error: Failed to change bio!")
	}
}

func changeName(newName, password string) {
	data := map[string]string{"newName": newName, "password": password}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	body, err := makeRequest("POST", URL_CHANGE_NAME, bytes.NewBuffer(jsonData), "2")
	if err != nil {
		logErr(err)
		return
	}

	var dataKN KNResponse
	err = json.Unmarshal(body, &dataKN)
	if err != nil {
		logErr(err)
		return
	}

	if dataKN.Status == "success" {
		fmt.Println("Success! Name changed!")
	} else {
		fmt.Println("Error: Failed to change name!")
	}
}

func changePass(oldpass, newpass string) {
	data := map[string]string{"old_password": oldpass, "password": newpass}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	body, err := makeRequest("POST", URL_CHANGE_PASS, bytes.NewBuffer(jsonData), "2")
	if err != nil {
		logErr(err)
		return
	}

	var dataKN KNResponse
	err = json.Unmarshal(body, &dataKN)
	if err != nil {
		logErr(err)
		return
	}

	if dataKN.Status == "success" {
		fmt.Println("Success! password changed! You'll have to login again with your new credentials!")
		logout()
	} else {
		fmt.Println("Error: Failed to change password!")
	}
}

func changeEmail(email, password string) {

	formData := u.Values{}
	formData.Set("email", email)
	formData.Set("password", password)

	data := []byte(formData.Encode())

	body, err := makeRequest("POST", URL_CHANGE_EMAIL, bytes.NewBuffer(data), "1")
	if err != nil {
		logErr(err)
		return
	}

	var dataKN KNResponse
	err = json.Unmarshal(body, &dataKN)
	if err != nil {
		logErr(err)
		return
	}

	if dataKN.Status == "success" {
		fmt.Println("Success! Email changed!")
	} else {
		fmt.Println("Error: Failed to change email!")
	}
}

func resetPass(email string) {
	_, hasToken := readToken()

	if hasToken {
		fmt.Println("You have to be logged out to perform this action!")
		return
	}

	formData := u.Values{}
	formData.Set("email", email)

	data := []byte(formData.Encode())

	body, err := makeRequest("POST", URL_CHANGE_PASS, bytes.NewBuffer(data), "1")
	if err != nil {
		logErr(err)
		return
	}

	var dataKN KNResponse
	err = json.Unmarshal(body, &dataKN)
	if err != nil {
		logErr(err)
		return
	}

	fmt.Println(dataKN.Data)
}

func deleteUser() {
	fmt.Println("Are you sure to delete your account? (Y/N)")
	var resp string
	fmt.Scan(&resp)

	if resp != "Y" {
		return
	}

	fmt.Println("Do you understand that this action is irreversible? (Y/N)")
	fmt.Scan(&resp)

	if resp != "Y" {
		return
	}

	fmt.Println("Do you understand that you'll lose all of your data? (Y/N)")
	fmt.Scan(&resp)

	if resp != "Y" {
		return
	}

	fmt.Println("Okay. Deleting account...")

}

func resendEmail() {
	body, err := makeRequest("POST", URL_RESEND_MAIL, nil, "2")
	if err != nil {
		logErr(err)
		return
	}

	var dataKN KNResponse
	err = json.Unmarshal(body, &dataKN)
	if err != nil {
		logErr(err)
		return
	}

	fmt.Println(dataKN.Data)
}
