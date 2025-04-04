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
	URL_CHANGE_NAME                = "https://kilonova.ro/api/user/updateName"
	URL_USER                       = "https://kilonova.ro/api/user/byID/%s"
	URL_USER_PROBLEMS              = "https://kilonova.ro/api/user/byID/%s/solvedProblems"
	URL_LANGS_PB                   = "https://kilonova.ro/api/problem/%s/languages"
	URL_SUBMIT                     = "https://kilonova.ro/api/submissions/submit"
	URL_LATEST_SUBMISSION          = "https://kilonova.ro/api/submissions/getByID?id=%d"
	URL_SUBMISSION_LIST            = "https://kilonova.ro/api/submissions/get?ascending=false&limit=50&offset=%d&ordering=id&problem_id=%s&user_id=%s"
	URL_SUBMISSION_LIST_NO_FILTER  = "https://kilonova.ro/api/submissions/get?ascending=false&limit=50&offset=%d&ordering=id"
	URL_SUBMISSION_LIST_NO_PROBLEM = "https://kilonova.ro/api/submissions/get?ascending=false&limit=50&offset=%d&ordering=id&user_id=%s"
	URL_SUBMISSION_LIST_NO_USER    = "https://kilonova.ro/api/submissions/get?ascending=false&limit=50&offset=%d&ordering=id&problem_id=%s"
	STAT_FILENAME_RO               = "statement-ro.md"
	STAT_FILENAME_EN               = "statement-en.md"
	URL_STATEMENT                  = "https://kilonova.ro/api/problem/%s/get/attachmentByName/%s"
	URL_ASSETS                     = "https://kilonova.ro/assets/problem/%s/problemArchive?tests=true&attachments=true&private_attachments=false&details=true&tags=true&editors=true&submissions=false&all_submissions=false"
	URL_RESEND_MAIL                = "https://kilonova.ro/api/user/resendEmail"
	URL_DELETE_USER                = "https://kilonova.ro/api/user/moderation/deleteUser"
	userAgent                      = "KilonovaCLIClient/0.1"
	help                           = "Kilonova CLI - ver 0.1.5\n\n-signin <USERNAME> <PASSWORD>\n-langs <ID>\n-search <PROBLEM ID or NAME>\n-submit <PROBLEM ID> <LANGUAGE> <solution>\n-submissions <ID>\n-statement <PROBLEM ID> <RO or EN>\n-logout"
	XMLCBPStruct                   = `<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>
<CodeBlocks_project_file>
	<FileVersion major="1" minor="6" />
	<Project>
		<Option title="%s" />
		<Option pch_mode="2" />
		<Option compiler="gcc" />
		<Build>
			<Target title="Debug">
				<Option output="bin/Debug/probleme" prefix_auto="1" extension_auto="1" />
				<Option object_output="obj/Debug/" />
				<Option type="1" />
				<Option compiler="gcc" />
				<Compiler>
					<Add option="-g" />
				</Compiler>
			</Target>
			<Target title="Release">
				<Option output="bin/Release/probleme" prefix_auto="1" extension_auto="1" />
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
}

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
	return lipgloss.NewStyle().Margin(1, 2).Render(m.table.View()) + "\n(Use ↑/↓ to navigate, 'q' to quit, 'enter' to get the statement)"
}

// Other functions

type userid struct {
	Data struct {
		ID int `json:"id"`
	} `json:"data"`
}

type KNResponseRaw struct {
	Status string          `json:"status"`
	Data   json.RawMessage `json:"data"`
}

type KNResponse struct {
	Status string `json:"status"`
	Data   string `json:"data"`
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
	homedir, err := os.UserHomeDir()
	if err != nil {
		logErr(err)
	}
	configDir := filepath.Join(homedir, ".config", "kn-cli")
	tokenFile := filepath.Join(configDir, "token")

	token, err := os.ReadFile(tokenFile)
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
	case use_case == "1" || use_case == "3": //form: 1 - logged, 3 - guest
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	case use_case == "2": //json
		req.Header.Set("Content-Type", "application/json")
	case use_case == "4": //download zip
		req.Header.Set("Content-Type", "application/zip")
		fmt.Println("Trying to obtain archive...")
	default:
		if use_case != "0" { //other use cases
			req.Header.Set("Content-Type", use_case)
		}
	}

	if hasToken {
		req.Header.Set("Authorization", token)
	} else {
		if use_case == "1" || use_case == "4" {
			log.Fatal("You must be authenticated to do this")
		}
	}

	if use_case == "4" {
		req.Header.Set("Accept", "application/zip")
		cookie := &http.Cookie{
			Name:  "kn-sessionid",
			Value: token,
		}
		req.AddCookie(cookie)
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
		var respKN KNResponseRaw
		if err := json.Unmarshal(data, &respKN); err != nil {
			logErr(err)
		}
		fmt.Println("Error: " + string(respKN.Data))
		os.Exit(1)
	}

	return data, nil
}
