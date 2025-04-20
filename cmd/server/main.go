package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"smolink/internal/config"
	"smolink/internal/controller"
	"smolink/internal/migration"
	"smolink/internal/repository"
	"smolink/internal/service"
	"smolink/pkg/database"
	"smolink/pkg/logger"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error loading config: %v", err)
	}

	// Set Gin mode based on config
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	pgDB, err := database.NewPostgresDB(cfg.PostgresDSN)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	defer pgDB.Close()

	// Run migrations
	if err := migration.RunMigrations(pgDB.Pool); err != nil {
		log.Fatal("Migration error:", err)
	}

	redisClient, err := database.NewRedisDB(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatal("Redis connection error:", err)

	}

	pgRepo := repository.NewPostgresRepository(pgDB.Pool)
	redisRepo := repository.NewRedisRepository(redisClient.Client)
	urlService := service.NewURLService(pgRepo, redisRepo)
	urlController := controller.NewURLController(urlService)

	router := gin.New()
	router.Use(
		gin.Recovery(),      // panic recovery
		logger.Middleware(), // custom logging
	)

	router.POST("/shorten", urlController.ShortenURL)
	router.GET("/:code", urlController.ResolveURL)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Ah, you don reach home. Welcome to the smolink service. We dey for you!",
		})
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Omo, this route no dey exist o! You don miss road. Go back jare!",
		})
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"message": "Abeg, no try that method here. Na wrong move!",
		})
	})

	// Start the server
	server := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: router,
	}

	go func() {
		log.Println("Starting server on", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error:", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited properly")
}
