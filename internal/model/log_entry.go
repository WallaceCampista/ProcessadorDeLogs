package model

import "time"

// LogEntry representa a estrutura de um log recebido e processado.
type LogEntry struct {
	// ID único para o log (gerado no processamento).
	ID string `json:"id"`
	// Conteúdo original do log.
	Message string `json:"message"`
	// Nível de severidade (ex: INFO, WARN, ERROR).
	Severity string `json:"severity"`
	// Fonte do log (ex: "web-server-1", "api-gateway").
	Source string `json:"source"`
	// Timestamp original do log, se disponível.
	Timestamp time.Time `json:"timestamp"`
	// Timestamp de processamento pelo nosso sistema.
	ProcessedAt time.Time `json:"processed_at"`
}

// RawLogEntry é a estrutura do log como ele chega na API.
type RawLogEntry struct {
	Message   string    `json:"message"`
	Severity  string    `json:"severity,omitempty"`
	Source    string    `json:"source,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
