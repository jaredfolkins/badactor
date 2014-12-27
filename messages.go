package badactor

type Incoming struct {
	ActorName  string
	RuleName   string
	Type       int
	Infraction *Infraction
	Outgoing   chan *Outgoing
}

type Outgoing struct {
	Message string
	Error   error
}

func NewIncoming(an string, rn string, t int) *Incoming {
	return &Incoming{
		ActorName: an,
		RuleName:  rn,
		Type:      t,
		Outgoing:  make(chan *Outgoing),
	}
}
