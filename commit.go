package parser

// Commit represents a commit that adheres to the conventional commits specification
type Commit struct {
	Header Header
	Body   string
	Footer Footer

	BreakingChange bool

	FullCommit string
}

// Header represents Header in commit message
type Header struct {
	Type        string
	Scope       string
	Description string

	FullHeader string
}

// Footer represents Footer in commit message
type Footer struct {
	Notes []FooterNote

	FullFooter string
}

// FooterNote represents one footer note in Footer
type FooterNote struct {
	Token string
	Value string
}

func newFooterNote(token, value string) FooterNote {
	return FooterNote{Token: token, Value: value}
}
