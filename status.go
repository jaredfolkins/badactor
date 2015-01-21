package badactor

type status struct {
	outgoing chan *message
}

func newStatus() *status {
	return &status{
		outgoing: make(chan *message),
	}
}
