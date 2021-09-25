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

// regex Compile everytime
// BenchmarkParser-4   	    7255	    156239 ns/op	  126473 B/op	     761 allocs/op

// regex Compile once
// BenchmarkParser-4   	  179227	      6531 ns/op	    1478 B/op	      23 allocs/op

// header regex cleaned
// BenchmarkParser-4   	  206452	      5199 ns/op	    1414 B/op	      23 allocs/op

func BenchmarkParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Parse(sampleCommit)
		if err != nil {
			b.Fatal(err)
		}
	}
}
