package badactor

import "time"

type sentence struct {
	rule      *Rule
	releaseBy time.Time
	start     time.Time
}

func newSentence(r *Rule, sen time.Duration) *sentence {
	return &sentence{
		rule:      r,
		releaseBy: time.Now().Add(sen),
		start:     time.Now(),
	}
}
