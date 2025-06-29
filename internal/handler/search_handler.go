package handler

import (
	"log-processor/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SearchHandler gerencia a busca de logs via HTTP.
type SearchHandler struct {
	MySQLStorager *storage.MySQLStorager
}

// NewSearchHandler cria uma nova instância de SearchHandler.
func NewSearchHandler(s *storage.MySQLStorager) *SearchHandler {
	return &SearchHandler{
		MySQLStorager: s,
	}
}

// SearchLogs é o endpoint que recebe os parâmetros de busca e retorna os logs.
func (h *SearchHandler) SearchLogs(c *gin.Context) {
	// Lê os parâmetros da query da URL.
	source := c.Query("source")
	severity := c.Query("severity")
	// NOVO: Lê os novos parâmetros de busca.
	message := c.Query("message")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Chama o método de busca da camada de armazenamento com os novos parâmetros.
	logs, err := h.MySQLStorager.SearchLogs(c.Request.Context(), source, severity, message, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve logs", "details": err.Error()})
		return
	}

	// Retorna os logs encontrados em formato JSON.
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": logs})
}
