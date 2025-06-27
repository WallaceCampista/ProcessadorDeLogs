package main

import (
	"context"
	"fmt"
	"log"
	"log-processor/internal/config"
	"log-processor/internal/handler"
	"log-processor/internal/processor"
	"log-processor/internal/storage"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Carregue as configurações.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Crie a camada de armazenamento (agora MySQL Storager).
	// DSN (Data Source Name) para um banco de dados MySQL local.
	// Altere 'root_password' e 'log_db' se você usa credenciais diferentes.
	mysqlDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", "user_go", "user1234.", "127.0.0.1", "3306", "log_db")
	mysqlStorager, err := storage.NewMySQLStorager(mysqlDSN)
	if err != nil {
		log.Fatalf("Failed to create MySQL storager: %v", err)
	}

	// 3. Crie e inicie o processador de logs.
	bufferSize := 1000
	logProcessor := processor.NewLogProcessor(bufferSize)
	logProcessor.Start()

	// 4. Crie uma goroutine para consumir os logs...
	go func() {
		ctx := context.Background()
		for processedLog := range logProcessor.Output {
			log.Printf("--> Processed log: ID=%s | Source=%s | Severity=%s | Message='%s'\n",
				processedLog.ID, processedLog.Source, processedLog.Severity, processedLog.Message)

			if err := mysqlStorager.StoreLog(ctx, processedLog); err != nil {
				log.Printf("Failed to store log in MySQL: %v", err)
			}
		}
	}()

	// 5. Configure o servidor Gin e os handlers.
	router := gin.Default()
	logHandler := handler.NewLogHandler(logProcessor)

	// NOVO: Crie uma instância do handler de busca.
	searchHandler := handler.NewSearchHandler(mysqlStorager)

	// Defina as rotas.
	router.POST("/logs", logHandler.IngestLog)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "uptime": time.Since(time.Now()).Round(time.Second).String()})
	})

	// NOVO: Adicione a rota de busca.
	router.GET("/search", searchHandler.SearchLogs)

	// 6. Inicie o servidor em uma goroutine para lidar com o desligamento.
	go func() {
		log.Printf("Starting server on port %s...", cfg.ServerPort)
		if err := router.Run(":" + cfg.ServerPort); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	// 7. Configure o graceful shutdown (desligamento suave).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := mysqlStorager.Close(); err != nil {
		log.Printf("Error closing MySQL storager: %v", err)
	}
	log.Println("Server gracefully stopped.")
}
