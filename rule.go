package badactor

import "time"

// A Rule is a worker and the programmer uses them
// to enforce the desired behavior
type Rule struct {
	Name        string
	Message     string
	StrikeLimit int
	ExpireBase  time.Duration
	Sentence    time.Duration
}

func NewClassicRule(n string, m string) *Rule {
	return &Rule{
		Name:        n,
		Message:     m,
		StrikeLimit: 10,
		ExpireBase:  time.Minute * 10,
		Sentence:    time.Minute * 10,
	}
}
