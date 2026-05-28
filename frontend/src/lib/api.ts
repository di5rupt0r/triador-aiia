const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:9000";

export interface Analysis {
  id: number;
  candidate_name: string;
  skills: string[];
  years_experience: number;
  fit_score: number;
  summary: string;
  created_at: string;
}

export async function createAnalysis(
  resume: string,
  jobDescription: string
): Promise<Analysis> {
  const res = await fetch(`${API_URL}/analyses`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ resume, job_description: jobDescription }),
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error ?? "analysis failed");
  }
  return res.json();
}

export async function listAnalyses(): Promise<Analysis[]> {
  const res = await fetch(`${API_URL}/analyses`, { cache: "no-store" });
  if (!res.ok) throw new Error("failed to fetch analyses");
  return res.json();
}
