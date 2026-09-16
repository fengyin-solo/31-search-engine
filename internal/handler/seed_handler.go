package handler

import (
	"net/http"

	"searchengine/pkg/httpx"
)

func (s *Server) registerSeedRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/seed/demo", s.seedDemo)
}

func (s *Server) seedDemo(w http.ResponseWriter, r *http.Request) {
	res, err := s.svc.SeedDemoData()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, res)
}
