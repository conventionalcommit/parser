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

	p := parser.New()
	commit, err := p.Parse(msg)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	fmt.Printf("%#v", commit)

	// Output: &parser.Commit{message:"feat(scope): description\n\nthis is first line in body\n\nthis is second line in body\n\nRef #123\nDate: 01-01-2021\nBy: John Doe", header:"feat(scope): description", body:"this is first line in body\n\nthis is second line in body", footer:"Ref #123\nDate: 01-01-2021\nBy: John Doe", commitType:"feat", scope:"scope", description:"description", notes:[]parser.Note{parser.Note{token:"Ref", value:"123"}, parser.Note{token:"Date", value:"01-01-2021"}, parser.Note{token:"By", value:"John Doe"}}, isBreakingChange:false}
}
