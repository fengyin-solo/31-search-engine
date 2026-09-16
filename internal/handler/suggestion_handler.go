package handler

import (
	"net/http"
	"strconv"

	"searchengine/internal/model"
	"searchengine/pkg/httpx"
)

func (s *Server) registerSuggestionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/suggestions", s.createSuggestion)
	mux.HandleFunc("GET /api/suggestions", s.listSuggestions)
	mux.HandleFunc("GET /api/suggestions/{id}", s.getSuggestion)
	mux.HandleFunc("PUT /api/suggestions/{id}", s.updateSuggestion)
	mux.HandleFunc("DELETE /api/suggestions/{id}", s.deleteSuggestion)
	mux.HandleFunc("GET /api/suggest", s.suggest)
}

type suggestionRequest struct {
	Term   string `json:"term"`
	Weight int    `json:"weight"`
}

func (s *Server) createSuggestion(w http.ResponseWriter, r *http.Request) {
	var req suggestionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	v, err := s.svc.CreateSuggestion(model.Suggestion{Term: req.Term, Weight: req.Weight})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, v)
}

func (s *Server) listSuggestions(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	items, total, err := s.svc.ListSuggestions(pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSuggestion(w http.ResponseWriter, r *http.Request) {
	v, err := s.svc.GetSuggestion(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, v)
}

func (s *Server) updateSuggestion(w http.ResponseWriter, r *http.Request) {
	var req suggestionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	v, err := s.svc.UpdateSuggestion(r.PathValue("id"), model.Suggestion{Term: req.Term, Weight: req.Weight})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, v)
}

func (s *Server) deleteSuggestion(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteSuggestion(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) suggest(w http.ResponseWriter, r *http.Request) {
	prefix := r.URL.Query().Get("prefix")
	limit := 10
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	items := s.svc.Suggest(prefix, limit)
	httpx.OK(w, map[string]interface{}{"items": items})
}
