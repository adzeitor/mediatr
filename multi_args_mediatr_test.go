package mediatr

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestMultiArgsMediatr_PublishSubscriptions(t *testing.T) {
	type Foo struct{ Value string }
	type Bar struct{ Value int }
	type subscriberCall struct {
		Subscriber string
		Foo        Foo
		Bar        Bar
	}
	foo := Foo{Value: "foo"}
	bar := Bar{Value: 42}

	t.Run("Publish matches reversed events to handler arguments", func(t *testing.T) {
		mediator := NewMultiArgsMediatr(New())
		var calls []subscriberCall
		mediator.Subscribe(func(foo Foo, bar Bar) {
			calls = append(calls, subscriberCall{Subscriber: "Foo, Bar", Foo: foo, Bar: bar})
		})

		err := mediator.Publish(context.Background(), bar, foo)

		if err != nil {
			t.Fatal(err)
		}
		want := []subscriberCall{{Subscriber: "Foo, Bar", Foo: foo, Bar: bar}}
		if !reflect.DeepEqual(calls, want) {
			t.Fatalf("got calls %+v, want %+v", calls, want)
		}
	})

	t.Run("Publish calls both subscribers with different argument order", func(t *testing.T) {
		mediator := NewMultiArgsMediatr(New())
		var calls []subscriberCall
		mediator.Subscribe(func(foo Foo, bar Bar) {
			calls = append(calls, subscriberCall{Subscriber: "Foo, Bar", Foo: foo, Bar: bar})
		})
		mediator.Subscribe(func(bar Bar, foo Foo) {
			calls = append(calls, subscriberCall{Subscriber: "Bar, Foo", Foo: foo, Bar: bar})
		})

		err := mediator.Publish(context.Background(), foo, bar)

		if err != nil {
			t.Fatal(err)
		}
		want := []subscriberCall{
			{Subscriber: "Foo, Bar", Foo: foo, Bar: bar},
			{Subscriber: "Bar, Foo", Foo: foo, Bar: bar},
		}
		if !reflect.DeepEqual(calls, want) {
			t.Fatalf("got calls %+v, want %+v", calls, want)
		}
	})

	t.Run("Publish calls both subscribers with reversed events", func(t *testing.T) {
		mediator := NewMultiArgsMediatr(New())
		var calls []subscriberCall
		mediator.Subscribe(func(foo Foo, bar Bar) {
			calls = append(calls, subscriberCall{Subscriber: "Foo, Bar", Foo: foo, Bar: bar})
		})
		mediator.Subscribe(func(bar Bar, foo Foo) {
			calls = append(calls, subscriberCall{Subscriber: "Bar, Foo", Foo: foo, Bar: bar})
		})

		err := mediator.Publish(context.Background(), bar, foo)

		if err != nil {
			t.Fatal(err)
		}
		want := []subscriberCall{
			{Subscriber: "Foo, Bar", Foo: foo, Bar: bar},
			{Subscriber: "Bar, Foo", Foo: foo, Bar: bar},
		}
		if !reflect.DeepEqual(calls, want) {
			t.Fatalf("got calls %+v, want %+v", calls, want)
		}
	})

	t.Run("Publish calls both subscribers when only Foo is published", func(t *testing.T) {
		mediator := NewMultiArgsMediatr(New())
		var calls []subscriberCall
		mediator.Subscribe(func(foo Foo, bar Bar) {
			calls = append(calls, subscriberCall{Subscriber: "Foo, Bar", Foo: foo, Bar: bar})
		})
		mediator.Subscribe(func(bar Bar, foo Foo) {
			calls = append(calls, subscriberCall{Subscriber: "Bar, Foo", Foo: foo, Bar: bar})
		})

		err := mediator.Publish(context.Background(), foo)

		if err != nil {
			t.Fatal(err)
		}
		want := []subscriberCall{
			{Subscriber: "Foo, Bar", Foo: foo, Bar: Bar{}},
			{Subscriber: "Bar, Foo", Foo: foo, Bar: Bar{}},
		}
		if !reflect.DeepEqual(calls, want) {
			t.Fatalf("got calls %+v, want %+v", calls, want)
		}
	})

	t.Run("Publish calls both subscribers when only Bar is published", func(t *testing.T) {
		mediator := NewMultiArgsMediatr(New())
		var calls []subscriberCall
		mediator.Subscribe(func(foo Foo, bar Bar) {
			calls = append(calls, subscriberCall{Subscriber: "Foo, Bar", Foo: foo, Bar: bar})
		})
		mediator.Subscribe(func(bar Bar, foo Foo) {
			calls = append(calls, subscriberCall{Subscriber: "Bar, Foo", Foo: foo, Bar: bar})
		})

		err := mediator.Publish(context.Background(), bar)

		if err != nil {
			t.Fatal(err)
		}
		want := []subscriberCall{
			{Subscriber: "Foo, Bar", Foo: Foo{}, Bar: bar},
			{Subscriber: "Bar, Foo", Foo: Foo{}, Bar: bar},
		}
		if !reflect.DeepEqual(calls, want) {
			t.Fatalf("got calls %+v, want %+v", calls, want)
		}
	})

	t.Run("Publish calls both single-event subscribers", func(t *testing.T) {
		mediator := NewMultiArgsMediatr(New())
		var calls []subscriberCall
		mediator.Subscribe(func(foo Foo) {
			calls = append(calls, subscriberCall{Subscriber: "Foo", Foo: foo})
		})
		mediator.Subscribe(func(bar Bar) {
			calls = append(calls, subscriberCall{Subscriber: "Bar", Bar: bar})
		})

		err := mediator.Publish(context.Background(), foo, bar)

		if err != nil {
			t.Fatal(err)
		}
		want := []subscriberCall{
			{Subscriber: "Foo", Foo: foo},
			{Subscriber: "Bar", Bar: bar},
		}
		if !reflect.DeepEqual(calls, want) {
			t.Fatalf("got calls %+v, want %+v", calls, want)
		}
	})
}

func TestMultiArgsMediatr_Publish(t *testing.T) {
	t.Run("Publish calls handler with multiple events", func(t *testing.T) {
		mediator := NewMultiArgsMediatr(New())
		calls := 0
		mediator.Subscribe(func(event1 string, event2 int) {
			calls++
			if event1 != "some" || event2 != 2 {
				t.Fatalf("got events (%q, %d), want (%q, %d)", event1, event2, "some", 2)
			}
		})

		err := mediator.Publish(context.Background(), "some", 2)

		if err != nil {
			t.Fatal(err)
		}
		if calls != 1 {
			t.Fatalf("got %d calls, want 1", calls)
		}
	})

	t.Run("Publish supplies nil for a missing pointer argument", func(t *testing.T) {
		mediator := NewMultiArgsMediatr(New())
		called := false
		mediator.Subscribe(func(event1 *string, event2 *int) {
			called = true
			if event1 != nil {
				t.Fatalf("got first event %v, want nil", event1)
			}
			if event2 == nil || *event2 != 2 {
				t.Fatalf("got second event %v, want a pointer to 2", event2)
			}
		})
		value := 2

		err := mediator.Publish(context.Background(), &value)

		if err != nil {
			t.Fatal(err)
		}
		if !called {
			t.Fatal("Subscriber is not triggered on events")
		}
	})

	t.Run("Publish passes context to handler", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mediator := NewMultiArgsMediatr(New())
		called := false
		mediator.Subscribe(func(event1 string, contxt context.Context, event2 int) {
			called = true
			if event1 != "some" || event2 != 2 {
				t.Fatalf("got events (%q, %d), want (%q, %d)", event1, event2, "some", 2)
			}
			if contxt != ctx {
				t.Fatal("Subscriber did not receive the publication context")
			}
		})

		err := mediator.Publish(ctx, "some", 2)

		if err != nil {
			t.Fatal(err)
		}
		if !called {
			t.Fatal("Subscriber is not triggered on events")
		}
	})

	t.Run("Publish supplies zero values for missing struct arguments", func(t *testing.T) {
		type firstEvent struct{ Value string }
		type secondEvent struct{ Value int }
		type thirdEvent struct{ Value bool }
		type arguments struct {
			First  firstEvent
			Second secondEvent
			Third  thirdEvent
		}
		tests := []struct {
			name   string
			events []interface{}
			want   arguments
		}{
			{
				name:   "first event only",
				events: []interface{}{firstEvent{Value: "some"}},
				want:   arguments{First: firstEvent{Value: "some"}},
			},
			{
				name:   "second event only",
				events: []interface{}{secondEvent{Value: 2}},
				want:   arguments{Second: secondEvent{Value: 2}},
			},
			{
				name:   "third event only",
				events: []interface{}{thirdEvent{Value: true}},
				want:   arguments{Third: thirdEvent{Value: true}},
			},
			{
				name:   "all events",
				events: []interface{}{thirdEvent{Value: true}, firstEvent{Value: "some"}, secondEvent{Value: 2}},
				want: arguments{
					First:  firstEvent{Value: "some"},
					Second: secondEvent{Value: 2},
					Third:  thirdEvent{Value: true},
				},
			},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				mediator := NewMultiArgsMediatr(New())
				calls := 0
				var got arguments
				mediator.Subscribe(func(first firstEvent, second secondEvent, third thirdEvent) {
					calls++
					got = arguments{First: first, Second: second, Third: third}
				})

				err := mediator.Publish(context.Background(), test.events...)

				if err != nil {
					t.Fatal(err)
				}
				if calls != 1 || got != test.want {
					t.Fatalf("got %d calls with %+v, want 1 with %+v", calls, got, test.want)
				}
			})
		}
	})

	t.Run("Publish does not call subscribers for unrelated events or empty batches", func(t *testing.T) {
		mediator := NewMultiArgsMediatr(New())
		calls := 0
		mediator.Subscribe(func(string, int) { calls++ })

		unrelatedErr := mediator.Publish(context.Background(), true)
		emptyErr := mediator.Publish(context.Background())

		if unrelatedErr != nil || emptyErr != nil {
			t.Fatalf("unexpected publication errors: %v, %v", unrelatedErr, emptyErr)
		}
		if calls != 0 {
			t.Fatalf("got %d calls, want 0", calls)
		}
	})

	t.Run("Publish calls subscribers in subscription order", func(t *testing.T) {
		mediator := NewMultiArgsMediatr(New())
		var calls []string
		mediator.Subscribe(func(FooEvent) { calls = append(calls, "first") })
		mediator.Subscribe(func(FooEvent, BarEvent) error {
			calls = append(calls, "combined")
			return nil
		})
		mediator.Subscribe(func(BarEvent) { calls = append(calls, "last") })
		want := []string{"first", "combined", "last"}

		err := mediator.Publish(context.Background(), BarEvent{}, FooEvent{})

		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(calls, want) {
			t.Fatalf("got call order %v, want %v", calls, want)
		}
	})

	t.Run("Publish stops at the first subscriber error", func(t *testing.T) {
		mediator := NewMultiArgsMediatr(New())
		wantErr := errors.New("handler failed")
		var calls []string
		mediator.Subscribe(func(FooEvent) error {
			calls = append(calls, "first")
			return nil
		})
		mediator.Subscribe(func(FooEvent, BarEvent) error {
			calls = append(calls, "combined")
			return wantErr
		})
		mediator.Subscribe(func(BarEvent) { calls = append(calls, "last") })
		wantCalls := []string{"first", "combined"}

		err := mediator.Publish(context.Background(), BarEvent{}, FooEvent{})

		if err != wantErr {
			t.Fatalf("got error %v, want %v", err, wantErr)
		}
		if !reflect.DeepEqual(calls, wantCalls) {
			t.Fatalf("got calls %v, want %v", calls, wantCalls)
		}
	})
}
