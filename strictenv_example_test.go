package strictenv

import "fmt"

type exampleConfig struct {
	Host string `env:"APP_HOST"`
	Port int    `env:"APP_PORT"`
}

func ExampleParseAs_from() {
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

func ExampleParseAs_from_invalidValue() {
	_, err := ParseAsFrom[exampleConfig](map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "banana",
	})
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// strictenv: field Port (env APP_PORT): parse int: strconv.ParseInt: parsing "banana": invalid syntax
}

func ExampleMissingError() {
	_, err := ParseAsFrom[exampleConfig](map[string]string{})
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// missing env vars:
	//   APP_HOST (field Host)
	//   APP_PORT (field Port)
}
