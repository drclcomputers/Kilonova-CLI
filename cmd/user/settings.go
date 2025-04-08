// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	utility "kilocli/cmd/utility"
	u "net/url"
)

func setUserBio(bio string) {
	payload := map[string]string{"bio": bio}
	resp, err := utility.PostJSON[utility.KilonovaResponse](utility.URL_SELF_SET_BIO, payload)
	if err != nil {
		utility.LogError(err)
		return
	}

	if resp.Status == utility.SUCCESS {
		fmt.Println("Success! Bio changed!")
	} else {
		fmt.Println("Error: Failed to change bio!")
	}
}

func changeName(newName, password string) {
	payload := map[string]string{
		"newName":  newName,
		"password": password,
	}
	resp, err := utility.PostJSON[utility.KilonovaResponse](utility.URL_CHANGE_NAME, payload)
	if err != nil {
		utility.LogError(err)
		return
	}

	if resp.Status == utility.SUCCESS {
		fmt.Println("Success! Name changed!")
	} else {
		utility.LogError(fmt.Errorf("failed to change name"))
		return
	}
}

func changePass(oldPass, newPass string) {
	payload := map[string]string{
		"old_password": oldPass,
		"password":     newPass,
	}
	resp, err := utility.PostJSON[utility.KilonovaResponse](utility.URL_CHANGE_PASS, payload)
	if err != nil {
		utility.LogError(err)
		return
	}

	if resp.Status == utility.SUCCESS {
		fmt.Println("Success! Password changed! You'll need to login again.")
		logout()
	} else {
		utility.LogError(fmt.Errorf("failed to change password"))
		return
	}
}

func changeEmail(email, password string) {
	formData := u.Values{}
	formData.Set("email", email)
	formData.Set("password", password)

	ResponseBody, err := utility.MakePostRequest(utility.URL_CHANGE_EMAIL, bytes.NewBufferString(formData.Encode()), utility.RequestFormAuth)
	if err != nil {
		utility.LogError(err)
		return
	}

	var res utility.KilonovaResponse
	if err := json.Unmarshal(ResponseBody, &res); err != nil {
		utility.LogError(err)
		return
	}

	if res.Status == utility.SUCCESS {
		fmt.Println("Success! Email changed!")
	} else {
		utility.LogError(fmt.Errorf("failed to change email"))
		return
	}
}

func resetPass(email string) {
	if _, loggedIn := utility.ReadToken(); loggedIn {
		fmt.Println("You must be logged out to reset your password.")
		return
	}

	form := u.Values{}
	form.Set("email", email)

	ResponseBody, err := utility.MakePostRequest(utility.URL_CHANGE_PASS, bytes.NewBufferString(form.Encode()), utility.RequestFormAuth)
	if err != nil {
		utility.LogError(err)
		return
	}

	var res utility.KilonovaResponse
	if err := json.Unmarshal(ResponseBody, &res); err != nil {
		utility.LogError(err)
		return
	}

	fmt.Println(res.Data)
}

func resendEmail() {
	ResponseBody, err := utility.MakePostRequest(utility.URL_RESEND_MAIL, nil, utility.RequestFormAuth)
	if err != nil {
		utility.LogError(err)
		return
	}

	var res utility.KilonovaResponse
	if err := json.Unmarshal(ResponseBody, &res); err != nil {
		utility.LogError(err)
		return
	}

	fmt.Println(res.Data)
}
