package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Mirrors []MirrorConfig `yaml:"mirrors"`
}

type MirrorConfig struct {
	Name   string     `yaml:"name"`
	Target Repository `yaml:"target"`
	Source Repository `yaml:"source"`
}

type Repository struct {
	URL  string `yaml:"url"`
	Auth Auth   `yaml:"auth"`
}
type Auth struct {
	Method   string `yaml:"method"` // "token", "ssh", "basic", "none"
	Token    string `yaml:"token,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	SSHKey   string `yaml:"ssh_key,omitempty"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func validateConfig(config *Config) error {
	if len(config.Mirrors) == 0 {
		return fmt.Errorf("no mirror configurations found")
	}

	for i, mirror := range config.Mirrors {
		if mirror.Target.URL == "" {
			return fmt.Errorf("mirror %d: target URL is required", i+1)
		}
		if mirror.Source.URL == "" {
			return fmt.Errorf("mirror %d: source URL is required", i+1)
		}
		if err := validateAuth(&mirror.Target.Auth, fmt.Sprintf("mirror %d target", i+1)); err != nil {
			return err
		}
		if err := validateAuth(&mirror.Source.Auth, fmt.Sprintf("mirror %d source", i+1)); err != nil {
			return err
		}
	}

	return nil
}

func validateAuth(auth *Auth, context string) error {
	switch auth.Method {
	case "token":
		if auth.Token == "" {
			return fmt.Errorf("%s: token is required when using token auth", context)
		}
	case "basic":
		if auth.Username == "" || auth.Password == "" {
			return fmt.Errorf("%s: username and password are required when using basic auth", context)
		}
	case "ssh":
		if auth.SSHKey == "" {
			return fmt.Errorf("%s: ssh_key is required when using ssh auth", context)
		}
	case "none", "":
		fmt.Println("skipping auth...")
	default:
		return fmt.Errorf("%s: unsupported auth method '%s'", context, auth.Method)
	}
	return nil
}
