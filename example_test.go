package parser_test

import (
	"fmt"

	"github.com/conventionalcommit/parser"
)

func ExampleParse() {
	var msg = `feat(scope): description

this is first line in body

this is second line in body

Ref #123
Date: 01-01-2021
By: John Doe`

	commit, err := parser.Parse(msg)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	fmt.Printf("%#v", commit)

	// Output: &parser.Commit{Header:parser.Header{Type:"feat", Scope:"scope", Description:"description", FullHeader:"feat(scope): description"}, Body:"this is first line in body\n\nthis is second line in body", Footer:parser.Footer{Notes:[]parser.FooterNote{parser.FooterNote{Token:"Ref", Value:"123"}, parser.FooterNote{Token:"Date", Value:"01-01-2021"}, parser.FooterNote{Token:"By", Value:"John Doe"}}, FullFooter:"Ref #123\nDate: 01-01-2021\nBy: John Doe"}, BreakingChange:false, FullCommit:"feat(scope): description\n\nthis is first line in body\n\nthis is second line in body\n\nRef #123\nDate: 01-01-2021\nBy: John Doe"}
}
