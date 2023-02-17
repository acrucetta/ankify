package parser

import (
	"encoding/csv"
	"fmt"
	"log"
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
	Long: `Parses a PDF and generates Anki cards, which are then printed to the console and saved as a JSON file in your output folder. 
	You may use the flag "type" or "t" to specify the input file type.
	You may use the flag "page" or "p" to specify the page numbers to parse.
	You may use the flag "tag" or "T" to specify the tags to add to the cards.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		file_type, _ := cmd.Flags().GetString("type")

		var page_numbers []int
		page_numbers, _ = cmd.Flags().GetIntSlice("page")

		var tag string = ""
		tag, _ = cmd.Flags().GetString("tag")

		if len(page_numbers) == 0 {
			page_numbers = []int{1}
		}

		var res map[int]string
		var err error

		if file_type == "txt" {
			res, _ = docparser.ParseTxt(args[0])
		} else if file_type == "pdf" {
			res, _ = docparser.ParsePdf(args[0], page_numbers)
		} else if file_type == "url" {
			res, _ = docparser.ParseUrl(args[0])
		} else {
			fmt.Println("Flag 'type' must be either 'txt' or 'pdf, defaulting to 'txt'")
			res, _ = docparser.ParsePdf(args[0], page_numbers)
		}

		anki_cards, _ := ankify.Ankify(res)

		// Save string as txt using os package
		// Create file name based on date and time
		const folder string = "output"
		var file_name string = time.Now().Format("2006-01-02_15-04-05") + ".csv"

		// Create output folder if it doesn't exist
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			os.Mkdir(folder, 0755)
		}

		// Create file path
		file_name = folder + "/" + file_name

		// Create txt file
		file, err := os.OpenFile(file_name, os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			log.Fatal(err)
		}

		// Create a new CSV writer
		writer := csv.NewWriter(file)

		// Write the data rows based on the AnkiQuestion struct
		for _, card := range anki_cards.Questions {
			writer.Write([]string{card.Question, card.Answer, tag})
		}

		// Flush the writer
		writer.Flush()
	},
}

func init() {
	rootCmd.AddCommand(AnkifyCmd)
	AnkifyCmd.Flags().StringP("type", "t", "", "Type of file to parse, either 'txt', 'pdf', or 'url' (default is 'txt')")
	AnkifyCmd.Flags().IntSliceP("pages", "p", []int{}, "Page numbers to parse, e.g., '1,2,3' (default is 1)")
	AnkifyCmd.Flags().StringP("tag", "T", "", "Tags to add to the cards, e.g., 'tag1' (default is no tags)")
}
