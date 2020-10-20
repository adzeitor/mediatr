package mediatr

import (
	"context"
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
	valueOf := reflect.ValueOf(subscription)
	typeOf := reflect.TypeOf(subscription)
	argKind := typeOf.In(0)

	if typeOf.NumIn() > 1 {
		if argIsContext(argKind) {
			argKind = typeOf.In(1)
		}
	}

	m.subscriptions[argKind] = append(m.subscriptions[argKind], valueOf)
}

// Publish publishes specified domain event to subscribers.
func (m Mediator) Publish(ctx context.Context, event interface{}) error {
	for _, subscription := range m.subscriptions[reflect.TypeOf(event)] {
		arguments := []reflect.Value{
			reflect.ValueOf(event),
		}

		if subscription.Type().NumIn() == 2 {
			arguments = append(
				[]reflect.Value{reflect.ValueOf(ctx)},
				arguments...,
			)
		}

		result := subscription.Call(arguments)
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
	typeOf := reflect.TypeOf(handler)
	argKind := typeOf.In(0)

	if typeOf.NumIn() > 1 {
		if argIsContext(argKind) {
			argKind = typeOf.In(1)
		}
	}

	_, exist := m.registrations[argKind]
	if exist {
		return fmt.Errorf("handler already registered for command %T", argKind)
	}

	m.registrations[argKind] = reflect.ValueOf(handler)
	return nil
}

// Send sent command to handler.
func (m Mediator) Send(ctx context.Context, command interface{}) (interface{}, error) {
	handler, ok := m.registrations[reflect.TypeOf(command)]
	if !ok {
		return nil, fmt.Errorf("no handlers for command %T", command)
	}

	arguments := []reflect.Value{
		reflect.ValueOf(command),
	}

	if handler.Type().NumIn() == 2 {
		arguments = append(
			[]reflect.Value{reflect.ValueOf(ctx)},
			arguments...,
		)
	}

	result := handler.Call(arguments)
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

var contextType = reflect.TypeOf(new(context.Context)).Elem()

func argIsContext(typeOf reflect.Type) bool {
	return contextType == typeOf
}
