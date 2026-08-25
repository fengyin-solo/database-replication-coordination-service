package handler

import (
	"net/http"

	"datasync/internal/model"
	"datasync/pkg/httpx"
)

func (s *Server) registerCheckpointRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/checkpoints", s.createCheckpoint)
	mux.HandleFunc("GET /api/checkpoints", s.listCheckpoints)
	mux.HandleFunc("GET /api/checkpoints/{task_id}", s.getCheckpointByTaskID)
	mux.HandleFunc("PUT /api/checkpoints/{task_id}", s.advanceCheckpoint)
}

type createCheckpointRequest struct {
	TaskID        string `json:"task_id"`
	Position      string `json:"position"`
	RowsProcessed int64  `json:"rows_processed"`
}

func (s *Server) createCheckpoint(w http.ResponseWriter, r *http.Request) {
	var req createCheckpointRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.CreateCheckpoint(model.Checkpoint{
		TaskID:        req.TaskID,
		Position:      req.Position,
		RowsProcessed: req.RowsProcessed,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) listCheckpoints(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.CheckpointFilter{
		TaskID: r.URL.Query().Get("task_id"),
	}
	items, total, err := s.svc.ListCheckpoints(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getCheckpointByTaskID(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("task_id")
	c, err := s.svc.GetCheckpointByTaskID(taskID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

type advanceCheckpointRequest struct {
	Position  string `json:"position"`
	RowsDelta int64  `json:"rows_delta"`
}

func (s *Server) advanceCheckpoint(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("task_id")
	var req advanceCheckpointRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.AdvanceCheckpoint(taskID, req.Position, req.RowsDelta)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) registerRunRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/run", s.runSyncTask)
	mux.HandleFunc("POST /api/record-change", s.recordChange)
}

type runSyncTaskRequest struct {
	TaskID string `json:"task_id"`
}

func (s *Server) runSyncTask(w http.ResponseWriter, r *http.Request) {
	var req runSyncTaskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.RunSyncTask(req.TaskID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

type recordChangeRequest struct {
	TaskID    string `json:"task_id"`
	TableName string `json:"table_name"`
	Operation string `json:"operation"`
	Payload   string `json:"payload"`
	BinlogPos string `json:"binlog_pos"`
}

func (s *Server) recordChange(w http.ResponseWriter, r *http.Request) {
	var req recordChangeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rec, err := s.svc.RecordChangeAndAdvance(req.TaskID, req.TableName, req.Operation, req.Payload, req.BinlogPos)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rec)
}
