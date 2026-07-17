package strictenv

import "fmt"

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
		"APP_PORT": "8080", //nolint:goconst
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
