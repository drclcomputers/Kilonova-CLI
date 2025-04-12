// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package internal

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"path/filepath"
)

// Utility Functions

func ProblemExists(ID string) bool {
	url := fmt.Sprintf(URL_PROBLEM, ID)
	body, _ := MakeGetRequest(url, nil, RequestNone)
	if string(body) == "notfound" {
		return false
	}
	return true
}

func FileExists(filename string) bool {
	FileName := filepath.Join(GetConfigDir(), filename)
	if _, err := os.Stat(FileName); err == nil {
		return true
	}
	return false
}

func DecodeBase64Text(EncodedText string) (string, error) {
	DecodedText, err := b64.StdEncoding.DecodeString(EncodedText)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 text: %w", err)
	}
	return string(DecodedText), nil
}

func EncodeBase64Text(DecodedText string) (string, error) {
	EncodedText := b64.StdEncoding.EncodeToString([]byte(DecodedText))
	return EncodedText, nil
}

func ValidateBoolean(input string) (bool, error) {
	if input != BOOLTRUE && input != BOOLFALSE {
		return false, fmt.Errorf("value must be either 'true' or 'false'")
	}
	return input == BOOLTRUE, nil
}

func ValidateInt(input string) (int, error) {
	var err error
	var nr int
	if nr, err = strconv.Atoi(input); err == nil {
		return nr, nil
	}
	return 0, err
}

func ParseTime(timeStr string) (string, error) {
	parsedTime, err := time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse time %q: %w", timeStr, err)
	}
	return parsedTime.Format("2006-01-02 15:04:05"), nil
}

func GetUserID() string {
	ResponseBody, err := MakeGetRequest(URL_SELF, nil, RequestFormAuth)
	if err != nil {
		LogError(fmt.Errorf("failed to retrieve user info: %w", err))
		return ""
	}

	var user UserId
	if err := json.Unmarshal(ResponseBody, &user); err != nil {
		LogError(fmt.Errorf("failed to parse user ID from response: %w", err))
		return ""
	}

	return strconv.Itoa(user.Data.ID)
}

func GetAProblemName(problemID string) (string, error) {
	url := fmt.Sprintf(URL_PROBLEM, problemID)
	ResponseBody, err := MakeGetRequest(url, nil, RequestNone)
	if err != nil {
		return "", fmt.Errorf("failed to fetch problem details: %w", err)
	}

	var info ProblemInfo
	if err := json.Unmarshal(ResponseBody, &info); err != nil {
		return "", fmt.Errorf("failed to parse problem info: %w", err)
	}
	return info.Data.Name, nil
}

func GetConfigDir() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		LogError(err)
		return "error"
	}
	configDir := filepath.Join(homedir, CONFIGFOLDER, KNCLIFOLDER)
	err = os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		LogError(err)
		return "error"
	}
	return configDir
}

// Read Token Function

func ReadToken() (string, bool) {
	tokenPath := filepath.Join(GetConfigDir(), TOKENFILENAME)

	data, err := os.ReadFile(tokenPath)
	if err != nil {
		return "", false
	}

	decryptText, err := Decrypt(string(bytes.TrimSpace(data)))
	if err != nil {
		LogError(fmt.Errorf("failed to decrypt token: %w", err))
		return "", false
	}

	return decryptText, true
}

// Log Error Function

func LogError(err error) {
	log.Fatalf("%s%s%s", RED, err.Error(), WHITE)
}
