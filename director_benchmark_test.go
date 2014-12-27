package badactor

import (
	"strconv"
	"testing"
	"time"
)

func BenchmarkInfraction(b *testing.B) {
	var err error
	d := NewDirector()
	d.Run()
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
		b.Errorf("AddRule for Actor [%s] should not fail", an)
	}

	for i := 0; i < b.N; i++ {
		d.Infraction(an, rn)
	}
}

func BenchmarkInfractionWithIsJailed(b *testing.B) {
	var err error
	d := NewDirector()
	d.Run()
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
		b.Errorf("AddRule for Actor [%s] should not fail", an)
	}

	for i := 0; i < b.N; i++ {
		if !d.IsJailed(an) {
			d.Infraction(an, rn)
		}
	}
}

func BenchmarkInfractionWithIsJailed2000(b *testing.B) {
	var err error
	d := NewDirector()
	d.Run()
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
		b.Errorf("AddRule [%v] should not fail", rn)
	}

	for i := 0; i < b.N; i++ {
		for a := 0; a < 2000; a++ {
			an := "an_" + strconv.FormatInt(time.Now().UnixNano(), 10)
			if !d.IsJailed(an) {
				d.Infraction(an, rn)
			}
		}
	}
}
