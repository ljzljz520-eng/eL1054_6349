package httpapi

import (
	"net/http"
	"schoolbusauth/internal/domain"
)

func (s *Server) dailyReport(w http.ResponseWriter, r *http.Request) {
	summary, err := s.query.Summary(r.URL.Query().Get("date"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

func (s *Server) exportCSV(w http.ResponseWriter, r *http.Request) {
	filter := domain.SearchFilter{Date: r.URL.Query().Get("date"), StudentName: r.URL.Query().Get("student"), RouteName: r.URL.Query().Get("route"), Status: domain.AuthorizationStatus(r.URL.Query().Get("status"))}
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="authorizations.csv"`)
	w.Write([]byte{0xef, 0xbb, 0xbf})
	if err := s.query.ExportCSV(w, filter); err != nil {
		return
	}
}
