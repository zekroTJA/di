# di

An experimental package providing a dependency injection container based on type parameters and reflection.

The implementation is using a multitree DAG to allow multiple references of same dependencies while disallowing circular dependencies.

Primarily, this package is designed to replace the currently used dependency injection container [sarulabs/di](https://github.com/sarulabs/di) in my Discord bot project [shinpuru](https://github.com/zekroTJA/shinpuru).

> **Warning**
> The state of the package is in a very early experimental state and the API might experience breaking changes in future updates.

## Example

```go
type Database interface {
    // ...
}

type Config interface {
    // ...
}

type WebServer interface {
    // ...
}

type MySql struct {
    Config Config `di:"singleton"`
}

type EnvConfig struct {}

type WebServerImpl struct {}

func main() {
    c := di.NewContainer()

    MustRegister[WebServer, WebServerImpl](c)
    MustRegister[Database, MySql](c)
    MustRegister[Config, EnvConfig](c, RegisterOptions[EnvConfig]{
		Setup: func(c *Container) (*EnvConfig, error) {
			return readFromEnv("MYAPP_"), nil
		},
	})

    ws := MustGet[WebServer](c)
    ws.Run()
}
```

## Advantages

- Strongly typed API due to the usage of type parameters.
- Simple API.
- Ability to use singleton as well as transistent dependencies.

## Disadvantages

- Teardown currently not implemented for transistent dependencies.
- Experimental state.
- Dependency fields must be public fields to be accessible via reflection.