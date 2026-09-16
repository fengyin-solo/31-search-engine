package handler

import (
	"net/http"
	"strconv"

	"searchengine/pkg/httpx"
)

func (s *Server) registerSearchRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/search", s.search)
	mux.HandleFunc("GET /api/search/facets", s.searchFaceted)
	mux.HandleFunc("POST /api/index-documents", s.indexDocument)
}

func (s *Server) search(w http.ResponseWriter, r *http.Request) {
	indexID := r.URL.Query().Get("index_id")
	query := r.URL.Query().Get("q")
	topK := 10
	if v := r.URL.Query().Get("top_k"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			topK = n
		}
	}
	results, total, err := s.svc.Search(indexID, query, topK)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{
		"items": results,
		"total": total,
		"query": query,
	})
}

type indexDocumentRequest struct {
	IndexID string `json:"index_id"`
	DocID   string `json:"doc_id"`
}

func (s *Server) indexDocument(w http.ResponseWriter, r *http.Request) {
	var req indexDocumentRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if err := s.svc.IndexDocument(req.IndexID, req.DocID); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]bool{"indexed": true})
}

func (s *Server) searchFaceted(w http.ResponseWriter, r *http.Request) {
	indexID := r.URL.Query().Get("index_id")
	query := r.URL.Query().Get("q")
	facets, total, err := s.svc.SearchFaceted(indexID, query)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"facets": facets, "total": total})
}
