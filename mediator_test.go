package mediatr

import (
	"errors"
	"testing"
)

type FooEvent struct{}
type BarEvent struct{}

func TestReflectMediator_Events(t *testing.T) {
	t.Run("Publish triggers subscriber", func(t *testing.T) {
		mediator := New()

		triggered := false
		mediator.Subscribe(func(event FooEvent) {
			triggered = true
		})

		_ = mediator.Publish(FooEvent{})
		if !triggered {
			t.Fatal("Subscribes is not triggered on event")
		}
	})

	t.Run("Publish not subscribed event", func(t *testing.T) {
		mediator := New()

		triggered := false
		mediator.Subscribe(func(event BarEvent) {
			triggered = true
		})

		_ = mediator.Publish(FooEvent{})
		if triggered {
			t.Fatal("Subscribes was triggered on incorrect event")
		}
	})

	t.Run("Publish return error if handler returns error", func(t *testing.T) {
		mediator := New()

		wantErr := errors.New("bus is busy")
		mediator.Subscribe(func(event BarEvent) error {
			return wantErr
		})

		err := mediator.Publish(BarEvent{})
		if err != wantErr {
			t.Fatal("Publish doesn't return proper error from handler")
		}
	})

	t.Run("Second subscriber returns error", func(t *testing.T) {
		mediator := New()

		wantErr := errors.New("connection refused")
		mediator.Subscribe(func(event BarEvent) error {
			return nil
		})
		mediator.Subscribe(func(event BarEvent) error {
			return wantErr
		})

		err := mediator.Publish(BarEvent{})
		if err != wantErr {
			t.Fatalf("Publish doesn't return proper error from handler %q != %q", err, wantErr)
		}
	})
}

func TestReflectMediator_Commands(t *testing.T) {
	t.Run("Send triggers handler", func(t *testing.T) {
		mediator := New()

		triggered := false
		_ = mediator.Register(func(command FooEvent) {
			triggered = true
		})

		_, _ = mediator.Send(FooEvent{})
		if !triggered {
			t.Fatal("Subscribes is not triggered on event")
		}
	})

	t.Run("Command handler can return error", func(t *testing.T) {
		mediator := New()
		wantErr := errors.New("db error")

		_ = mediator.Register(func(command FooEvent) error {
			return wantErr
		})

		_, err := mediator.Send(FooEvent{})
		if err != wantErr {
			t.Fatal("Send was not receive proper error from handler")
		}
	})

	t.Run("Command handler can return result", func(t *testing.T) {
		mediator := New()

		wantResult := "command result"
		_ = mediator.Register(func(command FooEvent) string {
			return wantResult
		})

		result, _ := mediator.Send(FooEvent{})
		if result != wantResult {
			t.Fatal("Send was not receive proper result from handler")
		}
	})

	t.Run("Command handler can return result and error", func(t *testing.T) {
		mediator := New()

		wantResult := "result"
		_ = mediator.Register(func(command FooEvent) (string, error) {
			return wantResult, nil
		})

		result, _ := mediator.Send(FooEvent{})
		if result != wantResult {
			t.Fatal("Send was not receive proper result from handler")
		}
	})

	t.Run("Command handler can return result and error", func(t *testing.T) {
		mediator := New()

		wantErr := errors.New("two value returns value error")
		_ = mediator.Register(func(command FooEvent) (string, error) {
			return "", wantErr
		})

		_, err := mediator.Send(FooEvent{})
		if err != wantErr {
			t.Fatal("Send was not receive proper error from handler")
		}
	})

	t.Run("Sent to not registered command", func(t *testing.T) {
		mediator := New()

		triggered := false
		_ = mediator.Register(func(event BarEvent) {
			triggered = true
		})

		_, _ = mediator.Send(FooEvent{})
		if triggered {
			t.Fatal("Handler was triggered on incorrect command")
		}
	})

	t.Run("Send command without registration", func(t *testing.T) {
		mediator := New()
		_, err := mediator.Send(FooEvent{})
		if err == nil {
			t.Fatal("Send should return error if command handler is not registered")
		}
	})

	t.Run("Second register is forbidden", func(t *testing.T) {
		mediator := New()
		commandHandler := func(event BarEvent) {}
		_ = mediator.Register(commandHandler)
		err := mediator.Register(commandHandler)
		if err == nil {
			t.Fatal("Command should only have one handler")
		}
	})
}
