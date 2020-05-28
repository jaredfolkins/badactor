package badactor

import (
	"strconv"
	"testing"
	"time"
)

const ia = 1024

type ActorMockAction struct{}

func (ama *ActorMockAction) WhenJailed(a *Actor, r *Rule) error {
	return nil
}
func (ama *ActorMockAction) WhenTimeServed(a *Actor, r *Rule) error {
	return nil
}

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
		Action:      &ActorMockAction{},
	}
	a := newClassicActor(an, r, d)

	//test
	for i := 0; i < 3; i++ {
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

func TestActorJail(t *testing.T) {
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
		Action:      &ActorMockAction{},
	}
	a := newClassicActor(an, r, d)

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
	d := NewDirector(ia)
	sl := 3
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: sl,
		ExpireBase:  time.Second * 60,
		Sentence:    time.Millisecond * 50,
		Action:      &ActorMockAction{},
	}
	a := newClassicActor(an, r, d)

	b = a.isJailedFor(rn)
	if b == true {
		t.Errorf("isJailedFor should be false instead [%v]", b)
	}

	// test
	for i := 0; i < 3; i++ {
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

	time.Sleep(time.Millisecond * 40)
	for _, s := range a.jails {
		a.timeServed(s)
	}
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
	d := NewDirector(ia)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Millisecond * 50,
		Sentence:    time.Millisecond * 50,
		Action:      &ActorMockAction{},
	}

	a := newClassicActor(an, r, d)

	err = a.expire(br)
	if err == nil {
		t.Errorf("Expire should error for [%v]", err)
	}

	a.infraction(rn)

	err = a.expire(br)
	if err == nil {
		t.Errorf("Expire should error for [%v]", err)
	}

	time.Sleep(dur)
	time.Sleep(dur)

	err = a.expire(rn)
	if err != nil {
		t.Errorf("Expire should delete [%v:%v]", an, rn)
	}

	if a.hasInfractions() {
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
		Action:      &ActorMockAction{},
	}

	a := newClassicActor(an, r, d)

	rde := "ruledoesntexist"
	err = a.infraction(rde)
	if err == nil {
		t.Errorf("Infraction should error [%v:%v:%v]", err, an, rde)
	}

	for i := 0; i < 3; i++ {
		err = a.infraction(rn)
		if err != nil {
			t.Errorf("Infraction should not error [%v:%v]", an, rn)
		}
	}

	err = a.jail(rn)
	if err == nil {
		t.Errorf("Jail() err should not be nil instead [%v]", err)
	}

	err = a.infraction(rn)
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
		Action:      &ActorMockAction{},
	}

	a := newClassicActor(an, r, d)

	inf := newInfraction(r)

	err = a.createInfraction(inf)
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
		Action:      &ActorMockAction{},
	}
	a := newClassicActor(an, r, d)

	inf := newInfraction(r)

	err = a.createInfraction(inf)
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
		Action:      &ActorMockAction{},
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

func TestActorShouldDelete(t *testing.T) {
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
		ExpireBase:  time.Millisecond * 1,
		Sentence:    time.Millisecond * 2,
		Action:      &ActorMockAction{},
	}
	a := newClassicActor(an, r, d)

	// test
	// assert all falsey
	b = a.shouldDelete()
	if b == true {
		t.Errorf("shouldDelete should be false instead [%v]", b)
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
	for i := 0; i < 3; i++ {
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

	for _, s := range a.jails {
		a.timeServed(s)
	}

	b = a.isJailedFor(rn)
	if b == true {
		t.Errorf("isJailedFor should be false instead [%v]", b)
	}

	b = a.hasJails()
	if b == true {
		t.Errorf("hasJails should be false instead [%v]", b)
	}

	b = a.hasInfractions()
	if b == true {
		t.Errorf("hasInfractions should be false instead [%v]", b)
	}

	time.Sleep(dur)

	b = a.shouldDelete()
	if b == false {
		t.Errorf("shouldDelete should be true instead [%v]", b)
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
		Action:      &ActorMockAction{},
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
		Action:      &ActorMockAction{},
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

func TestActorName(t *testing.T) {
	// setup
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
		Action:      &ActorMockAction{},
	}
	a := newClassicActor(an, r, d)

	// test
	name := a.Name()
	if name != an {
		t.Errorf("Name() should be %s instead [%s]", an, name)
	}
}
