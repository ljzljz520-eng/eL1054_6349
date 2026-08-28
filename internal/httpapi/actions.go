package httpapi

import (
	"net/http"
	"schoolbusauth/internal/domain"
	"strconv"
)

func (s *Server) create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	req := domain.CreateRequest{
		Student:        domain.Student{Name: r.FormValue("student_name"), ClassName: r.FormValue("class_name"), SchoolNumber: r.FormValue("school_number")},
		Route:          domain.Route{Name: r.FormValue("route_name"), BusNumber: r.FormValue("bus_number"), PickupPoint: r.FormValue("pickup_point"), GuardPost: r.FormValue("guard_post")},
		Guardian:       domain.Guardian{Name: r.FormValue("guardian_name"), Relationship: r.FormValue("relationship"), DocumentLastFour: r.FormValue("document_last_four")},
		AuthorizedDate: r.FormValue("authorized_date"), TeacherName: r.FormValue("teacher_name"), Reason: r.FormValue("reason"),
	}
	bundle, err := s.service.CreateDraft(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	http.Redirect(w, r, "/authorizations/"+bundle.Authorization.ID, http.StatusSeeOther)
}

func (s *Server) issue(w http.ResponseWriter, r *http.Request) {
	bundle, err := s.service.Issue(r.PathValue("id"), r.FormValue("actor"))
	if err != nil {
		writeError(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusOK, bundle)
}

func (s *Server) revoke(w http.ResponseWriter, r *http.Request) {
	bundle, err := s.service.Revoke(r.PathValue("id"), r.FormValue("actor"), r.FormValue("reason"))
	if err != nil {
		writeError(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusOK, bundle)
}

func (s *Server) preview(w http.ResponseWriter, r *http.Request) {
	preview, err := s.service.Preview(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(preview.HTML)
}

func (s *Server) download(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("format")
	var content []byte
	var filename string
	var err error
	if format == "text" {
		content, filename, err = s.service.DownloadText(r.PathValue("id"))
	} else {
		content, filename, err = s.service.DownloadHTML(r.PathValue("id"))
	}
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(content)
}

func (s *Server) recordPrint(w http.ResponseWriter, r *http.Request) {
	copies, err := strconv.Atoi(r.FormValue("copies"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	record, err := s.service.RecordPrint(r.PathValue("id"), r.FormValue("printer"), copies)
	if err != nil {
		writeError(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusCreated, record)
}
