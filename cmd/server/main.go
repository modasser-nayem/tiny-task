package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/modasser-nayem/tiny-task/internal/config"
	"github.com/modasser-nayem/tiny-task/internal/database"
	"github.com/modasser-nayem/tiny-task/internal/modules/auth"
	"github.com/modasser-nayem/tiny-task/internal/modules/todo"
	"github.com/modasser-nayem/tiny-task/internal/modules/user"
	"github.com/modasser-nayem/tiny-task/internal/router"
)

func main() {
	// Initialize structured logger
	logger := setupLogger()
	slog.SetDefault(logger)

	// 1. Load and validate configuration.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// 2. Connect to PostgreSQL.
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	// 3. Build dependencies.
	authRepository := auth.NewPostgresRepository(db)
	tokenManager := auth.NewTokenManager(cfg.JWTSecret)
	authService := auth.NewService(authRepository, tokenManager)
	authHandler := auth.NewHandler(authService)

	todoRepository := todo.NewPostgresRepository(db)
	todoService := todo.NewService(todoRepository)
	todoHandler := todo.NewHandler(todoService)

	userHandler := user.NewHandler()

	// 4. Set Gin mode.
	gin.SetMode(gin.ReleaseMode)

	// 5. Create router.
	r := router.Setup(cfg, authHandler, userHandler, todoHandler, tokenManager)

	// 6. Configure the HTTP server.
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// 7. Start server in a separate goroutine.
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Server listening on port %s", cfg.Port)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	// 8. Wait for OS shutdown signal or server failure.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case <-ctx.Done():
		log.Println("Shutdown signal received")
	case err := <-serverErrors:
		log.Printf("HTTP server failed: %v", err)
	}

	// 9. Give active requests time to finish.
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)

		if closeErr := server.Close(); closeErr != nil {
			log.Printf("Force close failed: %v", closeErr)
		}
	}

	log.Println("Server stopped")
}

func setupLogger() *slog.Logger {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	return logger
}