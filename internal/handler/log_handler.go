package handler

import (
	"log-processor/internal/model"
	"log-processor/internal/processor" // Importe o pacote do processador
	"net/http"

	"github.com/gin-gonic/gin"
)

// LogHandler gerencia a ingestão de logs via HTTP.
type LogHandler struct {
	// Adicione a dependência para o processador de logs.
	LogProcessor *processor.LogProcessor
}

// NewLogHandler cria uma nova instância de LogHandler.
func NewLogHandler(p *processor.LogProcessor) *LogHandler {
	return &LogHandler{
		LogProcessor: p,
	}
}

// IngestLog é o endpoint que recebe um log via POST.
func (h *LogHandler) IngestLog(c *gin.Context) {
	var rawLog model.RawLogEntry
	if err := c.ShouldBindJSON(&rawLog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de log inválido", "detalhe": err.Error()})
		return
	}

	// Envie o log para o canal de entrada do processador.
	// Isso não bloqueia a requisição HTTP.
	h.LogProcessor.Input <- rawLog

	// Responda ao cliente imediatamente.
	c.JSON(http.StatusOK, gin.H{"message": "Log recebido e fila para processamento"})
}
