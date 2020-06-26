package mediatr

import (
	"fmt"
	"reflect"
)

// Mediator represents mediator for commands and domain events.
type Mediator struct {
	subscriptions map[reflect.Type][]reflect.Value
	registrations map[reflect.Type]reflect.Value
}

// New return new instance of mediator.
func New() Mediator {
	return Mediator{
		subscriptions: make(map[reflect.Type][]reflect.Value),
		registrations: make(map[reflect.Type]reflect.Value),
	}
}

// Subscribe add subscription for domain event.
// Type of event is detected by arguments of handler.
func (m Mediator) Subscribe(subscription interface{}) {
	fn := reflect.ValueOf(subscription)
	argKind := reflect.TypeOf(subscription).In(0)
	m.subscriptions[argKind] = append(m.subscriptions[argKind], fn)
}

// Publish publishes specified domain event to subscribers.
func (m Mediator) Publish(event interface{}) error {
	for _, subscription := range m.subscriptions[reflect.TypeOf(event)] {
		result := subscription.Call([]reflect.Value{reflect.ValueOf(event)})
		if len(result) == 0 || result[0].IsNil() {
			continue
		}
		return result[0].Interface().(error)
	}

	return nil
}

// Register registers command handler.
// Command type is detected by argument of handler.
func (m Mediator) Register(handler interface{}) error {
	argKind := reflect.TypeOf(handler).In(0)
	_, exist := m.registrations[argKind]
	if exist {
		return fmt.Errorf("Handler already registered for command %T", argKind)
	}

	m.registrations[argKind] = reflect.ValueOf(handler)
	return nil
}

// Send sent command to handler.
func (m Mediator) Send(command interface{}) (interface{}, error) {
	handler, ok := m.registrations[reflect.TypeOf(command)]
	if !ok {
		return nil, fmt.Errorf("No handlers for command %T", command)
	}

	result := handler.Call([]reflect.Value{reflect.ValueOf(command)})
	switch len(result) {
	case 0:
		return nil, nil
	case 1:
		return oneReturnValuesCommand(result)
	case 2:
		return twoReturnValuesCommand(result)
	}
	return nil, nil
}

func oneReturnValuesCommand(result []reflect.Value) (interface{}, error) {
	err, isError := result[0].Interface().(error)
	if isError {
		return nil, err
	}
	return result[0].Interface(), nil
}

func twoReturnValuesCommand(result []reflect.Value) (interface{}, error) {
	var err error
	if !result[1].IsNil() {
		err = result[1].Interface().(error)
	}
	return result[0].Interface(), err
}
