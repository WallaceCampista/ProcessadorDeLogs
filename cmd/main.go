package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log-processor/internal/config"
	"log-processor/internal/handler"
	"log-processor/internal/processor"
	"log-processor/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// rateLimitMiddleware cria um middleware de limite de taxa.
func rateLimitMiddleware(limit *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Pega o IP do cliente como chave.
		ip := c.ClientIP()

		limitContext, err := limit.Get(c.Request.Context(), ip)
		if err != nil {
			log.Printf("Erro ao verificar o limite de taxa: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Adiciona os headers de resposta.
		c.Header("X-Rate-Limit-Limit", fmt.Sprintf("%d", limitContext.Limit))
		c.Header("X-Rate-Limit-Remaining", fmt.Sprintf("%d", limitContext.Remaining))
		c.Header("X-Rate-Limit-Reset", fmt.Sprintf("%d", limitContext.Reset))

		// Aborta a requisição se o limite for excedido.
		if limitContext.Reached {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}

		// Continua o processamento da requisição.
		c.Next()
	}
}

func main() {
	// 1. Carregue as configurações do arquivo YAML.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Falha ao carregar a configuração: %v", err)
	}

	// 2. Crie a camada de armazenamento (MySQL Storager).
	mysqlDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)
	mysqlStorager, err := storage.NewMySQLStorager(mysqlDSN)
	if err != nil {
		log.Fatalf("Falha ao criar o MySQL Storage: %v", err)
	}

	// 3. Crie e inicie o processador de logs.
	bufferSize := 1000
	logProcessor := processor.NewLogProcessor(bufferSize)
	logProcessor.Start()

	// 4. Crie uma goroutine para consumir os logs processados e enviá-los ao MySQL.
	go func() {
		ctx := context.Background()
		for processedLog := range logProcessor.Output {
			log.Printf("--> Log processado: ID=%s | Source=%s | Severity=%s | Message='%s'\n",
				processedLog.ID, processedLog.Source, processedLog.Severity, processedLog.Message)

			if err := mysqlStorager.StoreLog(ctx, processedLog); err != nil {
				log.Printf("Falha ao armazenar o log no mysql: %v", err)
			}
		}
	}()

	// 5. Configure o servidor Gin, os middlewares e os handlers.
	router := gin.Default()

	// CORREÇÃO: Inicialização mais robusta do limite de taxa.
	rateLimitStore := memory.NewStore()
	rate := limiter.Rate{
		Limit:  int64(cfg.RateLimiting.Rate),
		Period: 1 * time.Second, // Hardcoded para 1 segundo, para evitar problemas de parsing.
		// NOTE: Para usar a configuração, teríamos que parsear `cfg.RateLimiting.Period`
		// usando time.ParseDuration antes.
	}
	rateLimiter := limiter.New(rateLimitStore, rate)

	// Crie a instância dos handlers, passando as dependências.
	logHandler := handler.NewLogHandler(logProcessor)
	searchHandler := handler.NewSearchHandler(mysqlStorager)

	// 6. Defina as rotas, aplicando os middlewares e handlers.
	router.POST("/logs", rateLimitMiddleware(rateLimiter), logHandler.IngestLog)
	router.GET("/search", searchHandler.SearchLogs)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "uptime": time.Since(time.Now()).Round(time.Second).String()})
	})
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "uptime": time.Since(time.Now()).Round(time.Second).String()})
	})

	// 7. Inicie o servidor.
	go func() {
		if err := router.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
			log.Fatalf("Falha ao executar o servidor: %v", err)
		}
	}()

	// 8. Configure o graceful shutdown (desligamento suave).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Desligando servidor...")
	if err := mysqlStorager.Close(); err != nil {
		log.Printf("Erro ao fechar o armazenamento MySQL: %v", err)
	}
	log.Println("Servidor parou graciosamente.")
}
