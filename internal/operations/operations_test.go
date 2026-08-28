package operations

import (
	"schoolbusauth/internal/domain"
	"testing"
	"time"
)

func operationBundle() domain.AuthorizationBundle {
	return domain.AuthorizationBundle{Authorization: domain.Authorization{ID: "a1", StudentID: "s1", RouteID: "r1", GuardianID: "g1", AuthorizedDate: "2026-09-01", TeacherName: "王老师", Status: domain.StatusIssued, Revision: 2, IssuedAt: ptrTime(time.Date(2026, 9, 1, 8, 0, 0, 0, time.UTC))}, Student: domain.Student{ID: "s1", Name: "李明", ClassName: "三班", SchoolNumber: "S1"}, Route: domain.Route{ID: "r1", Name: "东线", BusNumber: "B1", PickupPoint: "南门", GuardPost: "东门"}, Guardian: domain.Guardian{ID: "g1", Name: "李华", Relationship: "父亲", DocumentLastFour: "4832"}}
}

func ptrTime(value time.Time) *time.Time { return &value }

func TestOperationalViews(t *testing.T) {
	now := time.Date(2026, 9, 1, 16, 0, 0, 0, time.UTC)
	bundle := operationBundle()
	dashboard := BuildDashboard([]domain.AuthorizationBundle{bundle}, now)
	if dashboard.Total != 1 || dashboard.Usable != 1 {
		t.Fatal("dashboard totals wrong")
	}
	handoff := BuildHandoffSheet([]domain.AuthorizationBundle{bundle}, "东门", now)
	if len(handoff.Rows) != 1 || handoff.PlainText() == "" {
		t.Fatal("handoff missing")
	}
	schedule := BuildGateSchedule([]domain.AuthorizationBundle{bundle}, []GateShift{{GuardPost: "东门", Start: "15:00", End: "18:00", GuardName: "赵师傅"}}, now)
	if err := schedule.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestGateReconciliationAndCompliance(t *testing.T) {
	now := time.Date(2026, 9, 1, 16, 0, 0, 0, time.UTC)
	bundle := operationBundle()
	entry := GateEntry{AuthorizationID: "a1", StudentName: "李明", GuardianName: "李华", DocumentSuffix: "4832", GuardPost: "东门", ReleasedAt: now, GuardName: "赵师傅"}
	report := ReconcileGateEntries([]domain.AuthorizationBundle{bundle}, []GateEntry{entry}, now)
	if report.Exceptions != 0 || report.Error() != nil {
		t.Fatal("valid gate entry did not reconcile")
	}
	audits := []domain.AuditEvent{{Action: "created", OccurredAt: now}, {Action: "issued", OccurredAt: now}}
	compliance := ReviewCompliance(bundle, audits, nil, now)
	if !compliance.Consistent {
		t.Fatalf("compliance findings: %#v", compliance.Findings)
	}
}
