package schoolbusauth

import (
	"path/filepath"
	"schoolbusauth/internal/domain"
	"schoolbusauth/internal/query"
	"strings"
	"testing"
	"time"
)

func TestWorkflowRevokeAndAudit(t *testing.T) {
	now := time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC)
	application, err := Open(Config{DatabasePath: filepath.Join(t.TempDir(), "revoke.db"), BusinessDate: now, OnComplete: func(domain.AuthorizationBundle) error { return nil }})
	if err != nil {
		t.Fatal(err)
	}
	defer application.Close()
	draft, err := application.Service.CreateDraft(DemoRequest("2026-09-01"))
	if err != nil {
		t.Fatal(err)
	}
	issued, err := application.Service.Issue(draft.Authorization.ID, "王老师")
	if err != nil {
		t.Fatal(err)
	}
	revoked, err := application.Service.Revoke(issued.Authorization.ID, "王老师", "家长取消临时安排")
	if err != nil {
		t.Fatal(err)
	}
	preview, err := application.Service.Preview(revoked.Authorization.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(preview.HTML), "授权已撤销") {
		t.Fatal("revoked warning missing")
	}
	detail, err := query.New(application.Store).Detail(revoked.Authorization.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(detail.AuditEvents) != 3 {
		t.Fatalf("audit count %d", len(detail.AuditEvents))
	}
}
