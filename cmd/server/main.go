package main

import (
	"log"

	"github.com/di5rupt0r/triador-aiia/config"
	"github.com/di5rupt0r/triador-aiia/internal/repository"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	db, err := repository.NewSQLiteDB(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
	defer db.Close()

	_ = db
}
