package mediatr

import "testing"

var testEvents = []Event{
	FooEvent{},
	BarEvent{},
	QuxEvent{},
	FooEvent{},
	BarEvent{},
	QuxEvent{},
	FooEvent{},
	BarEvent{},
	QuxEvent{},
	FooEvent{},
	BarEvent{},
	QuxEvent{},
}

func Benchmark_staticMediator(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, event := range testEvents {
			dispatchToStaticMediator(event)
		}
	}
}

func Benchmark_Mediator(b *testing.B) {
	mediator := NewMediator()
	mediator.Subscribe(FooEvent{}, FooEventHandler)
	mediator.Subscribe(BarEvent{}, BarEventHandler)
	mediator.Subscribe(QuxEvent{}, QuxEventHandler)

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		for _, event := range testEvents {
			mediator.Publish(event)
		}
	}
	b.StopTimer()
}

func Benchmark_ReflectMediator(b *testing.B) {
	mediator := NewReflectMediator()
	mediator.Subscribe(FooEventTypedHandler)
	mediator.Subscribe(BarEventTypedHandler)
	mediator.Subscribe(QuxEventTypedHandler)

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		for _, event := range testEvents {
			mediator.Publish(event)
		}
	}
	b.StopTimer()
}
