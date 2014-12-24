package badactor

import (
	"fmt"
	"log"
	"strconv"
	"sync"
)

const (
	QUIT              = "quit"
	CREATE_ACTOR      = "create_actor"
	INFRACTION        = "infraction"
	STRIKES           = "strikes"
	IS_JAILED         = "is_jailed"
	IS_JAILED_FOR     = "is_jailed_for"
	REMOVE_ACTOR      = "remove_actor"
	KEEP_ALIVE        = "keep_alive"
	CREATE_INFRACTION = "create_infraction"
	INFRACTION_EXISTS = "infraction_exists"
)

type Director struct {
	Actors    map[string]*Actor
	Rules     map[string]*Rule
	delete_me chan *Incoming
	rwmu      sync.RWMutex
}

func NewDirector() *Director {
	d := &Director{
		Actors:    make(map[string]*Actor),
		Rules:     make(map[string]*Rule),
		delete_me: make(chan *Incoming),
	}
	return d
}

// spin up all the Rules worker
func (d *Director) Run() {
	go func() {
		for {
			select {
			case in := <-d.delete_me:
				d.rwmu.Lock()
				log.Printf("director.Actors[%v] is deleted\n", in.ActorName)
				delete(d.Actors, in.ActorName)
				d.rwmu.Unlock()
			}
		}
	}()
}

/*
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
*/

// Infraction does the most
func (d *Director) Infraction(an string, rn string) error {

	var res bool
	var err error

	log.Println("leastCostly: start")
	d.rwmu.RLock()
	res, err = d.leastCostlyInfraction(an, rn)
	d.rwmu.RUnlock()
	if res {
		log.Println("leastCostly: end")
		return err
	}

	log.Println("mostCostly: start")
	d.rwmu.Lock()
	res, err = d.mostCostlyInfraction(an, rn)
	d.rwmu.Unlock()
	if res {
		log.Println("mostCostly: end")
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

	if _, ok := d.Actors[an]; !ok {
		return fmt.Errorf("createActor failed for Actor [%s] and Rule [%s]", an, rn)
	}

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
	log.Println("incrementInfraction")
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
	log.Println("infractionExists")
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

func (d *Director) mostCostlyInfraction(an string, rn string) (bool, error) {
	var res bool
	var err error

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
		err := d.createInfraction(an, rn)
		if err != nil {
			return false, err
		}
	}

	err = d.incrementInfraction(an, rn)
	return true, err
}

func (d *Director) leastCostlyInfraction(an string, rn string) (bool, error) {
	if !d.actorExists(an) {
		return false, fmt.Errorf("director.leastCostlyInfraction() actorExists() failed [%v:%v]", an, rn)
	}

	if !d.infractionExists(an, rn) {
		return false, fmt.Errorf("director.leastCostlyInfraction() infractionExists() failed [%v:%v]", an, rn)
	}

	return true, d.incrementInfraction(an, rn)
}
