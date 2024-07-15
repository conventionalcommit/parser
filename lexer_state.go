package parser

import (
	"errors"
	"fmt"
	"unicode"
)

const (
	breakingTokenSpace  = "BREAKING CHANGE"
	breakingTokenHyphen = "BREAKING-CHANGE"
)

var (
	errMissingScopeOrDesc     = errors.New("header: missing scope or description")
	errScopeMissingParen      = errors.New("scope should end with ')'")
	errScopeEmpty             = errors.New("scope is empty")
	errDescMissingDelimiter   = errors.New("scope must be followed by ': '")
	errHeaderMissingEmptyLine = errors.New("at least one empty line required after header")
	errBodyEmptyLine          = errors.New("at least one empty line required after body")

	errScopeInvalidChar = "scope: invalid character '%c'"
	errTypeInvalidChar  = "type: invalid character '%c'"
)

// all lexer token types that are emitted.
const (
	_ tokenType = iota

	headerTypeToken
	headerScopeToken

	leftScopeDelimiterToken
	rightScopeDelimiterToken
	breakingChangeToken

	descDelimiterToken
	descriptionToken

	bodyToken

	footerDelimterToken
	footerKeyToken
	footerValueToken
)

func typeState(l *lexer) stateFunc {
	for {
		r := l.Peek()

		if r == eof {
			l.Error(errMissingScopeOrDesc)
			return nil
		}

		if r == ':' || r == '!' {
			l.Emit(headerTypeToken)
			return descriptionDelimiterState
		}

		if r == '(' {
			l.Emit(headerTypeToken)
			l.TakeNext('(')
			l.Emit(leftScopeDelimiterToken)
			return scopeState
		}

		if !isValidTypeChar(r) {
			l.Error(fmt.Errorf(errTypeInvalidChar, r))
			return nil
		}

		l.Next()
	}
}

func scopeState(l *lexer) stateFunc {
	for {
		r := l.Peek()

		if r == eof {
			l.Error(errScopeMissingParen)
			return nil
		}

		if r == ')' {
			if l.Current() == "" {
				l.Error(errScopeEmpty)
				return nil
			}

			l.Emit(headerScopeToken)
			l.TakeNext(')')
			l.Emit(rightScopeDelimiterToken)

			return descriptionDelimiterState
		}

		if !isValidScopeChar(r) {
			l.Error(fmt.Errorf(errScopeInvalidChar, r))
			return nil
		}

		l.Next()
	}
}

func descriptionDelimiterState(l *lexer) stateFunc {
	if l.Peek() == '!' {
		l.Next()
		l.Emit(breakingChangeToken)
	}

	l.Next()

	if l.Current() != ":" || l.Peek() != ' ' {
		l.Error(errDescMissingDelimiter)
		return nil
	}
	l.Next()

	l.Emit(descDelimiterToken)

	return descriptionState
}

func descriptionState(l *lexer) stateFunc {
	for {
		r := l.Peek()

		if r == eof {
			l.Emit(descriptionToken)
			return nil
		}

		if r == '\n' {
			l.Emit(descriptionToken)
			return headerDelimeterState
		}

		l.Next()
	}
}

func headerDelimeterState(l *lexer) stateFunc {
	l.Take("\n")

	if len(l.Current()) < 2 {
		l.Error(errHeaderMissingEmptyLine)
		return nil
	}

	l.Ignore()

	return bodyOrFooterState
}

func bodyOrFooterState(l *lexer) stateFunc {
	count, isFooter := checkIfFooterToken(l)

	// there is no body
	if isFooter {
		rewind(l, count)
		return footerTokenState
	}

	return bodyState
}

func bodyState(l *lexer) stateFunc {
	found := takeUntilFirstFooterToken(l)
	if !found {
		l.Emit(bodyToken)
		return nil
	}

	// go back to the last newline character
	for {
		l.Rewind()

		if l.Peek() != '\n' {
			break
		}
	}
	l.Next()
	l.Emit(bodyToken)

	return bodyDelimiterState
}

func bodyDelimiterState(l *lexer) stateFunc {
	l.Take("\n")

	if len(l.Current()) < 2 {
		l.Error(errBodyEmptyLine)
		return nil
	}

	l.Ignore()

	return footerTokenState
}

func footerTokenState(l *lexer) stateFunc {
	l.Take("\n")
	l.Ignore()

	checkIfFooterToken(l)
	l.Emit(footerKeyToken)

	return footerDelimiterState
}

func footerValueState(l *lexer) stateFunc {
	if l.Peek() == eof {
		return nil
	}

	found := takeUntilFirstFooterToken(l)
	l.Emit(footerValueToken)

	if !found {
		return nil
	}

	return footerTokenState
}

func footerDelimiterState(l *lexer) stateFunc {
	l.Take(": #")
	l.Emit(footerDelimterToken)

	return footerValueState
}

// takeUntilFirstFooter takes all characters until a footer token is detected
func takeUntilFirstFooterToken(l *lexer) bool {
	for {
		r := l.Next()

		if r == eof {
			return false
		}

		// a footer token has to begin at the start of a line
		if r == '\n' {
			count, isFooter := checkIfFooterToken(l)
			if isFooter {
				// if count is > 0 we are at the end of the footer token
				// for i := 0; i < count; i++ {
				// 	l.Rewind()
				// }
				rewind(l, count)

				return true
			}
		}
	}
}

func rewind(l *lexer, count int) {
	for i := count; i > 0; i-- {
		l.Rewind()
	}
}

// takeFooterToken continues over each consecutive character in the source
// until an invalid footer token character is detected. The method returns
// the length if there is a valid footer token found. If it is not a valid footer
// token, 0 is returned.

// token: "BREAKING CHANGE" | "BREAKING-CHANGE" | <any UTF8-octets except newline or parens or ":" or "!:" or whitespace>+
func checkIfFooterToken(l *lexer) (int, bool) {
	// handle BREAKING CHANGE
	if l.Peek() == 'B' {
		// BREAKING-CHANGE: or BREAKING CHANGE:
		candidate := peekString(l, len(breakingTokenSpace)+2)
		if candidate == breakingTokenSpace+": " || candidate == breakingTokenHyphen+": " {
			l.Emit(breakingChangeToken)
			l.Take(breakingTokenSpace + "-")
			return len(breakingTokenSpace), true
		}
	}

	count := 0
	r := l.Next()

	for unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
		count++
		r = l.Next()
	}

	nextChar := l.Peek()

	// :<space> or <space># delimiter
	if (r == ':' && nextChar == ' ') || (r == ' ' && nextChar == '#') {

		l.Rewind()
		return count, true
	}

	l.Rewind()

	return 0, false
}

func peekString(l *lexer, count int) string {
	s := ""
	i := 0

	defer func() {
		for j := 0; j < i; j++ {
			l.Rewind()
		}
	}()

	for i = 0; i < count; i++ {
		if l.Peek() == eof {
			return s
		}
		s += string(l.Next())
	}

	return s
}

// from https://github.com/conventional-commits/parser#the-grammar
// <header/summary> ::= <type>, "(", <scope>, ")", ["!"], ":", <whitespace>*, <text> <type>, ["!"], ":", <whitespace>*, <text>

// <type>  ::= <any UTF8-octets except newline or parens or ":" or "!:" or whitespace>+
func isValidTypeChar(r rune) bool {
	switch r {
	case '\n', '(', ')', ':':
		return false
	default:
		if unicode.IsSpace(r) {
			return false
		}
		return true
	}
}

// <scope> ::= <any UTF8-octets except newline or parens>+
func isValidScopeChar(r rune) bool {
	switch r {
	case '\n', '(', ')':
		return false
	default:
		return true
	}
}
