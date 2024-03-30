package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ffigari/stored-strings/internal/oos"
)

type config struct {
	GoogleAppPassword              string `json:"google_app_password"`
	JWTSecret                      string `json:"jwt_secret"`
	WebPassword                    string `json:"web_password"`
	PostgresServerConnectionString string `json:"postgres_server_connection_string"`
}

func (c *config) Validate() error {
	var errors []string

	if c.GoogleAppPassword == "" {
		errors = append(errors, "missing google app password")
	}

	if c.JWTSecret == "" {
		errors = append(errors, "missing jwt secret")
	}

	if c.WebPassword == "" {
		errors = append(errors, "missing web password")
	}

	if c.PostgresServerConnectionString == "" {
		errors = append(errors, "missing postgres server connection string")
	}

	if len(errors) != 0 {
		return fmt.Errorf("invalid config: %s", strings.Join(errors, "; "))
	}

	return nil
}

var c *config

func Get() (*config, error) {
	if c != nil {
		return c, nil
	}

	bytes, err := oos.ReadFileAtRoot("config.json")

	c = &config{}

	err = json.Unmarshal(bytes, c)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling read stdin: %w", err)
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return c, nil
}
