package handler

import (
	"net/http"

	"searchengine/internal/model"
	"searchengine/pkg/httpx"
)

func (s *Server) registerSynonymRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/synonyms", s.createSynonym)
	mux.HandleFunc("GET /api/synonyms", s.listSynonyms)
	mux.HandleFunc("GET /api/synonyms/{id}", s.getSynonym)
	mux.HandleFunc("PUT /api/synonyms/{id}", s.updateSynonym)
	mux.HandleFunc("DELETE /api/synonyms/{id}", s.deleteSynonym)
}

type synonymRequest struct {
	Word     string   `json:"word"`
	Synonyms []string `json:"synonyms"`
	Status   string   `json:"status"`
}

func (s *Server) createSynonym(w http.ResponseWriter, r *http.Request) {
	var req synonymRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sy, err := s.svc.CreateSynonym(model.Synonym{Word: req.Word, Synonyms: req.Synonyms, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, sy)
}

func (s *Server) listSynonyms(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	items, total, err := s.svc.ListSynonyms(pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSynonym(w http.ResponseWriter, r *http.Request) {
	sy, err := s.svc.GetSynonym(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sy)
}

func (s *Server) updateSynonym(w http.ResponseWriter, r *http.Request) {
	var req synonymRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sy, err := s.svc.UpdateSynonym(r.PathValue("id"), model.Synonym{Word: req.Word, Synonyms: req.Synonyms, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sy)
}

func (s *Server) deleteSynonym(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteSynonym(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
