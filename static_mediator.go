package mediatr

func dispatchToStaticMediator(event Event) {
	switch e := event.(type) {
	case FooEvent:
		FooEventTypedHandler(e)
	case BarEvent:
		BarEventTypedHandler(e)
	case QuxEvent:
		QuxEventTypedHandler(e)
	}
}
