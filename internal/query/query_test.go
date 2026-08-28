package query

import (
	"bytes"
	"path/filepath"
	"schoolbusauth/internal/domain"
	"schoolbusauth/internal/store"
	"testing"
	"time"
)

func TestSearchSummaryExportAndCalendar(t *testing.T) {
	repository, err := store.Open(filepath.Join(t.TempDir(), "query.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer repository.Close()
	bundle := domain.AuthorizationBundle{Authorization: domain.Authorization{ID: "a1", StudentID: "s1", RouteID: "r1", GuardianID: "g1", AuthorizedDate: "2026-09-01", TeacherName: "王老师", Status: domain.StatusDraft, CreatedAt: time.Unix(1, 0), Revision: 1}, Student: domain.Student{ID: "s1", Name: "李明", ClassName: "三班", SchoolNumber: "S1"}, Route: domain.Route{ID: "r1", Name: "东线", BusNumber: "B1", PickupPoint: "南门"}, Guardian: domain.Guardian{ID: "g1", Name: "李华", Relationship: "父亲", DocumentLastFour: "4832"}}
	if err := repository.CreateBundle(bundle, domain.AuditEvent{ID: "e1", AuthorizationID: "a1", Action: "created", OccurredAt: time.Unix(1, 0)}); err != nil {
		t.Fatal(err)
	}
	service := New(repository)
	values, err := service.Search(domain.SearchFilter{StudentName: "李"})
	if err != nil || len(values) != 1 {
		t.Fatal("search failed")
	}
	summary, err := service.Summary("2026-09-01")
	if err != nil || summary.Drafts != 1 {
		t.Fatal("summary failed")
	}
	var csv bytes.Buffer
	if err := service.ExportCSV(&csv, domain.SearchFilter{}); err != nil || csv.Len() == 0 {
		t.Fatal("export failed")
	}
	calendar, err := service.Calendar(time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 9, 3, 0, 0, 0, 0, time.UTC))
	if err != nil || len(calendar.Days) != 3 || calendar.Days[0].Total != 1 {
		t.Fatal("calendar failed")
	}
}
