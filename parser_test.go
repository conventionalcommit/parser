package parser_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	. "github.com/conventionalcommit/parser"
)

const (
	commitBody = `This is a multiline commit body.

This is the second line`
)

const (
	commitDescription = "description message"
	commitScope       = "scope"
	commitType        = "type"

	testDataDir = "testdata"
)

var breakingChangeFooter = Footer{
	Notes: []FooterNote{
		{
			Token: "BREAKING CHANGE",
			Value: "reason",
		},
	},
}

var commitFooters = Footer{
	Notes: []FooterNote{
		{
			Token: "footer",
			Value: "simple",
		},
		{
			Token: "hash-footer",
			Value: "123",
		},
	},
}

func TestParser(t *testing.T) {
	ps := &parserSuite{}
	suite.Run(t, ps)
}

type parserSuite struct {
	suite.Suite
}

func (s *parserSuite) TestDescription() {
	expectedCommit := &Commit{
		Header: Header{
			Type:        commitType,
			Description: commitDescription,
		},
	}
	s.parseMsgAndCompare("description", expectedCommit)
}

func (s *parserSuite) TestDescriptionScope() {
	expectedCommit := &Commit{
		Header: Header{
			Type:        commitType,
			Scope:       commitScope,
			Description: commitDescription,
		},
	}
	s.parseMsgAndCompare("description_scope", expectedCommit)
}

func (s *parserSuite) TestBreakingChangeDescription() {
	expectedCommit := &Commit{
		Header: Header{
			Type:        commitType,
			Description: commitDescription,
		},
		BreakingChange: true,
	}
	s.parseMsgAndCompare("breaking_change_description", expectedCommit)
}

func (s *parserSuite) TestBreakingChangeDescriptionScope() {
	expectedCommit := &Commit{
		Header: Header{
			Type:        commitType,
			Scope:       commitScope,
			Description: commitDescription,
		},
		BreakingChange: true,
	}
	s.parseMsgAndCompare("breaking_change_description_scope", expectedCommit)
}

func (s *parserSuite) TestDescriptionBody() {
	expectedCommit := &Commit{
		Header: Header{
			Type:        commitType,
			Description: commitDescription,
		},
		Body: commitBody,
	}
	s.parseMsgAndCompare("description_body", expectedCommit)
}

func (s *parserSuite) TestDescriptionScopeBody() {
	expectedCommit := &Commit{
		Header: Header{
			Type:        commitType,
			Scope:       commitScope,
			Description: commitDescription,
		},
		Body: commitBody,
	}
	s.parseMsgAndCompare("description_scope_body", expectedCommit)
}

func (s *parserSuite) TestBreakingChangeDescriptionBody() {
	expectedCommit := &Commit{
		Header: Header{
			Type:        commitType,
			Description: commitDescription,
		},
		Body:           commitBody,
		BreakingChange: true,
	}
	s.parseMsgAndCompare("breaking_change_description_body", expectedCommit)
}

func (s *parserSuite) TestBreakingChangeDescriptionScopeBody() {
	expectedCommit := &Commit{
		Header: Header{
			Type:        commitType,
			Scope:       commitScope,
			Description: commitDescription,
		},
		Body:           commitBody,
		BreakingChange: true,
	}
	s.parseMsgAndCompare("breaking_change_description_scope_body", expectedCommit)
}

func (s *parserSuite) TestDescriptionFooters() {
	expectedCommit := &Commit{
		Header: Header{
			Type:        commitType,
			Description: commitDescription,
		},
		Footer: commitFooters,
	}
	s.parseMsgAndCompare("description_footers", expectedCommit)
}

func (s *parserSuite) TestDescriptionScopeFooters() {
	expectedCommit := &Commit{
		Header: Header{
			Type:        commitType,
			Scope:       commitScope,
			Description: commitDescription,
		},
		Footer: commitFooters,
	}
	s.parseMsgAndCompare("description_scope_footers", expectedCommit)
}

func (s *parserSuite) TestDescriptionBodyFooters() {
	expectedCommit := &Commit{
		Header: Header{
			Type:        commitType,
			Description: commitDescription,
		},
		Body:   commitBody,
		Footer: commitFooters,
	}
	s.parseMsgAndCompare("description_body_footers", expectedCommit)
}

func (s *parserSuite) TestDescriptionScopeBodyFooters() {
	expectedCommit := &Commit{
		Header: Header{
			Type:        commitType,
			Scope:       commitScope,
			Description: commitDescription,
		},
		Body:   commitBody,
		Footer: commitFooters,
	}
	s.parseMsgAndCompare("description_scope_body_footers", expectedCommit)
}

func (s *parserSuite) TestDescriptionFootersBreakingChange() {
	expectedCommit := &Commit{
		Header: Header{
			Type:        commitType,
			Description: commitDescription,
		},
		Footer:         breakingChangeFooter,
		BreakingChange: true,
	}
	s.parseMsgAndCompare("description_footers_breaking_change", expectedCommit)
}

func (s *parserSuite) TestBreakingChangeDescriptionFooters() {
	expectedCommit := &Commit{
		BreakingChange: true,
		Header: Header{
			Type:        commitType,
			Description: commitDescription,
		},
		Footer: commitFooters,
	}
	s.parseMsgAndCompare("breaking_change_description_footers", expectedCommit)
}

func (s *parserSuite) TestBreakingChangeDescriptionBodyFooters() {
	expectedCommit := &Commit{
		BreakingChange: true,
		Header: Header{
			Type:        commitType,
			Description: commitDescription,
		},
		Body:   commitBody,
		Footer: commitFooters,
	}
	s.parseMsgAndCompare("breaking_change_description_body_footers", expectedCommit)
}

func (s *parserSuite) TestBreakingChangeDescriptionScopeFooters() {
	expectedCommit := &Commit{
		Header: Header{
			Type:        commitType,
			Scope:       commitScope,
			Description: commitDescription,
		},
		Footer:         commitFooters,
		BreakingChange: true,
	}
	s.parseMsgAndCompare("breaking_change_description_scope_footers", expectedCommit)
}

func (s *parserSuite) TestBreakingChangeDescriptionScopeBodyFooters() {
	expectedCommit := &Commit{
		Header: Header{
			Type:        commitType,
			Scope:       commitScope,
			Description: commitDescription,
		},
		Body:           commitBody,
		Footer:         commitFooters,
		BreakingChange: true,
	}
	s.parseMsgAndCompare("breaking_change_description_scope_body_footers", expectedCommit)
}

func (s *parserSuite) parseMsgAndCompare(fileName string, expectedCommit *Commit) {
	t := s.T()
	t.Helper()

	commitMsg := s.loadCommitMsgFromFile(filepath.Join(testDataDir, fileName))
	actualCommit, err := Parse(commitMsg)
	if err != nil {
		t.Errorf("Received unexpected error:\n%+v", err)
		return
	}

	if !s.compareCommit(actualCommit, expectedCommit) {
		t.Errorf("Commit not equal :\n\tExpected: %v,\n\tActual: %v", expectedCommit, actualCommit)
		return
	}
}

// loadCommitMsgFromFile loads a file and returns the entire contents as a string. Any
// leading or trailing whitespace is removed
func (s *parserSuite) loadCommitMsgFromFile(fileName string) string {
	t := s.T()
	t.Helper()

	out, err := os.ReadFile(fileName)
	if err != nil {
		assert.Failf(t, "error in test setup", "unable to load file %s", fileName)
	}
	return strings.TrimSpace(string(out))
}

func (s *parserSuite) compareCommit(a, b *Commit) bool {
	t := s.T()

	if a.Header.Type != b.Header.Type {
		t.Log("Header Type Not Equal")
		return false
	}
	if a.Header.Description != b.Header.Description {
		t.Log("Header Description Not Equal")
		return false
	}
	if a.Header.Scope != b.Header.Scope {
		t.Log("Header Scope Not Equal")
		return false
	}

	if a.Body != b.Body {
		t.Log("Body Not Equal")
		return false
	}

	notesA := a.Footer.Notes
	notesB := b.Footer.Notes

	if len(notesA) != len(notesB) {
		t.Log("Footer Notes Not Equal")
		return false
	}

	for index, aFoot := range notesA {
		bFoot := notesB[index]
		if aFoot.Token != bFoot.Token {
			t.Log("Footer Notes Token Not Equal", index, aFoot.Token, bFoot.Token)
			return false
		}
		if aFoot.Value != bFoot.Value {
			t.Log("Footer Notes Value Not Equal", index, aFoot.Value, bFoot.Value)
			return false
		}
	}

	return true
}
