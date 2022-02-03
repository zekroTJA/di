# di

A ***(very)*** WIP Go 1.18+ generic dependency injection package based on type reflection.

## Usage

Below, you can see a very simplified demonstration on how this package works. If you want to see a more comprehensive example, please take a look at the [examples](example/) listed.
```go
type Service1 interface {}

type Service2 interface {}

type service1Impl struct {
    S2 Service2
}

type service2Impl struct {
    S1 Service1
}

func main() {
    c := di.NewContainer()

    di.MustRegister[Service1, service1Impl](c)
    di.MustRegister[Service2, service2Impl](c)

    s1 := di.MustGet[Service1](c)
}
```

## Limitations

There are some severe limitations with this package.

- To set fields automatically via the injection system, they **must** be exported. All unexported fields are ignored.
- Currently, only interface fields can be automatically assigned to instances of registered services.
- Currently, there are only [singleton](https://docs.microsoft.com/en-us/dotnet/core/extensions/dependency-injection#service-lifetimes) instances because of simplicity.