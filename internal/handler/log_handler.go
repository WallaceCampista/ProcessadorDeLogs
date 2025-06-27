package handler

import (
	"log-processor/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LogHandler gerencia a ingestão de logs via HTTP.
type LogHandler struct {
	// Adicione aqui a dependência para o processador de logs
	// Ex: processor *processor.LogProcessor
}

// NewLogHandler cria uma nova instância de LogHandler.
func NewLogHandler() *LogHandler {
	return &LogHandler{}
}

// IngestLog é o endpoint que recebe um log via POST.
func (h *LogHandler) IngestLog(c *gin.Context) {
	var rawLog model.RawLogEntry
	// Tenta fazer o bind do JSON para a struct RawLogEntry.
	if err := c.ShouldBindJSON(&rawLog); err != nil {
		// Se o JSON for inválido, retorna um erro.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid log format", "details": err.Error()})
		return
	}

	// TODO: Envie o log para o processador.
	// O log `rawLog` foi recebido com sucesso.
	// Aqui, você passaria `rawLog` para o próximo passo no pipeline (o processador).

	c.JSON(http.StatusOK, gin.H{"message": "Log received successfully", "log": rawLog})
}
