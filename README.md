# Parser

A simple go parser for [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)

[![PkgGoDev](https://pkg.go.dev/badge/github.com/conventionalcommit/parser)](https://pkg.go.dev/github.com/conventionalcommit/parser)

### Usage

```go
var msg = `feat(scope): description

this is first line in body

this is second line in body

Ref: #123
Date: 01-01-2021
By: John Doe`

commit, err := Parse(msg)
if err != nil {
    fmt.Printf("Error: %s", err.Error())
}
fmt.Printf("%#v", commit)

/*
commitMsg = &parser.Commit{
  Header: parser.Header{
    Type: "feat",
    Scope: "scope",
    Description: "description",
    FullHeader: "feat(scope): description",
  },
  Body: "this is first line in body\n\nthis is second line in body",
  Footer: parser.Footer{
    Notes: []parser.FooterNote{
      parser.FooterNote{
        Token: "Ref",
        Value: "#123",
      },
      parser.FooterNote{
        Token: "Date",
        Value: "01-01-2021",
      },
      parser.FooterNote{
        Token: "By",
        Value: "John Doe",
      },
    },
    FullFooter: "Ref: #123\nDate: 01-01-2021\nBy: John Doe",
  },
  BreakingChange: false,
  FullCommit: "feat(scope): description\n\nthis is first line in body\n\nthis is second line in body\n\nRef: #123\nDate: 01-01-2021\nBy: John Doe",
}
*/
```

### Fork

This parser is a fork of [cov-commit-parser](github.com/mbamber/cov-commit-parser) by [Matthew Bamber](github.com/mbamber/)

### License

[MIT License](https://github.com/conventionalcommit/parser/tree/master/LICENSE.md)

