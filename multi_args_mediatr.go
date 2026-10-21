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
// the events they subscribe to has arrived. Events are combined in the order
// they are published: an event that comes earlier in the handler signature
// than the previous one, or repeats a filled argument, starts a new call.
func (m MultiArgsMediatr) Publish(ctx context.Context, events ...interface{}) error {
	handlerForCall := make(map[reflect.Value][][]reflect.Value, 0)
	for _, subscription := range *m.orderSubscriptions {
		callArgumentLists := buildCallArgumentLists(ctx, subscription, events)
		if callArgumentLists != nil {
			handlerForCall[subscription.Handler] = callArgumentLists
		}
	}

	return m.callByOrder(handlerForCall)
}

// buildCallArgumentLists builds the argument lists a single handler should be
// called with for a batch of published events. It returns nil when no
// published event matches any of the handler's argument types.
func buildCallArgumentLists(
	ctx context.Context,
	subscription handlerSignature,
	events []interface{},
) [][]reflect.Value {
	var callArgumentLists [][]reflect.Value
	var callArguments []reflect.Value
	lastArgIndex := -1

	for _, event := range events {
		argIndex := -1
		for index, argType := range subscription.Args {
			if argIsContext(argType) || reflect.TypeOf(event) != argType {
				continue
			}
			if argIndex < 0 {
				argIndex = index
			}
			if index > lastArgIndex {
				argIndex = index
				break
			}
		}
		if argIndex < 0 {
			continue
		}

		if callArguments == nil || argIndex <= lastArgIndex {
			callArguments = make([]reflect.Value, len(subscription.Args))
			for index, argType := range subscription.Args {
				if argIsContext(argType) {
					callArguments[index] = reflect.ValueOf(ctx)
				} else {
					callArguments[index] = reflect.Zero(argType)
				}
			}
			callArgumentLists = append(callArgumentLists, callArguments)
		}

		callArguments[argIndex] = reflect.ValueOf(event)
		lastArgIndex = argIndex
	}

	return callArgumentLists
}

// callByOrder calls handlers by order given in Subscribe.
func (m MultiArgsMediatr) callByOrder(
	handlerForCall map[reflect.Value][][]reflect.Value,
) error {
	for _, handler := range *m.orderSubscriptions {
		argumentLists, exists := handlerForCall[handler.Handler]
		if !exists {
			continue
		}
		for _, arguments := range argumentLists {
			result := handler.Handler.Call(arguments)
			if len(result) == 0 || result[0].IsNil() {
				continue
			}
			return result[0].Interface().(error)
		}
	}

	return nil
}
