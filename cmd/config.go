package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	URL_LOGIN           = "https://kilonova.ro/api/auth/login"
	URL_LOGOUT          = "https://kilonova.ro/api/auth/logout"
	URL_SEARCH          = "https://kilonova.ro/api/problem/search"
	URL_PROBLEM         = "https://kilonova.ro/api/problem/%s/"
	URL_SELF            = "https://kilonova.ro/api/user/self/"
	URL_LANGS_PB        = "https://kilonova.ro/api/problem/%s/languages"
	URL_SUBMIT          = "https://kilonova.ro/api/submissions/submit"
	URL_SUBMISSION_LIST = "https://kilonova.ro/api/submissions/get?ascending=false&limit=500&offset=0&ordering=id&problem_id=%s&user_id=%d"
	STAT_FILENAME_RO    = "statement-ro.md"
	STAT_FILENAME_EN    = "statement-en.md"
	URL_STATEMENT       = "https://kilonova.ro/api/problem/%s/get/attachmentByName/%s"
	userAgent           = "KilonovaCLIClient/0.1"
	help                = "Kilonova CLI - ver 0.1.0\n\n-signin <USERNAME> <PASSWORD>\n-langs <ID>\n-search <PROBLEM ID or NAME>\n-submit <PROBLEM ID> <LANGUAGE> <solution>\n-submissions <ID>\n-statement <PROBLEM ID> <RO or EN>\n-logout"
)

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
		case "q": // Quit the program
			return m, tea.Quit
		case "esc": // Also quit
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
