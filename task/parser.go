package task

import (
	"flag"
	p "github.com/its-my-data/doubak/proto"
)

// Parser contains the information used by the parser.
type Parser struct {
	user       string
	categories []string
}

// NewParser returns a new parser task and initialise it.
func NewParser(categories []string) *Parser {
	return &Parser{
		user:       flag.Lookup(p.Flag_user.String()).Value.(flag.Getter).Get().(string),
		categories: categories,
	}
}

// Precheck validates the flags.
func (task *Parser) Precheck() error {
	// TODO: check user existence, etc.
	return nil
}

// Execute starts parsing.
func (task *Parser) Execute() error {
	return nil
}
