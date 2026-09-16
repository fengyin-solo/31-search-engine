package handler

import (
	"net/http"

	"searchengine/internal/model"
	"searchengine/pkg/httpx"
)

func (s *Server) registerIndexRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/indexes", s.createIndex)
	mux.HandleFunc("GET /api/indexes", s.listIndexes)
	mux.HandleFunc("GET /api/indexes/{id}", s.getIndex)
	mux.HandleFunc("PUT /api/indexes/{id}", s.updateIndex)
	mux.HandleFunc("DELETE /api/indexes/{id}", s.deleteIndex)
	mux.HandleFunc("POST /api/indexes/{id}/activate", s.activateIndex)
	mux.HandleFunc("POST /api/indexes/{id}/deactivate", s.deactivateIndex)
	mux.HandleFunc("POST /api/indexes/{id}/rebuild", s.rebuildIndex)
}

type indexRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Fields      []string `json:"fields"`
	AnalyzerID  string   `json:"analyzer_id"`
}

func (s *Server) createIndex(w http.ResponseWriter, r *http.Request) {
	var req indexRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	idx, err := s.svc.CreateIndex(model.Index{
		Name: req.Name, Description: req.Description, Fields: req.Fields, AnalyzerID: req.AnalyzerID,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, idx)
}

func (s *Server) listIndexes(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	items, total, err := s.svc.ListIndexes(pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getIndex(w http.ResponseWriter, r *http.Request) {
	idx, err := s.svc.GetIndex(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, idx)
}

func (s *Server) updateIndex(w http.ResponseWriter, r *http.Request) {
	var req indexRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	idx, err := s.svc.UpdateIndex(r.PathValue("id"), model.Index{
		Name: req.Name, Description: req.Description, Fields: req.Fields, AnalyzerID: req.AnalyzerID,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, idx)
}

func (s *Server) deleteIndex(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteIndex(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) activateIndex(w http.ResponseWriter, r *http.Request) {
	idx, err := s.svc.ActivateIndex(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, idx)
}

func (s *Server) deactivateIndex(w http.ResponseWriter, r *http.Request) {
	idx, err := s.svc.DeactivateIndex(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, idx)
}

func (s *Server) rebuildIndex(w http.ResponseWriter, r *http.Request) {
	count, err := s.svc.RebuildIndex(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]int{"reindexed": count})
}
