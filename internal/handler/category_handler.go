package handler

import (
	"net/http"

	"searchengine/internal/model"
	"searchengine/pkg/httpx"
)

func (s *Server) registerCategoryRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/categories", s.createCategory)
	mux.HandleFunc("GET /api/categories", s.listCategories)
	mux.HandleFunc("GET /api/categories/{id}", s.getCategory)
	mux.HandleFunc("PUT /api/categories/{id}", s.updateCategory)
	mux.HandleFunc("DELETE /api/categories/{id}", s.deleteCategory)
	mux.HandleFunc("POST /api/documents/{id}/category", s.setDocumentCategory)
}

type categoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (s *Server) createCategory(w http.ResponseWriter, r *http.Request) {
	var req categoryRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.CreateCategory(model.Category{Name: req.Name, Description: req.Description, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) listCategories(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	items, total, err := s.svc.ListCategories(pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getCategory(w http.ResponseWriter, r *http.Request) {
	c, err := s.svc.GetCategory(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) updateCategory(w http.ResponseWriter, r *http.Request) {
	var req categoryRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.UpdateCategory(r.PathValue("id"), model.Category{Name: req.Name, Description: req.Description, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) deleteCategory(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteCategory(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type setCategoryRequest struct {
	CategoryID string `json:"category_id"`
}

func (s *Server) setDocumentCategory(w http.ResponseWriter, r *http.Request) {
	var req setCategoryRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.SetDocumentCategory(r.PathValue("id"), req.CategoryID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}
