package cli

import (
	"github.com/spf13/cobra"
)

// These will be set by the linker
var (
	Version   string
	GitCommit string
	BuildTime string
)

func init() {
	// Set defaults for development builds
	if Version == "" {
		Version = "dev"
	}
	if GitCommit == "" {
		GitCommit = "unknown"
	}
	if BuildTime == "" {
		BuildTime = "unknown"
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Printf("prepf version %s\n", Version)
		cmd.Printf("Commit: %s\n", GitCommit)
		cmd.Printf("Built: %s\n", BuildTime)
		return nil
	},
}
