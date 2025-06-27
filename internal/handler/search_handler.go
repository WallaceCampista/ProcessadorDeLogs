package handler

import (
	"log-processor/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SearchHandler gerencia a busca de logs via HTTP.
type SearchHandler struct {
	// O handler precisa de uma referência para a camada de armazenamento.
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
	// Lê os parâmetros da query da URL (ex: ?source=web-app&severity=ERROR).
	source := c.Query("source")
	severity := c.Query("severity")

	// Chama a função de busca da camada de armazenamento.
	logs, err := h.MySQLStorager.SearchLogs(c.Request.Context(), source, severity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve logs", "details": err.Error()})
		return
	}

	// Retorna os logs encontrados em formato JSON.
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": logs})
}
