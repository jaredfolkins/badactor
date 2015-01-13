package badactor

import (
	"strconv"
	"testing"
	"time"
)

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

	st.AddRule(r)

	_, ok := st.rules[r.Name]
	if ok == false {
		t.Errorf("AddRule for Actor [%s] should not fail", an)
	}
}
