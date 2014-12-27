package badactor

import (
	"fmt"
	"strconv"
	"sync"
)

const (
	QUIT              = 1
	CREATE_ACTOR      = 2
	INFRACTION        = 3
	STRIKES           = 4
	IS_JAILED         = 5
	IS_JAILED_FOR     = 6
	REMOVE_ACTOR      = 7
	KEEP_ALIVE        = 8
	CREATE_INFRACTION = 9
	INFRACTION_EXISTS = 10
)

type Director struct {
	Actors map[string]*Actor
	Rules  map[string]*Rule
	remove chan *Incoming
	rwmu   sync.RWMutex
}

func NewDirector() *Director {
	d := &Director{
		Actors: make(map[string]*Actor),
		Rules:  make(map[string]*Rule),
		remove: make(chan *Incoming),
	}
	return d
}

// spin up all the Rules worker
func (d *Director) Run() {
	go func() {
		for {
			select {
			case in := <-d.remove:
				d.rwmu.Lock()
				delete(d.Actors, in.ActorName)
				d.rwmu.Unlock()
			}
		}
	}()
}

func (d *Director) MostCostlyInfraction(an string, rn string) (bool, error) {
	d.rwmu.Lock()
	defer d.rwmu.Unlock()
	var res bool

	res = d.actorExists(an)
	if res == false {
		err := d.createActor(an, rn)
		if err != nil {
			return false, err
		}
	}

	// create infraction if needed

	res = d.infractionExists(an, rn)
	if res == false {
		// i can't detect a way for this to fail with
		// all the previous checks in place
		d.createInfraction(an, rn)
		/*
			if err != nil {
				return false, err
			}
		*/
	}

	// i can't detect a way for this to fail with
	// all the previous checks in place
	d.incrementInfraction(an, rn)
	/*
		err = d.incrementInfraction(an, rn)
		if err != nil {
			return false, err
		}
	*/
	return true, nil
}

func (d *Director) LeastCostlyInfraction(an string, rn string) (bool, error) {
	d.rwmu.RLock()
	defer d.rwmu.RUnlock()
	if !d.actorExists(an) {
		return false, fmt.Errorf("director.leastCostlyInfraction() actorExists() failed [%v:%v]", an, rn)
	}

	if !d.infractionExists(an, rn) {
		return false, fmt.Errorf("director.leastCostlyInfraction() infractionExists() failed [%v:%v]", an, rn)
	}

	return true, d.incrementInfraction(an, rn)
}

// takes an ActorName and RuleName and creates an Infraction
func (d *Director) CreateInfraction(an string, rn string) error {
	d.rwmu.Lock()
	defer d.rwmu.Unlock()
	if !d.actorExists(an) {
		return fmt.Errorf("director.CreateInfraction() failed, Actor does not exists")
	}
	return d.createInfraction(an, rn)
}

func (d *Director) CreateActor(an string, rn string) error {
	d.rwmu.Lock()
	defer d.rwmu.Unlock()
	if d.actorExists(an) {
		return fmt.Errorf("director.CreateActor() failed, Actor already exists")
	}
	return d.createActor(an, rn)
}

func (d *Director) KeepAlive(an string) error {
	d.rwmu.RLock()
	defer d.rwmu.RUnlock()
	if !d.actorExists(an) {
		return fmt.Errorf("director.KeepAlive() failed, Actor does not exists")
	}
	return d.keepAlive(an)
}

func (d *Director) ActorExists(an string) bool {
	d.rwmu.RLock()
	defer d.rwmu.RUnlock()
	return d.actorExists(an)
}

func (d *Director) InfractionExists(an string, rn string) bool {
	d.rwmu.RLock()
	defer d.rwmu.RUnlock()
	if !d.actorExists(an) {
		return false
	}
	return d.infractionExists(an, rn)
}

func (d *Director) IsJailedFor(an string, rn string) bool {
	d.rwmu.RLock()
	defer d.rwmu.RUnlock()
	if !d.actorExists(an) {
		return false
	}
	return d.isJailedFor(an, rn)
}

func (d *Director) IsJailed(an string) bool {
	d.rwmu.RLock()
	defer d.rwmu.RUnlock()
	if !d.actorExists(an) {
		return false
	}
	return d.isJailed(an)
}

func (d *Director) Strikes(an string, rn string) (int, error) {
	d.rwmu.RLock()
	defer d.rwmu.RUnlock()

	if !d.actorExists(an) {
		return 0, fmt.Errorf("director.Strikes() failed, Actor does not exists")
	}

	if !d.infractionExists(an, rn) {
		return 0, fmt.Errorf("director.Strikes() failed, Infraction does not exists")
	}

	return d.strikes(an, rn)
}

// Infraction does the most
func (d *Director) Infraction(an string, rn string) error {

	var res bool
	var err error

	//try
	res, err = d.LeastCostlyInfraction(an, rn)
	if res {
		return err
	}

	// then try
	res, err = d.MostCostlyInfraction(an, rn)
	if res {
		return err
	}

	return fmt.Errorf("director.Infraction() failed for [ActorName:%v, RuleName:%v]", an, rn)
}

// returns the total strikes of an Actor
func (d *Director) AddRule(r *Rule) error {
	d.rwmu.Lock()
	defer d.rwmu.Unlock()

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
	a := NewActor(an, d)
	a.Run()
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
	in := NewIncoming(an, rn, INFRACTION)
	d.Actors[an].Incoming <- in
	out := <-in.Outgoing
	close(in.Outgoing)
	return out.Error
}

// takes an ActorName and RuleName and creates an Infraction
func (d *Director) createInfraction(an string, rn string) error {
	in := NewIncoming(an, rn, CREATE_INFRACTION)
	inf := NewInfraction(d.Rules[rn])
	in.Infraction = inf
	d.Actors[an].Incoming <- in
	out := <-in.Outgoing
	close(in.Outgoing)
	return out.Error
}

func (d *Director) infractionExists(an string, rn string) bool {
	in := NewIncoming(an, rn, INFRACTION_EXISTS)
	d.Actors[an].Incoming <- in
	out := <-in.Outgoing
	res, _ := strconv.ParseBool(out.Message)
	return res
}

func (d *Director) isJailed(an string) bool {
	in := NewIncoming(an, "", IS_JAILED)
	d.Actors[an].Incoming <- in
	out := <-in.Outgoing
	res, _ := strconv.ParseBool(out.Message)
	return res
}

func (d *Director) isJailedFor(an string, rn string) bool {
	in := NewIncoming(an, rn, IS_JAILED_FOR)
	d.Actors[an].Incoming <- in
	out := <-in.Outgoing
	res, _ := strconv.ParseBool(out.Message)
	return res
}

func (d *Director) strikes(an string, rn string) (int, error) {
	in := NewIncoming(an, rn, STRIKES)
	d.Actors[an].Incoming <- in
	out := <-in.Outgoing
	close(in.Outgoing)
	i, _ := strconv.Atoi(out.Message)
	return i, nil
}

func (d *Director) keepAlive(an string) error {
	in := NewIncoming(an, "", KEEP_ALIVE)
	d.Actors[an].Incoming <- in
	out := <-in.Outgoing
	close(in.Outgoing)
	return out.Error
}
