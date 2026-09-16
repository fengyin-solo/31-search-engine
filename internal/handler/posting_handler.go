package handler

import (
	"net/http"

	"searchengine/pkg/httpx"
)

func (s *Server) registerPostingRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/postings", s.listPostings)
	mux.HandleFunc("GET /api/postings/{id}", s.getPosting)
}

func (s *Server) listPostings(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	indexID := r.URL.Query().Get("index_id")
	term := r.URL.Query().Get("term")
	docID := r.URL.Query().Get("doc_id")
	items, total, err := s.svc.ListPostings(indexID, term, docID, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getPosting(w http.ResponseWriter, r *http.Request) {
	p, err := s.svc.GetPosting(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}
