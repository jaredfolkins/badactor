package badactor

import (
	"fmt"
	"sync"
	"time"
)

const (
	TTL          = 100
	REMOVE_ACTOR = 007
)

type Actor struct {
	sync.RWMutex
	Name        string
	Infractions map[string]*Infraction
	Jails       map[string]*Sentence
	director    *Director
	ttl         time.Time
}

func NewActor(n string, d *Director) *Actor {
	a := &Actor{
		director:    d,
		ttl:         time.Now().Add(time.Millisecond * TTL),
		Name:        n,
		Infractions: make(map[string]*Infraction),
		Jails:       make(map[string]*Sentence),
	}
	return a
}

func NewClassicActor(n string, r *Rule, d *Director) *Actor {
	a := &Actor{
		director:    d,
		ttl:         time.Now().Add(time.Millisecond * TTL),
		Name:        n,
		Infractions: make(map[string]*Infraction),
		Jails:       make(map[string]*Sentence),
	}

	a.Lock()
	a.Infractions[r.Name] = NewInfraction(r)
	a.Unlock()
	return a
}

func (a *Actor) Run() {
	ticker := time.NewTicker(time.Second * 20)
	go func() {
		// 1.4 means i refractor this
		for _ = range ticker.C {
			if a.Maintenance() {
				ticker.Stop()
				return
			}
		}
	}()
}

func (a *Actor) HasJails() bool {
	a.RLock()
	res := a.hasJails()
	a.RUnlock()
	return res
}

func (a *Actor) ShouldReturn() bool {
	a.Lock()
	res := a.shouldReturn()
	a.Unlock()
	return res
}

func (a *Actor) InfractionExists(rn string) bool {
	a.RLock()
	res := a.infractionExists(rn)
	a.RUnlock()
	return res
}

// Lockup the Actor if the Limit has been reached
func (a *Actor) Lockup(rn string) error {
	a.Lock()
	err := a.lockup(rn)
	a.Unlock()
	return err
}

func (a *Actor) HasInfractions() bool {
	a.RLock()
	res := a.hasInfractions()
	a.RUnlock()
	return res
}

func (a *Actor) Expire(rn string) error {
	a.Lock()
	err := a.expire(rn)
	a.Unlock()
	return err
}

func (a *Actor) Infraction(rn string) error {
	a.Lock()
	err := a.infraction(rn)
	a.Unlock()
	return err
}

func (a *Actor) CreateInfraction(inf *Infraction) error {
	a.Lock()
	err := a.createInfraction(inf)
	a.Unlock()
	return err
}

func (a *Actor) IsJailed() bool {
	a.RLock()
	res := a.isJailed()
	a.RUnlock()
	return res
}

func (a *Actor) IsJailedFor(rn string) bool {
	a.RLock()
	res := a.isJailedFor(rn)
	a.RUnlock()
	return res
}

func (a *Actor) Strikes(rn string) int {
	a.RLock()
	res := a.strikes(rn)
	a.RUnlock()
	return res
}

func (a *Actor) RebaseAll() error {
	a.Lock()
	res := a.rebaseAll()
	a.Unlock()
	return res
}

//
// NO LOCKS - UNSAFE - BECAREFUL
//
// below this point are helper functions that are dependant on the calling functions to preform the appropriate locking

func (a *Actor) rebaseAll() error {

	for _, inf := range a.Infractions {
		inf.Rebase()
	}

	a.ttl = time.Now().Add(time.Second * TTL)

	return nil
}

func (a *Actor) infraction(rn string) error {

	if a.isJailed() {
		return fmt.Errorf("Actor [%v] is already Jailed for [%v]", a.Name, rn)
	}

	if _, ok := a.Infractions[rn]; ok {
		inf := a.Infractions[rn]
		inf.Strikes++
		inf.Rebase()
		return a.lockup(rn)
	}

	return fmt.Errorf("Infraction against Actor [%v]", a.Name)
}

func (a *Actor) strikes(rn string) int {
	if _, ok := a.Infractions[rn]; ok {
		return a.Infractions[rn].Strikes
	}
	return 0
}

func (a *Actor) isJailedFor(rn string) bool {
	_, ok := a.Jails[rn]
	return ok
}

func (a *Actor) isJailed() bool {
	if len(a.Jails) > 0 {
		return true
	}
	return false
}

func (a *Actor) Maintenance() bool {
	a.Lock()
	res := a.maintenance()
	a.Unlock()
	return res
}

// maintenance does some background tasks
// Locksup, Expires, or Releases any Actors
func (a *Actor) maintenance() bool {
	for _, s := range a.Jails {
		a.timeServed(s)
	}

	for _, inf := range a.Infractions {
		a.lockup(inf.Rule.Name)
		a.expire(inf.Rule.Name)
	}

	return a.shouldReturn()
}

func (a *Actor) shouldReturn() bool {

	// if the Infractions OR Jails maps are not empty, we can return
	// as we are certain that we do not want the Actor to quit
	if a.hasInfractions() || a.hasJails() {
		return false
	}

	// the ttl is a time buffer, it allows a newly created actor
	// a few milliseconds to get its State setup
	// avoiding a potential errouneous quit
	if time.Now().After(a.ttl) {
		a.director.remove <- a.Name
		return true
	}

	return false
}

func (a *Actor) timeServed(s *Sentence) bool {

	if time.Now().After(s.ReleaseBy) && s != nil {
		delete(a.Jails, s.Rule.Name)
		a.rebaseAll()
		return true
	}
	return false
}

func (a *Actor) expire(rn string) error {

	// validate key exists
	if _, ok := a.Infractions[rn]; !ok {
		return fmt.Errorf("Infraction [%v] does not exist", rn)
	}

	if time.Now().After(a.Infractions[rn].ExpireBy) {
		delete(a.Infractions, rn)
		return nil
	}

	return nil
}

// Lockup the Actor if the Limit has been reached
func (a *Actor) lockup(rn string) error {

	if !a.infractionExists(rn) {
		return fmt.Errorf("Lockup failed, infraction [%v] does not exist", rn)
	}

	inf := a.Infractions[rn]

	if inf.Strikes >= inf.Rule.StrikeLimit {
		sen := NewSentence(inf.Rule, inf.Rule.Sentence)
		a.Jails[inf.Rule.Name] = sen
		delete(a.Infractions, inf.Rule.Name)
		a.rebaseAll()
	}

	return nil
}

func (a *Actor) createInfraction(inf *Infraction) error {
	if _, exists := a.Infractions[inf.Rule.Name]; !exists {
		a.Infractions[inf.Rule.Name] = inf
		return nil
	}
	return fmt.Errorf("Unable to create infraction [%v]", inf.Rule.Name)
}

func (a *Actor) hasInfractions() bool {
	if len(a.Infractions) > 0 {
		return true
	}
	return false
}

func (a *Actor) infractionExists(rn string) bool {
	_, ok := a.Infractions[rn]
	return ok
}

func (a *Actor) hasJails() bool {
	if len(a.Jails) > 0 {
		return true
	}
	return false
}

func (a *Actor) totalJails() int {
	return len(a.Jails)
}
