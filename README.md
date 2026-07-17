# strictenv

> Strict environment variable parsing for Go structs. No default values, no optional fields, no surprises.

[![Go Reference](https://pkg.go.dev/badge/github.com/wayneashleyberry/strictenv.svg)](https://pkg.go.dev/github.com/wayneashleyberry/strictenv)
[![test](https://github.com/wayneashleyberry/strictenv/actions/workflows/test.yaml/badge.svg)](https://github.com/wayneashleyberry/strictenv/actions/workflows/test.yaml)
[![lint](https://github.com/wayneashleyberry/strictenv/actions/workflows/lint.yaml/badge.svg)](https://github.com/wayneashleyberry/strictenv/actions/workflows/lint.yaml)
[![fmt](https://github.com/wayneashleyberry/strictenv/actions/workflows/fmt.yaml/badge.svg)](https://github.com/wayneashleyberry/strictenv/actions/workflows/fmt.yaml)

Inspired by [envconfig](https://github.com/kelseyhightower/envconfig) and [env](https://github.com/caarlos0/env), but with a very different set of personal opinions. Environment variables should be [simple, explicit, and predictable](https://12factor.net/config). `strictenv` enforces that at startup: if something is missing, you find out immediately, not at 3am when a nil pointer hits production.

## Install

```sh
go get github.com/wayneashleyberry/strictenv
```

## Usage

```go
type Config struct {
	Host string `env:"APP_HOST"`
	Port int    `env:"APP_PORT"`
}

cfg, err := strictenv.ParseAs[Config]()
if err != nil {
	log.Fatal(err)
}
```

If `APP_HOST` or `APP_PORT` is missing, empty, or contains a value that is invalid, `Parse` returns a single error listing every problem:

```
APP_HOST (field Host): value is missing or empty
APP_PORT (field Port): value is invalid: strconv.ParseInt: parsing "banana": invalid syntax
```

All issues are reported at once — missing variables, invalid values, and type mismatches — so you can fix everything in one pass.

Check for specific error types with `errors.Is`:

```go
errors.Is(err, strictenv.ErrMissingValue) // true if any variable is missing or empty
errors.Is(err, strictenv.ErrInvalidValue) // true if any variable has the wrong type
```

## Supported types

`string`, `bool`, `int8`–`int64`, `uint8`–`uint64`, `float32`, `float64`, `time.Duration`, `[]string` (comma-separated).

No maps, no custom decoders, no pointer types. Intentionally minimal.

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

## Philosophy

- **No defaults.** If a field is tagged, it must be set.
- **No optional tags.** There is no `required:"true"` because everything is required.
- **No prefix splitting, no `split_words`, no functional options.** One tag, one value.
- **Empty counts as missing.** An env var set to `""` is not accepted.

This follows the [12-factor app](https://12factor.net/config) approach: env vars are granular, orthogonal controls — not grouped into "environments", not bundled into config files, not hidden behind framework conventions. `strictenv` makes the contract explicit: set it, or the app won't start.

## Benchmarks

Comparison of `strictenv` vs `caarlos0/env` parsing a 6-field struct (string, int, bool, float64, time.Duration, []string) from a map.

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
