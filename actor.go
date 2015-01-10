package badactor

import (
	"fmt"
	"time"
)

// ttl is the time to live value for newly created actors
const ttl = 100

type actor struct {
	name        string
	infractions map[string]*infraction
	jails       map[string]*jail
	director    *Director
	ttl         time.Time
	accessedAt  time.Time
}

func newActor(n string, d *Director) *actor {
	a := &actor{
		director:    d,
		ttl:         time.Now().Add(time.Millisecond * ttl),
		name:        n,
		infractions: make(map[string]*infraction),
		jails:       make(map[string]*jail),
		accessedAt:  time.Now(),
	}
	return a
}

func newClassicActor(n string, r *Rule, d *Director) *actor {
	a := &actor{
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

//
// NO LOCKS - UNSAFE - BECAREFUL
//
// below this point are helper functions that are dependant on the calling functions to preform the appropriate locking

func (a *actor) rebaseAll() error {

	for _, inf := range a.infractions {
		inf.rebase()
	}

	a.ttl = time.Now().Add(time.Second * ttl)

	return nil
}

func (a *actor) infraction(rn string) error {

	if a.isJailed() {
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

func (a *actor) strikes(rn string) int {
	if _, ok := a.infractions[rn]; ok {
		return a.infractions[rn].strikes
	}
	return 0
}

func (a *actor) isJailedFor(rn string) bool {
	_, ok := a.jails[rn]
	return ok
}

func (a *actor) isJailed() bool {
	if len(a.jails) > 0 {
		return true
	}
	return false
}

// shouldDelete returns a bool if the Infractions and Jails maps are empty and the ttl is expired
func (a *actor) shouldDelete() bool {
	if !a.hasInfractions() && !a.hasJails() {
		if time.Now().After(a.ttl) {
			return true
		}
	}
	return false
}

func (a *actor) timeServed(j *jail) bool {
	if time.Now().After(j.releaseBy) && j != nil {
		s := NewStats(a, j.rule)
		if j.rule.Action != nil {
			j.rule.Action.WhenTimeServed(s)
		}
		delete(a.jails, j.rule.Name)
		a.rebaseAll()
		return true
	}
	return false
}

func (a *actor) expire(rn string) error {

	// validate key exists
	if _, ok := a.infractions[rn]; !ok {
		return fmt.Errorf("Infraction [%v] does not exist", rn)
	}

	if time.Now().After(a.infractions[rn].expireBy) {
		delete(a.infractions, rn)
		return nil
	}

	return nil
}

// jail the actor if the Limit has been reached
func (a *actor) jail(rn string) error {

	if !a.infractionExists(rn) {
		return fmt.Errorf("jail failed, infraction [%v] does not exist", rn)
	}

	inf := a.infractions[rn]

	if inf.strikes >= inf.rule.StrikeLimit {
		j := newJail(inf.rule, inf.rule.Sentence)
		a.jails[inf.rule.Name] = j
		delete(a.infractions, inf.rule.Name)
		if inf.rule.Action != nil {
			s := NewStats(a, inf.rule)
			inf.rule.Action.WhenJailed(s)
		}
		a.rebaseAll()
	}

	return nil
}

func (a *actor) createInfraction(inf *infraction) error {
	if _, exists := a.infractions[inf.rule.Name]; !exists {
		a.infractions[inf.rule.Name] = inf
		return nil
	}
	return fmt.Errorf("Unable to create infraction [%v]", inf.rule.Name)
}

func (a *actor) hasInfractions() bool {
	if len(a.infractions) > 0 {
		return true
	}
	return false
}

func (a *actor) infractionExists(rn string) bool {
	_, ok := a.infractions[rn]
	return ok
}

func (a *actor) hasJails() bool {
	if len(a.jails) > 0 {
		return true
	}
	return false
}

func (a *actor) totalJails() int {
	return len(a.jails)
}
