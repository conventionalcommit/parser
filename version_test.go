package main_test

import (
	main "cov-commit-parser"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNextVersion(t *testing.T) {
	const (
		v1   = "1.0.0"
		v11  = "1.1.0"
		v101 = "1.0.1"
		v2   = "2.0.0"

		minor = "feat: commit description"
		major = "feat!: commit description"
		bad   = "invalid commit message"
	)

	var (
		patch = fmt.Sprintf("%s: commit description", main.DefaultPatchTypes[0])
	)

	cases := map[string]struct {
		current    string
		messages   []string
		patchTypes []string
		expected   string
	}{
		"single patch change": {
			current: v1,
			messages: []string{
				patch,
			},
			patchTypes: main.DefaultPatchTypes,
			expected:   v101,
		},
		"single minor change": {
			current: v1,
			messages: []string{
				minor,
			},
			patchTypes: main.DefaultPatchTypes,
			expected:   v11,
		},
		"single major change": {
			current: v1,
			messages: []string{
				major,
			},
			patchTypes: main.DefaultPatchTypes,
			expected:   v2,
		},
		"no changes creates minor increment": {
			current:    v1,
			messages:   []string{},
			patchTypes: main.DefaultPatchTypes,
			expected:   v11,
		},
		"major change beats minor change": {
			current: v1,
			messages: []string{
				major,
				minor,
			},
			patchTypes: main.DefaultPatchTypes,
			expected:   v2,
		},
		"major change beats patch change": {
			current: v1,
			messages: []string{
				major,
				patch,
			},
			patchTypes: main.DefaultPatchTypes,
			expected:   v2,
		},
		"minor change beats patch change": {
			current: v1,
			messages: []string{
				minor,
				patch,
			},
			patchTypes: main.DefaultPatchTypes,
			expected:   v11,
		},
		"bad message ignored": {
			current: v1,
			messages: []string{
				bad,
				major,
			},
			patchTypes: main.DefaultPatchTypes,
			expected:   v2,
		},
		"only bad messages creates minor increment": {
			current: v1,
			messages: []string{
				bad,
			},
			patchTypes: main.DefaultPatchTypes,
			expected:   v11,
		},
	}

	for name, data := range cases {
		v, err := main.GetNextVersion(data.current, data.messages, data.patchTypes)
		if assert.NoErrorf(t, err, name) {
			assert.Equalf(t, data.expected, v, name)
		}
	}
}

func TestGetNextVersionError(t *testing.T) {
	const (
		v1 = "1.0.0"
	)

	var (
		patch = fmt.Sprintf("%s: commit description", main.DefaultPatchTypes[0])
	)

	cases := map[string]struct {
		current    string
		messages   []string
		patchTypes []string
	}{
		"invalid current version": {
			current: "bad version",
			messages: []string{
				patch,
			},
			patchTypes: main.DefaultPatchTypes,
		},
	}

	for name, data := range cases {
		_, err := main.GetNextVersion(data.current, data.messages, data.patchTypes)
		assert.Errorf(t, err, name)
	}
}

func TestGetNextVersionFromCommits(t *testing.T) {
	const (
		v1   = "1.0.0"
		v11  = "1.1.0"
		v101 = "1.0.1"
		v2   = "2.0.0"

		d = "commit description"
	)

	var (
		patch = main.ConventionalCommit{
			CommitType:  main.DefaultPatchTypes[0],
			Description: d,
		}
		minor = main.ConventionalCommit{
			CommitType:  "feat",
			Description: d,
		}
		major = main.ConventionalCommit{
			BreakingChange: true,
			CommitType:     "feat",
			Description:    d,
		}
	)

	cases := map[string]struct {
		current    string
		commits    []main.ConventionalCommit
		patchTypes []string
		expected   string
	}{
		"single patch change": {
			current: v1,
			commits: []main.ConventionalCommit{
				patch,
			},
			patchTypes: main.DefaultPatchTypes,
			expected:   v101,
		},
		"single minor change": {
			current: v1,
			commits: []main.ConventionalCommit{
				minor,
			},
			patchTypes: main.DefaultPatchTypes,
			expected:   v11,
		},
		"single major change": {
			current: v1,
			commits: []main.ConventionalCommit{
				major,
			},
			patchTypes: main.DefaultPatchTypes,
			expected:   v2,
		},
		"no changes creates minor increment": {
			current:    v1,
			commits:    []main.ConventionalCommit{},
			patchTypes: main.DefaultPatchTypes,
			expected:   v11,
		},
		"major change beats minor change": {
			current: v1,
			commits: []main.ConventionalCommit{
				major,
				minor,
			},
			patchTypes: main.DefaultPatchTypes,
			expected:   v2,
		},
		"major change beats patch change": {
			current: v1,
			commits: []main.ConventionalCommit{
				major,
				patch,
			},
			patchTypes: main.DefaultPatchTypes,
			expected:   v2,
		},
		"minor change beats patch change": {
			current: v1,
			commits: []main.ConventionalCommit{
				minor,
				patch,
			},
			patchTypes: main.DefaultPatchTypes,
			expected:   v11,
		},
	}

	for name, data := range cases {
		v, err := main.GetNextVersionFromCommits(data.current, data.commits, data.patchTypes)
		if assert.NoErrorf(t, err, name) {
			assert.Equalf(t, data.expected, v, name)
		}
	}
}

func TestGetNextVersionFromCommitsError(t *testing.T) {
	const (
		v1 = "1.0.0"

		d = "commit description"
	)

	var (
		patch = main.ConventionalCommit{
			CommitType:  main.DefaultPatchTypes[0],
			Description: d,
		}
	)

	cases := map[string]struct {
		current    string
		commits    []main.ConventionalCommit
		patchTypes []string
	}{
		"invalid current version": {
			current: "bad version",
			commits: []main.ConventionalCommit{
				patch,
			},
			patchTypes: main.DefaultPatchTypes,
		},
	}

	for name, data := range cases {
		_, err := main.GetNextVersionFromCommits(data.current, data.commits, data.patchTypes)
		assert.Errorf(t, err, name)
	}
}
