package parser_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/conventionalcommit/parser"
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

func TestParserDescription(t *testing.T) {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Description: commitDescription,
		},
	}
	parseMsgAndCompare(t, "description", expectedCommit)
}

func TestParserDescriptionScope(t *testing.T) {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Scope:       commitScope,
			Description: commitDescription,
		},
	}
	parseMsgAndCompare(t, "description_scope", expectedCommit)
}

func TestParserBreakingChangeDescription(t *testing.T) {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Description: commitDescription,
		},
		BreakingChange: true,
	}
	parseMsgAndCompare(t, "breaking_change_description", expectedCommit)
}

func TestParserBreakingChangeDescriptionScope(t *testing.T) {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Scope:       commitScope,
			Description: commitDescription,
		},
		BreakingChange: true,
	}
	parseMsgAndCompare(t, "breaking_change_description_scope", expectedCommit)
}

func TestParserDescriptionBody(t *testing.T) {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Description: commitDescription,
		},
		Body: commitBody,
	}
	parseMsgAndCompare(t, "description_body", expectedCommit)
}

func TestParserDescriptionScopeBody(t *testing.T) {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Scope:       commitScope,
			Description: commitDescription,
		},
		Body: commitBody,
	}
	parseMsgAndCompare(t, "description_scope_body", expectedCommit)
}

func TestParserBreakingChangeDescriptionBody(t *testing.T) {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Description: commitDescription,
		},
		Body:           commitBody,
		BreakingChange: true,
	}
	parseMsgAndCompare(t, "breaking_change_description_body", expectedCommit)
}

func TestParserBreakingChangeDescriptionScopeBody(t *testing.T) {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Scope:       commitScope,
			Description: commitDescription,
		},
		Body:           commitBody,
		BreakingChange: true,
	}
	parseMsgAndCompare(t, "breaking_change_description_scope_body", expectedCommit)
}

func TestParserDescriptionFooters(t *testing.T) {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Description: commitDescription,
		},
		Footer: commitFooters,
	}
	parseMsgAndCompare(t, "description_footers", expectedCommit)
}

func TestParserDescriptionScopeFooters(t *testing.T) {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Scope:       commitScope,
			Description: commitDescription,
		},
		Footer: commitFooters,
	}
	parseMsgAndCompare(t, "description_scope_footers", expectedCommit)
}

func TestParserDescriptionBodyFooters(t *testing.T) {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Description: commitDescription,
		},
		Body:   commitBody,
		Footer: commitFooters,
	}
	parseMsgAndCompare(t, "description_body_footers", expectedCommit)
}

func TestParserDescriptionScopeBodyFooters(t *testing.T) {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Scope:       commitScope,
			Description: commitDescription,
		},
		Body:   commitBody,
		Footer: commitFooters,
	}
	parseMsgAndCompare(t, "description_scope_body_footers", expectedCommit)
}

func TestParserDescriptionFootersBreakingChange(t *testing.T) {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Description: commitDescription,
		},
		Footer:         breakingChangeFooter,
		BreakingChange: true,
	}
	parseMsgAndCompare(t, "description_footers_breaking_change", expectedCommit)
}

func TestParserBreakingChangeDescriptionFooters(t *testing.T) {
	expectedCommit := &parser.Commit{
		BreakingChange: true,
		Header: parser.Header{
			Type:        commitType,
			Description: commitDescription,
		},
		Footer: commitFooters,
	}
	parseMsgAndCompare(t, "breaking_change_description_footers", expectedCommit)
}

func TestParserBreakingChangeDescriptionBodyFooters(t *testing.T) {
	expectedCommit := &parser.Commit{
		BreakingChange: true,
		Header: parser.Header{
			Type:        commitType,
			Description: commitDescription,
		},
		Body:   commitBody,
		Footer: commitFooters,
	}
	parseMsgAndCompare(t, "breaking_change_description_body_footers", expectedCommit)
}

func TestParserBreakingChangeDescriptionScopeFooters(t *testing.T) {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Scope:       commitScope,
			Description: commitDescription,
		},
		Footer:         commitFooters,
		BreakingChange: true,
	}
	parseMsgAndCompare(t, "breaking_change_description_scope_footers", expectedCommit)
}

func TestParserBreakingChangeDescriptionScopeBodyFooters(t *testing.T) {
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
	parseMsgAndCompare(t, "breaking_change_description_scope_body_footers", expectedCommit)
}

func TestParserFooterMultiLine(t *testing.T) {
	expectedCommit := &parser.Commit{
		Header: parser.Header{
			Type:        commitType,
			Description: commitDescription,
		},
		Footer: multiLineFooters,
	}
	parseMsgAndCompare(t, "footer_multi_line", expectedCommit)
}

func TestParserErrNoBlankLine(t *testing.T) {
	fileName := "err_no_blank_line"

	commitMsg, err := loadCommitMsgFromFile(filepath.Join(testDataDir, fileName))
	if err != nil {
		t.Error(err)
	}

	_, err = parser.Parse(commitMsg)
	if err == nil {
		t.Errorf("no error: test file %v passed", fileName)
		return
	}

	if !parser.IsNoBlankLineErr(err) {
		t.Error("error is not NoBlankLineErr error", err)
	}
}

func TestParserErrHeaderLine(t *testing.T) {
	fileName := "err_header_line"

	commitMsg, err := loadCommitMsgFromFile(filepath.Join(testDataDir, fileName))
	if err != nil {
		t.Error(err)
	}

	_, err = parser.Parse(commitMsg)
	if err == nil {
		t.Errorf("no error: test file %v passed", fileName)
		return
	}

	if !parser.IsHeaderErr(err) {
		t.Error("error is not HeaderErr error", err)
	}
}

func parseMsgAndCompare(t *testing.T, fileName string, expectedCommit *parser.Commit) {
	commitMsg, err := loadCommitMsgFromFile(filepath.Join(testDataDir, fileName))
	if err != nil {
		t.Errorf("Received unexpected error:\n%+v", err)
		return
	}

	actualCommit, err := parser.Parse(commitMsg)
	if err != nil {
		t.Errorf("Received unexpected error:\n%+v", err)
		return
	}

	if !compareCommit(t, actualCommit, expectedCommit) {
		t.Errorf("Commit not equal :\n\tExpected: %v,\n\tActual: %v", expectedCommit, actualCommit)
		return
	}
}

// loadCommitMsgFromFile loads a file and returns the entire contents as a string. Any
// leading or trailing whitespace is removed
func loadCommitMsgFromFile(fileName string) (string, error) {
	out, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func compareCommit(t *testing.T, a, b *parser.Commit) bool {
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
