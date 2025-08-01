package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// Config armazena as configurações da aplicação.
type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`

	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
	} `mapstructure:"database"`

	RateLimiting struct {
		Rate   int    `mapstructure:"rate"`
		Period string `mapstructure:"period"`
	} `mapstructure:"rate_limiting"`
}

// LoadConfig carrega as configurações do arquivo resource.yml.
func LoadConfig() (*Config, error) {
	viper.SetConfigName("resource")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
