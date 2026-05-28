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
			"candidate_name": "Ada Lovelace",
			"technical_skills": ["Go", "Python"],
			"years_experience": 5.0,
			"fit_score": 87,
			"summary": "Strong candidate with relevant experience."
		}`,
	})

	result, err := svc.Analyze(context.Background(), domain.AnalysisRequest{
		Resume:         "resume text",
		JobDescription: "job text",
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.CandidateName != "Ada Lovelace" {
		t.Errorf("expected Ada Lovelace, got %s", result.CandidateName)
	}
	if result.FitScore != 87 {
		t.Errorf("expected fit_score 87, got %d", result.FitScore)
	}
}

func TestAnalyze_MalformedJSON(t *testing.T) {
	svc := NewAnalysisService(&mockLLM{
		response: `this is not json {{{`,
	})

	_, err := svc.Analyze(context.Background(), domain.AnalysisRequest{
		Resume:         "resume text",
		JobDescription: "job text",
	})

	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
}

func TestAnalyze_InvalidFitScore(t *testing.T) {
	svc := NewAnalysisService(&mockLLM{
		response: `{
			"candidate_name": "Ada Lovelace",
			"technical_skills": ["Go"],
			"years_experience": 3.0,
			"fit_score": 150,
			"summary": "Some summary."
		}`,
	})

	_, err := svc.Analyze(context.Background(), domain.AnalysisRequest{
		Resume:         "resume text",
		JobDescription: "job text",
	})

	if err == nil {
		t.Fatal("expected error for fit_score out of range, got nil")
	}
}

func TestAnalyze_LLMProviderError(t *testing.T) {
	svc := NewAnalysisService(&mockLLM{
		err: errors.New("connection timeout"),
	})

	_, err := svc.Analyze(context.Background(), domain.AnalysisRequest{
		Resume:         "resume text",
		JobDescription: "job text",
	})

	if err == nil {
		t.Fatal("expected error for LLM provider failure, got nil")
	}
}
