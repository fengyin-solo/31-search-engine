package handler

import (
	"net/http"

	"searchengine/internal/model"
	"searchengine/pkg/httpx"
)

func (s *Server) registerDocumentRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/documents", s.createDocument)
	mux.HandleFunc("GET /api/documents", s.listDocuments)
	mux.HandleFunc("GET /api/documents/{id}", s.getDocument)
	mux.HandleFunc("PUT /api/documents/{id}", s.updateDocument)
	mux.HandleFunc("DELETE /api/documents/{id}", s.deleteDocument)
}

type documentRequest struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	URL    string `json:"url"`
	Source string `json:"source"`
}

func (s *Server) createDocument(w http.ResponseWriter, r *http.Request) {
	var req documentRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.CreateDocument(model.Document{
		Title: req.Title, Body: req.Body, URL: req.URL, Source: req.Source,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, d)
}

func (s *Server) listDocuments(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.DocumentFilter{
		Status:  r.URL.Query().Get("status"),
		Source:  r.URL.Query().Get("source"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListDocuments(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getDocument(w http.ResponseWriter, r *http.Request) {
	d, err := s.svc.GetDocument(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

func (s *Server) updateDocument(w http.ResponseWriter, r *http.Request) {
	var req documentRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.UpdateDocument(r.PathValue("id"), model.Document{
		Title: req.Title, Body: req.Body, URL: req.URL, Source: req.Source,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

func (s *Server) deleteDocument(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteDocument(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
