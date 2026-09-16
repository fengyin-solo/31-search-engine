package handler

import (
	"net/http"

	"searchengine/pkg/httpx"
)

func (s *Server) registerQueryLogRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/query-logs", s.listQueryLogs)
	mux.HandleFunc("GET /api/query-logs/{id}", s.getQueryLog)
	mux.HandleFunc("DELETE /api/query-logs/{id}", s.deleteQueryLog)
}

func (s *Server) listQueryLogs(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	indexID := r.URL.Query().Get("index_id")
	items, total, err := s.svc.ListQueryLogs(indexID, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getQueryLog(w http.ResponseWriter, r *http.Request) {
	q, err := s.svc.GetQueryLog(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, q)
}

func (s *Server) deleteQueryLog(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteQueryLog(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
