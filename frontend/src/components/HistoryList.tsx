import { Analysis } from "@/lib/api";

interface Props {
  analyses: Analysis[];
}

export default function HistoryList({ analyses }: Props) {
  if (analyses.length === 0) {
    return <p>Nenhuma análise realizada ainda.</p>;
  }

  return (
    <section>
      <h2>Histórico</h2>
      <ul>
        {analyses.map((a) => (
          <li key={a.id}>
            <strong>{a.candidate_name}</strong> — Fit: {a.fit_score}/100 —{" "}
            {new Date(a.created_at).toLocaleString("pt-BR")}
          </li>
        ))}
      </ul>
    </section>
  );
}
