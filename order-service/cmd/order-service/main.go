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

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/user/go-microservices/order-service/docs" // Generated docs
	delivery "github.com/user/go-microservices/order-service/internal/delivery/http"
	client "github.com/user/go-microservices/order-service/internal/infrastructure/client"
	repo "github.com/user/go-microservices/order-service/internal/infrastructure/db"
	"github.com/user/go-microservices/order-service/internal/usecase"
	"go.uber.org/zap"
)

// @title Order Service API
// @version 1.0
// @description This is an order management service.
// @host localhost:8082
// @BasePath /
func main() {
	logger.Init()
	defer logger.FromContext(context.Background()).Sync()

	log := logger.FromContext(context.Background())
	log.Info("Starting Order Service...")

	// Config
	serverPort := config.GetEnv("SERVER_PORT", "8082")
	productServiceURL := config.GetEnv("PRODUCT_SERVICE_URL", "http://localhost:8081")

	// DB Connection
	dbConn, err := repo.NewConnection()

	if err != nil {
		log.Fatal("Could not connect to database", zap.Error(err))
	}
	defer dbConn.Close()

	// Migration
	schemaPath := "schema/001_init.sql"
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		log.Warn("Schema file not found")
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
	orderRepo := repo.NewOrderRepository(dbConn)
	prodClient := client.NewProductClient(productServiceURL)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, prodClient, 5*time.Second)

	router := mux.NewRouter()
	delivery.NewOrderHandler(router, orderUsecase)

	// Swagger UI
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

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
