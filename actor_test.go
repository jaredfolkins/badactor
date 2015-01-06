package badactor

import (
	"strconv"
	"testing"
	"time"
)

const ia = 1024

func TestActorIsJailedFor(t *testing.T) {
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
	}
	a := newClassicActor(an, r, d)

	//test
	for i := 0; i < 3; i++ {
		a.lInfraction(rn)
	}

	b = a.lIsJailedFor(an)
	if b {
		t.Errorf("isJailedFor should be false instead [%v]", b)
	}

	for i := 0; i < 3; i++ {
		a.lInfraction(rn)
	}

	b = a.lIsJailedFor(rn)
	if !b {
		t.Errorf("isJailedFor should be false instead [%v]", b)
	}

}

func TestActorLockup(t *testing.T) {
	// setup
	var b bool
	an := "an_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	d := NewDirector(ia)
	sl := 3
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: sl,
		ExpireBase:  time.Second * 60,
		Sentence:    time.Millisecond * 3,
	}
	a := newClassicActor(an, r, d)

	// test
	for i := 0; i < 3; i++ {
		a.lInfraction(rn)
	}

	b = a.lIsJailedFor(rn)
	if b == false {
		t.Errorf("isJailedFor should be true instead [%v]", b)
	}

}

func TestActorTimeServed(t *testing.T) {
	// setup
	var err error
	var b bool
	an := "an_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	dur := time.Millisecond * 10
	d := NewDirector(ia)
	sl := 3
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: sl,
		ExpireBase:  time.Second * 60,
		Sentence:    time.Millisecond * 50,
	}
	a := newClassicActor(an, r, d)

	b = a.lIsJailedFor(rn)
	if b == true {
		t.Errorf("isJailedFor should be false instead [%v]", b)
	}

	// test
	for i := 0; i < 3; i++ {
		err = a.lInfraction(rn)
		if err != nil {
			t.Errorf("Infraction should not error instead [%v]", err)
		}
	}

	// .lockup should happen in the goroutine background
	time.Sleep(dur)
	//a.lMaintenance()
	b = a.lIsJailedFor(rn)
	if b == false {
		t.Errorf("isJailedFor should be true instead [%v]", b)
	}

	time.Sleep(time.Millisecond * 40)
	//a.lMaintenance()
	b = a.lIsJailedFor(rn)
	if b == true {
		t.Errorf("isJailedFor should be false instead [%v]", b)
	}

}

func TestActorlExpire(t *testing.T) {

	var err error
	br := "badrule"
	an := "an_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	dur := time.Millisecond * 40
	d := NewDirector(ia)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Second * 1,
		Sentence:    time.Millisecond * 50,
	}

	a := newClassicActor(an, r, d)

	err = a.lExpire(br)
	if err == nil {
		t.Errorf("Expire should error for [%v]", err)
	}

	a.lInfraction(rn)

	err = a.lExpire(rn)
	if err != nil {
		t.Errorf("Expire should not fail: %v", err)
	}

	time.Sleep(dur)

	err = a.lExpire(rn)
	if err != nil {
		t.Errorf("Expire should not delete [%v:%v]", an, rn)
	}

	if !a.lHasInfractions() {
		t.Errorf("Infractions should exist [%v:%v]", an, rn)
	}

	time.Sleep(dur)
	err = a.lExpire(rn)
	if err != nil {
		t.Errorf("Expire should delete [%v:%v]", an, rn)
	}

	if !a.lHasInfractions() {
		t.Errorf("Infractions should not exist [%v:%v]", an, rn)
	}

}

func TestActorInfraction(t *testing.T) {

	d := NewDirector(ia)
	var err error
	an := "an_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	sl := 3
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: sl,
		ExpireBase:  time.Second * 60,
		Sentence:    time.Millisecond * 10,
	}

	a := newClassicActor(an, r, d)

	rde := "ruledoesntexist"
	err = a.lInfraction(rde)
	if err == nil {
		t.Errorf("Infraction should error [%v:%v:%v]", err, an, rde)
	}

	for i := 0; i < 3; i++ {
		err = a.lInfraction(rn)
		if err != nil {
			t.Errorf("Infraction should not error [%v:%v]", an, rn)
		}
	}

	err = a.lLockup(rn)
	if err == nil {
		t.Errorf("Lockup err should not be nil instead [%v]", err)
	}

	err = a.lInfraction(rn)
	if err == nil {
		t.Errorf("Infraction should error [%v:%v]", an, rn)
	}

}

func TestActorCreateInfraction(t *testing.T) {

	d := NewDirector(ia)
	var err error
	an := "an_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Millisecond * 10,
		Sentence:    time.Millisecond * 10,
	}

	a := newClassicActor(an, r, d)

	inf := newInfraction(r)

	err = a.lCreateInfraction(inf)
	if err == nil {
		t.Errorf("CreateInfraction [%v:%v] should error %v", an, rn, err)
	}

	ninf := a.infractions[rn]
	if ninf.rule.Name != rn {
		t.Errorf("CreateInfraction Name incorrect [%v:%v]", an, rn)
	}

}

func TestActorRebaseAll(t *testing.T) {

	d := NewDirector(ia)
	var err error
	an := "an_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Millisecond * 10,
		Sentence:    time.Millisecond * 10,
	}
	a := newClassicActor(an, r, d)

	inf := newInfraction(r)

	err = a.lCreateInfraction(inf)
	if err == nil {
		t.Errorf("CreateInfraction [%v:%v] should error %v", an, rn, err)
	}

	ts := time.Now()
	a.rebaseAll()
	for _, inf := range a.infractions {
		if inf.expireBy.Before(ts) {
			t.Errorf("RebaseAll failed ExpireBy time is before TimeStamp [%v:%v:%v:%v]", an, rn, inf.expireBy, ts)
		}
	}

}

func TestActorStrikes(t *testing.T) {
	var err error
	var i int
	bn := "badname"
	d := NewDirector(ia)
	an := "an_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Millisecond * 10,
		Sentence:    time.Millisecond * 10,
	}
	a := newClassicActor(an, r, d)

	i = a.strikes(rn)
	if i != 0 {
		t.Errorf("Strike count for [%v:%v] must be 1", an, rn)
	}

	err = a.infraction(rn)
	if err != nil {
		t.Errorf("Infraction should not error")
	}

	i = a.strikes(rn)
	if i != 1 {
		t.Errorf("Strike count for [%v:%v] must be 1", an, rn)
	}

	i = a.strikes(bn)
	if i != 0 {
		t.Errorf("Badname should count for [%v] not be 0", bn)
	}

}

func TestActorShouldReturn(t *testing.T) {
	// setup
	var err error
	var b bool
	d := NewDirector(ia)
	dur := time.Millisecond * 100
	an := "an_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Second * 1,
		Sentence:    time.Millisecond * 50,
	}
	a := newClassicActor(an, r, d)

	// test
	// assert all falsey
	b = a.lShouldDelete()
	if b {
		t.Errorf("shouldDelete should be false instead [%v]", b)
	}

	b = a.lIsJailedFor(rn)
	if b == true {
		t.Errorf("isJailedFor should be false instead [%v]", b)
	}

	b = a.lHasJails()
	if b == true {
		t.Errorf("hasJails should be false instead [%v]", b)
	}

	// assert all truthy
	for i := 0; i < 3; i++ {
		err = a.lInfraction(rn)
		if err != nil {
			t.Errorf("Infraction [%v] should not error %v", rn, err)
		}
	}

	//a.lMaintenance()

	b = a.lIsJailedFor(rn)
	if b == false {
		t.Errorf("isJailedFor should be true instead [%v]", b)
	}

	b = a.lHasJails()
	if b == false {
		t.Errorf("hasJails should be true instead [%v]", b)
	}

	// sleep, quit, should NOT be jailed
	time.Sleep(dur)

	//a.lMaintenance()

	b = a.lIsJailedFor(rn)
	if b == true {
		t.Errorf("isJailedFor should be false instead [%v]", b)
	}

	b = a.lHasJails()
	if b == true {
		t.Errorf("hasJails should be false instead [%v]", b)
	}

}

func TestActorInfractionExists(t *testing.T) {
	// setup
	var b bool
	var err error
	d := NewDirector(ia)
	an := "an_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Second * 1,
		Sentence:    time.Millisecond * 50,
	}
	a := newActor(an, d)

	b = a.infractionExists(rn)
	if b == true {
		t.Errorf("infractionExists should be false instead [%v]", b)
	}

	inf := newInfraction(r)
	err = a.createInfraction(inf)
	if err != nil {
		t.Errorf("createInfraction should not error %v", err)
	}

	b = a.infractionExists(rn)
	if b == false {
		t.Errorf("infractionExists should be true instead [%v]", b)
	}
}

func TestActorTotalJails(t *testing.T) {
	// setup
	d := NewDirector(ia)
	an := "an_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Second * 1,
		Sentence:    time.Millisecond * 50,
	}
	a := newClassicActor(an, r, d)

	for i := 0; i < 3; i++ {
		a.infraction(rn)
	}

	i := a.totalJails()
	if i != 1 {
		t.Errorf("totalJails should be 1 instead [%v]", i)

	}

}
