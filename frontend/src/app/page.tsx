"use client";

import { useState, useEffect } from "react";
import { listAnalyses, Analysis } from "@/lib/api";
import AnalysisForm from "@/components/AnalysisForm";
import AnalysisResult from "@/components/AnalysisResult";
import HistoryList from "@/components/HistoryList";

export default function Home() {
  const [result, setResult] = useState<Analysis | null>(null);
  const [history, setHistory] = useState<Analysis[]>([]);

  useEffect(() => {
    listAnalyses().then(setHistory).catch(console.error);
  }, []);

  function handleSuccess(analysis: Analysis) {
    setResult(analysis);
    setHistory((prev) => [analysis, ...prev]);
  }

  return (
    <main>
      <h1>Triador de Currículos</h1>
      <AnalysisForm onSuccess={handleSuccess} />
      <AnalysisResult result={result} />
      <HistoryList analyses={history} />
    </main>
  );
}
