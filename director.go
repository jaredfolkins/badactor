package badactor

import (
	"fmt"
	"strconv"
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
)

type Director struct {
	Actors    map[string]*Actor
	Rules     map[string]*Rule
	delete_me chan *Incoming
}

func NewDirector() *Director {
	d := &Director{
		Actors:    make(map[string]*Actor),
		Rules:     make(map[string]*Rule),
		delete_me: make(chan *Incoming),
	}
	//d.Run()
	return d
}

// the classic Director has a predefined rule
// for easier comprehension and initial setup
func NewClassicDirector() *Director {
	// create director
	d := NewDirector()

	// create default login rule
	r := NewClassicRule("Login", "You exceeded the allowed login attempts")

	// assign default rule
	d.Rules[r.Name] = r
	//d.Run()
	return d
}

// spin up all the Rules worker
func (d *Director) Run() {
	go func() {
		for {
			select {
			case in := <-d.delete_me:
				if _, ok := d.Actors[in.ActorName]; ok {
					delete(d.Actors, in.ActorName)
					in.Outgoing <- &Outgoing{"success", nil}
				} else {
					err := fmt.Errorf("[%v] already deleted", in.ActorName)
					in.Outgoing <- &Outgoing{"fail", err}
				}
			}
		}
	}()
}

func (d *Director) switchBoard(in *Incoming) {
	switch in.Type {
	case REMOVE_ACTOR:
		if !d.ActorExists(in.ActorName) {
			delete(d.Actors, in.ActorName)
			out := &Outgoing{"Success", nil}
			in.Outgoing <- out
		}
	}
}

func (d *Director) KeepAlive(an string) error {
	// bail if no actor
	if !d.ActorExists(an) {
		return fmt.Errorf("KeepAlive failed, Actor does not exists")
	}

	in := NewIncoming(an, "", KEEP_ALIVE)
	d.Actors[an].Incoming <- in
	out := <-in.Outgoing
	close(in.Outgoing)
	return out.Error
}

func (d *Director) Infraction(an string, rn string) error {

	// create actor if needed
	if !d.ActorExists(an) {
		err := d.CreateActor(an, rn)
		if err != nil {
			return err
		}
	}

	// create infraction if needed
	if !d.InfractionExists(an, rn) {
		err := d.CreateInfraction(an, rn)
		if err != nil {
			return err
		}
	}

	// increment the infraction and return
	return d.IncrementInfraction(an, rn)
}

func (d *Director) AddRule(r *Rule) error {

	if d.Rules[r.Name] != nil {
		return fmt.Errorf("AddRule failed, Rule [%s] already exists", r.Name)
	}

	d.Rules[r.Name] = r

	return nil
}

func (d *Director) CreateActor(an string, rn string) error {

	if _, ok := d.Rules[rn]; !ok {
		return fmt.Errorf("CreateActor failed for Actor [%s], Rule [%s] does not exist", an, rn)
	}

	a := NewActor(an, d)
	a.Run()
	d.Actors[an] = a

	if _, ok := d.Actors[an]; !ok {
		return fmt.Errorf("CreateActor failed for Actor [%s] and Rule [%s]", an, rn)
	}

	return nil
}

func (d *Director) CreateClassicActor(an string, rn string) error {

	if d.Rules[rn] == nil {
		return fmt.Errorf("CreateActor failed for Actor [%s], Rule [%s] does not exist", an, rn)
	}

	a := NewActor(an, d)
	d.Actors[an] = a
	if d.Actors[an] == nil {
		return fmt.Errorf("CreateActor failed for Actor [%s] and Rule [%s]", an, rn)
	}

	return nil
}

// returns the total strikes of an Actor
func (d *Director) Strikes(an string, rn string) (int, error) {
	// return false if actor doesn't exist
	if !d.ActorExists(an) {
		return 0, fmt.Errorf("Actor does not exist [%v]", an)
	}

	if !d.InfractionExists(an, rn) {
		return 0, fmt.Errorf("Infraction does not exist [%v]", rn)
	}

	in := NewIncoming(an, rn, STRIKES)
	d.Actors[an].Incoming <- in
	out := <-in.Outgoing
	close(in.Outgoing)
	i, _ := strconv.Atoi(out.Message)
	return i, nil
}

func (d *Director) createActor(an string, rn string) {
	a := NewActor(an, d)
	d.Actors[an] = a
}

func (d *Director) IsJailed(an string) bool {
	// return false if actor doesn't exist
	if !d.ActorExists(an) {
		return false
	}

	in := NewIncoming(an, "", IS_JAILED)
	d.Actors[an].Incoming <- in
	out := <-in.Outgoing
	res, _ := strconv.ParseBool(out.Message)
	return res
}

func (d *Director) IsJailedFor(an string, rn string) bool {
	// return false if actor doesn't exist
	if !d.ActorExists(an) {
		return false
	}

	// return false if infraction still exists
	if d.InfractionExists(an, rn) {
		return false
	}

	in := NewIncoming(an, rn, IS_JAILED_FOR)
	d.Actors[an].Incoming <- in
	out := <-in.Outgoing
	res, _ := strconv.ParseBool(out.Message)
	return res
}

func (d *Director) IncrementInfraction(an string, rn string) error {
	if !d.ActorExists(an) {
		return fmt.Errorf("Actor [%v] does not exist", an)
	}

	if !d.InfractionExists(an, rn) {
		return fmt.Errorf("Infraction [%v] does not exist", rn)
	}

	if d.IsJailedFor(an, rn) {
		return fmt.Errorf("Actor already jailed for [%v:%v] ", an, rn)
	}

	in := NewIncoming(an, rn, INFRACTION)
	d.Actors[an].Incoming <- in
	out := <-in.Outgoing
	close(in.Outgoing)
	return out.Error
}

// takes an ActorName and RuleName and creates an Infraction
func (d *Director) CreateInfraction(an string, rn string) error {

	if !d.ActorExists(an) {
		return fmt.Errorf("CreateInfraction failed: Actor does not exist")
	}

	if d.InfractionExists(an, rn) {
		return fmt.Errorf("CreateInfraction failed: Infraction already exists")
	}

	in := NewIncoming(an, rn, CREATE_INFRACTION)
	inf := NewInfraction(d.Rules[rn])
	in.Infraction = inf
	d.Actors[an].Incoming <- in
	out := <-in.Outgoing
	close(in.Outgoing)
	return out.Error
}

// ActorExists asserts that you are checking for an Actor
// in anticipation to operate on said Actor in the near future
func (d *Director) ActorExists(an string) bool {
	_, ok := d.Actors[an]
	return ok
}

func (d *Director) InfractionExists(an string, rn string) bool {
	if !d.ActorExists(an) {
		return false
	}

	if _, ok := d.Actors[an].Infractions[rn]; ok {
		return true
	}
	return false
}
