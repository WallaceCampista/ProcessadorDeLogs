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
	// 1. Carregue as configurações do arquivo YAML.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Crie a camada de armazenamento (MySQL Storager).
	// Esta variável (mysqlStorager) precisa ser criada antes de ser usada.
	mysqlStorager, err := storage.NewMySQLStorager(
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
		),
	)
	if err != nil {
		log.Fatalf("Failed to create MySQL storager: %v", err)
	}

	// 3. Crie e inicie o processador de logs.
	// Esta variável (logProcessor) precisa ser criada antes de ser usada.
	bufferSize := 1000
	logProcessor := processor.NewLogProcessor(bufferSize)
	logProcessor.Start()

	// 4. Crie uma goroutine para consumir os logs processados e enviá-los ao MySQL.
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
	// Agora as variáveis logProcessor e mysqlStorager já foram declaradas.
	router := gin.Default()

	logHandler := handler.NewLogHandler(logProcessor)
	searchHandler := handler.NewSearchHandler(mysqlStorager)

	// Defina as rotas.
	router.POST("/logs", logHandler.IngestLog)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "uptime": time.Since(time.Now()).Round(time.Second).String()})
	})
	router.GET("/search", searchHandler.SearchLogs)

	// 6. Inicie o servidor em uma goroutine para lidar com o desligamento.
	go func() {
		// Acessa a porta da configuração aninhada.
		log.Printf("Starting server on port %d...", cfg.Server.Port)
		if err := router.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
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
