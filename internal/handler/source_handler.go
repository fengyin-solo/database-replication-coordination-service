package handler

import (
	"net/http"

	"datasync/internal/model"
	"datasync/pkg/httpx"
)

func (s *Server) registerSourceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/sources", s.createSource)
	mux.HandleFunc("GET /api/sources", s.listSources)
	mux.HandleFunc("GET /api/sources/{id}", s.getSource)
	mux.HandleFunc("PUT /api/sources/{id}", s.updateSource)
	mux.HandleFunc("DELETE /api/sources/{id}", s.deleteSource)
}

type createSourceRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Status   string `json:"status"`
}

func (s *Server) createSource(w http.ResponseWriter, r *http.Request) {
	var req createSourceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	src, err := s.svc.CreateSource(model.Source{
		Name:     req.Name,
		Type:     req.Type,
		Host:     req.Host,
		Port:     req.Port,
		Database: req.Database,
		Status:   req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, src)
}

func (s *Server) listSources(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.SourceFilter{
		Type:    r.URL.Query().Get("type"),
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListSources(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSource(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	src, err := s.svc.GetSource(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, src)
}

type updateSourceRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Status   string `json:"status"`
}

func (s *Server) updateSource(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateSourceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	src, err := s.svc.UpdateSource(id, model.Source{
		Name:     req.Name,
		Type:     req.Type,
		Host:     req.Host,
		Port:     req.Port,
		Database: req.Database,
		Status:   req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, src)
}

func (s *Server) deleteSource(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteSource(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
