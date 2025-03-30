package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
	userAgent           = "KilonovaCLIClient/1.0"
	help                = "Kilonova CLI - ver 0.1.0\n\n-signin <USERNAME> <PASSWORD>\n-langs <ID>\n-search <PROBLEM ID or NAME>\n-submit <PROBLEM ID> <LANGUAGE> <solution>\n-submissions <ID>\n-statement <PROBLEM ID> <RO or EN>\n-logout"
)

func readToken() string {
	token, err := os.ReadFile("token")
	if err != nil {
		log.Fatal("Could not read session ID. Make sure you are logged in!")
	}
	return string(token)
}

func logErr(err error) {
	log.Fatal(err)
}

func makeRequest(method, url string, body io.Reader, use_case string) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	token := readToken()

	req.Header.Set("User-Agent", userAgent)
	switch {
	case use_case == "1":
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	case use_case == "2":
		req.Header.Set("Content-Type", "application/json")
	default:
		if use_case != "0" {
			req.Header.Set("Content-Type", use_case)
		}
	}
	req.Header.Set("Authorization", token)

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
