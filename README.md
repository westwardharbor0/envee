# ENVEE

[pkg.go.dev](https://pkg.go.dev/github.com/westwardharbor0/envee#section-readme)

Yet another ENV configuration loader

```bash
go get github.com/westwardharbor0/envee
```

All the essential features for simple environment variable handling.
100% natural without any added packages.

**Supported features**: 
- **Required** - Set environment variable as required otherwise throw error
- **Default** - Set default value in case of environment variable missing
- **Prefix** - Set global prefix of environment variables or prefix of sub part of structure


## Example usage
```go
package main

import (
	"fmt"
	"github.com/westwardharbor0/envee"
	"time"
)

type ServiceConfig struct {
	Host  string `env:"HOST" default:"[::]"`
	Port  string `env:"PORT" default:"8080"`
	Debug bool   `env:"DEBUG" default:"false"`

	RequestTimeout time.Duration `env:"REQUEST_TIMEOUT,required"`
}

func main() {
	e := envee.New()
	e.SetPrefix("EXAMPLE_")

	var config ServiceConfig
	if err := e.Parse(&config); err != nil {
		panic(err.Error())
	}

	fmt.Println("Run on:", config.Host, ":", config.Port)
	fmt.Println("Debug:", config.Debug)
	fmt.Println("Timeout:", config.RequestTimeout.String())
}
```
Run: `env EXAMPLE_HOST="test.com" EXAMPLE_DEBUG="true" EXAMPLE_REQUEST_TIMEOUT="33s"  go run .`

## Example structs
```go
package example

type ExampleContact struct {
	Country string `env:"COUNTRY"` // <- Will be taken from CONTACT_COUNTRY due to the prefix set in tags.
	Phone string `env:"PHONE"`  // <- Will be taken from CONTACT_PHONE due to the prefix set in tags.
}

type Example struct {
	Name string `env:"NAME,required"` // <- Required name without default value. Returns error when not present
	Age int `env:"AGE" default:"42"` // <- Not required value with default that will be assigned if not present.
	Contact ExampleContact `prefix:"CONTACT"` // <- Prefix for all the environment vars.
}
```
