# Parser

A go parser for [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) messages

[![PkgGoDev](https://pkg.go.dev/badge/github.com/conventionalcommit/parser)](https://pkg.go.dev/github.com/conventionalcommit/parser)![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/conventionalcommit/parser)

### Usage

```go
var msg = `feat(scope): description

this is first line in body

this is second line in body

Ref #123
Date: 01-01-2021
By: John Doe`

commit, err := Parse(msg)
if err != nil {
    fmt.Printf("Error: %s", err.Error())
}
fmt.Printf("%#v", commit)

/*
commitMsg = &parser.Commit{
    message:     "feat(scope): description\n\nthis is first line in body\n\nthis is second line in body\n\nRef #123\nDate: 01-01-2021\nBy: John Doe",
    header:      "feat(scope): description",
    body:        "this is first line in body\n\nthis is second line in body",
    footer:      "Ref #123\nDate: 01-01-2021\nBy: John Doe",
    commitType:  "feat",
    scope:       "scope",
    description: "description",
    notes:       {
        {token:"Ref", value:"123"},
        {token:"Date", value:"01-01-2021"},
        {token:"By", value:"John Doe"},
    },
    isBreakingChange: false,
}
*/
```

### TODO

- [ ] More Test Cases
- [ ] Benchmark

### Attribution

This parser is inspired and forked from

- [cc](https://github.com/zbindenren/cc) by [Rene Zbinden](https://github.com/zbindenren)
- [go-lexer](https://github.com/bbuck/go-lexer) by [Brandon Buck](https://github.com/bbuck)
- [cov-commit-parser](https://github.com/mbamber/cov-commit-parser) by [Matthew Bamber](https://github.com/mbamber)

### License

[MIT License](https://github.com/conventionalcommit/parser/tree/master/LICENSE.md)
