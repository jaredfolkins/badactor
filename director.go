package badactor

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Director is a singleton, only one should be running
// it is engine of BadActor and exposes the primary set of Public Functions
type Director struct {
	sync.RWMutex
	maxActors int32
	actors    map[string]*actor
	rules     map[string]*Rule
}

// NewDirector instantiates a Director Struct
func NewDirector(ma int32) *Director {
	rand.Seed(time.Now().Unix())
	d := &Director{
		maxActors: ma,
		actors:    make(map[string]*actor),
		rules:     make(map[string]*Rule),
	}
	return d
}

func (d *Director) maintenance() {
	d.lReadMaintenance()
	d.lWriteMaintenance()
}

func (d *Director) lInfraction(an string, rn string) error {
	d.Lock()
	defer d.Unlock()

	if !d.actorExists(an) {
		err := d.createActor(an, rn)
		if err != nil {
			return err
		}
	}

	if d.isJailedFor(an, rn) {
		return fmt.Errorf("Actor [%v] is already jailed for [%v]", an, rn)
	}

	// create infraction if needed
	if !d.infractionExists(an, rn) {
		d.createInfraction(an, rn)
	}

	return d.incrementInfraction(an, rn)
}

func (d *Director) lWriteMaintenance() {
	d.Lock()
	defer d.Unlock()

	for _, a := range d.actors {
		if a.lShouldDelete() {
			delete(d.actors, a.name)
		}
	}
}

func (d *Director) lReadMaintenance() {
	d.RLock()
	defer d.RUnlock()

	for _, a := range d.actors {
		a.Lock()

		for _, s := range a.jails {
			a.timeServed(s)
		}

		for _, inf := range a.infractions {
			a.lockup(inf.rule.Name)
			a.expire(inf.rule.Name)
		}
		a.Unlock()
	}

}

func (d *Director) lCreateInfraction(an string, rn string) error {
	d.Lock()
	defer d.Unlock()

	if !d.actorExists(an) {
		return fmt.Errorf("director.CreateInfraction() failed, Actor does not exists")
	}
	return d.createInfraction(an, rn)
}

func (d *Director) lCreateActor(an string, rn string) error {
	d.Lock()
	defer d.Unlock()

	if d.actorExists(an) {
		return fmt.Errorf("director.CreateActor() failed, Actor already exists")
	}
	return d.createActor(an, rn)
}

func (d *Director) lKeepAlive(an string) error {
	d.RLock()
	defer d.RUnlock()
	if !d.actorExists(an) {
		return fmt.Errorf("director.KeepAlive() failed, Actor does not exists")
	}
	d.keepAlive(an)
	return nil
}

func (d *Director) lActorExists(an string) bool {
	d.RLock()
	defer d.RUnlock()
	return d.actorExists(an)
}

func (d *Director) lInfractionExists(an string, rn string) bool {
	d.RLock()
	defer d.RUnlock()
	if !d.actorExists(an) {
		return false
	}
	return d.infractionExists(an, rn)
}

func (d *Director) lIsJailedFor(an string, rn string) bool {
	d.RLock()
	defer d.RUnlock()

	if !d.actorExists(an) {
		return false
	}
	return d.isJailedFor(an, rn)
}

func (d *Director) lIsJailed(an string) bool {
	d.RLock()
	defer d.RUnlock()
	if !d.actorExists(an) {
		return false
	}
	return d.isJailed(an)
}

func (d *Director) lStrikes(an string, rn string) (int, error) {
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

func (d *Director) lAddRule(r *Rule) error {
	d.Lock()
	defer d.Unlock()

	if d.ruleExists(r.Name) {
		return fmt.Errorf("AddRule failed, Rule [%s] already exists", r.Name)
	}

	// add the rule
	d.rules[r.Name] = r

	return nil
}

//
// NO LOCKS - UNSAFE - BECAREFUL
//
// below this point are helper functions that are dependant on the calling functions to preform the appropriate locking

func (d *Director) createActor(an string, rn string) error {
	if _, ok := d.rules[rn]; !ok {
		return fmt.Errorf("createActor failed for Actor [%s], Rule [%s] does not exist", an, rn)
	}
	a := newActor(an, d)
	d.actors[an] = a
	return nil
}

func (d *Director) ruleExists(rn string) bool {
	_, ok := d.rules[rn]
	return ok
}

func (d *Director) actorExists(an string) bool {
	_, ok := d.actors[an]
	return ok
}

func (d *Director) incrementInfraction(an string, rn string) error {
	return d.actors[an].lInfraction(rn)
}

func (d *Director) createInfraction(an string, rn string) error {
	inf := newInfraction(d.rules[rn])
	return d.actors[an].lCreateInfraction(inf)
}

func (d *Director) infractionExists(an string, rn string) bool {
	return d.actors[an].lInfractionExists(rn)
}

func (d *Director) isJailed(an string) bool {
	return d.actors[an].lIsJailed()
}

func (d *Director) isJailedFor(an string, rn string) bool {
	return d.actors[an].lIsJailedFor(rn)
}

func (d *Director) strikes(an string, rn string) int {
	return d.actors[an].lStrikes(rn)
}

func (d *Director) keepAlive(an string) {
	d.actors[an].lRebaseAll()
}
