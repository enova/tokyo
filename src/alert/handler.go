package alert

// Handler is the interface for alert-handlers
type Handler interface {
	Handle(msg Message)
}
