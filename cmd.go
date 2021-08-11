package main

import (
	"github.com/spf13/cobra"
	"github.com/xfxdev/xlog"
)

// flags
var (
	current   string
	directory string
	since     string
	verbose   bool
)

var rootCmd = &cobra.Command{
	Use:     "ccp",
	Short:   "Cov Commit Parser is a simple tool for parsing conventional commits",
	Long:    "A simple tool for parsing conventional commits. The full specification for conventional commits is available at https://conventionalcommits.org/en/v1.0.0",
	Version: func() string { return applicationVersion }(), // Use a func() here so we can override the variable using linker
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		xlog.SetLevel(xlog.InfoLevel)
		if verbose {
			xlog.SetLevel(xlog.DebugLevel)
		}
		return nil
	},
}

func setup() {
	rootCmd.PersistentFlags().StringVarP(&current, "current", "c", "", "Current version number from which to base the version change. Defaults to the latest version tag in the repository")
	rootCmd.PersistentFlags().StringVarP(&directory, "directory", "d", "", "Directory of the git repository")
	rootCmd.PersistentFlags().StringVarP(&since, "since", "s", "", "Revision to track commits from. Defaults to the latest version tag in the repository")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose logging")
}
