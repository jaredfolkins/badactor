package badactor

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestNewDirector(t *testing.T) {
	d := NewDirector(ia)
	typ := reflect.TypeOf(d)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
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

	for i := 0; i < 3; i++ {
		err = d.lInfraction(an, rn)
		if err != nil {
			t.Errorf("lIsJailedFor lInfraction should not error [%v]", err)
		}
	}

	b = d.lIsJailed(an)
	if b == false {
		t.Errorf("lIsJailed should be true instead [%v]", b)
	}

	b = d.lIsJailedFor(an, rn)
	if b == false {
		t.Errorf("lIsJailedFor should be true instead [%v]", b)
	}

	// STATE CHANGE
	// sleep to make sure actor is jailed
	time.Sleep(time.Second * 3)

	d.maintenance()

	b = d.lIsJailed(an)
	if b == true {
		t.Errorf("lIsJailed should be true instead [%v]", b)
	}

	b = d.lIsJailedFor(an, rn)
	if b == true {
		t.Errorf("lIsJailedFor should be true instead [%v]", b)
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
