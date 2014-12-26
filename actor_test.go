package badactor

import (
	"strconv"
	"testing"
	"time"
)

func TestActorIsJailedFor(t *testing.T) {
	//setup
	var b bool
	d := NewDirector()
	d.Run()
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
	a := NewClassicActor(an, r, d)
	a.Run()

	//test
	for i := 1; i < 3; i++ {
		a.infraction(rn)
	}

	b = a.isJailedFor(an)
	if b {
		t.Errorf("isJailedFor should be false instead [%v]", b)
	}

	for i := 0; i < 3; i++ {
		a.infraction(rn)
	}

	b = a.isJailedFor(rn)
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
	d := NewDirector()
	d.Run()
	sl := 3
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: sl,
		ExpireBase:  time.Second * 60,
		Sentence:    time.Millisecond * 3,
	}
	a := NewClassicActor(an, r, d)
	a.Run()

	// test
	for i := 0; i < 3; i++ {
		a.infraction(rn)
	}

	b = a.isJailedFor(rn)
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
	d := NewDirector()
	d.Run()
	sl := 3
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: sl,
		ExpireBase:  time.Second * 60,
		Sentence:    time.Millisecond * 50,
	}
	a := NewClassicActor(an, r, d)
	a.Run()

	b = a.isJailedFor(rn)
	if b == true {
		t.Errorf("isJailedFor should be false instead [%v]", b)
	}

	// test
	for i := 1; i < 3; i++ {
		err = a.infraction(rn)
		if err != nil {
			t.Errorf("Infraction should not error instead [%v]", err)
		}
	}

	// .lockup should happen in the goroutine background
	time.Sleep(dur)
	b = a.isJailedFor(rn)
	if b == false {
		t.Errorf("isJailedFor should be true instead [%v]", b)
	}

	// timeServed should fire
	time.Sleep(time.Millisecond * 40)
	b = a.isJailedFor(rn)
	if b == true {
		t.Errorf("isJailedFor should be false instead [%v]", b)
	}

}

func TestActorExpire(t *testing.T) {

	var err error
	br := "badrule"
	an := "an_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	dur := time.Millisecond * 40
	d := NewDirector()
	d.Run()
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Second * 1,
		Sentence:    time.Millisecond * 50,
	}

	a := NewClassicActor(an, r, d)
	a.Run()

	err = a.expire(br)
	if err == nil {
		t.Errorf("Expire should error for [%v]", err)
	}

	a.infraction(rn)

	err = a.expire(rn)
	if err != nil {
		t.Errorf("Expire should not fail: %v", err)
	}

	time.Sleep(dur)

	err = a.expire(rn)
	if err != nil {
		t.Errorf("Expire should not delete [%v:%v]", an, rn)
	}

	if !a.hasInfractions() {
		t.Errorf("Infractions should exist [%v:%v]", an, rn)
	}

	time.Sleep(dur)
	err = a.expire(rn)
	if err != nil {
		t.Errorf("Expire should delete [%v:%v]", an, rn)
	}

	if !a.hasInfractions() {
		t.Errorf("Infractions should not exist [%v:%v]", an, rn)
	}

}

func TestActorInfraction(t *testing.T) {

	d := NewDirector()
	d.Run()
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

	a := NewClassicActor(an, r, d)
	a.Run()

	rde := "ruledoesntexist"
	err = a.infraction(rde)
	if err == nil {
		t.Errorf("Infraction should error [%v:%v:%v]", err, an, rde)
	}

	for i := 1; i < 3; i++ {
		err = a.infraction(rn)
		if err != nil {
			t.Errorf("Infraction should not error [%v:%v]", an, rn)
		}
	}

	err = a.lockup(rn)
	if err == nil {
		t.Errorf("Lockup err should not be nil instead [%v]", err)
	}

	err = a.infraction(rn)
	if err == nil {
		t.Errorf("Infraction should error [%v:%v]", an, rn)
	}

}

func TestActorCreateInfraction(t *testing.T) {

	d := NewDirector()
	d.Run()
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

	a := NewClassicActor(an, r, d)
	a.Run()

	inf := NewInfraction(r)

	err = a.createInfraction(inf)
	if err == nil {
		t.Errorf("CreateInfraction [%v:%v] should error %v", an, rn, err)
	}

	ninf := a.Infractions[rn]
	if ninf.Rule.Name != rn {
		t.Errorf("CreateInfraction Name incorrect [%v:%v]", an, rn)
	}

}

func TestActorRebaseAll(t *testing.T) {

	d := NewDirector()
	d.Run()
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
	a := NewClassicActor(an, r, d)
	a.Run()

	inf := NewInfraction(r)

	err = a.createInfraction(inf)
	if err == nil {
		t.Errorf("CreateInfraction [%v:%v] should error %v", an, rn, err)
	}

	ts := time.Now()
	err = a.rebaseAll()
	for _, inf := range a.Infractions {
		if inf.ExpireBy.Before(ts) {
			t.Errorf("RebaseAll failed ExpireBy time is before TimeStamp [%v:%v:%v:%v]", an, rn, inf.ExpireBy, ts)
		}
	}

}

func TestActorStrikes(t *testing.T) {
	var err error
	var i int
	bn := "badname"
	d := NewDirector()
	d.Run()
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
	a := NewClassicActor(an, r, d)
	a.Run()

	i = a.strikes(rn)
	if i != 1 {
		t.Errorf("Strike count for [%v:%v] must be 1", an, rn)
	}

	err = a.infraction(rn)
	if err != nil {
		t.Errorf("Infraction should not error")
	}

	i = a.strikes(rn)
	if i != 2 {
		t.Errorf("Strike count for [%v:%v] must be 2", an, rn)
	}

	i = a.strikes(bn)
	if i != 0 {
		t.Errorf("Badname should count for [%v] not be 0", bn)
	}

}

func TestActorQuit(t *testing.T) {
	// setup
	var err error
	var b bool
	d := NewDirector()
	d.Run()
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
	a := NewClassicActor(an, r, d)
	a.Run()

	// test
	// assert all falsey
	b = a.shouldReturn()
	if b {
		t.Errorf("shouldReturn should be false instead [%v]", b)
	}

	b = a.isJailedFor(rn)
	if b == true {
		t.Errorf("isJailedFor should be false instead [%v]", b)
	}

	b = a.hasJails()
	if b == true {
		t.Errorf("hasJails should be false instead [%v]", b)
	}

	// assert all truthy
	for i := 1; i < 3; i++ {
		err = a.infraction(rn)
		if err != nil {
			t.Errorf("Infraction [%v] should not error %v", rn, err)
		}
	}

	b = a.isJailedFor(rn)
	if b == false {
		t.Errorf("isJailedFor should be true instead [%v]", b)
	}

	b = a.hasJails()
	if b == false {
		t.Errorf("hasJails should be true instead [%v]", b)
	}

	// sleep, quit, should NOT be jailed
	time.Sleep(dur)

	b = a.isJailedFor(rn)
	if b == true {
		t.Errorf("isJailedFor should be false instead [%v]", b)
	}

	b = a.hasJails()
	if b == true {
		t.Errorf("hasJails should be false instead [%v]", b)
	}

}

func TestActorInfractionExists(t *testing.T) {
	// setup
	var b bool
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
		ExpireBase:  time.Second * 1,
		Sentence:    time.Millisecond * 50,
	}
	a := NewActor(an, d)
	a.Run()

	b = a.infractionExists(rn)
	if b == true {
		t.Errorf("infractionExists should be false instead [%v]", b)
	}

	inf := NewInfraction(r)
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
	d := NewDirector()
	d.Run()
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
	a := NewClassicActor(an, r, d)
	a.Run()

	for i := 1; i < 3; i++ {
		a.infraction(rn)
	}

	i := a.totalJails()
	if i != 1 {
		t.Errorf("totalJails should be 1 instead [%v]", i)

	}

}
