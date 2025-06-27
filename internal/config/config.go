package config

import (
	"fmt"
	"os"
)

// Config armazena as configurações da aplicação.
type Config struct {
	ServerPort        string
	ElasticsearchURL  string
	ElasticsearchUser string
	ElasticsearchPass string
}

// LoadConfig carrega as configurações das variáveis de ambiente.
func LoadConfig() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// ALTERAÇÃO AQUI: Mude de http para https.
	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		esURL = "https://localhost:9200"
	}

	esUser := os.Getenv("ELASTICSEARCH_USER")
	if esUser == "" {
		esUser = "elastic"
	}

	esPass := os.Getenv("ELASTICSEARCH_PASS")
	if esPass == "" {
		return nil, fmt.Errorf("ELASTICSEARCH_PASS environment variable is not set. Please set it to connect to Elasticsearch")
	}

	return &Config{
		ServerPort:        port,
		ElasticsearchURL:  esURL,
		ElasticsearchUser: esUser,
		ElasticsearchPass: esPass,
	}, nil
}
