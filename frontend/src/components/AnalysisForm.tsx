"use client";

import { useState } from "react";
import { createAnalysis, Analysis } from "@/lib/api";

interface Props {
  onSuccess: (result: Analysis) => void;
}

export default function AnalysisForm({ onSuccess }: Props) {
  const [resume, setResume] = useState("");
  const [jobDescription, setJobDescription] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError(null);
    try {
      const result = await createAnalysis(resume, jobDescription);
      onSuccess(result);
      setResume("");
      setJobDescription("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "unknown error");
    } finally {
      setLoading(false);
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label htmlFor="resume">Currículo</label>
        <textarea
          id="resume"
          value={resume}
          onChange={(e) => setResume(e.target.value)}
          rows={8}
          required
          disabled={loading}
        />
      </div>
      <div>
        <label htmlFor="job">Descrição da Vaga</label>
        <textarea
          id="job"
          value={jobDescription}
          onChange={(e) => setJobDescription(e.target.value)}
          rows={5}
          required
          disabled={loading}
        />
      </div>
      {error && <p style={{ color: "red" }}>{error}</p>}
      <button type="submit" disabled={loading}>
        {loading ? "Analisando..." : "Analisar"}
      </button>
    </form>
  );
}
