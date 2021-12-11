package parser

import "errors"

var (
	errHeader      = errors.New("unable to parse commit header")
	errNoBlankLine = errors.New("commit description not followed by an empty line")
)

// IsHeaderErr checks if given error is header parse error
func IsHeaderErr(err error) bool {
	return errors.Is(err, errHeader)
}

// IsNoBlankLineErr checks if given error is no new line error
func IsNoBlankLineErr(err error) bool {
	return errors.Is(err, errNoBlankLine)
}
