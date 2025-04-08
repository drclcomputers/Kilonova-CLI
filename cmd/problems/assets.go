// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package problems

import (
	"fmt"
	"os"

	utility "kilocli/cmd/utility"

	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
)

var GetAssetsCmd = &cobra.Command{
	Use:   "assets [Problem ID]",
	Short: "Download the assets for a problem.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		action := func() { GetAssets(args[0]) }
		if err := spinner.New().Title("Waiting ...").Action(action).Run(); err != nil {
			utility.LogError(err)
			return
		}
	},
}

func init() {
}

func saveToFile(filename string, data []byte) error {
	File, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer File.Close()

	_, err = File.Write(data)
	return err
}

func GetAssets(id string) error {
	url := fmt.Sprintf(utility.URL_ASSETS, id)

	OutputFile := fmt.Sprintf("%s.zip", id)
	DataToBeWritten, err := utility.MakeGetRequest(url, nil, utility.RequestDownloadZip)
	if err != nil {
		utility.LogError(fmt.Errorf("error making request: %v", err))
		return err
	}

	err = saveToFile(OutputFile, DataToBeWritten)
	if err != nil {
		utility.LogError(fmt.Errorf("error saving file: %v", err))
		return err
	}

	fmt.Println("ZIP file downloaded successfully:", OutputFile)

	return nil

}
