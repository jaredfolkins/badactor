package badactor

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func BenchmarklIsJailed(b *testing.B) {
	var err error
	d := NewDirector(ia)
	an := "an_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Second * 60,
		Sentence:    time.Minute * 5,
	}

	err = d.lAddRule(r)
	if err != nil {
		b.Errorf("lAddRule [%s] should not fail", rn)
	}

	for i := 0; i < 4; i++ {
		d.lInfraction(an, rn)
	}

	for i := 0; i < b.N; i++ {
		d.lIsJailed(an)
	}
}

func BenchmarklIsJailedFor(b *testing.B) {
	var err error
	d := NewDirector(ia)
	an := "an_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Second * 60,
		Sentence:    time.Minute * 5,
	}

	err = d.lAddRule(r)
	if err != nil {
		b.Errorf("lAddRule [%s] should not fail", rn)
	}

	for i := 0; i < 4; i++ {
		d.lInfraction(an, rn)
	}

	for i := 0; i < b.N; i++ {
		d.lIsJailedFor(an, rn)
	}
}

func BenchmarklInfraction(b *testing.B) {
	var err error
	d := NewDirector(ia)
	an := "an_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Second * 60,
		Sentence:    time.Minute * 5,
	}

	err = d.lAddRule(r)
	if err != nil {
		b.Errorf("lAddRule for [%v] should not fail", rn)
	}

	for i := 0; i < b.N; i++ {
		an = string(rand.Int63())
		d.lInfraction(an, rn)
	}
}

func BenchmarklInfractionMostCostly(b *testing.B) {
	var an string
	var err error
	d := NewDirector(ia)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Second * 60,
		Sentence:    time.Minute * 5,
	}

	err = d.lAddRule(r)
	if err != nil {
		b.Errorf("lAddRule for [%s] should not fail", rn)
	}

	// bench the Least Costly way
	for i := 0; i < b.N; i++ {
		an = strconv.FormatInt(rand.Int63(), 10)
		d.costlyInfraction(an, rn)
	}
}

func BenchmarklInfractionlIsJailed(b *testing.B) {
	var an string
	var err error
	d := NewDirector(ia)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Second * 60,
		Sentence:    time.Minute * 5,
	}

	err = d.lAddRule(r)
	if err != nil {
		b.Errorf("lAddRule [%s] should not fail", rn)
	}

	for i := 0; i < b.N; i++ {
		an = string(rand.Int63())
		for i := 0; i < 3; i++ {
			d.lInfraction(an, rn)
		}
		d.lIsJailed(an)
	}
}

func BenchmarklInfractionlIsJailedFor(b *testing.B) {
	var an string
	var err error
	d := NewDirector(ia)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Second * 60,
		Sentence:    time.Minute * 5,
	}

	err = d.lAddRule(r)
	if err != nil {
		b.Errorf("lAddRule [%s] should not fail", rn)
	}

	for i := 0; i < b.N; i++ {
		an = string(rand.Int63())
		for i := 0; i < 3; i++ {
			d.lInfraction(an, rn)
		}
		d.lIsJailedFor(an, rn)
	}
}

func Benchmark10000Actors1lInfraction(b *testing.B) {
	var an string
	var err error
	d := NewDirector(ia)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Second * 60,
		Sentence:    time.Minute * 5,
	}

	err = d.lAddRule(r)
	if err != nil {
		b.Errorf("lAddRule for [%v] should not fail", rn)
	}

	aN := 10000
	for i := 0; i < b.N; i++ {
		for a := 0; a < aN; a++ {
			an = string(rand.Int63())
			d.lInfraction(an, rn)
		}
	}
}

func Benchmark100000Actors1lInfraction(b *testing.B) {
	var an string
	var err error
	d := NewDirector(ia)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Second * 60,
		Sentence:    time.Minute * 5,
	}

	err = d.lAddRule(r)
	if err != nil {
		b.Errorf("lAddRule for [%v] should not fail", rn)
	}

	aN := 100000
	for i := 0; i < b.N; i++ {
		for a := 0; a < aN; a++ {
			an = string(rand.Int63())
			d.lInfraction(an, rn)
		}
	}
}

func Benchmark1000000Actors1lInfraction(b *testing.B) {
	var an string
	var err error

	d := NewDirector(ia)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 1,
		ExpireBase:  time.Second * 60,
		Sentence:    time.Minute * 5,
	}

	err = d.lAddRule(r)
	if err != nil {
		b.Errorf("lAddRule for [%v] should not fail", rn)
	}

	aN := 1000000
	for i := 0; i < b.N; i++ {
		for a := 0; a < aN; a++ {
			an = string(rand.Int63())
			d.lInfraction(an, rn)
		}
	}
}

func Benchmark10000Actors4lInfractions(b *testing.B) {
	var an string
	var err error
	d := NewDirector(ia)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Second * 60,
		Sentence:    time.Minute * 5,
	}

	err = d.lAddRule(r)
	if err != nil {
		b.Errorf("lAddRule for [%v] should not fail", rn)
	}

	aN := 10000

	for i := 0; i < b.N; i++ {
		for a := 0; a < aN; a++ {
			an = string(rand.Int63())
			for inf := 0; inf < 4; inf++ {
				d.lInfraction(an, rn)
			}
		}
	}
}
