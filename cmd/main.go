package main

import (
	"log"
	"log-processor/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inicializa o roteador do Gin.
	router := gin.Default()

	// Crie uma instância do handler de logs.
	logHandler := handler.NewLogHandler()

	// Defina a rota POST para ingestão de logs.
	router.POST("/logs", logHandler.IngestLog)

	// Defina uma rota de saúde para verificar se o serviço está rodando.
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Inicie o servidor na porta 8080.
	// Em um sistema real, essa porta viria de uma configuração.
	log.Println("Starting server on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
