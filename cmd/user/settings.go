// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package user

import (
	"fmt"
	"kncli/internal"
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
	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	resp, err := internal.PostJSON[internal.KilonovaResponse](internal.URL_CHANGE_EMAIL, payload)
	if err != nil {
		internal.LogError(err)
		return
	}

	if resp.Status == internal.SUCCESS {
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

	payload := map[string]string{"email": email}
	resp, err := internal.PostJSON[internal.KilonovaResponse](internal.URL_RESEND_MAIL, payload)
	if err != nil {
		internal.LogError(err)
		return
	}

	if resp.Status == internal.SUCCESS {
		fmt.Println("Password reset email sent! Check your inbox.")
	} else {
		fmt.Println(resp.Data)
	}
}

func resendEmail() {
	resp, err := internal.PostJSON[internal.KilonovaResponse](internal.URL_RESEND_MAIL, nil)
	if err != nil {
		internal.LogError(err)
		return
	}

	if resp.Status == internal.SUCCESS {
		fmt.Println("Verification email resent! Check your inbox.")
	} else {
		fmt.Println(resp.Data)
	}
}
