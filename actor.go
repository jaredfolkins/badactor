package badactor

import (
	"fmt"
	"sync"
	"time"
)

// ttl is the time to live value for newly created actors
const ttl = 100

type actor struct {
	sync.RWMutex
	name        string
	infractions map[string]*infraction
	jails       map[string]*sentence
	director    *Director
	ttl         time.Time
}

func newActor(n string, d *Director) *actor {
	a := &actor{
		director:    d,
		ttl:         time.Now().Add(time.Millisecond * ttl),
		name:        n,
		infractions: make(map[string]*infraction),
		jails:       make(map[string]*sentence),
	}
	return a
}

func newClassicActor(n string, r *Rule, d *Director) *actor {
	a := &actor{
		director:    d,
		ttl:         time.Now().Add(time.Millisecond * ttl),
		name:        n,
		infractions: make(map[string]*infraction),
		jails:       make(map[string]*sentence),
	}

	a.Lock()
	a.infractions[r.Name] = newInfraction(r)
	a.Unlock()
	return a
}

func (a *actor) lHasJails() bool {
	a.RLock()
	res := a.hasJails()
	a.RUnlock()
	return res
}

func (a *actor) lShouldDelete() bool {
	a.Lock()
	res := a.shouldDelete()
	a.Unlock()
	return res
}

func (a *actor) lInfractionExists(rn string) bool {
	a.RLock()
	res := a.infractionExists(rn)
	a.RUnlock()
	return res
}

// Lockup the actor if the Limit has been reached
func (a *actor) lLockup(rn string) error {
	a.Lock()
	err := a.lockup(rn)
	a.Unlock()
	return err
}

func (a *actor) lHasInfractions() bool {
	a.RLock()
	res := a.hasInfractions()
	a.RUnlock()
	return res
}

func (a *actor) lExpire(rn string) error {
	a.Lock()
	err := a.expire(rn)
	a.Unlock()
	return err
}

func (a *actor) lInfraction(rn string) error {
	a.Lock()
	err := a.infraction(rn)
	a.Unlock()
	return err
}

func (a *actor) lCreateInfraction(inf *infraction) error {
	a.Lock()
	err := a.createInfraction(inf)
	a.Unlock()
	return err
}

func (a *actor) lIsJailed() bool {
	a.RLock()
	res := a.isJailed()
	a.RUnlock()
	return res
}

func (a *actor) lIsJailedFor(rn string) bool {
	a.RLock()
	res := a.isJailedFor(rn)
	a.RUnlock()
	return res
}

func (a *actor) lStrikes(rn string) int {
	a.RLock()
	res := a.strikes(rn)
	a.RUnlock()
	return res
}

func (a *actor) lRebaseAll() error {
	a.Lock()
	res := a.rebaseAll()
	a.Unlock()
	return res
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
		return a.lockup(rn)
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

func (a *actor) timeServed(s *sentence) bool {
	if time.Now().After(s.releaseBy) && s != nil {
		delete(a.jails, s.rule.Name)
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

// Lockup the actor if the Limit has been reached
func (a *actor) lockup(rn string) error {

	if !a.infractionExists(rn) {
		return fmt.Errorf("Lockup failed, infraction [%v] does not exist", rn)
	}

	inf := a.infractions[rn]

	if inf.strikes >= inf.rule.StrikeLimit {
		sen := newSentence(inf.rule, inf.rule.Sentence)
		a.jails[inf.rule.Name] = sen
		delete(a.infractions, inf.rule.Name)
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
