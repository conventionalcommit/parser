package parser_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/conventionalcommit/parser"
	"github.com/stretchr/testify/suite"
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

var breakingChangeFooter = parser.Footer{
	Notes: []parser.FooterNote{
		{
			Token: "BREAKING CHANGE",
			Value: "reason",
		},
	},
}

var commitFooters = parser.Footer{
	Notes: []parser.FooterNote{
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

var multiLineFooters = parser.Footer{
	Notes: []parser.FooterNote{
		{
			Token: "footer",
			Value: `multi line footer
message is here
`,
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
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Description: commitDescription,
		},
	}
	s.parseMsgAndCompare("description", expectedCommit)
}

func (s *parserSuite) TestDescriptionScope() {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Scope:       commitScope,
			Description: commitDescription,
		},
	}
	s.parseMsgAndCompare("description_scope", expectedCommit)
}

func (s *parserSuite) TestBreakingChangeDescription() {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Description: commitDescription,
		},
		BreakingChange: true,
	}
	s.parseMsgAndCompare("breaking_change_description", expectedCommit)
}

func (s *parserSuite) TestBreakingChangeDescriptionScope() {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Scope:       commitScope,
			Description: commitDescription,
		},
		BreakingChange: true,
	}
	s.parseMsgAndCompare("breaking_change_description_scope", expectedCommit)
}

func (s *parserSuite) TestDescriptionBody() {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Description: commitDescription,
		},
		Body: commitBody,
	}
	s.parseMsgAndCompare("description_body", expectedCommit)
}

func (s *parserSuite) TestDescriptionScopeBody() {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Scope:       commitScope,
			Description: commitDescription,
		},
		Body: commitBody,
	}
	s.parseMsgAndCompare("description_scope_body", expectedCommit)
}

func (s *parserSuite) TestBreakingChangeDescriptionBody() {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Description: commitDescription,
		},
		Body:           commitBody,
		BreakingChange: true,
	}
	s.parseMsgAndCompare("breaking_change_description_body", expectedCommit)
}

func (s *parserSuite) TestBreakingChangeDescriptionScopeBody() {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
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
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Description: commitDescription,
		},
		Footer: commitFooters,
	}
	s.parseMsgAndCompare("description_footers", expectedCommit)
}

func (s *parserSuite) TestDescriptionScopeFooters() {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Scope:       commitScope,
			Description: commitDescription,
		},
		Footer: commitFooters,
	}
	s.parseMsgAndCompare("description_scope_footers", expectedCommit)
}

func (s *parserSuite) TestDescriptionBodyFooters() {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Description: commitDescription,
		},
		Body:   commitBody,
		Footer: commitFooters,
	}
	s.parseMsgAndCompare("description_body_footers", expectedCommit)
}

func (s *parserSuite) TestDescriptionScopeBodyFooters() {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
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
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Description: commitDescription,
		},
		Footer:         breakingChangeFooter,
		BreakingChange: true,
	}
	s.parseMsgAndCompare("description_footers_breaking_change", expectedCommit)
}

func (s *parserSuite) TestBreakingChangeDescriptionFooters() {
	expectedCommit := &parser.Commit{
		BreakingChange: true,
		Header: parser.Header{
			Type:        commitType,
			Description: commitDescription,
		},
		Footer: commitFooters,
	}
	s.parseMsgAndCompare("breaking_change_description_footers", expectedCommit)
}

func (s *parserSuite) TestBreakingChangeDescriptionBodyFooters() {
	expectedCommit := &parser.Commit{
		BreakingChange: true,
		Header: parser.Header{
			Type:        commitType,
			Description: commitDescription,
		},
		Body:   commitBody,
		Footer: commitFooters,
	}
	s.parseMsgAndCompare("breaking_change_description_body_footers", expectedCommit)
}

func (s *parserSuite) TestBreakingChangeDescriptionScopeFooters() {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
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
	expectedCommit := &parser.Commit{
		Header: parser.Header{
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

func (s *parserSuite) TestFooterMultiLine() {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Description: commitDescription,
		},
		Footer: multiLineFooters,
	}
	s.parseMsgAndCompare("footer_multi_line", expectedCommit)
}

func (s *parserSuite) TestErrNoBlankLine() {
	t := s.T()

	fileName := "err_no_blank_line"

	commitMsg := s.loadCommitMsgFromFile(filepath.Join(testDataDir, fileName))
	_, err := parser.Parse(commitMsg)
	if err == nil {
		t.Errorf("no error: test file %v passed", fileName)
		return
	}

	if !parser.IsNoBlankLineErr(err) {
		t.Error("error is not NoBlankLineErr error", err)
	}
}

func (s *parserSuite) TestErrHeaderLine() {
	t := s.T()

	fileName := "err_header_line"

	commitMsg := s.loadCommitMsgFromFile(filepath.Join(testDataDir, fileName))
	_, err := parser.Parse(commitMsg)
	if err == nil {
		t.Errorf("no error: test file %v passed", fileName)
		return
	}

	if !parser.IsHeaderErr(err) {
		t.Error("error is not HeaderErr error", err)
	}
}

func (s *parserSuite) parseMsgAndCompare(fileName string, expectedCommit *parser.Commit) {
	t := s.T()
	t.Helper()

	commitMsg := s.loadCommitMsgFromFile(filepath.Join(testDataDir, fileName))
	actualCommit, err := parser.Parse(commitMsg)
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
		t.Errorf("error in test setup; unable to load file %s", fileName)
	}
	return strings.TrimSpace(string(out))
}

func (s *parserSuite) compareCommit(a, b *parser.Commit) bool {
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
