package badactor

import "time"

// Rule struct is used as a basic ruleset to Judge and Jail an Actor by
type Rule struct {
	Name        string
	Message     string
	StrikeLimit int
	ExpireBase  time.Duration
	Sentence    time.Duration
	Action      Action
}

// NewClassicRule returns a Rule with basic default values
func NewClassicRule(n string, m string) *Rule {
	return &Rule{
		Name:        n,
		Message:     m,
		StrikeLimit: 10,
		ExpireBase:  time.Minute * 10,
		Sentence:    time.Minute * 10,
	}
}
