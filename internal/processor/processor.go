package processor

import (
	"log"
	"log-processor/internal/model"
	"time"

	"github.com/google/uuid" // Você precisará desta biblioteca para gerar IDs
)

// LogProcessor processa e enriquece os logs.
type LogProcessor struct {
	// O canal de entrada para receber logs crus.
	Input chan model.RawLogEntry
	// O canal de saída para enviar logs processados.
	Output chan model.LogEntry
}

// NewLogProcessor cria uma nova instância de LogProcessor.
func NewLogProcessor(bufferSize int) *LogProcessor {
	return &LogProcessor{
		Input:  make(chan model.RawLogEntry, bufferSize),
		Output: make(chan model.LogEntry, bufferSize),
	}
}

// Start inicia o processador.
// Ele lê logs do canal de entrada e os processa em goroutines.
func (p *LogProcessor) Start() {
	// Inicia um loop infinito para processar logs do canal de entrada.
	go func() {
		log.Println("Log processor started. Waiting for logs...")
		for rawLog := range p.Input {
			// Para cada log recebido no canal, inicie uma goroutine para processá-lo.
			// Isso permite processar múltiplos logs concorrentemente.
			go p.processLog(rawLog)
		}
		log.Println("Log processor stopped.")
	}()
}

// processLog executa a lógica de enriquecimento em um log.
func (p *LogProcessor) processLog(rawLog model.RawLogEntry) {
	// LINHA DE DEBUG ANTERIOR
	log.Printf("DEBUG: Processing log message: '%s' from source '%s'", rawLog.Message, rawLog.Source)

	// ... (código de enriquecimento, sem mudanças) ...
	id := uuid.New().String()
	severity := rawLog.Severity
	if severity == "" {
		severity = "INFO"
	}
	source := rawLog.Source
	if source == "" {
		source = "unknown"
	}
	timestamp := rawLog.Timestamp
	if timestamp.IsZero() {
		timestamp = time.Now()
	}

	processedLog := model.LogEntry{
		ID:          id,
		Message:     rawLog.Message,
		Severity:    severity,
		Source:      source,
		Timestamp:   timestamp,
		ProcessedAt: time.Now(),
	}

	// Envie o log processado para o canal de saída.
	p.Output <- processedLog

	// NOVA LINHA DE DEBUG: Adicione esta linha para confirmar que o envio para o canal de saída não travou.
	log.Printf("DEBUG: Log ID %s sent to output channel.", processedLog.ID)
}

// Close fecha os canais de entrada e saída.
func (p *LogProcessor) Close() {
	close(p.Input)
	close(p.Output)
}
