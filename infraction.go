package badactor

import "time"

type infraction struct {
	rule     *Rule
	strikes  int
	expireBy time.Time
}

func NewInfraction(r *Rule) *infraction {
	return &infraction{
		rule:     r,
		strikes:  0,
		expireBy: time.Now().Add(r.ExpireBase),
	}
}

func (inf *infraction) rebase() {
	inf.expireBy = time.Now().Add(inf.rule.ExpireBase)
}
