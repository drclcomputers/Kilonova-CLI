// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

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

	"path/filepath"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	API_URL                        = "https://kilonova.ro/api/"
	URL_LOGIN                      = API_URL + "auth/login"
	URL_LOGOUT                     = API_URL + "auth/logout"
	URL_EXTEND_SESSION             = API_URL + "auth/extendSession"
	URL_SEARCH                     = API_URL + "problem/search"
	URL_PROBLEM                    = API_URL + "problem/%s/"
	URL_SELF                       = API_URL + "user/self/"
	URL_SELF_PROBLEMS              = API_URL + "user/self/solvedProblems"
	URL_SELF_SET_BIO               = API_URL + "user/self/setBio"
	URL_CHANGE_EMAIL               = API_URL + "user/changeEmail"
	URL_CHANGE_PASS                = API_URL + "user/changePassword"
	URL_CHANGE_NAME                = API_URL + "user/updateName"
	URL_USER                       = API_URL + "user/byID/%s"
	URL_USER_PROBLEMS              = API_URL + "user/byID/%s/solvedProblems"
	URL_LANGS_PB                   = API_URL + "problem/%s/languages"
	URL_SUBMIT                     = API_URL + "submissions/submit"
	URL_LATEST_SUBMISSION          = API_URL + "submissions/getByID?id=%d"
	URL_SUBMISSION_LIST            = API_URL + "submissions/get?ascending=false&limit=50&offset=%d&ordering=id&problem_id=%s&user_id=%s"
	URL_SUBMISSION_LIST_NO_FILTER  = API_URL + "submissions/get?ascending=false&limit=50&offset=%d&ordering=id"
	URL_SUBMISSION_LIST_NO_PROBLEM = API_URL + "submissions/get?ascending=false&limit=50&offset=%d&ordering=id&user_id=%s"
	URL_SUBMISSION_LIST_NO_USER    = API_URL + "submissions/get?ascending=false&limit=50&offset=%d&ordering=id&problem_id=%s"
	STAT_FILENAME_RO               = "statement-ro.md"
	STAT_FILENAME_EN               = "statement-en.md"
	URL_STATEMENT                  = API_URL + "problem/%s/get/attachmentByName/%s"
	URL_ASSETS                     = "https://kilonova.ro/assets/problem/%s/problemArchive?tests=true&attachments=true&private_attachments=false&details=true&tags=true&editors=true&submissions=false&all_submissions=false"
	URL_RESEND_MAIL                = API_URL + "user/resendEmail"
	URL_DELETE_USER                = API_URL + "user/moderation/deleteUser"
	userAgent                      = "KilonovaCLIClient/0.1"
	XMLCBPStruct                   = `<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>
<CodeBlocks_project_file>
	<FileVersion major="1" minor="6" />
	<Project>
		<Option title="%s" />
		<Option pch_mode="2" />
		<Option compiler="gcc" />
		<Build>
			<Target title="Debug">
				<Option output="bin/Debug/%s" prefix_auto="1" extension_auto="1" />
				<Option object_output="obj/Debug/" />
				<Option type="1" />
				<Option compiler="gcc" />
				<Compiler>
					<Add option="-g" />
				</Compiler>
			</Target>
			<Target title="Release">
				<Option output="bin/Release/%s" prefix_auto="1" extension_auto="1" />
				<Option object_output="obj/Release/" />
				<Option type="1" />
				<Option compiler="gcc" />
				<Compiler>
					<Add option="-O2" />
				</Compiler>
				<Linker>
					<Add option="-s" />
				</Linker>
			</Target>
		</Build>
		<Compiler>
			<Add option="-Wall" />
			<Add option="-fexceptions" />
		</Compiler>
		<Unit filename="Source.cpp" />
		<Extensions />
	</Project>
</CodeBlocks_project_file>
`
	CMAKEStruct = `cmake_minimum_required(VERSION 3.10)
project(%s VERSION 1.0 LANGUAGES CXX)
add_executable(%s Source.cpp)`
)

var helloWorldprog = []string{
	`#include <stdio.h>
int main() {
    printf("Hello, World!\n");
    return 0;
}`,
	`#include <iostream>
int main() {
    std::cout << "Hello, World!" << std::endl;
    return 0;
}`,
	`package main
import "fmt"
func main() {
    fmt.Println("Hello, World!")
}`,
	`fun main() {
    println("Hello, World!")
}`,
	`console.log("Hello, World!");`,
	`program HelloWorld;
begin
    writeln('Hello, World!');
end.
`,
	`<?php
echo "Hello, World!\n";
?>
`, `print("Hello, World!")
`, `fn main() {
    println!("Hello, World!");
}
`,
	`#include <stdio.h>

int main() {
    FILE *file = fopen("example.txt", "r");
    char ch;
    while ((ch = fgetc(file)) != EOF) {
        putchar(ch);
    }
    return 0;
}`,
	`#include <iostream>
#include <fstream>

int main() {
    std::ifstream file("example.txt");
    char ch;
    while (file.get(ch)) {
        std::cout << ch;
    }
    return 0;
}`,
}

type RequestType int

const (
	RequestNone          RequestType = iota
	RequestFormAuth                  // 1
	RequestJSON                      // 2
	RequestFormGuest                 // 3
	RequestDownloadZip               // 4
	RequestMultipartForm             // 5
)

// TEXT MODEL

type TextModel struct {
	viewport viewport.Model
}

func (m *TextModel) Init() tea.Cmd {
	return nil
}

func (m *TextModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "q", "esc":
			return m, tea.Quit
		case "up", "k":
			m.viewport.LineUp(1)
		case "down", "j":
			m.viewport.LineDown(1)
		}
	}
	return m, cmd
}

func (m *TextModel) View() string {
	style := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1)
	footer := "\n(Use ↑/↓ to scroll, 'q' to quit)"
	return style.Render(m.viewport.View()) + footer
}

func NewTextModel(text string) *TextModel {
	vp := viewport.New(80, 25)
	vp.SetContent(text)
	return &TextModel{viewport: vp}
}

type Model struct {
	table table.Model
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "q", "esc":
			return m, tea.Quit
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	footer := "\n(Use ↑/↓ to navigate, 'q' to quit, 'enter' to get the statement)"
	return lipgloss.NewStyle().Margin(1, 2).Render(m.table.View()) + footer
}

type UserId struct {
	Data struct {
		ID int `json:"id"`
	} `json:"data"`
}

type RawKilonovaResponse struct {
	Status string          `json:"status"`
	Data   json.RawMessage `json:"data"`
}

type KilonovaResponse struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

func getUserID() string {
	body, err := MakeGetRequest(URL_SELF, nil, RequestFormAuth)
	if err != nil {
		logError(fmt.Errorf("failed to retrieve user info: %w", err))
		return ""
	}

	var user UserId
	if err := json.Unmarshal(body, &user); err != nil {
		logError(fmt.Errorf("failed to parse user ID from response: %w", err))
		return ""
	}

	return strconv.Itoa(user.Data.ID)
}

func readToken() (string, bool) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		logError(fmt.Errorf("failed to get user home directory: %w", err))
		return "", false
	}

	tokenPath := filepath.Join(homedir, ".config", "kn-cli", "token")
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		logError(fmt.Errorf("failed to read token file: %w", err))
		return "", false
	}

	return string(bytes.TrimSpace(data)), true
}

func logError(err error) {
	log.Fatalf("\033[31m%v\033[0m\n", err)
}

func MakeRequest(method, url string, body io.Reader, reqType RequestType, contentType ...string) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", userAgent)
	token, hasToken := readToken()

	switch reqType {
	case RequestFormAuth, RequestFormGuest:
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	case RequestJSON:
		req.Header.Set("Content-Type", "application/json")
	case RequestDownloadZip:
		req.Header.Set("Content-Type", "application/zip")
		req.Header.Set("Accept", "application/zip")
		fmt.Println("Trying to obtain archive...")
		cookie := &http.Cookie{
			Name:  "kn-sessionid",
			Value: token,
		}
		req.AddCookie(cookie)
	case RequestMultipartForm:
		if len(contentType) > 0 {
			req.Header.Set("Content-Type", contentType[0])
		} else {
			logError(fmt.Errorf("missing content type for multipart form request"))
		}
	}

	if hasToken {
		req.Header.Set("Authorization", token)
	} else if reqType == RequestFormAuth || reqType == RequestDownloadZip {
		logError(fmt.Errorf("you must be authenticated to do this"))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var res RawKilonovaResponse
		if err := json.Unmarshal(data, &res); err != nil {
			logError(err)
		}
		logError(fmt.Errorf("error: %s", string(res.Data)))
	}

	return data, nil
}

func MakeGetRequest(url string, body io.Reader, reqType RequestType, contentType ...string) ([]byte, error) {
	return MakeRequest("GET", url, body, reqType)
}

func MakePostRequest(url string, body io.Reader, reqType RequestType, contentType ...string) ([]byte, error) {
	return MakeRequest("POST", url, body, reqType)
}

func PostJSON[T any](url string, payload any) (T, error) {
	var result T

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return result, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	body, err := MakePostRequest(url, bytes.NewBuffer(jsonData), RequestJSON)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}
