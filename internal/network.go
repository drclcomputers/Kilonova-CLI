// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func CreateRequest(req http.Request, reqType RequestType, contentType ...string) *http.Request {

	req.Header.Set("User-Agent", UserAgent)
	token, hasToken := ReadToken()

	switch reqType {
	case RequestFormAuth, RequestFormGuest:
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	case RequestJSON:
		req.Header.Set("Content-Type", "application/json")
	case RequestDownloadZip, RequestInfo:
		req.Header.Set("Content-Type", "application/zip")
		req.Header.Set("Accept", "application/zip")
		if reqType == RequestDownloadZip {
			fmt.Println("Trying to obtain archive...")
		}
		cookie := &http.Cookie{
			Name:  "kn-sessionid",
			Value: token,
		}
		req.AddCookie(cookie)
	case RequestMultipartForm:
		if len(contentType) > 0 {
			req.Header.Set("Content-Type", contentType[0])
		} else {
			LogError(fmt.Errorf("missing content type for multipart form request"))
		}
	default:
	}

	if hasToken {
		req.Header.Set("Authorization", token)
	} else if reqType == RequestFormAuth || reqType == RequestDownloadZip {
		LogError(fmt.Errorf("you must be authenticated to do this"))
	}

	return &req
}

func MakeRequest(method, url string, ResponseBody io.Reader, reqType RequestType, contentType ...string) ([]byte, error) {
	req, err := http.NewRequest(method, url, ResponseBody)
	if err != nil {
		LogError(fmt.Errorf("error creating request: %w", err))
		return nil, err
	}
	req = CreateRequest(*req, reqType, contentType...)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		LogError(fmt.Errorf("error making request: %w", err))
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		LogError(fmt.Errorf("error reading response ResponseBody: %w", err))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		if reqType == RequestDatabase {
			if strings.Contains(string(data), "not") {
				return []byte("notfound"), nil
			}
		}
		var res RawKilonovaResponse
		if err := json.Unmarshal(data, &res); err != nil {
			LogError(err)
		}
		LogError(fmt.Errorf("error: %s", string(res.Data)))
	}

	return data, nil
}

func MakeGetRequest(url string, ResponseBody io.Reader, reqType RequestType, contentType ...string) ([]byte, error) {
	return MakeRequest("GET", url, ResponseBody, reqType, contentType...)
}

func MakePostRequest(url string, ResponseBody io.Reader, reqType RequestType, contentType ...string) ([]byte, error) {
	return MakeRequest("POST", url, ResponseBody, reqType, contentType...)
}

func PostJSON[T any](url string, payload any) (T, error) {
	var result T

	jsonData, err := json.Marshal(payload)
	if err != nil {
		LogError(fmt.Errorf("failed to marshal JSON: %w", err))
		return result, err
	}

	ResponseBody, err := MakePostRequest(url, bytes.NewBuffer(jsonData), RequestJSON)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(ResponseBody, &result)
	if err != nil {
		LogError(fmt.Errorf("failed to decode response: %w", err))
		return result, err
	}

	return result, nil
}
