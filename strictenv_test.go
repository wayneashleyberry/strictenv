package strictenv

import (
	"errors"
	"testing"
	"time"
)

func TestParseString(t *testing.T) {
	type cfg struct {
		Name string `env:"TEST_NAME"`
	}
	t.Setenv("TEST_NAME", "hello")
	got, err := ParseAs[cfg]()
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "hello" {
		t.Errorf("got %q, want %q", got.Name, "hello")
	}
}

func TestParseBool(t *testing.T) {
	type cfg struct {
		Verbose bool `env:"TEST_VERBOSE"`
	}
	t.Setenv("TEST_VERBOSE", "true")
	got, err := ParseAs[cfg]()
	if err != nil {
		t.Fatal(err)
	}
	if !got.Verbose {
		t.Error("expected true")
	}
}

func TestParseInt(t *testing.T) {
	type cfg struct {
		Port int `env:"TEST_PORT"`
	}
	t.Setenv("TEST_PORT", "8080")
	got, err := ParseAs[cfg]()
	if err != nil {
		t.Fatal(err)
	}
	if got.Port != 8080 {
		t.Errorf("got %d, want 8080", got.Port)
	}
}

func TestParseIntSizes(t *testing.T) {
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
	t.Setenv("TEST_INT8", "127")
	t.Setenv("TEST_INT16", "32000")
	t.Setenv("TEST_INT32", "2000000")
	t.Setenv("TEST_INT64", "9000000000")
	t.Setenv("TEST_UINT", "42")
	t.Setenv("TEST_UINT8", "255")
	t.Setenv("TEST_UINT16", "65535")
	t.Setenv("TEST_UINT32", "4000000")
	t.Setenv("TEST_UINT64", "9000000000")
	got, err := ParseAs[cfg]()
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
	type cfg struct {
		Rate float64 `env:"TEST_RATE"`
	}
	t.Setenv("TEST_RATE", "3.14")
	got, err := ParseAs[cfg]()
	if err != nil {
		t.Fatal(err)
	}
	if got.Rate != 3.14 {
		t.Errorf("got %f, want 3.14", got.Rate)
	}
}

func TestParseDuration(t *testing.T) {
	type cfg struct {
		Timeout time.Duration `env:"TEST_TIMEOUT"`
	}
	t.Setenv("TEST_TIMEOUT", "5s")
	got, err := ParseAs[cfg]()
	if err != nil {
		t.Fatal(err)
	}
	if got.Timeout != 5*time.Second {
		t.Errorf("got %s, want 5s", got.Timeout)
	}
}

func TestParseStringSlice(t *testing.T) {
	type cfg struct {
		Tags []string `env:"TEST_TAGS"`
	}
	t.Setenv("TEST_TAGS", "a,b,c")
	got, err := ParseAs[cfg]()
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Tags) != 3 || got.Tags[0] != "a" || got.Tags[1] != "b" || got.Tags[2] != "c" {
		t.Errorf("got %v, want [a b c]", got.Tags)
	}
}

func TestParseMissing(t *testing.T) {
	type cfg struct {
		Name string `env:"TEST_MISSING_VAR"`
	}
	var c cfg
	err := Parse(&c)
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
	type cfg struct {
		Name string `env:"TEST_EMPTY_VAR"`
	}
	t.Setenv("TEST_EMPTY_VAR", "")
	var c cfg
	err := Parse(&c)
	if err == nil {
		t.Fatal("expected error for empty env var")
	}
	var me *MissingError
	if !errors.As(err, &me) {
		t.Fatalf("expected MissingError, got %T", err)
	}
}

func TestParseMultipleMissing(t *testing.T) {
	type cfg struct {
		A string `env:"TEST_MISS_A"`
		B string `env:"TEST_MISS_B"`
	}
	var c cfg
	err := Parse(&c)
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
	type cfg struct {
		Port int `env:"TEST_BAD_PORT"`
	}
	t.Setenv("TEST_BAD_PORT", "abc")
	var c cfg
	err := Parse(&c)
	if err == nil {
		t.Fatal("expected error for invalid int")
	}
}

func TestParseNonStructPointer(t *testing.T) {
	s := "not a struct"
	err := Parse(&s)
	if err == nil {
		t.Fatal("expected error for non-struct pointer")
	}
}

func TestParseUnexportedField(t *testing.T) {
	type cfg struct {
		visible string `env:"TEST_UNEXPORTED"`
		OK      string `env:"TEST_OK"`
	}
	t.Setenv("TEST_OK", "yes")
	got, err := ParseAs[cfg]()
	if err != nil {
		t.Fatal(err)
	}
	if got.OK != "yes" {
		t.Errorf("got %q, want %q", got.OK, "yes")
	}
}

func TestParseNoTag(t *testing.T) {
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
	type cfg struct {
		Str  string        `env:"TEST_ALL_STR"`
		Bool bool          `env:"TEST_ALL_BOOL"`
		Int  int           `env:"TEST_ALL_INT"`
		Uint uint          `env:"TEST_ALL_UINT"`
		Flt  float64       `env:"TEST_ALL_FLT"`
		Dur  time.Duration `env:"TEST_ALL_DUR"`
		Sli  []string      `env:"TEST_ALL_SLI"`
	}
	t.Setenv("TEST_ALL_STR", "hello")
	t.Setenv("TEST_ALL_BOOL", "1")
	t.Setenv("TEST_ALL_INT", "42")
	t.Setenv("TEST_ALL_UINT", "42")
	t.Setenv("TEST_ALL_FLT", "1.5")
	t.Setenv("TEST_ALL_DUR", "10m")
	t.Setenv("TEST_ALL_SLI", "x,y")
	got, err := ParseAs[cfg]()
	if err != nil {
		t.Fatal(err)
	}
	if got.Str != "hello" {
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
