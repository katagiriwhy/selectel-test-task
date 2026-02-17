package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

const (
	DefaultConfigName = ".sloglint.yaml"
)

var (
	cfg  *Config
	once sync.Once
)

type Config struct {
	Rules             Rules    `yaml:"rules"`
	SensitivePatterns []string `yaml:"sensitive_patterns"`
	ForbiddenSymbols  string   `yaml:"forbidden_symbols"`
}

type Rules struct {
	Lowercase        bool `yaml:"lowercase"`
	EnglishOnly      bool `yaml:"english_only"`
	NoSpecialSymbols bool `yaml:"no_special_symbols"`
	SensitiveData    bool `yaml:"sensitive_data"`
}

func Load(path string) *Config {
	once.Do(func() {
		cfg = defaultConfig()

		if path == "" {
			path = DefaultConfigName
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return
		}

		_ = yaml.Unmarshal(data, cfg)
	})

	return cfg
}

func defaultConfig() *Config {
	return &Config{
		Rules: Rules{
			Lowercase:        true,
			EnglishOnly:      true,
			NoSpecialSymbols: true,
			SensitiveData:    true,
		},
		SensitivePatterns: []string{
			"password",
			"pwd",
			"token",
			"api_key",
			"secret",
		},
		ForbiddenSymbols: "!?:*",
	}
}
