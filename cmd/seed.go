package cmd

import (
	"log"

	"github.com/senatroxx/filmix-backend/internal/config"
	"github.com/senatroxx/filmix-backend/internal/database"
	"github.com/senatroxx/filmix-backend/internal/database/seeder"
	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed database with initial data",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()

		// Connect DB
		db, err := database.Connect(&cfg.Database)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()

		s := seeder.NewSeeder(db, &cfg)

		if err := s.SeedAll(); err != nil {
			log.Fatalf("Failed to seed database: %v", err)
		}

		log.Println("Database seeded successfully!")
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}
