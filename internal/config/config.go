package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// Config armazena as configurações da aplicação.
type Config struct {
	// NOVO: Estrutura aninhada para a seção 'server' do YAML.
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
}

// LoadConfig carrega as configurações do arquivo config.yml.
func LoadConfig() (*Config, error) {
	// Define o nome do arquivo de configuração (sem a extensão).
	viper.SetConfigName("resource")
	// Adiciona o caminho onde procurar o arquivo.
	viper.AddConfigPath(".")
	// Define o tipo do arquivo.
	viper.SetConfigType("yaml")

	// Lê o arquivo de configuração.
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	// Mapeia os valores do arquivo para a struct.
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
