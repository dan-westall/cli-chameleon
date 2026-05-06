package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const FileName = "chameleon.yaml"

type Config struct {
	Name     string    `yaml:"name"`
	Commands []Command `yaml:"commands"`
}

type Command struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Run         RunField `yaml:"-"`
	Stream      bool     `yaml:"stream"`
	Parallel    bool     `yaml:"parallel"`
}

// RunField handles string or []string in YAML.
type RunField []string

func (c *Command) UnmarshalYAML(node *yaml.Node) error {
	// Decode into a raw map to handle the `run` field specially.
	var raw struct {
		Name        string   `yaml:"name"`
		Description string   `yaml:"description"`
		Run         yaml.Node `yaml:"run"`
		Stream      bool     `yaml:"stream"`
		Parallel    bool     `yaml:"parallel"`
	}
	if err := node.Decode(&raw); err != nil {
		return err
	}

	c.Name = raw.Name
	c.Description = raw.Description
	c.Stream = raw.Stream
	c.Parallel = raw.Parallel

	switch raw.Run.Kind {
	case yaml.ScalarNode:
		c.Run = RunField{raw.Run.Value}
	case yaml.SequenceNode:
		var items []string
		if err := raw.Run.Decode(&items); err != nil {
			return fmt.Errorf("invalid run field: %w", err)
		}
		c.Run = items
	case 0:
		// run field not provided
		c.Run = nil
	default:
		return fmt.Errorf("run must be a string or array of strings")
	}

	return nil
}

func Load() (*Config, error) {
	data, err := os.ReadFile(FileName)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid %s: %w", FileName, err)
	}

	return &cfg, nil
}

func Exists() bool {
	_, err := os.Stat(FileName)
	return err == nil
}
