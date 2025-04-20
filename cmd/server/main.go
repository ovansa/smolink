package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"smolink/internal/app"
	"smolink/internal/config"
	"smolink/internal/migration"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	appInstance, err := app.NewApp(cfg, true)
	if err != nil {
		log.Fatalf("App init failed: %v", err)
	}
	defer appInstance.DBCloser()

	if err := migration.RunMigrations(appInstance.PGRepo.DB()); err != nil {
		log.Fatal("Migration error:", err)
	}

	server := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: appInstance.Router,
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
