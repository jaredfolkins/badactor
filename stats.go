package badactor

import "time"

type Stats struct {
	ActorName  string
	RuleName   string
	AccessedAt time.Time
}

func NewStats(a *actor, r *Rule) *Stats {
	return &Stats{
		ActorName:  a.name,
		RuleName:   r.Name,
		AccessedAt: time.Now(),
	}
}
