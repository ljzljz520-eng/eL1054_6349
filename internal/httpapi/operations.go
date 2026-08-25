package httpapi

import (
	"net/http"
	"schoolbusauth/internal/domain"
	"schoolbusauth/internal/operations"
	"time"
)

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	values, err := s.query.Search(domain.SearchFilter{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	date := r.URL.Query().Get("date")
	now, err := time.Parse("2006-01-02", date)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, operations.BuildDashboard(values, now))
}

func (s *Server) handoff(w http.ResponseWriter, r *http.Request) {
	date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	values, err := s.query.Search(domain.SearchFilter{Date: date.Format("2006-01-02")})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	sheet := operations.BuildHandoffSheet(values, r.URL.Query().Get("guard_post"), date)
	if r.URL.Query().Get("format") == "text" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(sheet.PlainText()))
		return
	}
	writeJSON(w, http.StatusOK, sheet)
}

func (s *Server) compliance(w http.ResponseWriter, r *http.Request) {
	detail, err := s.query.Detail(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	report := operations.ReviewCompliance(detail.Bundle, detail.AuditEvents, detail.PrintRecords, s.service.Now())
	writeJSON(w, http.StatusOK, report)
}
