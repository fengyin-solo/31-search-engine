package handler

import (
	"net/http"

	"searchengine/internal/model"
	"searchengine/pkg/httpx"
)

func (s *Server) registerSpellCorrectionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/spell-corrections", s.createSpellCorrection)
	mux.HandleFunc("GET /api/spell-corrections", s.listSpellCorrections)
	mux.HandleFunc("GET /api/spell-corrections/{id}", s.getSpellCorrection)
	mux.HandleFunc("PUT /api/spell-corrections/{id}", s.updateSpellCorrection)
	mux.HandleFunc("DELETE /api/spell-corrections/{id}", s.deleteSpellCorrection)
	mux.HandleFunc("GET /api/correct", s.correct)
}

type spellCorrectionRequest struct {
	Word    string `json:"word"`
	Correct string `json:"correct"`
	Status  string `json:"status"`
}

func (s *Server) createSpellCorrection(w http.ResponseWriter, r *http.Request) {
	var req spellCorrectionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	v, err := s.svc.CreateSpellCorrection(model.SpellCorrection{Word: req.Word, Correct: req.Correct, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, v)
}

func (s *Server) listSpellCorrections(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	items, total, err := s.svc.ListSpellCorrections(pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSpellCorrection(w http.ResponseWriter, r *http.Request) {
	v, err := s.svc.GetSpellCorrection(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, v)
}

func (s *Server) updateSpellCorrection(w http.ResponseWriter, r *http.Request) {
	var req spellCorrectionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	v, err := s.svc.UpdateSpellCorrection(r.PathValue("id"), model.SpellCorrection{Word: req.Word, Correct: req.Correct, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, v)
}

func (s *Server) deleteSpellCorrection(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteSpellCorrection(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) correct(w http.ResponseWriter, r *http.Request) {
	word := r.URL.Query().Get("word")
	corrected, ok := s.svc.Correct(word)
	httpx.OK(w, map[string]interface{}{"word": word, "corrected": corrected, "changed": ok})
}
