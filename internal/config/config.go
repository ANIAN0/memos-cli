// Package config provides Memos-specific configuration.
package config

import (
	"fmt"

	"gopkg.in/yaml.v3"

	sharedconfig "github.com/ANIAN0/memos-cli/pkg/config"
)

// Config is the Memos-specific configuration.
// It embeds the shared Config and adds Memos-specific fields.
type Config struct {
	sharedconfig.Config `yaml:",inline"`

	// InstanceURL is the Memos server URL (e.g., "https://memos.example.com").
	InstanceURL string `yaml:"instance_url"`

	// AccessToken is the Memos access token (supports ${ENV_VAR} interpolation).
	AccessToken string `yaml:"access_token"`

	// DefaultPageSize is the default page size for list operations.
	DefaultPageSize int `yaml:"default_page_size"`

	// DefaultVisibility is the default visibility for new memos (PRIVATE, PUBLIC, PROTECTED).
	DefaultVisibility string `yaml:"default_visibility"`
}

// LoadFromBytes parses a YAML config into c.
func LoadFromBytes(data []byte) (*Config, error) {
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &c, nil
}

// Validate checks if the config is valid.
func (c *Config) Validate() error {
	if c.InstanceURL == "" {
		return fmt.Errorf("instance_url is required")
	}
	if c.AccessToken == "" {
		return fmt.Errorf("access_token is required")
	}
	return nil
}