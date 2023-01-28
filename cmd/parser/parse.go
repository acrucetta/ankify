package parser

import (
	"fmt"

	"github.com/acrucetta/pkg/parser"
	"github.com/spf13/cobra"
)

var ParseCmd = &cobra.Command{
	Use:     "parse",
	Aliases: []string{"p"},
	Short:   "Parses a PDF and generates Anki cards",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		res := parser.ParsePdf(args[0])
		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(ParseCmd)
}
