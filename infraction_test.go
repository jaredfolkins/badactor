package badactor

import (
	"strconv"
	"testing"
	"time"
)

func TestNewInfraction(t *testing.T) {
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	sl := 3
	eb := time.Second * 60
	s := time.Second * 60
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: sl,
		ExpireBase:  eb,
		Sentence:    s,
	}

	inf := NewInfraction(r)

	if inf.strikes != 0 {
		t.Errorf("Infraction.Strikes should be [%v]", 0)
	}

}

func TestInfractionRebase(t *testing.T) {
	rn := "rn_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	rm := "rm_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	sl := 3
	eb := time.Second * 60
	s := time.Second * 60
	r := &Rule{
		Name:        rn,
		Message:     rm,
		StrikeLimit: sl,
		ExpireBase:  eb,
		Sentence:    s,
	}

	inf := NewInfraction(r)
	ot := inf.expireBy
	inf.rebase()
	if !inf.expireBy.After(ot) {
		t.Errorf("Infraction.ExpireBy should be new, greater value, instead [%v:%v]", inf.expireBy, ot)
	}

}
