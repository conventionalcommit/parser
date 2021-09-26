// Package parser provides a simple parser for conventional commits
package parser

import (
	"errors"
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

var (
	headerRegexp = regexp.MustCompile(headRegExStr)
	footerRegexp = regexp.MustCompile(footRegExStr)
)

var (
	errHeader      = errors.New("unable to parse commit header")
	errNoBlankLine = errors.New("commit description not followed by an empty line")
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
		// First Line
		if i == 0 {
			head, isBreak, err := parseHeader(msgLine)
			if err != nil {
				return nil, err
			}
			commit.Header = head
			commit.BreakingChange = isBreak
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
		key, value, isFooter := parseLineAsFooter(msgLine)
		if isFooter {
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
				currFooterValue = currFooterValue + "\n" + msgLine
			} else {
				if commit.Body == "" {
					commit.Body = msgLine
				} else {
					commit.Body += "\n" + msgLine
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

	// Check if a footer contains a breaking change
	for _, footer := range commit.Footer.Notes {
		if footer.Token == "BREAKING CHANGE" || footer.Token == "BREAKING-CHANGE" {
			commit.BreakingChange = true
			break
		}
	}

	return commit, nil
}

// parseLineAsFooter attempts to parse the given line as a footer, returning both the key and the value of the header.
// If the line cannot be parsed then isFooter is false
func parseLineAsFooter(line string) (key, value string, isFooter bool) {
	matches := footerRegexp.FindStringSubmatch(line)
	if len(matches) != 4 {
		return "", "", false
	}

	if matches[1] == "" {
		return matches[2], matches[3], true
	}
	return matches[1], matches[3], true
}

// parseHeader attempts to parse the commit description line and set the appropriate values in the the given commit
func parseHeader(header string) (Header, bool, error) {
	matches := headerRegexp.FindStringSubmatch(header)
	if matches == nil {
		return Header{}, false, errHeader
	}

	head := Header{
		FullHeader: header,
	}

	isBreakingChange := false

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
			isBreakingChange = (match == "!")
		}
	}

	return head, isBreakingChange, nil
}

// IsHeaderErr checks if given error is header parse error
func IsHeaderErr(err error) bool {
	return errors.Is(err, errHeader)
}

// IsNoBlankLineErr checks if given error is no new line error
func IsNoBlankLineErr(err error) bool {
	return errors.Is(err, errNoBlankLine)
}
