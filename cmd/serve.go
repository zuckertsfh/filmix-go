package cmd

import (
	"github.com/senatroxx/filmix-backend/internal/config"
	"github.com/senatroxx/filmix-backend/internal/database"
	"github.com/senatroxx/filmix-backend/internal/http"
	"github.com/senatroxx/filmix-backend/internal/utilities"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		utilities.InitLogger(cfg.Mode)

		db, err := database.Connect(&cfg.Database)
		if err != nil {
			utilities.Logger.Fatal().Err(err).Msg("Failed to connect to database")
		}
		defer db.Close()

		hr := config.InitializeHandlers(config.InitializeServices(config.InitializeRepositories(db)))
		srv := http.InitializeAPI(&cfg, hr, db, utilities.Logger)
		srv.Run()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
