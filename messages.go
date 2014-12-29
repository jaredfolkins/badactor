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
