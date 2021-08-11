// Package parser provides a simple parser for conventional commits
package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	errHeader    = errors.New("unable to parse commit header")
	errNoNewLine = errors.New("commit description not followed by an empty line")
)

// Commit represents a commit that adheres to the conventional commits specification
type Commit struct {
	Header Header
	Body   string
	Footer Footer

	BreakingChange bool

	FullCommit string
}

// Header represents Header in commit message
type Header struct {
	Type        string
	Scope       string
	Description string
	FullHeader  string
}

// Footer represents Footer in commit message
type Footer struct {
	Notes      []FooterNote
	FullFooter string
}

// FooterNote represents one footer note in Footer
type FooterNote struct {
	Token string
	Value string
}

func newFooterNote(token, value string) FooterNote {
	return FooterNote{Token: token, Value: value}
}

// Parse attempts to parse a commit message to a conventional commit
func Parse(message string) (*Commit, error) {
	message = strings.TrimRight(message, "\n\t ")
	messageLines := strings.Split(message, "\n")

	commit := &Commit{
		FullCommit: message,
	}
	currKeyValue := ""
	currFooterValue := ""

	foot := Footer{}

	inFooters := false
	for i, msgLine := range messageLines {
		switch i {
		case 0:
			err := parseHeader(msgLine, commit)
			if err != nil {
				return commit, err
			}
		case 1:
			if msgLine != "" {
				return commit, errNoNewLine
			}
		default:
			key, value := parseLineAsFooter(msgLine)

			if key != "" && value != "" {
				inFooters = true

				// Check if we have previously found a footer. If we have, set the current footer,
				// otherwise just record it.
				if currKeyValue != "" {
					foot.Notes = append(foot.Notes, newFooterNote(currKeyValue, currFooterValue))
					foot.FullFooter += messageLines[i-1] + "\n" // add previous line to FullFooter
				}
				currKeyValue = key
				currFooterValue = value
			} else {
				if inFooters {
					currFooterValue = fmt.Sprintf("%s\n%s", currFooterValue, msgLine)
				} else {
					if commit.Body == "" {
						commit.Body = msgLine
					} else {
						commit.Body += fmt.Sprintf("\n%s", msgLine)
					}
				}
			}
		}
	}

	// We reached the end of the commit message, so check if we need to record the footers
	if inFooters {
		foot.Notes = append(foot.Notes, newFooterNote(currKeyValue, currFooterValue))
		foot.FullFooter += messageLines[len(messageLines)-1]
	}

	// Remove whitespace in the Full Footer
	foot.FullFooter = strings.TrimSpace(foot.FullFooter)

	// Remove whitespace in the commit body
	commit.Body = strings.TrimSpace(commit.Body)
	commit.Footer = foot

	// Check if a footer contained a breaking change
	for _, footer := range commit.Footer.Notes {
		if footer.Token == "BREAKING CHANGE" || footer.Token == "BREAKING-CHANGE" {
			commit.BreakingChange = true
			break
		}
	}

	return commit, nil
}

// parseLineAsFooter attempts to parse the given line as a footer, returning both the key and the value of the header.
// If the line cannot be parsed then both return values will be empty.
func parseLineAsFooter(line string) (key, value string) {
	footerRegexp := regexp.MustCompile(`^(?:(BREAKING[- ]CHANGE|(?:[A-Za-z-])+): |((?:[A-Za-z-])+) #)(.+)$`)
	matches := footerRegexp.FindStringSubmatch(line)
	if len(matches) != 4 {
		return "", ""
	}

	if matches[1] == "" {
		return matches[2], matches[3]
	}
	return matches[1], matches[3]
}

// parseHeader attempts to parse the commit description line and set the appropriate values in the the given commit
func parseHeader(header string, commit *Commit) error {
	// from https://github.com/conventional-commits/parser#the-grammar

	// <header/summary> ::= <type>, "(", <scope>, ")", ["!"], ":", <whitespace>*, <text> <type>, ["!"], ":", <whitespace>*, <text>
	// <type>  ::= <any UTF8-octets except newline or parens or ":" or "!:" or whitespace>+
	// <scope> ::= <any UTF8-octets except newline or parens>+

	headerRegexp := regexp.MustCompile(`^(?P<type>[^\n\(\)(:|!:| )]+)(?:\((?P<scope>[^\n\(\)]+)\))?(?P<breaking>!)?: (?P<description>[^\n]+)(?:\n\s*\n(?P<body>(?:.|\n)*)(?:\n\s+\n(?P<footers>(?:[A-Za-z-]+: (?:.|\n)*)|(?:BREAKING CHANGE: (?:.|\n)*)|(?:[A-Za-z]+ \#(?:.|\n)*)))?)?$`)
	// TODO: comma separated multiple scopes?
	matches := headerRegexp.FindStringSubmatch(header)
	if matches == nil {
		return errHeader
	}

	head := Header{
		FullHeader: header,
	}

	names := headerRegexp.SubexpNames()
	for i, match := range matches {
		switch names[i] {
		case "type":
			head.Type = match
		case "scope":
			head.Scope = match
		case "description":
			head.Description = match
		case "breaking":
			commit.BreakingChange = (match == "!")
		}
	}

	commit.Header = head
	return nil
}
