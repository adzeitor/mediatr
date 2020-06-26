package mediatr

type Event interface {
	Name() string
}

type subscription func(Event)

type Mediator struct {
	subscriptions map[string][]subscription
}

func NewMediator() Mediator {
	return Mediator{subscriptions: make(map[string][]subscription)}
}

func (m Mediator) Subscribe(event Event, subscription subscription) {
	m.subscriptions[event.Name()] = append(m.subscriptions[event.Name()], subscription)
}

func (m Mediator) Publish(event Event) {
	for _, subscription := range m.subscriptions[event.Name()] {
		subscription(event)
	}
}
