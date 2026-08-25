package handler

import (
	"net/http"
	"strconv"

	"datasync/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/overview", s.getOverviewStats)
	mux.HandleFunc("GET /api/stats/tasks-top", s.getTasksTop)
}

func (s *Server) getOverviewStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.GetOverviewStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) getTasksTop(w http.ResponseWriter, r *http.Request) {
	n, _ := strconv.Atoi(r.URL.Query().Get("n"))
	if n <= 0 {
		n = 10
	}
	items, err := s.svc.GetTasksTop(n)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}
