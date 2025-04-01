package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	u "net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
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
	Use:   "user [User ID or me (get self ID)]",
	Short: "Get details about a user.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userGetDetails(args[0])
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

func init() {
	rootCmd.AddCommand(signinCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(userGetDetailsCmd)
	rootCmd.AddCommand(extendSessionCmd)
	rootCmd.AddCommand(userSolvedProblemsCmd)
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

func userGetDetails(user_id string) {
	var url string
	if user_id == "me" {
		url = URL_SELF
	} else {
		url = fmt.Sprintf(URL_USER, user_id)
	}

	body, err := makeRequest("GET", url, nil, "3")
	if err != nil {
		logErr(err)
		return
	}

	var dataUser userDetailResp

	err = json.Unmarshal(body, &dataUser)
	if err != nil {
		logErr(err)
		return
	}

	if dataUser.Data.DisplayName == "" {
		dataUser.Data.DisplayName = "-"
	}

	fmt.Printf("ID: %d\nName: %s\nA.K.A: %s\nAdmin: %t\nProposer: %t\n\n",
		dataUser.Data.Id, dataUser.Data.Name, dataUser.Data.DisplayName,
		dataUser.Data.Admin, dataUser.Data.Proposer)
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
