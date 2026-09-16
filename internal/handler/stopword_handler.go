package handler

import (
	"net/http"

	"searchengine/internal/model"
	"searchengine/pkg/httpx"
)

func (s *Server) registerStopWordRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/stop-words", s.createStopWord)
	mux.HandleFunc("GET /api/stop-words", s.listStopWords)
	mux.HandleFunc("DELETE /api/stop-words/{id}", s.deleteStopWord)
}

type stopWordRequest struct {
	Word string `json:"word"`
}

func (s *Server) createStopWord(w http.ResponseWriter, r *http.Request) {
	var req stopWordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sw, err := s.svc.CreateStopWord(model.StopWord{Word: req.Word})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, sw)
}

func (s *Server) listStopWords(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	items, total, err := s.svc.ListStopWords(pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) deleteStopWord(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteStopWord(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
