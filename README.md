
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
