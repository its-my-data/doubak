package task

import (
	"flag"
	p "github.com/its-my-data/doubak/proto"
)

// Publisher contains the information used by the web publisher.
type Publisher struct {
	user       string
	categories []string
}

// NewPublisher returns a new parser task and initialise it.
func NewPublisher(categories []string) *Publisher {
	return &Publisher{
		user:       flag.Lookup(p.Flag_user.String()).Value.(flag.Getter).Get().(string),
		categories: categories,
	}
}

// Precheck validates the flags.
func (task *Publisher) Precheck() error {
	// TODO: check user existence, etc.
	return nil
}

// Execute starts publishing.
func (task *Publisher) Execute() error {
	return nil
}
