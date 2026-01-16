package config

import (
	"database/sql"

	"github.com/senatroxx/filmix-backend/internal/http/handlers"
	"github.com/senatroxx/filmix-backend/internal/repositories"
	"github.com/senatroxx/filmix-backend/internal/services"
)

func InitializeHandlers(s *services.Services) *handlers.Handlers {
	return handlers.RegisterHandlers(s)
}

func InitializeRepositories(db *sql.DB) *repositories.Repositories {
	return repositories.RegisterRepositories(db)
}

func InitializeServices(r *repositories.Repositories) *services.Services {
	return services.RegisterServices(r)
}
