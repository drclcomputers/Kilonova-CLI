// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package user

import (
	"fmt"
	utility "kncli/internal"

	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
)

var SettingsCmd = &cobra.Command{
	Use:   "settings [command] ...",
	Short: "Modify your account.",
}

var SigninCmd = &cobra.Command{
	Use:   "signin [username] [password]",
	Short: "Sign in to your account",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		username, password := loginForm()
		action := func() { login(username, password) }
		if err := spinner.New().Title("Logging in...").Action(action).Run(); err != nil {
			utility.LogError(err)
			return
		}
	},
}

var LogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of your account",
	Run: func(cmd *cobra.Command, args []string) {
		action := func() { logout() }
		if err := spinner.New().Title("Waiting ...").Action(action).Run(); err != nil {
			utility.LogError(err)
			return
		}
	},
}

var UserGetDetailsCmd = &cobra.Command{
	Use:   "user [User ID or me (get self ID)]",
	Short: "Get details about a user.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userGetDetails(args[0], "user")
	},
}

var UserSolvedProblemsCmd = &cobra.Command{
	Use:   "solvedproblems [User ID or me (get self ID)]",
	Short: "Get list of solved problems by user.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userGetSolvedProblems(args[0])
	},
}

var ExtendSessionCmd = &cobra.Command{
	Use:   "extendsession",
	Short: "Extend the current session for 30 days more.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		extendSession()
	},
}

var SetBioCmd = &cobra.Command{
	Use:   "setbio [bio]",
	Short: "Set your profile's bio.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		setUserBio(args[0])
	},
}

var ChangeNameCmd = &cobra.Command{
	Use:   "changename [new name] [password]",
	Short: "Change your profile name.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		changeName(args[0], args[1])
	},
}

var ChangePassCmd = &cobra.Command{
	Use:   "changepass [old password] [new password]",
	Short: "Change your account password.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		changePass(args[0], args[1])
	},
}

var ResetPassCmd = &cobra.Command{
	Use:   "resetpass [email]",
	Short: "Reset password via email when forgotten.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resetPass(args[0])
	},
}

var DeleteUserCmd = &cobra.Command{
	Use:   "deleteuser",
	Short: "Delete your Kilonova account. (Currently not working in API V1)",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		//deleteUser()
	},
}

var ChangeEmailCmd = &cobra.Command{
	Use:   "changemail [new email] [password]",
	Short: "Change your account email.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		changeEmail(args[0], args[1])
	},
}

var ResendEmailCmd = &cobra.Command{
	Use:   "resendemail",
	Short: "Resend verification mail.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		resendEmail()
	},
}

var AmILoggedInCmd = &cobra.Command{
	Use:   "amilogged",
	Short: "Check wether you're logged in or not.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(isCurrentUserLoggedIn())
	},
}

var AmIAdminCmd = &cobra.Command{
	Use:   "amiadmin",
	Short: "Check wether you're an admin.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(isAdmin("me"))
	},
}
