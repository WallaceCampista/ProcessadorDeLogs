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

	log.Println("Conectado com sucesso ao banco de dados MySQL!")
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
		return fmt.Errorf("falha ao inserir login no banco de dados: %w", err)
	}

	log.Printf("DEBUG:ID de log armazenado com sucesso %s em MySQL.", logEntry.ID)
	return nil
}

// SearchLogs busca logs na tabela 'logs' com base em critérios de filtro.
// NOVO: Adiciona filtros por 'message', 'startDate' e 'endDate'.
func (s *MySQLStorager) SearchLogs(ctx context.Context, source, severity, message, startDate, endDate string) ([]model.LogEntry, error) {
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

	// NOVO: Filtro por palavra-chave na mensagem (busca por substring).
	if message != "" {
		query += " AND message LIKE ?"
		args = append(args, "%"+message+"%") // O '%' permite buscar por substring.
	}

	// NOVO: Filtro por intervalo de datas.
	layout := "2006-01-02" // Layout de data esperado (ex: 2025-06-27)

	if startDate != "" {
		// Converte a string de data para o formato de data/hora.
		query += " AND timestamp >= ?"
		args = append(args, startDate) // MySQL pode comparar strings de data diretamente.
	}

	if endDate != "" {
		// Adiciona um dia à data final para incluir todo o dia.
		// Isso é crucial para queries de data.
		parsedEndDate, err := time.Parse(layout, endDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date format: %w", err)
		}
		endDateExclusive := parsedEndDate.AddDate(0, 0, 1).Format(layout)
		query += " AND timestamp < ?"
		args = append(args, endDateExclusive)
	}

	// Adiciona uma ordenação para obter os logs mais recentes primeiro.
	query += " ORDER BY processed_at DESC"

	// ... (restante do código permanece o mesmo: QueryContext, Scan, etc.) ...

	// Executa a consulta no banco de dados.
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs from database: %w", err)
	}
	defer rows.Close()

	var logs []model.LogEntry
	for rows.Next() {
		var logEntry model.LogEntry
		var timestampStr, processedAtStr []uint8

		if err := rows.Scan(
			&logEntry.ID,
			&logEntry.Message,
			&logEntry.Severity,
			&logEntry.Source,
			&timestampStr,
			&processedAtStr,
		); err != nil {
			return nil, fmt.Errorf("failed to scan log row: %w", err)
		}

		loc, err := time.LoadLocation("Local")
		if err != nil {
			return nil, fmt.Errorf("failed to load local timezone: %w", err)
		}

		logEntry.Timestamp, err = time.ParseInLocation("2006-01-02 15:04:05", string(timestampStr), loc)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp: %w", err)
		}

		logEntry.ProcessedAt, err = time.ParseInLocation("2006-01-02 15:04:05", string(processedAtStr), loc)
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
	log.Println("Fechando a conexão do banco de dados MySQL...")
	return s.db.Close()
}
