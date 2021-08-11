package parser

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// DefaultPatchTypes is the default list of commit types that should be treates as a patch change
	DefaultPatchTypes = []string{"fix"}
)

// ConventionalCommit represents a commit that adheres to the conventional commits specification
type ConventionalCommit struct {
	Body           string            `json:"body"`
	BreakingChange bool              `json:"breaking_change"`
	CommitScope    string            `json:"scope"`
	CommitType     string            `json:"type"`
	Description    string            `json:"description"`
	Footers        map[string]string `json:"footers"`
}

// ParseMessages attempts to parse a slice of commit messages to a slice of
// conventional commits. Returns a slice of errors to indicate all errors
// occurred during parsing
func ParseMessages(messages []string) ([]ConventionalCommit, []error) {
	commits := []ConventionalCommit{}
	errs := []error{}

	for _, m := range messages {
		c, err := ParseMessage(m)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		commits = append(commits, c)
	}
	return commits, errs
}

// ParseMessage attempts to parse a commit message to a conventional commit
func ParseMessage(message string) (ConventionalCommit, error) {
	messageLines := strings.Split(strings.TrimRight(message, "\n\t "), "\n")

	commit := ConventionalCommit{
		Footers: make(map[string]string),
	}
	currKeyValue := ""
	currFooterValue := ""

	inFooters := false
	for i, line := range messageLines {
		switch i {
		case 0:
			if err := parseHeader(line, &commit); err != nil {
				return commit, err
			}
		case 1:
			if line != "" {
				return commit, fmt.Errorf("commit description not followed by an empty line")
			}
		default:
			key, value := parseLineAsFooter(line)

			if key != "" && value != "" {
				inFooters = true

				// Check if we have previously found a footer. If we have, set the current footer,
				// otherwise just record it.
				if currKeyValue != "" {
					commit.Footers[currKeyValue] = currFooterValue
				}
				currKeyValue = key
				currFooterValue = value
			} else {
				if inFooters {
					currFooterValue = fmt.Sprintf("%s\n%s", currFooterValue, line)
				} else {
					if commit.Body == "" {
						commit.Body = line
					} else {
						commit.Body = commit.Body + fmt.Sprintf("\n%s", line)
					}
				}
			}
		}
	}

	// We reached the end of the commit message, so check if we need to record the footers
	if inFooters {
		commit.Footers[currKeyValue] = currFooterValue
	}

	// Remove whitespace in the commit body
	commit.Body = strings.TrimSpace(commit.Body)

	// Check if a footer contained a breaking change
	for footer := range commit.Footers {
		if footer == "BREAKING CHANGE" || footer == "BREAKING-CHANGE" {
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
func parseHeader(header string, commit *ConventionalCommit) error {
	headerRegexp := regexp.MustCompile(`^(?P<type>[A-Za-z]+)(?:\((?P<scope>[A-Za-z]+)\))?(?P<breaking>!)?: (?P<description>[\w| ]+)(?:\n\s*\n(?P<body>(?:.|\n)*)(?:\n\s+\n(?P<footers>(?:[A-Za-z-]+: (?:.|\n)*)|(?:BREAKING CHANGE: (?:.|\n)*)|(?:[A-Za-z]+ \#(?:.|\n)*)))?)?$`)

	matches := headerRegexp.FindStringSubmatch(header)
	if matches == nil {
		return fmt.Errorf("unable to parse commit header: %s", header)
	}

	names := headerRegexp.SubexpNames()
	for i, match := range matches {
		switch names[i] {
		case "type":
			commit.CommitType = match
		case "scope":
			commit.CommitScope = match
		case "description":
			commit.Description = match
		case "breaking":
			commit.BreakingChange = (match == "!")
		}
	}

	return nil
}
