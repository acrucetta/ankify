package parser

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "Anki Builder",
	Short: "Anki Builder - a simple CLI to parse PDFs and generate Anki cards",
	Long:  `Anki Builder - a simple CLI to parse PDFs and generate Anki cards`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Anki Builder - a simple CLI to parse PDFs and generate Anki cards")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
