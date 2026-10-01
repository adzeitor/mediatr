package mediatr

import (
	"context"
	"reflect"
)

type handlerSignature struct {
	Handler reflect.Value
	Args    []reflect.Type
}

// MultiArgsMediatr represents mediator for domain events.
type MultiArgsMediatr struct {
	Mediator

	subscriptions      map[reflect.Type][]handlerSignature
	orderSubscriptions *[]handlerSignature
}

// NewMultiArgsMediatr is a constructor.
func NewMultiArgsMediatr(mediator Mediator) MultiArgsMediatr {
	orderSubscriptions := make([]handlerSignature, 0)
	return MultiArgsMediatr{
		Mediator:           mediator,
		subscriptions:      make(map[reflect.Type][]handlerSignature),
		orderSubscriptions: &orderSubscriptions,
	}
}

// Subscribe add subscription for domain events.
// Type of events is detected by arguments of handler.
func (m MultiArgsMediatr) Subscribe(subscription interface{}) {
	valueOf := reflect.ValueOf(subscription)
	typeOf := reflect.TypeOf(subscription)

	var args []reflect.Type
	for i := 0; i < typeOf.NumIn(); i++ {
		args = append(args, typeOf.In(i))
	}

	signature := handlerSignature{
		Handler: valueOf,
		Args:    args,
	}

	for _, arg := range args {
		if !argIsContext(arg) {
			m.subscriptions[arg] = append(m.subscriptions[arg], signature)
		}
	}

	*m.orderSubscriptions = append(*m.orderSubscriptions, signature)
}

// Publish publishes specified domain events to subscribers.
// This method calls handlers that have only one event as an argument.
// It also calls handlers with multiple event arguments if at least one of
// the events they subscribe to has arrived.
func (m MultiArgsMediatr) Publish(ctx context.Context, events ...interface{}) error {
	handlerForCall := make(map[reflect.Value][]reflect.Value)

	for _, event := range events {
		subscriptions := m.subscriptions[reflect.TypeOf(event)]

		for _, subscription := range subscriptions {
			arguments := make([]reflect.Value, 0, len(subscription.Args))

			for _, arg := range subscription.Args {
				if argIsContext(arg) {
					arguments = append(arguments, reflect.ValueOf(ctx))
					continue
				}
				eventFound := false
				for _, event := range events {
					eventType := reflect.TypeOf(event)
					if eventType == arg {
						arguments = append(arguments, reflect.ValueOf(event))
						eventFound = true
						break
					}
				}
				if !eventFound {
					arguments = append(arguments, reflect.Zero(arg))
				}
			}

			handlerForCall[subscription.Handler] = arguments
		}
	}

	return m.callByOrder(handlerForCall)
}

// callByOrder calls handlers by order given in Subscribe.
func (m MultiArgsMediatr) callByOrder(
	handlerForCall map[reflect.Value][]reflect.Value,
) error {
	for _, handler := range *m.orderSubscriptions {
		arguments, exists := handlerForCall[handler.Handler]
		if !exists {
			continue
		}
		result := handler.Handler.Call(arguments)
		if len(result) == 0 || result[0].IsNil() {
			continue
		}
		return result[0].Interface().(error)
	}

	return nil
}
