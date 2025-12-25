package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/trankhanh040147/prepf/internal/config"
)

type configKey struct{}
type flagsKey struct{}

type Flags struct {
	ConfigPath  string
	ProfilePath string
	Verbose     bool
	Quiet       bool
}

var (
	flagConfigPath  string
	flagProfilePath string
	flagVerbose     bool
	flagQuiet       bool
)

var rootCmd = &cobra.Command{
	Use:   "prepf",
	Short: "Technical Interview Coach CLI",
	Long:  "prepf - A CLI tool for technical interview preparation with AI-powered coaching",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Parse flags into struct
		flags := &Flags{
			ConfigPath:  flagConfigPath,
			ProfilePath: flagProfilePath,
			Verbose:     flagVerbose,
			Quiet:       flagQuiet,
		}
		cmd.SetContext(context.WithValue(cmd.Context(), flagsKey{}, flags))

		// Load config with flag overrides
		cfg, err := config.LoadWithOverrides(flags.ConfigPath, flags.ProfilePath)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		cmd.SetContext(context.WithValue(cmd.Context(), configKey{}, cfg))

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

// GetConfig returns the config from context, or nil if not set
func GetConfig(cmd *cobra.Command) *config.Config {
	if cfg, ok := cmd.Context().Value(configKey{}).(*config.Config); ok {
		return cfg
	}
	return nil
}

// GetFlags returns the flags from context, or nil if not set
func GetFlags(cmd *cobra.Command) *Flags {
	if flags, ok := cmd.Context().Value(flagsKey{}).(*Flags); ok {
		return flags
	}
	return nil
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&flagConfigPath, FlagConfig, FlagConfigShort, "", "override config file path")
	rootCmd.PersistentFlags().StringVarP(&flagProfilePath, FlagProfile, FlagProfileShort, "", "override profile file path")
	rootCmd.PersistentFlags().BoolVarP(&flagVerbose, FlagVerbose, FlagVerboseShort, false, "enable verbose logging")
	rootCmd.PersistentFlags().BoolVarP(&flagQuiet, FlagQuiet, FlagQuietShort, false, "suppress non-error output")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
}
