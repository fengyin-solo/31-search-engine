package handler

import (
	"net/http"

	"searchengine/internal/model"
	"searchengine/pkg/httpx"
)

func (s *Server) registerBoostRuleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/boost-rules", s.createBoostRule)
	mux.HandleFunc("GET /api/boost-rules", s.listBoostRules)
	mux.HandleFunc("GET /api/boost-rules/{id}", s.getBoostRule)
	mux.HandleFunc("PUT /api/boost-rules/{id}", s.updateBoostRule)
	mux.HandleFunc("DELETE /api/boost-rules/{id}", s.deleteBoostRule)
}

type boostRuleRequest struct {
	Field  string  `json:"field"`
	Term   string  `json:"term"`
	Boost  float64 `json:"boost"`
	Status string  `json:"status"`
}

func (s *Server) createBoostRule(w http.ResponseWriter, r *http.Request) {
	var req boostRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.CreateBoostRule(model.BoostRule{Field: req.Field, Term: req.Term, Boost: req.Boost, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, b)
}

func (s *Server) listBoostRules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	items, total, err := s.svc.ListBoostRules(pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getBoostRule(w http.ResponseWriter, r *http.Request) {
	b, err := s.svc.GetBoostRule(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

func (s *Server) updateBoostRule(w http.ResponseWriter, r *http.Request) {
	var req boostRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.UpdateBoostRule(r.PathValue("id"), model.BoostRule{Field: req.Field, Term: req.Term, Boost: req.Boost, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

func (s *Server) deleteBoostRule(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteBoostRule(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
