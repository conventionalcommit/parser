package parser

// Commit represents a commit that adheres to the conventional commits specification
type Commit struct {
	message string

	header string
	body   string
	footer string

	commitType  string
	scope       string
	description string
	notes       []Note

	isBreakingChange bool
}

// Message returns input commit message
func (c *Commit) Message() string {
	return c.message
}

// Header returns header of the commit
func (c *Commit) Header() string {
	return c.header
}

// Body returns body of the commit
func (c *Commit) Body() string {
	return c.body
}

// Footer returns footer of the commit
func (c *Commit) Footer() string {
	return c.footer
}

// Type returns type of the commit
func (c *Commit) Type() string {
	return c.commitType
}

// Scope returns scope of the commit
func (c *Commit) Scope() string {
	return c.scope
}

// Description returns description of the commit
func (c *Commit) Description() string {
	return c.description
}

// Notes returns footer notes of the commit
func (c *Commit) Notes() []Note {
	return c.notes
}

// IsBreakingChange returns true if breaking change
func (c *Commit) IsBreakingChange() bool {
	return c.isBreakingChange
}

// Note represents one footer note
type Note struct {
	token string
	value string
}

func newNote(token, value string) Note {
	return Note{
		token: token,
		value: value,
	}
}

func (n *Note) Token() string {
	return n.token
}

func (n *Note) Value() string {
	return n.value
}
