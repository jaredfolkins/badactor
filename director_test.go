package badactor

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestNewDirector(t *testing.T) {
	d := NewDirector()
	d.Run()
	typ := reflect.TypeOf(d)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
}

func TestActorExists(t *testing.T) {
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
		ExpireBase:  time.Millisecond * 10,
		Sentence:    time.Millisecond * 10,
	}

	err = d.AddRule(r)
	if err != nil {
		t.Errorf("AddRule for Actor [%s] should not fail", an)
	}

	if d.ActorExists(an) {
		t.Errorf("Actor [%s] should not be found", an)
	}

	err = d.CreateActor(an, rn)
	if err != nil {
		t.Errorf("Actor [%s] should be created %v", an, err)
	}

	if !d.ActorExists(an) {
		t.Errorf("Actor [%s] should be found", an)
	}
}

func TestKeepAlive(t *testing.T) {
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
		ExpireBase:  time.Second * 10,
		Sentence:    time.Second * 10,
	}

	err = d.AddRule(r)
	if err != nil {
		t.Errorf("AddRule for Actor [%s] should not fail", an)
	}

	err = d.KeepAlive(an)
	if err == nil {
		t.Errorf("Keep Alive for Actor [%s] should fail", an)
	}

	err = d.CreateActor(an, rn)
	if err != nil {
		t.Errorf("CreateActor for Actor [%s] should not fail", an)
	}

	err = d.KeepAlive(an)
	if err != nil {
		t.Errorf("Keep Alive should not fail : %v ", err)
	}
}

func TestStrikes(t *testing.T) {
	var i int
	var err error
	d := NewDirector()
	d.Run()
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

	err = d.AddRule(r)
	if err != nil {
		t.Errorf("AddRule for Actor [%s] should not fail", an)
	}

	// setup valid infraction
	err = d.Infraction(an, rn)
	if err != nil {
		t.Errorf("Ifraction should not fail : %v ", err)
	}

	i, err = d.Strikes(ban, brn)
	if err == nil {
		t.Errorf("Strikes should fail : %v ", err)
	}

	if i != 0 {
		t.Errorf("Strikes should be [%v]:[%v] ", 0, err)
	}

	i, err = d.Strikes(an, brn)
	if err == nil {
		t.Errorf("Strikes should fail : %v ", err)
	}

	if i != 0 {
		t.Errorf("Strikes should be [%v]:[%v] ", 0, err)
	}

	i, err = d.Strikes(an, rn)
	if err != nil {
		t.Errorf("Strikes should not fail : %v ", err)
	}

	if i != 1 {
		t.Errorf("Strikes should be [%v]:[%v] ", 1, err)
	}

}

func TestMostCostlyInfraction(t *testing.T) {
	var b bool
	var err error
	d := NewDirector()
	d.Run()
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

	err = d.AddRule(r)
	if err != nil {
		t.Errorf("AddRule for Actor [%s] should not fail", an)
	}

	b, err = d.MostCostlyInfraction(an, brn)
	if err == nil {
		t.Errorf("MostCostlyIfraction should not fail : %v %v", err, b)
	}

	b, err = d.MostCostlyInfraction(ban, rn)
	if err != nil {
		t.Errorf("MostCostlyIfraction should fail : %v  %v", err, b)
	}

}

func TestInfraction(t *testing.T) {
	var err error
	d := NewDirector()
	d.Run()
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

	err = d.AddRule(r)
	if err != nil {
		t.Errorf("AddRule for Actor [%s] should not fail", an)
	}

	err = d.Infraction(ban, brn)
	if err == nil {
		t.Errorf("Ifraction should fail : %v ", err)
	}

	err = d.Infraction(an, brn)
	if err == nil {
		t.Errorf("Ifraction should fail : %v ", err)
	}

	err = d.Infraction(an, rn)
	if err != nil {
		t.Errorf("Ifraction should not fail : %v ", err)
	}

	i, err := d.Strikes(an, rn)
	if err != nil {
		t.Errorf("Strikes should not fail : %v ", err)
	}

	if i != 1 {
		t.Errorf("Strikes return value is %d should equal %d : %v ", i, 1, err)
	}

}

func TestInfractionIncrement(t *testing.T) {
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
		ExpireBase:  time.Millisecond * 10,
		Sentence:    time.Millisecond * 10,
	}

	err = d.AddRule(r)
	if err != nil {
		t.Errorf("AddRule for Actor [%s] should not fail", an)
	}

	err = d.Infraction(an, rn)
	if err != nil {
		t.Errorf("Infraction should not be err : %v", err)
	}

	i, err := d.Strikes(an, rn)
	if err != nil {
		t.Errorf("Strikes should no fail : %v ", err)
	}

	if i != 1 {
		t.Errorf("Strikes should return %d instead %d", 1, i)
	}

}

func TestIsJailedFor(t *testing.T) {
	var b bool
	var err error
	expectFalse := false
	expectTrue := true
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

	err = d.AddRule(r)
	if err != nil {
		t.Errorf("AddRule for Actor [%s] should not fail", an)
	}

	b = d.IsJailedFor(an, rn)
	if b != expectFalse {
		t.Errorf("IsJailedFor() should be [%v] instead %v", expectFalse, b)
	}

	for i := 0; i < 3; i++ {
		err = d.Infraction(an, rn)
		if err != nil {
			t.Errorf("IsJailedFor() Infraction should not error [%v]", err)
		}
	}

	b = d.IsJailedFor(an, rn)
	if b != expectTrue {
		t.Errorf("IsJailedFor() should be [%v] instead %v", expectTrue, b)
	}

}

func TestIsJailed(t *testing.T) {
	var b bool
	var err error
	expectFalse := false
	expectTrue := true
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

	err = d.AddRule(r)
	if err != nil {
		t.Errorf("AddRule for Actor [%s] should not fail", an)
	}

	b = d.IsJailed(an)
	if b != expectFalse {
		t.Errorf("IsJailed should be should be [%v] instead %v", expectFalse, b)
	}

	for i := 0; i < 3; i++ {
		err = d.Infraction(an, rn)
		if err != nil {
			t.Errorf("IsJailedFor Infraction should not error [%v]", err)
		}
	}

	b = d.IsJailed(an)
	if b != expectTrue {
		t.Errorf("IsJailed should be should be [%v] instead %v", expectTrue, b)
	}

}

func TestInfractionExists(t *testing.T) {
	var b bool
	var err error
	expectFalse := false
	expectTrue := true
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
	err = d.AddRule(r)
	if err != nil {
		t.Errorf("AddRule for Actor [%s] should not fail", an)
	}

	b = d.InfractionExists(an, rn)
	if b != expectFalse {
		t.Errorf("Infraction should not exist: expected %v instead %v", expectFalse, b)
	}

	d.Infraction(an, rn)
	b = d.InfractionExists(an, rn)
	if b != expectTrue {
		t.Errorf("Infraction should exist: expected %v instead %v", expectTrue, b)
	}

}

func TestCreateInfraction(t *testing.T) {
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
		ExpireBase:  time.Millisecond * 10,
		Sentence:    time.Millisecond * 10,
	}
	err = d.AddRule(r)
	if err != nil {
		t.Errorf("AddRule for Actor [%s] should not fail", an)
	}

	br := "badrule"
	ba := "badactor"

	err = d.CreateInfraction(an, br)
	if err == nil {
		t.Errorf("Should error, Rule does not exist: %v", err)
	}

	err = d.CreateInfraction(ba, rn)
	if err == nil {
		t.Errorf("Should error, Actor does not exist: %v", err)
	}

	err = d.CreateActor(an, rn)
	if err != nil {
		t.Errorf("Should not error, Actor and Rule exist: %v", err)
	}

	err = d.CreateInfraction(an, rn)
	if err != nil {
		t.Errorf("Should not error, Actor and Rule exist: %v", err)
	}

}

func TestAddRule(t *testing.T) {
	var err error
	d := NewDirector()
	d.Run()

	r := NewClassicRule("PasswordReset", "You have requested a password reset too often")

	err = d.AddRule(r)
	if err != nil {
		t.Errorf("Should not fail, Rule shouldn't exist: %v", err)
	}

	err = d.AddRule(r)
	if err == nil {
		t.Errorf("Should fail, Rule does exist: %v", err)
	}

}

func TestDirectorIsJailedFor(t *testing.T) {
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
		ExpireBase:  time.Second * 60,
		Sentence:    time.Millisecond * 10,
	}
	d.AddRule(r)

	for i := 0; i < 3; i++ {
		err = d.Infraction(an, rn)
		if err != nil {
			t.Errorf("IsJailedFor Infraction should not error [%v]", err)
		}
	}

	b = d.IsJailed(an)
	if b == false {
		t.Errorf("IsJailed should be true instead [%v]", b)
	}

	b = d.IsJailedFor(an, rn)
	if b == false {
		t.Errorf("IsJailedFor should be true instead [%v]", b)
	}

	// STATE CHANGE
	// sleep to make sure actor is jailed
	time.Sleep(time.Second * 3)

	b = d.IsJailed(an)
	if b == true {
		t.Errorf("IsJailed should be true instead [%v]", b)
	}

	b = d.IsJailedFor(an, rn)
	if b == true {
		t.Errorf("IsJailedFor should be true instead [%v]", b)
	}

}

func TestCreateActor(t *testing.T) {
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
		ExpireBase:  time.Second * 10,
		Sentence:    time.Second * 10,
	}

	err = d.AddRule(r)
	if err != nil {
		t.Errorf("AddRule for Actor [%s] should not fail", an)
	}

	err = d.CreateActor(an, rn)
	if err != nil {
		t.Errorf("director.CreateActor() for [%s] should not fail", an)
	}

	err = d.CreateActor(an, rn)
	if err == nil {
		t.Errorf("director.CreateActor() for [%s] should fail", an)
	}
}
