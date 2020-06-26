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
