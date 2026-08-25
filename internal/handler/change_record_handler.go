package handler

import (
	"net/http"

	"datasync/internal/model"
	"datasync/pkg/httpx"
)

func (s *Server) registerChangeRecordRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/records", s.createChangeRecord)
	mux.HandleFunc("GET /api/records", s.listChangeRecords)
	mux.HandleFunc("GET /api/records/{id}", s.getChangeRecord)
}

type createChangeRecordRequest struct {
	TaskID    string `json:"task_id"`
	TableName string `json:"table_name"`
	Operation string `json:"operation"`
	Payload   string `json:"payload"`
	BinlogPos string `json:"binlog_pos"`
}

func (s *Server) createChangeRecord(w http.ResponseWriter, r *http.Request) {
	var req createChangeRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.CreateChangeRecord(model.ChangeRecord{
		TaskID:    req.TaskID,
		TableName: req.TableName,
		Operation: req.Operation,
		Payload:   req.Payload,
		BinlogPos: req.BinlogPos,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) listChangeRecords(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ChangeRecordFilter{
		TaskID:    r.URL.Query().Get("task_id"),
		TableName: r.URL.Query().Get("table_name"),
		Operation: r.URL.Query().Get("operation"),
	}
	items, total, err := s.svc.ListChangeRecords(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getChangeRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := s.svc.GetChangeRecord(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}
