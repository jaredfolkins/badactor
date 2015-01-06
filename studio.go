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

func (st *Studio) AddRule(r *Rule) {
	st.rules[r.Name] = r
}

func (st *Studio) ApplyRules() {
	for _, r := range st.rules {
		for _, d := range st.directors {
			d.AddRule(r)
		}
	}
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

func (st *Studio) Infraction(an string, rn string) error {
	d := st.GetDirector(an)
	return d.Infraction(an, rn)
}

func (st *Studio) maintenance() {
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
