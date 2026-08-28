package schoolbusauth

import (
	"path/filepath"
	"schoolbusauth/internal/domain"
	"schoolbusauth/internal/query"
	"strings"
	"testing"
	"time"
)

func TestWorkflowGatePrint(t *testing.T) {
	now := time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC)
	called := false
	application, err := Open(Config{DatabasePath: filepath.Join(t.TempDir(), "print.db"), BusinessDate: now, OnComplete: func(domain.AuthorizationBundle) error { called = true; return nil }})
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
	content, filename, err := application.Service.DownloadHTML(issued.Authorization.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(filename, ".html") || !strings.Contains(string(content), "仅限当日使用") {
		t.Fatal("invalid download")
	}
	if _, err := application.Service.RecordPrint(issued.Authorization.ID, "门卫赵师傅", 2); err != nil {
		t.Fatal(err)
	}
	detail, err := query.New(application.Store).Detail(issued.Authorization.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(detail.PrintRecords) != 1 || detail.PrintRecords[0].Copies != 2 {
		t.Fatal("print record missing")
	}
	if err := application.Service.Complete(issued.Authorization.ID); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("completion callback not called")
	}
}
