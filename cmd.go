package main

import (
	"ccp/git"
	"fmt"
	"os"

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

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Determine the next build number to use",
	Long:  "Determine the next build number to use based on git tags or specified verisons",
	RunE:  runVersionCmd,
}

func runVersionCmd(cmd *cobra.Command, args []string) (err error) {
	// Get the directory of the git repo. Default to the current working directory
	if directory == "" {
		directory, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	xlog.Debugf("Using directory %s", directory)

	// Get the latest version tag in the git repo
	latestVersion, err := git.GetLatestVersionInDirectory(directory)
	if err != nil {
		if current == "" || since == "" { // We only need this if either `current` or `since` were not provided
			xlog.Warn("Unable to determine current version so suggesting initial version of 0.1.0")
			fmt.Println("0.1.0")
			return nil
		}
	}

	// Get the current version number. Default to the latest tag on the branch
	if current == "" {
		current = latestVersion
	}
	xlog.Debugf("Using current version %s", current)

	// Get the start point
	if since == "" {
		since = latestVersion
	}
	xlog.Debugf("Discovering commits since %s", since)

	// Discover the commits
	commits, err := git.GetCommitsInDirectory(since, "HEAD", directory)
	if err != nil {
		return err
	}

	// Compute the next version
	v, err := GetNextVersion(current, commits, DefaultPatchTypes)
	if err != nil {
		return err
	}

	// Print the version
	fmt.Println(v)

	return nil
}

func setup() {
	rootCmd.PersistentFlags().StringVarP(&current, "current", "c", "", "Current version number from which to base the version change. Defaults to the latest version tag in the repository")
	rootCmd.PersistentFlags().StringVarP(&directory, "directory", "d", "", "Directory of the git repository")
	rootCmd.PersistentFlags().StringVarP(&since, "since", "s", "", "Revision to track commits from. Defaults to the latest version tag in the repository")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose logging")

	rootCmd.AddCommand(versionCmd)
}
