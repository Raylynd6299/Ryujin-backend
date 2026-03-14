package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Raylynd6299/Ryujin-backend/internal/config"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("✓ Configuration loaded")

	// Initialize dependencies
	deps, err := NewAppDependencies(&config.App)
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}
	defer func() {
		if err := deps.Close(); err != nil {
			log.Printf("Error closing dependencies: %v", err)
		}
	}()

	// Setup router
	router := SetupRouter(deps)

	// Configure HTTP server
	srv := &http.Server{
		Addr:    ":" + config.App.Server.Port,
		Handler: router,
	}

	// Worker context — cancelled on graceful shutdown
	workerCtx, cancelWorker := context.WithCancel(context.Background())
	defer cancelWorker()

	// Start background worker
	deps.PriceRefreshWorker.Start(workerCtx)
	log.Println("✓ Price refresh worker started")

	// Start server in a goroutine
	go func() {
		log.Printf("✓ Server starting on port %s", config.App.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for graceful shutdown
	gracefulShutdown(srv, cancelWorker)
}

// gracefulShutdown waits for interrupt signals and shuts down the server gracefully
func gracefulShutdown(srv *http.Server, cancelWorker context.CancelFunc) {
	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Stop the background worker first
	cancelWorker()

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✓ Server stopped gracefully")
}
