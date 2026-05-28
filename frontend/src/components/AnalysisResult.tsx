import { Analysis } from "@/lib/api";

interface Props {
  result: Analysis | null;
}

export default function AnalysisResult({ result }: Props) {
  if (!result) return null;

  return (
    <section>
      <h2>Resultado</h2>
      <p><strong>Candidato:</strong> {result.candidate_name}</p>
      <p><strong>Fit Score:</strong> {result.fit_score}/100</p>
      <p><strong>Experiência:</strong> {result.years_experience} anos</p>
      <p><strong>Habilidades:</strong> {result.skills.join(", ")}</p>
      <p><strong>Resumo:</strong> {result.summary}</p>
    </section>
  );
}
