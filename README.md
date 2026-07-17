# strictenv

> Strict environment variable parsing for Go structs. No default values, no optional fields, no surprises.

[![test](https://github.com/wayneashleyberry/strictenv/actions/workflows/test.yaml/badge.svg)](https://github.com/wayneashleyberry/strictenv/actions/workflows/test.yaml)
[![lint](https://github.com/wayneashleyberry/strictenv/actions/workflows/lint.yaml/badge.svg)](https://github.com/wayneashleyberry/strictenv/actions/workflows/lint.yaml)
[![fmt](https://github.com/wayneashleyberry/strictenv/actions/workflows/fmt.yaml/badge.svg)](https://github.com/wayneashleyberry/strictenv/actions/workflows/fmt.yaml)

Inspired by [envconfig](https://github.com/kelseyhightower/envconfig) and [env](https://github.com/caarlos0/env), but with a very different set of personal opinions. Environment variables should be [simple, explicit, and predictable](https://12factor.net/config). `strictenv` enforces that at startup: if something is missing, you find out immediately, not at 3am when a nil pointer hits production.

## Install

```
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

If `APP_HOST` or `APP_PORT` is missing or empty, `Parse` returns a single
error listing every missing variable:

```
missing env vars:
  APP_HOST (field Host)
  APP_PORT (field Port)
```

Inspect the missing variables with `errors.As`:

```go
var me *strictenv.MissingError
if errors.As(err, &me) {
	for _, m := range me.Missing {
		fmt.Printf("missing: %s (field %s)\n", m.Env, m.Field)
	}
}
```

## Supported types

`string`, `bool`, `int8`–`int64`, `uint8`–`uint64`, `float32`, `float64`,
`time.Duration`, `[]string` (comma-separated).

No maps, no custom decoders, no pointer types. Intentionally minimal.

## Testing

Use `ParseFrom` and `ParseAsFrom` to pass an explicit env map.
Tests can run in parallel without touching the real environment:

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

This follows the [12-factor app](https://12factor.net/config) approach:
env vars are granular, orthogonal controls — not grouped into "environments",
not bundled into config files, not hidden behind framework conventions.
`strictenv` makes the contract explicit: set it, or the app won't start.
