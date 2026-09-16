package handler

import (
	"net/http"

	"searchengine/pkg/httpx"
)

func (s *Server) registerHotSearchRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/hot-searches", s.listHotSearches)
	mux.HandleFunc("GET /api/hot-searches/{id}", s.getHotSearch)
	mux.HandleFunc("DELETE /api/hot-searches/{id}", s.deleteHotSearch)
}

func (s *Server) listHotSearches(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	items, total, err := s.svc.ListHotSearches(pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getHotSearch(w http.ResponseWriter, r *http.Request) {
	h, err := s.svc.GetHotSearch(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, h)
}

func (s *Server) deleteHotSearch(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteHotSearch(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
