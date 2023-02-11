package parser

import (
	"fmt"

	"github.com/acrucetta/anki-builder/pkg/docparser"
	"github.com/spf13/cobra"
)

var ParseCmd = &cobra.Command{
	Use:     "parse",
	Aliases: []string{"p"},
	Short:   "Parses a PDF and generates Anki cards",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		res, err := docparser.ParsePdf(args[0])
		if err != nil {
			fmt.Println(err)
		}

		// Generate Anki cards
		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(ParseCmd)
}
