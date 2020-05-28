package badactor

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

// Director is a bucket/shard and contains many Actors
type Director struct {
	sync.Mutex

	index  *list.List
	actors map[string]*list.Element
	rules  map[string]*Rule

	capacity int32
	size     int32
}

// NewDirector instantiates a Director Struct
func NewDirector(ma int32) *Director {
	d := &Director{
		capacity: ma,
		index:    list.New(),
		actors:   make(map[string]*list.Element),
		rules:    make(map[string]*Rule),
	}
	return d
}

func (d *Director) lMaintenance() {
	d.Lock()
	defer d.Unlock()
	for e := d.index.Front(); e != nil; e = e.Next() {
		an := e.Value.(*Actor).name
		d.maintenance(an)
	}
}

func (d *Director) lInfraction(an string, rn string) error {
	d.Lock()
	defer d.Unlock()

	if !d.ruleExists(rn) {
		return fmt.Errorf("Rule [%v] does not exists", rn)
	}

	if !d.actorExists(an) {
		d.createActor(an, rn)
		d.deleteOldest()
	}

	// bail if actor is already jailed
	if d.isJailedFor(an, rn) {
		return fmt.Errorf("Actor [%v] is already jailed for [%v]", an, rn)
	}

	// move the infraction up the stack
	// else create infraction if needed
	if d.infractionExists(an, rn) {
		d.up(an)
	} else {
		d.createInfraction(an, rn)
	}

	return d.incrementInfraction(an, rn)
}

func (d *Director) lCreateInfraction(an string, rn string) error {
	d.Lock()
	defer d.Unlock()

	if !d.actorExists(an) {
		return fmt.Errorf("director.CreateInfraction() failed, Actor does not exists")
	}
	return d.createInfraction(an, rn)
}

func (d *Director) lTimeToLive(an string) (time.Time, error) {
	d.Lock()
	defer d.Unlock()

	if !d.actorExists(an) {
		return time.Now(), fmt.Errorf("director.CreateInfraction() failed, Actor does not exists")
	}
	return d.timeToLive(an), nil
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
	d.Lock()
	defer d.Unlock()
	if !d.actorExists(an) {
		return fmt.Errorf("director.KeepAlive() failed, Actor does not exists")
	}
	d.keepAlive(an)
	return nil
}

func (d *Director) lActorExists(an string) bool {
	d.Lock()
	defer d.Unlock()
	return d.actorExists(an)
}

func (d *Director) lInfractionExists(an string, rn string) bool {
	d.Lock()
	defer d.Unlock()

	if !d.actorExists(an) {
		return false
	}
	return d.infractionExists(an, rn)
}

func (d *Director) lIsJailedFor(an string, rn string) bool {
	d.Lock()
	defer d.Unlock()

	if !d.actorExists(an) {
		return false
	}
	return d.isJailedFor(an, rn)
}

func (d *Director) lIsJailed(an string) bool {
	d.Lock()
	defer d.Unlock()
	if !d.actorExists(an) {
		return false
	}
	return d.isJailed(an)
}

func (d *Director) lStrikes(an string, rn string) (int, error) {
	d.Lock()
	defer d.Unlock()

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

func (d *Director) maintenance(an string) {
	if av := d.actors[an]; av != nil {
		a := av.Value.(*Actor)

		for _, j := range a.jails {
			a.timeServed(j)
		}
		for _, inf := range a.infractions {
			a.expire(inf.rule.Name)
		}
		if a.shouldDelete() {
			delete(d.actors, a.name)
			d.size--
		}
	}
}

func (d *Director) createActor(an string, rn string) error {

	if !d.ruleExists(rn) {
		return fmt.Errorf("createActor failed for Actor [%s], Rule [%s] does not exist", an, rn)
	}

	a := newActor(an, d)
	e := d.index.PushFront(a)
	d.actors[an] = e
	d.size++

	return nil
}

func (d *Director) deleteOldest() {
	for d.isFull() {
		e := d.index.Back()
		a := e.Value.(*Actor)
		d.index.Remove(e)
		delete(d.actors, a.name)
		d.size--
	}
}

func (d *Director) isFull() bool {
	if d.size > d.capacity {
		return true
	}
	return false
}

func (d *Director) ruleExists(rn string) bool {
	_, ok := d.rules[rn]
	return ok
}

func (d *Director) actorExists(an string) bool {
	if _, ok := d.actors[an]; ok {
		d.maintenance(an)
	}
	_, ok := d.actors[an]
	return ok
}

func (d *Director) up(an string) {
	if a := d.actors[an]; a != nil {
		a.Value.(*Actor).accessedAt = time.Now()
		d.index.MoveToFront(a)
		d.deleteOldest()
	}
}

func (d *Director) incrementInfraction(an string, rn string) error {
	return d.actors[an].Value.(*Actor).infraction(rn)
}

func (d *Director) createInfraction(an string, rn string) error {
	inf := newInfraction(d.rules[rn])
	return d.actors[an].Value.(*Actor).createInfraction(inf)
}

func (d *Director) infractionExists(an string, rn string) bool {
	return d.actors[an].Value.(*Actor).infractionExists(rn)
}

func (d *Director) isJailed(an string) bool {
	return d.actors[an].Value.(*Actor).isJailed()
}

func (d *Director) isJailedFor(an string, rn string) bool {
	return d.actors[an].Value.(*Actor).isJailedFor(rn)
}

func (d *Director) strikes(an string, rn string) int {
	return d.actors[an].Value.(*Actor).strikes(rn)
}

func (d *Director) timeToLive(an string) time.Time {
	return d.actors[an].Value.(*Actor).timeToLive()
}

func (d *Director) keepAlive(an string) {
	d.actors[an].Value.(*Actor).rebaseAll()
}

func (d *Director) actor(an string) *Actor {
	return d.actors[an].Value.(*Actor)
}
