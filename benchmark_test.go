package badactor

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func BenchmarkIsJailed(b *testing.B) {
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

	err = d.AddRule(r)
	if err != nil {
		b.Errorf("AddRule [%s] should not fail", rn)
	}

	for i := 0; i < 4; i++ {
		d.Infraction(an, rn)
	}

	for i := 0; i < b.N; i++ {
		d.IsJailed(an)
	}
}

func BenchmarkIsJailedFor(b *testing.B) {
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

	err = d.AddRule(r)
	if err != nil {
		b.Errorf("AddRule [%s] should not fail", rn)
	}

	for i := 0; i < 4; i++ {
		d.Infraction(an, rn)
	}

	for i := 0; i < b.N; i++ {
		d.IsJailedFor(an, rn)
	}
}

func BenchmarkInfraction(b *testing.B) {
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

	err = d.AddRule(r)
	if err != nil {
		b.Errorf("AddRule for [%v] should not fail", rn)
	}

	for i := 0; i < b.N; i++ {
		an = string(rand.Int63())
		d.Infraction(an, rn)
	}
}

func BenchmarkInfractionMostCostly(b *testing.B) {
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

	err = d.AddRule(r)
	if err != nil {
		b.Errorf("AddRule for [%s] should not fail", rn)
	}

	// bench the Least Costly way
	for i := 0; i < b.N; i++ {
		an = strconv.FormatInt(rand.Int63(), 10)
		d.MostCostlyInfraction(an, rn)
	}
}

func BenchmarkInfractionIsJailed(b *testing.B) {
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

	err = d.AddRule(r)
	if err != nil {
		b.Errorf("AddRule [%s] should not fail", rn)
	}

	for i := 0; i < b.N; i++ {
		an = string(rand.Int63())
		for i := 0; i < 3; i++ {
			d.Infraction(an, rn)
		}
		d.IsJailed(an)
	}
}

func BenchmarkInfractionIsJailedFor(b *testing.B) {
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

	err = d.AddRule(r)
	if err != nil {
		b.Errorf("AddRule [%s] should not fail", rn)
	}

	for i := 0; i < b.N; i++ {
		an = string(rand.Int63())
		for i := 0; i < 3; i++ {
			d.Infraction(an, rn)
		}
		d.IsJailedFor(an, rn)
	}
}

func Benchmark10000Actors1Infraction(b *testing.B) {
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

	err = d.AddRule(r)
	if err != nil {
		b.Errorf("AddRule for [%v] should not fail", rn)
	}

	aN := 10000
	for i := 0; i < b.N; i++ {
		for a := 0; a < aN; a++ {
			an = string(rand.Int63())
			d.Infraction(an, rn)
		}
	}
}

func Benchmark100000Actors1Infraction(b *testing.B) {
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

	err = d.AddRule(r)
	if err != nil {
		b.Errorf("AddRule for [%v] should not fail", rn)
	}

	aN := 100000
	for i := 0; i < b.N; i++ {
		for a := 0; a < aN; a++ {
			an = string(rand.Int63())
			d.Infraction(an, rn)
		}
	}
}

func Benchmark1000000Actors1Infraction(b *testing.B) {
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

	err = d.AddRule(r)
	if err != nil {
		b.Errorf("AddRule for [%v] should not fail", rn)
	}

	aN := 1000000
	for i := 0; i < b.N; i++ {
		for a := 0; a < aN; a++ {
			an = string(rand.Int63())
			d.Infraction(an, rn)
		}
	}
}

func Benchmark10000Actors4Infractions(b *testing.B) {
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

	err = d.AddRule(r)
	if err != nil {
		b.Errorf("AddRule for [%v] should not fail", rn)
	}

	aN := 10000

	for i := 0; i < b.N; i++ {
		for a := 0; a < aN; a++ {
			an = string(rand.Int63())
			for inf := 0; inf < 4; inf++ {
				d.Infraction(an, rn)
			}
		}
	}
}
