package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log-processor/internal/model"

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

// Close fecha a conexão com o banco de dados.
func (s *MySQLStorager) Close() error {
	log.Println("Closing MySQL database connection...")
	return s.db.Close()
}
