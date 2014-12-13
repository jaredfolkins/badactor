package badactor

import "time"

type Sentence struct {
	Rule      *Rule
	ReleaseBy time.Time
	Start     time.Time
}

func NewSentence(r *Rule, sen time.Duration) *Sentence {
	return &Sentence{
		Rule:      r,
		ReleaseBy: time.Now().Add(sen),
		Start:     time.Now(),
	}
}
