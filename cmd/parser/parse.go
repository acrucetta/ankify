package parser

import (
	"fmt"

	"github.com/acrucetta/anki-builder/pkg/pdfparser"
	"github.com/spf13/cobra"
)

var ParseCmd = &cobra.Command{
	Use:     "parse",
	Aliases: []string{"p"},
	Short:   "Parses a PDF and generates Anki cards",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		res, err := pdfparser.ParsePdf(args[0])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(ParseCmd)
}
