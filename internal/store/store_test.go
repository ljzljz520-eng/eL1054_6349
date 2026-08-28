package store

import (
	"path/filepath"
	"schoolbusauth/internal/domain"
	"testing"
	"time"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "reopen.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	bundle := domain.AuthorizationBundle{
		Authorization: domain.Authorization{ID: "a1", StudentID: "s1", RouteID: "r1", GuardianID: "g1", AuthorizedDate: "2026-09-01", TeacherName: "王老师", Status: domain.StatusDraft, CreatedAt: time.Unix(1, 0), Revision: 1},
		Student:       domain.Student{ID: "s1", Name: "李明", ClassName: "三年级二班", SchoolNumber: "S1"},
		Route:         domain.Route{ID: "r1", Name: "东线", BusNumber: "B1", PickupPoint: "南门"},
		Guardian:      domain.Guardian{ID: "g1", Name: "李华", Relationship: "父亲", DocumentLastFour: "4832"},
	}
	audit := domain.AuditEvent{ID: "e1", AuthorizationID: "a1", Action: "created", Actor: "王老师", OccurredAt: time.Unix(1, 0)}
	if err := s.CreateBundle(bundle, audit); err != nil {
		t.Fatal(err)
	}
	if err := s.SavePrintRecord(domain.PrintRecord{ID: "p1", AuthorizationID: "a1", Printer: "赵师傅", Copies: 1, PrintedAt: time.Unix(2, 0)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	s, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	loaded, err := s.Bundle("a1")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Student.Name != "李明" || loaded.Guardian.DocumentLastFour != "4832" {
		t.Fatal("data changed after reopen")
	}
	prints, err := s.PrintRecords("a1")
	if err != nil || len(prints) != 1 {
		t.Fatal("print record missing after reopen")
	}
}
