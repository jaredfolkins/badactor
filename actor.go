package badactor

import (
	"fmt"
	"time"
)

// ttl is the time to live value for newly created actors
const ttl = 100

type Actor struct {
	name        string
	infractions map[string]*infraction
	jails       map[string]*jail
	director    *Director
	ttl         time.Time
	accessedAt  time.Time
}

func newActor(n string, d *Director) *Actor {
	a := &Actor{
		director:    d,
		ttl:         time.Now().Add(time.Millisecond * ttl),
		name:        n,
		infractions: make(map[string]*infraction),
		jails:       make(map[string]*jail),
		accessedAt:  time.Now(),
	}
	return a
}

func newClassicActor(n string, r *Rule, d *Director) *Actor {
	a := &Actor{
		director:    d,
		ttl:         time.Now().Add(time.Millisecond * ttl),
		name:        n,
		infractions: make(map[string]*infraction),
		jails:       make(map[string]*jail),
		accessedAt:  time.Now(),
	}

	a.infractions[r.Name] = newInfraction(r)
	return a
}

func (a *Actor) rebaseAll() error {

	for _, inf := range a.infractions {
		inf.rebase()
	}

	a.ttl = time.Now().Add(time.Millisecond * ttl)

	return nil
}

func (a *Actor) infraction(rn string) error {

	if a.isJailedFor(rn) {
		return fmt.Errorf("actor [%v] is already Jailed for [%v]", a.name, rn)
	}

	if _, ok := a.infractions[rn]; ok {
		inf := a.infractions[rn]
		inf.strikes++
		inf.rebase()
		return a.jail(rn)
	}

	return fmt.Errorf("Infraction against actor [%v]", a.name)
}

func (a *Actor) strikes(rn string) int {
	if _, ok := a.infractions[rn]; ok {
		return a.infractions[rn].strikes
	}
	return 0
}

func (a *Actor) isJailedFor(rn string) bool {
	_, ok := a.jails[rn]
	return ok
}

func (a *Actor) isJailed() bool {
	if len(a.jails) > 0 {
		return true
	}
	return false
}

// shouldDelete returns a bool if the Infractions and Jails maps are empty and the ttl is expired
func (a *Actor) shouldDelete() bool {
	if !a.hasInfractions() && !a.hasJails() {
		if time.Now().After(a.ttl) {
			return true
		}
	}
	return false
}

func (a *Actor) timeServed(j *jail) bool {
	if time.Now().After(j.releaseBy) && j != nil {
		if j.rule.Action != nil {
			ca := a      // copy actor
			cr := j.rule // copy rule
			j.rule.Action.WhenTimeServed(ca, cr)
		}
		delete(a.jails, j.rule.Name)
		a.rebaseAll()
		return true
	}
	return false
}

func (a *Actor) expire(rn string) error {

	// validate key exists
	if _, ok := a.infractions[rn]; !ok {
		return fmt.Errorf("Infraction [%v] does not exist", rn)
	}

	if time.Now().After(a.infractions[rn].expireBy) {
		delete(a.infractions, rn)
		return nil
	}

	return fmt.Errorf("Could not expire [%v]", rn)
}

// jail the actor if the Limit has been reached
func (a *Actor) jail(rn string) error {

	if !a.infractionExists(rn) {
		return fmt.Errorf("jail failed, infraction [%v] does not exist", rn)
	}

	inf := a.infractions[rn]

	if inf.strikes >= inf.rule.StrikeLimit {
		j := newJail(inf.rule, inf.rule.Sentence)
		a.jails[inf.rule.Name] = j
		delete(a.infractions, inf.rule.Name)
		if inf.rule.Action != nil {
			ca := a        // copy actor
			cr := inf.rule // copy rule
			inf.rule.Action.WhenJailed(ca, cr)
		}
		a.rebaseAll()
	}

	return nil
}

func (a *Actor) createInfraction(inf *infraction) error {
	if _, exists := a.infractions[inf.rule.Name]; !exists {
		a.infractions[inf.rule.Name] = inf
		return nil
	}
	return fmt.Errorf("Unable to create infraction [%v]", inf.rule.Name)
}

func (a *Actor) hasInfractions() bool {
	if len(a.infractions) > 0 {
		return true
	}
	return false
}

func (a *Actor) infractionExists(rn string) bool {
	_, ok := a.infractions[rn]
	return ok
}

func (a *Actor) hasJails() bool {
	if len(a.jails) > 0 {
		return true
	}
	return false
}

func (a *Actor) totalJails() int {
	return len(a.jails)
}

func (a *Actor) timeToLive() time.Time {
	return a.ttl
}
