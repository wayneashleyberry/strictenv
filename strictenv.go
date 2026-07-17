package strictenv

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// MissingError reports environment variables that are missing or empty.
type MissingError struct {
	Missing []MissingVar
}

// MissingVar is a single missing environment variable.
type MissingVar struct {
	Field string
	Env   string
}

func (e *MissingError) Error() string {
	if len(e.Missing) == 1 {
		return fmt.Sprintf("missing env var %s (field %s)", e.Missing[0].Env, e.Missing[0].Field)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "missing env vars:")

	for _, m := range e.Missing {
		fmt.Fprintf(&b, "\n  %s (field %s)", m.Env, m.Field)
	}

	return b.String()
}

// Parse populates dst, a pointer to a struct, from environment variables.
// Every exported field with an "env" tag is required. Missing or empty
// values are collected and returned as a single error.
func Parse(dst any) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("strictenv: Parse requires a pointer to a struct, got %T", dst)
	}

	v = v.Elem()
	t := v.Type()

	var missing []MissingVar

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		envKey := field.Tag.Get("env")
		if envKey == "" {
			continue
		}

		val := os.Getenv(envKey)
		if val == "" {
			missing = append(missing, MissingVar{Field: field.Name, Env: envKey})

			continue
		}

		err := setField(v.Field(i), field.Type, val)
		if err != nil {
			return fmt.Errorf("strictenv: field %s (env %s): %w", field.Name, envKey, err)
		}
	}

	if len(missing) > 0 {
		return &MissingError{Missing: missing}
	}

	return nil
}

// ParseAs parses environment variables into a new value of type T.
func ParseAs[T any]() (T, error) {
	var zero T

	err := Parse(&zero)
	if err != nil {
		return zero, err
	}

	return zero, nil
}

func setField(f reflect.Value, t reflect.Type, val string) error {
	if t.Kind() == reflect.Slice {
		if t.Elem().Kind() != reflect.String {
			return fmt.Errorf("unsupported slice element type %s", t.Elem())
		}

		f.Set(reflect.ValueOf(strings.Split(val, ",")))

		return nil
	}

	switch t.Kind() {
	case reflect.String:
		f.SetString(val)
	case reflect.Bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}

		f.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if t == reflect.TypeFor[time.Duration]() {
			d, err := time.ParseDuration(val)
			if err != nil {
				return err
			}

			f.SetInt(int64(d))

			return nil
		}

		n, err := strconv.ParseInt(val, 0, t.Bits())
		if err != nil {
			return err
		}

		f.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(val, 0, t.Bits())
		if err != nil {
			return err
		}

		f.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(val, t.Bits())
		if err != nil {
			return err
		}

		f.SetFloat(n)
	default:
		return fmt.Errorf("unsupported type %s", t)
	}

	return nil
}
