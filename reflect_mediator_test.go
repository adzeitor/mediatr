package mediatr

import "testing"

func TestReflectMediator_Events(t *testing.T) {
	mediator := NewReflectMediator()

	t.Run("Publish triggers subscriber", func(t *testing.T) {
		triggered := false
		mediator.Subscribe(func(event FooEvent) {
			triggered = true
		})

		mediator.Publish(FooEvent{})
		if !triggered {
			t.Fatal("Subscribes is not triggered on event")
		}
	})
	t.Run("Publish not subscribed event", func(t *testing.T) {
		triggered := false
		mediator.Subscribe(func(event BarEvent) {
			triggered = true
		})

		mediator.Publish(FooEvent{})
		if triggered {
			t.Fatal("Subscribes was triggered on incorrect event")
		}
	})
}

func TestReflectMediator_Commands(t *testing.T) {
	t.Run("Send triggers handler", func(t *testing.T) {
		mediator := NewReflectMediator()

		triggered := false
		_ = mediator.Register(func(command FooEvent) {
			triggered = true
		})

		_ = mediator.Send(FooEvent{})
		if !triggered {
			t.Fatal("Subscribes is not triggered on event")
		}
	})

	t.Run("Publish not subscribed event", func(t *testing.T) {
		mediator := NewReflectMediator()

		triggered := false
		_ = mediator.Register(func(event BarEvent) {
			triggered = true
		})

		_ = mediator.Send(FooEvent{})
		if triggered {
			t.Fatal("Handler was triggered on incorrect event")
		}
	})

	t.Run("Send command without registration", func(t *testing.T) {
		mediator := NewReflectMediator()
		err := mediator.Send(FooEvent{})
		if err == nil {
			t.Fatal("Send should return error if command handler is not registered")
		}
	})

	t.Run("Send command without registration", func(t *testing.T) {
		mediator := NewReflectMediator()
		_ = mediator.Register(FooEventTypedHandler)
		err := mediator.Register(FooEventTypedHandler)
		if err == nil {
			t.Fatal("Command should only have one handler")
		}
	})
}
