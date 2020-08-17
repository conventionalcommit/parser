# Cov Commit Parser
![CD](https://github.com/mbamber/cov-commit-parser/workflows/CD/badge.svg)

A simple parser for [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).

## Usage
When run without any arguments, `ccp version` will parse the commits at the current HEAD and output a single line containing a version number. This version number is the recommended version to use for the next build, based on the commit messages included since the latest tag on the branch.

The current version number is determined by the tag representing the latest semantic version. This can be overridden using the `--current` flag.

Use the `--since` flag to specify a commit hash, branch name or tag to adjust the commits used during the parsing.

### Use as a library
The tool can also be embedded into existing Go programs. The example below returns a new version based on the starting version `1.0.0` and using the single commit on the `HEAD` of a git repository in the directory `repo_path`:
```go
commitMessages, err := git.GetCommitsInDirectory("repo_path", "HEAD~1", "HEAD")
if err != nil {
    fmt.Printf("Error: %s", err.Error())
}

v, err := ccp.GetNextVersion("1.0.0", commitMessages, ccp.DefaultPatchTypes)
if err != nil {
    fmt.Printf("Error: %s", err.Error())
}

fmt.Println(v)
```
