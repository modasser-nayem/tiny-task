package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/modasser-nayem/tiny-task/internal/config"
	"github.com/modasser-nayem/tiny-task/internal/database"
	"github.com/modasser-nayem/tiny-task/internal/modules/auth"
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

		r := router.Setup(cfg, authHandler)

		log.Println("Server running on", cfg.Port)

		if err := r.Run(":" + cfg.Port); err != nil {
			log.Fatal(err)
		}
}