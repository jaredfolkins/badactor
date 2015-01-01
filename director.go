package badactor

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Director struct {
	sync.RWMutex
	Actors map[string]*actor
	Rules  map[string]*Rule
	remove chan string
}

func NewDirector() *Director {
	rand.Seed(time.Now().Unix())
	d := &Director{
		Actors: make(map[string]*actor),
		Rules:  make(map[string]*Rule),
		remove: make(chan string),
	}
	return d
}

// spin up all the Rules worker
func (d *Director) Run() {
	go func() {
		for {
			select {
			case an := <-d.remove:
				d.Lock()
				delete(d.Actors, an)
				d.Unlock()
			}
		}
	}()
}

// takes an ActorName and RuleName and creates an Infraction
func (d *Director) CreateInfraction(an string, rn string) error {
	d.Lock()
	defer d.Unlock()

	if !d.actorExists(an) {
		return fmt.Errorf("director.CreateInfraction() failed, Actor does not exists")
	}
	return d.createInfraction(an, rn)
}

func (d *Director) CreateActor(an string, rn string) error {
	d.Lock()
	defer d.Unlock()

	if d.actorExists(an) {
		return fmt.Errorf("director.CreateActor() failed, Actor already exists")
	}
	return d.createActor(an, rn)
}

func (d *Director) KeepAlive(an string) error {
	d.RLock()
	defer d.RUnlock()
	if !d.actorExists(an) {
		return fmt.Errorf("director.KeepAlive() failed, Actor does not exists")
	}
	d.keepAlive(an)
	return nil
}

func (d *Director) ActorExists(an string) bool {
	d.RLock()
	defer d.RUnlock()
	return d.actorExists(an)
}

func (d *Director) InfractionExists(an string, rn string) bool {
	d.RLock()
	defer d.RUnlock()
	if !d.actorExists(an) {
		return false
	}
	return d.infractionExists(an, rn)
}

func (d *Director) IsJailedFor(an string, rn string) bool {
	d.RLock()
	defer d.RUnlock()

	if !d.actorExists(an) {
		return false
	}
	return d.isJailedFor(an, rn)
}

func (d *Director) IsJailed(an string) bool {
	d.RLock()
	defer d.RUnlock()
	if !d.actorExists(an) {
		return false
	}
	return d.isJailed(an)
}

func (d *Director) Strikes(an string, rn string) (int, error) {
	d.RLock()
	defer d.RUnlock()

	if !d.actorExists(an) {
		return 0, fmt.Errorf("director.Strikes() failed, Actor does not exists")
	}

	if !d.infractionExists(an, rn) {
		return 0, fmt.Errorf("director.Strikes() failed, Infraction does not exists")
	}

	return d.strikes(an, rn), nil
}

// Infraction does the most
func (d *Director) Infraction(an string, rn string) error {

	var res bool
	var err error

	if d.IsJailedFor(an, rn) {
		return fmt.Errorf("Actor [%v] is already jailed for [%v]", an, rn)
	}

	res, err = d.MostCostlyInfraction(an, rn)
	if res {
		return err
	}

	return fmt.Errorf("director.Infraction() failed for [ActorName:%v, RuleName:%v]", an, rn)
}

func (d *Director) MostCostlyInfraction(an string, rn string) (bool, error) {
	d.Lock()
	defer d.Unlock()

	if !d.actorExists(an) {
		err := d.createActor(an, rn)
		if err != nil {
			return false, err
		}
	}

	// create infraction if needed
	if !d.infractionExists(an, rn) {
		d.createInfraction(an, rn)
	}

	return true, d.incrementInfraction(an, rn)
}

// returns the total strikes of an Actor
func (d *Director) AddRule(r *Rule) error {
	d.Lock()
	defer d.Unlock()

	if d.ruleExists(r.Name) {
		return fmt.Errorf("AddRule failed, Rule [%s] already exists", r.Name)
	}

	// add the rule
	d.Rules[r.Name] = r

	return nil
}

//
// NO LOCKS - UNSAFE - BECAREFUL
//
// below this point are helper functions that are dependant on the calling functions to preform the appropriate locking

func (d *Director) createActor(an string, rn string) error {
	if _, ok := d.Rules[rn]; !ok {
		return fmt.Errorf("createActor failed for Actor [%s], Rule [%s] does not exist", an, rn)
	}
	a := newActor(an, d)
	a.run()
	d.Actors[an] = a
	return nil
}

func (d *Director) ruleExists(rn string) bool {
	_, ok := d.Rules[rn]
	return ok
}

func (d *Director) actorExists(an string) bool {
	_, ok := d.Actors[an]
	return ok
}

func (d *Director) incrementInfraction(an string, rn string) error {
	return d.Actors[an].lInfraction(rn)
}

// takes an ActorName and RuleName and creates an Infraction
func (d *Director) createInfraction(an string, rn string) error {
	inf := NewInfraction(d.Rules[rn])
	return d.Actors[an].lCreateInfraction(inf)
}

func (d *Director) infractionExists(an string, rn string) bool {
	return d.Actors[an].lInfractionExists(rn)
}

func (d *Director) isJailed(an string) bool {
	return d.Actors[an].lIsJailed()
}

func (d *Director) isJailedFor(an string, rn string) bool {
	return d.Actors[an].lIsJailedFor(rn)
}

func (d *Director) strikes(an string, rn string) int {
	return d.Actors[an].lStrikes(rn)
}

func (d *Director) keepAlive(an string) {
	d.Actors[an].lRebaseAll()
}
