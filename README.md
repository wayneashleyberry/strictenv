# strictenv

> ⚡ _The fastest, most memory-efficient strict environment variable parser for Go structs. Zero implicit defaults, zero dependencies, zero surprises, tons of opinions._

[![Go Reference](https://pkg.go.dev/badge/github.com/wayneashleyberry/strictenv.svg)](https://pkg.go.dev/github.com/wayneashleyberry/strictenv)
[![test](https://github.com/wayneashleyberry/strictenv/actions/workflows/test.yaml/badge.svg)](https://github.com/wayneashleyberry/strictenv/actions/workflows/test.yaml)
[![lint](https://github.com/wayneashleyberry/strictenv/actions/workflows/lint.yaml/badge.svg)](https://github.com/wayneashleyberry/strictenv/actions/workflows/lint.yaml)
[![fmt](https://github.com/wayneashleyberry/strictenv/actions/workflows/fmt.yaml/badge.svg)](https://github.com/wayneashleyberry/strictenv/actions/workflows/fmt.yaml)

Inspired by [envconfig](https://github.com/kelseyhightower/envconfig) and [env](https://github.com/caarlos0/env), but with a very different set of personal opinions. Environment variables should be [simple, explicit, and predictable](https://12factor.net/config). `strictenv` enforces that at startup: if something is missing, you find out immediately, not at 3am when a nil pointer hits production.

## Why `strictenv`?

Most Go environment variable parsers fail silently or make assumptions when an environment variable is missing or empty. They silently fall back to Go's type zero-values (`0`, `false`, `""`).

This creates a dangerous runtime ambiguity. For example:

- If `TIMEOUT=""`, does the user want a timeout of `0` (never timeout), or did they just forget to fill it out, intending to use a safe default?
- If `JWT_SECRET=""`, does your app boot up with an empty string as a signature key, exposing you to critical security vulnerabilities?

`strictenv` fixes this by treating environment variables as **deterministic and explicit**. If a variable is declared in your struct, it must exist and be valid—unless you explicitly define it as a pointer.

## Features

- **No Implicit Zero-Values:** An empty string (`PORT=""`) will fail to parse as an integer rather than defaulting to `0`.
- **Fail-Fast Initialization:** Clear, highly descriptive errors on application startup so configuration issues never leak into runtime.
- **Explicit Optionality:** Uses standard Go pointers (`*int`, `*string`) to distinguish between a missing/null value and an explicit zero.
- **Zero Dependencies:** Built entirely on top of the Go standard library.
- **Supported Types:** `string`, `bool`, `int`, `int8`–`int64`, `uint`, `uint8`–`uint64`, `float32`, `float64`, `time.Duration`, `[]string` (comma-separated).

## Install

```sh
go get github.com/wayneashleyberry/strictenv
```

## Quick start

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/wayneashleyberry/strictenv"
)

type Config struct {
	// Required fields (Value must be present and valid)
	Env         string        `env:"APP_ENV"`
	Port        int           `env:"PORT"`
	Debug       bool          `env:"DEBUG"`
	Timeout     time.Duration `env:"TIMEOUT"`

	// Optional fields (Uses pointers to safely handle missing/empty values)
	DatabaseURL *string       `env:"DATABASE_URL"`
	MaxConns    *int          `env:"MAX_CONNECTIONS"`
}

func main() {
	cfg, err := strictenv.ParseAs[Config]()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	fmt.Printf("App booted in %s mode on port %d\n", cfg.Env, cfg.Port)

	// Safely consume optional pointer values
	if cfg.MaxConns != nil {
		fmt.Printf("Max connections limited to: %d\n", *cfg.MaxConns)
	} else {
		fmt.Println("Max connections: unlimited")
	}
}
```

## Best Practices

### Handling Defaults

Because `strictenv` purposefully does not support a default struct tag, default configuration logic should live explicitly in your application code where it is visible and testable.

Note that this only works with **pointer** fields. Non-pointer fields are always required: if the env var is missing or empty, `Parse` returns `ErrMissingValue` regardless of any value you preset on the struct, since there is no way to distinguish "left as the zero value" from "deliberately defaulted".

```go
type Config struct {
    Port *int `env:"PORT"` // optional, falls back to the preset default below
}

port := 8080
cfg := Config{
    // Define your defaults upfront in standard Go
    Port: &port,
}

// PORT, if present in the environment, overwrites the default.
// If PORT is absent, cfg.Port keeps pointing at the preset value.
if err := strictenv.Parse(&cfg); err != nil {
    log.Fatalf("Invalid config: %v", err)
}
```

### Avoiding Nil Pointer Dereferences

When using optional fields (pointers), always perform a `nil` check before extracting the value to prevent runtime panics:

```go
// Bad: Might panic if DATABASE_URL was missing from the environment
connect(*cfg.DatabaseURL)

// Good: Checked and handled explicitly
if cfg.DatabaseURL != nil {
    connect(*cfg.DatabaseURL)
}
```

## Testing

Use `ParseFrom` and `ParseAsFrom` to pass an explicit env map. Tests can run in parallel without touching the real environment:

```go
func TestConfig(t *testing.T) {
	t.Parallel()

	cfg, err := strictenv.ParseAsFrom[Config](map[string]string{
		"APP_ENV": "test",
		"PORT":    "8080",
		"DEBUG":   "true",
		"TIMEOUT": "5s",
	})
	if err != nil {
		t.Fatal(err)
	}
	// ...
}
```

## Benchmarks

| Library                                                                   |   ns/op |   B/op | allocs/op |
| ------------------------------------------------------------------------- | ------: | -----: | --------: |
| **strictenv**                                                             | **452** | **64** |     **1** |
| [syntaqx/env](https://github.com/syntaqx/env)                             |   1,255 |    436 |        15 |
| [vrischmann/envconfig](https://github.com/vrischmann/envconfig)           |   1,381 |    576 |        19 |
| [go-simpler/env](https://github.com/junk1tm/env)                          |   1,560 |  1,856 |        19 |
| [cleanenv](https://github.com/ilyakaznacheev/cleanenv)                    |   2,300 |  2,488 |        41 |
| [kelseyhightower/envconfig](https://github.com/kelseyhightower/envconfig) |   2,567 |  1,544 |        62 |
| [caarlos0/env](https://github.com/caarlos0/env)                           |   5,624 |  7,280 |        71 |
| [cristalhq/aconfig](https://github.com/cristalhq/aconfig)                 |   6,135 |  5,754 |       125 |

`strictenv` is the fastest and most memory-efficient Go env parser benchmarked — **3–13× faster** and **7–113× less memory** than alternatives, with a single allocation per parse.

> Full benchmarks and methodology: [strictenv-benchmarks](https://github.com/wayneashleyberry/strictenv-benchmarks)
