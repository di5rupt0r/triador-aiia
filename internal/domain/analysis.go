package domain

import "time"

// Analysis representa uma análise persistida no banco.
type Analysis struct {
	ID              int64
	CandidateName   string
	Skills          []string
	YearsExperience float64
	FitScore        int
	Summary         string
	RawResume       string
	JobDescription  string
	CreatedAt       time.Time
}

// AnalysisRequest é o input recebido do usuário para criar uma nova análise.
type AnalysisRequest struct {
	Resume         string
	JobDescription string
}

// AnalysisResult é a saída estruturada retornada pelo LLM, antes de persistir.
type AnalysisResult struct {
	CandidateName   string
	Skills          []string
	YearsExperience float64
	FitScore        int
	Summary         string
}
