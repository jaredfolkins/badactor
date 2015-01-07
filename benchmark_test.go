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

func BenchmarkInfraction(b *testing.B) {
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

	for i := 0; i < b.N; i++ {
		an := strconv.FormatInt(rand.Int63(), 10)
		d.lInfraction(an, rn)
	}
}

func BenchmarkInfractionlIsJailed(b *testing.B) {
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
		an := strconv.FormatInt(rand.Int63(), 10)
		for i := 0; i < 3; i++ {
			d.lInfraction(an, rn)
		}
		d.lIsJailed(an)
	}
}

func BenchmarkInfractionlIsJailedFor(b *testing.B) {
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
		an := strconv.FormatInt(rand.Int63(), 10)
		for i := 0; i < 3; i++ {
			d.lInfraction(an, rn)
		}
		d.lIsJailedFor(an, rn)
	}
}
