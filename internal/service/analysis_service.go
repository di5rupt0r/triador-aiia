package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/di5rupt0r/triador-aiia/internal/domain"
	"github.com/di5rupt0r/triador-aiia/internal/prompts"
)

// llmClient define o contrato que o serviço espera do cliente LLM.
// Interface privada — permite substituição por mock em testes.
type llmClient interface {
	Complete(ctx context.Context, systemPrompt, userMessage string) (string, error)
}

// llmResponse representa a saída esperada do LLM antes de mapear para o domínio.
// Tipo interno do serviço — não vaza para o domínio nem para o handler.
type llmResponse struct {
	CandidateName   string   `json:"candidate_name"`
	TechnicalSkills []string `json:"technical_skills"`
	YearsExperience float64  `json:"years_experience"`
	FitScore        int      `json:"fit_score"`
	Summary         string   `json:"summary"`
}

func (r *llmResponse) validate() error {
	if r.CandidateName == "" {
		return fmt.Errorf("candidate_name is empty")
	}
	if len(r.TechnicalSkills) == 0 {
		return fmt.Errorf("technical_skills is empty")
	}
	if r.FitScore < 0 || r.FitScore > 100 {
		return fmt.Errorf("fit_score %d is out of range [0, 100]", r.FitScore)
	}
	if r.Summary == "" {
		return fmt.Errorf("summary is empty")
	}
	return nil
}

type AnalysisService struct {
	llm llmClient
}

func NewAnalysisService(llm llmClient) *AnalysisService {
	return &AnalysisService{llm: llm}
}

func (s *AnalysisService) Analyze(ctx context.Context, req domain.AnalysisRequest) (*domain.AnalysisResult, error) {
	userMessage := prompts.BuildPrompt(req.Resume, req.JobDescription)

	raw, err := s.llm.Complete(ctx, prompts.SystemPrompt(), userMessage)
	if err != nil {
		return nil, fmt.Errorf("llm call failed: %w", err)
	}

	var llmResp llmResponse
	if err := json.Unmarshal([]byte(raw), &llmResp); err != nil {
		return nil, fmt.Errorf("llm response is not valid JSON: %w", err)
	}

	if err := llmResp.validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrValidation, err)
	}

	return &domain.AnalysisResult{
		CandidateName:   llmResp.CandidateName,
		Skills:          llmResp.TechnicalSkills,
		YearsExperience: llmResp.YearsExperience,
		FitScore:        llmResp.FitScore,
		Summary:         llmResp.Summary,
	}, nil
}
