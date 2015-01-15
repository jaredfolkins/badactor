package badactor

import (
	"strconv"
	"testing"
	"time"
)

func TestStudioCreateActorAndInfraction(t *testing.T) {
	var err error

	st := NewStudio(256)
	an := "actorname"
	rn := "rulename"
	rm := "rulemessage"
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Second * 10,
		Sentence:    time.Minute * 10,
	}

	// add rule
	st.AddRule(r)

	// creat directors
	err = st.CreateDirectors(1024)
	if err != nil {
		t.Errorf("CreateDirectors failed %v", an, rn)
	}

	if st.ActorExists(an) {
		t.Errorf("ActorExists should be false")
	}

	if st.InfractionExists(an, rn) {
		t.Errorf("InfractionExists should be false")
	}

	err = st.CreateActor(an, rn)
	if err != nil {
		t.Errorf("CreateActor should be nil [%v]", err)
	}

	err = st.CreateInfraction(an, rn)
	if err != nil {
		t.Errorf("CreateInfraction should be nil [%v]", err)
	}

	if !st.ActorExists(an) {
		t.Errorf("InfractionExists should be true")
	}

	if !st.InfractionExists(an, rn) {
		t.Errorf("InfractionExists should be true")
	}
}

func TestStudioStrikes(t *testing.T) {
	var si int
	var err error

	st := NewStudio(256)
	an := "actorname"
	rn := "rulename"
	rm := "rulemessage"
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Second * 10,
		Sentence:    time.Minute * 10,
	}

	// add rule
	st.AddRule(r)

	// creat directors
	err = st.CreateDirectors(1024)
	if err != nil {
		t.Errorf("CreateDirectors failed %v", an, rn)
	}

	si, err = st.Strikes(an, rn)
	if si != 0 {
		t.Errorf("Strikes for Actor [%v] and Rule [%v] should not be %v, %v", an, rn, si, err)
	}

	// 1st inf
	err = st.Infraction(an, rn)
	if err != nil {
		t.Errorf("Infraction failed for Actor [%v] and Rule [%v] should not be %v", an, rn, err)
	}

	si, err = st.Strikes(an, rn)
	if si != 1 {
		t.Errorf("Strikes for Actor [%v] and Rule [%v] should not be %v", an, rn, si, err)
	}

	// 2nd inf
	st.Infraction(an, rn)
	si, err = st.Strikes(an, rn)
	if si != 2 {
		t.Errorf("Strikes for Actor [%v] and Rule [%v] should not be %v", an, rn, si, err)
	}

	// 3rd inf, jail, strikes for that infraction name should be 0
	st.Infraction(an, rn)
	si, err = st.Strikes(an, rn)
	if si != 0 {
		t.Errorf("Strikes for Actor [%v] and Rule [%v] should not be %v", an, rn, si, err)
	}

	// should still be jailed
	st.Infraction(an, rn)
	si, err = st.Strikes(an, rn)
	if si != 0 {
		t.Errorf("Strikes for Actor [%v] and Rule [%v] should not be %v", an, rn, si, err)
	}
}

func TestStudioAddRule(t *testing.T) {
	st := NewStudio(256)
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

	// add rule
	st.AddRule(r)

	// test rule exists
	_, ok := st.rules[r.Name]
	if ok == false {
		t.Errorf("AddRule for Actor [%s] should not fail", an)
	}
}

func TestStudioAddRules(t *testing.T) {
	st := NewStudio(2)
	r1 := &Rule{
		Name:        "rule1",
		Message:     "message1",
		StrikeLimit: 3,
		ExpireBase:  time.Millisecond * 10,
		Sentence:    time.Millisecond * 10,
	}

	r2 := &Rule{
		Name:        "rule2",
		Message:     "message2",
		StrikeLimit: 3,
		ExpireBase:  time.Millisecond * 10,
		Sentence:    time.Millisecond * 10,
	}

	// add rule safety is of no concern
	st.AddRule(r1)
	st.AddRule(r2)

	// apply rules
	err := st.CreateDirectors(1024)
	if err != nil {
		t.Errorf("CreateDirectors failed %v", err)
	}

	// range of rules and directors and make sure the rule exists for each
	for di, d := range st.directors {
		if !d.ruleExists(r1.Name) {
			t.Errorf("ApplyRules for director [%v] is missing rule %v", di, r1.Name)
		}
		if !d.ruleExists(r2.Name) {
			t.Errorf("ApplyRules for director [%v] is missing rule %v", di, r2.Name)
		}
	}

}

func TestStudioCreateDirectors(t *testing.T) {

	st := NewStudio(256)
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: 3,
		ExpireBase:  time.Millisecond * 10,
		Sentence:    time.Millisecond * 10,
	}

	// add rule safety is of no concern
	st.AddRule(r)

	if st.capacity != 256 {
		t.Errorf("Capacity for Studio [%s] should not be 256 instead %v", st.capacity)
	}

	st.CreateDirectors(256)

	var i int32
	for i = 0; i < st.capacity; i++ {
		_, ok := st.directors[i]
		if !ok {
			t.Errorf("Director [%v] for Studio was not created", i)
		}
	}

}
