package handler

import (
	"net/http"

	"searchengine/pkg/httpx"
)

func (s *Server) registerExportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/export/summary", s.exportSummary)
}

// exportSummary 导出系统全量快照汇总。
func (s *Server) exportSummary(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, map[string]interface{}{
		"overview":     s.svc.Overview(),
		"documents":    s.svc.ListAllDocumentsForExport(),
		"indexes":      s.svc.ListAllIndexesForExport(),
		"analyzers":    s.svc.ListAllAnalyzersForExport(),
		"stop_words":   s.svc.ListAllStopWordsForExport(),
		"terms":        s.svc.ListAllTermsForExport(),
		"postings":     s.svc.ListAllPostingsForExport(),
		"synonyms":     s.svc.ListAllSynonymsForExport(),
		"query_logs":   s.svc.ListAllQueryLogsForExport(),
		"hot_searches": s.svc.ListAllHotSearchesForExport(),
		"categories":   s.svc.ListAllCategoriesForExport(),
		"suggestions":  s.svc.ListAllSuggestionsForExport(),
		"spell_corrections": s.svc.ListAllSpellCorrectionsForExport(),
		"boost_rules":  s.svc.ListAllBoostRulesForExport(),
	})
}
