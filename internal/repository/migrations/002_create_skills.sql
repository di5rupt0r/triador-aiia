CREATE TABLE IF NOT EXISTS analysis_skills (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    analysis_id INTEGER NOT NULL REFERENCES analyses(id) ON DELETE CASCADE,
    skill       TEXT    NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_analysis_skills_analysis_id
    ON analysis_skills(analysis_id);
