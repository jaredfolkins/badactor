package badactor

import (
	"strconv"
	"testing"
	"time"
)

var whenJailedCounter = 0
var whenTimeServedCounter = 0

type MockAction struct{}

func (ma *MockAction) Log(a *Actor, r *Rule) {}

func (ma *MockAction) WhenJailed(a *Actor, r *Rule) error {
	whenJailedCounter++
	ma.Log(a, r)
	return nil
}

func (ma *MockAction) WhenTimeServed(a *Actor, r *Rule) error {
	whenTimeServedCounter++
	ma.Log(a, r)
	return nil
}

func TestActionWhenJailed(t *testing.T) {
	//setup
	var b bool
	d := NewDirector(ia)
	an := "an_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	sl := 3
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: sl,
		ExpireBase:  time.Second * 60,
		Sentence:    time.Second * 60,
		Action:      &MockAction{},
	}
	a := newClassicActor(an, r, d)

	//test
	for i := 0; i < 3; i++ {
		a.infraction(rn)
	}

	b = a.isJailedFor(an)
	if b == true {
		t.Errorf("isJailedFor should be false instead [%v]", b)
	}

	if whenJailedCounter != 1 {
		t.Errorf("whenJailedCounter should be 1 instead [%v]", whenJailedCounter)
	}

}

func TestActionWhenTimeServed(t *testing.T) {
	//setup
	var b bool

	d := NewDirector(ia)
	an := "an_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	sl := 3
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: sl,
		ExpireBase:  time.Second * 1,
		Sentence:    time.Millisecond * 10,
		Action:      &MockAction{},
	}
	a := newClassicActor(an, r, d)

	if whenTimeServedCounter != 0 {
		t.Errorf("whenTimeServedCounter should be 0 instead", whenTimeServedCounter)
	}
	//test
	for i := 0; i < 3; i++ {
		err := a.infraction(rn)
		if err != nil {
			t.Errorf("infraction %v should not error", err)
		}
	}

	if whenTimeServedCounter != 0 {
		t.Errorf("whenTimeServedCounter should be 0 instead", whenTimeServedCounter)
	}

	b = a.isJailed()
	if b == false {
		t.Errorf("isJailed should be true instead [%v]", b)
	}

	b = a.isJailedFor(rn)
	if b == false {
		t.Errorf("isJailedFor should be true instead [%v]", b)
	}

	if whenTimeServedCounter != 0 {
		t.Errorf("whenTimeServedCounter should be 0 instead", whenTimeServedCounter)
	}

	dur := time.Duration(time.Second * 1)

	time.Sleep(dur)

	for _, j := range a.jails {
		a.timeServed(j)
	}

	if whenTimeServedCounter != 1 {
		t.Errorf("whenTimeServedCounter should be 1 instead", whenTimeServedCounter)
	}

}
