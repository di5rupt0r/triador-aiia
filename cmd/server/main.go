package main

import (
	"log"
	"net/http"

	"github.com/di5rupt0r/triador-aiia/config"
	"github.com/di5rupt0r/triador-aiia/internal/handler"
	"github.com/di5rupt0r/triador-aiia/internal/llm"
	"github.com/di5rupt0r/triador-aiia/internal/repository"
	"github.com/di5rupt0r/triador-aiia/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := repository.NewSQLiteDB(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	repo := repository.NewAnalysisRepository(db)
	llmClient := llm.NewOpenAIClient(cfg.OpenAIAPIKey, cfg.OpenAIBaseURL, cfg.OpenAIModel)
	svc := service.NewAnalysisService(llmClient)
	h := handler.NewAnalysisHandler(svc, repo)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /analyses", h.Create)
	mux.HandleFunc("GET /analyses", h.List)

	log.Printf("server listening on :%s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, handler.CORSMiddleware(mux)); err != nil {
		log.Fatalf("server: %v", err)
	}
}
