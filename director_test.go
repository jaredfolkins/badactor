package badactor

import (
	"strconv"
	"testing"
	"time"
)

func TestDirectorlMaintenance(t *testing.T) {
	var b bool
	var err error
	d := NewDirector(ia)
	an := "an_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	r1n := "r1"
	r2n := "r2"
	r1 := &Rule{
		Name:        r1n,
		Message:     "r1 message",
		StrikeLimit: 3,
		ExpireBase:  time.Minute * 1,
		Sentence:    time.Minute * 5,
	}

	r2 := &Rule{
		Name:        r2n,
		Message:     "r2 message",
		StrikeLimit: 3,
		ExpireBase:  time.Minute * 1,
		Sentence:    time.Minute * 5,
	}

	d.lAddRule(r1)
	d.lAddRule(r2)

	if d.size != 0 {
		t.Errorf("d.size should be 0, instead [%v]", d.size)
	}

	for i := 0; i < 3; i++ {
		err = d.lInfraction(an, r1n)
		if err != nil {
			t.Errorf("lIsJailedFor lInfraction should not error [%v]", err)
		}
	}

	for i := 0; i < 1; i++ {
		err = d.lInfraction(an, r2n)
		if err != nil {
			t.Errorf("lIsJailedFor lInfraction should not error [%v]", err)
		}
	}

	if d.size != 1 {
		t.Errorf("d.size should be 1, instead [%v]", d.size)
	}

	b = d.lIsJailed(an)
	if !b {
		t.Errorf("lIsJailed should be true instead [%v]", b)
	}

	b = d.lIsJailedFor(an, r1n)
	if !b {
		t.Errorf("lIsJailedFor should be true instead [%v]", b)
	}

	// MOCK STATE CHANGE
	// Remove 10 minutes from the world
	// get actor
	a := d.actor(an)
	// set time to be one hour ago
	dur := time.Now().Add(-time.Hour * 1)

	// change time.Time of infraction and jail
	a.jails[r1n].releaseBy = dur
	a.infractions[r2n].expireBy = dur

	// perform maintenance
	d.lMaintenance()

	// change time to live
	a.ttl = dur

	// perform maintenance
	d.lMaintenance()

	b = d.lInfractionExists(an, r2n)
	if b {
		t.Errorf("lInfractionExists should be false instead [%v]", b)
	}

	b = d.lIsJailed(an)
	if b {
		t.Errorf("lIsJailed should be false instead [%v]", b)
	}

	b = d.lIsJailedFor(an, r1n)
	if b {
		t.Errorf("lIsJailedFor should be false instead [%v]", b)
	}

	if d.size != 0 {
		t.Errorf("d.size should be 0, instead [%v]", d.size)
	}

}

func TestActorExists(t *testing.T) {
	var err error
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

	err = d.lAddRule(r)
	if err != nil {
		t.Errorf("lAddRule for Actor [%s] should not fail", an)
	}

	if d.lActorExists(an) {
		t.Errorf("Actor [%s] should not be found", an)
	}

	err = d.lCreateActor(an, rn)
	if err != nil {
		t.Errorf("Actor [%s] should be created %v", an, err)
	}

	if !d.lActorExists(an) {
		t.Errorf("Actor [%s] should be found", an)
	}
}

func TestKeepAlive(t *testing.T) {
	var err error
	d := NewDirector(ia)
	an := "an_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Second * 10,
		Sentence:    time.Second * 10,
	}

	err = d.lAddRule(r)
	if err != nil {
		t.Errorf("lAddRule for Actor [%s] should not fail", an)
	}

	err = d.lKeepAlive(an)
	if err == nil {
		t.Errorf("Keep Alive for Actor [%s] should fail", an)
	}

	err = d.lCreateActor(an, rn)
	if err != nil {
		t.Errorf("lCreateActor for Actor [%s] should not fail", an)
	}

	err = d.lKeepAlive(an)
	if err != nil {
		t.Errorf("Keep Alive should not fail : %v ", err)
	}
}

func TestStrikes(t *testing.T) {
	var i int
	var err error
	d := NewDirector(ia)
	ban := "badname"
	brn := "badrule"
	an := "an_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Minute * 10,
		Sentence:    time.Minute * 10,
	}

	err = d.lAddRule(r)
	if err != nil {
		t.Errorf("lAddRule for Actor [%s] should not fail", an)
	}

	// setup valid lInfraction
	err = d.lInfraction(an, rn)
	if err != nil {
		t.Errorf("Ifraction should not fail : %v ", err)
	}

	i, err = d.lStrikes(ban, brn)
	if err == nil {
		t.Errorf("lStrikes should fail : %v ", err)
	}

	if i != 0 {
		t.Errorf("lStrikes should be [%v]:[%v] ", 0, err)
	}

	i, err = d.lStrikes(an, brn)
	if err == nil {
		t.Errorf("lStrikes should fail : %v ", err)
	}

	if i != 0 {
		t.Errorf("lStrikes should be [%v]:[%v] ", 0, err)
	}

	i, err = d.lStrikes(an, rn)
	if err != nil {
		t.Errorf("lStrikes should not fail : %v ", err)
	}

	if i != 1 {
		t.Errorf("lStrikes should be [%v]:[%v] ", 1, err)
	}

}

func TestInfraction(t *testing.T) {
	var err error
	d := NewDirector(ia)
	ban := "badname"
	brn := "badrule"
	an := "an_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Minute * 10,
		Sentence:    time.Minute * 10,
	}

	err = d.lAddRule(r)
	if err != nil {
		t.Errorf("lAddRule for Actor [%s] should not fail", an)
	}

	err = d.lInfraction(ban, brn)
	if err == nil {
		t.Errorf("Ifraction should fail : %v ", err)
	}

	err = d.lInfraction(an, brn)
	if err == nil {
		t.Errorf("Ifraction should fail : %v ", err)
	}

	err = d.lInfraction(an, rn)
	if err != nil {
		t.Errorf("Ifraction should not fail : %v ", err)
	}

	i, err := d.lStrikes(an, rn)
	if err != nil {
		t.Errorf("lStrikes should not fail : %v ", err)
	}

	if i != 1 {
		t.Errorf("lStrikes return value is %d should equal %d : %v ", i, 1, err)
	}

}

func TestInfractionIncrement(t *testing.T) {
	var err error
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

	err = d.lAddRule(r)
	if err != nil {
		t.Errorf("lAddRule for Actor [%s] should not fail", an)
	}

	err = d.lInfraction(an, rn)
	if err != nil {
		t.Errorf("lInfraction should not be err : %v", err)
	}

	i, err := d.lStrikes(an, rn)
	if err != nil {
		t.Errorf("lStrikes should no fail : %v ", err)
	}

	if i != 1 {
		t.Errorf("lStrikes should return %d instead %d", 1, i)
	}

}

func TestIsJailedFor(t *testing.T) {
	var b bool
	var err error
	expectFalse := false
	expectTrue := true
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

	err = d.lAddRule(r)
	if err != nil {
		t.Errorf("lAddRule for Actor [%s] should not fail", an)
	}

	b = d.lIsJailedFor(an, rn)
	if b != expectFalse {
		t.Errorf("lIsJailedFor() should be [%v] instead %v", expectFalse, b)
	}

	for i := 0; i < 3; i++ {
		err = d.lInfraction(an, rn)
		if err != nil {
			t.Errorf("lIsJailedFor() lInfraction should not error [%v]", err)
		}
	}

	b = d.lIsJailedFor(an, rn)
	if b != expectTrue {
		t.Errorf("lIsJailedFor() should be [%v] instead %v", expectTrue, b)
	}

}

func TestIsJailed(t *testing.T) {
	var b bool
	var err error
	expectFalse := false
	expectTrue := true
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

	err = d.lAddRule(r)
	if err != nil {
		t.Errorf("lAddRule for Actor [%s] should not fail", an)
	}

	b = d.lIsJailed(an)
	if b != expectFalse {
		t.Errorf("lIsJailed should be should be [%v] instead %v", expectFalse, b)
	}

	for i := 0; i < 3; i++ {
		err = d.lInfraction(an, rn)
		if err != nil {
			t.Errorf("lIsJailedFor lInfraction should not error [%v]", err)
		}
	}

	b = d.lIsJailed(an)
	if b != expectTrue {
		t.Errorf("lIsJailed should be should be [%v] instead %v", expectTrue, b)
	}

}

func TestInfractionExists(t *testing.T) {
	var b bool
	var err error
	expectFalse := false
	expectTrue := true
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
	err = d.lAddRule(r)
	if err != nil {
		t.Errorf("lAddRule for Actor [%s] should not fail", an)
	}

	b = d.lInfractionExists(an, rn)
	if b != expectFalse {
		t.Errorf("lInfraction should not exist: expected %v instead %v", expectFalse, b)
	}

	d.lInfraction(an, rn)
	b = d.lInfractionExists(an, rn)
	if b != expectTrue {
		t.Errorf("lInfraction should exist: expected %v instead %v", expectTrue, b)
	}

}

func TestCreateInfraction(t *testing.T) {
	var err error
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
	err = d.lAddRule(r)
	if err != nil {
		t.Errorf("lAddRule for Actor [%s] should not fail", an)
	}

	br := "badrule"
	ba := "badactor"

	err = d.lCreateInfraction(an, br)
	if err == nil {
		t.Errorf("Should error, Rule does not exist: %v", err)
	}

	err = d.lCreateInfraction(ba, rn)
	if err == nil {
		t.Errorf("Should error, Actor does not exist: %v", err)
	}

	err = d.lCreateActor(an, rn)
	if err != nil {
		t.Errorf("Should not error, Actor and Rule exist: %v", err)
	}

	err = d.lCreateInfraction(an, rn)
	if err != nil {
		t.Errorf("Should not error, Actor and Rule exist: %v", err)
	}

}

func TestAddRule(t *testing.T) {
	var err error
	d := NewDirector(ia)

	r := NewClassicRule("PasswordReset", "You have requested a password reset too often")

	err = d.lAddRule(r)
	if err != nil {
		t.Errorf("Should not fail, Rule shouldn't exist: %v", err)
	}

	err = d.lAddRule(r)
	if err == nil {
		t.Errorf("Should fail, Rule does exist: %v", err)
	}

}

func TestDirectorlIsJailedFor(t *testing.T) {
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
		ExpireBase:  time.Second * 60,
		Sentence:    time.Millisecond * 10,
	}
	d.lAddRule(r)

	if d.size != 0 {
		t.Errorf("d.size should be 0, instead [%v]", d.size)
	}

	for i := 0; i < 3; i++ {
		err = d.lInfraction(an, rn)
		if err != nil {
			t.Errorf("lIsJailedFor lInfraction should not error [%v]", err)
		}
	}

	if d.size != 1 {
		t.Errorf("d.size should be 1, instead [%v]", d.size)
	}

	b = d.lIsJailed(an)
	if !b {
		t.Errorf("lIsJailed should be true instead [%v]", b)
	}

	b = d.lIsJailedFor(an, rn)
	if !b {
		t.Errorf("lIsJailedFor should be true instead [%v]", b)
	}

	// MOCK STATE CHANGE
	// Remove 5 minutes from the world
	a := d.actors[an].Value.(*Actor)
	dur := time.Now().Add(-time.Minute * 5)
	a.jails[rn].releaseBy = dur
	a.ttl = dur

	d.lMaintenance()

	b = d.lIsJailed(an)
	if b {
		t.Errorf("lIsJailed should be false instead [%v]", b)
	}

	b = d.lIsJailedFor(an, rn)
	if b {
		t.Errorf("lIsJailedFor should be false instead [%v]", b)
	}

}

func TestCreateActor(t *testing.T) {
	var err error
	d := NewDirector(ia)
	an := "an_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Second * 10,
		Sentence:    time.Second * 10,
	}

	err = d.lCreateActor(an, rn)
	if err == nil {
		t.Errorf("Actor [%s] should be created %v", an, err)
	}

	err = d.lAddRule(r)
	if err != nil {
		t.Errorf("lAddRule for Actor [%s] should not fail", an)
	}

	err = d.lCreateActor(an, rn)
	if err != nil {
		t.Errorf("director.lCreateActor() for [%s] should not fail", an)
	}

	err = d.lCreateActor(an, rn)
	if err == nil {
		t.Errorf("director.lCreateActor() for [%s] should fail", an)
	}
}

func TestDeleteOldest(t *testing.T) {
	var ok bool
	var err error
	var size int32
	var i int64
	size = 10
	d := NewDirector(size)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Minute * 10,
		Sentence:    time.Minute * 10,
	}

	err = d.lAddRule(r)
	if err != nil {
		t.Errorf("lAddRule [%s] should not fail", rn)
	}

	ok = d.isFull()
	if ok == true {
		t.Errorf("isFull should be false", ok)
	}

	for i = 0; i < 20; i++ {
		an := strconv.FormatInt(i, 10)
		d.createActor(an, rn)
	}

	ok = d.isFull()
	if ok == false {
		t.Errorf("isFull should be true", ok)
	}

	d.deleteOldest()

	ok = d.isFull()
	if ok == true {
		t.Errorf("isFull should be false", ok)
	}

}
