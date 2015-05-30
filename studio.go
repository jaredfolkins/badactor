package badactor

import (
	"hash/fnv"
	"sync"
	"time"
)

// Studio is the singleton instance, it contains the Directors(buckets) who have many Actors(points)
type Studio struct {
	sync.Mutex
	capacity  int32
	directors map[int32]*Director
	rules     map[string]*Rule
	status    chan *status
}

// NewStudio returns a init'd Studio struct, you pass it an int32 value which is the capacity and informs the Studio how many Directors will be created, it is also the value that jumpHash uses to mod
func NewStudio(md int32) *Studio {
	return &Studio{
		capacity:  md,
		directors: make(map[int32]*Director, md),
		rules:     make(map[string]*Rule),
		status:    make(chan *status),
	}
}

// AddRule accepts a Rule struct and adds it to the rules map if it doesn't exist
func (st *Studio) AddRule(r *Rule) {
	st.Lock()
	st.rules[r.Name] = r
	st.Unlock()
	return
}

// ApplyRules takes the currently stored rules map and applies it to all Directors
func (st *Studio) ApplyRules() error {
	for _, d := range st.directors {
		for _, r := range st.rules {
			d.lAddRule(r)
		}
	}
	return nil
}

// CreateDirectors creates and adds the Directors to the director map
func (st *Studio) CreateDirectors(ma int32) error {
	var dk int32
	for dk = 0; dk < st.capacity; dk++ {
		d := NewDirector(ma)
		st.directors[dk] = d
	}
	return st.ApplyRules()
}

// Infraction accepts an ActorName and RuleName and either creates, increments, or increments and jails the Actor
func (st *Studio) Infraction(an string, rn string) error {
	d := st.Director(an)
	return d.lInfraction(an, rn)
}

// Strikes accepts an ActorName and a RuleName and returns the total strikes an Actor holds for a particular Rule
func (st *Studio) Strikes(an string, rn string) (int, error) {
	d := st.Director(an)
	return d.lStrikes(an, rn)
}

// CreateInfraction takes an ActorName and RuleName and creates an Infraction
func (st *Studio) CreateInfraction(an string, rn string) error {
	d := st.Director(an)
	return d.lCreateInfraction(an, rn)
}

// CreateActor takes an ActorName and RuleName and creates an Actor
func (st *Studio) CreateActor(an string, rn string) error {
	d := st.Director(an)
	return d.lCreateActor(an, rn)
}

// KeepAlive accepts an ActorName and allows you to rebase the TTL for the Actor so that it isn't removed from the stack as scheduled, keeping it alive.
func (st *Studio) KeepAlive(an string) error {
	d := st.Director(an)
	return d.lKeepAlive(an)
}

// ActorExists accepts an ActorName and returns a bool if the Actor is found
func (st *Studio) ActorExists(an string) bool {
	d := st.Director(an)
	return d.lActorExists(an)
}

// InfractionExists accepts an ActorName and RuleName and returns a bool if the Infraction is found
func (st *Studio) InfractionExists(an string, rn string) bool {
	d := st.Director(an)
	return d.lInfractionExists(an, rn)
}

// IsJailedFor accepts an ActorName and a RuleName and returns a bool if the Actor is Jailed for that particular Rule
func (st *Studio) IsJailedFor(an string, rn string) bool {
	d := st.Director(an)
	return d.lIsJailedFor(an, rn)
}

// IsJailed accepts an ActorName and returns a bool if the Actor is Jailed for ANY Rule
func (st *Studio) IsJailed(an string) bool {
	d := st.Director(an)
	return d.lIsJailed(an)
}

// StartReaper starts the reaping goroutine and takes a time.Duration on how often you want the Reaper to run
func (st *Studio) StartReaper(dur time.Duration) {
	ticker := time.NewTicker(dur)
	go func() {
		for {
			select {
			case <-ticker.C:
				for _, d := range st.directors {
					d.lMaintenance()
				}
			case stat := <-st.status:
				m := &message{
					reaperAlive: true,
				}
				stat.outgoing <- m
			}
		}
	}()
}

func (st *Studio) Status() *message {
	stat := newStatus()
	defer close(stat.outgoing)
	st.status <- stat
	return <-stat.outgoing
}

// Director takes the name of an Actor as a string, serializes it, uses the jumpHash aglo to determine the Director that the Actor belongs to
func (st Studio) Director(an string) *Director {
	dk := st.jumpHash(st.serialize(an), st.capacity)
	return st.directors[dk]
}

// serialize a string to an uint64
func (st Studio) serialize(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func (st Studio) jumpHash(key uint64, numBuckets int32) int32 {

	var b int64 = -1
	var j int64

	for j < int64(numBuckets) {
		b = j
		key = key*2862933555777941757 + 1
		j = int64(float64(b+1) * (float64(int64(1)<<31) / float64((key>>33)+1)))
	}
	return int32(b)
}
