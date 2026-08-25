package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"schoolbusauth/internal/query"
	"schoolbusauth/internal/service"
)

type Server struct {
	service *service.Service
	query   *query.Service
	mux     *http.ServeMux
}

func New(serviceLayer *service.Service) *Server {
	server := &Server{service: serviceLayer, query: query.New(serviceLayer.Store()), mux: http.NewServeMux()}
	server.routes()
	return server
}

func (s *Server) Handler() http.Handler { return securityHeaders(s.mux) }

func (s *Server) routes() {
	s.mux.HandleFunc("GET /", s.home)
	s.mux.HandleFunc("POST /authorizations", s.create)
	s.mux.HandleFunc("GET /authorizations/{id}", s.detail)
	s.mux.HandleFunc("POST /authorizations/{id}/issue", s.issue)
	s.mux.HandleFunc("POST /authorizations/{id}/revoke", s.revoke)
	s.mux.HandleFunc("GET /authorizations/{id}/preview", s.preview)
	s.mux.HandleFunc("GET /authorizations/{id}/download", s.download)
	s.mux.HandleFunc("POST /authorizations/{id}/print", s.recordPrint)
	s.mux.HandleFunc("GET /reports/daily", s.dailyReport)
	s.mux.HandleFunc("GET /reports/export.csv", s.exportCSV)
	s.mux.HandleFunc("GET /operations/dashboard", s.dashboard)
	s.mux.HandleFunc("GET /operations/handoff", s.handoff)
	s.mux.HandleFunc("GET /authorizations/{id}/compliance", s.compliance)
	s.mux.HandleFunc("GET /health", s.health)
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func requireMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return false
	}
	return true
}
