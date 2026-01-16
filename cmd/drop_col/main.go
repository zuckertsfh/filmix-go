package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/senatroxx/filmix-backend/internal/config"
	"github.com/senatroxx/filmix-backend/internal/database"
)

func main() {
	cfg := config.Load()
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Alter table to drop column
	_, err = db.Exec("ALTER TABLE users DROP COLUMN IF EXISTS phone_number")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("phone_number column dropped successfully.")
}
