// Package parser provides a parser for conventional commits
package parser

import (
	"strings"
)

// Parser represent a conventional commits parser
type Parser struct{}

// New returns a new Parser instance
func New() *Parser {
	return &Parser{}
}

// Parse parses the conventional commit. If it fails, an error is returned.
func (p *Parser) Parse(input string) (*Commit, error) {
	input = strings.TrimSpace(input)
	return p.parse(input)
}

func (p *Parser) parse(input string) (*Commit, error) {
	lex := newLexer(input, typeState, func(error) {})
	lex.Start()

	c := &Commit{
		message: input,
	}

	footerCount := 0

	footerStartPos := 0
	footerEndPos := 0

	for {
		t, done := lex.NextToken()
		if done {
			break
		}

		switch t.Type {
		case breakingChangeToken:
			c.isBreakingChange = true
		case headerTypeToken:
			c.commitType = t.Value
		case headerScopeToken:
			c.scope = t.Value
		case descriptionToken:
			c.description = t.Value
		case bodyToken:
			c.header = strings.TrimSpace(lex.Get(0, t.Start))
			c.body = strings.TrimSpace(t.Value)
		case footerKeyToken:
			if footerStartPos == 0 {
				footerStartPos = t.Start
			}
			n := Note{
				token: t.Value,
			}
			c.notes = append(c.notes, n)
		case footerValueToken:
			c.notes[footerCount].value = strings.TrimSpace(t.Value)
			footerCount++
			footerEndPos = t.End
		}
	}

	if lex.Err() != nil {
		return nil, lex.Err()
	}

	if footerStartPos != 0 {
		c.footer = strings.TrimSpace(lex.Get(footerStartPos, footerEndPos))
	}

	return c, nil
}
