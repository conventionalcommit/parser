package parser

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHeaderValid(t *testing.T) {
	var validCases = []string{
		`feat: description with name.txt`,
		`feat: description with question?`,
		`feat: description with numbers 1, 2, 3 and 4?`,
		"feat: !@#$%^&*() ??>?///||| /\\", // just characters. why ?
		"feat: 123 description \n\n body 1, 2, 3 and 4?",
		"feat: ?123 description \n\n body 1, 2, 3 and 4?",
		"feat: description with body 1, \n\n2, 3 and 4?",
		"feat1234(@scope/scope1,scope2): description, \n\n body 1 2, 3 and 4?",
		"1245#feat1234(@scope/scope1,scope2): description, \n\n body 1 2, 3 and 4?",
	}

	for index, validCase := range validCases {
		testName := "case#" + strconv.Itoa(index+1)
		t.Run(testName, func(innerT *testing.T) {
			commit := &Commit{}
			err := parseHeader(validCase, commit)
			assert.NoError(innerT, err, validCase)
		})
	}
}

func TestParseHeaderInvalid(t *testing.T) {
	var validCases = []string{
		`feat:() description with name.txt`,
		`feat:1 description with name.txt`,
		`feat:! description with name.txt`,
		`feat:A description with name.txt`,
		`feat123:A description with name.txt`,
		`feat!:A description with name.txt`,
		`feat())!:A description with name.txt`,
		`feat(scope1)!:A description with name.txt`,
		`!feat(scope1)!:A description with name.txt`,
		`feat(scope))!: A description with name.txt`,
	}

	for index, validCase := range validCases {
		testName := "case#" + strconv.Itoa(index+1)
		t.Run(testName, func(innerT *testing.T) {
			commit := &Commit{}
			err := parseHeader(validCase, commit)
			assert.Error(innerT, err, validCase)
		})
	}
}
