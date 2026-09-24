package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found, using system environment variables")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize migration engine: %v", err)
	}
	defer m.Close()

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "up":
		log.Println("Applying migrations (up)...")
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Println("No new migrations to apply (database is already up to date).")
				return
			}
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("Successfully applied migrations!")

	case "down":
		log.Println("Rolling back last migration (down 1 step)...")
		if err := m.Steps(-1); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Println("No migrations to rollback.")
				return
			}
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("Successfully rolled back 1 migration step.")

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			if errors.Is(err, migrate.ErrNilVersion) {
				log.Println("No migrations have been applied yet.")
				return
			}
			log.Fatalf("Failed to get version: %v", err)
		}
		log.Printf("Current migration version: %d (dirty: %t)\n", version, dirty)

	case "force":
		forceCmd := flag.NewFlagSet("force", flag.ExitOnError)
		if len(os.Args) < 3 {
			log.Fatal("Usage: go run ./cmd/migrate force <version>")
		}
		var version int
		_, err := fmt.Sscanf(os.Args[2], "%d", &version)
		if err != nil {
			log.Fatalf("Invalid version number: %v", os.Args[2])
		}
		_ = forceCmd
		if err := m.Force(version); err != nil {
			log.Fatalf("Failed to force version: %v", err)
		}
		log.Printf("Forced migration version to: %d\n", version)

	default:
		log.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Usage:
  go run ./cmd/migrate <command>

Commands:
  up       Apply all pending migrations
  down     Roll back the most recent migration (1 step)
  version  Show current applied migration version
  force <v> Force set migration version (useful if dirty state occurs)`)
}
