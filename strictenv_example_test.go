package strictenv

import (
	"fmt"
	"time"
)

type exampleConfig struct {
	Host string `env:"APP_HOST"`
	Port int    `env:"APP_PORT"`
}

type exampleOptionalConfig struct {
	Host       string  `env:"APP_HOST"`
	Port       int     `env:"APP_PORT"`
	DBPassword *string `env:"DB_PASSWORD"`
}

func ExampleParseAsFrom() {
	cfg, err := ParseAsFrom[exampleConfig](map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
	})
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Printf("Host: %s\nPort: %d\n", cfg.Host, cfg.Port)
	// Output:
	// Host: localhost
	// Port: 8080
}

func ExampleParseAsFrom_invalidValue() {
	_, err := ParseAsFrom[exampleConfig](map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "banana",
	})
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// APP_PORT (field Port): value is invalid: parse int: strconv.ParseInt: parsing "banana": invalid syntax
}

func ExampleParseAsFrom_missingValue() {
	_, err := ParseAsFrom[exampleConfig](map[string]string{})
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// APP_HOST (field Host): value is missing or empty
	// APP_PORT (field Port): value is missing or empty
}

func ExampleParseAsFrom_missingAndInvalidValues() {
	_, err := ParseAsFrom[exampleConfig](map[string]string{
		"APP_PORT": "banana",
	})
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// APP_HOST (field Host): value is missing or empty
	// APP_PORT (field Port): value is invalid: parse int: strconv.ParseInt: parsing "banana": invalid syntax
}

func ExampleParseAsFrom_optionalField() {
	cfg, err := ParseAsFrom[exampleOptionalConfig](map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
	})
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Printf("Host: %s\n", cfg.Host)

	if cfg.DBPassword != nil {
		fmt.Printf("DBPassword: %s\n", *cfg.DBPassword)
	} else {
		fmt.Println("DBPassword: not set")
	}
	// Output:
	// Host: localhost
	// DBPassword: not set
}

// exampleQuickStartConfig mirrors the Config struct in the README's Quick
// Start section, so that example stays backed by a real, tested example.
type exampleQuickStartConfig struct {
	Env         string        `env:"APP_ENV"`
	Port        int           `env:"PORT"`
	Debug       bool          `env:"DEBUG"`
	Timeout     time.Duration `env:"TIMEOUT"`
	DatabaseURL *string       `env:"DATABASE_URL"`
	MaxConns    *int          `env:"MAX_CONNECTIONS"`
}

func ExampleParseAsFrom_quickStart() {
	cfg, err := ParseAsFrom[exampleQuickStartConfig](map[string]string{
		"APP_ENV": "production",
		"PORT":    "8080",
		"DEBUG":   "false",
		"TIMEOUT": "30s",
	})
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Printf("App booted in %s mode on port %d\n", cfg.Env, cfg.Port)

	if cfg.MaxConns != nil {
		fmt.Printf("Max connections limited to: %d\n", *cfg.MaxConns)
	} else {
		fmt.Println("Max connections: unlimited")
	}
	// Output:
	// App booted in production mode on port 8080
	// Max connections: unlimited
}

func ExampleParseAsFrom_quickStartMaxConns() {
	cfg, err := ParseAsFrom[exampleQuickStartConfig](map[string]string{
		"APP_ENV":         "production",
		"PORT":            "8080",
		"DEBUG":           "false",
		"TIMEOUT":         "30s",
		"MAX_CONNECTIONS": "100",
	})
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Printf("App booted in %s mode on port %d\n", cfg.Env, cfg.Port)

	if cfg.MaxConns != nil {
		fmt.Printf("Max connections limited to: %d\n", *cfg.MaxConns)
	} else {
		fmt.Println("Max connections: unlimited")
	}
	// Output:
	// App booted in production mode on port 8080
	// Max connections limited to: 100
}

type exampleDefaultConfig struct {
	Host string `env:"APP_HOST"`
	Port *int   `env:"APP_PORT"`
}

func ExampleParseFrom_defaults() {
	defaultPort := 8080
	cfg := exampleDefaultConfig{
		Host: "localhost",
		Port: &defaultPort,
	}

	// Only APP_HOST is required. APP_PORT overrides the default if present.
	err := ParseFrom(&cfg, map[string]string{
		"APP_HOST": "example.com",
	})
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Printf("Host: %s\n", cfg.Host)

	if cfg.Port != nil {
		fmt.Printf("Port: %d\n", *cfg.Port)
	}
	// Output:
	// Host: example.com
	// Port: 8080
}

func ExampleParseFrom_defaultOverridden() {
	port := 8080
	cfg := exampleDefaultConfig{
		Host: "localhost",
		Port: &port,
	}

	// APP_PORT is present and overrides the default.
	err := ParseFrom(&cfg, map[string]string{
		"APP_HOST": "example.com",
		"APP_PORT": "9090",
	})
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Printf("Host: %s\n", cfg.Host)

	if cfg.Port != nil {
		fmt.Printf("Port: %d\n", *cfg.Port)
	}
	// Output:
	// Host: example.com
	// Port: 9090
}
