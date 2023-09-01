package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-version-bumper",
	Short: "A CLI tool to upgrade version of your package.json",
	Long: `A CLI tool to upgrade version of your package.json depending on the name of your current git branch. For example:

	$ version-bumper bump`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
