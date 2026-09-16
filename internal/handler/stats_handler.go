package handler

import (
	"net/http"

	"searchengine/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/overview", s.statsOverview)
	mux.HandleFunc("GET /api/stats/indexes", s.statsIndexes)
	mux.HandleFunc("GET /api/stats/queries", s.statsQueries)
	mux.HandleFunc("GET /api/stats/terms", s.statsTerms)
	mux.HandleFunc("GET /api/stats/categories", s.statsCategories)
}

func (s *Server) statsCategories(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, map[string]interface{}{
		"by_category": s.svc.CategoryDocStats(),
		"by_source":   s.svc.SourceDocStats(),
	})
}

func (s *Server) statsOverview(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.Overview())
}

func (s *Server) statsIndexes(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, map[string]interface{}{
		"indexes": s.svc.IndexStats(),
	})
}

func (s *Server) statsQueries(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, map[string]interface{}{
		"top_query_terms":  s.svc.TopQueryTerms(10),
		"avg_duration_ms":  s.svc.AvgQueryDuration(),
		"top_hot_searches": s.svc.TopHotSearches(10),
	})
}

func (s *Server) statsTerms(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, map[string]interface{}{
		"top_terms": s.svc.TopTerms(20),
	})
}
