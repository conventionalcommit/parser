// Package parser provides a simple parser for conventional commits
package parser

import (
	"regexp"
	"strings"
)

// from https://github.com/conventional-commits/parser#the-grammar

// <header/summary> ::= <type>, "(", <scope>, ")", ["!"], ":", <whitespace>*, <text> <type>, ["!"], ":", <whitespace>*, <text>
// <type>  ::= <any UTF8-octets except newline or parens or ":" or "!:" or whitespace>+
// <scope> ::= <any UTF8-octets except newline or parens>+
// <description> ::= <any UTF8-octets except newline>*

const (
	headRegExStr = `^(?P<type>[^\n\(\)(:|!:| )]+)(?:\((?P<scope>[^\n\(\)]+)\))?(?P<breaking>!)?: (?P<description>[^\n]+)$`
	footRegExStr = `^(?:(BREAKING[- ]CHANGE|(?:[A-Za-z-])+): |((?:[A-Za-z-])+) #)(.+)$`
)

// Parser represent a conventional commit message parser
type Parser struct {
	headerRegex, footerRegex *regexp.Regexp
}

// New returns a new parser
func New() *Parser {
	headerRegex := regexp.MustCompile(headRegExStr)
	footerRegex := regexp.MustCompile(footRegExStr)

	return &Parser{
		headerRegex: headerRegex,
		footerRegex: footerRegex,
	}
}

// Parse attempts to parse a commit message to a conventional commit
func (p *Parser) Parse(message string) (*Commit, error) {
	return p.parse(message)
}

func (p *Parser) parse(message string) (*Commit, error) {
	c := &Commit{
		message: message,
	}

	message = strings.TrimRight(message, "\n\t ")
	messageLines := strings.Split(message, "\n")

	currKeyValue := ""
	currFooterValue := ""

	inFooters := false
	for i, msgLine := range messageLines {
		// First Line
		if i == 0 {
			err := p.parseHeader(c, msgLine)
			if err != nil {
				return nil, err
			}
			continue
		}

		// Second Line
		if i == 1 {
			if msgLine != "" {
				return nil, errNoBlankLine
			}
			continue
		}

		// Remaining Line
		key, value, isFooter := p.parseLineAsFooter(msgLine)
		// Is Footer
		if isFooter {
			inFooters = true

			// Check if we have previously found a footer. If we have, set the current footer,
			// otherwise just record it.
			if currKeyValue != "" {
				c.notes = append(c.notes, newNote(currKeyValue, currFooterValue))
				c.footer += messageLines[i-1] + "\n" // add previous line to FullFooter
			}

			currKeyValue = key
			currFooterValue = value
			continue
		}

		// Not a Footer Line
		if inFooters {
			currFooterValue = currFooterValue + "\n" + msgLine
		} else {
			if c.body == "" {
				c.body = msgLine
			} else {
				c.body += "\n" + msgLine
			}
		}
	}

	// We reached the end of the commit message, so check if we need to record the footers
	if inFooters {
		c.notes = append(c.notes, newNote(currKeyValue, currFooterValue))
		c.footer += messageLines[len(messageLines)-1]
	}

	// Remove whitespace in the Full Footer
	c.footer = strings.TrimSpace(c.footer)

	// Remove whitespace in the commit body
	c.body = strings.TrimSpace(c.body)

	// Check if a footer contains a breaking change
	for _, note := range c.notes {
		if note.Token() == "BREAKING CHANGE" || note.Token() == "BREAKING-CHANGE" {
			c.isBreakingChange = true
			break
		}
	}

	return c, nil
}

// parseLineAsFooter attempts to parse the given line as a footer, returning both the key and the value of the header.
// If the line cannot be parsed then isFooter is false
func (p *Parser) parseLineAsFooter(line string) (key, value string, isFooter bool) {
	matches := p.footerRegex.FindStringSubmatch(line)
	if len(matches) != 4 {
		return "", "", false
	}

	if matches[1] == "" {
		return matches[2], matches[3], true
	}
	return matches[1], matches[3], true
}

// parseHeader attempts to parse the commit description line and set the appropriate values in the the given commit
func (p *Parser) parseHeader(c *Commit, header string) error {
	matches := p.headerRegex.FindStringSubmatch(header)
	if matches == nil {
		return errHeader
	}

	c.header = header

	names := p.headerRegex.SubexpNames()
	for i, match := range matches {
		switch names[i] {
		case "type":
			c.commitType = match
		case "scope":
			// TODO: comma separated multiple scopes?
			c.scope = match
		case "description":
			c.description = match
		case "breaking":
			c.isBreakingChange = (match == "!")
		}
	}

	return nil
}
