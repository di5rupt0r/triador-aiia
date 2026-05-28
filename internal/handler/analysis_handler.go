package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/di5rupt0r/triador-aiia/internal/domain"
)

type analysisService interface {
	Analyze(ctx context.Context, req domain.AnalysisRequest) (*domain.AnalysisResult, error)
}

type analysisRepository interface {
	Save(ctx context.Context, a *domain.Analysis) error
	List(ctx context.Context) ([]domain.Analysis, error)
}

type AnalysisHandler struct {
	svc  analysisService
	repo analysisRepository
}

func NewAnalysisHandler(svc analysisService, repo analysisRepository) *AnalysisHandler {
	return &AnalysisHandler{svc: svc, repo: repo}
}

type createRequest struct {
	Resume         string `json:"resume"`
	JobDescription string `json:"job_description"`
}

type analysisResponse struct {
	ID              int64    `json:"id"`
	CandidateName   string   `json:"candidate_name"`
	Skills          []string `json:"skills"`
	YearsExperience float64  `json:"years_experience"`
	FitScore        int      `json:"fit_score"`
	Summary         string   `json:"summary"`
	CreatedAt       string   `json:"created_at"`
}

func toResponse(a *domain.Analysis) analysisResponse {
	skills := a.Skills
	if skills == nil {
		skills = []string{}
	}
	return analysisResponse{
		ID:              a.ID,
		CandidateName:   a.CandidateName,
		Skills:          skills,
		YearsExperience: a.YearsExperience,
		FitScore:        a.FitScore,
		Summary:         a.Summary,
		CreatedAt:       a.CreatedAt.Format(time.RFC3339),
	}
}

func (h *AnalysisHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Resume == "" || req.JobDescription == "" {
		writeError(w, http.StatusBadRequest, "resume and job_description are required")
		return
	}

	result, err := h.svc.Analyze(r.Context(), domain.AnalysisRequest{
		Resume:         req.Resume,
		JobDescription: req.JobDescription,
	})
	if err != nil {
		if errors.Is(err, domain.ErrValidation) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		} else {
			log.Printf("ERROR create analysis: %v", err)
			writeError(w, http.StatusInternalServerError, "analysis failed")
		}
		return
	}

	analysis := domain.Analysis{
		CandidateName:   result.CandidateName,
		Skills:          result.Skills,
		YearsExperience: result.YearsExperience,
		FitScore:        result.FitScore,
		Summary:         result.Summary,
		RawResume:       req.Resume,
		JobDescription:  req.JobDescription,
	}

	if err := h.repo.Save(r.Context(), &analysis); err != nil {
		log.Printf("ERROR persist analysis: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to persist analysis")
		return
	}

	writeJSON(w, http.StatusCreated, toResponse(&analysis))
}

func (h *AnalysisHandler) List(w http.ResponseWriter, r *http.Request) {
	analyses, err := h.repo.List(r.Context())
	if err != nil {
		log.Printf("ERROR list analyses: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to fetch analyses")
		return
	}

	resp := make([]analysisResponse, 0, len(analyses))
	for i := range analyses {
		resp = append(resp, toResponse(&analyses[i]))
	}

	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
