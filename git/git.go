package git

import (
	"errors"

	"github.com/Masterminds/semver"
	"github.com/xfxdev/xlog"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// GetCommitsInDirectory returns all the commit messages between the two refs
// (exluding `from`, but including `to`). The order of the returned commits is
// not guaranteed
func GetCommitsInDirectory(from, to, directory string) ([]string, error) {
	// Load the repository
	r, err := git.PlainOpen(directory)
	if err != nil {
		return nil, err
	}

	return getCommits(r, from, to)
}

// getCommits returns all the commit messages for the given repository between the two
// refs (exluding `from`, but including `to`)
func getCommits(r *git.Repository, from, to string) ([]string, error) {
	// Resolve the revisions to hashes
	xlog.Debugf("Resolving revision %s", to)
	toHash, err := r.ResolveRevision(plumbing.Revision(to))
	if err != nil {
		return nil, err
	}

	xlog.Debugf("Resolving revision %s", from)
	fromHash, err := r.ResolveRevision(plumbing.Revision(from))
	if err != nil {
		return nil, err
	}

	// Get an iterator for the commit messages
	iter, err := r.Log(&git.LogOptions{
		From: *toHash,
	})
	if err != nil {
		return nil, err
	}

	// Record all the commit messages
	messages := []string{}
	for {
		c, err := iter.Next()
		if err != nil {
			xlog.Debug("No more messages to parse")
			break
		}

		if c.Hash.String() == fromHash.String() {
			break
		}

		messages = append(messages, c.Message)
	}
	xlog.Debugf("Got messages: %v", messages)
	return messages, nil
}

// GetLatestVersionInDirectory retrieves the latest version tag for the repository in the given directory
func GetLatestVersionInDirectory(directory string) (string, error) {
	// Load the repository
	r, err := git.PlainOpen(directory)
	if err != nil {
		return "", err
	}

	return getLatestVersion(r)
}

// getLatestVersion retrieves the latest version tag for the given repository
func getLatestVersion(r *git.Repository) (tag string, err error) {
	iter, err := r.Tags()
	if err != nil {
		return "", err
	}

	latestTag := semver.MustParse("0.0.0")
	err = iter.ForEach(func(ref *plumbing.Reference) error {
		v, err := semver.NewVersion(ref.Name().Short())
		if err != nil {
			xlog.Debugf("Ignoring tag %s as it is not a valid semver", ref.Name().Short())
			return nil
		}

		if v.GreaterThan(latestTag) {
			latestTag = v
			return nil
		}

		// version was older than current version
		return nil
	})
	if err != nil {
		return "", err
	}

	if latestTag.String() == "0.0.0" {
		return "", errors.New("Unable to find latest version")
	}

	return latestTag.String(), nil
}
