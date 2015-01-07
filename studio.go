package badactor

import (
	"hash/fnv"
	"math/rand"
	"sync"
	"time"
)

// maxNs is used as the max value in the range for the Ticker
// minNs is used as the min value in the range for the Ticker
const (
	maxNs = 3000000000
	minNs = 1000000000
)

type Studio struct {
	sync.Mutex
	maxDirectors int32
	directors    map[int32]*Director
	rules        map[string]*Rule
}

func NewStudio(md int32) *Studio {
	st := &Studio{
		maxDirectors: md,
		directors:    make(map[int32]*Director, md),
		rules:        make(map[string]*Rule),
	}
	return st
}

// AddRule accepts a Rule struct and adds it to the rules map if it doesn't exist
func (st *Studio) AddRule(r *Rule) {
	st.Lock()
	st.addRule(r)
	st.Unlock()
	return
}

func (st *Studio) addRule(r *Rule) {
	st.rules[r.Name] = r
	return
}

func (st *Studio) ApplyRules() {
	for _, r := range st.rules {
		for _, d := range st.directors {
			d.lAddRule(r)
		}
	}
	return
}

func (st *Studio) CreateDirectors(ma int32) error {
	var dk int32

	for dk = 0; dk < st.maxDirectors; dk++ {
		d := NewDirector(ma)
		st.directors[dk] = d
	}

	st.ApplyRules()

	return nil
}

// Infraction accepts an ActorName and RuleName and either creates, increments, or increments and jails the Actor
func (st *Studio) Infraction(an string, rn string) error {
	d := st.GetDirector(an)
	return d.lInfraction(an, rn)
}

// Strikes accepts an ActorName and a RuleName and returns the total strikes an Actor holds for a particular Rule
func (st *Studio) Strikes(an string, rn string) (int, error) {
	d := st.GetDirector(an)
	return d.lStrikes(an, rn)
}

// CreateInfraction takes and ActorName and RuleName and creates an Infraction
func (st *Studio) CreateInfraction(an string, rn string) error {
	d := st.GetDirector(an)
	return d.lCreateInfraction(an, rn)
}

// CreateActor takes and ActorName and RuleName and creates an Actor
func (st *Studio) CreateActor(an string, rn string) error {
	d := st.GetDirector(an)
	return d.lCreateInfraction(an, rn)
}

// KeepAlive accepts an ActorName and allows you to rebase the TTL for the Actor so that it isn't removed from the stack as scheduled
func (st *Studio) KeepAlive(an string) error {
	d := st.GetDirector(an)
	return d.lKeepAlive(an)
}

// ActorExists accepts an ActorName and returns a bool if the Actor is found
func (st *Studio) ActorExists(an string) bool {
	d := st.GetDirector(an)
	return d.lActorExists(an)
}

// InfractionExists accepts an ActorName and RuleName and returns a bool if the Infraction is found
func (st *Studio) InfractionExists(an string, rn string) bool {
	d := st.GetDirector(an)
	return d.lInfractionExists(an, rn)
}

// IsJailedFor accepts an ActorName and a RuleName and returns a bool if the Actor is Jailed for that particular Rule
func (st *Studio) IsJailedFor(an string, rn string) bool {
	d := st.GetDirector(an)
	return d.lIsJailedFor(an, rn)
}

// IsJailed accepts an ActorName and returns a bool if the Actor is Jailed for ANY Rule
func (st *Studio) IsJailed(an string) bool {
	d := st.GetDirector(an)
	return d.lIsJailed(an)
}

func (st *Studio) Run() {
	r := time.Duration(rand.Intn(maxNs-minNs) + 1)
	ticker := time.NewTicker(time.Nanosecond * r)
	go func() {
		// 1.4 means i refractor this
		for _ = range ticker.C {
			for _, d := range st.directors {
				d.Maintenance()
			}
		}
	}()
}

// serialize a string to an uint64
func (st Studio) serializeId(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func (st Studio) GetDirector(an string) *Director {
	dk := st.jHash(st.serializeId(an), st.maxDirectors)
	return st.directors[dk]
}

func (st Studio) jHash(key uint64, numBuckets int32) int32 {

	var b int64 = -1
	var j int64

	for j < int64(numBuckets) {
		b = j
		key = key*2862933555777941757 + 1
		j = int64(float64(b+1) * (float64(int64(1)<<31) / float64((key>>33)+1)))
	}
	return int32(b)
}
