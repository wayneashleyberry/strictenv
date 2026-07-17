package strictenv

import (
	"errors"
	"testing"
	"time"
)

const testVal = "hello"

func TestParseString(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Name string `env:"TEST_NAME"`
	}

	got, err := ParseAsFrom[cfg](map[string]string{"TEST_NAME": testVal})
	if err != nil {
		t.Fatal(err)
	}

	if got.Name != testVal {
		t.Errorf("got %q, want %q", got.Name, testVal)
	}
}

func TestParseBool(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Verbose bool `env:"TEST_VERBOSE"`
	}

	got, err := ParseAsFrom[cfg](map[string]string{"TEST_VERBOSE": "true"})
	if err != nil {
		t.Fatal(err)
	}

	if !got.Verbose {
		t.Error("expected true")
	}
}

func TestParseInt(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Port int `env:"TEST_PORT"`
	}

	got, err := ParseAsFrom[cfg](map[string]string{"TEST_PORT": "8080"})
	if err != nil {
		t.Fatal(err)
	}

	if got.Port != 8080 {
		t.Errorf("got %d, want 8080", got.Port)
	}
}

func TestParseIntSizes(t *testing.T) {
	t.Parallel()

	type cfg struct {
		A int8   `env:"TEST_INT8"`
		B int16  `env:"TEST_INT16"`
		C int32  `env:"TEST_INT32"`
		D int64  `env:"TEST_INT64"`
		E uint   `env:"TEST_UINT"`
		F uint8  `env:"TEST_UINT8"`
		G uint16 `env:"TEST_UINT16"`
		H uint32 `env:"TEST_UINT32"`
		I uint64 `env:"TEST_UINT64"`
	}

	got, err := ParseAsFrom[cfg](map[string]string{
		"TEST_INT8":   "127",
		"TEST_INT16":  "32000",
		"TEST_INT32":  "2000000",
		"TEST_INT64":  "9000000000",
		"TEST_UINT":   "42",
		"TEST_UINT8":  "255",
		"TEST_UINT16": "65535",
		"TEST_UINT32": "4000000",
		"TEST_UINT64": "9000000000",
	})
	if err != nil {
		t.Fatal(err)
	}

	if got.A != 127 || got.B != 32000 || got.C != 2000000 || got.D != 9000000000 {
		t.Errorf("int sizes wrong: %+v", got)
	}

	if got.E != 42 || got.F != 255 || got.G != 65535 || got.H != 4000000 || got.I != 9000000000 {
		t.Errorf("uint sizes wrong: %+v", got)
	}
}

func TestParseFloat(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Rate float64 `env:"TEST_RATE"`
	}

	got, err := ParseAsFrom[cfg](map[string]string{"TEST_RATE": "3.14"})
	if err != nil {
		t.Fatal(err)
	}

	if got.Rate != 3.14 {
		t.Errorf("got %f, want 3.14", got.Rate)
	}
}

func TestParseDuration(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Timeout time.Duration `env:"TEST_TIMEOUT"`
	}

	got, err := ParseAsFrom[cfg](map[string]string{"TEST_TIMEOUT": "5s"})
	if err != nil {
		t.Fatal(err)
	}

	if got.Timeout != 5*time.Second {
		t.Errorf("got %s, want 5s", got.Timeout)
	}
}

func TestParseStringSlice(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Tags []string `env:"TEST_TAGS"`
	}

	got, err := ParseAsFrom[cfg](map[string]string{"TEST_TAGS": "a,b,c"})
	if err != nil {
		t.Fatal(err)
	}

	if len(got.Tags) != 3 || got.Tags[0] != "a" || got.Tags[1] != "b" || got.Tags[2] != "c" {
		t.Errorf("got %v, want [a b c]", got.Tags)
	}
}

func TestParseMissing(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Name string `env:"TEST_MISSING_VAR"`
	}

	var c cfg

	err := ParseFrom(&c, nil)
	if err == nil {
		t.Fatal("expected error for missing env var")
	}

	var me *MissingError
	if !errors.As(err, &me) {
		t.Fatalf("expected MissingError, got %T", err)
	}

	if len(me.Missing) != 1 {
		t.Fatalf("expected 1 missing, got %d", len(me.Missing))
	}

	if me.Missing[0].Field != "Name" || me.Missing[0].Env != "TEST_MISSING_VAR" {
		t.Errorf("wrong missing: %+v", me.Missing[0])
	}
}

func TestParseEmpty(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Name string `env:"TEST_EMPTY_VAR"`
	}

	var c cfg

	err := ParseFrom(&c, map[string]string{"TEST_EMPTY_VAR": ""})
	if err == nil {
		t.Fatal("expected error for empty env var")
	}

	if _, ok := errors.AsType[*MissingError](err); !ok {
		t.Fatalf("expected MissingError, got %T", err)
	}
}

func TestParseMultipleMissing(t *testing.T) {
	t.Parallel()

	type cfg struct {
		A string `env:"TEST_MISS_A"`
		B string `env:"TEST_MISS_B"`
	}

	var c cfg

	err := ParseFrom(&c, nil)
	if err == nil {
		t.Fatal("expected error")
	}

	var me *MissingError
	if !errors.As(err, &me) {
		t.Fatalf("expected MissingError, got %T", err)
	}

	if len(me.Missing) != 2 {
		t.Errorf("expected 2 missing, got %d", len(me.Missing))
	}
}

func TestParseInvalidValue(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Port int `env:"TEST_BAD_PORT"`
	}

	var c cfg

	err := ParseFrom(&c, map[string]string{"TEST_BAD_PORT": "abc"})
	if err == nil {
		t.Fatal("expected error for invalid int")
	}
}

func TestParseNonStructPointer(t *testing.T) {
	t.Parallel()

	s := "not a struct"

	err := Parse(&s)
	if err == nil {
		t.Fatal("expected error for non-struct pointer")
	}
}

func TestParseUnexportedField(t *testing.T) {
	t.Parallel()

	type cfg struct {
		visible string //nolint: unused
		OK      string `env:"TEST_UNEXP_OK"`
	}

	got, err := ParseAsFrom[cfg](map[string]string{"TEST_UNEXP_OK": "yes"})
	if err != nil {
		t.Fatal(err)
	}

	if got.OK != "yes" {
		t.Errorf("got %q, want %q", got.OK, "yes")
	}
}

func TestParseNoTag(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Ignored string
	}

	var c cfg

	err := Parse(&c)
	if err != nil {
		t.Fatalf("expected no error for untagged field, got %v", err)
	}
}

func TestParseAllTypes(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Str  string        `env:"TEST_ALL_STR"`
		Bool bool          `env:"TEST_ALL_BOOL"`
		Int  int           `env:"TEST_ALL_INT"`
		Uint uint          `env:"TEST_ALL_UINT"`
		Flt  float64       `env:"TEST_ALL_FLT"`
		Dur  time.Duration `env:"TEST_ALL_DUR"`
		Sli  []string      `env:"TEST_ALL_SLI"`
	}

	got, err := ParseAsFrom[cfg](map[string]string{
		"TEST_ALL_STR": testVal,
		"TEST_ALL_BOOL": "1",
		"TEST_ALL_INT":  "42",
		"TEST_ALL_UINT": "42",
		"TEST_ALL_FLT":  "1.5",
		"TEST_ALL_DUR":  "10m",
		"TEST_ALL_SLI":  "x,y",
	})
	if err != nil {
		t.Fatal(err)
	}

	if got.Str != testVal {
		t.Errorf("Str: got %q", got.Str)
	}

	if !got.Bool {
		t.Error("Bool: expected true")
	}

	if got.Int != 42 {
		t.Errorf("Int: got %d", got.Int)
	}

	if got.Uint != 42 {
		t.Errorf("Uint: got %d", got.Uint)
	}

	if got.Flt != 1.5 {
		t.Errorf("Flt: got %f", got.Flt)
	}

	if got.Dur != 10*time.Minute {
		t.Errorf("Dur: got %s", got.Dur)
	}

	if len(got.Sli) != 2 || got.Sli[0] != "x" || got.Sli[1] != "y" {
		t.Errorf("Sli: got %v", got.Sli)
	}
}
