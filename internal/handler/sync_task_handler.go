package handler

import (
	"net/http"

	"datasync/internal/model"
	"datasync/pkg/httpx"
)

func (s *Server) registerSyncTaskRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/tasks", s.createSyncTask)
	mux.HandleFunc("GET /api/tasks", s.listSyncTasks)
	mux.HandleFunc("GET /api/tasks/{id}", s.getSyncTask)
	mux.HandleFunc("PUT /api/tasks/{id}", s.updateSyncTask)
	mux.HandleFunc("DELETE /api/tasks/{id}", s.deleteSyncTask)
	mux.HandleFunc("POST /api/tasks/{id}/start", s.startSyncTask)
	mux.HandleFunc("POST /api/tasks/{id}/pause", s.pauseSyncTask)
	mux.HandleFunc("POST /api/tasks/{id}/fail", s.failSyncTask)
}

type createSyncTaskRequest struct {
	Name        string `json:"name"`
	SourceID    string `json:"source_id"`
	TargetID    string `json:"target_id"`
	Mode        string `json:"mode"`
	IntervalSec int    `json:"interval_sec"`
}

func (s *Server) createSyncTask(w http.ResponseWriter, r *http.Request) {
	var req createSyncTaskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.CreateSyncTask(model.SyncTask{
		Name:        req.Name,
		SourceID:    req.SourceID,
		TargetID:    req.TargetID,
		Mode:        req.Mode,
		IntervalSec: req.IntervalSec,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, t)
}

func (s *Server) listSyncTasks(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.SyncTaskFilter{
		Mode:    r.URL.Query().Get("mode"),
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListSyncTasks(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSyncTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := s.svc.GetSyncTask(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

type updateSyncTaskRequest struct {
	Name        string `json:"name"`
	SourceID    string `json:"source_id"`
	TargetID    string `json:"target_id"`
	Mode        string `json:"mode"`
	IntervalSec int    `json:"interval_sec"`
	Status      string `json:"status"`
}

func (s *Server) updateSyncTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateSyncTaskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.UpdateSyncTask(id, model.SyncTask{
		Name:        req.Name,
		SourceID:    req.SourceID,
		TargetID:    req.TargetID,
		Mode:        req.Mode,
		IntervalSec: req.IntervalSec,
		Status:      req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) deleteSyncTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteSyncTask(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) startSyncTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := s.svc.StartSyncTask(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) pauseSyncTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := s.svc.PauseSyncTask(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) failSyncTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := s.svc.FailSyncTask(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}
