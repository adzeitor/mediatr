package mediatr

import "testing"

func TestMediator_Publish(t *testing.T) {
	mediator := NewMediator()

	t.Run("Publish triggers subscriber", func(t *testing.T) {
		triggered := false
		mediator.Subscribe(FooEvent{}, func(event Event) {
			triggered = true
		})

		mediator.Publish(FooEvent{})
		if !triggered {
			t.Fatal("Subscribes is not triggered on event")
		}
	})
	t.Run("Publish not subscribed event", func(t *testing.T) {
		triggered := false
		mediator.Subscribe(BarEvent{}, func(event Event) {
			triggered = true
		})

		mediator.Publish(FooEvent{})
		if triggered {
			t.Fatal("Subscribes was triggered on incorrect event")
		}
	})
}
