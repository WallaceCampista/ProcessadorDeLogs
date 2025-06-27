package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log-processor/internal/model"
	"time"

	_ "github.com/go-sql-driver/mysql" // Importa o driver MySQL
)

// MySQLStorager armazena logs no MySQL.
type MySQLStorager struct {
	db *sql.DB
}

// NewMySQLStorager cria uma nova instância de MySQLStorager e testa a conexão.
func NewMySQLStorager(dsn string) (*MySQLStorager, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Testa a conexão com o banco de dados.
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to MySQL database!")
	return &MySQLStorager{db: db}, nil
}

// StoreLog insere um log na tabela 'logs' do MySQL.
func (s *MySQLStorager) StoreLog(ctx context.Context, logEntry model.LogEntry) error {
	query := `INSERT INTO logs (id, message, severity, source, timestamp, processed_at) VALUES (?, ?, ?, ?, ?, ?)`

	// Executa a query de inserção.
	_, err := s.db.ExecContext(
		ctx,
		query,
		logEntry.ID,
		logEntry.Message,
		logEntry.Severity,
		logEntry.Source,
		logEntry.Timestamp,
		logEntry.ProcessedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert log into database: %w", err)
	}

	log.Printf("DEBUG: Successfully stored log ID %s in MySQL.", logEntry.ID)
	return nil
}

// SearchLogs busca logs na tabela 'logs' com base em critérios de filtro.
func (s *MySQLStorager) SearchLogs(ctx context.Context, source, severity string) ([]model.LogEntry, error) {
	// A consulta base.
	query := "SELECT id, message, severity, source, timestamp, processed_at FROM logs WHERE 1=1"

	// Um slice para os argumentos da consulta.
	args := []interface{}{}

	// Adiciona filtros dinamicamente se os parâmetros de busca forem fornecidos.
	if source != "" {
		query += " AND source = ?"
		args = append(args, source)
	}
	if severity != "" {
		query += " AND severity = ?"
		args = append(args, severity)
	}

	// Adiciona uma ordenação para obter os logs mais recentes primeiro.
	query += " ORDER BY processed_at DESC"

	// Executa a consulta no banco de dados.
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs from database: %w", err)
	}
	defer rows.Close()

	// Um slice para armazenar os logs encontrados.
	var logs []model.LogEntry

	// Itera sobre as linhas do resultado e escaneia cada uma para a struct LogEntry.
	for rows.Next() {
		var logEntry model.LogEntry

		// NOVO: Use variáveis intermediárias para escanear os campos de data e hora.
		var timestampStr, processedAtStr []uint8

		if err := rows.Scan(
			&logEntry.ID,
			&logEntry.Message,
			&logEntry.Severity,
			&logEntry.Source,
			&timestampStr,   // Escaneia para a variável intermediária
			&processedAtStr, // Escaneia para a variável intermediária
		); err != nil {
			return nil, fmt.Errorf("failed to scan log row: %w", err)
		}

		// NOVO: Converte as strings para time.Time
		// O formato de layout deve corresponder ao formato de data/hora do MySQL.
		layout := "2006-01-02 15:04:05" // Formato padrão do DATETIME do MySQL

		logEntry.Timestamp, err = time.Parse(layout, string(timestampStr))
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp: %w", err)
		}

		logEntry.ProcessedAt, err = time.Parse(layout, string(processedAtStr))
		if err != nil {
			return nil, fmt.Errorf("failed to parse processed_at: %w", err)
		}

		logs = append(logs, logEntry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return logs, nil
}

// Close fecha a conexão com o banco de dados.
func (s *MySQLStorager) Close() error {
	log.Println("Closing MySQL database connection...")
	return s.db.Close()
}
