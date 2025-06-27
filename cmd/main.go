package main

import (
	"log"
	"log-processor/internal/handler"
	"log-processor/internal/processor"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Crie o processador de logs.
	// O buffer de 1000 logs permite que a API receba logs rapidamente sem bloquear,
	// enquanto o processador os consome em segundo plano.
	bufferSize := 1000
	logProcessor := processor.NewLogProcessor(bufferSize)

	// 2. Inicie o processador em uma goroutine.
	// Ele começará a ouvir logs no seu canal de entrada.
	logProcessor.Start()

	// 3. (Opcional) Crie uma goroutine para "consumir" os logs processados.
	// Por enquanto, apenas os imprimiremos.
	go func() {
		for processedLog := range logProcessor.Output {
			log.Printf("--> Processed log: ID=%s | Source=%s | Severity=%s | Message='%s'\n",
				processedLog.ID, processedLog.Source, processedLog.Severity, processedLog.Message)
		}
	}()

	// 4. Configure o servidor Gin.
	router := gin.Default()

	// 5. Crie uma instância do handler de logs, passando o processador.
	logHandler := handler.NewLogHandler(logProcessor)

	// Defina as rotas.
	router.POST("/logs", logHandler.IngestLog)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "uptime": time.Since(time.Now()).String()})
	})

	// 6. Inicie o servidor.
	log.Println("Starting server on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

	// Feche os canais quando a aplicação for encerrada (usando Ctrl+C).
	defer logProcessor.Close()
}
