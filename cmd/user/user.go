// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package user

import (
	"fmt"
	"os"
	"text/template"

	utility "kilocli/cmd/utility"

	"github.com/charmbracelet/bubbles/table"
)

type UserDetailResponse struct {
	Data struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Admin       bool   `json:"admin"`
		Proposer    bool   `json:"proposer"`
		DisplayName string `json:"display_name"`
	} `json:"data"`
}

type UserSolvedProblems struct {
	Data []struct {
		ProblemId int     `json:"id"`
		Name      string  `json:"name"`
		Source    string  `json:"source_credits"`
		Score     float64 `json:"score_scale"`
	} `json:"data"`
}

func init() {
	SettingsCmd.AddCommand(ExtendSessionCmd)
	SettingsCmd.AddCommand(SetBioCmd)
	SettingsCmd.AddCommand(ChangeNameCmd)
	SettingsCmd.AddCommand(ChangePassCmd)
	SettingsCmd.AddCommand(ChangeEmailCmd)
	SettingsCmd.AddCommand(ResetPassCmd)
	SettingsCmd.AddCommand(DeleteUserCmd)
	SettingsCmd.AddCommand(ResendEmailCmd)
	SettingsCmd.AddCommand(AmILoggedInCmd)
	SettingsCmd.AddCommand(AmIAdminCmd)

}

func printUserDetails(dataUser UserDetailResponse) {
	userTemplate := `ID: {{.Id}}
Name: {{.Name}}
A.K.A: {{.DisplayName}}
Bio: {{.Bio}}
Admin: {{.Admin}}
Proposer: {{.Proposer}}`

	bio := getUserBio(dataUser.Data.Name)

	userData := struct {
		Id          int
		Name        string
		DisplayName string
		Bio         string
		Admin       bool
		Proposer    bool
	}{
		Id:          dataUser.Data.Id,
		Name:        dataUser.Data.Name,
		DisplayName: dataUser.Data.DisplayName,
		Bio:         bio,
		Admin:       dataUser.Data.Admin,
		Proposer:    dataUser.Data.Proposer,
	}

	tmpl, err := template.New("userDetails").Parse(userTemplate)
	if err != nil {
		utility.LogError(fmt.Errorf("error parsing template: %w", err))
		return
	}

	if err := tmpl.Execute(os.Stdout, userData); err != nil {
		utility.LogError(fmt.Errorf("error executing template: %w", err))
		return
	}
}

func prepareTableRows(dataUser UserSolvedProblems) []table.Row {
	var rows []table.Row
	for _, problem := range dataUser.Data {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", problem.ProblemId),
			problem.Name,
			problem.Source,
			fmt.Sprintf("%.0f", problem.Score),
		})
	}
	return rows
}

func isCurrentUserLoggedIn() bool {
	_, hasToken := utility.ReadToken()
	return hasToken
}

func isAdmin(userId string) bool {
	return userGetDetails(userId, "isadmin")
}

/*
func deleteUser() {
	fmt.Print("Are you sure to delete your account? (Y/N) ")
	var resp string
	fmt.Scan(&resp)

	if resp != "Y" {
		return
	}

	fmt.Print("Do you understand that this action is irreversible? (Y/N) ")
	fmt.Scan(&resp)

	if resp != "Y" {
		return
	}

	fmt.Print("Do you understand that you'll lose all of your data? (Y/N) ")
	fmt.Scan(&resp)

	if resp != "Y" {
		return
	}

	fmt.Println("Deleting account...")
	fmt.Println("Currently the KN api does not support deleting an account if you're not an admin.")

	if !isAdmin("me") {
		return
	}

	ResponseBody, err := MakePostRequest(URL_CHANGE_EMAIL, nil, RequestFormAuth)
	if err != nil {
		LogError(err)
		return
	}

	var dataKN KilonovaResponse
	err = json.Unmarshal(ResponseBody, &dataKN)
	if err != nil {
		LogError(err)
		return
	}

	fmt.Println(dataKN.Data)
}
*/
