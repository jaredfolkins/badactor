package badactor

import "testing"

func TestNewIncoming(t *testing.T) {
	la := "localactor"
	lr := "localrule"
	in := NewIncoming(la, lr, INFRACTION)
	if in.ActorName != la {
		t.Errorf("NewIncoming ActorName expected [%v] was [%v]", la, in.ActorName)
	}

	if in.RuleName != lr {
		t.Errorf("NewIncoming RuleName expected [%v] was [%v]", lr, in.RuleName)
	}

}
