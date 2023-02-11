package parser

import (
	"fmt"
	"os"
	"time"

	"github.com/acrucetta/anki-builder/pkg/ankify"
	"github.com/acrucetta/anki-builder/pkg/docparser"
	"github.com/spf13/cobra"
)

var AnkifyCmd = &cobra.Command{
	Use:     "ankify [file path]",
	Aliases: []string{"a"},
	Short:   "Parses a PDF and generates Anki cards",
	Long:    `Parses a PDF and generates Anki cards, which are then printed to the console and saved as a txt file in your output folder. You may use the flag "type" or "t" to specify the input file type.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// Get file type
		fileType, _ := cmd.Flags().GetString("type")

		var res string
		var err error

		if fileType == "txt" {
			// Parse text file
			res, _ = docparser.ParseTxt(args[0])
		} else if fileType == "pdf" {
			res, _ = docparser.ParsePdf(args[0])
		} else {
			fmt.Println("Flag 'type' must be either 'txt' or 'pdf, defaulting to 'txt'")
			res, _ = docparser.ParsePdf(args[0])
		}

		anki_cards, err := ankify.Ankify(res)

		if err != nil {
			fmt.Println(err)
		}

		// Save string as txt using os package
		// Create file name based on date and time
		const folder string = "output"
		var file_name string = time.Now().Format("2006-01-02_15-04-05") + ".txt"

		// Create output folder if it doesn't exist
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			os.Mkdir(folder, 0755)
		}

		// Create file path
		file_name = folder + "/" + file_name

		// Create txt file
		file, err := os.OpenFile(file_name, os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			fmt.Println(err)
		}

		// Write to file
		_, err = file.WriteString(anki_cards)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Anki cards saved to " + file_name)

	},
}

func init() {
	rootCmd.AddCommand(AnkifyCmd)
	AnkifyCmd.Flags().StringP("type", "t", "", "Type of file to parse, either 'txt' or 'pdf'")
}
