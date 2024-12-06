package configfile

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// Pod represents the pod configuration
type Pod struct {
	Name    string `toml:"name"`
	Quadlet string `toml:"quadlet"`
}

// Quadlet represents individual quadlet configurations
type Quadlet struct {
	Name    string `toml:"name"`
	Type    string `toml:"type"`
	Quadlet string `toml:"quadlet"`
}

type Replaceables struct {
	Id   string `toml:"id"`
	Dev  string `toml:"dev"`
	Prod string `toml:"prod"`
}

// Config represents the overall configuration
type Config struct {
	Pod          Pod            `toml:"main_pod"`
	Quadlets     []Quadlet      `toml:"quadlets"`
	Replaceables []Replaceables `toml:"replaceables"`
}

// ParseConfigFromTOML parses TOML configuration into a Config struct and errors on extra keys
func ParseConfigFromTOML(tomlData string) (*Config, *toml.MetaData, error) {
	var config Config

	metaData, err := toml.Decode(tomlData, &config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse TOML: %w", err)
	}

	return &config, &metaData, nil
}
