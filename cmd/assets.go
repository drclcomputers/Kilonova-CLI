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
			logErr(err)
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
	data, err := makeRequest("GET", url, nil, "4")
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}

	err = saveToFile(outputFile, data)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}

	fmt.Println("ZIP file downloaded successfully:", outputFile)
}
