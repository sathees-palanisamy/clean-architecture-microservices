package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/user/go-microservices/pkg/config"
	"github.com/user/go-microservices/pkg/logger"

	delivery "github.com/user/go-microservices/product-service/internal/delivery/http"
	repo "github.com/user/go-microservices/product-service/internal/infrastructure/db"
	"github.com/user/go-microservices/product-service/internal/usecase"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	defer logger.FromContext(context.Background()).Sync()

	log := logger.FromContext(context.Background())
	log.Info("Starting Product Service...")

	// Config
	serverPort := config.GetEnv("SERVER_PORT", "8081")

	// DB Connection
	dbConn, err := repo.NewConnection()

	if err != nil {
		log.Fatal("Could not connect to database", zap.Error(err))
	}
	defer dbConn.Close()

	// Simple Migration (For demo purposes)
	// In production, use a proper migration tool like golang-migrate
	schemaPath := "schema/001_init.sql"
	// Check if running in container or local might change path, simplistic check
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		log.Warn("Schema file not found (ok if already migrated or wrong path)")
	} else {
		content, err := ioutil.ReadFile(schemaPath)
		if err == nil {
			if _, err := dbConn.Exec(string(content)); err != nil {
				log.Error("Migration failed", zap.Error(err))
			} else {
				log.Info("Migration applied successfully")
			}
		}
	}

	// Layers
	productRepo := repo.NewPostgresRepository(dbConn)
	productUsecase := usecase.NewProductUsecase(productRepo, 2*time.Second)

	router := mux.NewRouter()
	delivery.NewProductHandler(router, productUsecase)

	// Server
	srv := &http.Server{
		Addr:    ":" + serverPort,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server start failed", zap.Error(err))
		}
	}()
	log.Info("Server started", zap.String("port", serverPort))

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}
	log.Info("Server exiting")
}
