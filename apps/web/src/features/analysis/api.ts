import { apiFetch } from "../../lib/api/client";
import type { AnalysisResult } from "../../components/CompletePanel";

export async function analyzeContribution(): Promise<AnalysisResult> {
  return apiFetch<AnalysisResult>("/analysis/contribution", {
    method: "POST",
  });
}
