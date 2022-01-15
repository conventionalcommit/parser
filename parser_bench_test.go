package parser

import (
	"testing"
)

var sampleCommit = `feat(scope): description

this is first line in body

this is second line in body

Ref: #123
Date: 01-01-2021
By: John Doe`

// regex based parser - last version
// BenchmarkParser-4   	  229952	      4875 ns/op	    1365 B/op	      21 allocs/op

// lexer based parser
// BenchmarkParser-4   	   77244	     15265 ns/op	    6448 B/op	     306 allocs/op

// lexer: remove linked list stack
// BenchmarkParser-4   	  152005	      7352 ns/op	    2928 B/op	      60 allocs/op

var dumpRes *Commit

var p = New()

func BenchmarkParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r, err := p.Parse(sampleCommit)
		if err != nil {
			b.Error(err)
			return
		}
		dumpRes = r
	}
}
