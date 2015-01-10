package badactor

import "time"

type jail struct {
	rule      *Rule
	releaseBy time.Time
	start     time.Time
}

func newJail(r *Rule, sen time.Duration) *jail {
	return &jail{
		rule:      r,
		releaseBy: time.Now().Add(sen),
		start:     time.Now(),
	}
}
