package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "prepf",
	Short: "Technical Interview Coach CLI",
	Long:  "prepf - A CLI tool for technical interview preparation with AI-powered coaching",
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
}
