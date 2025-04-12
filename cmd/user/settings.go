// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kncli/internal"
	u "net/url"
)

func setUserBio(bio string) {
	payload := map[string]string{"bio": bio}
	resp, err := internal.PostJSON[internal.KilonovaResponse](internal.URL_SELF_SET_BIO, payload)
	if err != nil {
		internal.LogError(err)
		return
	}

	if resp.Status == internal.SUCCESS {
		fmt.Println("Success! Bio changed!")
		return
	}
	fmt.Println("Error: Failed to change bio!")

}

func changeName(newName, password string) {
	payload := map[string]string{
		"newName":  newName,
		"password": password,
	}
	resp, err := internal.PostJSON[internal.KilonovaResponse](internal.URL_CHANGE_NAME, payload)
	if err != nil {
		internal.LogError(err)
		return
	}

	if resp.Status == internal.SUCCESS {
		fmt.Println("Success! Name changed!")
		return
	}
	internal.LogError(fmt.Errorf("failed to change name"))
}

func changePass(oldPass, newPass string) {
	payload := map[string]string{
		"old_password": oldPass,
		"password":     newPass,
	}
	resp, err := internal.PostJSON[internal.KilonovaResponse](internal.URL_CHANGE_PASS, payload)
	if err != nil {
		internal.LogError(err)
		return
	}

	if resp.Status == internal.SUCCESS {
		fmt.Println("Success! Password changed! You'll need to login again.")
		logout()
		return
	}
	internal.LogError(fmt.Errorf("failed to change password"))
}

func changeEmail(email, password string) {
	formData := u.Values{}
	formData.Set("email", email)
	formData.Set("password", password)

	ResponseBody, err := internal.MakePostRequest(internal.URL_CHANGE_EMAIL, bytes.NewBufferString(formData.Encode()), internal.RequestFormAuth)
	if err != nil {
		internal.LogError(err)
		return
	}

	var res internal.KilonovaResponse
	if err := json.Unmarshal(ResponseBody, &res); err != nil {
		internal.LogError(err)
		return
	}

	if res.Status == internal.SUCCESS {
		fmt.Println("Success! Email changed!")
		return
	}
	internal.LogError(fmt.Errorf("failed to change email"))
}

func resetPass(email string) {
	if _, loggedIn := internal.ReadToken(); loggedIn {
		fmt.Println("You must be logged out to reset your password.")
		return
	}

	form := u.Values{}
	form.Set("email", email)

	ResponseBody, err := internal.MakePostRequest(internal.URL_CHANGE_PASS, bytes.NewBufferString(form.Encode()), internal.RequestFormAuth)
	if err != nil {
		internal.LogError(err)
		return
	}

	var res internal.KilonovaResponse
	if err := json.Unmarshal(ResponseBody, &res); err != nil {
		internal.LogError(err)
		return
	}

	fmt.Println(res.Data)
}

func resendEmail() {
	ResponseBody, err := internal.MakePostRequest(internal.URL_RESEND_MAIL, nil, internal.RequestFormAuth)
	if err != nil {
		internal.LogError(err)
		return
	}

	var res internal.KilonovaResponse
	if err := json.Unmarshal(ResponseBody, &res); err != nil {
		internal.LogError(err)
		return
	}

	fmt.Println(res.Data)
}
