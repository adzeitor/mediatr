package mediatr

import (
	"fmt"
	"reflect"
)

type ReflectMediator struct {
	subscriptions map[reflect.Type][]reflect.Value
	registrations map[reflect.Type]reflect.Value
}

func NewReflectMediator() ReflectMediator {
	return ReflectMediator{
		subscriptions: make(map[reflect.Type][]reflect.Value),
		registrations: make(map[reflect.Type]reflect.Value),
	}
}

func (m ReflectMediator) Subscribe(subscription interface{}) {
	fn := reflect.ValueOf(subscription)
	argKind := reflect.TypeOf(subscription).In(0)
	m.subscriptions[argKind] = append(m.subscriptions[argKind], fn)
}

func (m ReflectMediator) Publish(event interface{}) {
	for _, subscription := range m.subscriptions[reflect.TypeOf(event)] {
		subscription.Call([]reflect.Value{reflect.ValueOf(event)})
	}
}

func (m ReflectMediator) Register(handler interface{}) error {
	argKind := reflect.TypeOf(handler).In(0)
	_, exist := m.registrations[argKind]
	if exist {
		return fmt.Errorf("Handler already registered for command %T", argKind)
	}

	m.registrations[argKind] = reflect.ValueOf(handler)
	return nil
}

func (m ReflectMediator) Send(command interface{}) error {
	handler, ok := m.registrations[reflect.TypeOf(command)]
	if !ok {
		return fmt.Errorf("No handlers for command %T", command)
	}

	handler.Call([]reflect.Value{reflect.ValueOf(command)})
	return nil
}
