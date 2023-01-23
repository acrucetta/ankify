package parser

import (
	"fmt"

	"github.com/spf13/cobra"
)

var parseCmd = &cobra.Command{
	Use:     "parse",
	Aliases: []string{"p"},
	Short:   "Parses a PDF and generates Anki cards",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		res := parser.parse_pdf(args[0])
		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)
}
