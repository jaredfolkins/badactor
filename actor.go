package badactor

import (
	"fmt"
	"strconv"
	"time"
)

const TTL = 100

type Actor struct {
	Name        string
	Infractions map[string]*Infraction
	Jails       map[string]*Sentence
	director    *Director
	Incoming    chan *Incoming
	ttl         time.Time
}

func NewActor(n string, d *Director) *Actor {
	a := &Actor{
		Name:        n,
		Infractions: make(map[string]*Infraction),
		Jails:       make(map[string]*Sentence),
		Incoming:    make(chan *Incoming),
		director:    d,
		ttl:         time.Now().Add(time.Millisecond * TTL),
	}
	return a
}

func NewClassicActor(n string, r *Rule, d *Director) *Actor {
	a := &Actor{
		Name:        n,
		Infractions: make(map[string]*Infraction),
		Jails:       make(map[string]*Sentence),
		Incoming:    make(chan *Incoming),
		director:    d,
		ttl:         time.Now().Add(time.Second * TTL),
	}

	a.Infractions[r.Name] = NewInfraction(r)
	return a
}

func (a *Actor) Run() {
	go func() {
		for {
			select {
			case in := <-a.Incoming:
				a.switchBoard(in)
			default:
				if a.maintenance() {
					return
				}
			}
		}
	}()
}

func (a *Actor) rebaseAll() error {
	for _, inf := range a.Infractions {
		inf.Rebase()
	}
	a.ttl = time.Now().Add(time.Second * TTL)
	return nil
}

func (a *Actor) switchBoard(in *Incoming) {
	switch in.Type {
	case KEEP_ALIVE:
		err := a.rebaseAll()
		in.Outgoing <- &Outgoing{Error: err}
	case INFRACTION:
		err := a.infraction(in.RuleName)
		in.Outgoing <- &Outgoing{Error: err}
	case CREATE_INFRACTION:
		err := a.createInfraction(in.Infraction)
		in.Outgoing <- &Outgoing{Error: err}
	case STRIKES:
		total := a.strikes(in.RuleName)
		in.Outgoing <- &Outgoing{Message: strconv.Itoa(total)}
	case IS_JAILED:
		res := a.isJailed()
		in.Outgoing <- &Outgoing{Message: strconv.FormatBool(res), Error: nil}
	case IS_JAILED_FOR:
		res := a.isJailedFor(in.RuleName)
		in.Outgoing <- &Outgoing{Message: strconv.FormatBool(res), Error: nil}
	case INFRACTION_EXISTS:
		res := a.infractionExists(in.RuleName)
		in.Outgoing <- &Outgoing{Message: strconv.FormatBool(res), Error: nil}
	}
}

func (a *Actor) infraction(rn string) error {

	if a.isJailed() {
		return fmt.Errorf("Actor [%v] is already Jailed for [%v]", a.Name, rn)
	}

	if _, ok := a.Infractions[rn]; ok {
		inf := a.Infractions[rn]
		inf.Strikes = inf.Strikes + 1
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

// maintenance does some background tasks
// Locksup, Expires, or Releases any Actors
func (a *Actor) maintenance() bool {
	for _, s := range a.Jails {
		a.timeServed(s)
	}

	for _, inf := range a.Infractions {
		a.expire(inf.Rule.Name)
		a.lockup(inf.Rule.Name)
	}

	return a.shouldReturn()
}

func (a *Actor) shouldReturn() bool {
	// if the Infractions OR Jails maps are not empty, we can return
	// certain, that we do not want the Actor to quit
	if a.hasInfractions() || a.hasJails() {
		return false
	}

	// the ttl is a time buffer, it allows a newly created actor
	// a few milliseconds to get its State setup
	// avoiding a potential errouneous quit
	if time.Now().After(a.ttl) {
		in := &Incoming{
			ActorName: a.Name,
			Type:      REMOVE_ACTOR,
		}
		a.director.delete_me <- in
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
