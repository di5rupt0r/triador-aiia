CREATE TABLE IF NOT EXISTS analyses (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    candidate_name   TEXT    NOT NULL,
    years_experience REAL    NOT NULL,
    fit_score        INTEGER NOT NULL CHECK (fit_score BETWEEN 0 AND 100),
    summary          TEXT    NOT NULL,
    raw_resume       TEXT    NOT NULL,
    job_description  TEXT    NOT NULL,
    created_at       INTEGER NOT NULL DEFAULT (unixepoch())
);
