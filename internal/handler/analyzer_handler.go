package handler

import (
	"net/http"

	"searchengine/internal/model"
	"searchengine/pkg/httpx"
)

func (s *Server) registerAnalyzerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/analyzers", s.createAnalyzer)
	mux.HandleFunc("GET /api/analyzers", s.listAnalyzers)
	mux.HandleFunc("GET /api/analyzers/{id}", s.getAnalyzer)
	mux.HandleFunc("PUT /api/analyzers/{id}", s.updateAnalyzer)
	mux.HandleFunc("DELETE /api/analyzers/{id}", s.deleteAnalyzer)
}

type analyzerRequest struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	StopWords []string `json:"stop_words"`
	Status    string   `json:"status"`
}

func (s *Server) createAnalyzer(w http.ResponseWriter, r *http.Request) {
	var req analyzerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.CreateAnalyzer(model.Analyzer{
		Name: req.Name, Type: req.Type, StopWords: req.StopWords, Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, a)
}

func (s *Server) listAnalyzers(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	items, total, err := s.svc.ListAnalyzers(pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAnalyzer(w http.ResponseWriter, r *http.Request) {
	a, err := s.svc.GetAnalyzer(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) updateAnalyzer(w http.ResponseWriter, r *http.Request) {
	var req analyzerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.UpdateAnalyzer(r.PathValue("id"), model.Analyzer{
		Name: req.Name, Type: req.Type, StopWords: req.StopWords, Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) deleteAnalyzer(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteAnalyzer(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
