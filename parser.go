// Package parser provides a simple parser for conventional commits
package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// from https://github.com/conventional-commits/parser#the-grammar

// <header/summary> ::= <type>, "(", <scope>, ")", ["!"], ":", <whitespace>*, <text> <type>, ["!"], ":", <whitespace>*, <text>
// <type>  ::= <any UTF8-octets except newline or parens or ":" or "!:" or whitespace>+
// <scope> ::= <any UTF8-octets except newline or parens>+

const (
	headRegExStr = `^(?P<type>[^\n\(\)(:|!:| )]+)(?:\((?P<scope>[^\n\(\)]+)\))?(?P<breaking>!)?: (?P<description>[^\n]+)(?:\n\s*\n(?P<body>(?:.|\n)*)(?:\n\s+\n(?P<footers>(?:[A-Za-z-]+: (?:.|\n)*)|(?:BREAKING CHANGE: (?:.|\n)*)|(?:[A-Za-z]+ \#(?:.|\n)*)))?)?$`
	footRegExStr = `^(?:(BREAKING[- ]CHANGE|(?:[A-Za-z-])+): |((?:[A-Za-z-])+) #)(.+)$`
)

var (
	headerRegexp = regexp.MustCompile(headRegExStr)
	footerRegexp = regexp.MustCompile(footRegExStr)
)

var (
	errHeader    = errors.New("unable to parse commit header")
	errNoNewLine = errors.New("commit description not followed by an empty line")
)

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
			// TODO: comma separated multiple scopes?
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

// IsHeaderErr checks if given error is header parse error
func IsHeaderErr(err error) bool {
	return errors.Is(err, errHeader)
}
