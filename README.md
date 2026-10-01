![Go](https://github.com/adzeitor/mediatr/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/adzeitor/mediatr/badge.svg)](https://codecov.io/gh/adzeitor/mediatr)
[![Go Report Card](https://goreportcard.com/badge/github.com/adzeitor/mediatr)](https://goreportcard.com/report/github.com/adzeitor/mediatr)

###### Command handler without return values

```go
err := mediator.Register(func(command FooEvent){
    return nil
})
```

###### Command handler can return error

```go
err := mediator.Register(func(command FooEvent) error{
    return errors.New("db error")
})
```

###### Command handler can return error and result(any type)

```go
err := mediator.Register(func(command FooEvent) (string,error){
    return "command executed", nil
})
```

###### Event handlers with multiple arguments

Use `MultiArgsMediatr` to handle multiple event types in one call. The handler
runs when at least one subscribed event is published; missing pointer arguments
are `nil`.

```go
type FooEvent struct{ Value string }
type BarEvent struct{ Value int }

mediator := mediatr.NewMultiArgsMediatr(mediatr.New())
mediator.Subscribe(func(foo *FooEvent, bar *BarEvent) {
    fmt.Printf("Foo: %v, Bar: %v\n", foo, bar)
})

foo := &FooEvent{Value: "foo"}
bar := &BarEvent{Value: 42}

if err := mediator.Publish(context.Background(), foo); err != nil {
    panic(err)
}
if err := mediator.Publish(context.Background(), bar, foo); err != nil {
    panic(err)
}
```

Output:

```text
Foo: &{foo}, Bar: <nil>
Foo: &{foo}, Bar: &{42}
```
