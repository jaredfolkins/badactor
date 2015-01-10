package badactor

// Action is the inferface the Programmer implements to perform event based actions
type Action interface {
	WhenJailed(s *Stats) error     // When an Actor isJailed, do this
	WhenTimeServed(s *Stats) error // When an Actor is relased because of timeServed, do this
}
