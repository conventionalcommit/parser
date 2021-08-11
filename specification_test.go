package parser_test

import (
	"testing"

	parser "github.com/conventionalcommit/parser"
	"github.com/stretchr/testify/assert"
)

func parseMessageHelper(t *testing.T, dir string, expected parser.ConventionalCommit) {
	t.Helper()

	commit := loadStringFromFile(t, dir)
	out, err := parser.ParseMessage(commit)
	if assert.NoError(t, err) {
		assert.Equal(t, expected, out)
	}
}

func TestParseMessage(t *testing.T) {
	const (
		commitBody = `This is a multiline commit body.

This is the second line`
		commitDescription = "description message"
		commitScope       = "scope"
		commitType        = "type"
		dir               = "test_fixtures/example_commits"
	)

	var (
		breakingChangeFooter = map[string]string{
			"BREAKING CHANGE": "reason",
		}
		commitFooters = map[string]string{
			"footer":      "simple",
			"hash-footer": "123",
		}
		emptyFooters = map[string]string{}
	)

	t.Run("description", func(t *testing.T) {
		parseMessageHelper(t, dir, parser.ConventionalCommit{
			CommitType:  commitType,
			Description: commitDescription,
			Footers:     emptyFooters,
		})
	})
	t.Run("description scope", func(t *testing.T) {
		parseMessageHelper(t, dir, parser.ConventionalCommit{
			CommitScope: commitScope,
			CommitType:  commitType,
			Description: commitDescription,
			Footers:     emptyFooters,
		})
	})
	t.Run("breaking change description", func(t *testing.T) {
		parseMessageHelper(t, dir, parser.ConventionalCommit{
			BreakingChange: true,
			CommitType:     commitType,
			Description:    commitDescription,
			Footers:        emptyFooters,
		})
	})
	t.Run("breaking change description scope", func(t *testing.T) {
		parseMessageHelper(t, dir, parser.ConventionalCommit{
			BreakingChange: true,
			CommitScope:    commitScope,
			CommitType:     commitType,
			Description:    commitDescription,
			Footers:        emptyFooters,
		})
	})
	t.Run("description body", func(t *testing.T) {
		parseMessageHelper(t, dir, parser.ConventionalCommit{
			Body:        commitBody,
			CommitType:  commitType,
			Description: commitDescription,
			Footers:     emptyFooters,
		})
	})
	t.Run("description scope body", func(t *testing.T) {
		parseMessageHelper(t, dir, parser.ConventionalCommit{
			Body:        commitBody,
			CommitScope: commitScope,
			CommitType:  commitType,
			Description: commitDescription,
			Footers:     emptyFooters,
		})
	})
	t.Run("breaking change description body", func(t *testing.T) {
		parseMessageHelper(t, dir, parser.ConventionalCommit{
			Body:           commitBody,
			BreakingChange: true,
			CommitType:     commitType,
			Description:    commitDescription,
			Footers:        emptyFooters,
		})
	})
	t.Run("breaking change description scope body", func(t *testing.T) {
		parseMessageHelper(t, dir, parser.ConventionalCommit{
			Body:           commitBody,
			BreakingChange: true,
			CommitScope:    commitScope,
			CommitType:     commitType,
			Description:    commitDescription,
			Footers:        emptyFooters,
		})
	})
	t.Run("description footers", func(t *testing.T) {
		parseMessageHelper(t, dir, parser.ConventionalCommit{
			CommitType:  commitType,
			Description: commitDescription,
			Footers:     commitFooters,
		})
	})
	t.Run("description scope footers", func(t *testing.T) {
		parseMessageHelper(t, dir, parser.ConventionalCommit{
			CommitScope: commitScope,
			CommitType:  commitType,
			Description: commitDescription,
			Footers:     commitFooters,
		})
	})
	t.Run("description body footers", func(t *testing.T) {
		parseMessageHelper(t, dir, parser.ConventionalCommit{
			Body:        commitBody,
			CommitType:  commitType,
			Description: commitDescription,
			Footers:     commitFooters,
		})
	})
	t.Run("description scope body footers", func(t *testing.T) {
		parseMessageHelper(t, dir, parser.ConventionalCommit{
			Body:        commitBody,
			CommitScope: commitScope,
			CommitType:  commitType,
			Description: commitDescription,
			Footers:     commitFooters,
		})
	})
	t.Run("breaking change description footers", func(t *testing.T) {
		parseMessageHelper(t, dir, parser.ConventionalCommit{
			BreakingChange: true,
			CommitType:     commitType,
			Description:    commitDescription,
			Footers:        commitFooters,
		})
	})
	t.Run("breaking change description body footers", func(t *testing.T) {
		parseMessageHelper(t, dir, parser.ConventionalCommit{
			Body:           commitBody,
			BreakingChange: true,
			CommitType:     commitType,
			Description:    commitDescription,
			Footers:        commitFooters,
		})
	})
	t.Run("breaking change description scope footers", func(t *testing.T) {
		parseMessageHelper(t, dir, parser.ConventionalCommit{
			BreakingChange: true,
			CommitScope:    commitScope,
			CommitType:     commitType,
			Description:    commitDescription,
			Footers:        commitFooters,
		})
	})
	t.Run("breaking change description scope body footers", func(t *testing.T) {
		parseMessageHelper(t, dir, parser.ConventionalCommit{
			Body:           commitBody,
			BreakingChange: true,
			CommitScope:    commitScope,
			CommitType:     commitType,
			Description:    commitDescription,
			Footers:        commitFooters,
		})
	})
	t.Run("description footers breaking change", func(t *testing.T) {
		parseMessageHelper(t, dir, parser.ConventionalCommit{
			BreakingChange: true,
			CommitType:     commitType,
			Description:    commitDescription,
			Footers:        breakingChangeFooter,
		})
	})
}
