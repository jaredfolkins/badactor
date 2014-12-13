package badactor

type Incoming struct {
	ActorName  string
	RuleName   string
	Type       string
	Infraction *Infraction
	Outgoing   chan *Outgoing
}

type Outgoing struct {
	Message string
	Error   error
}

func NewIncoming(an string, rn string, t string) *Incoming {
	return &Incoming{
		ActorName: an,
		RuleName:  rn,
		Type:      t,
		Outgoing:  make(chan *Outgoing),
	}
}
