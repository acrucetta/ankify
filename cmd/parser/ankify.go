package parser

import (
	"fmt"

	"github.com/acrucetta/anki-builder/pkg/ankify"
	"github.com/acrucetta/anki-builder/pkg/pdfparser"
	"github.com/spf13/cobra"
)

var AnkifyCmd = &cobra.Command{
	Use:     "ankify",
	Aliases: []string{"a"},
	Short:   "Parses a PDF and generates Anki cards",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// Parse PDF
		res, err := pdfparser.ParsePdf(args[0])

		if err != nil {
			fmt.Println(err)
		}

		res, err = ankify.Ankify(res)

		if err != nil {
			fmt.Println(err)
		}

		// Generate Anki cards
		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(AnkifyCmd)
}
