package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/di5rupt0r/triador-aiia/internal/domain"
)

type AnalysisRepository struct {
	db *sql.DB
}

func NewAnalysisRepository(db *sql.DB) *AnalysisRepository {
	return &AnalysisRepository{db: db}
}

func (r *AnalysisRepository) Save(ctx context.Context, a *domain.Analysis) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `
		INSERT INTO analyses (candidate_name, years_experience, fit_score, summary, raw_resume, job_description)
		VALUES (?, ?, ?, ?, ?, ?)`,
		a.CandidateName,
		a.YearsExperience,
		a.FitScore,
		a.Summary,
		a.RawResume,
		a.JobDescription,
	)
	if err != nil {
		return fmt.Errorf("insert analysis: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}
	a.ID = id

	for _, skill := range a.Skills {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO analysis_skills (analysis_id, skill) VALUES (?, ?)`,
			id, skill,
		); err != nil {
			return fmt.Errorf("insert skill %q: %w", skill, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	var createdAtUnix int64
	if err := r.db.QueryRowContext(ctx, `SELECT created_at FROM analyses WHERE id = ?`, a.ID).Scan(&createdAtUnix); err != nil {
		return fmt.Errorf("read created_at: %w", err)
	}
	a.CreatedAt = time.Unix(createdAtUnix, 0).UTC()

	return nil
}

func (r *AnalysisRepository) List(ctx context.Context) ([]domain.Analysis, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, candidate_name, years_experience, fit_score, summary,
		       raw_resume, job_description, created_at
		FROM analyses
		ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("query analyses: %w", err)
	}
	defer rows.Close()

	var analyses []domain.Analysis

	for rows.Next() {
		var a domain.Analysis
		var createdAtUnix int64

		if err := rows.Scan(
			&a.ID,
			&a.CandidateName,
			&a.YearsExperience,
			&a.FitScore,
			&a.Summary,
			&a.RawResume,
			&a.JobDescription,
			&createdAtUnix,
		); err != nil {
			return nil, fmt.Errorf("scan analysis row: %w", err)
		}

		a.CreatedAt = time.Unix(createdAtUnix, 0).UTC()

		skills, err := r.fetchSkills(ctx, a.ID)
		if err != nil {
			return nil, err
		}
		a.Skills = skills

		analyses = append(analyses, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate analysis rows: %w", err)
	}

	return analyses, nil
}

func (r *AnalysisRepository) fetchSkills(ctx context.Context, analysisID int64) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT skill FROM analysis_skills
		WHERE analysis_id = ?
		ORDER BY id ASC`,
		analysisID,
	)
	if err != nil {
		return nil, fmt.Errorf("query skills for analysis %d: %w", analysisID, err)
	}
	defer rows.Close()

	var skills []string

	for rows.Next() {
		var skill string
		if err := rows.Scan(&skill); err != nil {
			return nil, fmt.Errorf("scan skill row: %w", err)
		}
		skills = append(skills, skill)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate skill rows: %w", err)
	}

	return skills, nil
}
