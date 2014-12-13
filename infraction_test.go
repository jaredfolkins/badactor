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

	if inf.Strikes != 1 {
		t.Errorf("Infraction.Strikes should be [%v]", 1)
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
	ot := inf.ExpireBy
	inf.Rebase()
	if !inf.ExpireBy.After(ot) {
		t.Errorf("Infraction.ExpireBy should be new, greater value, instead [%v:%v]", inf.ExpireBy, ot)
	}

}
