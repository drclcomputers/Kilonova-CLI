package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	URL_LOGIN                      = "https://kilonova.ro/api/auth/login"
	URL_LOGOUT                     = "https://kilonova.ro/api/auth/logout"
	URL_EXTEND_SESSION             = "https://kilonova.ro/api/auth/extendSession"
	URL_SEARCH                     = "https://kilonova.ro/api/problem/search"
	URL_PROBLEM                    = "https://kilonova.ro/api/problem/%s/"
	URL_SELF                       = "https://kilonova.ro/api/user/self/"
	URL_SELF_PROBLEMS              = "https://kilonova.ro/api/user/self/solvedProblems"
	URL_SELF_SET_BIO               = "https://kilonova.ro/api/user/self/setBio"
	URL_CHANGE_EMAIL               = "https://kilonova.ro/api/user/changeEmail"
	URL_CHANGE_PASS                = "https://kilonova.ro/api/user/changePassword"
	URL_UPDATE_NAME                = "https://kilonova.ro/api/user/updateName"
	URL_USER                       = "https://kilonova.ro/api/user/byID/%s"
	URL_USER_PROBLEMS              = "https://kilonova.ro/api/user/byID/%s/solvedProblems"
	URL_LANGS_PB                   = "https://kilonova.ro/api/problem/%s/languages"
	URL_SUBMIT                     = "https://kilonova.ro/api/submissions/submit"
	URL_SUBMISSION_LIST            = "https://kilonova.ro/api/submissions/get?ascending=false&limit=50&offset=%d&ordering=id&problem_id=%s&user_id=%s"
	URL_SUBMISSION_LIST_NO_FILTER  = "https://kilonova.ro/api/submissions/get?ascending=false&limit=50&offset=%d&ordering=id"
	URL_SUBMISSION_LIST_NO_PROBLEM = "https://kilonova.ro/api/submissions/get?ascending=false&limit=50&offset=%d&ordering=id&user_id=%s"
	URL_SUBMISSION_LIST_NO_USER    = "https://kilonova.ro/api/submissions/get?ascending=false&limit=50&offset=%d&ordering=id&problem_id=%s"
	STAT_FILENAME_RO               = "statement-ro.md"
	STAT_FILENAME_EN               = "statement-en.md"
	URL_STATEMENT                  = "https://kilonova.ro/api/problem/%s/get/attachmentByName/%s"
	userAgent                      = "KilonovaCLIClient/0.1"
	help                           = "Kilonova CLI - ver 0.1.0\n\n-signin <USERNAME> <PASSWORD>\n-langs <ID>\n-search <PROBLEM ID or NAME>\n-submit <PROBLEM ID> <LANGUAGE> <solution>\n-submissions <ID>\n-statement <PROBLEM ID> <RO or EN>\n-logout"
)

// TEXT MODEL

type textModel struct {
	viewport viewport.Model
}

func (m textModel) Init() tea.Cmd {
	return nil
}

func (m textModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
		case "up", "k":
			m.viewport.LineUp(1)
		case "down", "j":
			m.viewport.LineDown(1)
		}
	}
	return m, nil
}

func (m textModel) View() string {
	style := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1)
	return style.Render(m.viewport.View()) + "\n(Use ↑/↓ to scroll, 'q' to quit)"
}

func newTextModel(text string) textModel {
	vp := viewport.New(80, 25)
	vp.SetContent(text)
	return textModel{viewport: vp}
}

// TABLE MODEL

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "esc":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return lipgloss.NewStyle().Margin(1, 2).Render(m.table.View()) + "\n(Use ↑/↓ to navigate, 'q' to quit)"
}

// Other functions

type userid struct {
	Data struct {
		ID int `json:"id"`
	} `json:"data"`
}

func getUserID() string {
	//get user id

	jsonData := []byte(`{"key": "value"}`)
	body, err := makeRequest("GET", URL_SELF, bytes.NewBuffer(jsonData), "1")
	if err != nil {
		logErr(err)
		return ""
	}

	var data userid
	if err := json.Unmarshal(body, &data); err != nil {
		logErr(err)
		return ""
	}

	return strconv.Itoa(data.Data.ID)
}

func readToken() (string, bool) {
	token, err := os.ReadFile("token")
	if err != nil {
		return "", false
	}
	return string(bytes.TrimSpace(token)), true
}

func logErr(err error) {
	log.Fatal(err)
}

func makeRequest(method, url string, body io.Reader, use_case string) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	token, hasToken := readToken()

	req.Header.Set("User-Agent", userAgent)
	switch {
	case use_case == "1" || use_case == "3":
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	case use_case == "2":
		req.Header.Set("Content-Type", "application/json")
	default:
		if use_case != "0" {
			req.Header.Set("Content-Type", use_case)
		}
	}

	if hasToken {
		req.Header.Set("Authorization", token)
	} else {
		if use_case == "1" {
			log.Fatal("You must be authenticated to do this")
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return data, nil
}
