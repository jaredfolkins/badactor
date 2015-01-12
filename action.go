package badactor

// Action is the inferface the Programmer implements to perform event based actions
type Action interface {
	WhenJailed(a *Actor, r *Rule) error     // When an Actor isJailed, do this
	WhenTimeServed(a *Actor, r *Rule) error // When an Actor is relased because of timeServed, do this
}
