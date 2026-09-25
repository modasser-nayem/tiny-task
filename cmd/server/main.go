package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/modasser-nayem/tiny-task/internal/config"
	"github.com/modasser-nayem/tiny-task/internal/database"
	"github.com/modasser-nayem/tiny-task/internal/modules/auth"
	"github.com/modasser-nayem/tiny-task/internal/modules/todo"
	"github.com/modasser-nayem/tiny-task/internal/modules/user"
	"github.com/modasser-nayem/tiny-task/internal/router"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}
		defer db.Close()

		authRepository := auth.NewPostgresRepository(db)

		tokenManager := auth.NewTokenManager(cfg.JWTSecret)

		authService := auth.NewService(authRepository, tokenManager)

		authHandler := auth.NewHandler(authService)

		userHandler := user.NewHandler()

		todoRepository := todo.NewPostgresRepository(db)
		todoService := todo.NewService(todoRepository)
		todoHandler := todo.NewHandler(todoService)

		r := router.Setup(
			cfg,
			authHandler,
			userHandler,
			todoHandler,
			tokenManager,
		)

		log.Println("Server running on", cfg.Port)

		if err := r.Run(":" + cfg.Port); err != nil {
			log.Fatal(err)
		}
}