// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package utility

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

func DecodeBase64Text(EncodedText string) (string, error) {
	DecodedText, err := b64.StdEncoding.DecodeString(EncodedText)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 text: %w", err)
	}
	return string(DecodedText), nil
}

func ValidateBoolean(input string) (bool, error) {
	if input != BOOLTRUE && input != BOOLFALSE {
		return false, fmt.Errorf("value must be either 'true' or 'false'")
	}
	return input == BOOLTRUE, nil
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

// Read Token Function

func ReadToken() (string, bool) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		LogError(fmt.Errorf("failed to get user home directory: %w", err))
	}

	tokenPath := filepath.Join(homedir, CONFIGFOLDER, KNCLIFOLDER, TOKENFILENAME)
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
