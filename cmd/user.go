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
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
)

type UserDetailResponse struct {
	Data struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Admin       bool   `json:"admin"`
		Proposer    bool   `json:"proposer"`
		DisplayName string `json:"display_name"`
	} `json:"data"`
}

type UserSolvedProblems struct {
	Data []struct {
		ProblemId int     `json:"id"`
		Name      string  `json:"name"`
		Source    string  `json:"source_credits"`
		Score     float64 `json:"score_scale"`
	} `json:"data"`
}

// cobra vars

var settingsCmd = &cobra.Command{
	Use:   "settings [command] ...",
	Short: "Modify your account.",
}

var signinCmd = &cobra.Command{
	Use:   "signin [username] [password]",
	Short: "Sign in to your account",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		username, password := loginForm()
		action := func() { login(username, password) }
		if err := spinner.New().Title("Logging in...").Action(action).Run(); err != nil {
			logError(err)
		}
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of your account",
	Run: func(cmd *cobra.Command, args []string) {
		action := func() { logout() }
		if err := spinner.New().Title("Waiting ...").Action(action).Run(); err != nil {
			logError(err)
		}
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
	Short: "Delete your Kilonova account. (Currently not working in API V1)",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		//deleteUser()
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

var amILoggedInCmd = &cobra.Command{
	Use:   "amilogged",
	Short: "Check wether you're logged in or not.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(isCurrentUserLoggedIn())
	},
}

var amIAdminCmd = &cobra.Command{
	Use:   "amiadmin",
	Short: "Check wether you're an admin.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(isAdmin("me"))
	},
}

func init() {
	rootCmd.AddCommand(signinCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(userGetDetailsCmd)
	settingsCmd.AddCommand(extendSessionCmd)
	rootCmd.AddCommand(userSolvedProblemsCmd)
	settingsCmd.AddCommand(setBioCmd)
	settingsCmd.AddCommand(changeNameCmd)
	settingsCmd.AddCommand(changePassCmd)
	settingsCmd.AddCommand(changeEmailCmd)
	settingsCmd.AddCommand(resetPassCmd)
	settingsCmd.AddCommand(deleteUserCmd)
	settingsCmd.AddCommand(resendEmailCmd)
	settingsCmd.AddCommand(amILoggedInCmd)
	settingsCmd.AddCommand(amIAdminCmd)

	rootCmd.AddCommand(settingsCmd)
}

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
		logError(err)
	}

	return username, password
}

func login(username, password string) {
	formData := u.Values{
		"username": {username},
		"password": {password},
	}

	body, err := MakePostRequest(URL_LOGIN, bytes.NewBufferString(formData.Encode()), RequestFormGuest)
	if err != nil {
		logError(fmt.Errorf("login failed: %v", err))
	}

	if !bytes.Contains(body, []byte("success")) {
		logError(fmt.Errorf("login failed: invalid credentials"))
	}

	var response KilonovaResponse
	if err := json.Unmarshal(body, &response); err != nil {
		logError(fmt.Errorf("error parsing response: %v", err))
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		logError(err)
	}
	configDir := filepath.Join(homedir, ".config", "kn-cli")
	err = os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		logError(err)
	}
	tokenFile := filepath.Join(configDir, "token")

	file, err := os.Create(tokenFile)
	if err != nil {
		logError(fmt.Errorf("error creating file: %v", err))
	}
	defer file.Close()

	if err := os.WriteFile(tokenFile, []byte(response.Data), 0644); err != nil {
		logError(fmt.Errorf("error writing auth token to file: %v", err))
	}

	fmt.Println("Login successful!")
}

// logout
func logout() {
	jsonData := []byte(`{"key": "value"}`)
	body, err := MakePostRequest(URL_LOGOUT, bytes.NewBuffer(jsonData), RequestFormAuth)
	if err != nil {
		logError(err)
	}

	if bytes.Contains(body, []byte("success")) {
		fmt.Println("Logged out successfully!")
		homedir, err := os.UserHomeDir()
		if err != nil {
			logError(err)
		}
		configDir := filepath.Join(homedir, ".config", "kn-cli")
		tokenFile := filepath.Join(configDir, "token")
		_ = os.Remove(tokenFile)
	} else {
		log.Println("Logout failed: You must be logged in to do this!")
	}
}

// get info about user
func getBio(name string) string {
	res, err := http.Get(fmt.Sprintf("https://kilonova.ro/profile/%s", name))
	if err != nil {
		logError(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logError(err)
	}

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

	// Fetch user data
	body, err := MakeGetRequest(url, nil, RequestFormGuest)
	if err != nil {
		logError(fmt.Errorf("error fetching user details: %w", err))
		return false
	}

	var dataUser UserDetailResponse
	if err := json.Unmarshal(body, &dataUser); err != nil {
		logError(fmt.Errorf("error unmarshalling user data: %w", err))
		return false
	}

	// Set a default value for DisplayName if empty
	if dataUser.Data.DisplayName == "" {
		dataUser.Data.DisplayName = "-"
	}

	// Handle use_case
	switch use_case {
	case "isadmin":
		return dataUser.Data.Admin
	default:
		// Generate user details output using a template
		userTemplate := `ID: {{.Id}}
Name: {{.Name}}
A.K.A: {{.DisplayName}}
Bio: {{.Bio}}
Admin: {{.Admin}}
Proposer: {{.Proposer}}
`

		// Prepare user bio
		bio := getBio(dataUser.Data.Name)

		// Create a map for template data
		userData := struct {
			Id          int
			Name        string
			DisplayName string
			Bio         string
			Admin       bool
			Proposer    bool
		}{
			Id:          dataUser.Data.Id,
			Name:        dataUser.Data.Name,
			DisplayName: dataUser.Data.DisplayName,
			Bio:         bio,
			Admin:       dataUser.Data.Admin,
			Proposer:    dataUser.Data.Proposer,
		}

		// Parse and execute the template
		tmpl, err := template.New("userDetails").Parse(userTemplate)
		if err != nil {
			logError(fmt.Errorf("error parsing template: %w", err))
			return false
		}

		if err := tmpl.Execute(os.Stdout, userData); err != nil {
			logError(fmt.Errorf("error executing template: %w", err))
			return false
		}

		return false
	}
}

func userGetSolvedProblems(userId string) {
	var url string
	if userId == "me" {
		url = URL_SELF_PROBLEMS
	} else {
		url = fmt.Sprintf(URL_USER_PROBLEMS, userId)
	}

	body, err := MakeGetRequest(url, nil, RequestFormGuest)
	if err != nil {
		logError(fmt.Errorf("error fetching solved problems: %w", err))
		return
	}

	var dataUser UserSolvedProblems
	if err := json.Unmarshal(body, &dataUser); err != nil {
		logError(fmt.Errorf("error unmarshalling solved problems: %w", err))
		return
	}

	rows := prepareTableRows(dataUser)

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

	p := tea.NewProgram(&Model{table: t}, tea.WithAltScreen())
	if t, err := p.Run(); err != nil {
		logError(fmt.Errorf("error running program: %v %v", t, err))
	}
}

func prepareTableRows(dataUser UserSolvedProblems) []table.Row {
	var rows []table.Row
	for _, problem := range dataUser.Data {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", problem.ProblemId),
			problem.Name,
			problem.Source,
			fmt.Sprintf("%.0f", problem.Score),
		})
	}
	return rows
}

func isCurrentUserLoggedIn() bool {
	_, hasToken := readToken()
	return hasToken
}

func extendSession() {
	body, err := MakePostRequest(URL_EXTEND_SESSION, nil, RequestFormAuth)
	if err != nil {
		logError(err)
	}

	var resp KilonovaResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		logError(fmt.Errorf("error unmarshalling response: %s", err))
	}

	if resp.Status == "success" {
		parsedTime, err := time.Parse(time.RFC3339Nano, resp.Data)
		if err != nil {
			logError(err)
		}
		formattedTime := parsedTime.Format("2006-01-02 15:04:05")
		fmt.Println("Your session has been extended until ", formattedTime)
	} else {
		fmt.Println(resp.Data)
	}

}

func isAdmin(userId string) bool {
	return userGetDetails(userId, "isadmin")
}

func setBio(bio string) {
	payload := map[string]string{"bio": bio}
	resp, err := PostJSON[KilonovaResponse](URL_SELF_SET_BIO, payload)
	if err != nil {
		logError(err)
	}

	if resp.Status == "success" {
		fmt.Println("Success! Bio changed!")
	} else {
		fmt.Println("Error: Failed to change bio!")
	}
}

func changeName(newName, password string) {
	payload := map[string]string{
		"newName":  newName,
		"password": password,
	}
	resp, err := PostJSON[KilonovaResponse](URL_CHANGE_NAME, payload)
	if err != nil {
		logError(err)
	}

	if resp.Status == "success" {
		fmt.Println("Success! Name changed!")
	} else {
		logError(fmt.Errorf("failed to change name"))
	}
}

func changePass(oldPass, newPass string) {
	payload := map[string]string{
		"old_password": oldPass,
		"password":     newPass,
	}
	resp, err := PostJSON[KilonovaResponse](URL_CHANGE_PASS, payload)
	if err != nil {
		logError(err)
	}

	if resp.Status == "success" {
		fmt.Println("Success! Password changed! You'll need to login again.")
		logout()
	} else {
		logError(fmt.Errorf("failed to change password"))
	}
}

func changeEmail(email, password string) {
	formData := u.Values{}
	formData.Set("email", email)
	formData.Set("password", password)

	body, err := MakePostRequest(URL_CHANGE_EMAIL, bytes.NewBufferString(formData.Encode()), RequestFormAuth)
	if err != nil {
		logError(err)
	}

	var res KilonovaResponse
	if err := json.Unmarshal(body, &res); err != nil {
		logError(err)
	}

	if res.Status == "success" {
		fmt.Println("Success! Email changed!")
	} else {
		logError(fmt.Errorf("failed to change email"))
	}
}

func resetPass(email string) {
	if _, loggedIn := readToken(); loggedIn {
		fmt.Println("You must be logged out to reset your password.")
		return
	}

	form := u.Values{}
	form.Set("email", email)

	body, err := MakePostRequest(URL_CHANGE_PASS, bytes.NewBufferString(form.Encode()), RequestFormAuth)
	if err != nil {
		logError(err)
	}

	var res KilonovaResponse
	if err := json.Unmarshal(body, &res); err != nil {
		logError(err)
	}

	fmt.Println(res.Data)
}

func resendEmail() {
	body, err := MakePostRequest(URL_RESEND_MAIL, nil, RequestFormAuth)
	if err != nil {
		logError(err)
	}

	var res KilonovaResponse
	if err := json.Unmarshal(body, &res); err != nil {
		logError(err)
	}

	fmt.Println(res.Data)
}

/*
func deleteUser() {
	fmt.Print("Are you sure to delete your account? (Y/N) ")
	var resp string
	fmt.Scan(&resp)

	if resp != "Y" {
		return
	}

	fmt.Print("Do you understand that this action is irreversible? (Y/N) ")
	fmt.Scan(&resp)

	if resp != "Y" {
		return
	}

	fmt.Print("Do you understand that you'll lose all of your data? (Y/N) ")
	fmt.Scan(&resp)

	if resp != "Y" {
		return
	}

	fmt.Println("Deleting account...")
	fmt.Println("Currently the KN api does not support deleting an account if you're not an admin.")

	if !isAdmin("me") {
		return
	}

	body, err := MakePostRequest(URL_CHANGE_EMAIL, nil, RequestFormAuth)
	if err != nil {
		LogError(err)
		return
	}

	var dataKN KilonovaResponse
	err = json.Unmarshal(body, &dataKN)
	if err != nil {
		LogError(err)
		return
	}

	fmt.Println(dataKN.Data)
}
*/
