package service

import (
	"context"
	"errors"
	"testing"

	"github.com/di5rupt0r/triador-aiia/internal/domain"
)

type mockLLM struct {
	response string
	err      error
}

func (m *mockLLM) Complete(_ context.Context, _, _ string) (string, error) {
	return m.response, m.err
}

func TestAnalyze_ValidResponse(t *testing.T) {
	svc := NewAnalysisService(&mockLLM{
		response: `{
			"candidate_name": "João Silva",
			"technical_skills": ["Go", "PostgreSQL"],
			"years_experience": 6,
			"fit_score": 87,
			"summary": "Forte aderência à vaga."
		}`,
	})

	result, err := svc.Analyze(context.Background(), domain.AnalysisRequest{
		Resume:         "...",
		JobDescription: "...",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.CandidateName != "João Silva" {
		t.Errorf("expected João Silva, got %s", result.CandidateName)
	}
	if result.FitScore != 87 {
		t.Errorf("expected fit_score 87, got %d", result.FitScore)
	}
}

func TestAnalyze_MalformedJSON(t *testing.T) {
	svc := NewAnalysisService(&mockLLM{
		response: `not json at all`,
	})

	_, err := svc.Analyze(context.Background(), domain.AnalysisRequest{
		Resume:         "...",
		JobDescription: "...",
	})
	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
}

func TestAnalyze_MissingCandidateName(t *testing.T) {
	svc := NewAnalysisService(&mockLLM{
		response: `{
			"candidate_name": "",
			"technical_skills": ["Go"],
			"years_experience": 3,
			"fit_score": 70,
			"summary": "Bom candidato."
		}`,
	})

	_, err := svc.Analyze(context.Background(), domain.AnalysisRequest{
		Resume:         "...",
		JobDescription: "...",
	})
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestAnalyze_FitScoreOutOfRange(t *testing.T) {
	svc := NewAnalysisService(&mockLLM{
		response: `{
			"candidate_name": "Ana Costa",
			"technical_skills": ["Python"],
			"years_experience": 2,
			"fit_score": 150,
			"summary": "Candidata promissora."
		}`,
	})

	_, err := svc.Analyze(context.Background(), domain.AnalysisRequest{
		Resume:         "...",
		JobDescription: "...",
	})
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected ErrValidation for fit_score out of range, got %v", err)
	}
}

func TestAnalyze_EmptySkills(t *testing.T) {
	svc := NewAnalysisService(&mockLLM{
		response: `{
			"candidate_name": "Carlos Lima",
			"technical_skills": [],
			"years_experience": 4,
			"fit_score": 60,
			"summary": "Perfil genérico."
		}`,
	})

	_, err := svc.Analyze(context.Background(), domain.AnalysisRequest{
		Resume:         "...",
		JobDescription: "...",
	})
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected ErrValidation for empty skills, got %v", err)
	}
}

func TestAnalyze_LLMError(t *testing.T) {
	svc := NewAnalysisService(&mockLLM{
		err: errors.New("connection refused"),
	})

	_, err := svc.Analyze(context.Background(), domain.AnalysisRequest{
		Resume:         "...",
		JobDescription: "...",
	})
	if err == nil {
		t.Fatal("expected error for LLM failure, got nil")
	}
	if errors.Is(err, domain.ErrValidation) {
		t.Fatal("LLM infra error should NOT be ErrValidation")
	}
}
