# strictenv

> Strict environment variable parsing for Go structs. No default values, no optional tags, no implicit zero-values, no surprises.

[![Go Reference](https://pkg.go.dev/badge/github.com/wayneashleyberry/strictenv.svg)](https://pkg.go.dev/github.com/wayneashleyberry/strictenv)
[![test](https://github.com/wayneashleyberry/strictenv/actions/workflows/test.yaml/badge.svg)](https://github.com/wayneashleyberry/strictenv/actions/workflows/test.yaml)
[![lint](https://github.com/wayneashleyberry/strictenv/actions/workflows/lint.yaml/badge.svg)](https://github.com/wayneashleyberry/strictenv/actions/workflows/lint.yaml)
[![fmt](https://github.com/wayneashleyberry/strictenv/actions/workflows/fmt.yaml/badge.svg)](https://github.com/wayneashleyberry/strictenv/actions/workflows/fmt.yaml)

Inspired by [envconfig](https://github.com/kelseyhightower/envconfig) and [env](https://github.com/caarlos0/env), but with a very different set of personal opinions. Environment variables should be [simple, explicit, and predictable](https://12factor.net/config). `strictenv` enforces that at startup: if something is missing, you find out immediately, not at 3am when a nil pointer hits production.

## Why `strictenv`?

Most Go environment variable parsers fail silently or make assumptions when an environment variable is missing or empty. They silently fall back to Go's type zero-values (`0`, `false`, `""`). 

This creates a dangerous runtime ambiguity. For example:
* If `TIMEOUT=""`, does the user want a timeout of `0` (never timeout), or did they just forget to fill it out, intending to use a safe default?
* If `JWT_SECRET=""`, does your app boot up with an empty string as a signature key, exposing you to critical security vulnerabilities?

`strictenv` fixes this by treating environment variables as **deterministic and explicit**. If a variable is declared in your struct, it must exist and be valid—unless you explicitly define it as a pointer.

## Features

- **No Implicit Zero-Values:** An empty string (`PORT=""`) will fail to parse as an integer rather than defaulting to `0`.
- **Fail-Fast Initialization:** Clear, highly descriptive errors on application startup so configuration issues never leak into runtime.
- **Explicit Optionality:** Uses standard Go pointers (`*int`, `*string`) to distinguish between a missing/null value and an explicit zero.
- **Zero Dependencies:** Built entirely on top of the Go standard library.

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
	var cfg Config

	// Parse environment variables directly into the struct
	if err := strictenv.Parse(&cfg); err != nil {
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

Because `strictenv` purposefully does not support a default struct tag, default configuration logic should live explicitly in your application code where it is visible and testable:

```go
cfg := Config{
    // Define your defaults upfront in standard Go
    Port: 8080, 
}

// Any environment variables present will strictly overwrite these values.
// Any missing non-pointer fields will throw an error.
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

## Supported types

`string`, `bool`, `int8`–`int64`, `uint8`–`uint64`, `float32`, `float64`, `time.Duration`, `[]string` (comma-separated).

Pointer types (`*string`, `*int`, etc.) are supported for optional fields — see below.

## Testing

Use `ParseFrom` and `ParseAsFrom` to pass an explicit env map. Tests can run in parallel without touching the real environment:

```go
func TestConfig(t *testing.T) {
	t.Parallel()

	cfg, err := strictenv.ParseAsFrom[Config](map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
	})
	if err != nil {
		t.Fatal(err)
	}
	// ...
}
```

## Benchmarks

Comparison of `strictenv` vs `caarlos0/env` parsing a 6-field struct (`string`, `int`, `bool`, `float64`, `time.Duration`, `[]string`) from a map.

```
goos: darwin
goarch: arm64
pkg: github.com/wayneashleyberry/strictenv
cpu: Apple M2 Pro
```

| Benchmark            | ns/op | B/op | allocs/op |
| -------------------- | ----- | ---- | --------- |
| `Strictenv`          | 478   | 72   | 2         |
| `StrictenvGeneric`   | 508   | 152  | 3         |
| `Caarlos0Env`        | 5807  | 7560 | 81        |
| `Caarlos0EnvGeneric` | 5838  | 7640 | 82        |

`strictenv` is ~12× faster and uses ~50× fewer allocations.
