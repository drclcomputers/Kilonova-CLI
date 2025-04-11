// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package problems

import (
	"fmt"
	"kncli/internal"
	"os"

	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
)

var GetAssetsCmd = &cobra.Command{
	Use:   "assets [Problem ID]",
	Short: "Download the assets for a problem. (online)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		action := func() { _ = GetAssets(args[0]) }
		if err := spinner.New().Title("Please wait...").Action(action).Run(); err != nil {
			internal.LogError(err)
			return
		}
	},
}

func init() {
}

func saveToFile(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}

	return nil
}

func GetAssets(id string) error {
	url := fmt.Sprintf(internal.URL_ASSETS, id)

	OutputFile := fmt.Sprintf("%s.zip", id)
	DataToBeWritten, err := internal.MakeGetRequest(url, nil, internal.RequestDownloadZip)
	if err != nil {
		internal.LogError(fmt.Errorf("error making request: %v", err))
		return err
	}

	if err := saveToFile(OutputFile, DataToBeWritten); err != nil {
		internal.LogError(fmt.Errorf("error saving file: %v", err))
		return err
	}

	fmt.Println("ZIP file downloaded successfully:", OutputFile)

	return nil

}
