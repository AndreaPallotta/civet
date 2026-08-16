package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the .civet.yml configuration file.
type Config struct {
	AI    AIConfig    `yaml:"ai,omitempty"`
	Rules RulesConfig `yaml:"rules,omitempty"`
}

// AIConfig holds AI provider settings.
type AIConfig struct {
	Provider  string `yaml:"provider,omitempty"`
	Model     string `yaml:"model,omitempty"`
	APIKeyEnv string `yaml:"api_key_env,omitempty"`
	Endpoint  string `yaml:"endpoint,omitempty"`
}

// RulesConfig holds rule customization settings.
type RulesConfig struct {
	Disabled          []string          `yaml:"disabled,omitempty"`
	SeverityOverrides map[string]string `yaml:"severity_overrides,omitempty"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{}
}

// Load reads a .civet.yml config from the given path or directory.
// If the file does not exist, returns the default config.
func Load(targetPath string) (*Config, error) {
	path := targetPath
	if fi, err := os.Stat(targetPath); err == nil && fi.IsDir() {
		path = filepath.Join(targetPath, ".civet.yml")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, err
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// IsDisabled checks if a rule ID is in the disabled list.
func (c *Config) IsDisabled(ruleID string) bool {
	for _, id := range c.Rules.Disabled {
		if id == ruleID {
			return true
		}
	}
	return false
}

// Generate writes a default .civet.yml to the given directory.
func Generate(dir string) (string, error) {
	path := filepath.Join(dir, ".civet.yml")

	content := `# Civet Configuration
# https://github.com/AndreaPallotta/civet

# AI provider settings (optional).
# Civet works fully offline without AI configuration.
# ai:
#   provider: claude           # claude, openai, gemini, ollama
#   model: claude-sonnet-4-20250514       # provider-specific model name
#   api_key_env: ANTHROPIC_API_KEY  # env var name containing the API key
#   endpoint: ""               # custom endpoint (required for ollama)

# Rule customization.
# rules:
#   disabled:                  # list of rule IDs to skip
#     - GL-003
#   severity_overrides:        # override default severity levels
#     UNI-004: warning
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", err
	}
	return path, nil
}
