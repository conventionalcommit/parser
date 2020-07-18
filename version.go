package main

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/xfxdev/xlog"
)

// GetNextVersion returns the next version to use based on the current version
// and a slice of commit messages
func GetNextVersion(currentVersion string, messages, patchTypes []string) (string, error) {
	commits, errs := ParseMessages(messages)
	for _, err := range errs {
		xlog.Warn(err) // Log any errors, but dont actually terminate
	}

	return GetNextVersionFromCommits(currentVersion, commits, patchTypes)
}

// GetNextVersionFromCommits returns the next version to use based on the current
// version and a slice of `ConventionalCommit`s
func GetNextVersionFromCommits(currentVersion string, commits []ConventionalCommit, patchTypes []string) (string, error) {
	current, err := semver.NewVersion(currentVersion)
	if err != nil {
		return "", fmt.Errorf("Unable to parse version %s", currentVersion)
	}

	var newVersion semver.Version
	for _, commit := range commits {
		if commit.BreakingChange {
			xlog.Debugf("Commit message with description \"%s\" has resulted in a major version increment", commit.Description)
			newVersion = current.IncMajor()
			break
		}

		if !typeIsPatch(commit.CommitType, patchTypes) {
			xlog.Debugf("Commit message with type \"%s\" and description \"%s\" has resulted in a minor version increment", commit.CommitType, commit.Description)
			newVersion = current.IncMinor()
			break
		}

		newVersion = current.IncPatch()
	}

	if newVersion.Equal(&semver.Version{}) {
		xlog.Warn("No conventional commits parsed so defaulting to minor increment")
		newVersion = current.IncMinor()
	}

	return newVersion.String(), nil
}
