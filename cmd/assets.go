package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
)

var getAssetsCmd = &cobra.Command{
	Use:   "assets [Problem ID]",
	Short: "Download the assets for a problem.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		action := func() { getAssets(args[0]) }
		if err := spinner.New().Title("Waiting ...").Action(action).Run(); err != nil {
			logError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(getAssetsCmd)
}

func saveToFile(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func getAssets(id string) {
	url := fmt.Sprintf(URL_ASSETS, id)

	outputFile := fmt.Sprintf("%s.zip", id)
	data, err := MakeGetRequest(url, nil, RequestDownloadZip)
	if err != nil {
		logError(fmt.Errorf("error making request: %v", err))
		return
	}

	err = saveToFile(outputFile, data)
	if err != nil {
		logError(fmt.Errorf("error saving file: %v", err))
		return
	}

	fmt.Println("ZIP file downloaded successfully:", outputFile)
}
