package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"

	"github.com/senatroxx/filmix-backend/internal/config"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all pending migrations",
	Run: func(cmd *cobra.Command, args []string) {
		m := initMigrate()
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("migrate up error:", err)
		}
		fmt.Println("Migration up complete")
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the last migration batch",
	Run: func(cmd *cobra.Command, args []string) {
		m := initMigrate()
		if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
			log.Fatal("migrate down error:", err)
		}
		fmt.Println("Migration down complete")
	},
}

var freshCmd = &cobra.Command{
	Use:   "fresh",
	Short: "Drop all tables and re-run all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		m := initMigrate()
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("migrate down error:", err)
		}
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("migrate up error:", err)
		}
		fmt.Println("Migration fresh complete")
	},
}

var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new migration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		migrationsPath := filepath.Join("internal", "database", "migrations")
		if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
			if err := os.MkdirAll(migrationsPath, os.ModePerm); err != nil {
				log.Fatal("failed to create migrations directory:", err)
			}
		}

		timestamp := time.Now().Format("20060102150405")
		upFile := filepath.Join(migrationsPath, fmt.Sprintf("%s_%s.up.sql", timestamp, name))
		downFile := filepath.Join(migrationsPath, fmt.Sprintf("%s_%s.down.sql", timestamp, name))

		if _, err := os.Create(upFile); err != nil {
			log.Fatal("failed to create up migration:", err)
		}
		if _, err := os.Create(downFile); err != nil {
			log.Fatal("failed to create down migration:", err)
		}

		fmt.Println("Created migration:", upFile, "and", downFile)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(upCmd, downCmd, freshCmd, createCmd)
}

func initMigrate() *migrate.Migrate {
	cfg := config.Load()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	m, err := migrate.New(
		"file://internal/database/migrations",
		dsn,
	)
	if err != nil {
		log.Fatal("migrate init error:", err)
	}
	return m
}
