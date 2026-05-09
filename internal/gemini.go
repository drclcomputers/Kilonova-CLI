// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/genai"
)

const (
	GeminiModel       = "gemini-3-flash-preview"
	GeminiModelPro    = "gemini-2.5-pro"
	GeminiAPIKeyEnv   = "GEMINI_API_KEY"
	GeminiTimeout     = 30 * time.Second
	GeminiMaxTokens   = 8192
	GeminiTemperature = 0.3
)

var geminiClient *genai.Client

// InitGemini initializes the Gemini client using the GEMINI_API_KEY environment variable.
func InitGemini() (*genai.Client, error) {
	if geminiClient != nil {
		return geminiClient, nil
	}
	apiKey := os.Getenv(GeminiAPIKeyEnv)
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set. Get your key at https://aistudio.google.com/apikey")
	}
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Gemini client: %w", err)
	}
	geminiClient = client
	return geminiClient, nil
}

// GeminiGenerate sends a prompt to Gemini and returns the response text.
func GeminiGenerate(systemPrompt, userPrompt string) (string, error) {
	client, err := InitGemini()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), GeminiTimeout)
	defer cancel()

	fullPrompt := userPrompt
	if systemPrompt != "" {
		fullPrompt = systemPrompt + "\n\n" + userPrompt
	}

	temp := float32(GeminiTemperature)
	resp, err := client.Models.GenerateContent(
		ctx,
		GeminiModel,
		genai.Text(fullPrompt),
		&genai.GenerateContentConfig{
			MaxOutputTokens: int32(GeminiMaxTokens),
			Temperature:     &temp,
		},
	)
	if err != nil {
		return "", fmt.Errorf("Gemini API error: %w", err)
	}

	var result strings.Builder
	for _, candidate := range resp.Candidates {
		for _, part := range candidate.Content.Parts {
			if part.Text != "" {
				result.WriteString(part.Text)
			}
		}
	}

	if result.Len() == 0 {
		return "", fmt.Errorf("empty response from Gemini")
	}

	return result.String(), nil
}

// CacheEntry defines the structure for cached problem context.
type CacheEntry struct {
	InfoText      string `json:"info_text"`
	StatementText string `json:"statement_text"`
}

type problemInfoResult struct {
	info ProblemInfo
	err  error
}

type statementResult struct {
	raw []byte
	err error
}

// GetProblemContext fetches problem details and statement for Gemini context, utilizing a local cache.
func GetProblemContext(problemID string) (infoText, statementText string, err error) {
	home, err := os.UserHomeDir()
	if err == nil {
		cachePath := filepath.Join(home, ".kncli", "cache", fmt.Sprintf("%s.json", problemID))
		if data, err := os.ReadFile(cachePath); err == nil {
			var entry CacheEntry
			if err := json.Unmarshal(data, &entry); err == nil {
				return entry.InfoText, entry.StatementText, nil
			}
		}
	}

	infoCh := make(chan problemInfoResult, 1)
	statementCh := make(chan statementResult, 1)

	go func() {
		url := fmt.Sprintf(URL_PROBLEM, problemID)
		res, err := MakeGetRequest(url, nil, RequestNone)
		if err != nil {
			infoCh <- problemInfoResult{err: fmt.Errorf("failed to fetch problem %s: %w", problemID, err)}
			return
		}

		var info ProblemInfo
		if err := json.Unmarshal(res, &info); err != nil {
			infoCh <- problemInfoResult{err: fmt.Errorf("failed to parse problem info: %w", err)}
			return
		}
		infoCh <- problemInfoResult{info: info}
	}()

	go func() {
		statementURL := fmt.Sprintf(URL_STATEMENT, problemID, STAT_FILENAME_EN)
		res, err := MakeGetRequest(statementURL, nil, RequestNone)
		if err != nil {
			statementURL = fmt.Sprintf(URL_STATEMENT, problemID, STAT_FILENAME_RO)
			res, err = MakeGetRequest(statementURL, nil, RequestNone)
			if err != nil {
				statementCh <- statementResult{err: fmt.Errorf("failed to fetch statement (EN/RO) for %s: %w", problemID, err)}
				return
			}
		}
		statementCh <- statementResult{raw: res}
	}()

	infoResult := <-infoCh
	statementResult := <-statementCh
	if infoResult.err != nil {
		return "", "", infoResult.err
	}
	if statementResult.err != nil {
		// Non-fatal: provide partial context
		statementText = "(Statement could not be fetched)"
	} else {
		// Decode statement
		type StatementResp struct {
			Status string `json:"status"`
			Data   struct {
				Data string `json:"data"`
			} `json:"data"`
		}
		var stmt StatementResp
		if err := json.Unmarshal(statementResult.raw, &stmt); err != nil {
			statementText = string(statementResult.raw)
		} else {
			decoded, err := DecodeBase64Text(stmt.Data.Data)
			if err != nil {
				statementText = stmt.Data.Data
			} else {
				statementText = decoded
			}
		}
	}

	infoText = fmt.Sprintf("Problem: #%d - %s\nTime Limit: %.2fs | Memory Limit: %dKB | Max Score: %d",
		infoResult.info.Data.Id, infoResult.info.Data.Name, infoResult.info.Data.Time, infoResult.info.Data.MemoryLimit, infoResult.info.Data.MaxScore)

	// Save to cache
	if err == nil && home != "" {
		cachePath := filepath.Join(home, ".kncli", "cache", fmt.Sprintf("%s.json", problemID))
		entry := CacheEntry{InfoText: infoText, StatementText: statementText}
		if cachedData, err := json.Marshal(entry); err == nil && os.MkdirAll(filepath.Dir(cachePath), 0755) == nil {
			_ = os.WriteFile(cachePath, cachedData, 0644)
		}
	}

	return infoText, statementText, nil
}
