package badactor

import "time"

type Infraction struct {
	Rule     *Rule
	Strikes  int
	ExpireBy time.Time
}

func NewInfraction(r *Rule) *Infraction {
	return &Infraction{
		Rule:     r,
		Strikes:  1,
		ExpireBy: time.Now().Add(r.ExpireBase),
	}
}

func (inf *Infraction) Rebase() {
	inf.ExpireBy = time.Now().Add(inf.Rule.ExpireBase)
}
