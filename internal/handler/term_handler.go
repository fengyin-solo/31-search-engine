package handler

import (
	"net/http"

	"searchengine/pkg/httpx"
)

func (s *Server) registerTermRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/terms", s.listTerms)
	mux.HandleFunc("GET /api/terms/{id}", s.getTerm)
}

func (s *Server) listTerms(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	indexID := r.URL.Query().Get("index_id")
	items, total, err := s.svc.ListTerms(indexID, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getTerm(w http.ResponseWriter, r *http.Request) {
	t, err := s.svc.GetTerm(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}
