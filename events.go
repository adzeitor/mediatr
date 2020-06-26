package mediatr

type FooEvent struct{}

func (FooEvent) Name() string {
	return "foo"
}

type BarEvent struct{}

func (BarEvent) Name() string {
	return "bar"
}

type QuxEvent struct{}

func (QuxEvent) Name() string {
	return "qux"
}

func FooEventTypedHandler(_ FooEvent) {

}

func BarEventTypedHandler(_ BarEvent) {

}

func QuxEventTypedHandler(_ QuxEvent) {

}

func FooEventHandler(event Event) {
	_ = event.(FooEvent)
}

func BarEventHandler(event Event) {
	_ = event.(BarEvent)
}

func QuxEventHandler(event Event) {
	_ = event.(QuxEvent)
}
