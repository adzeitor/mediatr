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
		Bar2       Bar
	}
	foo := Foo{Value: "foo"}
	bar := Bar{Value: 42}
	bar2 := Bar{Value: 1142}

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
		want := []subscriberCall{
			{Subscriber: "Foo, Bar", Bar: bar},
			{Subscriber: "Foo, Bar", Foo: foo},
		}
		if !reflect.DeepEqual(calls, want) {
			t.Fatalf("got calls %+v, want %+v", calls, want)
		}
	})

	t.Run("Publish starts another call for a repeated argument", func(t *testing.T) {
		mediator := NewMultiArgsMediatr(New())
		var calls []subscriberCall
		mediator.Subscribe(func(foo Foo, bar Bar) {
			calls = append(calls, subscriberCall{Subscriber: "Foo, Bar", Foo: foo, Bar: bar})
		})

		err := mediator.Publish(context.Background(), foo, bar, bar2)

		if err != nil {
			t.Fatal(err)
		}
		want := []subscriberCall{
			{Subscriber: "Foo, Bar", Foo: foo, Bar: bar},
			{Subscriber: "Foo, Bar", Foo: Foo{}, Bar: bar2},
		}
		if !reflect.DeepEqual(calls, want) {
			t.Fatalf("got calls %+v, want %+v", calls, want)
		}
	})

	t.Run("Publish supplies zero for a missing repeated argument", func(t *testing.T) {
		mediator := NewMultiArgsMediatr(New())
		var calls []subscriberCall
		mediator.Subscribe(func(foo Foo, bar Bar, bar2 Bar) {
			calls = append(calls, subscriberCall{Subscriber: "Foo, Bar, Bar", Foo: foo, Bar: bar, Bar2: bar2})
		})

		err := mediator.Publish(context.Background(), foo, bar)

		if err != nil {
			t.Fatal(err)
		}
		want := []subscriberCall{
			{Subscriber: "Foo, Bar, Bar", Foo: foo, Bar: bar, Bar2: Bar{}},
		}
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
			{Subscriber: "Bar, Foo", Foo: foo},
			{Subscriber: "Bar, Foo", Bar: bar},
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
			{Subscriber: "Foo, Bar", Bar: bar},
			{Subscriber: "Foo, Bar", Foo: foo},
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
			want   []arguments
		}{
			{
				name:   "first event only",
				events: []interface{}{firstEvent{Value: "some"}},
				want:   []arguments{{First: firstEvent{Value: "some"}}},
			},
			{
				name:   "second event only",
				events: []interface{}{secondEvent{Value: 2}},
				want:   []arguments{{Second: secondEvent{Value: 2}}},
			},
			{
				name:   "third event only",
				events: []interface{}{thirdEvent{Value: true}},
				want:   []arguments{{Third: thirdEvent{Value: true}}},
			},
			{
				name:   "all events",
				events: []interface{}{thirdEvent{Value: true}, firstEvent{Value: "some"}, secondEvent{Value: 2}},
				want: []arguments{
					{Third: thirdEvent{Value: true}},
					{
						First:  firstEvent{Value: "some"},
						Second: secondEvent{Value: 2},
					},
				},
			},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				mediator := NewMultiArgsMediatr(New())
				var calls []arguments
				mediator.Subscribe(func(first firstEvent, second secondEvent, third thirdEvent) {
					calls = append(calls, arguments{First: first, Second: second, Third: third})
				})

				err := mediator.Publish(context.Background(), test.events...)

				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(calls, test.want) {
					t.Fatalf("got calls %+v, want %+v", calls, test.want)
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
		want := []string{"first", "combined", "combined", "last"}

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

func TestMultiArgsMediatr_PublishEventOrder(t *testing.T) {
	type FirstEvent struct{ Value int }
	type SecondEvent struct{ Value int }
	type subscriberCall struct {
		First  FirstEvent
		Second SecondEvent
	}

	tests := []struct {
		name             string
		reverseArguments bool
		events           []interface{}
		want             []subscriberCall
	}{
		{
			name:   "Publish handles alternating event types in argument order",
			events: []interface{}{FirstEvent{Value: 1}, SecondEvent{Value: 1}, FirstEvent{Value: 2}},
			want: []subscriberCall{
				{First: FirstEvent{Value: 1}, Second: SecondEvent{Value: 1}},
				{First: FirstEvent{Value: 2}},
			},
		},
		{
			name: "Publish handles consecutive events of the same type",
			events: []interface{}{
				FirstEvent{Value: 1}, FirstEvent{Value: 2},
				SecondEvent{Value: 1}, SecondEvent{Value: 2},
			},
			want: []subscriberCall{
				{First: FirstEvent{Value: 1}},
				{First: FirstEvent{Value: 2}, Second: SecondEvent{Value: 1}},
				{Second: SecondEvent{Value: 2}},
			},
		},
		{
			name:             "Publish preserves event order with reversed handler arguments",
			reverseArguments: true,
			events:           []interface{}{FirstEvent{Value: 1}, SecondEvent{Value: 1}, FirstEvent{Value: 2}},
			want: []subscriberCall{
				{First: FirstEvent{Value: 1}},
				{Second: SecondEvent{Value: 1}, First: FirstEvent{Value: 2}},
			},
		},
		{
			name:             "Publish handles consecutive events with reversed handler arguments",
			reverseArguments: true,
			events: []interface{}{
				FirstEvent{Value: 1}, FirstEvent{Value: 2},
				SecondEvent{Value: 1}, SecondEvent{Value: 2},
			},
			want: []subscriberCall{
				{First: FirstEvent{Value: 1}},
				{First: FirstEvent{Value: 2}},
				{Second: SecondEvent{Value: 1}},
				{Second: SecondEvent{Value: 2}},
			},
		},
		{
			name:   "Publish starts separate calls when events reverse argument order",
			events: []interface{}{SecondEvent{Value: 1}, FirstEvent{Value: 1}},
			want: []subscriberCall{
				{Second: SecondEvent{Value: 1}},
				{First: FirstEvent{Value: 1}},
			},
		},
		{
			name:   "Publish delivers zero-valued events",
			events: []interface{}{FirstEvent{}, SecondEvent{}},
			want:   []subscriberCall{{}},
		},
		{
			name:   "Publish ignores unrelated events without splitting a call",
			events: []interface{}{FirstEvent{Value: 1}, "unrelated", SecondEvent{Value: 1}},
			want: []subscriberCall{
				{First: FirstEvent{Value: 1}, Second: SecondEvent{Value: 1}},
			},
		},
		{
			name: "Publish delivers equal events across calls in the same batch",
			events: []interface{}{
				FirstEvent{Value: 1}, FirstEvent{Value: 2},
				FirstEvent{Value: 1}, SecondEvent{Value: 1},
			},
			want: []subscriberCall{
				{First: FirstEvent{Value: 1}},
				{First: FirstEvent{Value: 2}},
				{First: FirstEvent{Value: 1}, Second: SecondEvent{Value: 1}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mediator := NewMultiArgsMediatr(New())
			var calls []subscriberCall
			if test.reverseArguments {
				mediator.Subscribe(func(secondEvent SecondEvent, firstEvent FirstEvent) {
					calls = append(calls, subscriberCall{First: firstEvent, Second: secondEvent})
				})
			} else {
				mediator.Subscribe(func(firstEvent FirstEvent, secondEvent SecondEvent) {
					calls = append(calls, subscriberCall{First: firstEvent, Second: secondEvent})
				})
			}

			err := mediator.Publish(context.Background(), test.events...)

			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(calls, test.want) {
				t.Fatalf("got calls %+v, want %+v", calls, test.want)
			}
		})
	}

	t.Run("Publish delivers every distinct event to a single-event handler", func(t *testing.T) {
		mediator := NewMultiArgsMediatr(New())
		var calls []SecondEvent
		mediator.Subscribe(func(event SecondEvent) { calls = append(calls, event) })

		err := mediator.Publish(context.Background(), SecondEvent{Value: 1}, SecondEvent{Value: 2})

		if err != nil {
			t.Fatal(err)
		}
		want := []SecondEvent{{Value: 1}, {Value: 2}}
		if !reflect.DeepEqual(calls, want) {
			t.Fatalf("got calls %+v, want %+v", calls, want)
		}
	})

	t.Run("Publish supplies context to each ordered call", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mediator := NewMultiArgsMediatr(New())
		var calls []subscriberCall
		mediator.Subscribe(func(firstEvent FirstEvent, receivedContext context.Context, secondEvent SecondEvent) {
			if receivedContext != ctx {
				t.Fatal("subscriber did not receive the publication context")
			}
			calls = append(calls, subscriberCall{First: firstEvent, Second: secondEvent})
		})

		err := mediator.Publish(ctx, SecondEvent{Value: 1}, FirstEvent{Value: 1}, SecondEvent{Value: 2})

		if err != nil {
			t.Fatal(err)
		}
		want := []subscriberCall{
			{Second: SecondEvent{Value: 1}},
			{First: FirstEvent{Value: 1}, Second: SecondEvent{Value: 2}},
		}
		if !reflect.DeepEqual(calls, want) {
			t.Fatalf("got calls %+v, want %+v", calls, want)
		}
	})

	t.Run("Publish calls a repeated subscription once per registration and event", func(t *testing.T) {
		mediator := NewMultiArgsMediatr(New())
		var calls []FirstEvent
		handler := func(event FirstEvent) { calls = append(calls, event) }
		mediator.Subscribe(handler)
		mediator.Subscribe(handler)

		err := mediator.Publish(context.Background(), FirstEvent{Value: 1}, FirstEvent{Value: 2})

		if err != nil {
			t.Fatal(err)
		}
		want := []FirstEvent{{Value: 1}, {Value: 2}, {Value: 1}, {Value: 2}}
		if !reflect.DeepEqual(calls, want) {
			t.Fatalf("got calls %+v, want %+v", calls, want)
		}
	})

	t.Run("Publish fills repeated argument types in signature order", func(t *testing.T) {
		mediator := NewMultiArgsMediatr(New())
		var calls [][3]int
		mediator.Subscribe(func(first FirstEvent, second FirstEvent, secondEvent SecondEvent) {
			calls = append(calls, [3]int{first.Value, second.Value, secondEvent.Value})
		})

		err := mediator.Publish(context.Background(), FirstEvent{Value: 1}, FirstEvent{Value: 2}, SecondEvent{Value: 3})

		if err != nil {
			t.Fatal(err)
		}
		want := [][3]int{{1, 2, 3}}
		if !reflect.DeepEqual(calls, want) {
			t.Fatalf("got calls %v, want %v", calls, want)
		}
	})

	t.Run("Publish fills repeated argument types with equal events", func(t *testing.T) {
		mediator := NewMultiArgsMediatr(New())
		var calls [][3]int
		mediator.Subscribe(func(first FirstEvent, second FirstEvent, secondEvent SecondEvent) {
			calls = append(calls, [3]int{first.Value, second.Value, secondEvent.Value})
		})

		err := mediator.Publish(context.Background(), FirstEvent{Value: 1}, FirstEvent{Value: 1}, SecondEvent{Value: 2})

		if err != nil {
			t.Fatal(err)
		}
		want := [][3]int{{1, 1, 2}}
		if !reflect.DeepEqual(calls, want) {
			t.Fatalf("got calls %v, want %v", calls, want)
		}
	})

	t.Run("Publish supplies zero for a missing middle argument in a later call", func(t *testing.T) {
		type ThirdEvent struct{ Value int }
		mediator := NewMultiArgsMediatr(New())
		var calls [][3]int
		mediator.Subscribe(func(firstEvent FirstEvent, secondEvent SecondEvent, thirdEvent ThirdEvent) {
			calls = append(calls, [3]int{firstEvent.Value, secondEvent.Value, thirdEvent.Value})
		})

		err := mediator.Publish(context.Background(),
			FirstEvent{Value: 1}, SecondEvent{Value: 1}, ThirdEvent{Value: 1},
			FirstEvent{Value: 2}, ThirdEvent{Value: 2},
		)

		if err != nil {
			t.Fatal(err)
		}
		want := [][3]int{{1, 1, 1}, {2, 0, 2}}
		if !reflect.DeepEqual(calls, want) {
			t.Fatalf("got calls %v, want %v", calls, want)
		}
	})

	t.Run("Publish fills repeated argument types across multiple calls", func(t *testing.T) {
		mediator := NewMultiArgsMediatr(New())
		var calls [][3]int
		mediator.Subscribe(func(first FirstEvent, secondEvent SecondEvent, last FirstEvent) {
			calls = append(calls, [3]int{first.Value, secondEvent.Value, last.Value})
		})

		err := mediator.Publish(context.Background(),
			FirstEvent{Value: 1}, SecondEvent{Value: 1}, FirstEvent{Value: 2},
			FirstEvent{Value: 3}, SecondEvent{Value: 2}, FirstEvent{Value: 4},
		)

		if err != nil {
			t.Fatal(err)
		}
		want := [][3]int{{1, 1, 2}, {3, 2, 4}}
		if !reflect.DeepEqual(calls, want) {
			t.Fatalf("got calls %v, want %v", calls, want)
		}
	})

	t.Run("Publish delivers repeated events to both multi-event and single-event subscribers", func(t *testing.T) {
		type subscriberCall struct {
			Subscriber string
			First      FirstEvent
			Second     SecondEvent
		}
		mediator := NewMultiArgsMediatr(New())
		var calls []subscriberCall
		mediator.Subscribe(func(firstEvent FirstEvent, secondEvent SecondEvent) {
			calls = append(calls, subscriberCall{Subscriber: "combined", First: firstEvent, Second: secondEvent})
		})
		mediator.Subscribe(func(event FirstEvent) {
			calls = append(calls, subscriberCall{Subscriber: "single", First: event})
		})

		err := mediator.Publish(context.Background(), FirstEvent{Value: 1}, SecondEvent{Value: 1}, FirstEvent{Value: 2})

		if err != nil {
			t.Fatal(err)
		}
		want := []subscriberCall{
			{Subscriber: "combined", First: FirstEvent{Value: 1}, Second: SecondEvent{Value: 1}},
			{Subscriber: "combined", First: FirstEvent{Value: 2}},
			{Subscriber: "single", First: FirstEvent{Value: 1}},
			{Subscriber: "single", First: FirstEvent{Value: 2}},
		}
		if !reflect.DeepEqual(calls, want) {
			t.Fatalf("got calls %+v, want %+v", calls, want)
		}
	})

	t.Run("Publish delivers equal events containing slices and maps", func(t *testing.T) {
		type Event struct {
			Values []int
			Labels map[string]int
		}
		tests := []struct {
			name   string
			events []interface{}
			want   []Event
		}{
			{
				name: "slice values",
				events: []interface{}{
					Event{Values: []int{1, 2}},
					Event{Values: []int{1, 2}},
					Event{Values: []int{3}},
				},
				want: []Event{
					{Values: []int{1, 2}},
					{Values: []int{1, 2}},
					{Values: []int{3}},
				},
			},
			{
				name: "map values",
				events: []interface{}{
					Event{Labels: map[string]int{"value": 1}},
					Event{Labels: map[string]int{"value": 1}},
					Event{Labels: map[string]int{"value": 2}},
				},
				want: []Event{
					{Labels: map[string]int{"value": 1}},
					{Labels: map[string]int{"value": 1}},
					{Labels: map[string]int{"value": 2}},
				},
			},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				mediator := NewMultiArgsMediatr(New())
				var calls []Event
				mediator.Subscribe(func(event Event) { calls = append(calls, event) })

				err := mediator.Publish(context.Background(), test.events...)

				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(calls, test.want) {
					t.Fatalf("got calls %+v, want %+v", calls, test.want)
				}
			})
		}
	})

	t.Run("Publish delivers equal pointer contents and repeated pointers across publications", func(t *testing.T) {
		mediator := NewMultiArgsMediatr(New())
		first := &FirstEvent{Value: 1}
		equal := &FirstEvent{Value: 1}
		second := &FirstEvent{Value: 2}
		var calls []*FirstEvent
		mediator.Subscribe(func(event *FirstEvent) { calls = append(calls, event) })

		err := mediator.Publish(context.Background(), first, equal, first, second)
		if err != nil {
			t.Fatal(err)
		}
		err = mediator.Publish(context.Background(), first)
		if err != nil {
			t.Fatal(err)
		}

		want := []*FirstEvent{first, equal, first, second, first}
		if len(calls) != len(want) {
			t.Fatalf("got %d calls, want %d", len(calls), len(want))
		}
		for index, event := range want {
			if calls[index] != event {
				t.Fatalf("call %d: got %p, want %p", index, calls[index], event)
			}
		}
	})

	t.Run("Publish stops remaining calls and subscribers on an error", func(t *testing.T) {
		mediator := NewMultiArgsMediatr(New())
		wantErr := errors.New("handler failed")
		var calls []FirstEvent
		mediator.Subscribe(func(event FirstEvent) error {
			calls = append(calls, event)
			if event.Value == 2 {
				return wantErr
			}
			return nil
		})
		mediator.Subscribe(func(FirstEvent) { t.Fatal("later subscriber must not be called") })

		err := mediator.Publish(context.Background(), FirstEvent{Value: 1}, FirstEvent{Value: 2}, FirstEvent{Value: 3})

		if err != wantErr {
			t.Fatalf("got error %v, want %v", err, wantErr)
		}
		want := []FirstEvent{{Value: 1}, {Value: 2}}
		if !reflect.DeepEqual(calls, want) {
			t.Fatalf("got calls %+v, want %+v", calls, want)
		}
	})

	t.Run("Publish delivers events again after a previous publication failed", func(t *testing.T) {
		mediator := NewMultiArgsMediatr(New())
		wantErr := errors.New("handler failed")
		var calls []FirstEvent
		mediator.Subscribe(func(event FirstEvent) error {
			calls = append(calls, event)
			if event.Value == 2 {
				return wantErr
			}
			return nil
		})

		err := mediator.Publish(context.Background(), FirstEvent{Value: 1}, FirstEvent{Value: 2}, FirstEvent{Value: 3})

		if err != wantErr {
			t.Fatalf("got error %v, want %v", err, wantErr)
		}
		want := []FirstEvent{{Value: 1}, {Value: 2}}
		if !reflect.DeepEqual(calls, want) {
			t.Fatalf("got calls %+v, want %+v", calls, want)
		}

		err = mediator.Publish(context.Background(), FirstEvent{Value: 1}, FirstEvent{Value: 3})

		if err != nil {
			t.Fatal(err)
		}
		want = []FirstEvent{{Value: 1}, {Value: 2}, {Value: 1}, {Value: 3}}
		if !reflect.DeepEqual(calls, want) {
			t.Fatalf("got calls %+v, want %+v", calls, want)
		}
	})
}
