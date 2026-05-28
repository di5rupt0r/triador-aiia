package main

import (
	"log"

	"github.com/di5rupt0r/triador-aiia/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}
	_ = cfg
}
