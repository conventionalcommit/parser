package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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

var breakingChangeFooter = []Note{
	newNote("BREAKING CHANGE", "reason"),
}

var commitFooters = []Note{
	newNote("footer", "simple"),
	newNote("hash-footer", "123"),
}

var multiLineFooters = []Note{
	newNote("footer", `multi line footer
message is here
`),
	newNote("hash-footer", "123"),
}

func TestParserDescription(t *testing.T) {
	expectedCommit := &Commit{
		commitType:  commitType,
		description: commitDescription,
	}
	parseMsgAndCompare(t, "description", expectedCommit)
}

func TestParserDescriptionScope(t *testing.T) {
	expectedCommit := &Commit{
		commitType:  commitType,
		scope:       commitScope,
		description: commitDescription,
	}
	parseMsgAndCompare(t, "description_scope", expectedCommit)
}

func TestParserBreakingChangeDescription(t *testing.T) {
	expectedCommit := &Commit{
		commitType:       commitType,
		description:      commitDescription,
		isBreakingChange: true,
	}
	parseMsgAndCompare(t, "breaking_change_description", expectedCommit)
}

func TestParserBreakingChangeDescriptionScope(t *testing.T) {
	expectedCommit := &Commit{
		commitType:       commitType,
		scope:            commitScope,
		description:      commitDescription,
		isBreakingChange: true,
	}
	parseMsgAndCompare(t, "breaking_change_description_scope", expectedCommit)
}

func TestParserDescriptionBody(t *testing.T) {
	expectedCommit := &Commit{
		commitType:  commitType,
		description: commitDescription,
		body:        commitBody,
	}
	parseMsgAndCompare(t, "description_body", expectedCommit)
}

func TestParserDescriptionScopeBody(t *testing.T) {
	expectedCommit := &Commit{
		commitType:  commitType,
		scope:       commitScope,
		description: commitDescription,
		body:        commitBody,
	}
	parseMsgAndCompare(t, "description_scope_body", expectedCommit)
}

func TestParserBreakingChangeDescriptionBody(t *testing.T) {
	expectedCommit := &Commit{
		commitType:       commitType,
		description:      commitDescription,
		body:             commitBody,
		isBreakingChange: true,
	}
	parseMsgAndCompare(t, "breaking_change_description_body", expectedCommit)
}

func TestParserBreakingChangeDescriptionScopeBody(t *testing.T) {
	expectedCommit := &Commit{
		commitType:       commitType,
		scope:            commitScope,
		description:      commitDescription,
		body:             commitBody,
		isBreakingChange: true,
	}
	parseMsgAndCompare(t, "breaking_change_description_scope_body", expectedCommit)
}

func TestParserDescriptionFooters(t *testing.T) {
	expectedCommit := &Commit{
		commitType:  commitType,
		description: commitDescription,
		notes:       commitFooters,
	}
	parseMsgAndCompare(t, "description_footers", expectedCommit)
}

func TestParserDescriptionScopeFooters(t *testing.T) {
	expectedCommit := &Commit{
		commitType:  commitType,
		scope:       commitScope,
		description: commitDescription,
		notes:       commitFooters,
	}
	parseMsgAndCompare(t, "description_scope_footers", expectedCommit)
}

func TestParserDescriptionBodyFooters(t *testing.T) {
	expectedCommit := &Commit{
		commitType:  commitType,
		description: commitDescription,
		body:        commitBody,
		notes:       commitFooters,
	}
	parseMsgAndCompare(t, "description_body_footers", expectedCommit)
}

func TestParserDescriptionScopeBodyFooters(t *testing.T) {
	expectedCommit := &Commit{
		commitType:  commitType,
		scope:       commitScope,
		description: commitDescription,
		body:        commitBody,
		notes:       commitFooters,
	}
	parseMsgAndCompare(t, "description_scope_body_footers", expectedCommit)
}

func TestParserDescriptionFootersBreakingChange(t *testing.T) {
	expectedCommit := &Commit{
		commitType:       commitType,
		description:      commitDescription,
		notes:            breakingChangeFooter,
		isBreakingChange: true,
	}
	parseMsgAndCompare(t, "description_footers_breaking_change", expectedCommit)
}

func TestParserBreakingChangeDescriptionFooters(t *testing.T) {
	expectedCommit := &Commit{
		isBreakingChange: true,
		commitType:       commitType,
		description:      commitDescription,
		notes:            commitFooters,
	}
	parseMsgAndCompare(t, "breaking_change_description_footers", expectedCommit)
}

func TestParserBreakingChangeDescriptionBodyFooters(t *testing.T) {
	expectedCommit := &Commit{
		isBreakingChange: true,
		commitType:       commitType,
		description:      commitDescription,
		body:             commitBody,
		notes:            commitFooters,
	}
	parseMsgAndCompare(t, "breaking_change_description_body_footers", expectedCommit)
}

func TestParserBreakingChangeDescriptionScopeFooters(t *testing.T) {
	expectedCommit := &Commit{
		commitType:       commitType,
		scope:            commitScope,
		description:      commitDescription,
		notes:            commitFooters,
		isBreakingChange: true,
	}
	parseMsgAndCompare(t, "breaking_change_description_scope_footers", expectedCommit)
}

func TestParserBreakingChangeDescriptionScopeBodyFooters(t *testing.T) {
	expectedCommit := &Commit{
		commitType:       commitType,
		scope:            commitScope,
		description:      commitDescription,
		body:             commitBody,
		notes:            commitFooters,
		isBreakingChange: true,
	}
	parseMsgAndCompare(t, "breaking_change_description_scope_body_footers", expectedCommit)
}

func TestParserFooterMultiLine(t *testing.T) {
	expectedCommit := &Commit{
		commitType:  commitType,
		description: commitDescription,
		notes:       multiLineFooters,
	}
	parseMsgAndCompare(t, "footer_multi_line", expectedCommit)
}

func TestParserErrNoBlankLine(t *testing.T) {
	fileName := "err_no_blank_line"

	commitMsg, err := loadCommitMsgFromFile(filepath.Join(testDataDir, fileName))
	if err != nil {
		t.Error(err)
	}

	p := New()
	_, err = p.Parse(commitMsg)
	if err == nil {
		t.Errorf("no error: test file %v passed", fileName)
		return
	}

	if !IsNoBlankLineErr(err) {
		t.Error("error is not NoBlankLineErr error", err)
	}
}

func TestParserErrHeaderLine(t *testing.T) {
	fileName := "err_header_line"

	commitMsg, err := loadCommitMsgFromFile(filepath.Join(testDataDir, fileName))
	if err != nil {
		t.Error(err)
	}

	p := New()
	_, err = p.Parse(commitMsg)
	if err == nil {
		t.Errorf("no error: test file %v passed", fileName)
		return
	}

	if !IsHeaderErr(err) {
		t.Error("error is not HeaderErr error", err)
	}
}

func parseMsgAndCompare(t *testing.T, fileName string, expectedCommit *Commit) {
	commitMsg, err := loadCommitMsgFromFile(filepath.Join(testDataDir, fileName))
	if err != nil {
		t.Errorf("Received unexpected error:\n%+v", err)
		return
	}

	p := New()
	actualCommit, err := p.Parse(commitMsg)
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

func compareCommit(t *testing.T, a, b *Commit) bool {
	if a.commitType != b.commitType {
		t.Log("Header Type Not Equal")
		return false
	}
	if a.description != b.description {
		t.Log("Header Description Not Equal")
		return false
	}
	if a.scope != b.scope {
		t.Log("Header Scope Not Equal")
		return false
	}

	if a.body != b.body {
		t.Log("Body Not Equal")
		return false
	}

	notesA := a.notes
	notesB := b.notes

	if len(notesA) != len(notesB) {
		t.Log("Footer Notes Not Equal")
		return false
	}

	for index, aFoot := range notesA {
		bFoot := notesB[index]
		if aFoot.Token() != bFoot.Token() {
			t.Log("Footer Notes Token Not Equal", index, aFoot.Token(), bFoot.Token())
			return false
		}
		if aFoot.Value() != bFoot.Value() {
			t.Log("Footer Notes Value Not Equal", index, aFoot.Value(), bFoot.Value())
			return false
		}
	}

	return true
}
