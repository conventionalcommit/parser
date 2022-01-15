package parser

import (
	"strconv"
	"strings"
	"testing"
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

	p := New()

	for i, validCase := range validCases {
		testName := "case#" + strconv.Itoa(i+1)
		t.Run(testName, func(innerT *testing.T) {
			headerLine := strings.Split(validCase, "\n")[0]
			_, err := p.Parse(headerLine)
			if err != nil {
				innerT.Error("parseHeader failed for", headerLine, err)
				return
			}
		})
	}
}

func TestParseHeaderInvalid(t *testing.T) {
	var invalidCases = []string{
		`feat:() description with name.txt`,
		`feat:1 description with name.txt`,
		`feat:! description with name.txt`,
		`feat:A description with name.txt`,
		`feat123:A description with name.txt`,
		`feat!:A description with name.txt`,
		`feat())!:A description with name1.txt`,
		`feat(()!:A description with name2.txt`,
		`feat(scope1)!:A description with name.txt`,
		`!feat(scope1)!:A description with name.txt`,
		`feat(scope))!: A description with name.txt`,
		`feat((scope))!: A description with name.txt`,
		`feat((scope)!: A description with name.txt`,
		`feat((`,
		`feat():`,
		`feat):`,
	}

	p := New()

	for i, invalidCase := range invalidCases {
		testName := "case#" + strconv.Itoa(i+1)
		t.Run(testName, func(innerT *testing.T) {
			headerLine := strings.Split(invalidCase, "\n")[0]
			_, err := p.Parse(headerLine)
			if err == nil {
				innerT.Error("parseHeader passed without error for", headerLine)
			}
		})
	}
}
